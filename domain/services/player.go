package services

import (
	"errors"

	"sts/domain/essence"
)

var (
	ErrPlayerNotEnoughPoints = errors.New("Player not enough points")
)

type PlayerService struct {
	players essence.PlayerRepository
}

func NewPlayerService(repo essence.PlayerRepository) *PlayerService {
	return &PlayerService{players: repo}
}

func (ps PlayerService) IncreasePoints(playerID string, points int64) error {
	player, err := ps.players.FindByID(playerID)
	if err != nil {
		if err == essence.ErrPlayerNotFound {
			player = &essence.Player{ID: playerID}
			if err := ps.players.Create(player); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	amount := points * 100
	if points < 0 && player.Balance+amount < 0 {
		return ErrPlayerNotEnoughPoints
	}

	if err := ps.players.UpdateBalance(player, amount); err != nil {
		return err
	}

	return nil
}
