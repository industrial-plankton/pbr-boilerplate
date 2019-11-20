package utility

import (
	"fmt"
	"log"
	"os"
	"time"
)

func Log(info interface{}) {
	// If the file doesn't exist, create it, or append to the file
	log.Println(info)
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(info)
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Println(name, " took ", elapsed)
}