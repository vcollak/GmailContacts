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
	"github.com/vcollak/GmailContacts/db"
	"github.com/vcollak/GmailContacts/gmail"
	"log"
)

func main() {

	var knownEmails = []string{"vlad@collak.net", "vcollak@gmail.com", "vcollak@ignitedev.com",
		"info@slovacihouston.com", "vlad@openkloud.com", "vladimir.collak@ignitemediallc.com",
		"vladimir.collak@ignitemediahosting.com"}

	const (
		server      = "127.0.0.1"         //DB server address
		dbName      = "gmailContacts"     //DB name
		accountName = "vcollak@gmail.com" //the user's account name
	)

	log.Println("Starting...")

	//mongo DB
	var mongo = new(mongo.MongoDB)
	mongo.NewMongo(server, dbName, accountName)

	//gmail
	var gmail = new(gmail.Gmail)
	gmail.NewGmail(knownEmails, mongo)
	gmail.ProcessMessages()

}
