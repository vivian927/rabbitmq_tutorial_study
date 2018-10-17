package main

import (
	"time"
	"fmt"
	"log"
	"bytes"
	
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	// 1、连接Rabbitmq服务
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	// 2、创建管道
	ch, err := conn.Channel()
	failOnError(err, "Failed to declare channel")

	// 3、定义队列
	q, err := ch.QueueDeclare(
		"task_queue", 			// name
		true,					// durable
		false,					// deleted when unused
		false,					// exclusive
		false,					// no-wait
		nil,					// args
	)
	failOnError(err, "Failed to declare a queue")

	// 4、订阅
	msgs, err := ch.Consume(
		q.Name,				// name
		"",					// consumer
		true,				// auto-ack
		false,				// exclusive
		false,				// no-local
		false,				// no-wait
		nil,				// args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	// 5、读取消息
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")
		}
	}()
	log.Printf(" [*] Waiting for message. To exit press CTRL+C")
	<-forever
}