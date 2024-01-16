package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// Buat limiter dengan laju 50 request per detik
	limiter := rate.NewLimiter(rate.Every(time.Millisecond*20), 1)

	// Buat WaitGroup untuk menunggu selesainya semua permintaan
	var wg sync.WaitGroup

	// Loop untuk membuat goroutine sebanyak 50 permintaan
	for i := 1; i <= 50; i++ {
		wg.Add(1)
		go func() {
			// Tunggu hingga laju permintaan terpenuhi
			limiter.Wait(context.Background())

			// Buat permintaan ke API
			err := makeRequest()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Request succeeded!")
			}

			// Kurangi WaitGroup ketika permintaan selesai
			wg.Done()
		}()
	}

	// Tunggu hingga semua permintaan selesai
	wg.Wait()
}

func makeRequest() error {
	url := "http://localhost:8080/api/books?subject=love"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
