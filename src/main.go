package main

import (
	"fmt"
	"go-toggl-api-client/src/baseClient"
)

func main() {
	client, err := baseClient.NewClient()
	if err != nil {
		return
	}
	body, err := client.GetRequest("https://api.track.toggl.com/api/v9/me/time_entries")
	if err != nil {
		return
	}
	fmt.Print(string(*body))
}
