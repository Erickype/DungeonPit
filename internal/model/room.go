package model

import (
	"errors"
	"github.com/Erickype/DungeonPit/internal/database"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Room struct {
	X          int            `gorm:"primaryKey" json:"x"`
	Y          int            `gorm:"primaryKey" json:"y"`
	Z          int            `gorm:"primaryKey" json:"z"`
	RoomID     uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()" json:"roomID"`
	Discovered bool           `gorm:"default:false" json:"discovered"`
	Data       datatypes.JSON `gorm:"type:jsonb" json:"data"`
}

func (room *Room) GetZeroRoom() error {
	*room = Room{X: 0, Y: 0, Z: 0}
	err := database.DB.First(room, "x = ? AND y = ? AND z = ?", room.X, room.Y, room.Z).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("room not found")
		}
		return err
	}
	return nil
}

func (room *Room) GetRoom() (err error) {
	err = database.DB.First(room, "x = ? AND y = ? AND z = ?", room.X, room.Y, room.Z).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return errors.New("room not found")
		}
		return err
	}
	return nil
}

func (room *Room) GetByRoomID(roomID uuid.UUID) (err error) {
	err = database.DB.Where("room_id = ?", roomID).First(&room).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return errors.New("room not found")
		}
		return err
	}
	return nil
}

func (room *Room) Create() (err error) {
	err = database.DB.Create(&room).Error
	if err != nil {
		if errors.Is(gorm.ErrDuplicatedKey, err) {
			return errors.New("room already exists")
		}
		return err
	}
	return nil
}
