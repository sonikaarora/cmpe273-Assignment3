# cmpe273-Assignment3

Configurations required to execute project
1)	In server/services.go file
-	Go to method connectionWithDB
func connectionWithDB() (sess*mgo.Session, dbName string) {
-	Replace uri value with mongolab configured database
uri="mongodb://admin:admin@ds043694.mongolab.com:43694/addressbook"
-	Replace database(dbName) name with your configured database name like dbName = "addressbook"


----------------------
Command Line Arguments

1) Start server using the following command, server will start listening on port 8080
    Go run startserver.go

2) Start postman or any other REST client 
1.	Create Multiple New Locations using - POST  /locations 
Provide JSON in format as below: 
{
   "name" : "John Smith",
   "address" : "123 Main St",
   "city" : "San Francisco",
   "state" : "CA",
   "zip" : "94113"
}
Expected output, in below format:
{
  "Id": "562b17c52ff8612e86000001",
  "Name": "John Smith",
  "Address": "123 Main St",
  "City": "San Francisco",
  "State": "CA",
  "Zip": "94113",
  "Coordinate": {
    "Lat": 37.7917618,
    "Lng": -122.3943405
  }
}

2.	Get shortest path possible using POST  /trips   API
Provide JSON in format as below: 
{
       "starting_from_location_id" : "565039432ff861eca0000001",
 "location_ids" : ["565039fa2ff861eca0000002","56503a272ff861eca0000003","56503a662ff861eca0000004"]   
}
Expected output:
     {
      "Id":"56503aff2ff861eca0000005",
"Status":"planning",
"Starting_from_location_id":"565039432ff861eca0000001",
"Best_route_location_ids":["565039fa2ff861eca0000002","56503a272ff861eca0000003","56503a662ff861eca0000004"],
"Total_uber_costs":148,
"Total_uber_duration":6153,
"Total_distance":83.94
      }

3.	GET        /trips/{trip_id} # Check the trip details and status
Expected output:

     {
      "Id":"56503aff2ff861eca0000005",
      "Status":"planning",
      "Starting_from_location_id":"565039432ff861eca0000001",
"Best_route_location_ids":["565039fa2ff861eca0000002","56503a272ff861eca0000003","56503a662ff861eca0000004"],
"Total_uber_costs":148,
"Total_uber_duration":6153,
"Total_distance":83.94
     }

4.	PUT        /trips/{trip_id}/request 

Expected output:
      {
 "Id":"56503aff2ff861eca0000005",
 "Status":"requesting",
 "Starting_from_location_id":"565039432ff861eca0000001",
       "Next_destination_location_id":"565039fa2ff861eca0000002",
  "Best_route_location_ids":["565039fa2ff861eca0000002","56503a272ff861eca0000003","56503a662ff861eca0000004"],
       "Total_uber_costs":148,
       "Total_uber_duration":6153,
       "Total_distance":83.94,
       "Uber_wait_time_eta":8
      }

     {
      "Id":"56503aff2ff861eca0000005",
      "Status":"requesting",
      "Starting_from_location_id":"565039432ff861eca0000001",
      "Next_destination_location_id":"56503a272ff861eca0000003",
   "Best_route_location_ids":["565039fa2ff861eca0000002","56503a272ff861eca0000003","56503a662ff861eca0000004"],
     "Total_uber_costs":148,
     "Total_uber_duration":6153,
     "Total_distance":83.94,
     "Uber_wait_time_eta":8}

    {
     "Id":"56503aff2ff861eca0000005",
     "Status":"requesting",
     "Starting_from_location_id":"565039432ff861eca0000001",
     "Next_destination_location_id":"56503a662ff861eca0000004",
  "Best_route_location_ids":["565039fa2ff861eca0000002","56503a272ff861eca00000 03","56503a662ff861eca0000004"],
    "Total_uber_costs":148,
    "Total_uber_duration":6153,
    "Total_distance":83.94,
    "Uber_wait_time_eta":8}

  {
   "Id":"56503aff2ff861eca0000005",
   "Status":"requesting",
   "Starting_from_location_id":"565039432ff861eca0000001",
   "Next_destination_location_id":"565039432ff861eca0000001",
 "Best_route_location_ids":["565039fa2ff861eca0000002","56503a272ff861eca0000003","56503a662ff861eca0000004"],
  "Total_uber_costs":148,
  "Total_uber_duration":6153,
  "Total_distance":83.94,
  "Uber_wait_time_eta":8}

 {
  "Id":"56503aff2ff861eca0000005",
  "Status":"finished",
  "Starting_from_location_id":"565039432ff861eca0000001",
  "Next_destination_location_id":"",
  "Best_route_location_ids":["565039fa2ff861eca0000002","56503a272ff861eca0000003","56503a662ff861eca0000004"],
  "Total_uber_costs":148,
  "Total_uber_duration":6153,
  "Total_distance":83.94,
  "Uber_wait_time_eta":0}












