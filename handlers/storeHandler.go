package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/devum25/cloudnativego/models"
	"github.com/devum25/cloudnativego/store"
	"github.com/devum25/cloudnativego/transactionlog"
	"github.com/gorilla/mux"
)

var logger transactionlog.TransactionLogger

func InitializeTransactionLog() error{
	var err error

	logger,err = transactionlog.NewFileTransactionLogger("transaction.log")
	if err != nil{
	  return fmt.Errorf("failed to create event logger: %w",err)
	}

	events,errors:=logger.ReadEvents()
	e,ok := models.Event{},true

	for ok && err == nil{
		select{
		case err,ok = <-errors:
			  if ok{
			  log.Fatalf("some error occured on scan: %e",err)
			  }
		case e,ok = <-events:
			  switch e.EventType{
			  case models.EventDelete:
				  err = store.Delete(e.Key)
			  case models.EventPut:
			     err = store.Put(e.Key,e.Value)
			  }
		}
	}

	logger.Run()

	return err
}

func KeyValuePutHandler(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	key := vars["key"]

	val,err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	err = store.Put(key,string(val))
	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	logger.WritePut(key,string(val))

	w.WriteHeader(http.StatusCreated)
}


func KeyValueGetHandler(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	key := vars["key"]

	

	val,err := store.Get(key)
	if errors.Is(err,store.ErrorNoSuchKey){
		http.Error(w,err.Error(),http.StatusNotFound)
		return
	}
	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}

	w.Write([]byte(val))
}

func KeyValueDeleteHandler(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	key := vars["key"]

	err := store.Delete(key)
	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}

	logger.WriteDelete(key)

	w.WriteHeader(http.StatusOK)
}
