package main

import (
	"fmt"
)

// Options defines options needed to do something.
type Options struct {
	fieldA  string
	fieldB  int
	switchA bool
}

// WithFieldA sets the fieldA option.
// The default value is empty string.
func WithFieldA(fieldA string) func(opts *Options) {
	return func(opts *Options) {
		opts.fieldA = fieldA
	}
}

// WithFieldB sets the fieldB option.
// The default value is 0.
func WithFieldB(fieldB int) func(opts *Options) {
	return func(opts *Options) {
		opts.fieldB = fieldB
	}
}

// WithSwitchA sets the switchA option.
// The default value is true.
func WithSwitchA(switchA bool) func(opts *Options) {
	return func(opts *Options) {
		opts.switchA = switchA
	}
}

// Client ...
type Client struct {
	Addr   string
	FieldA string
	FieldB int
}

// New returns a new Client implementation.
func New(addr string, options ...func(*Options)) *Client {
	opts := &Options{
		switchA: true,
	}
	for _, f := range options {
		f(opts)
	}

	if opts.switchA {
		addr = fmt.Sprintf("[%s]", addr)
	}

	return &Client{
		Addr:   addr,
		FieldA: opts.fieldA,
		FieldB: opts.fieldB,
	}
}

func main() {
	client1 := New("localhost")
	client2 := New("localhost", WithFieldA("hello"))
	client3 := New("localhost", WithFieldA("hello"), WithFieldB(1))
	client4 := New(
		"localhost",
		WithFieldA("hello"), WithFieldB(1), WithSwitchA(false),
	)

	fmt.Println(client1, client2, client3, client4)
}
