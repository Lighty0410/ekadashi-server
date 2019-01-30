package mongo

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

// GetMongoSession establish connection with MongoDB.
func GetMongoSession() *mgo.Session {
	mgoSession, err := mgo.Dial("mongodb://localhost:27017")
	mgoSession.SetMode(mgo.Monotonic, true)
	if err != nil {
		log.Fatal("unable to start mongo session")
	}
	return mgoSession.Clone()
}

//CreateUser creates users and insert it in MongoDB
func CreateUser(u *Users) error {
	session := GetMongoSession()

	c := session.DB("info").C("users")
	defer session.Close()

	err := c.Insert(u)
	if err != nil {
		log.Fatal("failed to insert user")
	}
	return err
}
