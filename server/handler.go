package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dyluth/goMud/gomud"
	"github.com/gin-gonic/gin"
)

var (
	game *gomud.LoadedGame
)

// Start starts handling traffic
func Start(lg *gomud.LoadedGame) {
	port := os.Getenv("PORT")
	game = lg

	if port == "" {
		port = "8080"
		fmt.Printf("Defaulting to port %s", port)
	}

	// Starts a new Gin instance with no middle-ware
	r := gin.New()

	// Define handlers
	r.GET("/v1/getdescription", GetDescription)
	r.GET("/v1/move/:to", Move)
	r.GET("/v1/take/:item", Take)

	r.Run(":" + port)
}

//ErrorResponse - if we are returning an error that the client code should understand, we should use this struct
type ErrorResponse struct {
	Msg string
}

//GetDescription  returns the description of where the player is
func GetDescription(c *gin.Context) {
	_, room, ok := GetPlayerAndRoom(c)
	if !ok {
		return
	}
	//game.GetRoom(player.CurrentRoom)
	c.String(http.StatusOK, room.Describe())
}

//Move tries to move the player through the exit
func Move(c *gin.Context) {
	player, room, ok := GetPlayerAndRoom(c)
	if !ok {
		return
	}
	to := c.Param("to")
	// try to move the player to the new room
	newRoomName, ok := room.Exits[to]
	if !ok {
		c.String(http.StatusNotFound, fmt.Sprintf("exit %v not found", to))
		return
	}
	newRoom, ok := game.GetRoom(newRoomName)
	if !ok {
		//	look for a door instead
		door, ok := game.GetDoor(newRoomName)
		if !ok {
			c.JSON(http.StatusServiceUnavailable, ErrorResponse{Msg: fmt.Sprintf("uh-oh, %v is not in a valid room OR door", newRoomName)})
			return
		}
		throughDoorRoomName, description, moved := door.Enter(room.Name, player)
		if moved {
			newRoom, ok = game.GetRoom(throughDoorRoomName)
			if !ok {
				c.JSON(http.StatusServiceUnavailable, ErrorResponse{Msg: fmt.Sprintf("uh-oh, %v is not in a valid room OR door", throughDoorRoomName)})
				return
			}
			player.CurrentRoom = newRoom.Name
			c.String(http.StatusOK, fmt.Sprintf("%v\n%v\n", description, newRoom.Describe()))
			return
		}
		fmt.Printf("door: %v\n", description)
		c.String(http.StatusOK, description)
		return
	}

	player.CurrentRoom = newRoom.Name
	c.String(http.StatusOK, fmt.Sprintf("you move into %v\n%v", newRoom.Name, newRoom.Describe()))
}

// Take tries to take an item in the current room
func Take(c *gin.Context) {
	player, room, ok := GetPlayerAndRoom(c)
	if !ok {
		return
	}

	itemName := c.Param("item")
	if itemName == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Msg: "uh-oh, item should not be an empty string"})
		return
	}
	for i, item := range room.Contents {
		if item == itemName {
			player.Carrying = append(player.Carrying, item)
			copy(room.Contents[i:], room.Contents[i+1:])         // Shift a[i+1:] left one index.
			room.Contents[len(room.Contents)-1] = ""             // Erase last element (write zero value).
			room.Contents = room.Contents[:len(room.Contents)-1] // Truncate slice.
			c.String(http.StatusOK, fmt.Sprintf("you take the %v\n", item))
			return
		}
	}
	c.String(http.StatusNotFound, fmt.Sprintf("you cannot see a %v in the room to take\n", itemName))

}

// GetPlayerAndRoom gets the player specified by the AuthToken
// and the room they are in
// will return false if there was a problem and the calling method should just return
func GetPlayerAndRoom(c *gin.Context) (*gomud.Player, *gomud.Room, bool) {
	tok := c.GetHeader("token")
	// go through the players to find the right one
	player, ok := game.GetPlayer(tok)
	if !ok {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("Not a valid token"))
		return nil, nil, false
	}
	//fmt.Printf("player `%v` in room `%v`\n", player.Name, player.CurrentRoom)
	room, ok := game.GetRoom(player.CurrentRoom)
	if !ok {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Msg: fmt.Sprintf("uh-oh, %v is not in a valid room", player.CurrentRoom)})
		return nil, nil, false
	}

	// if not found return error
	return player, room, true

}
