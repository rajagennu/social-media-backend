package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	path string
}

func NewClient(path string) Client {
	return Client{path: path}
}

type databaseSchema struct {
	Users map[string]User `json:"users"`
	Posts map[string]Post `json:"posts"`
}

type User struct {
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
}

type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserEmail string    `json:"userEmail"`
	Text      string    `json:"text"`
}

func (c Client) createDB() error {
	data, err := json.Marshal(databaseSchema{
		Users: make(map[string]User),
		Posts: make(map[string]Post),
	})

	if err != nil {
		return err
	}

	err = os.WriteFile(c.path, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) EnsureDB() error {
	_, err := os.ReadFile(c.path)
	if err != nil {
		log.Println("DB File not existed so creating one.")
		err = c.createDB()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Client) updateDB(db databaseSchema) error {
	data, err := json.Marshal(db)
	if err != nil {
		log.Fatal("Errors while JSON Marshaling ", err)
	}
	err = os.WriteFile(c.path, data, 0600)
	fmt.Println("Written content to file ", c.path)
	if err != nil {
		log.Fatal("Errors while writing to db path ", err)
	}
	return nil
}

func (c Client) readDB() (databaseSchema, error) {
	db := databaseSchema{
		Users: make(map[string]User),
		Posts: make(map[string]Post),
	}
	data, err := os.ReadFile(c.path)
	if err != nil {
		return db, err
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return db, err
	}
	return db, nil

}

func (c Client) CreateUser(email, password, name string, age int) (User, error) {
	u := User{
		CreatedAt: time.Now().UTC(),
		Email:     email,
		Password:  password,
		Name:      name,
		Age:       age,
	}

	err := c.EnsureDB()
	if err != nil {
		return u, err
	}

	currentData, err := c.readDB()
	if err != nil {
		return u, err
	}

	if _, ok := currentData.Users[email]; !ok {
		currentData.Users[email] = u
	}
	log.Println("Writing new user to db ", currentData)

	err = c.updateDB(currentData)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (c Client) UpdateUser(email, password, name string, age int) (User, error) {
	currentData, err := c.readDB()
	if err != nil {
		return User{}, err
	}
	if _, ok := currentData.Users[email]; ok {
		user := User{
			Email:    email,
			Password: password,
			Name:     name,
			Age:      age,
		}
		currentData.Users[email] = user

		err = c.updateDB(currentData)
		if err != nil {
			return user, err
		}
		return user, nil
	}
	return User{}, errors.New("user doesn't exits")

}

func (c Client) GetUser(email string) (User, error) {

	currentData, err := c.readDB()
	if err != nil {
		return User{}, err
	}
	if _, ok := currentData.Users[email]; ok {
		return currentData.Users[email], nil
	}
	return User{}, nil

}

func (c Client) DeleteUser(email string) error {
	currentData, err := c.readDB()
	if err != nil {
		return err
	}
	if _, ok := currentData.Users[email]; ok {
		delete(currentData.Users, email)
		err = c.updateDB(currentData)
		if err != nil {
			return err
		}
	}
	return nil

}

func (c Client) CreatePost(userEmail, text string) (Post, error) {
	currentData, err := c.readDB()
	if err != nil {
		return Post{}, err
	}
	uid := uuid.New().String()

	_, err = c.GetUser(userEmail)
	if err != nil {
		return Post{}, err
	}

	post := Post{
		ID:        uid,
		Text:      text,
		CreatedAt: time.Now().UTC(),
		UserEmail: userEmail,
	}

	currentData.Posts[uid] = post
	err = c.updateDB(currentData)
	return post, err

}

func (c Client) GetPosts(userEmail string) ([]Post, error) {

	var posts []Post
	currentData, err := c.readDB()
	if err != nil {
		return posts, err
	}

	for uid, data := range currentData.Posts {
		fmt.Println("uid ", uid, "data ", data)
		if data.UserEmail == userEmail {
			posts = append(posts, currentData.Posts[uid])
		}
	}
	return posts, nil
}

func (c Client) DeletePost(id string) error {

	currentData, err := c.readDB()
	if err != nil {
		return errors.New("Issue while reading data from readDB")
	}
	if _, ok := currentData.Posts[id]; ok {
		delete(currentData.Posts, id)
		c.updateDB(currentData)
		return nil
	}
	return nil

}
