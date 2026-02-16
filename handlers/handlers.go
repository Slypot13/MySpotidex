package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"spotidex/models"
	"spotidex/services"
	"strconv"
)

type PageData struct {
	Title     string
	Artists   []models.Artist
	Artist    *models.Artist
	Albums    []models.Album
	Query     string
	Offset    int
	Favorites []models.Artist
}

var spotifyService *services.SpotifyService

func Init(s *services.SpotifyService) {
	spotifyService = s
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Recherche par défaut pour afficher des artistes sur l'accueil
	resp, err := spotifyService.SearchArtists("genre:pop", 0)
	var artists []models.Artist
	if err == nil {
		artists = resp.Artists.Items
	} else {
		log.Println("Erreur chargement accueil:", err)
	}

	renderTemplate(w, "home", PageData{Title: "Accueil", Artists: artists})
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	genre := r.URL.Query().Get("genre")
	year := r.URL.Query().Get("year") // Année de début de carrière ou autre (artificiel pour l'exercice)

	offsetStr := r.URL.Query().Get("offset")
	offset, _ := strconv.Atoi(offsetStr)

	data := PageData{Title: "Recherche", Query: query, Offset: offset}

	if query != "" {
		searchQuery := query
		if genre != "" {
			searchQuery += " genre:\"" + genre + "\""
		}
		if year != "" {
			searchQuery += " year:" + year
		}

		resp, err := spotifyService.SearchArtists(searchQuery, offset)
		if err != nil {
			// Si ça plante, on affiche l'erreur dans la console
			log.Println("Erreur recherche:", err)
			// TODO: Faire une page d'erreur plus jolie un jour
			http.Error(w, "Erreur lors de la recherche", http.StatusInternalServerError)
			if err.Error() == "EOF" {
				log.Println("C'est bizarre, l'erreur est EOF")
			}
			return
		}
		data.Artists = resp.Artists.Items
	}

	renderTemplate(w, "search_results", data)
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}

	artist, err := spotifyService.GetArtist(id)
	if err != nil {
		http.Error(w, "Erreur récupération artiste", http.StatusInternalServerError)
		return
	}

	albums, err := spotifyService.GetArtistAlbums(id)
	if err != nil {
		// On continue même sans albums, c'est pas critique
	}

	data := PageData{
		Title:  artist.Name,
		Artist: artist,
		Albums: albums,
	}

	renderTemplate(w, "artist_detail", data)
}

func FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	favs, err := loadFavorites()
	if err != nil {
		favs = []models.Artist{}
	}
	renderTemplate(w, "favorites", PageData{Title: "Favoris", Favorites: favs})
}

func AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var artist models.Artist
	if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	favs, _ := loadFavorites()
	favs = append(favs, artist)
	saveFavorites(favs)

	w.WriteHeader(http.StatusOK)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	t, err := template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/"+tmpl+".html")
	if err != nil {
		http.Error(w, "Erreur template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "layout", data)
}

func loadFavorites() ([]models.Artist, error) {
	file, err := os.Open("data/favorites.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var favs []models.Artist
	json.NewDecoder(file).Decode(&favs)
	return favs, nil
}

func saveFavorites(favs []models.Artist) {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Println("Erreur création dossier data:", err)
		return
	}
	file, err := os.Create("data/favorites.json")
	if err != nil {
		log.Println("Erreur création fichier favoris:", err)
		return
	}
	defer file.Close()
	if err := json.NewEncoder(file).Encode(favs); err != nil {
		log.Println("Erreur écriture favoris:", err)
	}
}
