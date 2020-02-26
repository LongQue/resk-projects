package main

import "fmt"

type Status int

const (
	StatusOk Status=iota
	StatusFailed
	StatusTimeout
)

func main() {
	var s Status
	s=StatusFailed
	fmt.Println(s)
}