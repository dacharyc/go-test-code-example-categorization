package main

import (
	"fmt"
	"time"
)

func LogStartInfoToConsole(startTime time.Time, fileCount int) {
	fmt.Printf("Processing %d files for %s project\n", fileCount, ProjectName)
	fmt.Println("Starting at ", startTime)
}

func LogFinishInfoToConsole(startTime time.Time, filesProcessed int) {
	endTime := time.Now()
	fmt.Println("Finished at ", endTime)
	fmt.Println("Completed in ", endTime.Sub(startTime))
	fmt.Println("Total snippets processed: ", filesProcessed)
}
