package main

import (
	"fmt"

	pg "gopkg.in/pg.v4"
)

type User struct {
	ID     int64
	Name   string
	Emails []string
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.ID, u.Name, u.Emails)
}

type Story struct {
	ID       int64
	Title    string
	AuthorID int64
	Author   *User
}

func (s Story) String() string {
	return fmt.Sprintf("Story<%d %s %s>", s.ID, s.Title, s.Author)
}

func createSchema(db *pg.DB) error {
	queries := []string{
		`CREATE TEMP TABLE users (id serial, name text, emails jsonb)`,
		`CREATE TEMP TABLE stories (id serial, title text, author_id bigint)`,
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "1123581321",
		Database: "foosio",
	})

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	user1 := &User{
		Name:   "admin",
		Emails: []string{"admin1@admin", "admin2@admin"},
	}
	err = db.Create(user1)
	if err != nil {
		panic(err)
	}

	err = db.Create(&User{
		Name:   "root",
		Emails: []string{"root1@root", "root2@root"},
	})
	if err != nil {
		panic(err)
	}

	story1 := &Story{
		Title:    "Cool story",
		AuthorID: user1.ID,
	}
	err = db.Create(story1)
	if err != nil {
		panic(err)
	}

	// Select user by primary key.
	user := User{ID: user1.ID}
	err = db.Select(&user)
	if err != nil {
		panic(err)
	}

	// Select all users.
	var users []User
	err = db.Model(&users).Select()
	if err != nil {
		panic(err)
	}

	// Select story and associated author in one query.
	var story Story
	err = db.Model(&story).
		Column("story.*", "Author").
		Where("story.id = ?", story1.ID).
		Select()
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
	fmt.Println(users)
	fmt.Println(story)
	// Output: User<1 admin [admin1@admin admin2@admin]>
	// [User<1 admin [admin1@admin admin2@admin]> User<2 root [root1@root root2@root]>]
	// Story<1 Cool story User<1 admin [admin1@admin admin2@admin]>>
}
