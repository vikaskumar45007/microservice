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
	http.ListenAndServe(":8080",mux)
	
}