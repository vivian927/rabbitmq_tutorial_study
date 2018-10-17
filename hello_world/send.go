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
	// 连接到rabbitMQ服务上
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 创建一个管道，用于完成API获取数据。
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 定义一个队列,队列的定义是一个幂等的过程，只有当队列不存在是才会定义
	// 参数： “hello” 队列名；durable：持久；
	
	q, err := ch.QueueDeclare(
		"hello",	// name
		false, 		// durable
		false,		// delete when unused
		false,		// exclusive
		false,		// no-wait
		nil,		// arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 发送消息，消息的内容是一个字节数组，因此发送放可以向其中放入任何内容。
	body := "Hello Vivi"
	err = ch.Publish(
		"",			// exchange
		q.Name,		// 路由键值(routing key)
		false,		// mandatory
		false,		// immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:		 []byte(body),  // 消息
		})
	failOnError(err, "Failed to publish a message")

}
