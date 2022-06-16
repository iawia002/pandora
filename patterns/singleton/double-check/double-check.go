package main

import (
	"sync"
)

// Singleton ...
type Singleton struct{}

var (
	instance *Singleton
	lock     sync.Mutex
)

// GetInstance returns the instance.
func GetInstance() *Singleton {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()

		if instance == nil {
			instance = &Singleton{}
		}
	}
	return instance
}
