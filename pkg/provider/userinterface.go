package provider

// User is a struct that provides json information
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Handler is an interface
type Handler interface {
	Create(u *User) error
}
