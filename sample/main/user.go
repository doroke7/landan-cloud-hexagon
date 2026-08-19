package main

func NewUser(nAge int) *User {
	return &User{
		Age: nAge,
	}
}

type User struct {
	Age int
}
