package main

import "testing"
import "fmt"

func TestM(t *testing.T) {
	str := "asdf"
	fmt.Println(str[len(str)-1])
}
