package subscriber

import (
	"encoding/json"

	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"github.com/sirupsen/logrus"
	app "github.com/xenia11111/WB_TASKL0"
	service "github.com/xenia11111/WB_TASKL0/pkg/service"
)

type Subscriber struct {
	sc      stan.Conn
	service *service.Service
}

func NewSubscriber(sc stan.Conn, service *service.Service) *Subscriber {
	st := &Subscriber{
		sc:      sc,
		service: service,
	}
	st.Subscriber()

	return st
}

func (s *Subscriber) Subscriber() {

	s.sc.Subscribe("order", func(msg *stan.Msg) {

		var order_body app.OrderBody

		err := json.Unmarshal(msg.Data, &order_body)
		if err != nil {
			logrus.Errorf("error parse order from stream: %s", err.Error())
			return
		}

		if order_body.Order_uid == "" {
			logrus.Error("error order_uid cannot be empty")
			return
		}

		order := app.Order{
			Order_uid:  order_body.Order_uid,
			Order_body: order_body,
		}

		err = s.service.Create(order)

		if err != nil {
			logrus.Errorf("cannot be create order: %s", err.Error())
			return
		}
	}, stan.StartAt(pb.StartPosition_NewOnly))
}
