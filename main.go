package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var client ApiClient
var reader *bufio.Reader

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("TMDB_API_RA_TOKEN")
	client = ApiClient{}.New(apiKey)

	reader = bufio.NewReader(os.Stdin)
}

func main() {
	var currentActor Actor
	var endActor Actor
	var currentMovie Movie

	// 1. Get random start and end actors
	currentActor, endActor, err := client.GetActors()
	check(err)
	for {
		for {
			// 1a. Display start and end
			fmt.Printf("Actor : %s\n", currentActor.Name)
			fmt.Printf("Goal  : %s\n", endActor.Name)

			// 2. Ask for movie from user
			movieName, err := input("\nGuess movie:")
			check(err)

			// 3. Search for movies with the name provided
			movies, err := client.SearchMovies(movieName)
			check(err)

			// 4. For each currentMovie from search, get credits from id and
			//    see whether the actor is credited or not.
			var success bool = false
			for _, m := range movies {
				success = client.MovieHasCast(m.Id, currentActor.Id)
				if success {
					currentMovie = m
					break
				}
			}

			// 5. Report success/fail
			if success {
				fmt.Print("\nCorrect!\n\n")
				break
			} else {
				fmt.Print("\nIncorrect. Try again.\n\n")
			}
		}

		for {
			// 6. Print new current movie

			fmt.Printf("Movie : %s (%s)\n", currentMovie.Title, strings.Split(currentMovie.ReleaseDate, "-")[0])
			fmt.Printf("Goal  : %s\n", endActor.Name)

			// 7. Ask to guess actor in new movie
			actorName, err := input("\nGuess actor:")
			check(err)

			// 8. Search for actors with name provided
			actors, err := client.SearchActors(actorName)
			check(err)

			// 9. For each actor, see if they're in the movie
			var success bool = false
			for _, a := range actors {
				success = client.MovieHasCast(currentMovie.Id, a.Id)
				if success {
					currentActor = a
					break
				}
			}

			if success {
				fmt.Print("\nCorrect!\n\n")
				break
			} else {
				fmt.Print("\nIncorrect. Try again.\n\n")
			}
		}

		if currentActor.Id == endActor.Id {
			fmt.Println("You won!")
			os.Exit(0)
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func input(question string) (string, error) {
	fmt.Printf("%s\n> ", question)
	resp, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	resp = resp[:len(resp)-1] // Remove the newline char

	return resp, nil
}
