package utility

import (
	"fmt"
	"log"

	// "os"
	"time"
)

func Log(info interface{}) { //Prints to logfile
	// If the file doesn't exist, create it, or append to the file
	// f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close() //ensure close

	// log.SetOutput(f)
	log.Println(info) //print to file
}

func TimeTrack(start time.Time, name string) { //logs function time. to use-> defer timeTrack(time.Now(), "label") at start of function
	elapsed := time.Since(start)
	fmt.Println(name, " took ", elapsed)
}
