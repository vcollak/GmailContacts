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
	"google.golang.org/api/gmail/v1"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

//converts the UNIX epoch string to a time
func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

//returns the last date we processed
func getLastDate() string {

	lastDate, err := getLastDateFromMongo()
	if err != nil {
		log.Println("Unable to get last date from DB. Error:", err)
	}

	return lastDate

}

//see if the email is one of the known emails
func isKnownEmail(email string) bool {

	knownEmails := []string{"vlad@collak.net", "vcollak@gmail.com", "vcollak@ignitedev.com", "info@slovacihouston.com", "vlad@openkloud.com"}
	for _, e := range knownEmails {

		if strings.ToUpper(email) == strings.ToUpper(e) {
			return true
		}
	}

	return false
}

func parseAndSave(headerValue string) {

	e, err := mail.ParseAddress(headerValue)
	if err != nil {
		log.Println("Unable to parse:", headerValue)
	} else {

		name := e.Name
		email := e.Address

		if !isKnownEmail(email) {
			saveContact(name, email)
		}

	}
}

func main() {

	svc, err := getGmailClient()
	if err != nil {
		log.Fatal("Error:", err)
	}

	//get messages
	pageToken := ""
	for {

		var req *gmail.UsersMessagesListCall
		lastDate := getLastDate()
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
			log.Fatalf("Unable to retrieve messages: %v", err)
		}

		log.Printf("--------------")
		log.Printf("Processing %v messages...\n", len(r.Messages))
		for _, m := range r.Messages {

			msg, err := svc.Users.Messages.Get("me", m.Id).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve message %v: %v", m.Id, err)
			}

			internalDate, err := msToTime(strconv.FormatInt(msg.InternalDate, 10))
			if err != nil {
				log.Fatalln("Unable to parse message date", err)
			}

			//message date
			fmt.Println(internalDate)

			for _, h := range msg.Payload.Headers {

				//prints all header values
				//fmt.Println(h.Name + ":" + h.Value)

				if h.Name == "From" {

					fmt.Println("From:" + h.Value)
					parseAndSave(h.Value)

				} else if h.Name == "To" {

					fmt.Println("To:" + h.Value)
					parseAndSave(h.Value)

				} else if h.Name == "Cc" {

					fmt.Println("Cc:" + h.Value)
					parseAndSave(h.Value)
				} else if h.Name == "Subject" {
					fmt.Println("Subject:" + h.Value)
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

	//set the last known date
	t := time.Now()
	currentDate := t.Format("2006/01/02")
	err = updateLastDateInMongo(currentDate)

}
