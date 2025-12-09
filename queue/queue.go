package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/streadway/amqp"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/repo"
	"github.com/vviveksharma/auth/models"
)

// Task types for routing
const (
	TaskTypeCreateMessage = "create_message"
	TaskTypeSendEmail     = "send_email"
	TaskTypeCleanupTokens = "cleanup_tokens"
	// Add more task types as needed
)

// Generic task wrapper
type Task struct {
	Type    string      `json:"type"`    // Task type for routing
	Payload interface{} `json:"payload"` // Actual data
}

type IQueueService interface {
	Connect() (conn *amqp.Connection, err error)
	DeclareQueue(conn *amqp.Connection) (qu amqp.Queue, err error)
	PublishMessage(qu amqp.Queue, conn *amqp.Connection, message interface{}) error
	ConsumeMessages(conn *amqp.Connection, queue amqp.Queue) error
	CreateMessageRequest(value interface{}) error
	PublishMessageTask(qu amqp.Queue, conn *amqp.Connection, message models.DBMessage) error
}

type QueueService struct{}

func NewQueueRequest() (IQueueService, error) {
	return &QueueService{}, nil
}

func (q *QueueService) Connect() (conn *amqp.Connection, err error) {
	host := os.Getenv("QUEUE_HOST")
	port := os.Getenv("QUEUE_PORT")
	addr := host + ":" + port
	conn, err = amqp.Dial("amqp://user:password@" + addr + "/")
	if err != nil {
		log.Println("error while connecting to RabbitMQ: ", err)
	}
	fmt.Println("Connected to RabbitMQ Successfully!")
	return conn, err
}

func (q *QueueService) DeclareQueue(conn *amqp.Connection) (qu amqp.Queue, err error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while creating a channel: ", err)
		return amqp.Queue{}, err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"task-queue",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	fmt.Println("Declared Queue Successfully!")
	return queue, nil
}

// PublishMessage - Generic publish that accepts any type
func (q *QueueService) PublishMessage(qu amqp.Queue, conn *amqp.Connection, message interface{}) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while creating a channel: ", err)
		return err
	}
	defer ch.Close()

	// Convert message to JSON bytes
	var body []byte
	switch msg := message.(type) {
	case string:
		body = []byte(msg)
	case []byte:
		body = msg
	default:
		// Marshal any struct/map to JSON
		body, err = json.Marshal(msg)
		if err != nil {
			log.Println("error while marshaling message to JSON: ", err)
			return fmt.Errorf("failed to marshal message: %v", err)
		}
	}

	err = ch.Publish(
		"",      // exchange
		qu.Name, // routing key (queue name)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make messages persistent
		},
	)
	if err != nil {
		log.Println("error while publishing the queue message: ", err)
		return err
	}
	log.Printf("✅ Published message to queue: %s", qu.Name)
	return nil
}

// ConsumeMessages - Processes messages based on task type
func (q *QueueService) ConsumeMessages(conn *amqp.Connection, queue amqp.Queue) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while creating a channel for consumer: ", err)
		return err
	}
	defer ch.Close()

	// QoS - process one message at a time
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Println("error while setting QoS: ", err)
		return err
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer tag
		false,      // auto-ack (manual acknowledgment)
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Println("error while registering consumer: ", err)
		return err
	}

	log.Println("🚀 Consumer started. Waiting for messages...")

	// Process messages in a blocking loop
	for msg := range msgs {
		log.Printf("📨 Received message: %s", string(msg.Body))
		var task Task
		err := json.Unmarshal(msg.Body, &task)

		var processErr error

		if err == nil && task.Type != "" {
			processErr = q.routeTask(task)
		} else {
			var dbMsg models.DBMessage
			if jsonErr := json.Unmarshal(msg.Body, &dbMsg); jsonErr == nil {
				processErr = q.CreateMessageRequest(dbMsg)
			} else {
				processErr = fmt.Errorf("unable to parse message: %v", jsonErr)
			}
		}

		// Acknowledge or reject message
		if processErr != nil {
			log.Printf("❌ Error processing message: %v", processErr)
			msg.Nack(false, true)
		} else {
			log.Printf("✅ Message processed successfully")
			msg.Ack(false)
		}
	}

	log.Println("⚠️  Consumer stopped - channel closed")
	return nil
}

// routeTask - Routes tasks based on type
func (q *QueueService) routeTask(task Task) error {
	log.Printf("🔀 Routing task type: %s", task.Type)

	switch task.Type {
	case TaskTypeCreateMessage:
		// Convert payload to DBMessage
		payloadBytes, err := json.Marshal(task.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %v", err)
		}

		var dbMsg models.DBMessage
		if err := json.Unmarshal(payloadBytes, &dbMsg); err != nil {
			return fmt.Errorf("failed to unmarshal to DBMessage: %v", err)
		}

		return q.CreateMessageRequest(dbMsg)

	case TaskTypeSendEmail:
		// Future: Handle email sending
		log.Println("📧 Email task received (not yet implemented)")
		return nil

	case TaskTypeCleanupTokens:
		// Future: Handle token cleanup
		log.Println("🧹 Cleanup task received (not yet implemented)")
		return nil

	default:
		return fmt.Errorf("unknown task type: %s", task.Type)
	}
}

// CreateMessageRequest - Creates a message in the database
func (q *QueueService) CreateMessageRequest(value interface{}) error {
	iMessage, err := repo.NewMessageRepository(db.DB)
	if err != nil {
		log.Printf("❌ Error connecting to message repository: %v", err)
		return fmt.Errorf("error while connecting to the message repo layer: %v", err)
	}

	var messageRequest models.DBMessage

	switch v := value.(type) {
	case models.DBMessage:
		messageRequest = v
	case *models.DBMessage:
		messageRequest = *v
	default:
		// Fallback: try JSON unmarshal
		payloadBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %v", err)
		}

		if err := json.Unmarshal(payloadBytes, &messageRequest); err != nil {
			return fmt.Errorf("unable to typecast to DBMessage: %v", err)
		}
	}

	// Create in database
	err = iMessage.Create(&messageRequest)
	if err != nil {
		log.Printf("❌ Error creating message in database: %v", err)
		return fmt.Errorf("error while creating the message entry in the database: %v", err)
	}

	log.Println("✅ Message entry created successfully in database")
	return nil
}

// Helper functions for publishing specific task types

// PublishMessageTask - Publishes a message creation task
func (q *QueueService) PublishMessageTask(qu amqp.Queue, conn *amqp.Connection, message models.DBMessage) error {
	task := Task{
		Type:    TaskTypeCreateMessage,
		Payload: message,
	}
	return q.PublishMessage(qu, conn, task)
}

// PublishEmailTask - Publishes an email task (for future use)
func (q *QueueService) PublishEmailTask(qu amqp.Queue, conn *amqp.Connection, emailData interface{}) error {
	task := Task{
		Type:    TaskTypeSendEmail,
		Payload: emailData,
	}
	return q.PublishMessage(qu, conn, task)
}

// PublishCleanupTask - Publishes a cleanup task (for future use)
func (q *QueueService) PublishCleanupTask(qu amqp.Queue, conn *amqp.Connection, cleanupData interface{}) error {
	task := Task{
		Type:    TaskTypeCleanupTokens,
		Payload: cleanupData,
	}
	return q.PublishMessage(qu, conn, task)
}
