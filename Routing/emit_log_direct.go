package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", err, msg))
	}
}

func main() {
	// 1、连接到RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect RabbitMQ server.")

	// 2、打开管道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel.")

	// 3、定义exchange
	err = ch.ExchangeDeclare(
		"logs_direct",			// exchange name
		"direct",				// tyoe
		true,					// durable
		false,					// auto-deleted
		false,					// internal
		false,					// no-wait
		nil,					// arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// 4、推送消息到Rabbitmq server
	body := bodyFrom(os.Args)
	err = ch.Publish(
		"logs_direct",			// exchange name
		serverityFrom(os.Args),	// routing key
		false,					// mandatory
		false,					// immediate
		amqp.Publishing{
			ContentType:	"text/plain",
			Body:			[]byte(body),
		})
	failOnError(err, "Failed to publish message to exchange")
	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func serverityFrom(args []string) string {
	var s string
	if(len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}