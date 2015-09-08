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
	"google.golang.org/api/gmail/v1"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

//mongo DB
var mongo = &MongoDB{}

//converts the UNIX epoch string to a time
func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

//see if the email is one of the known emails
func isKnownEmail(email string) bool {

	knownEmails := []string{"vlad@collak.net", "vcollak@gmail.com", "vcollak@ignitedev.com", "info@slovacihouston.com", "vlad@openkloud.com", "vladimir.collak@ignitemediallc.com", "vladimir.collak@ignitemediahosting.com"}
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
				err = mongo.SetContact(name, email)
				if err != nil {
					log.Println("Unable to save:", email)
				}
			} else {
				log.Println("Known email. Ignoring:", email)
			}
		}
	}
}

func main() {

	//get the mongo object
	mongo = NewMongo("127.0.0.1", "gmailContacts", "vcollak@gmail.com")

	svc, err := getGmailClient()
	if err != nil {
		log.Fatal("Error:", err)
	}

	//get messages
	pageToken := ""
	firstMessage := true

	for {

		var req *gmail.UsersMessagesListCall
		lastDate, err := mongo.LastDate()

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
				err = mongo.SetLastDate(currentDate)
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
