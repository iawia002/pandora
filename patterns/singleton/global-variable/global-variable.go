package main

// Singleton ...
type Singleton struct{}

var instance *Singleton

func init() {
	instance = &Singleton{}
}

// GetInstance returns the instance.
func GetInstance() *Singleton {
	return instance
}
