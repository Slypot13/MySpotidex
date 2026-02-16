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

// Structure pour les données de la page
// J'ai mis tout dedans c'est plus simple
type DonneesPage struct {
	Title     string
	Artists   []models.Artist
	Artist    *models.Artist
	Albums    []models.Album
	Query     string
	Offset    int
	Favorites []models.Artist
}

// Plus besoin d'init le service car on utilise des variables globales dans services
func Init() {
	// Vide
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// On charge des artistes pour l'accueil
	resp, err := services.RechercheArtists("genre:pop", 0)
	var artistes []models.Artist
	if err == nil {
		artistes = resp.Artists.Items
	} else {
		log.Println("Bug chargement accueil:", err)
	}

	renderTemplate(w, "home", DonneesPage{Title: "Accueil", Artists: artistes})
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	genre := r.URL.Query().Get("genre")
	year := r.URL.Query().Get("year")

	offsetStr := r.URL.Query().Get("offset")
	offset, _ := strconv.Atoi(offsetStr)

	data := DonneesPage{Title: "Recherche", Query: query, Offset: offset}

	if query != "" {
		recherche := query
		// Gestion des filtres un peu moche mais ça marche
		if genre != "" {
			recherche += " genre:\"" + genre + "\""
		}
		if year != "" {
			recherche += " year:" + year
		}

		resultat, err := services.RechercheArtists(recherche, offset)
		if err != nil {
			log.Println("Erreur recherche:", err)
			http.Error(w, "Erreur search", 500)
			return
		}
		data.Artists = resultat.Artists.Items
	}

	renderTemplate(w, "search_results", data)
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		// Pas d'id, erreur
		http.Error(w, "Y'a pas d'ID", 400)
		return
	}

	artiste := services.GetArtiste(id)
	if artiste == nil {
		http.Error(w, "Artiste pas trouvé", 500)
		return
	}

	albums := services.GetAlbumsArtiste(id)

	// On remplit les données
	data := DonneesPage{
		Title:  artiste.Name,
		Artist: artiste,
		Albums: albums,
	}

	renderTemplate(w, "artist_detail", data)
}

func FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	favs, err := loadFavorites()
	if err != nil {
		favs = []models.Artist{}
	}
	renderTemplate(w, "favorites", DonneesPage{Title: "Favoris", Favorites: favs})
}

func AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Pas autorisé", 405)
		return
	}

	var atels models.Artist
	if err := json.NewDecoder(r.Body).Decode(&atels); err != nil {
		http.Error(w, "Pas bon JSON", 400)
		return
	}

	// J'utilise _ pour ignorer l'erreur car ça devrait aller
	favs, _ := loadFavorites()
	favs = append(favs, atels)
	saveFavorites(favs)

	w.WriteHeader(200)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data DonneesPage) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	t, err := template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/"+tmpl+".html")
	if err != nil {
		http.Error(w, "Erreur template: "+err.Error(), 500)
		return
	}
	t.ExecuteTemplate(w, "layout", data)
}

// Fonction pour lire le JSON
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

// Pour sauvegarder
func saveFavorites(favs []models.Artist) {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Println("Erreur dossier data:", err)
		return
	}
	file, err := os.Create("data/favorites.json")
	if err != nil {
		log.Println("Erreur fichier favoris:", err)
		return
	}
	defer file.Close()
	json.NewEncoder(file).Encode(favs)
}
