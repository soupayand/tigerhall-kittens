package worker

import (
	"github.com/IBM/sarama"
	"gopkg.in/gomail.v2"
	"strings"
	"tigerhall-kittens/logger"
)

const (
	group = "emailgroup"
	topic = "email_animal_sighted"
)

func StartConsumer() {
	brokerList := []string{"localhost:9092"} // Add your Kafka broker addresses

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		logger.LogError(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			logger.LogError(err)
		}
	}()
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		logger.LogError(err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			logger.LogError(err)
		}
	}()

	logger.LogInfo("Kafka worker started...")

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			emailIDs := strings.Split(string(msg.Value), ",")
			for _, emailID := range emailIDs {
				if err := sendEmail(emailID); err != nil {
					logger.LogError(err)
				} else {
					logger.LogInfo("Email sent to %s", emailID)
				}
			}
		case err := <-partitionConsumer.Errors():
			logger.LogError(err)
		}
	}
}

// Replace mailtrap smtp credentials
func sendEmail(emailID string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "from@example.com")
	m.SetHeader("To", emailID)
	m.SetHeader("Subject", "Animal Sighting Update")
	m.SetBody("text/plain", "An animal you reported has been sighted again")

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "$username", "$password") // Replace with your SMTP server and credentials
	return d.DialAndSend(m)
}
