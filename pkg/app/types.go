package app

import (
	"database/sql"
)

type Model interface {
	User | Customer
}

type Repo[T any] interface {
	Read() T
	Write(T) int
}

type User struct {
	id int
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo() UserRepo {
	return UserRepo{db: nil}
}

func (r UserRepo) Read() User {
	return User{id: 1}
}

func (r UserRepo) Write(u User) int {
	return u.id
}

type Customer struct {
	id int
}

type CustomerRepo struct {
	db *sql.DB
}

func NewCustomerRepo() CustomerRepo {
	return CustomerRepo{db: nil}
}

func (r CustomerRepo) Read() Customer {
	return Customer{id: 1}
}

func (r CustomerRepo) Write(c Customer) int {
	return c.id

}

//var ur Repo[User] = NewUserRepo()
//var cr Repo[Customer] = NewCustomerRepo()

var Repos = make(map[string]Repo[any])
