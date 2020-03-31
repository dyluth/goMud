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
	Name             string    `json:"name"`
	EnterDescription string    `json:"description"`
	LeadsTo          []string  `json:"leads-to"` // only makes sense to have 1 or 2 rooms / doors in this list
	Openable         *Openable `json:"openable"`
}

// Enter describes trying to enter through a door and entering the new room
func (d *Door) Enter(from string, player *Player) (newRoomName, description string, moved bool) {
	moved = true
	ok := false
	// if the door is openable, check to see if we can get through it
	if d.Openable != nil {
		if !d.Openable.IsOpen {
			if d.Openable.AutoOpens {
				description, ok = d.Open(player)
				if !ok {
					return from, description, false
				}
			} else {
				return from, d.Openable.LockedDescription, false
			}
		}
	}
	// if we get here the door is now open and we can move through it
	newRoomName = from
	for _, room := range d.LeadsTo {
		if room != from {
			newRoomName = room
			description = fmt.Sprintf("%v%v", description, d.EnterDescription)
			break
		}
	}
	return
}

// Open tries to open the door
func (d *Door) Open(player *Player) (description string, success bool) {
	fmt.Printf("trying to open door %v\n", d.Name)
	if d.Openable == nil {
		return fmt.Sprintf("the %v doesnt open or close", d.Name), false
	}
	// make sure we have all the things we need to open the door
	foundCount := 0
	for _, need := range d.Openable.NeedsToOpen {
		for _, have := range player.Carrying {
			if have == need {
				foundCount++
				break
			}
		}
	}
	if foundCount < len(d.Openable.NeedsToOpen) {
		return d.Openable.LockedDescription, false
	}
	return d.Openable.OpenDescription, true
}

// Openable describes how to open a door
type Openable struct {
	IsOpen            bool     `json:"is-open,omitempty"`
	AutoOpens         bool     `json:"auto-opens,omitempty"`  // if true automatically opens when you have the needsToOpen items
	NeedsToOpen       []string `json:"open-needs,omitempty"`  // name of the thing you need to open the door
	NeedsToClose      []string `json:"close-needs,omitempty"` // name of the thing you need to close the door
	OpenDescription   string   `json:"open-description"`      // description of the door opening
	CloseDescription  string   `json:"close-description"`     // description of the door Closing
	LockedDescription string   `json:"locked-description"`    // description when you try to open but can't

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
		r := room // make a copy
		lg.rooms[room.Name] = &r
	}
	for _, player := range gf.Players {
		p := player // make a copy
		lg.players[player.AuthToken] = &p
	}
	for _, door := range gf.Doors {
		d := door // make a copy
		lg.doors[door.Name] = &d
	}
	return lg
}

// LoadedGame represents the game ready to run
type LoadedGame struct {
	rooms   map[string]*Room
	players map[string]*Player
	doors   map[string]*Door
}

// GetRoom gets a room by its name
func (lg *LoadedGame) GetRoom(roomName string) (*Room, bool) {
	room, ok := lg.rooms[roomName]
	fmt.Printf("getting room `%v`: %v ok %v\n", roomName, room, ok)
	return room, ok
}

// GetDoor gets a door by its name
func (lg *LoadedGame) GetDoor(doorName string) (*Door, bool) {
	door, ok := lg.doors[doorName]
	fmt.Printf("getting door `%v`: %v ok %v\n", doorName, door, ok)
	return door, ok
}

// GetPlayer gets a player by its name
func (lg *LoadedGame) GetPlayer(authToken string) (*Player, bool) {
	p, ok := lg.players[authToken]
	return p, ok
}
