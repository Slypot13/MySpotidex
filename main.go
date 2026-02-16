package main

import (
	"log"
	"net/http"
	"spotidex/handlers"
	"spotidex/services"
)

func main() {
	// Initialiser le service Spotify
	spotifyService := services.NewSpotifyService()
	handlers.Init(spotifyService)

	// Servir les fichiers statiques (CSS, JS, Images)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Définir les routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/artist", handlers.ArtistHandler)
	http.HandleFunc("/favorites", handlers.FavoritesHandler)
	http.HandleFunc("/api/favorite", handlers.AddFavoriteHandler)

	log.Println("Serveur démarré sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
