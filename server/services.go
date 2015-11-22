package server


import (
      "fmt"
      "strings"
      "net/http"
      "io/ioutil"
      "encoding/json"
      "gopkg.in/mgo.v2"
      "os"
      "gopkg.in/mgo.v2/bson"
      "encoding/hex"
      )


func convertToObjectId(id string) bson.ObjectId {
	objectIdFormat, err := hex.DecodeString(id)
	if err != nil || len(objectIdFormat) != 12 {
		panic(fmt.Sprintf("Invalid input to ObjectIdHex: %q", id))
	}
	return bson.ObjectId(objectIdFormat)
}

func MakeUrl(request LocationRequest) string{
  var url="http://maps.google.com/maps/api/geocode/json?address="
   AddressArray := strings.Split(request.Address," ")
   CityArray := strings.Split(  request.City," ")
   StateArray := strings.Split(request.State," ")
   ZipCode :=   request.Zip
   for count:=0 ; count< len(AddressArray); count++ {
     url = url+AddressArray[count]+"+"
   }
   for count:=0 ; count< len(CityArray); count++ {
     url = url+CityArray[count]+"+"
   }
   for count:=0 ; count< len(StateArray); count++ {
     url = url+StateArray[count]+"+"
   }
   url = url+ZipCode
   fmt.Println(url)
   return url
}

func DataFromGoogelApi(url string,locationCoordinates *LocationCoordinates) {
  resp, err := http.Get(url)
  if err!=nil {
    fmt.Println("Error from DataFromGoogelApi function: ",err)
  } else {
  defer resp.Body.Close()
         jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
         if err != nil {
                 panic(err)
         }
         err = json.Unmarshal([]byte(jsonDataFromHttp), &locationCoordinates)

         if err != nil {
           fmt.Println("error: ",err)
                 panic(err)
         }
   }
}

func connectionWithDB() (sess*mgo.Session, dbName string) {
  uri := os.Getenv("MONGOHQ_URL")
   if uri == "" {
     fmt.Println("no connection string provided")
     os.Exit(1)
   }
   uri = "mongodb://admin:admin@ds043694.mongolab.com:43694/addressbook"
   dbName = "addressbook"
   session, err := mgo.Dial(uri)
   if err != nil {
     fmt.Printf("Can't connect to mongo, go error %v\n", err)
     os.Exit(1)
   }
  return session, dbName
}

//Get a Location
func getFromDB(id bson.ObjectId) (LocationResponse){
  session, dbName := connectionWithDB()
  defer session.Close()
  session.SetSafe(&mgo.Safe{})
    c := session.DB(dbName).C("user")
  var resultFromDb LocationResponse
  err := c.Find(bson.M{"id": id}).One(&resultFromDb)
   if err != nil {
       panic(err)
   }

  return resultFromDb
}

//Delete a Location
func delete(id bson.ObjectId) (status int){
  session,dbName := connectionWithDB()
  defer session.Close()
  session.SetSafe(&mgo.Safe{})
    c := session.DB(dbName).C("user")
    err :=  c.Remove(bson.M{"id": id})
     if err != nil {
         panic(err)
     }
     return 200
}

//Create New Location
func SaveInDB(locationRequest *LocationRequest, lat float64, lng float64) (LocationResponse){
  session, dbName := connectionWithDB()
  defer session.Close()
  session.SetSafe(&mgo.Safe{})
  c := session.DB(dbName).C("user")

  data := LocationResponse {
          // _Id: bson.NewObjectId(),
           Id: bson.NewObjectId(),
           Name: locationRequest.Name,
           Address: locationRequest.Address,
           City: locationRequest.City,
           State: locationRequest.State,
           Zip:locationRequest.Zip,
           Coordinate: Coordinates{
                         Lat:lat,
                         Lng:lng,
                    },
  }

err := c.Insert(data)
if err != nil {
  fmt.Println("error while inserting..........")
  panic(err)
}
return data
}



//Update a Location
func UpdateInDB(locationRequest *LocationRequest, lat float64, lng float64, idValue bson.ObjectId) (LocationResponse){
  session,dbName := connectionWithDB()
  defer session.Close()
  session.SetSafe(&mgo.Safe{})
  c := session.DB(dbName).C("user")

  data := LocationResponse {
           Id: bson.NewObjectId(),
           Name: locationRequest.Name,
           Address: locationRequest.Address,
           City: locationRequest.City,
           State: locationRequest.State,
           Zip:locationRequest.Zip,
           Coordinate: Coordinates{
                         Lat:lat,
                         Lng:lng,
                    },
  }
  var resultFromDb LocationResponse
  err := c.Find(bson.M{"id": idValue}).One(&resultFromDb)
  data.Id = resultFromDb.Id
  data.Name = resultFromDb.Name
  err = c.Update(bson.M{"id": idValue}, data)
    if err != nil {
      fmt.Printf("Can't update document %v\n", err)
      os.Exit(1)
    }
    return data
}
