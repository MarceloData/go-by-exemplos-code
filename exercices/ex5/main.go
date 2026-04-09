package main

import (
	"fmt"
	"slices"
	"sync"
)

type User struct {
	ID   int
	Name string
}

type Database []User

func GetUserByID(db Database, id int) *User {
	for i := range db {
		if db[i].ID == id {
			return &db[i]
		}
	}
	return nil
}

func main() {
	db := Database{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Jane Doe"},
		{ID: 3, Name: "Alice"},
		{ID: 4, Name: "Bob"},
	}

	ch := make(chan int)
	bufferedCh := make(chan int)
	result := make(chan *User)

	var producersWG sync.WaitGroup

	list := []int{}

	go func() {
		for id := range ch {
			if !slices.Contains(list, id) {
				list = append(list, id)
			}
		}

		for _, id := range list {
			bufferedCh <- id
		}
		close(bufferedCh)
	}()

	go func() {
		for id := range bufferedCh {
			user := GetUserByID(db, id)
			result <- user
		}
		close(result)
	}()

	ids := []int{1, 1, 1, 2, 4}

	producersWG.Add(len(ids))
	for _, id := range ids {
		go func(id int) {
			defer producersWG.Done()
			ch <- id
		}(id)
	}

	producersWG.Wait()
	close(ch)

	for user := range result {
		if user != nil {
			fmt.Println("User found:", user.Name)
		} else {
			fmt.Println("User not found")
		}
	}
}
