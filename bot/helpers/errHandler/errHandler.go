/*
Package "errHandler" is a package that handles all kinds of error(s).
This package should handle all error(s).
*/
package errHandler

import (
	log "github.com/sirupsen/logrus"
)

// Error function returns nothing as it only handles error and log it.
func Error(err error) {
	if err != nil {
		log.Error(err)
	}
}

// Fatal function returns nothing as it only handles error and log it.
func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// SendError function will send an error message to the chat.
func SendError(err error) {
	if err != nil {

	}
}
