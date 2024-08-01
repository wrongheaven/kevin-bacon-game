package main

type MovieSearchResponse struct {
	Results []Movie `json:"results"`
}

type PopularActorsResponse struct {
	Results []Actor `json:"results"`
}

type MovieCastResponse struct {
	Results []Actor `json:"cast"`
}

type ActorSearchResponse struct {
	Results []Actor `json:"results"`
}

type Actor struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Movie struct {
	Id          int    `json:"id"`
	Title       string `json:"Title"`
	ReleaseDate string `json:"release_date"`
}
