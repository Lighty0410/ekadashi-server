package database

import (
	"log"

	"github.com/Lighty0410/microservice-test/pkg/handleserver"
	mgo "gopkg.in/mgo.v2"
)

func connect() {
	users := &handleserver.AuthName{}
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("users")
	//err = c.Insert(users{username, password}) ??
}
