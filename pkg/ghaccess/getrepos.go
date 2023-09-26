package ghaccess

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rudifa/goutil/fjson"
	"github.com/rudifa/goutil/util"
)

// Mode is the output mode -----------------------------------------------
type Mode string

// Valid modes
const (
	String Mode = "string"
	JSON   Mode = "json"
	Data   Mode = "data"
)

// Modes is the list of valid Mode values
var Modes = []Mode{String, JSON, Data}

// ParseMode parses a string into a Mode value
func ParseMode(modeStr string) (Mode, error) {
    for _, mode := range Modes {
        if string(mode) == modeStr {
            return mode, nil
        }
    }
    return "", fmt.Errorf("invalid mode: %s", modeStr)
}

// ModeStrings returns the list of valid Mode values as strings
func ModeStrings() []string {
	return util.Map(Modes, func(m Mode) string {
		return string(m)
	})
}

// Repository represents a GitHub repository
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
}

// GetRepos gets the list of repositories for the given user and prints them
func GetRepos(user string, mode Mode) {
	// Set up the HTTP client
	client := &http.Client{}

	// Set the initial page number and per-page limit
	page := 1
	perPage := 100

	// Loop over all pages of the repository list
	for {
		// Create the request
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&per_page=%d", user, page, perPage), nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		log.Println("req:", req)

		// Set the User-Agent header
		req.Header.Set("User-Agent", "my-app")
		// Set the Authorization header with your access token
		token := getEnvVar("GITHUB_TOKEN")
		// fmt.Println("token:", token)
		req.Header.Set("Authorization", "Bearer "+token)

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

		switch mode {
		case String:
			err = printAsString(resp)
		case JSON:
			err = printAsJSON(resp)
		case Data:
			err = printAsData(resp)
		}
		if err != nil {
			return
		}

		// Check if there are more pages
		linkHeader := resp.Header.Get("Link")
		log.Println("linkHeader:", linkHeader)
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

func printAsData(resp *http.Response) error {
	var repos []Repository
	err := json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return err
	}

	for _, repo := range repos {
		fmt.Println("Name:", repo.Name)
		fmt.Println("Description:", repo.Description)
		fmt.Println("URL:", repo.URL)
		fmt.Println()
	}
	return nil
}

func printAsJSON(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return err
	}

	log.Println("len(body):", len(body))

	pretty, err := fjson.Prettyfmt(string(body))
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	fmt.Println(string(pretty))
	return nil
}

func printAsString(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return err
	}
	fmt.Println(string(body))
	return nil
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
