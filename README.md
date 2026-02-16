# My SpotiDex - Projet Groupie Tracker

Voici My SpotiDex, un projet du module Groupie-Tracker, on avait des règles a suivre comme le système de recherche, de favoris,..
    
 Dans ce document tout sera indiqué :

## Description

C'est un site qui utilise l'API Spotify pour chercher des artistes et voir leurs infos. On peut voir les albums, la popularité etc.

## Comment ça marche

Le projet est fait en Go avec des templates HTML. J'ai utilisé l'API Spotify pour récupérer les données, on avait deja un document pour faire validé notre projet, avec deja les endpoints,...

### Fichiers importants :
- `main.go` : Le fichier principal qui lance le serveur
- `handlers/` : Les fonctions pour gérer les pages
- `services/` : Le code pour communiquer avec Spotify
- `templates/` : Les pages HTML
- `static/css/` : Le fichier CSS pour le style

## Installation

Il faut avoir Go installé sur votre ordinateur.

1. Cloner le projet
```bash
git clone https://github.com/Slypot13/MySpotidex.git
cd MySpotidex
```

2. Lancer le serveur
```bash
go run main.go
```

3. Ouvrir dans le navigateur : `http://localhost:8080`

## Fonctionnalités

- Recherche d'artistes
- Affichage des détails (albums, popularité, followers)
- Page de favoris (pas fini )
- Filtres par genre et année ( pour ca il faut d'abord rechercher quelque chose pour etre sur la page de recherche et avoir la fonctionnalité)

## Problèmes connus

- Parfois l'authentification Spotify plante, faut relancer
- Les favoris marchent pas encore très bien
- Le CSS est un peu basique mais ça marche

## Langage utilisé

- Go (langage backend)
- HTML/CSS (frontend)
- API Spotify Developpers (pour les données)

## Auteur

Pottier Sylvain 

---

**Note** : Pour que ça marche il faut avoir une connexion internet pour accéder à l'API Spotify.
