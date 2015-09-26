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
	"log"
)

func main() {

	log.Println("Starting...")
	processMessages()

}
