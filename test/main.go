package main

import (
	"fmt"
	mqtts "github.com/sacca97/unitn-mqtts"
)

func main() {
	mqtts.init()
	fmt.Println("Hello World")
}
