# Publish Subscribe
发布和订阅：RabbitMQ可以向多个consumer分发同一个message。

RabbitMQ的核心思想：Producer 不需要直接向队列发送任何消息，而且，producer也不知道一个消息是否已经被发送到任何队列。取而代之的是，Producer 可以将消息发送到一个 exchange。

Exchange：exchange的功能比较简单，一边从producers接收消息，另一边来将消息pushes 到相应的队列中，但exchange必须确切的知道每个消息的作用。exchange时通过器exchange type来实行规则，并将消息进行广播的。常见的exchange type 有： direct，topic， headers，fanout。

其中，fanout类型的exchange工作方式比较简单：它就是将其收到的所有消息广播给所有它知道的队列。

	列出server中的exchanges，可以使用下面命令：
		sudo rabbitmqctl list_exchanges

在amqp client代码中，当我们将一个队列命名为空字符串时，我们就相当于使用一个生成的名字(如，amq.gen-JzTY20BRgKO-HjmUJj0wLg）来创建了一个非持久化的队列。
		q,err := ch.QueueDechare (
			"", 		// name
			false,		// durable
			false, 	// delete when usused
			true,		// exclusive
			false, 	// no-wait
			nil,		// arguments
		)

## 1、绑定（Bindings）
绑定格式为：
	
	err = ch.QueueBind(
		q.Name,		// queue name
		"",				// routing key
		"logs",		// exchange
		false,	
		nil,
	)
这样，log exchange就会将消息发送到q 队列中。实际应用中，可以用下面命令来列出所有绑定关系：

	rabbitmqctl list_bindings

exchange 定义的方法如 emit_logs.go所示，在建立了和RabbitMQ之间的连接之后，定义了一个 exchange，注意，这一步必需存在，因为RabbitMQ不允许向一个不存在的exchange发送消息。还需要注意的是，如果没有绑定一个队列到该exchange，会发生消息丢失的现象。

## 2、总结
exchange是一个用于从producer向queue根据exchange type转发消息的中间件。它需要在producer 和consumer双方的程序中定义，并在 consumer中将其定义的队列与要从接收消息的exchange绑定，其他的与前面将的代码相同。在实现了绑定后，就可以由producer产生消息，并向exchange发送消息。然后consumer 定义的队列绑定到exchange上，就可以接收消息，并进一步处理。但在本节中，consumer收到的是producer发出的所有消息，如果只需要接收部分指定的消息，则会在下面小节中介绍。

