package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err);
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// 接收数据的consumer，需要一直跑着，监听队列中的消息，并取出进行下一步处理

func main () {
	// 1、打开一个和消息队列的连接，需要与publisher连接的是一个 rabbitmq服务
	// 2、打开一个管道
	// 3、定义一个队列，这是因为consumer需要在publisher之前启动，因此需要保证消费者可以从队列中获得消息。
	// 上面的这三步都需要跟要订阅的publisher 匹配。
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err,"Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare (
		"hello",		// name
		false,			// durable
		false,			// delete when unused
		false,			// excuisive
		false,			// no-wait
		nil,			// arguments
	)
	failOnError(err, "Failed to declear a queue")

	// 消费者将通知amqp服务通过队列向自己传递消息；
	// amqp 是以异步的方式向消费者发送消息；
	// 消费者开启一个goroutine 来从 管道读消息
	msgs, err := ch.Consume(
		q.Name,			// queue
		"",				// consumer
		true,			// auto-ack
		false,			// exclusive
		false,			// no-local
		false,			// no-wait
		nil,			// args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)  // 定义了一个bool类型的双向 管道。

	go func () {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for message. To exit press CTRL+C")
	
	// 同步，等待goroutine结束，就在通知main goroutine任务完成。
	<-forever
}