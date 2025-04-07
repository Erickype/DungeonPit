package logic

import (
	"github.com/Erickype/DungeonPit/internal/model"
	"log"
)

func LoginPlayer(username string) (*model.Player, error) {
	log.Printf("Login requested for username: %s", username)
	player := &model.Player{}
	player.Username = username
	err := player.Login()
	if err != nil {
		return nil, err
	}
	log.Printf("Player %s succesfully loged in!", player.Username)

	return player, nil
}

func SetInitialPlayerRoom(player *model.Player) (room *model.Room, err error) {
	initialRoom := &model.Room{}
	currentRoom := &model.Room{}
	err = currentRoom.GetByRoomID(player.CurrentRoomID)
	if err != nil {
		log.Println("Player current room not found, teleporting to Room Zero")
		err = initialRoom.GetZeroRoom()
		if err != nil {
			return nil, err
		}
	} else {
		log.Printf("Player current room found, teleporting to Room (%d,%d,%d)",
			currentRoom.X, currentRoom.Y, currentRoom.Z)
		initialRoom = currentRoom
	}

	return initialRoom, nil
}
