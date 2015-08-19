/*

App that connects to Gmail via gmail api and lists all messages and their To, From, Cc

Resources:

https://developers.google.com/gmail/api/quickstart/go
https://console.developers.google.com
https://godoc.org/google.golang.org/api/gmail/v1
https://tools.ietf.org/html/rfc4021


*/

package main

import (
	"fmt"
	"log"
)

func main() {

	svc, err := getGmailClient()
	if err != nil {
		log.Fatal("Error:", err)
	}

	//get messages
	pageToken := ""
	for {

		req := svc.Users.Messages.List("me")

		if pageToken != "" {
			req.PageToken(pageToken)
		}
		r, err := req.Do()

		if err != nil {
			log.Fatalf("Unable to retrieve messages: %v", err)
		}

		log.Printf("--------------")
		log.Printf("Processing %v messages...\n", len(r.Messages))
		for _, m := range r.Messages {

			msg, err := svc.Users.Messages.Get("me", m.Id).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve message %v: %v", m.Id, err)
			}

			for _, h := range msg.Payload.Headers {
				//fmt.Println(h.Name + ":" + h.Value)

				if h.Name == "Subject" {

					log.Println("Subject:" + h.Value)

				} else if h.Name == "From" {

					log.Println("From:" + h.Value)

				} else if h.Name == "To" {

					log.Println("To:" + h.Value)

				} else if h.Name == "Cc" {

					log.Println("Cc:" + h.Value)
				}

			}
			fmt.Println("")

		}

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}
}
