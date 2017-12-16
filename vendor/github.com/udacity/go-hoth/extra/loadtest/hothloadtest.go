package main

// Basic load test example

import (
	"fmt"
	"sync"

	hoth "github.com/udacity/go-hoth"
)

func main() {
	fmt.Println("Go Load Tests")

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go loadtest(&wg, i)
	}
	wg.Wait()

}

func loadtest(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	// TODO: CLI args via flag package
	jwt, err := hoth.GetJWT("<UdacityStaffUserAccount>", "<Password")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Attempt #%v, Got: %v\n", id, jwt)

}
