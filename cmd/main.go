package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	app "github.com/xenia11111/WB_TASKL0"
	"github.com/xenia11111/WB_TASKL0/pkg/cache"
	"github.com/xenia11111/WB_TASKL0/pkg/handler"
	"github.com/xenia11111/WB_TASKL0/pkg/repository"
	service "github.com/xenia11111/WB_TASKL0/pkg/service"
	subscriber "github.com/xenia11111/WB_TASKL0/subscriber"
)

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variable: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	sc, err := stan.Connect(viper.GetString("nats.clusterID"),
		viper.GetString("nats.clientSubscriber"),
		stan.NatsURL(viper.GetString("nats.serverURL")))

	defer sc.Close()

	repos := repository.NewRepository(db)
	cache := cache.NewCache(0, 0, repos, 100)
	defer cache.GC()
	service := service.NewService(repos, cache)
	handler := handler.NewHandler(service)

	subscriber.NewSubscriber(sc, service)

	srv := new(app.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
