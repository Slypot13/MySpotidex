package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"spotidex/models"
	"strings"
	"time"
)

// Variables globales (c'est plus simple)
var token string
var tokenExpiration time.Time
var client_http *http.Client

const (
	AuthURL      = "https://accounts.spotify.com/api/token"
	BaseAPIURL   = "https://api.spotify.com/v1"
	ClientID     = "b2758b6fc111451ea08499f71d2ec221"
	ClientSecret = "1a6559f3d05c43359f63a4ede6cd1b8e"
)

// On initialise le client
func InitService() {
	client_http = &http.Client{Timeout: 10 * time.Second}
}

// Fonction pour avoir le token
func GetToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	// On prépare la requete
	req, err := http.NewRequest("POST", AuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Encodage en base64 pour l'auth
	authHeader := base64.StdEncoding.EncodeToString([]byte(ClientID + ":" + ClientSecret))
	req.Header.Add("Authorization", "Basic "+authHeader)

	resp, err := client_http.Do(req)
	if err != nil {
		fmt.Println("Erreur requete token:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Erreur status code token:", resp.StatusCode)
		return fmt.Errorf("Erreur token")
	}

	var resultat models.TokenResponse
	json.NewDecoder(resp.Body).Decode(&resultat)

	token = resultat.AccessToken
	// On ajoute le temps d'expiration
	tokenExpiration = time.Now().Add(time.Duration(resultat.ExpiresIn) * time.Second)

	return nil
}

// Verif si le token est bon
func CheckTokenValid() {
	if token == "" || time.Now().After(tokenExpiration) {
		log.Println("Le token est vide ou expiré, on le refait")
		GetToken()
	}
}

// Recherche des artistes
func RechercheArtists(query string, offset int) (*models.SearchResponse, error) {
	CheckTokenValid()

	// On construit l'url
	url_search := fmt.Sprintf("%s/search?q=%s&type=artist&limit=20&offset=%d", BaseAPIURL, url.QueryEscape(query), offset)

	req, _ := http.NewRequest("GET", url_search, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client_http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Erreur décodage JSON search:", err)
		return nil, err
	}

	return &result, nil
}

// Récupérer un artiste par son ID
func GetArtiste(id string) *models.Artist {
	CheckTokenValid()

	url_artiste := fmt.Sprintf("%s/artists/%s", BaseAPIURL, id)
	req, _ := http.NewRequest("GET", url_artiste, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client_http.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	var artiste models.Artist
	json.NewDecoder(resp.Body).Decode(&artiste)

	return &artiste
}

// Récupérer les albums
func GetAlbumsArtiste(id string) []models.Album {
	CheckTokenValid()

	url_albums := fmt.Sprintf("%s/artists/%s/albums?limit=10", BaseAPIURL, id)
	req, _ := http.NewRequest("GET", url_albums, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client_http.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Structure temporaire juste pour ici
	var reponse struct {
		Items []models.Album `json:"items"`
	}
	json.NewDecoder(resp.Body).Decode(&reponse)

	return reponse.Items
}
