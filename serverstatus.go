package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/getlantern/systray"
)

func pingServer(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			color.Red(url + ": Connection refused")
			sendNotification("Server Status", url+": Connection refused")
			return url + ": Connection refused"
		}
		color.Red("Error: %v", err)
		sendNotification("Server Status", "Error: "+err.Error())
		return url + ": Error"
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		color.Green(url + ": Server online")
		return url + ": Server online"
	default:
		color.Yellow(url + ": Server not responding")
		sendNotification("Server Status", url+": Server not responding")
		return url + ": Server not responding"
	}
}

func sendNotification(title string, message string) {
	cmd := exec.Command("osascript", "-e", `display notification "`+message+`" with title "`+title+`"`)
	err := cmd.Run()
	if err != nil {
		log.Println("Error sending notification:", err)
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
		menuItem := systray.AddMenuItem(website, "Checking status...")
		menuItems[website] = menuItem

		go func(website string, menuItem *systray.MenuItem) {
			for {
				<-menuItem.ClickedCh
				cmd := exec.Command("open", website)
				err := cmd.Run()
				if err != nil {
					log.Println("Error opening website:", err)
				}
			}
		}(website, menuItem)
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

	systray.SetIcon(loadIcon("your-path"))

	go func() {
		for {
			updateStatuses(websites, menuItems)
			time.Sleep(1 * time.Minute)
		}
	}()
}

func onExit() {}

func main() {
	systray.Run(onReady, onExit)
}
