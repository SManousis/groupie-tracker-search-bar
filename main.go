package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

type SearchResult struct {
	Matches string `json:"matches"`
}

type Artist struct {
	ID           int               `json:"id"`
	Image        string            `json:"image"`
	Name         string            `json:"name"`
	Members      []string          `json:"members"`
	CreationDate int               `json:"creationDate"`
	FirstAlbum   string            `json:"firstAlbum"`
	Locations    []string          `json:"-"`
	Dates        []string          `json:"-"`
	Relation     map[string]string `json:"-"`
}

const baseUrl = "https://groupietrackers.herokuapp.com/api/"

var endpoints = []string{"locations/", "dates/", "relation/"}

// create a var FetchArtist so that it points to the actual implementation when used in unit testing
var FetchArtists = fetchArtists

func fetchArtists(artists *[]Artist) error {
	response, err := http.Get(baseUrl + "artists")
	if err != nil {
		return err
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(artists)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for counter := range *artists {
		for _, endpoint := range endpoints {
			wg.Add(1)

			go func(artist *Artist, endpoint string) {
				defer wg.Done()
				if err := updateArtistStruct(artist, endpoint); err != nil {
					fmt.Println("Error updating artist:", err)
				}
			}(&(*artists)[counter], endpoint)

		}
	}
	wg.Wait()
	return nil
}

func updateArtistStruct(artist *Artist, endpoint string) error {
	response, err := http.Get(baseUrl + endpoint + strconv.Itoa(artist.ID))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var rawData map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&rawData); err != nil {
		return err
	}

	switch endpoint {
	case "locations/":
		if locations, exists := rawData["locations"].([]interface{}); exists {
			artist.Locations = make([]string, len(locations))
			for i, v := range locations {
				artist.Locations[i] = strings.ToTitle(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(v), "_", " "), "-", ", "))

			}
		}
	case "dates/":
		if dates, exists := rawData["dates"].([]interface{}); exists {
			artist.Dates = make([]string, len(dates))
			for i, v := range dates {
				date := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(v), "[", ""), "]", ""), " ", ", ")
				artist.Dates[i] = date
			}
		}
	case "relation/":
		if datesLocations, exists := rawData["datesLocations"].(map[string]interface{}); exists {
			artist.Relation = make(map[string]string)
			for k, v := range datesLocations {
				location := strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(k), "_", " "), "-", ", ")
				date := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(v), "[", ""), "]", ""), " ", ", ")
				artist.Relation[strings.ToTitle(location)] = date
			}
		}
	}

	return nil
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	// Fetch the artists and dynamically update them
	artists := []Artist{}
	err := FetchArtists(&artists)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching artists: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare the data to be passed to the template
	data := struct {
		Artists []Artist
	}{
		Artists: artists,
	}
	if r.URL.Path != "/" {
		http.ServeFile(w, r, "templates/error.html")
		return
	}
	// Prepare the template with embedded JavaScript
	tmpl, err := template.New("index.html").ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	// Execute the template with artist data
	//err = tmpl.Execute(w, data)
	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve the homepage on /
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/result", resultHome)
	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//just test

	// Start the server
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	query = strings.ToLower(query)

	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	var list []string
	var Bands []Artist

	// Fetch the list of artists
	err := FetchArtists(&Bands)
	if err != nil {
		fmt.Println("Error fetching artists:", err)
	}

	// Check if the search query matches any details
	for _, v := range Bands {
		// Search by artist name
		if strings.Contains(strings.ToLower(v.Name), query) {
			list = append(list, v.Name+" - artist/band")
		}

		// Search by band members
		for _, member := range v.Members {
			if strings.Contains(strings.ToLower(member), query) {
				list = append(list, member+" - member of "+v.Name)
			}
		}

		// Now check for location-based search separately
		for _, location := range v.Locations {
			if strings.Contains(strings.ToLower(location), query) {
				list = append(list, v.Name+" - artist/band in location: "+location)
			}
		}

		// Check if the creation date matches
		if strings.Contains(strconv.Itoa(v.CreationDate), query) {
			list = append(list, v.Name+" - artist/band created in "+strconv.Itoa(v.CreationDate))
		}

		// Check if the first album date matches
		if strings.Contains(strings.ToLower(v.FirstAlbum), query) {
			list = append(list, v.Name+" - artist/band first album released in "+strings.ToLower(v.FirstAlbum))
		}
	}

	// Return the list of suggestions as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func resultHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Extract the main query
	parts := strings.Split(data.Query, " - ")
	query := parts[0]
	relatedInfo := ""
	if len(parts) > 1 {
		relatedInfo = parts[1] // This will capture "member of <band>" or other info
	}
	queryLower := strings.ToLower(query)

	var Bands []Artist
	err := FetchArtists(&Bands)
	if err != nil {
		fmt.Println("Error fetching artists:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var matchedArtist *Artist
OuterLoop:
	for _, v := range Bands {
		// Check if the query matches the band or solo artist directly
		if strings.ToLower(v.Name) == queryLower ||
			strings.ToLower(strconv.Itoa(v.CreationDate)) == queryLower ||
			strings.ToLower(v.FirstAlbum) == queryLower {
			matchedArtist = &v
			break OuterLoop
		}

		// If the query matches a band member, ensure we match the correct band
		for _, member := range v.Members {
			if strings.ToLower(member) == queryLower {
				if relatedInfo != "" && strings.Contains(strings.ToLower(relatedInfo), strings.ToLower(v.Name)) {
					matchedArtist = &v
					break OuterLoop
				}
			}
		}

		// Check for exact matches for locations or other criteria
		for _, location := range v.Locations {
			if strings.ToLower(location) == queryLower {
				matchedArtist = &v
				break OuterLoop
			}
		}
	}

	// If no match found, return error
	if matchedArtist == nil {
		http.Error(w, "No artist found", http.StatusNotFound)
		return
	}

	// Return the matched artist's information as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matchedArtist)
}
