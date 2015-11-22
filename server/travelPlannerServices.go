package server

import (
      "fmt"
      "net/http"
      "io/ioutil"
      "encoding/json"
      "strconv"
      "gopkg.in/mgo.v2/bson"
      "gopkg.in/mgo.v2"
      "bytes"
    )


type PriceResponse struct {

  Prices[] struct {
    Localized_display_name string
    Duration int
    Distance float64
    Low_estimate int
    Product_id string
  }
}

type EtaStruct struct {
   Eta int
}

var bestRouteArray[] string
var startLocation string
var status = "requesting"

func findShortestPath(tripRequest TripRequest) (TripResponse) {

  var tripResponse TripResponse
  startLocation := tripRequest.Starting_from_location_id
  var total_uber_costs int
  var total_uber_duration int
  var total_distance float64
  destinationArray := tripRequest.Location_ids
  resultArray := make([]string, len(destinationArray))
  index := -1

  for len(destinationArray) >= 1 {
    array,minimumCost,dur,dist,idWithMinEstimate := compareCost(startLocation, destinationArray)
    total_uber_costs += minimumCost
    total_uber_duration += dur
    total_distance += dist
    index = index +1
    resultArray[index] = idWithMinEstimate
    destinationArray = array
    startLocation = idWithMinEstimate
  }
  arrayForStart := make([]string, 1)
  arrayForStart[0] = tripRequest.Starting_from_location_id
  _,minimumCost,dur,dist,_ := compareCost(startLocation, arrayForStart)

   total_uber_costs += minimumCost
   total_uber_duration += dur
   total_distance += dist

   tripResponse.Starting_from_location_id = tripRequest.Starting_from_location_id
   tripResponse.Best_route_location_ids = resultArray
   tripResponse.Total_uber_costs = total_uber_costs
   tripResponse.Total_uber_duration = total_uber_duration
   tripResponse.Total_distance = total_distance
   response := SaveBestRouteInDB(&tripResponse)
   return response

}

func compareCost(startLocation string, destinationArray[] string) (array[] string,costReturned int,durReturned int,distReturned float64, idWithMinEstimate string ){
  var minimumCost int
  var dur int
  var dist float64
  var removeIndex int
  var idWithMinCost string

  for count:=0 ; count< len(destinationArray); count++ {
       minEstimate,duration,distance:= getDataFromUber(startLocation,destinationArray[count])
         if minimumCost == 0{
               minimumCost = minEstimate
               removeIndex = count
               dur = duration
               dist = distance
               removeIndex = count
               idWithMinCost = destinationArray[count]
         } else {
              if minimumCost > minEstimate {
                    minimumCost = minEstimate
                    dur = duration
                    dist = distance
                    removeIndex = count
                    idWithMinCost = destinationArray[count]
              }
            }
    }
    result := []string{}
    result = append(result, destinationArray[0:removeIndex]...)
    // Append part after the removed element.
    result = append(result, destinationArray[removeIndex+1:]...)
    return result,minimumCost,dur,dist,idWithMinCost

}

func SaveBestRouteInDB(tripResponse *TripResponse) (TripResponse){
  session, dbName := connectionWithDB()
  defer session.Close()
  session.SetSafe(&mgo.Safe{})
  c := session.DB(dbName).C("uber")

  data := TripResponse {
           Id: bson.NewObjectId(),
           Status: "planning",
           Starting_from_location_id: tripResponse.Starting_from_location_id,
           Best_route_location_ids: tripResponse.Best_route_location_ids,
           Total_uber_costs: tripResponse.Total_uber_costs,
           Total_uber_duration:tripResponse.Total_uber_duration,
           Total_distance: tripResponse.Total_distance,
  }
err := c.Insert(data)
if err != nil {
  panic(err)
}
return data
}

func getBestRouteFromDB(id bson.ObjectId) (TripResponse){
  session, dbName := connectionWithDB()
  defer session.Close()
  session.SetSafe(&mgo.Safe{})
    c := session.DB(dbName).C("uber")
  var resultFromDb TripResponse
  err := c.Find(bson.M{"id": id}).One(&resultFromDb)
   if err != nil {
       panic(err)
   }

  return resultFromDb
}


func FetchUberPrice(start LocationResponse, end LocationResponse) (PriceResponse){
  var priceResponse PriceResponse

   var url="https://sandbox-api.uber.com/v1/estimates/price?"
   startLat := strconv.FormatFloat(start.Coordinate.Lat, 'G', -1, 64)
   startLng := strconv.FormatFloat(start.Coordinate.Lng, 'G', -1, 64)
   endLat := strconv.FormatFloat(end.Coordinate.Lat, 'G', -1, 64)
   endLng := strconv.FormatFloat(end.Coordinate.Lng, 'G', -1, 64)

   url = url+"start_latitude="+startLat+"&start_longitude="+startLng+
   "&end_latitude="+endLat+"&end_longitude="+endLng+"&server_token=8CsictSgQ-7dGckffNxRJfYlEIouRJTOGb7SLBXB"
   resp, err := http.Get(url)
   if err!=nil {
     fmt.Println("Error from Uber API: ",err)
   } else {
   defer resp.Body.Close()
          jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
          if err != nil {
                  panic(err)
          }
          err = json.Unmarshal([]byte(jsonDataFromHttp), &priceResponse)

          if err != nil {
            fmt.Println("error: ",err)
                  panic(err)
          }
    }
        return priceResponse
}

