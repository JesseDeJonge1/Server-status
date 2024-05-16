package main

import (
	"net/http"
	"sync"
	"time"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/fatih/color"
)

func pingServer(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		color.Red("Error: %v", err)
		return url + ": Error"
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			color.Green(url + ": Server online")
			return url + ": Server online"
		} else {
			color.Yellow(url + ": Server not responding")
			return url + ": Server not responding"
		}
	}
}

func main() {
	websites := []string{
		"http://app.ditkanik.nu",
		"http://bfi.folioo.at",
		"http://ust.ditkanik.nu",
	}

	var wg sync.WaitGroup
	statuses := make([]string, len(websites))

	for i, website := range websites {
		wg.Add(1)
		go func(i int, website string) {
			defer wg.Done()
			statuses[i] = pingServer(website)
		}(i, website)
	}

	wg.Wait()

	// Send macOS notification
	note := gosxnotifier.NewNotification("Server statuses: \n")
	note.Title = "Server Status"
	note.Subtitle = "Multiple servers"
	note.Sound = gosxnotifier.Default
	for _, status := range statuses {
		note.Message += status + "\n"
	}
	err := note.Push()

	// If there was an error, print it out
	if err != nil {
		color.Red("Error: %v", err)
	}

	time.Sleep(5 * time.Minute)
}
