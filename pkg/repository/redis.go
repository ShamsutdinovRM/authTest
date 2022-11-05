package repository

import (
	"authTest/model"
	"context"
	"github.com/go-redis/redis/v9"
	"log"
	"time"
)

type OperationsRedis interface {
	LoginUser(user model.LogUser) (model.LogUser, error)
	LogoutUser(user model.Name) error
	CheckUser(user model.Name) error
}

type RedisConnect struct {
	Conn *redis.Client
}

const (
	tokenTTL = 5 * time.Minute
)

var ctx = context.Background()

func NewRedisConnect() *RedisConnect {
	client := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "",
		DB:       0,
	})
	return &RedisConnect{Conn: client}
}

func (r *RedisConnect) LoginUser(user model.LogUser) (model.LogUser, error) {

	err := r.Conn.Set(ctx, user.Username, user.Token, tokenTTL).Err()
	if err != nil {
		log.Printf("Error save data user login: %s", err)
		return model.LogUser{}, err
	}

	return user, nil
}

func (r *RedisConnect) LogoutUser(user model.Name) error {
	err := r.Conn.Del(ctx, user.Username).Err()
	if err != nil {
		log.Printf("error del user cache: %s", err)
		return err
	}

	return nil
}

func (r *RedisConnect) CheckUser(user model.Name) error {
	_, err := r.Conn.Get(ctx, user.Username).Result()
	if err == redis.Nil {
		log.Printf("User doesn't login: %s", err)
		return err
	} else if err != nil {
		log.Printf("Error check user: %s", err)
		return err
	}

	return nil
}
