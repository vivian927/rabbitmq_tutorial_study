package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	// 1、连接到 RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect RabbitMQ server")
	defer conn.Close()

	// 2、打开一个管道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 3、定义一个 Exchange
	err = ch.ExchangeDeclare(
		"logs",   // exchange name
		"fanout", // exchange type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a exchange")

	// 4、定义一个 队列
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // deleted when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 5、将队列绑定到 exchange上
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange name
		false,
		nil,
	)
	failOnError(err, "Failed to bindi to exchange")

	// 6、订阅消息
	msgs, err := ch.Consume(
		q.Name, // 队列名称， 这里并没有为队列设置名字，是一个生成到队列名
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to consume the queue")

	// 7、接收消息
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [x] Waiting for logs. To exit press CTRL+C.")
	<-forever

}
