package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// very simple program that willc heck once every 10min if the ticker api has a newer timestamp that what it has seen.
// if it does it invokes a discord webhook.

func main() {
	webhookURL := "https://discordapp.com/api/webhooks/506006968327208960/PaEfQRa6inOdrT6xEn7wmDu-6LJWxX6yX0x20p42BLnQs2Jpt1W6UWH9pwJMTAIi1ZlZ"
	latestTimestamp := 0
	for {
		fmt.Println("scanning")
		// ask for latest added timestamp
		resp, err := http.Get("https://paragliding-a2-131348.herokuapp.com/paragliding/api/ticker/latest")
		if err != nil {
			fmt.Println(err)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			} else {
				// invoke webhook and post the new timestamp if a newer track is found
				timestampReceivedInt, _ := strconv.Atoi(string(body))
				if timestampReceivedInt > latestTimestamp {
					latestTimestamp = timestampReceivedInt
					var jsonStr = []byte(`{"content":"` + string(body) + `"}`)
					_, postErr := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonStr))
					if postErr != nil {
						fmt.Println(postErr)
					}
				}

			}
		}
		defer resp.Body.Close()

		time.Sleep(10 * time.Minute)
	}
}
