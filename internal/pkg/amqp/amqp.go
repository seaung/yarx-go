package amqp

import (
	"fmt"

	"github.com/RichardKnop/machinery/v2"
	backends "github.com/RichardKnop/machinery/v2/backends/amqp"
	brokers "github.com/RichardKnop/machinery/v2/brokers/amqp"
	"github.com/RichardKnop/machinery/v2/config"
	eager "github.com/RichardKnop/machinery/v2/locks/eager"
)

func newAmqp(username, password, host string, port int) *machinery.Server {
	cfg := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)

	cnf := &config.Config{
		Broker:          cfg,
		DefaultQueue:    "machinery_queue",
		ResultBackend:   cfg,
		ResultsExpireIn: 3600,
		AMQP: &config.AMQPConfig{
			Exchange:      "machinery_exchange",
			ExchangeType:  "direct",
			BindingKey:    "machinery_task",
			PrefetchCount: 3,
		},
	}

	broker := brokers.New(cnf)
	backend := backends.New(cnf)
	lock := eager.New()

	amqpServer := machinery.NewServer(cnf, broker, backend, lock)

	return amqpServer
}

func NewAMQPWoker() *machinery.Server {
	return newAmqp("", "", "", 6738)
}

func NewAMQPProduct() *machinery.Server {
	return newAmqp("", "", "", 6738)
}
