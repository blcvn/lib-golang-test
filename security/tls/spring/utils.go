package spring

import "log"

var logmsgs bool

func init() {
	logmsgs = true
}

// Output a message if logmsgs is true.
func Output(format string, args ...interface{}) {
	if logmsgs {
		log.Printf(format, args...)
	}
}
