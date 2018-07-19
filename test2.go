package main

import (
	"github.com/satori/go.uuid"
	"fmt"
)

//var name = "cloudware-" + uuid.Must(uuid.NewV4()).String()

func main() {
	for i := 0; i < 200; i++ {
		fmt.Println(uuid.Must(uuid.NewV4()).String())
	}
}
