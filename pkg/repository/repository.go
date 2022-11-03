package repository

import (
	"authTest/model"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Operations interface {
	CreateUser(body model.User) (model.User, error)
	LoginUser(body model.User) (model.User, error)
}

type DBModel struct {
	DB *sql.DB
}

func New(cfg model.DB) (*DBModel, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	log.Println(conn)
	db, err := sql.Open(cfg.Schema, conn)
	if err != nil {
		log.Printf("Error open connect DB: %s", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Printf("Error Ping DB: %s", err)
		return nil, err
	}

	return &DBModel{DB: db}, nil
}

func (db *DBModel) CreateUser(body model.User) (model.User, error) {

	_, err := db.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", body.Username, body.Password)
	if err != nil {
		log.Printf("error create user: %s", err)
		return model.User{}, fmt.Errorf("error create user: %s", err)
	}

	return model.User{Username: body.Username}, nil
}

func (db *DBModel) LoginUser(body model.User) (model.User, error) {

	rowUser := db.DB.QueryRow("SELECT password FROM users WHERE username=$1", body.Username)

	var userInf model.User
	err := rowUser.Scan(&userInf.Password)
	if err != nil {
		log.Printf("error select password: %s", err)
		return model.User{}, fmt.Errorf("error, user not found: %s", err)
	}

	userInf.Username = body.Username

	return userInf, nil
}
