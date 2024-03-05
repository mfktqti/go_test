package main

import "fmt"

type User struct {
	Name   string
	Age    int
	Gender string
}

func NewUser(name string, opts ...UserOption) *User {
	user := &User{
		Name: name,
	}
	for _, opt := range opts {
		opt(user)
	}
	return user
}

func WithAge(age int) UserOption {
	return func(u *User) {
		u.Age = age
	}
}

func WithGender(gender string) UserOption {
	return func(u *User) {
		u.Gender = gender
	}
}

type UserOption func(*User)

func main() {
	u := NewUser("pig", WithAge(12), WithGender("xx"))
	fmt.Printf("u: %#v\n", u)
	fmt.Printf("u: %+v\n", u)
}
