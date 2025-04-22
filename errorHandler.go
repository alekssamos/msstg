package main

import (
	"log"
)

func LogError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func Must(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
