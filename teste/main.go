package main

import "fmt"

type Config struct {
	Port string
}

func (c Config) Print() {
	c.Port = "9000"
}

func NewConfig() *Config {
	cfg := Config{}
	return &cfg
}

func (c *Config) SetPort(p string) {
	c.Port = p
}

func main() {
	c1 := Config{
		Port: "8000",
	}
	fmt.Println(c1)
	c1.Print()
	fmt.Println(c1)
	c1.SetPort("9000")
	fmt.Println(c1)

	str := "hello"
	fmt.Println(str)
	bt := []byte(str)
	fmt.Println(string(bt))
}
