package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Email string
}

type LastDate struct {
	lastDate string
}

const (
	mongoServer             = "localhost:27017"
	mongoDB                 = "gmailContacts"
	mongoContactsCollection = "contacts"
	mongoSettingsCollection = "settings"
)

func getLastDateFromMongo() (string, error) {

	session, err := mgo.Dial(mongoServer)
	if err != nil {
		return "", err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//get a session
	c := session.DB(mongoDB).C(mongoSettingsCollection)

	result := ""
	err = c.Find(bson.M{"email": "vcollak@gmail.com"}).Select(bson.M{"lastDate": 0}).One(&result)
	if err != nil {
		return "", err
	}

	return result, nil

}

func saveLastDate(date string) error {

	session, err := mgo.Dial(mongoServer)
	if err != nil {
		return err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//get a session
	c := session.DB(mongoDB).C(mongoSettingsCollection)

	// Update
	colQuerier := bson.M{"account": "vcollak@gmail.com"}
	change := bson.M{"$set": bson.M{"lastDate": date}}
	err = c.Update(colQuerier, change)
	if err != nil {
		return err
	}

	return nil

}

func saveContact(name string, email string) error {
	session, err := mgo.Dial(mongoServer)
	if err != nil {
		return err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//get a session
	c := session.DB(mongoDB).C(mongoContactsCollection)

	//find the contact
	result := Person{}
	err = c.Find(bson.M{"email": email}).One(&result)
	if err != nil {

		//insert the contact

		err = c.Insert(&Person{name, email})
		if err != nil {
			return err
		}

	} else {
		//do nothing. the contact already exists
	}

	return nil

}
