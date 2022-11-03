package model

type DB struct {
	Username string `json:"username"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Schema   string `json:"schema"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogUser struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type DefaultError struct {
	Text string `json:"text"`
}
