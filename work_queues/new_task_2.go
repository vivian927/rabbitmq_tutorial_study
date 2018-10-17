package main

import (
	"os"
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"strings"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg,err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	// 1.连接 rabbitmq服务
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect rabitmq")

	// 2.创建管道
	ch, err := conn.Channel()
	failOnError(err, "Failed to create channel.")

	// 3.定义队列
	q, err := ch.QueueDeclare(
		"new_task2",			// name
		true,					// durable
		false,					// deleted when unused
		false,					// exclusive
 Body: []byte(body),
		})
	failOnError(err,"Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}

func bodyForm(args []string) string {
	var s string
	if(len(args) < 2) || os.Args[1] == "" {
		s = "Hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}