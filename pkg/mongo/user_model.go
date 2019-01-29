package mongo

import (
	"github.com/Lighty0410/ekadashi-server/pkg/provider"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//UserModel is a struct for MongoDB
type UserModel struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Username string
	Password string
}

//UserCollection is a collection for MongoDB
type UserCollection struct {
	collection *mgo.Collection
}

//NewUserModel transform User to UserModel
func NewUserModel(u provider.User) *UserModel {
	return &UserModel{
		Username: u.Username,
		Password: u.Password,
	}
}

//Create insert and information to mgo collection
func (col *UserCollection) Create(u provider.User) error {
	user := NewUserModel(u)
	return col.collection.Insert(&user)
}
