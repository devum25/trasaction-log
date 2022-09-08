package store

import (
	"errors"
	"sync"
)

var persistentStore = struct{
   sync.RWMutex 
   m map[string]string
}{m:make(map[string]string)}
var ErrorNoSuchKey = errors.New("no such key")

func Put(key string,value string) error{
	persistentStore.Lock()
	persistentStore.m[key] = value
	persistentStore.Unlock()
	return nil
}

func Get(key string) (string,error){
	persistentStore.RLock()
	defer persistentStore.RUnlock()
	val,ok := persistentStore.m[key]
	if !ok{
         return "",ErrorNoSuchKey
	}

	return val,nil
}

func Delete(key string) error{
	persistentStore.Lock()
	defer persistentStore.RUnlock()
	delete(persistentStore.m,key)
	return nil
}