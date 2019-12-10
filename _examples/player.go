package main

import "github.com/17xande/go-omxplayer"

func main() {
	omx := go-omxplayer.NewOMXPlayer("file.mp4", false, false, nil)
	waiting := make(chan string)
	defer close(waiting)
	if err := omx.Open(waiting); err != nil {
		log.Println("Error trying to open OMXPlayer: ", err)
	}
	done := <-waiting
}
