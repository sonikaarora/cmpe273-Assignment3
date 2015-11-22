package main
import (
  "fmt"
  "location/server"
)

func main() {
 go server.Server()

 var input int
 fmt.Scanln(&input)


}
