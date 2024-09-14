package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"main.go/internal"
)

var (
	tmpl          *template.Template
	artists       []internal.Artist
	locations     internal.Location
	artistDetails internal.ArtistDetails
	relations     map[int]internal.Relation
	combinedData  internal.CombinedData
)

func main() {
	fetchInitialData()
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", Home)
	http.HandleFunc("/details", Details)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func fetchInitialData() {
	var err error
	artists, err = internal.FetchArtists()
	if err != nil {
		fmt.Println("Error fetching artists:", err)
		return
	}

	locations, err = internal.FetchLocations()
	if err != nil {
		fmt.Println("Error fetching locations:", err)
		return
	}

	// Fetch relations for all artists
	relations = make(map[int]internal.Relation)
	for _, artist := range artists {
		relation, err := internal.FetchRelation(strconv.Itoa(artist.ID))
		if err != nil {
			fmt.Println("Error fetching relation for artist:", artist.ID, err)
			continue
		}
		relations[artist.ID] = relation
	}
	// fmt.Println(relations)
}
func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/" {
		var err error
		artists, err = internal.FetchArtists()
		if err != nil {
			fmt.Println("Error handling request:", err)
			return
		}

		locations, err = internal.FetchLocations()
		if err != nil {
			fmt.Println("Error handling request:", err)
			return
		}

		combinedData.Artists = artists
		combinedData.Locations = locations

	} else if r.Method == http.MethodGet && r.URL.Path == "/search" {
		Search(w, r)
		return
	}
	tmpl.ExecuteTemplate(w, "index.html", combinedData)
}

func Details(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var selectedArtist internal.Artist
	for _, artist := range artists {
		if strconv.Itoa(artist.ID) == id {
			selectedArtist = artist
			break
		}
	}

	relation := relations[selectedArtist.ID]
	artistDetails.RetArtist = selectedArtist
	artistDetails.RetRelation = relation

	tmpl.ExecuteTemplate(w, "details.html", artistDetails)
}

func Search(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))

	var filteredArtists []internal.Artist

	for _, artist := range artists {
		creationDateString := strconv.Itoa(artist.CreationDate)

		matched := false
		if strings.Contains(strings.ToLower(artist.Name), query) ||
			strings.Contains(strings.ToLower(artist.FirstAlbum), query) ||
			strings.Contains(strings.ToLower(creationDateString), query) {
			matched = true
		}

		if !matched {
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), query) {
					matched = true
					break
				}
			}
		}

		// Retrieve relation data from pre-fetched data
		relation := relations[artist.ID]
		for location := range relation.DatesLocations {
			if strings.Contains(strings.ToLower(location), query) {
				matched = true
				break
			}
		}

		if matched {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	combinedData.Artists = filteredArtists

	tmpl.ExecuteTemplate(w, "index.html", combinedData)
}
