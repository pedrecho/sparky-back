package main

import (
	"os"
	"sparky-back/internal/app"
)

const configEnv = "CONFIG"

func main() {
	val, ok := os.LookupEnv(configEnv)
	if !ok {
		panic("no config env")
	}
	if err := app.Run(val); err != nil {
		panic(err)
	}
}
