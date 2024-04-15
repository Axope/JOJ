package rabbitmq

import (
	"context"
	"fmt"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
)

func InitMQ() error {
	cfg := configs.GetRBTConfig()
	log.LoggerSugar.Debugf("cfg = %+v", cfg)
	username := cfg.Username
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
	log.Logger.Debug(url)

	var err error
	if conn, err = amqp.Dial(url); err != nil {
		return err
	}
	if ch, err = conn.Channel(); err != nil {
		return err
	}
	if queue, err = ch.QueueDeclare(
		"publisher",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	return nil
}

func SendMsgByJson(msg []byte) error {
	if err := ch.PublishWithContext(
		context.TODO(),
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         msg,
		}); err != nil {
		return err
	}
	return nil
}

func SendMsgByProtobuf(msg []byte) error {
	if err := ch.PublishWithContext(
		context.TODO(),
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/x-protobuf",
			Body:         msg,
		}); err != nil {
		return err
	}
	return nil
}
