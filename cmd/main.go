package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	todoappgo "todo-app-go"
	"todo-app-go/pkg/handler"
	"todo-app-go/pkg/repository"
	"todo-app-go/pkg/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing congifs %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("error initializing DB %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(todoappgo.Server)

	go func ()  {
		if err := srv.Run(viper.GetString("port"), handlers.InitHandler()); err != nil {
			logrus.Fatalf("failed to start server: %v", err.Error())
		}
	}()
	logrus.Print("TodoApp Started...")
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit

	logrus.Print("TodoApp Shutting down")
	if err :=srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on the shutting down server %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on the connection %s", err.Error())
	}

	
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
