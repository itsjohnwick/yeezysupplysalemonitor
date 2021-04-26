package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var webhook string

func main() {

	webhook = readTxt() //pull webhook link from txt file and store it in global var

monitorLoop:
	for {

		httpMethod := "GET"
		url := "https://www.yeezysupply.com/hpl/content/yeezy-supply/config/US/waitingRoomConfig.json" // sets the url variable as yeezysupply api url
		client := http.Client{}                                                                        // sets the client variable as http.Client{}
		request, err := http.NewRequest(httpMethod, url, nil)                                          // sets the request as a new request with GET method, pointing to the yeezysupply url, with an empty body

		if err != nil { // err handling (print error to console, wait 5 sec, try again)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue monitorLoop
		}

		request.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.128 Safari/537.36")    // sets the key header as user agent, and the value as the user agent
		request.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9") // sets the key header as accept, and value

		response, err := client.Do(request) // executes the request

		if err != nil { // err handling (print error to console, wait 5 sec, try again)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue monitorLoop
		}

		// fmt.Println(response.StatusCode) // print the response (200, 429, etc)

		responseString, err := ioutil.ReadAll(response.Body) // reads the response body, returns responseString (variable) and an err if there is one

		if err != nil { // err handling (print error to console, wait 5 sec, try again)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue monitorLoop
		}

		// fmt.Println(string(responseString)) // print the response body as a string

		if strings.Contains(string(responseString), "sale_started") { // if the responseString contains sale_started, print sale has started. Else, print sale has not started
			discordWebhook()
			fmt.Println("Sale is live.")
		} else {
			fmt.Println("Sale is not live.")
		}

		response.Body.Close()       // close the request
		time.Sleep(5 * time.Second) // sleep for 5 seconds
		continue monitorLoop        // do the loop again
	}
}

func discordWebhook() {

discordLoop:
	for {
		httpMethod := "POST"
		url := webhook // url == webhook link
		client := http.Client{}

		request, err := http.NewRequest(httpMethod, url, bytes.NewBufferString("{\"content\":\"Sale has started. Start tasks.\"}")) // sets body to the sale has started message

		if err != nil { // err handling (print error to console, wait 5 sec, try again)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue discordLoop
		}

		request.Header.Add("Content-Type", "application/json") // adds header key as content type and value as application/json

		response, err := client.Do(request) // executes the request

		if err != nil { // err handling (print error to console, wait 5 sec, try again)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue discordLoop
		}

		fmt.Println(response.StatusCode) // print the response (200, 429, etc)

		responseString, err := ioutil.ReadAll(response.Body) // reads the response body, returns responseString (variable) and an err if there is one

		if err != nil { // err handling (print error to console, wait 5 sec, try again)
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue discordLoop
		}

		fmt.Println(string(responseString)) // prints response from discord api
		response.Body.Close()               // closes request
		break discordLoop                   // breaks the discord loop
	}

}

func readTxt() string {
	b, err := ioutil.ReadFile("webhook.txt") // reads webhook txt file

	if err != nil { // error handling
		log.Fatal(err)
	}

	webhook := string(b) // convert content to a string

	if webhook == "webhookhere" || webhook == "" {
		fmt.Println("Please put your webhook in webhook.txt")
		time.Sleep(5 * time.Second)
		main()
	}

	return webhook
}
