package main

import (
	"fmt"
	"log"
	"os"
	"tutgo/pkg/api"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate)

	var port string
	port = os.Getenv("PORT")
	if port == "" {
		fmt.Println("env port empting, setting in file")
		port = "8080"
	}
	server := api.NewServer(fmt.Sprintf(":%s", port))
	server.Run()

}
