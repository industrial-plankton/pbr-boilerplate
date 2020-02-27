package utility

import (
	"fmt"
	"log"

	// "os"
	"time"
)

//Log logs info to logfile
func Log(info interface{}) { //Prints to logfile
	// If the file doesn't exist, create it, or append to the file
	// f, err := os.OpenFile("/home/cameron/Go_Server/logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close() //ensure close

	// log.SetOutput(f)
	log.Println(info) //print to file
}

//TimeTrack logs function time. to use-> defer timeTrack(time.Now(), "label") at start of function
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Println(name, " took ", elapsed)
}
