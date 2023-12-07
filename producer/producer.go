package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	app "github.com/xenia11111/WB_TASKL0"
)

const path = "models/model2.json"

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	sc, _ := stan.Connect(viper.GetString("nats.clusterID"),
		viper.GetString("nats.clientProducer"),
		stan.NatsURL(viper.GetString("nats.serverURL")))
	defer sc.Close()

	PublishMessadge(sc, path)

	time.Sleep(1 * time.Second)

}

func PublishMessadge(sc stan.Conn, path string) {

	jsonFile, err := os.Open(path)

	if err != nil {
		logrus.Fatalf("error open json file: %s", err.Error())
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var order_body app.OrderBody
	json.Unmarshal(byteValue, &order_body)
	fmt.Println(order_body)

	sc.Publish("order", byteValue)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
