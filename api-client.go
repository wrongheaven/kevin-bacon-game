package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

type ApiClient struct {
	baseUrl string
	apiKey  string
}

func (ac ApiClient) New(apiKey string) ApiClient {
	return ApiClient{
		baseUrl: "https://api.themoviedb.org/3/",
		apiKey:  apiKey,
	}
}

func (ac ApiClient) GetActors() (Actor, Actor, error) {
	body, err := ac.doRequest("person/popular", "")
	if err != nil {
		return Actor{}, Actor{}, err
	}

	var apiResponse PopularActorsResponse
	json.Unmarshal(body, &apiResponse)

	ai := rand.Intn(len(apiResponse.Results))
	bi := rand.Intn(len(apiResponse.Results))
	for ai == bi {
		bi = rand.Intn(len(apiResponse.Results))
	}

	return apiResponse.Results[ai], apiResponse.Results[bi], nil
}

func (ac ApiClient) SearchMovies(title string) ([]Movie, error) {
	query := strings.ReplaceAll("?query="+title, " ", "+")
	body, err := ac.doRequest("search/movie", html.EscapeString(query))
	if err != nil {
		return []Movie{}, err
	}

	var response MovieSearchResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []Movie{}, err
	}

	return response.Results, nil
}

func (ac ApiClient) SearchActors(name string) ([]Actor, error) {
	query := strings.ReplaceAll("?query="+name, " ", "+")
	body, err := ac.doRequest("search/person", html.EscapeString(query))
	if err != nil {
		return []Actor{}, err
	}

	var response ActorSearchResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []Actor{}, err
	}

	return response.Results, nil
}

func (ac ApiClient) MovieHasCast(movieId int, actorId int) bool {
	body, err := ac.doRequest(fmt.Sprintf("movie/%d/credits", movieId), "")
	if err != nil {
		return false
	}

	var response MovieCastResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false
	}

	for _, actor := range response.Results {
		if actorId == actor.Id {
			return true
		}
	}

	return false
}

func (ac ApiClient) doRequest(path string, query string) ([]byte, error) {
	url, err := url.JoinPath(ac.baseUrl, path)
	if err != nil {
		return nil, err
	}
	url += query

	// fmt.Println("url:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+ac.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
