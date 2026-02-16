package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"spotidex/models"
	"strings"
	"time"
)

// Constant pour l'API Spotify
const (
	AuthURL             = "https://accounts.spotify.com/api/token"
	BaseAPIURL          = "https://api.spotify.com/v1"
	defaultClientID     = "b2758b6fc111451ea08499f71d2ec221"
	defaultClientSecret = "1a6559f3d05c43359f63a4ede6cd1b8e"
)

type SpotifyService struct {
	Client      *http.Client
	AccessToken string
	ExpiresAt   time.Time
}

func NewSpotifyService() *SpotifyService {
	return &SpotifyService{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Authenticate récupère le token d'accès
func (s *SpotifyService) Authenticate() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", AuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	auth := base64.StdEncoding.EncodeToString([]byte(defaultClientID + ":" + defaultClientSecret))
	req.Header.Add("Authorization", "Basic "+auth)

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		log.Printf("Erreur Auth Spotify (%d): %s - %s\n", resp.StatusCode, errorResponse.Error, errorResponse.ErrorDescription)
		return errors.New("échec de l'authentification spotify: " + errorResponse.ErrorDescription)
	}

	var tokenResp models.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	s.AccessToken = tokenResp.AccessToken
	s.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return nil
}

// CheckToken vérifie si le token est valide
func (s *SpotifyService) CheckToken() error {
	if s.AccessToken == "" {
		// Pas de token, on authentifie
		return s.Authenticate()
	}
	if time.Now().After(s.ExpiresAt) {
		// Token expiré, on recommence
		return s.Authenticate()
	}
	// Tout est bon
	return nil
}

// SearchArtists recherche des artistes
func (s *SpotifyService) SearchArtists(query string, offset int) (*models.SearchResponse, error) {
	if err := s.CheckToken(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/search?q=%s&type=artist&limit=20&offset=%d", BaseAPIURL, url.QueryEscape(query), offset)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+s.AccessToken)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur api: %d", resp.StatusCode)
	}

	var searchResp models.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}

// GetArtist récupère les détails d'un artiste
func (s *SpotifyService) GetArtist(id string) (*models.Artist, error) {
	if err := s.CheckToken(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/artists/%s", BaseAPIURL, id)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+s.AccessToken)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artist models.Artist
	if err := json.NewDecoder(resp.Body).Decode(&artist); err != nil {
		return nil, err
	}

	return &artist, nil
}

// GetArtistAlbums récupère les albums d'un artiste
func (s *SpotifyService) GetArtistAlbums(id string) ([]models.Album, error) {
	if err := s.CheckToken(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/artists/%s/albums?limit=10", BaseAPIURL, id)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+s.AccessToken)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var albumResp struct {
		Items []models.Album `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&albumResp); err != nil {
		return nil, err
	}

	return albumResp.Items, nil
}
