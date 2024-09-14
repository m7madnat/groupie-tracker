package internal

type ArtistDetails struct {
	RetArtist   Artist
	RetRelation Relation
}

type CombinedData struct {
	Artists   []Artist
	Relation  Relation
	Locations Location
}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Location struct {
	Index []struct {
		Locations []string `json:"locations"`
	} `json:"index"`
}
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
