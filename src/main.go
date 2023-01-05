package main

import (
	"flag"
	"fmt"
	"go-toggl-api-client/src/baseClient"
)

func main() {
	apiKey := flag.String("ak", "", "api key")
	flag.Parse()
	if len(*apiKey) == 0  {
		fmt.Println("api key is nothing")
		return
	}

	client, err := baseClient.NewClient(*apiKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, err := client.GetRequest("https://api.track.toggl.com/api/v9/me/time_entries")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(*body))
}
