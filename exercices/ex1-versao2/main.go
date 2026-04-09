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
	mu      sync.Mutex
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid deposit amount")
	}

	// simula processamento fora da região crítica
	time.Sleep(1 * time.Second)

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
	if amount <= 0 {
		return fmt.Errorf("invalid withdraw amount")
	}

	// simula processamento fora da região crítica
	time.Sleep(1 * time.Second)

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
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.Balance
}

func (a *Account) GetHistory() []Operation {
	a.mu.Lock()
	defer a.mu.Unlock()

	h := make([]Operation, len(a.Hist))
	copy(h, a.Hist)
	return h
}

func main() {
	account := &Account{
		ID:      "12345",
		Balance: 100.0,
	}

	fmt.Printf("Initial Balance: %.2f\n", account.GetBalance())

	var wg sync.WaitGroup

	operations := []func(){
		func() {
			defer wg.Done()
			if err := account.Deposit(50.0); err != nil {
				fmt.Printf("Deposit failed: %v\n", err)
				return
			}
			fmt.Printf("Deposit 50 successful. Current Balance: %.2f\n", account.GetBalance())
		},
		func() {
			defer wg.Done()
			if err := account.Withdraw(40.0); err != nil {
				fmt.Printf("Withdrawal failed: %v\n", err)
				return
			}
			fmt.Printf("Withdraw 40 successful. Current Balance: %.2f\n", account.GetBalance())
		},
		func() {
			defer wg.Done()
			if err := account.Deposit(60.0); err != nil {
				fmt.Printf("Deposit failed: %v\n", err)
				return
			}
			fmt.Printf("Deposit 60 successful. Current Balance: %.2f\n", account.GetBalance())
		},
		func() {
			defer wg.Done()
			if err := account.Withdraw(30.0); err != nil {
				fmt.Printf("Withdrawal failed: %v\n", err)
				return
			}
			fmt.Printf("Withdraw 30 successful. Current Balance: %.2f\n", account.GetBalance())
		},
		func() {
			defer wg.Done()
			if err := account.Withdraw(150.0); err != nil {
				fmt.Printf("Withdrawal failed: %v\n", err)
				return
			}
			fmt.Printf("Withdraw 150 successful. Current Balance: %.2f\n", account.GetBalance())
		},
	}

	wg.Add(len(operations))
	for _, op := range operations {
		go op()
	}

	wg.Wait()

	fmt.Println("\nFinal history:")
	for _, op := range account.GetHistory() {
		fmt.Printf("ID: %s, Operation: %s, Amount: %.2f, Date: %s\n",
			op.ID, op.Type, op.Amount, op.Date.Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("\nFinal Balance: %.2f\n", account.GetBalance())
}
