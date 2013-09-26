package logfile

import (
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"testing"
	"time"
)

// LogFile will reopen the log file when the os signal USR1
// is received.
func TestReopenFile(t *testing.T) {
	defer func() {
		// restore defaults
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	}()

	log.SetFlags(0)

	file, err := ioutil.TempFile("/tmp", "logfile-")
	if err != nil {
		t.Fatal(err)
	}
	filename := file.Name()
	file.Close()

	lf, err := New(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer lf.Close()

	log.Print("hi")
	log.Print("igor")

	if content := read(filename); content != "hi\nigor\n" {
		t.Fatalf("Invalid log file content: %#v", content)
	}

	err = os.Rename(filename, filename+"_old")
	if err != nil {
		t.Fatal(err)
	}

	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(100 * time.Millisecond) // :(

	log.Print("keep")
	log.Print("going")

	if content := read(filename + "_old"); content != "hi\nigor\n" {
		t.Fatalf("Old file should keep same content, found: %#v", content)
	}
	if content := read(filename); content != "keep\ngoing\n" {
		t.Fatalf("Invalid content after rotate: %#v", content)
	}

}

func read(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(content)
}
