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

var isThrowing = false
var throwingCounter int

func throwing(me PlayerState) string {
	var action string

	for isThrowing {
		if throwingCounter > 0 {
			action = "T"
		} else {
			//move away
			isRunning = true
			runningCounter = 2
			if isRunning {
				return runningAway(me)
			}
			isThrowing = false
		}
		throwingCounter--
	}
	prevRunningAction = action
	return action
}

var prevPrevScore int
var prevScore int

var retryThrow = false
var isRunning = false
var runningCounter int
var prevRunningAction string

func runningAway(me PlayerState) string {
	var action string

	for isRunning {
		if runningCounter > 0 {
			if prevRunningAction == "R" || prevRunningAction == "L" {
				action = "F"
			} else {
				var commands = []string{"F", "R", "L"}
				var rand = rand2.Intn(3)
				action = commands[rand]
			}
			log.Printf("RUNNING %s\n", action)

			runningCounter--
		} else {
			if prevRunningAction == "R" || prevRunningAction == "L" {
				action = "F"
			} else {
				var commands = []string{"F", "R", "L"}
				var rand = rand2.Intn(3)
				action = commands[rand]
			}
			log.Printf("RUNNING %s\n", action)
			isRunning = false
		}
	}
	prevRunningAction = action
	return action
}

func play(input ArenaUpdate) (response string) {
	me := input.Arena.State["https://radiation70-zaiqduddka-uc.a.run.app"]
	prevPrevScore = prevScore
	prevScore = me.Score
	log.Printf("Scores : %d > %d \n", prevPrevScore, prevScore)

	arenaSize := input.Arena.Dimensions

	if retryThrow {
		retryThrow = false
		log.Println("THROWING again")
		return "T"
	}

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

	if isRunning {
		return runningAway(me)

	}
	//look for nearby

	if prevScore > prevPrevScore {
		retryThrow = true
		log.Println("THROWING again")
		return "T"
	}
	if me.WasHit {
		log.Printf("[HIT] WHERE AM I: x:%d y:%d, dir:%s , %#v\n", me.X, me.Y, me.Direction, me)

		isRunning = true
		runningCounter = 2
		if isRunning {
			return runningAway(me)
		}

		var commands = []string{"F", "R", "L"}
		var rand = rand2.Intn(3)
		var action = commands[rand]
		log.Printf("MOVING to %s", action)

		return action
	} else {
		log.Printf("[THROWING] WHERE AM I: x:%d y:%d, dir:%s , %#v\n", me.X, me.Y, me.Direction, me)
		otherPlayers := input.Arena.State

		for playerName := range otherPlayers {
			if me.Direction == "N" {
				if otherPlayers[playerName].Y == me.Y {
					if (otherPlayers[playerName].X)+1 == me.X || (otherPlayers[playerName].X+2) == me.X {
						log.Printf("FOUND other PLAYER, Throwing\n")
						retryThrow = true
						return "T"
					}
				}
			} else if me.Direction == "E" {
				if otherPlayers[playerName].X == me.X {
					if (otherPlayers[playerName].Y)-1 == me.X || (otherPlayers[playerName].Y-2) == me.Y {
						log.Printf("FOUND other PLAYER, Throwing\n")
						retryThrow = true

						return "T"
					}
				}
			} else if me.Direction == "W" {
				if otherPlayers[playerName].X == me.X {
					if (otherPlayers[playerName].Y)+1 == me.X || (otherPlayers[playerName].Y+2) == me.Y {
						log.Printf("FOUND other PLAYER, Throwing\n")
						retryThrow = true
						return "T"
					}
				}
			} else if me.Direction == "S" {
				if otherPlayers[playerName].Y == me.Y {
					if (otherPlayers[playerName].X)-1 == me.X || (otherPlayers[playerName].X-2) == me.X {
						log.Printf("FOUND other PLAYER, Throwing\n")
						retryThrow = true
						return "T"
					}
				}
			}
		}

		var commands = []string{"F", "R", "L"}
		var rand = rand2.Intn(3)
		var action = commands[rand]
		log.Printf("NOT FOUND anyone, MOVING to %s\n", action)
		return action
	}

}
