package worker

import (
	"fmt"
	"github.com/IBM/sarama"
	"gopkg.in/gomail.v2"
	"log"
	"strings"
)

func main() {
	brokerList := []string{"localhost:9092"} // Add your Kafka broker addresses
	topic := "email_animal_sighted"

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Error closing Kafka consumer: %v", err)
		}
	}()
	partitionList, err := consumer.Partitions("email_animal_sighted")
	if err != nil {
		fmt.Println(fmt.Sprintf("error while get Partitions %v", err))
		return
	}
	partitionConsumer, err := consumer.ConsumePartition(topic, partitionList[0], sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("Error closing partition consumer: %v", err)
		}
	}()

	log.Println("Kafka worker started...")

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			emailIDs := strings.Split(string(msg.Value), ",")
			for _, emailID := range emailIDs {
				if err := sendEmail(emailID); err != nil {
					log.Printf("Error sending email to %s: %v", emailID, err)
				} else {
					log.Printf("Email sent to %s", emailID)
				}
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Kafka consumer error: %v", err)
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

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "", "") // Replace with your SMTP server and credentials
	return d.DialAndSend(m)
}
