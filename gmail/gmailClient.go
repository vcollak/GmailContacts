/*

App that connects to Gmail via Gmail api and saves all emails from "To", "From", and "Cc" into a MongoDB

Resources:
https://developers.google.com/gmail/api/quickstart/go
https://console.developers.google.com
https://godoc.org/google.golang.org/api/gmail/v1
https://tools.ietf.org/html/rfc4021


*/
package gmail

import (
	"errors"
	"fmt"
	"github.com/vcollak/GmailContacts/db"
	"github.com/vcollak/GmailContacts/utils"
	"golang.org/x/net/context"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/mail"
	"strconv"
	"strings"
)

type Gmail struct {
	knownEmails []string
	db          *mongo.MongoDB
	svc         *gmail.Service
}

//creates a new gmail connection
func (g *Gmail) NewGmail(knownEmails []string, db *mongo.MongoDB) error {

	g.knownEmails = knownEmails
	g.db = db

	err := errors.New("")
	g.svc, err = g.getGmailClient()
	if err != nil {
		log.Fatal("Unable to access Gmail. Error:", err)
		return err
	} else {
		return nil
	}

}

func (g *Gmail) getGmailClient() (*gmail.Service, error) {
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
func (g *Gmail) isKnownEmail(email string) bool {

	for _, e := range g.knownEmails {

		if strings.ToUpper(email) == strings.ToUpper(e) {
			return true
		}
	}

	return false
}

func (g *Gmail) saveHeaderFields(headerValue string) {

	emails, err := mail.ParseAddressList(headerValue)
	if err != nil {
		log.Println("Unable to parse:", headerValue)
	} else {

		for _, v := range emails {

			name := v.Name
			email := v.Address

			if !g.isKnownEmail(email) {

				err := errors.New("")
				err = g.db.SetContact(name, email)

				if err != nil {
					log.Println("Unable to save email:", email)
				}
			} else {
				log.Println("Known email. Ignoring:", email)
			}
		}
	}
}

func (g *Gmail) ProcessMessages() {

	//close the sessions at the end
	defer g.db.Close()

	//get messages
	pageToken := ""
	firstMessage := true

	for {

		var req *gmail.UsersMessagesListCall
		lastDate, err := g.db.LastDate()

		if lastDate == "" {
			log.Println("Retrieving all messages.")
			req = g.svc.Users.Messages.List("me")

		} else {
			log.Println("Retrieving messages starting on", lastDate)
			req = g.svc.Users.Messages.List("me").Q("after: " + lastDate)
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

			msg, err := g.svc.Users.Messages.Get("me", m.Id).Do()
			if err != nil {
				log.Printf("Unable to retrieve message %v: %v", m.Id, err)
				continue
			}

			lastMessageRetrievedDate, err := utils.MsToTime(strconv.FormatInt(msg.InternalDate, 10))
			if err != nil {
				log.Println("Unable to parse message date", err)
			}

			//message date
			log.Println(lastMessageRetrievedDate)

			if firstMessage {

				//set the last known date
				currentDate := lastMessageRetrievedDate.Format("2006/01/02")
				err = g.db.SetLastDate(currentDate)
				if err != nil {
					log.Println("Unable to save last message date:", currentDate)
				} else {
					log.Println("Saved last message date:", currentDate)
					firstMessage = false
				}

			}

			for _, h := range msg.Payload.Headers {

				//prints all header values
				//fmt.Println(h.Name + ":" + h.Value)

				if h.Name == "From" {

					log.Println("From:" + h.Value)
					g.saveHeaderFields(h.Value)

				} else if h.Name == "To" {

					log.Println("To:" + h.Value)
					g.saveHeaderFields(h.Value)

				} else if h.Name == "Cc" {

					log.Println("Cc:" + h.Value)
					g.saveHeaderFields(h.Value)

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
