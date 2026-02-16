package main

import (
	"log"
	"net/http"
	"spotidex/handlers"
	"spotidex/services"
)

func main() {
	// On lance le truc pour spotify
	services.InitService()
	handlers.Init()

	// Dossier pour les css et tout
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Les routes du site
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/artist", handlers.ArtistHandler)
	http.HandleFunc("/favorites", handlers.FavoritesHandler)
	http.HandleFunc("/api/favorite", handlers.AddFavoriteHandler)

	log.Println("Le serveur marche sur http://localhost:8080")
	// On lance le serveur
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
