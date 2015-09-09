package main

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

//Contact object is saved into DB
type Contact struct {
	Name  string
	Email string
}

//Settings represnts a setings table
type Setting struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Saved   string
	Account string
}

type MongoDB struct {
	server         string
	db             string
	accountName    string
	collectionName string
	session        *mgo.Session
	collection     *mgo.Collection
}

func (m *MongoDB) NewMongo(server string, db string, accountName string, collectionName string) error {

	m.server = server
	m.db = db
	m.accountName = accountName
	m.collectionName = collectionName

	err := errors.New("")
	m.session, err = mgo.Dial(m.server)
	if err != nil {
		return err
	} else {
		m.collection = m.session.DB(m.db).C(m.collectionName)
	}

	return nil

}

//LastDate retrieves the last message date that was processed
func (m *MongoDB) LastDate() (string, error) {

	result := Setting{}
	err := m.collection.Find(bson.M{"account": "vcollak@gmail.com"}).One(&result)
	if err != nil {
		log.Println("Did not find last saved for:", m.accountName)
		return "", err
	} else {
		log.Printf("Found last saved for %s: %s", result.Account, result.Saved)
		return result.Saved, nil

	}

}

//SetLastDate sets the last message date that was processed
func (m *MongoDB) SetLastDate(lastSaved string) error {

	// Update
	colQuerier := bson.M{"account": m.accountName}
	change := bson.M{"$set": bson.M{"saved": lastSaved}}
	err := m.collection.Update(colQuerier, change)
	if err != nil {
		return err
	}

	return nil

}

//SetContact saves contact in DB
func (m *MongoDB) SetContact(name string, email string) error {

	//find the contact
	result := Contact{}
	err := m.collection.Find(bson.M{"email": email}).One(&result)
	log.Println("Found contact:", result)
	if err != nil {

		//insert the contact
		contact := &Contact{}
		contact.Name = name
		contact.Email = email

		err = m.collection.Insert(contact)
		if err != nil {
			return err
		}

	} else {
		log.Printf("Contact %s already exists", email)
		return err
	}

	return nil

}
