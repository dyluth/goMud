# goMUD 
a mud engine written in go!

seemed like a good idea at the time.

## run the server
from the root of the project run:
`go run cmd/server/main.go`
you will need to have a valid `game.json`  descriptor file - there is a default one to play with though

## interact with it - CURL
describe the current room:
`curl localhost:8080/v1/getdescription`
  
move from one room through an exit:
`curl localhost:8080/v1/move/:exitName`
  
take an item in a room
`curl localhost:8080/v1/take/:itemName`
