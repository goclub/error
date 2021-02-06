package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.Open("/nonexistent.txt")
	if err != nil {
		if err, ok := err.(*os.PathError); ok {
			log.Print(err.Op, err.Path, " failed")
			return
		}
		log.Print(err)
	}
	log.Print(f.Name(), "opened successfully")
}