func getDataFromUber(startLocation string, endLocation string)  (estimate int,dur int,dist float64){

  idValueStartLocation := convertToObjectId(startLocation)
  responseStartLocation := getFromDB(idValueStartLocation)
  idValueEndLocation := convertToObjectId(endLocation)
  responseEndLocation := getFromDB(idValueEndLocation)
  priceResponse := FetchUberPrice(responseStartLocation,responseEndLocation)
  minEstimate :=  priceResponse.Prices[0].Low_estimate
  duration := priceResponse.Prices[0].Duration
  distance := priceResponse.Prices[0].Distance

  return minEstimate,duration,distance
}

func putEtaData(tripResponse TripResponse) (TripPutResponse){
  var tripPutResponse TripPutResponse
  var endResponse LocationResponse
  tripPutResponse.Id = tripResponse.Id
  tripPutResponse.Starting_from_location_id = tripResponse.Starting_from_location_id
  tripPutResponse.Best_route_location_ids = tripResponse.Best_route_location_ids
  tripPutResponse.Total_uber_costs = tripResponse.Total_uber_costs
  tripPutResponse.Total_uber_duration = tripResponse.Total_uber_duration
  tripPutResponse.Total_distance = tripResponse.Total_distance

  if len(bestRouteArray) == 0 && status == "requesting" {
    bestRouteArray = tripResponse.Best_route_location_ids
    startLocation = tripResponse.Starting_from_location_id
    tripPutResponse.Status = status
    tripPutResponse.Next_destination_location_id = bestRouteArray[0]
    startResponse := getFromDB(convertToObjectId(startLocation))
    endResponse := getFromDB(convertToObjectId(tripPutResponse.Next_destination_location_id))

    fetchEtaFromUber(startResponse,endResponse,&tripPutResponse)
    data := bestRouteArray
    bestRouteArray = []string{}
    bestRouteArrayAfterDelition := append(bestRouteArray, data[1:]...)
    bestRouteArrayAfterDelition = append(bestRouteArrayAfterDelition,startLocation)

    bestRouteArray = bestRouteArrayAfterDelition
    startLocation = tripPutResponse.Next_destination_location_id
  } else {
    startResponse := getFromDB(convertToObjectId(startLocation))
    if len(bestRouteArray) !=0 {
      endResponse = getFromDB(convertToObjectId(bestRouteArray[0]))
      tripPutResponse.Next_destination_location_id = bestRouteArray[0]
      fetchEtaFromUber(startResponse,endResponse,&tripPutResponse)
      startLocation = bestRouteArray[0]
      data := bestRouteArray
      bestRouteArray = []string{}
      bestRouteArrayAfterDelition := append(bestRouteArray, data[1:]...)

      bestRouteArray = bestRouteArrayAfterDelition
    } else {
      tripPutResponse.Next_destination_location_id = ""
    }
    tripPutResponse.Status = status



    if len(bestRouteArray) ==  0 {
        status = "finished"
      }
  }
  return tripPutResponse
}

func fetchEtaFromUber(startResponse LocationResponse, endResponse LocationResponse, tripPutResponse *TripPutResponse) {

  var eta EtaStruct
  startLat := strconv.FormatFloat(startResponse.Coordinate.Lat, 'G', -1, 64)
  startLng := strconv.FormatFloat(startResponse.Coordinate.Lng, 'G', -1, 64)
  endLat := strconv.FormatFloat(endResponse.Coordinate.Lat, 'G', -1, 64)
  endLng := strconv.FormatFloat(endResponse.Coordinate.Lng, 'G', -1, 64)

  url := "https://sandbox-api.uber.com/v1/requests"
  priceResponse := FetchUberPrice(startResponse,endResponse)

  s := map[string]string{
    "start_latitude":  startLat,
    "start_longitude": startLng,
    "end_latitude": endLat,
    "end_longitude": endLng,
    "product_id": priceResponse.Prices[0].Product_id,
  }
    jsonStr, err := json.Marshal(s)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiIxMDEwNTVhNi0wMWE0LTRkMDYtOGJkOC1kYmE2MGI1YTVjNjkiLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6IjJkN2ViOTdmLWM2ODMtNDg1Zi1hOGFjLTk4ODIwYjYwOGY0NSIsImV4cCI6MTQ1MDY1Mjg3NywiaWF0IjoxNDQ4MDYwODc3LCJ1YWN0IjoiUUVDazdpQVRZV2N0RTJTSDVWT3Nqc044T2xuVVpaIiwibmJmIjoxNDQ4MDYwNzg3LCJhdWQiOiJCUFIwYjViMno5aUdCNk9nWm5LQmJCcFNxY3c5Y2xRZSJ9.MiabMKuCaMgYM9m9dPxBeKDLiCmzN8anZjQorRz4W6msW-sieCb3MRwXZRMbu1U1ZqqJvckc97eIgj-QCmObBZaRwhQblSlBEg7UeFLt2WRKMegIf2UrsdZYdEdD2-_YBVRsR30IGXus2gY3nXzhjJ4607bnhU6TU42dyGCRFGbRaUJ-SlDw2zYW_9V1z0TPMHnVt9o_uW8dndDqjSdSF21IDJewzYk1YM-z6s3s7qUDGhMYuIwJrI_U7TcgbS-ommm3bC7YwhJ-ec-OuX3f2XtpA14MVWuUVhm7_quQ_2hshH6WWGCoLpVlapJ4pElZwq29XHSLGQdwLIXbOUmGSQ")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    err = json.Unmarshal([]byte(body), &eta)
    tripPutResponse.Uber_wait_time_eta = eta.Eta
}
