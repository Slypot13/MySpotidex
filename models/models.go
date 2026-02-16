package models

type Artist struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Images     []Image  `json:"images"`
	Genres     []string `json:"genres"`
	Popularity int      `json:"popularity"`
	Followers  struct {
		Total int `json:"total"`
	} `json:"followers"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Album struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Images      []Image `json:"images"`
	ReleaseDate string  `json:"release_date"`
	TotalTracks int     `json:"total_tracks"`
}

type Track struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SearchResponse struct {
	Artists struct {
		Items []Artist `json:"items"`
		Total int      `json:"total"`
	} `json:"artists"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
