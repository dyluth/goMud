package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/dyluth/goMud/gomud"
	"github.com/dyluth/goMud/server"
)

// The main for running the mud Engine
func main() {
	fileName := os.Getenv("GAMEFILE")
	if fileName == "" {
		fileName = "./game.json"
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	gf := &gomud.GameFile{} //&gomud.GameFile{}

	err = json.Unmarshal(data, gf)
	if err != nil {
		panic(err)
	}

	// do some validation of the file
	err = validate(gf)
	if err != nil {
		panic(err)
	}
	g := gomud.LoadGame(*gf)

	//start the server
	server.Start(&g)

	// wait for ever
	<-make(chan bool)
}

func validate(gf *gomud.GameFile) error {
	// make sure there is at least 1 room and 1 player

	// go through all the rooms
	// make sure they all have a unique name
	// make sure all room exits go to a real place
	// make sure that all rooms are linked to another room

	// go through all doors
	// make sure all doors are connected to 2 different rooms

	// go through all players
	// make sure they are all in a valid room
	return nil
}
