package main

const (
	server      = "127.0.0.1"         //DB server address
	dbName      = "gmailContacts"     //DB name
	accountName = "vcollak@gmail.com" //the user's account name
)

//ignore these emails when harvesting emails from gmail
var knownEmails = []string{"vlad@collak.net", "vcollak@gmail.com", "vcollak@ignitedev.com",
	"info@slovacihouston.com", "vlad@openkloud.com", "vladimir.collak@ignitemediallc.com",
	"vladimir.collak@ignitemediahosting.com"}
