package storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Corray333/univer_cs/pkg/server/auth"
)

func (s *UserStorage) SetRefreshToken(id int, agent string, refresh string) error {
	var ctx = context.Background()
	if err := s.redis.Set(ctx, strconv.Itoa(id)+agent, refresh, auth.RefreshTokenLifeTime).Err(); err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) RefreshToken(id int, agent string, oldRefresh string) (string, string, error) {
	var ctx = context.Background()
	refresh, err := s.redis.Get(ctx, strconv.Itoa(id)+agent).Result()
	if err != nil {
		fmt.Println("1")
		return "", "", err
	}
	if refresh != oldRefresh {
		fmt.Println("2")
		return "", "", fmt.Errorf("invalid refresh token")
	}
	newRefresh, err := auth.CreateToken(id, auth.RefreshTokenLifeTime)
	if err != nil {
		fmt.Println("3")
		return "", "", err
	}
	newAccess, err := auth.CreateToken(id, auth.AccessTokenLifeTime)
	if err != nil {
		fmt.Println("4")
		return "", "", err
	}

	if err := s.redis.Set(ctx, strconv.Itoa(id)+agent, newRefresh, auth.RefreshTokenLifeTime).Err(); err != nil {
		fmt.Println("5")
		return "", "", err

	}
	return newAccess, newRefresh, nil

}
