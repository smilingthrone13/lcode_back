package main

import (
	"lcode/config"
	"lcode/internal/app"
	"log"
)

func main() {
	cfg, err := config.Init("./")
	if err != nil {
		log.Fatal(err)
	}

	app.Init(cfg).Run()
}
