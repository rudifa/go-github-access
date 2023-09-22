package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// type Repository struct {
//     Name        string `json:"name"`
//     Description string `json:"description"`
//     URL         string `json:"html_url"`
// }

func main4() {
	// Set up the HTTP client
	client := &http.Client{}

	// Set the initial page number and per-page limit
	page := 1
	perPage := 100

	// Loop over all pages of the repository list
	for {
		// Create the request
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/rudifa/repos?page=%d&per_page=%d", page, perPage), nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set the User-Agent header
		req.Header.Set("User-Agent", "my-app")
		// Set the Authorization header with your access token
		// token := getEnvVar("GITHUB_TOKEN")
		// fmt.Println("token:", token)
		// req.Header.Set("Authorization", "Bearer "+token)

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error: unexpected response status:", resp.Status)
			return
		}

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

		// Check if there are more pages
		linkHeader := resp.Header.Get("Link")
		if linkHeader == "" {
			break
		}
		links := parseLinkHeader(linkHeader)
		if links["next"] == "" {
			break
		}

		// Update the page number
		page++
	}
}

func parseLinkHeader(linkHeader string) map[string]string {
	links := make(map[string]string)
	for _, link := range strings.Split(linkHeader, ",") {
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) < 2 {
			continue
		}
		url := strings.Trim(parts[0], "<>")
		rel := strings.Trim(parts[1], " ")
		if rel == "rel=\"next\"" || rel == "rel=\"prev\"" || rel == "rel=\"first\"" || rel == "rel=\"last\"" {
			links[rel[5:len(rel)-1]] = url
		}
	}
	return links
}

func getEnvVar(varname string) string {
	// Load the environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the value of the environment variable
	value := os.Getenv(varname)

	// Return the environment variable value
	return value
}
