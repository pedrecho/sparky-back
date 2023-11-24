package main

import (
	"fmt"
	"time"
)

func main() {
	layout := "2006-01-02"
	str := "0"

	t, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(t)
}
