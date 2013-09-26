// Package logfile set the log output to a file and reopens this file on USR1
package logfile

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type LogFile struct {
	Name string
	File *os.File
}

// New creates a new LogFile.
//
// Will open (or create) the file in append mode and set it as the
// default log output.
//
// Start listening for USR1 signal to reopen the log file.
func New(name string) (*LogFile, error) {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}

	log.SetOutput(file)

	logfile := &LogFile{name, file}
	logfile.listen()

	return logfile, nil
}

func (lf *LogFile) Close() {
	lf.File.Close()
}

func (lf *LogFile) listen() {
	usr1 := make(chan os.Signal, 1)
	signal.Notify(usr1, syscall.SIGUSR1)

	go func() {
		for {
			<-usr1
			file, err := os.OpenFile(lf.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Printf("Failed to reopen log file %#v (%s). Will continue logging to stderr.", lf.Name, err)
				log.SetOutput(os.Stderr)
			} else {
				log.SetOutput(file)
			}
		}
	}()
}
