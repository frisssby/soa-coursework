package event

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"stats/model"

	"github.com/IBM/sarama"
)

func RunConsumer(topic string, db *sql.DB) {
	consumer, err := sarama.NewConsumer([]string{os.Getenv("KAFKA_URI")}, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err.Error())
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to create partition consumer: %v", err.Error())
	}
	defer partitionConsumer.Close()

	for {
		msg, ok := <-partitionConsumer.Messages()
		if !ok {
			log.Println("Partition consumer channel has been closed")
			return
		}

		var message model.Event
		json.Unmarshal(msg.Value, &message)

		tableName := topic
		_, err = db.Exec(
			fmt.Sprintf(`INSERT INTO %v (user_id, task_id, author_id) VALUES ($1, $2, $3)`, tableName),
			message.UserID,
			message.TaskID,
			message.AuthorID,
		)
		if err != nil {
			log.Printf("Error executing a query: %v", err)
		}
	}
}
