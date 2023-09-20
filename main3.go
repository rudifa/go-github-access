package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
}

// main3 returns the first page of repository list (30 items) for the user "rudifa" on GitHub.
// It uses an unauthenticated request, so the rate limit is 60 requests per hour.
func main3() {
	// Set up the HTTP client
	client := &http.Client{}

	// Create the request
	req, err := http.NewRequest("GET", "https://api.github.com/users/rudifa/repos", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", "my-app")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Parse the response
	var repos []Repository
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	// Print the repository information
	for _, repo := range repos {
		fmt.Println("Name:", repo.Name)
		fmt.Println("Description:", repo.Description)
		fmt.Println("URL:", repo.URL)
		fmt.Println()
	}
}
