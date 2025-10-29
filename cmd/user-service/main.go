package main

import (
	"log"
	"microservice/internal/db"
	"microservice/internal/user"
	"net/http"
)

func main(){
	dbconn,err := db.Connect()
	if err != nil{
		log.Fatal(err)
	}
	defer dbconn.Close()

	h := user.NewHandler(dbconn)

	mux := http.NewServeMux()
	mux.HandleFunc("/user/create",h.Create)
	log.Printf("Server is listening at 8080 port.")
	http.ListenAndServe(":8080",mux)
	
}