package queue

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"job-orchestrator/pkg"
	"job-orchestrator/utils"
	"log"
	"os"
)

func SendToSandbox(payload pkg.SandboxPayload) {
	err := godotenv.Load()
	utils.FailOnError(err, "Error loading .env file")

	url := os.Getenv("RABBIT_MQ_LISTENER_URL")
	conn, err := amqp.Dial(url)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open channel")
	defer ch.Close()

	body, err := json.Marshal(payload)
	utils.FailOnError(err, "Failed to marshal sandbox payload")

	err = ch.Publish(
		"", string(pkg.SandboxQueue), false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	utils.FailOnError(err, "Failed to publish sandbox message")

	log.Println("Sent job to sandbox.queue")

	logPayload := map[string]interface{}{
		"event":   "job_dispatched",
		"target":  string(pkg.SandboxQueue),
		"payload": payload,
	}
	logBody, err := json.Marshal(logPayload)
	utils.FailOnError(err, "Failed to marshal logger payload")

	err = ch.Publish(
		"", string(pkg.LoggerQueue), false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        logBody,
		})
	utils.FailOnError(err, "Failed to publish log message")
	log.Println("Sent log to log.queue")
}
