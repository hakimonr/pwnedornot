package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	haveIBeenPwnedAPI = "https://haveibeenpwned.com/api/v3/breachedaccount/%s"
	apiKey            = "API_KEY" // Replace with your actual API key
	userAgent         = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0"
	outputFile        = "pwnedHesaplar.txt"
	sleepDuration     = 6 * time.Second
)

func main() {
	filePath := "emails.txt" // Replace with the path to your email list file

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for scanner.Scan() {
		email := strings.TrimSpace(scanner.Text())
		if email != "" {
			wg.Add(1)
			go func(email string) {
				defer wg.Done()
				pwned, err := checkPwned(email)
				if err != nil {
					fmt.Printf("Failed to check pwned status for email %s: %s\n", email, err.Error())
					return
				}

				if pwned {
					mu.Lock()
					f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						fmt.Printf("Failed to open output file: %s\n", err.Error())
						mu.Unlock()
						return
					}
					defer f.Close()

					fmt.Printf("Pwned!! : %s\n", email)
					f.WriteString(email + "\n")
					mu.Unlock()
				} else {
					fmt.Printf("Not Pwned!! : %s\n", email)
				}
			}(email)

			// Delay between each goroutine to avoid hitting rate limits
			time.Sleep(sleepDuration)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %s\n", err.Error())
	}

	wg.Wait()
}

func checkPwned(email string) (bool, error) {
	url := fmt.Sprintf(haveIBeenPwnedAPI, email)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Host", "haveibeenpwned.com")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Hibp-Api-Key", apiKey)

	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return true, nil
	} else if response.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, fmt.Errorf("unexpected response status: %s", response.Status)
	}
}
