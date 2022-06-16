package main

import (
	"sync"
)

// Singleton ...
type Singleton struct{}

var (
	instance *Singleton
	once     sync.Once
)

// GetInstance returns the instance.
func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{}
	})
	return instance
}
