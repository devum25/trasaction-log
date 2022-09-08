package transactionlog

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/devum25/cloudnativego/models"
)

type FileTransactionLogger struct{
   events chan<- models.Event
   errors <-chan error
   lastSequence uint64
   file *os.File
   lock sync.RWMutex
}

func NewFileTransactionLogger(fileName string) (TransactionLogger,error){
	file,err := os.OpenFile(fileName,os.O_RDWR|os.O_APPEND|os.O_CREATE,0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
		}

		return &FileTransactionLogger{file: file}, nil
}


func (l *FileTransactionLogger) WritePut(key,val string){
      l.events <- models.Event{EventType: models.EventPut,Key: key,Value: val}
	  fmt.Print("devum")
}

func (l *FileTransactionLogger) WriteDelete(key string){
	l.events <- models.Event{EventType: models.EventDelete,Key: key}
}

func (l *FileTransactionLogger) Err() <-chan error{
	return l.errors
}

func (l *FileTransactionLogger) Run(){
	events := make(chan models.Event,16)
	l.events = events

	
	//l.events <- models.Event{}

	errors := make(chan error,1)
    l.errors = errors

	go func(){
		for e := range events{
			l.lastSequence++
            fmt.Printf("Writing log with value:%s,%s",e.Key,e.Value)
			l.lock.Lock()
			_,err := fmt.Fprintf(l.file,"%d\t%d\t%s\t%s\n",l.lastSequence,e.EventType,e.Key,e.Value)
			l.lock.Unlock()

			if err != nil{
				errors <- err
				return
			}
		}
	}()
}

func (l *FileTransactionLogger) ReadEvents() (<-chan models.Event,<-chan error){
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan models.Event)
	outError := make(chan error,1)


	go func(){
		var e models.Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan(){
			line := scanner.Text()

			l.lock.RLocker().Lock()
			if _,err := fmt.Sscanf(line,"%d\t%d\t%s\t%s",&e.Sequence,&e.EventType,&e.Key,&e.Value);err != nil{
				outError <- fmt.Errorf("input parse error: %w",err)

				return
			}
		l.lock.RLocker().Unlock()

		if l.lastSequence >= e.Sequence{
			outError <- fmt.Errorf("transaction numbers out of sequence")
			return
		}

		l.lastSequence =e.Sequence
		outEvent <- e
	}

	}()

	if err := scanner.Err();err != nil{
		outError <- fmt.Errorf("transaction log read failure:%w",err)
	}

	return outEvent,outError
}