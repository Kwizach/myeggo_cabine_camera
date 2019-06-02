package main

import (
	"fmt"
	"os"
	"time"
)

var logURL string

func log(err error) {
	if err == nil {
		return
	}

	// Get time and format it
	now := fmt.Sprintf("%s", time.Now().Format("2006-01-02T15:04:05"))
	// Stringify error
	strErr := fmt.Sprintf("%s", err)
	//  Get log file URL
	if logURL != "" {
		getLogURL()
	}

	// Open log file
	f, err := os.OpenFile(logURL, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		// Error opening log file write to Stderr
		fmt.Fprintf(os.Stderr, "%s - Error Can't open %s\nOriginal error: %s", now, logURL, strErr)
		return
	}
	defer f.Close()

	// Write to log file
	_, err = f.WriteString(now + " - " + strErr)
	if err != nil {
		// Error while writting to log file, write to Stderr
		fmt.Fprintf(os.Stderr, "%s - Error Can't wrtite to %s\nOriginal error: %s", now, logURL, strErr)
		return
	}
}

func getLogURL() {
	logURL = fmt.Sprintf("./%s.log", os.Args[0])
}
