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

	//throw
	//turn
	if prevAction == "R" {
		prevAction = "T"
	} else {
		prevAction = "F"
	}

	if me.WasHit {
		return "L"
	}

	log.Println("")
	log.Printf("IN: %#v\n", input)
	log.Printf("WHERE AM I: %#v\n", me)
	log.Printf("DOING: %s", prevAction)
	log.Println("")

	commands := []string{"F", "R", "L", "T"}
	rand := rand2.Intn(4)
	return commands[rand]
}
