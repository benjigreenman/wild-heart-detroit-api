package services

import (
	"log"
)

// TODO enumify labels

// Log logs an info message with a label and optional error
func Log(label string, message string, err error) {
	prefix := "[" + label + "] "
	if err != nil {
		message = message + ": " + err.Error()
	}
	log.Printf("%s%s", prefix, message)
}

// Error logs an error message with a label and error
func Error(label string, message string, err error) {
	prefix := "[" + label + "] "
	if err != nil {
		message = message + ": " + err.Error()
	}
	log.Printf("%sERROR: %s", prefix, message)
}
