package logline

import (
	"bytes"
	"log"
	"os"
	"sync"
	"testing"
)

func ExampleNew() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	ll := New()

	ll.Print("hi")
	ll.Printf("go%s", "lang")
	ll.Write()

	// Output:
	// hi. golang.
}

func ExampleNewWithInfo() {
	logger := log.New(os.Stdout, "main: ", 0)

	ll := NewWithInfo(Info{
		Logger:    logger,
		Separator: ", ",
	})

	ll.Print("hi")
	ll.Printf("go%s", "lang")
	ll.Write()

	// Output:
	// main: hi, golang,
}

// One Logline object can be used by multiple gotoutines.
// This test should not show warnings when using -race
func TestThreadSafety(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)

	ll := NewWithInfo(Info{
		Logger: logger,
	})

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		ll.Print("hi")
		wg.Done()
	}()
	go func() {
		ll.Printf("hi %s", "golang")
		wg.Done()
	}()
	go func() {
		ll.Write()
		wg.Done()
	}()

	wg.Wait()
}
