package model

import (
	"errors"
	"github.com/Erickype/DungeonPit/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Player struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email         string    `gorm:"unique;not null" json:"email"`
	Username      string    `gorm:"unique;not null" json:"username"`
	CurrentRoomID uuid.UUID `json:"current_room_id"`
}

func (player *Player) Login() (err error) {
	err = database.DB.Where("username = ?", player.Username).First(&player).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return errors.New("player not found")
		}
		return err
	}
	return nil
}

func (player *Player) Get() (err error) {
	err = database.DB.Where(&player).First(&player).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return errors.New("player not found")
		}
		return err
	}
	return nil
}

func (player *Player) Create() (err error) {
	err = database.DB.Create(&player).Error
	if err != nil {
		if errors.Is(gorm.ErrDuplicatedKey, err) {
			return errors.New("player already exists")
		}
		return err
	}
	return nil
}

func (player *Player) SetCurrentRoomID(currentRoomID uuid.UUID) (err error) {
	err = database.DB.Model(player).
		Update("current_room_id", currentRoomID).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return errors.New("player not found")
		}
		return err
	}
	return nil
}
