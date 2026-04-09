package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type User struct {
	ID   string
	Name string
}

type Database []User

func GetUserByID(ctx context.Context, db Database, id string) (*User, error) {
	fmt.Println("querying database for user", id)

	select {
	case <-time.After(1 * time.Second):
		for i := range db {
			if db[i].ID == id {
				return &db[i], nil
			}
		}
		return nil, fmt.Errorf("user %s not found", id)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type call struct {
	done chan struct{}
	user *User
	err  error
}

type UserService struct {
	db       Database
	mu       sync.Mutex
	inflight map[string]*call
}

func NewUserService(db Database) *UserService {
	return &UserService{
		db:       db,
		inflight: make(map[string]*call),
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()

	if c, ok := s.inflight[id]; ok {
		s.mu.Unlock()

		select {
		case <-c.done:
			return c.user, c.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	c := &call{
		done: make(chan struct{}),
	}
	s.inflight[id] = c
	s.mu.Unlock()

	c.user, c.err = GetUserByID(ctx, s.db, id)

	s.mu.Lock()
	delete(s.inflight, id)
	s.mu.Unlock()

	close(c.done)
	return c.user, c.err
}

func main() {
	db := Database{
		{ID: "1", Name: "John Doe"},
		{ID: "2", Name: "Jane Doe"},
		{ID: "3", Name: "Alice"},
		{ID: "4", Name: "Bob"},
	}

	service := NewUserService(db)

	var wg sync.WaitGroup
	ids := []string{"1", "1", "1", "2", "4", "2", "1"}

	wg.Add(len(ids))

	for _, id := range ids {
		go func(id string) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			user, err := service.GetUser(ctx, id)
			if err != nil {
				fmt.Println("error:", err)
				return
			}

			fmt.Printf("got user: id=%s name=%s\n", user.ID, user.Name)
		}(id)
	}

	wg.Wait()
}
