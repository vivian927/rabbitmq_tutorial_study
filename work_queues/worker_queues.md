# workers queue

工作队列，又名任务队列。用于多个在多个wokers中分配任务，其目标为：

*  立即执行且必须等待其执行完成的任务占用过多的资源，可以将其放在最后执行；
*  只要将任务作为消息发送到队列中即可；
*  工作在后台的worker进程只要从队列（queue）中pop出一个任务（或message），并执行即可；
*  当有多个worker时，队列中的这些tasks将会在它们之间共享；

默认情况下，RabbitMQ都会将message发送到下一个Consumer中。平均每个consumer会收到相同数量的message。rabbitmq的分配任务的方式为 round-robin。尤其是当有3+个consumer时。

## 1、消息确认（ack）

如果队列中存在一个较耗时的任务，当它被启动，并在执行完之谦，worker宕掉了。rabbitmq是否会知道，那个任务是否还保留在队列中？在new_task.go 和 worker.go 代码中。在一个task被Rabitmq分发给一个worker之后，会被立即从队列中删除。如果该任务在执行的过程中，worker死掉了，那么该任务包括rabitmq分发给该woker的所有task都会丢失。

为了防止上述情况的发生，RabbitMQ提供了 “消息 ack“(message acknowledgments) 机制。当RabbitMQ服务器收到了来自一个consumer的ack消息时，表明其分发给该consumer的某个任务（message）已经被接收并处理了，RabbitMQ可以将其删除。

如果一个consumer死掉了(如它的管道关闭了、rabbitmq connect关闭了，或TCP丢失连接)，RabbitMQ server 这边没有收到其发过来的ack，这样RabbitMQ这边就可以将该消息重新插入到队列中，以待重新分发。通过这样的方式，可以保证不会有消息丢失.

RabbitMQ不存在消息超时的情况。当有 Consumer进程死掉，RabbitMQ会重复消息，这样益于需要耗时较长的task。

消息确认可以设置为自动发送，在订阅时给auto-ack参数传 true，则表示在接收到消息并处理结束后，自动发送ack确认消息。也可以手动发送，这时需要给auto-ack参数传递 false，并在接收到消息后，调用d.Ack（false）来发送确认消息。详情见 worker_ack.go代码。

注意：如果忘记发送 ack 确认消息，则会造成严重的后果。当consumer死掉，则messager回重新传递消息，虽然不会影响消息的传递和处理，但RabbitMQ回因为没有任何接收的消息被释放掉而消耗掉来越多的内存。该类型的调试方法：
	
	sudo rabbitmqctl list_queues name messages_ready messages_unacknowledged
	// 使用rabbitmqctl来打印出 messages_unacknowledged字段
	Windows系统中：
	rabbitmqctl.bat list_queues name messages_ready messages_unacknowledged

## 2、Message durability
ack确认消息只能保证当worker宕机时不发生任务丢失问题。但却不能保证RabbitMQ服务停止，任务丢失情况的发生。因此，如果想要让 RabbitMQ服务在停止服务或宕机后仍能记住之前其保存的messages，则需要让程序通过某些设定来告诉它要记录。因此，就需要进行下面操作：
	
	1、在producer 和consumer代码中的队列的定义时，为参数durable传true值；
	2、在amqp.Publishing中为标记DeliveryMode 为 amqp.Persistent，即：
			amqp.Publishing {
				DeliveryMode: amqp.Persistent,
				ContentType: "text/plain",
				Body: []byte(body),
			}
## 3、按照任务量公平地分发任务
RabbitMQ会按照消息进入队列的顺序进行分发任务，而不会关注consumer未返回确认消息的数量。即忙目的对第n个consumer派发第n个task。

因此，可以通过设置 prefetch count的值来设置RabbitMQ每次为每个Worker分发的最多消息数。如果设为1，即表示如果没有烧到某个worker在处理的消息的ack，就不要想起分发新的message。设置方法如下：
	
	err := ch.Qos (
		1,			// prefetch count
		0,			// prefetch size
		false,		// global
	)

注意队列大小和 prefetch count 的关系，当所有的Worker都在忙，且不断的有消息进到队列中来，导致队列已经满了，这时就需要 增加 workers 的数量和采用其他策略来处理这个问题了。

2和3的用法如new_task_2.go和worker2.go中所示。







