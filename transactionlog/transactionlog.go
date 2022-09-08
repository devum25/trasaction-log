package transactionlog

import "github.com/devum25/cloudnativego/models"

type TransactionLogger interface{
	WriteDelete(key string)
	WritePut(Key string,val string)
	Err() <-chan error
	ReadEvents() (<-chan models.Event,<-chan error)
	Run()
}