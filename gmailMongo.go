package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Contact object is saved into DB
type Contact struct {
	Name  string
	Email string
}

//Settings represnts a setings table
type Settings struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	LastSaved string        `bson:"LastSaved"`
	Account   string        `bson:"Account"`
}

const (
	mongoServer             = "localhost:27017"
	mongoDB                 = "gmailContacts"
	mongoContactsCollection = "contacts"
	mongoSettingsCollection = "settings"
)

const (
	accountName = "vcollak@gmail.com"
)

func getLastDateFromMongo() (string, error) {

	session, err := mgo.Dial(mongoServer)
	if err != nil {
		return "", err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	//	session.SetMode(mgo.Monotonic, true)

	//get a session
	c := session.DB(mongoDB).C(mongoSettingsCollection)

	result := Settings{}
	err = c.Find(bson.M{"Account": accountName}).Select(bson.M{"LastSaved": 1}).One(&result)
	if err != nil {
		return "", err
	} else {
		return result.LastSaved, nil

	}

}

func updateLastDateInMongo(lastSaved string) error {

	session, err := mgo.Dial(mongoServer)
	if err != nil {
		return err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	//	session.SetMode(mgo.Monotonic, true)

	//get a session
	c := session.DB(mongoDB).C(mongoSettingsCollection)

	// Update
	colQuerier := bson.M{"Account": accountName}
	change := bson.M{"$set": bson.M{"LastSaved": lastSaved}}
	err = c.Update(colQuerier, change)
	if err != nil {
		return err
	}

	return nil

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
	result := Contact{}
	err = c.Find(bson.M{"email": email}).One(&result)
	if err != nil {

		//insert the contact

		err = c.Insert(&Contact{name, email})
		if err != nil {
			return err
		}

	} else {
		//do nothing. the contact already exists
	}

	return nil

}
