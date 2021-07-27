package main

import (
	"fmt"

	"play.ground/pkg/app"
)

func main() {
	fmt.Println("starting ... ")

	if err := app.Run(); err != nil {
		fmt.Println(" ┬───► load failed")
		return
	}
}
