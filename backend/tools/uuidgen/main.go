package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	id, err := uuid.NewV7()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return
	}
	fmt.Println(id.String())
}
