// Package amalgam provides a brutally simple webserver for no-nonsense news aggregation.
// It aims for lightweight hardware requirements, extensibility, and simplicity.

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/profile"
)

func frontPage(countC chan int) func(w http.ResponseWriter, r *http.Request) {
	// TODO all references to counts will end up being links and content.

	// Ideally this should only show for a second or so while we make our initial fetching of data.
	// Thereafter it will be replaced by new content
	latestContent := DefaultFrontPage()

	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case newCount := <-countC:
			currCount = newCount
			fmt.Println("new count! updating page")
		}

		fmt.Println("handling request")
		fmt.Fprintf(w, "<h1>%s</h1><div>%s. <b>You are on counter: %d</b></div>", []byte("Hello"), []byte("There aye!"), currCount)
	}
}

func timedLog(c chan int) {
	fmt.Println("starting timer...")
	i := 1

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Printf("on iteration: %d\n", i)
			i++

			select {
			case c <- i:
			default:
				// empty out the single-buffer channel to keep it updated with the latest data
				<-c
				c <- i
			}
		}
	}()
}

func main() {
	// cpu/memory profiling
	defer profile.Start(profile.MemProfile).Stop()

	// kick off links fetcher
	counterC := make(chan int, 1)
	go timedLog(counterC)

	http.HandleFunc("/view", frontPage(counterC))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
