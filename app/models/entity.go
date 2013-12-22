package models

// the strcuct stored in db
type User struct {
	Email    string
	Nickname string
	Password []byte
}

// user post to register, and then -> User
type MockUser struct {
	LoginUser
	Nickname        string
	ConfirmPassword string
}

// user post to log in
type LoginUser struct {
	Email    string
	Password string
}
