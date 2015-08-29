package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Contact object is saved into DB
type contact struct {
	Name  string
	Email string
}

//Settings represnts a setings table
type settings struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	LastSaved string        `bson:"LastSaved"`
	Account   string        `bson:"Account"`
}

type MongoDB struct {
	/*
		mongoServer             = "localhost:27017"
		mongoDB                 = "gmailContacts"
		mongoContactsCollection = "contacts"
		mongoSettingsCollection = "settings"
		accountName = "vcollak@gmail.com"
	*/

	server             string
	db                 string
	accountName        string
	contactsCollection string
	settingsCollection string
}

func NewMongo(server string, db string, accountName string) *MongoDB {

	s := MongoDB{
		server:             server,
		db:                 db,
		accountName:        accountName,
		contactsCollection: "contacts",
		settingsCollection: "settings",
	}

	return &s

}

//LastDate retrieves the last message date that was processed
func (m *MongoDB) LastDate() (string, error) {

	session, err := mgo.Dial(m.server)
	if err != nil {
		return "", err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//get a session
	c := session.DB(m.db).C(m.settingsCollection)

	result := settings{}
	err = c.Find(bson.M{"Account": m.accountName}).Select(bson.M{"LastSaved": 1}).One(&result)
	if err != nil {
		return "", err
	} else {
		return result.LastSaved, nil

	}

}

//SetLastDate sets the last message date that was processed
func (m *MongoDB) SetLastDate(lastSaved string) error {

	session, err := mgo.Dial(m.server)
	if err != nil {
		return err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//get a session
	c := session.DB(m.db).C(m.settingsCollection)

	// Update
	colQuerier := bson.M{"Account": m.accountName}
	change := bson.M{"$set": bson.M{"LastSaved": lastSaved}}
	err = c.Update(colQuerier, change)
	if err != nil {
		return err
	}

	return nil

}

//SetContact saves contact in DB
func (m *MongoDB) SetContact(name string, email string) error {
	session, err := mgo.Dial(m.server)
	if err != nil {
		return err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//get a session
	c := session.DB(m.db).C(m.contactsCollection)

	//find the contact
	result := contact{}
	err = c.Find(bson.M{"email": email}).One(&result)
	if err != nil {

		//insert the contact

		err = c.Insert(&contact{name, email})
		if err != nil {
			return err
		}

	} else {
		//do nothing. the contact already exists
	}

	return nil

}
