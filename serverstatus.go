package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/getlantern/systray"
)

func pingServer(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		color.Red("Error: %v", err)
		return url + ": Error"
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		color.Green(url + ": Server online")
		return url + ": Server online"
	default:
		color.Yellow(url + ": Server not responding")
		return url + ": Server not responding"
	}
}

func updateStatuses(websites []string, menuItems map[string]*systray.MenuItem) {
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

	for i, website := range websites {
		menuItems[website].SetTitle(statuses[i])
	}
}

func createMenuItems(websites []string) map[string]*systray.MenuItem {
	menuItems := make(map[string]*systray.MenuItem)
	for _, website := range websites {
		menuItems[website] = systray.AddMenuItem(website, "Checking...")
	}
	return menuItems
}

func loadIcon(path string) []byte {
	iconData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return iconData
}

func onReady() {
	websites := []string{
		"your-websites",
	}

	menuItems := createMenuItems(websites)

	systray.SetIcon(loadIcon("your/path"))

	go func() {
		for {
			updateStatuses(websites, menuItems)
			time.Sleep(5 * time.Minute)
		}
	}()
}

func onExit() {}

func main() {
	systray.Run(onReady, onExit)
}
