package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"encoding/json"
	"flashy-product/datamodels"
	"flashy-product/services"
	"sync"
)

//const MQURL = "amqp://flashyuser:flashyuser@172.31.96.59:5672/flashy"
const MQURL = "amqp://flashyuser:flashyuser@127.0.0.1:5672/flashy"

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	Mqurl     string
	sync.Mutex
}

func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabb"+
		"itmq!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// if consist
		false,
		// if deleted
		false,
		// uniqueness
		false,
		// block
		false,
		// extra property
		nil,
	)
	if err != nil {
		return err
	}
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// If true, the message will be returned to the sender if no qualified queue can be found according to its own exchange type and routekey rules
		false,
		// If true, when the exchange sends a message to the queue and finds that there is no consumer on the queue, it will return the message to the sender
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return nil
}

// Consumers in simple mode
func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	// 1. Apply for the queue, if the queue does not exist, it will be created automatically, if it exists, the creation will be skipped
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		// consistance
		false,
		// auto delete
		false,
		// unique
		false,
		// block
		false,
		// extra
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//消费者流控
	// customer block
	r.channel.Qos(
		// The maximum number of messages that the current consumer can accept at one time
		1,
		// The maximum capacity passed by the server in octets
		0,
		// If set to true, available to channel
		false,
	)

	//receive message
	msgs, err := r.channel.Consume(
		q.Name, // queue
		// Used to distinguish between multiple consumers
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}

			err = productService.SubNumberOne(message.ProductID)
			if err != nil {
				fmt.Println(err)
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
