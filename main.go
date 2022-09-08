package main

import (
	"log"
	"net/http"

	"github.com/devum25/cloudnativego/handlers"
	"github.com/gorilla/mux"
)


func main(){
	handlers.InitializeTransactionLog()
	r := mux.NewRouter()


	r.HandleFunc("/v1/{key}",handlers.KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}",handlers.KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}",handlers.KeyValueDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080",r))
}

