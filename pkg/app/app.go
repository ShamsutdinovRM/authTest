package app

import (
	"authTest/model"
	"authTest/pkg/handler"
	"authTest/pkg/repository"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func Run(path string) {
	if err := initConfig(path); err != nil {
		log.Printf("error initializaing configs: %s", err)
		return
	}

	DBSchema := model.DB{
		Username: viper.GetString("db.username"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Schema:   viper.GetString("db.schema"),
		Password: viper.GetString("db.password"),
	}

	db, err := repository.New(DBSchema)
	if err != nil {
		log.Printf("Error create DB connection: %s", err)
		return
	}
	defer db.DB.Close()

	dbr := repository.NewRedisConnect()
	//defer dbr.Conn.Close()

	hand := handler.Repos{
		Repository:      db,
		RedisRepository: dbr,
	}

	r := mux.NewRouter()
	r.Use(handler.CommonMiddleware)
	r.HandleFunc("/signup", hand.SignUp)
	r.HandleFunc("/login", hand.LogIn)
	r.HandleFunc("/logout", hand.LogOut)

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(hand.JWTMiddleware)
	s.HandleFunc("/hello", hand.Hello)

	port := viper.GetString("port")
	log.Printf("server started")
	log.Fatal(http.ListenAndServe(port, r))
}

func initConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}
