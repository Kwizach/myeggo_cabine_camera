package main

import (
	"fmt"
	"os"
)

var logURL string

func log(err error) {
	if err == nil {
		return
	}

	strErr := fmt.Sprintf("%s", err)
	if logURL != "" {
		getLogURL()
	}

	f, err := os.OpenFile(logURL, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Can't open %s\nOriginal error: %s", logURL, strErr)
		return
	}
	defer f.Close()

	_, err = f.WriteString(strErr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Can't wrtite to %s\nOriginal error: %s", logURL, strErr)
		return
	}
}

func getLogURL() {
	logURL = fmt.Sprintf("./%s.log", os.Args[0])
}
