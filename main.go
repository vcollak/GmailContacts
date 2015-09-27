/*

Copyright 2015 Vladimir Collak


Utility that will extract all contacts from Gmail emails. It takes From, To,
and Cc fields and saves them into a MongoDB database. Each time the utility
is executed it will scan only new email (the email the utility
has not processes yet) and add the contact as a new contact
(if it's not already in the DB).

*/

package main

import (
	"errors"
	"github.com/vcollak/GmailContacts/db"
	"github.com/vcollak/GmailContacts/gmail"
	"log"
)

func main() {

	var knownEmails = []string{"vlad@collak.net", "vcollak@gmail.com", "vcollak@ignitedev.com",
		"info@slovacihouston.com", "vlad@openkloud.com", "vladimir.collak@ignitemediallc.com",
		"vladimir.collak@ignitemediahosting.com"}

	const (
		server      = "127.0.0.1"       //DB server address
		dbName      = "GmailContactsA"  //DB name
		accountName = "vlad@collak.net" //the user's account name
	)

	log.Println("Starting...")

	//mongo DB
	var mongo = new(mongo.MongoDB)
	err := mongo.NewMongo(server, dbName, accountName)
	if err != nil {
		log.Printf("Unable to connect to DB. Server: %s  dbName: %s", server, dbName)
		log.Printf("Exiting...")
		return
	}

	//gmail
	err = errors.New("")
	var gmail = new(gmail.Gmail)
	err = gmail.NewGmail(knownEmails, mongo)
	if err != nil {
		log.Printf("Unable to connect to Gmail. Error:%s", err)

	}
	gmail.ProcessMessages()

}
