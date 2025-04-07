package util

import (
	pb "github.com/Erickype/DungeonPit/proto"
)

// MoveCoordinates returns the new coordinates after moving in a given direction.
func MoveCoordinates(x, y, z int, direction pb.Direction) (newX, newY, newZ int) {
	switch direction {
	case pb.Direction_NORTH:
		return x, y, z + 1
	case pb.Direction_SOUTH:
		return x, y, z - 1
	case pb.Direction_EAST:
		return x + 1, y, z
	case pb.Direction_WEST:
		return x - 1, y, z
	case pb.Direction_UP:
		return x, y + 1, z
	case pb.Direction_DOWN:
		return x, y - 1, z
	default:
		return x, y, z // No movement for unknown direction
	}
}
