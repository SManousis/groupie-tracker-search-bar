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
	tmpl, err := template.ParseFiles("templates/index.html")
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

func search(input string) []string {
	input = strings.ToLower(input)
	var list []string
	var Bands []Artist
	err := FetchArtists(&Bands)
	if err != nil {
		fmt.Println("Error fetching artists:", err)
	}
	for _, v := range Bands {
		if strings.Contains(strings.ToLower(v.Name), (input)) {
			list = append(list, v.Name+" - artist/band")

		}
		if strings.Contains(strconv.Itoa(v.CreationDate), (input)) {
			list = append(list, strconv.Itoa(v.CreationDate)+" - creation dated")
		}
		if strings.Contains(strings.ToLower(v.FirstAlbum), (input)) {
			list = append(list, v.FirstAlbum+" - first album date")
		}
		for _, member := range v.Members {
			if strings.Contains(strings.ToLower(member), input) {
				list = append(list, member+" - member")
			}
		}
		for _, location := range v.Locations {
			if strings.Contains(strings.ToLower(location), input) {
				list = append(list, location+" - member")
			}
		}
	}
	return list
}
