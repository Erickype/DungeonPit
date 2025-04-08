package logic

import (
	"fmt"
	"github.com/Erickype/DungeonPit/internal/log"
	"github.com/Erickype/DungeonPit/internal/model"
)

func LoginPlayer(username string) (*model.Player, error) {
	log.GetCoreInstance().Info("Login requested for username:", username)
	player := &model.Player{}
	player.Username = username
	err := player.Login()
	if err != nil {
		return nil, err
	}
	log.GetCoreInstance().Info("Player", player.Username, "succesfully logged in!")

	return player, nil
}

func SetInitialPlayerRoom(player *model.Player) (room *model.Room, err error) {
	initialRoom := &model.Room{}
	currentRoom := &model.Room{}
	err = currentRoom.GetByRoomID(player.CurrentRoomID)
	if err != nil {
		log.GetCoreInstance().Info("Player current room not found, teleporting to Room Zero")
		err = initialRoom.GetZeroRoom()
		if err != nil {
			return nil, err
		}
	} else {
		message := fmt.Sprintf("Player current room found, teleporting to Room (%d, %d, %d)",
			currentRoom.X, currentRoom.Y, currentRoom.Z)
		log.GetCoreInstance().Info(message)
		initialRoom = currentRoom
	}

	return initialRoom, nil
}
