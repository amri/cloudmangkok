package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
)

var prevAction string

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := play(v)
	fmt.Fprint(w, resp)
}

func play(input ArenaUpdate) (response string) {
	me := input.Arena.State["https://radiation70-zaiqduddka-uc.a.run.app"]

	arenaSize := input.Arena.Dimensions

	//CORRECTION
	//X:0 DIRECTION should be east
	if me.X == 0 {
		if me.Direction == "N" {
			return "R"
		} else if me.Direction == "S" {
			return "L"
		} else if me.Direction == "W" {
			return "R"
		}
	}

	if me.Y == 0 {
		if me.Direction == "N" {
			return "R"
		} else if me.Direction == "E" {
			return "R"
		} else if me.Direction == "W" {
			return "L"
		}
	}

	if me.Y == arenaSize[1]-1 {
		if me.Direction == "S" {
			return "L"
		} else if me.Direction == "E" {
			return "L"
		} else if me.Direction == "W" {
			return "R"
		}
	}

	//X:X DIRECTION should be WEST
	if me.X == arenaSize[0]-1 {
		if me.Direction == "N" {
			return "L"
		} else if me.Direction == "S" {
			return "R"
		} else if me.Direction == "E" {
			return "L"
		}
	}

	if me.WasHit {
		log.Printf("[HIT] WHERE AM I: x:%d y:%d, dir:%s , %#v\n", me.X, me.Y, me.Direction, me)

		var commands = []string{"F", "R", "L"}
		var rand = rand2.Intn(3)
		var action = commands[rand]
		log.Printf("MOVING to %s", action)

		return action
	} else {
		log.Printf("[THROWING] WHERE AM I: x:%d y:%d, dir:%s , %#v\n", me.X, me.Y, me.Direction, me)

		return "T"
	}

	log.Println("")
	// log.Printf("IN: %#v\n", input)
	// log.Printf("DOING: %s", prevAction)
	log.Println("")

	commands := []string{"F", "R", "L", "T"}
	rand := rand2.Intn(4)
	return commands[rand]
}
