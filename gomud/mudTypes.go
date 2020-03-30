package gomud

import (
	"fmt"
	"strings"
)

// place to put types for this project

// GameFile is the top level description of the game to run
type GameFile struct {
	Players []Player `json:"players"`
	Rooms   []Room   `json:"rooms"`
	Doors   []Door   `json:"doors"`
}

// Room is the basis for exploration
type Room struct { // zonetype: room
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits,omitempty"`    // map of keywords to enter/use to zones they lead to - can be doors or rooms
	Contents    []string          `json:"contents,omitempty"` // map of things in the room
}

// Describe describes the room
func (r *Room) Describe() string {
	contents := ""
	exits := ""

	if len(r.Contents) > 0 {
		contents = strings.Join(r.Contents, ", ")
		contents = fmt.Sprintf("you can see:\n%v\n", contents)
	}
	if len(r.Exits) > 0 {
		e := make([]string, len(r.Exits))
		i := 0
		for k := range r.Exits {
			e[i] = k
			i++
		}
		exits = strings.Join(e, ", ")
		exits = fmt.Sprintf("exits to the room are:\n%v\n", exits)
	}
	desc := fmt.Sprintf("%v\n%v%v", r.Description, contents, exits)
	return desc

}

// Door can be between 2 zones
// going through the door takes you into the zone you are not in
type Door struct {
	Name            string `json:"name"`
	OpenDescription string `json:"description"` // descriotion of the door opening
	Zone1           string `json:"zone1"`
	Zone2           string `json:"zone2"`
	StartsOpen      bool   `json:"is-open,omitempty"`
	NeedsToOpen     string `json:"needs,omitempty"` // name of the thing you need to open the door
}

// Player is a desccription of a player
type Player struct {
	Name        string   `json:"name"`
	AuthToken   string   `json:"token"` // this is their auth token.. secure i know..
	Description string   `json:"description"`
	Carrying    []string `json:"has,omitempty"`
	CurrentRoom string   `json:"room"`
}

// LoadGame loads the game from the gamefile struct
func LoadGame(gf GameFile) LoadedGame {
	lg := LoadedGame{
		rooms:   make(map[string]*Room),
		players: make(map[string]*Player),
		doors:   make(map[string]*Door),
	}
	for _, room := range gf.Rooms {
		if room.Contents == nil {
			room.Contents = []string{}
		}
		if room.Exits == nil {
			room.Exits = make(map[string]string)
		}
		r := room
		lg.rooms[room.Name] = &r
	}
	for _, player := range gf.Players {
		lg.players[player.AuthToken] = &player
	}
	for _, door := range gf.Doors {
		lg.doors[door.Name] = &door
	}

	return lg
}

// LoadedGame represents the game ready to run
type LoadedGame struct {
	rooms   map[string]*Room
	players map[string]*Player
	doors   map[string]*Door
}

// GetRoom gets a room by its ID
func (lg *LoadedGame) GetRoom(roomName string) (*Room, bool) {
	room, ok := lg.rooms[roomName]
	fmt.Printf("getting room `%v`: %v ok %v\n", roomName, room, ok)
	return room, ok
}

// GetPlayer gets a player by its ID
func (lg *LoadedGame) GetPlayer(authToken string) (*Player, bool) {
	p, ok := lg.players[authToken]
	return p, ok
}
