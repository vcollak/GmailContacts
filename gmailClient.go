/*

App that connects to Gmail via Gmail api and saves all emails from "To", "From", and "Cc" into a MongoDB

Resources:
https://developers.google.com/gmail/api/quickstart/go
https://console.developers.google.com
https://godoc.org/google.golang.org/api/gmail/v1
https://tools.ietf.org/html/rfc4021


*/
package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/mail"
	"strconv"
	"strings"
)

//mongo DB
var db = new(MongoDB)

func getGmailClient() (*gmail.Service, error) {
	ctx := context.Background()

	config, err := getGmailConfig()
	if err != nil {
		return nil, err
	}

	client := getClient(ctx, config)
	svc, err := gmail.New(client)
	if err != nil {
		return nil, err
	}

	return svc, nil

}

//see if the email is one of the known emails
func isKnownEmail(email string) bool {

	for _, e := range knownEmails {

		if strings.ToUpper(email) == strings.ToUpper(e) {
			return true
		}
	}

	return false
}

func saveHeaderFields(headerValue string) {

	emails, err := mail.ParseAddressList(headerValue)
	if err != nil {
		log.Println("Unable to parse:", headerValue)
	} else {

		for _, v := range emails {

			name := v.Name
			email := v.Address

			if !isKnownEmail(email) {
				err = db.SetContact(name, email)
				if err != nil {
					log.Println("Unable to save:", email)
				}
			} else {
				log.Println("Known email. Ignoring:", email)
			}
		}
	}
}

func processMessages() {

	//get the mongo object
	db.NewMongo(server, dbName, accountName)

	//close the sessions at the end
	defer db.session.Close()

	svc, err := getGmailClient()
	if err != nil {
		log.Fatal("Unable to access Gmail. Error:", err)
	}

	//get messages
	pageToken := ""
	firstMessage := true

	for {

		var req *gmail.UsersMessagesListCall
		lastDate, err := db.LastDate()

		if lastDate == "" {
			log.Println("Retrieving all messages.")
			req = svc.Users.Messages.List("me")

		} else {
			log.Println("Retrieving messages starting on", lastDate)
			req = svc.Users.Messages.List("me").Q("after: " + lastDate)
		}

		if pageToken != "" {
			req.PageToken(pageToken)
		}
		r, err := req.Do()

		if err != nil {
			log.Printf("Unable to retrieve messages: %v", err)
			continue
		}

		log.Printf("--------------")
		log.Printf("Processing %v messages...\n", len(r.Messages))
		for _, m := range r.Messages {

			msg, err := svc.Users.Messages.Get("me", m.Id).Do()
			if err != nil {
				log.Printf("Unable to retrieve message %v: %v", m.Id, err)
				continue
			}

			lastMessageRetrievedDate, err := msToTime(strconv.FormatInt(msg.InternalDate, 10))
			if err != nil {
				log.Println("Unable to parse message date", err)
			}

			//message date
			log.Println(lastMessageRetrievedDate)

			if firstMessage {

				//set the last known date
				currentDate := lastMessageRetrievedDate.Format("2006/01/02")
				err = db.SetLastDate(currentDate)
				if err != nil {
					log.Println("Unable to save:", currentDate)
				} else {
					log.Println("Saved:", currentDate)
					firstMessage = false
				}

			}

			for _, h := range msg.Payload.Headers {

				//prints all header values
				//fmt.Println(h.Name + ":" + h.Value)

				if h.Name == "From" {

					log.Println("From:" + h.Value)
					saveHeaderFields(h.Value)

				} else if h.Name == "To" {

					log.Println("To:" + h.Value)
					saveHeaderFields(h.Value)

				} else if h.Name == "Cc" {

					log.Println("Cc:" + h.Value)
					saveHeaderFields(h.Value)

				} else if h.Name == "Subject" {
					log.Println("Subject:" + h.Value)
				}

			}

			fmt.Println("")

		}

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken

		//break

	}

}
