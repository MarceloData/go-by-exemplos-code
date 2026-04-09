package main

import (
	"fmt"
	"sync"
	"time"
)

type Operation struct {
	ID     string
	Type   string
	Amount float64
	Date   time.Time
}

type Account struct {
	ID      string
	Balance float64
	Hist    []Operation
	mu      sync.RWMutex
}

func (a *Account) Deposit(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Balance += amount
	a.Hist = append(a.Hist, Operation{
		ID:     fmt.Sprintf("%d", len(a.Hist)+1),
		Type:   "deposit",
		Amount: amount,
		Date:   time.Now(),
	})
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	a.Balance -= amount
	a.Hist = append(a.Hist, Operation{
		ID:     fmt.Sprintf("%d", len(a.Hist)+1),
		Type:   "withdraw",
		Amount: amount,
		Date:   time.Now(),
	})
	return nil
}

func (a *Account) GetBalance() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Balance
}

func (a *Account) GetHistory() []Operation {
	a.mu.RLock()
	defer a.mu.RUnlock()

	h := make([]Operation, len(a.Hist))
	copy(h, a.Hist)
	return h
}

func main() {
	var wg sync.WaitGroup
	account := &Account{ID: "12345", Balance: 100.0}

	fmt.Printf("Initial Balance: %.2f\n", account.GetBalance())

	wg.Add(5)

	go func() {
		defer wg.Done()
		if err := account.Deposit(50.0); err == nil {
			fmt.Printf("Deposit successful. New Balance: %.2f\n", account.GetBalance())
		} else {
			fmt.Printf("Deposit failed: %v\n", err)
		}
		fmt.Printf("Current Balance: %.2f\n", account.GetBalance())
	}()

	go func() {
		defer wg.Done()
		if err := account.Withdraw(40.0); err == nil {
			fmt.Printf("Withdrawal successful. New Balance: %.2f\n", account.GetBalance())
		} else {
			fmt.Printf("Withdrawal failed: %v\n", err)
		}
		fmt.Printf("Current Balance: %.2f\n", account.GetBalance())
	}()

	go func() {
		defer wg.Done()
		if err := account.Deposit(60.0); err == nil {
			fmt.Printf("Deposit successful. New Balance: %.2f\n", account.GetBalance())
		} else {
			fmt.Printf("Deposit failed: %v\n", err)
		}
		fmt.Printf("Current Balance: %.2f\n", account.GetBalance())
	}()

	go func() {
		defer wg.Done()
		if err := account.Withdraw(30.0); err == nil {
			fmt.Printf("Withdrawal successful. New Balance: %.2f\n", account.GetBalance())
		} else {
			fmt.Printf("Withdrawal failed: %v\n", err)
		}
		fmt.Printf("Current Balance: %.2f\n", account.GetBalance())
	}()

	go func() {
		defer wg.Done()
		if err := account.Withdraw(150.0); err == nil {
			fmt.Printf("Withdrawal successful. New Balance: %.2f\n", account.GetBalance())
		} else {
			fmt.Printf("Withdrawal failed: %v\n", err)
		}
		fmt.Printf("Current Balance: %.2f\n", account.GetBalance())
	}()

	wg.Wait()

	fmt.Println("\nHistory:")
	for _, op := range account.GetHistory() {
		fmt.Printf("ID: %s, Operation: %s, Amount: %.2f, Date: %s\n",
			op.ID, op.Type, op.Amount, op.Date.Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("\nFinal Balance: %.2f\n", account.GetBalance())
}
