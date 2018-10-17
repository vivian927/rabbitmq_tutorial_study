package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/streadway/amqp"
)
// 使用：
//		1、go run new_task.go First message.
//		2、go run new_task.go First message..
//      3、go run new_task.go First message...
//		4、go run new_task.go First message....
//		5、...
//并启动两个worker程序

// work queues : 工作队列，又名任务队列。用于在多个workers中分配耗时的任务。功能：
// 1、立即执行且必须等待其执行完成的任务占用过多的资源，可以将其放在最后执行；
// 2、将任务作为消息发送到队列中；
// 3、工作在后台的worker进程只要从队列（queue）中pop出一个任务，并执行即可。
// 4、当有多个worker时，队列中的这些tasks将会在它们之间共享；

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err);
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main () {
	// 连接mq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect amqp server")

	// 打开amqp的管道
	ch, err := conn.Channel()
	failOnError(err, "Failed to declare channel")

	// 定义队列
	q, err := ch.QueueDeclare(
		"task_queue",			// name
		true,					// durable
		false,					// delete when unused
		false,					// exclusive
		false,					// no-waitar
		nil,					// arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 发送消息
	body := bodyForm(os.Args)
	err = ch.Publish(
		"",			// exchange
		q.Name,		// queue name
		false,		// mandatory
		false,		// immediate
		amqp.Publishing {
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body: []byte(body),
		})
	failOnError(err, "Failed to publish to queue")

	log.Printf("[x] %s sent: ",body)
}
func bodyForm(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
