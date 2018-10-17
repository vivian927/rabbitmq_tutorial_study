package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)
/*
	在建立了和RabbitMQ之间的连接之后，定义了一个 exchange，注意，这一步必需存在，因为RabbitMQ不允许向一个不存在的exchange发送消息。
	如果没有绑定一个队列到该exchange，会发生消息丢失的现象。
	在此程序中，如果没有consumer订阅，丢失消息也不会造成影响。
*/

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s",msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main () {
	// 1、连接 RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect RabbitMQ server")
	defer conn.Close()

	// 2、打开管道
	ch, err := conn.Channel()
	failOnError(err,"Failed to open a channel")
	defer ch.Close()

	// 3、定义 Exchange
	err = ch.ExchangeDeclare(
		"logs",			// name
		"fanout",		// exchange type
		true,			// durable
		false,			// auto-deleted
		false,			// internal
		false,			// no-wait
		nil,			// arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// 4、处理输入的参数，并将其发送出去
	body := bodyForm(os.Args)
	err = ch.Publish(
		"logs",			// exchange
		"",				// routing key
		false,			// mandatory
		false,			// imediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] sent %s", body)

}

func bodyForm(args []string) string {
	var s string
	if(len(args) <2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}