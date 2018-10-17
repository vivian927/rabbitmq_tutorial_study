## RabbitMQ 各函数及参数解释

### QueueDeclare
    QueueDeclare (
        string name,
        boolean durable,
        boolean auto-deleted,
        boolean exclusive,
        boolean no-wait,
        Map<string, Object> arguments,
    )

* name: 队列名称
* durable: 是（true）否（false）持久化，队列的声明默认是放在内存中的，如果rabbitmq重启会丢失。如果时持久化，则数据保存到Erlang自带的Mnesia数据库中，当rabbitmq重启会读取该数据库。
* auto-deleted: 是否自动删除，当最后一个消费者断开连接后，该队列是否自动被删除，可以通过RabbitMQ management，查看某个队列的消费者数量，当consumers=0时，队列就会被删除。
* exclusive: 是否排外的，有两个作用，1，当连接关闭时，connection.close()，该队列是否会自动删除；2，该队列是否时私有（private)的。如果不是排外的，可以使用两个消费者都访问同一个队列，没有任何问题，如果是排外的，会对当前队列加锁，其他channel不能访问，强制访问会报异常。如果该值为true，则一个队列只有一个消费者来消费的场景。
* no-wait: 是（true）否非阻塞。阻塞，表示创建交换器的请求发送后，阻塞等带的RMQ server。非阻塞，不会阻塞等待的RMQ server 的返回信息，而RMQ server也不会返回信息（不推荐使用非阻塞）。
* arguments
    * 可以直接写nil，表示没有参数。
    * Message TTL: 设置队列中所有消息的生存周期，也可以在发布消息时单独为某个消息指定剩余生存时间，单位ms。如果生存时间到了，消息会从队列里被删除。如果单独为某个消息设置生存时间，Features=TTL：
        1、AMQP.BasicProperties.Builder properties = new AMQP.BasicProperties().builder().expiration(“6000”); 
        2、channel.basicPublish(EXCHANGE_NAME, “”, properties.build(), message.getBytes(“UTF-8”));
    * Auto Expire: 当队列在指定的时间内没有被访问，就会被删掉。
    * Max length: 限定队列的消息的最大值长度，超过指定长度将会把最早的几条消息删掉。；
    * Max length Bytes: 限定队列最大占用空间的大小，一般受限于内存、磁盘的大小。
    * Dead letter exchange: 当队列的消息长度大于最大长度、或过期等，将从队列删除的消息推送到指定的exchange中，而不是丢弃。
    * Dead letter routing key: 将删除的消息推送到制定个交换机的制定路由键的队列中；
    * Maximum priority: 优先级队列，声明队列时定义最大优先级值（不能太大），在发布消息的时候指定该消息的优先级，优先级更高（数值大）的消息先被消费。 
    * lazy module: laze Queue,看将消息保存到磁盘中，不放在内存中，当消费者开始消费的时候，才加载到内存中；
    * master-locator：

### ExchangeDeclare
    ExchangeDeclare(
        string name,
        string type,
        boolean durable,
        boolean auto-deleted,
        boolean internal,
        boolean no-wait,
        Map<string,Object> arguments,
    )
* name: exchange name.
* type: exchange type，通常有四种：direct，fanout，topic，headers。
* durable: 是否持久化。
* auto-deleted: 是否自动删除。至少有一条绑定才可以触发自动删除，当所有的绑定都与交换器解绑后，会自动删除交换器。
* internal: 是否为内部。客户端无法直接发送msg到内部交换器，只有交换器可以发送message到内部交换器。
* no-wait: 是否阻塞。
* args: nil.

### ExchangeBind
    ExchangeBind(
        string des_exchange,
        string key,
        string source,
        boolean no-wait,
        Map<string, Object> arguments,
    )
* des_exchange: 目标exchange，通常时内部的exchange。
* key: bingd key,表示要绑定的key.
* source: 源交换器。
* no-wait: 是否非阻塞。
* arguments: 参数，nil。

### QueueBind
    QueueBind(
        string queue_name,
        string key,
        string exchange_name,
        boolean no-wait,
        Map<string,object> arguments,
    )
* queue_name: 队列名称。
* key: binding key，绑定的键。
* exchange_name:交换器名称。
* no-wait：是否则塞。
* arguments：参数，nil。

### Publish
    Publish(
        string exchange_name,
        string key,
        boolean mandatory,
        boolean immediate,
        msg Publishing,
    )
    type Publishing struct {
        ContentType string,         // 消息类型，如"text/plain"
        ContentEncoding string,     // 消息的编码，一般默认不写
        DeliveryMode uint8,         // 消息是否持久化，2表示持久化，0或1表示非持久化。
        Body []byte,                // 消息主体。 
        Priority uint8,             // 消息的优先级从0到9
        CorrelationId string,       // correlation id
        ReplyTo string,             // 用于RPC，处理结果返回地址。
        Expiration string,          // 该发送消息的有效期；
        MessageId string,           // message id
        Timestamp time.Time,        // message 的时间戳
        Type string,                // message 到type name。
        UserId string,              // 创建该消息的user id，如 guest。
        AppId string,               // 创建该消息的应用id
    }
* exchange_name：要发送到的exchange 的名称。
* key：routing key。
* mandatory: 建议直接false。表示如果当前消息无法通过exchange匹配到队列时，直接丢弃。如果为true，则会调用basic.return通知生产者。
* immediate：建议直接false，表示消息将一直缓存在队列中等待生产者。如果为true，当消息到达Queue后，发现队列上五消费者，通过basic.Return返回给生产者。
* msg：要发送的消息，对应上面到Publishing结构。

### Consume
    Consume(
        string queue_name,
        string consumer,
        boolean auto-ack,
        boolean exclusive,
        boolean no-local,
        boolean no-wait,
        Map<string,object> arguments
    )
* queue_name: 从该指定名称的队列中订阅消息。
* consumer：消费者到标签，可以为 “”。
* auto-ack：是否自动恢复，告诉服务器客户端已经收到消息。建议设为false，然后手动回复，这样可控性强。
* exclusive：是否排他。表示当前队列只能给一个消费者使用。
* no-local：如果为true，则表示生产者消费者不能是同一个connect。
* no-wait：是否阻塞。
* arguments：nil。

### Qos
    Qos(
        prefetchCount int,
        prefetchSize int,
        global bool,
    )
* prefetchCount: 消费者未确认消息的个数。
* prefetchSize: 消费者为确认消息的大小。
* global: 是否全局生效，true表示是。全局生效表示针对当前connect里的所有channel都生效。

### Get
    Get(
        string queue_name,
        boolean: auto-ack,
    )
* queue_name: 队列名称。
* auto-ack: 是否开启自动回复。

### Ack
    Ack(
        multiple bool
    )
* multiple: true表示回复当前信道所有未回复的ack，用于批量确认。false表示只回复当前条目。

### Reject
    Reject(
        boolean requeue
    )
* requeue: true,RMQ会把这条消息重新加入消息队列，如果为false，则丢弃本条消息。

### Close()
