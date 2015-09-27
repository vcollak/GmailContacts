package mongo

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

//Contact object is saved into DB
type Contact struct {
	Name    string
	Email   string
	Account string
}

//Settings represnts a setings table
type Setting struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Saved   string
	Account string
}

type contactsCollection struct {
	collection     *mgo.Collection
	collectionName string
}

type settingsCollection struct {
	collection     *mgo.Collection
	collectionName string
}

type MongoDB struct {
	server      string
	db          string
	accountName string
	session     *mgo.Session
	contacts    contactsCollection
	settings    settingsCollection
}

//creates a new mongo DB connection
func (m *MongoDB) NewMongo(server string, db string, accountName string) error {

	m.server = server
	m.db = db
	m.accountName = accountName
	m.settings.collectionName = "settings"
	m.contacts.collectionName = "contacts"

	err := errors.New("")
	m.session, err = mgo.Dial(m.server)
	if err != nil {
		return err
	} else {
		m.settings.collection = m.session.DB(m.db).C(m.settings.collectionName)
		m.contacts.collection = m.session.DB(m.db).C(m.contacts.collectionName)

	}

	return nil

}

//LastDate retrieves the last message date that was processed
func (m *MongoDB) LastDate() (string, error) {

	result := Setting{}
	err := m.settings.collection.Find(bson.M{"account": m.accountName}).One(&result)
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
	err := m.settings.collection.Update(colQuerier, change)
	if err != nil {

		//unable to update. let's insert
		setting := &Setting{}
		setting.Account = m.accountName
		setting.Saved = lastSaved

		err := errors.New("")
		err = m.settings.collection.Insert(setting)
		if err != nil {
			return err
		}

	}

	return nil

}

//SetContact saves contact in DB
func (m *MongoDB) SetContact(name string, email string) error {

	//find the contact
	result := Contact{}
	err := m.contacts.collection.Find(bson.M{"email": email}).One(&result)

	if err != nil {

		//insert the contact
		contact := &Contact{}
		contact.Name = name
		contact.Email = email
		contact.Account = m.accountName

		err := errors.New("")
		err = m.contacts.collection.Insert(contact)
		if err != nil {
			return err
		}

	} else {
		log.Printf("Contact %s already exists", email)
		return err
	}

	return nil

}

//close the db session
func (m *MongoDB) Close() {
	m.session.Close()
}
