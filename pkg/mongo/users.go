package mongo

// Users is a struct that contains user's parameters for MongoDB.
type Users struct {
	Username string `json:"firstname"`
	Password string `json:"lastname"`
}
