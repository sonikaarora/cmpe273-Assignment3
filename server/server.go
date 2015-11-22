package server

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/julienschmidt/httprouter"
    "gopkg.in/mgo.v2/bson"
)

type LocationRequest struct {
    Name string
    Address string
    City string
    State string
    Zip string
}

type Coordinates struct {
  Lat float64
  Lng float64
}

type LocationResponse struct {
  Id bson.ObjectId
  Name string
  Address string
  City string
  State string
  Zip string
  Coordinate Coordinates
}

type LocationCoordinates struct {
  Results[]struct {
    Geometry struct {
      Location struct {
        Lat float64
        Lng float64
      }
    }
  }
}

type TripRequest struct {
  Starting_from_location_id  string
  Location_ids[]  string
}

type TripResponse struct {
  Id bson.ObjectId
  Status string
  Starting_from_location_id string
  Best_route_location_ids[] string
  Total_uber_costs int
  Total_uber_duration int
  Total_distance float64
}

type TripPutResponse struct {
  Id bson.ObjectId
  Status string
  Starting_from_location_id string
  Next_destination_location_id string
  Best_route_location_ids[] string
  Total_uber_costs int
  Total_uber_duration int
  Total_distance float64
  Uber_wait_time_eta int
}

func updateLocationHandler(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {
  var locationRequest LocationRequest
  var locationCoordinates LocationCoordinates
  idValue := convertToObjectId(p.ByName("location_id"))
  jsonDataFromHttp, err := ioutil.ReadAll(request.Body)
  if err != nil {
          panic(err)
  }
  err = json.Unmarshal([]byte(jsonDataFromHttp), &locationRequest)

  url:= MakeUrl(locationRequest)
  DataFromGoogelApi(url,&locationCoordinates)
  lat := locationCoordinates.Results[0].Geometry.Location.Lat
  lng := locationCoordinates.Results[0].Geometry.Location.Lng
  response :=  UpdateInDB(&locationRequest,lat,lng,idValue)
  jsonResponse, _ := json.Marshal(response)
  rw.Header().Set("Content-Type", "application/json")
  rw.WriteHeader(201) // Status code for success
  fmt.Fprintf(rw, "%s", jsonResponse)

}

func getLocationHandler(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {
 idValue := convertToObjectId(p.ByName("location_id"))
 response := getFromDB(idValue)
 jsonResponse, _ := json.Marshal(response)
 rw.Header().Set("Content-Type", "application/json")
 rw.WriteHeader(200) // Status code for success
 fmt.Fprintf(rw, "%s", jsonResponse)
}

func createLocationHandler(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {
  var locationRequest LocationRequest
  var locationCoordinates LocationCoordinates
  jsonDataFromHttp, err := ioutil.ReadAll(request.Body)
  if err != nil {
          panic(err)
  }
  err = json.Unmarshal([]byte(jsonDataFromHttp), &locationRequest)

    url:= MakeUrl(locationRequest)
    DataFromGoogelApi(url,&locationCoordinates)
    lat := locationCoordinates.Results[0].Geometry.Location.Lat
    lng := locationCoordinates.Results[0].Geometry.Location.Lng
    response := SaveInDB(&locationRequest,lat,lng)
    jsonResponse, _ := json.Marshal(response)
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(201) // Status code for success
    fmt.Fprintf(rw, "%s", jsonResponse)
}

func deleteLocationHandler(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {
  idValue := convertToObjectId(p.ByName("location_id"))
  statusCode := delete(idValue)
  rw.WriteHeader(statusCode) // Status code for success
  fmt.Fprint(rw)
}

func createTripHandler(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {

  var tripRequest TripRequest
  jsonDataFromHttp, err := ioutil.ReadAll(request.Body)
  if err != nil {
          panic(err)
  }
  err = json.Unmarshal([]byte(jsonDataFromHttp), &tripRequest)
    response := findShortestPath(tripRequest)
    jsonResponse, _ := json.Marshal(response)
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(201) // Status code for success
    fmt.Fprintf(rw, "%s", jsonResponse)
}

func getTripHandler(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {
  idValue := convertToObjectId(p.ByName("trip_id"))
  response := getBestRouteFromDB(idValue)
  jsonResponse, _ := json.Marshal(response)
  rw.Header().Set("Content-Type", "application/json")
  rw.WriteHeader(200) // Status code for success
  fmt.Fprintf(rw, "%s", jsonResponse)
}

func updateTripHandler(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {
  idValue := convertToObjectId(p.ByName("trip_id"))
  dataFromDb := getBestRouteFromDB(idValue)
  response := putEtaData(dataFromDb)
  jsonResponse, _ := json.Marshal(response)
  rw.Header().Set("Content-Type", "application/json")
  rw.WriteHeader(200) // Status code for success
  fmt.Fprintf(rw, "%s", jsonResponse)
}

func Server() {
  fmt.Println("starting server...........")

  mux := httprouter.New()
  mux.GET("/locations/:location_id", getLocationHandler)
  mux.POST("/locations", createLocationHandler)
  mux.PUT("/locations/:location_id", updateLocationHandler)
  mux.DELETE("/locations/:location_id", deleteLocationHandler)

  //Uber calls
  mux.POST("/trips", createTripHandler)
  mux.GET("/trips/:trip_id", getTripHandler)
  mux.PUT("/trips/:trip_id/request", updateTripHandler)

  server := http.Server{
          Addr:        "0.0.0.0:8080",
          Handler: mux,
  }
  server.ListenAndServe()

}
