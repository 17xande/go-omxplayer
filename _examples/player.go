package main

import (
	"log"

	"github.com/17xande/go-omxplayer"
)

func main() {
	omx, err := omxplayer.NewOMXPlayer("file.mp4", false, false, nil)
	if err != nil {
		panic(err)
	}
	waiting := make(chan string)
	defer close(waiting)
	if err := omx.Open(waiting); err != nil {
		log.Println("Error trying to open OMXPlayer: ", err)
	}
	<-waiting
}
