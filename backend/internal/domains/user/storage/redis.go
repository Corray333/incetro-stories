package storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Corray333/univer_cs/pkg/server/auth"
)

var ctx = context.Background()

func (s *UserStorage) SetRefreshToken(id int, agent string, refresh string) error {
	if err := s.redis.Set(ctx, strconv.Itoa(id)+agent, refresh, auth.RefreshTokenLifeTime).Err(); err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) RefreshToken(id int, agent string, oldRefresh string) (string, string, error) {
	refresh, err := s.redis.Get(ctx, strconv.Itoa(id)+agent).Result()
	if err != nil {
		return "", "", err
	}
	if refresh != oldRefresh {
		return "", "", fmt.Errorf("invalid refresh token")
	}
	newRefresh, err := auth.CreateToken(id, auth.RefreshTokenLifeTime)
	if err != nil {
		return "", "", err
	}
	newAccess, err := auth.CreateToken(id, auth.AccessTokenLifeTime)
	if err != nil {
		return "", "", err
	}

	if err := s.redis.Set(ctx, strconv.Itoa(id)+agent, newRefresh, auth.RefreshTokenLifeTime).Err(); err != nil {
		return "", "", err

	}
	return newAccess, newRefresh, nil

}
