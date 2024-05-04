package rabbitmq

import (
	"context"
	"fmt"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	judgeSendConn  *amqp.Connection
	judgeSendCh    *amqp.Channel
	judgeSendQueue amqp.Queue
)

func InitJudgeSender(cfg configs.RabbitmqConfig) error {
	username := cfg.Username
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
	log.Logger.Debug(url)

	var err error
	if judgeSendConn, err = amqp.Dial(url); err != nil {
		return err
	}
	if judgeSendCh, err = judgeSendConn.Channel(); err != nil {
		return err
	}
	if judgeSendQueue, err = judgeSendCh.QueueDeclare(
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
func SendMsgByProtobuf(msg []byte) error {
	if err := judgeSendCh.PublishWithContext(
		context.TODO(),
		"",
		judgeSendQueue.Name,
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

// func SendMsgByJson(msg []byte) error {
// 	if err := judgeSendCh.PublishWithContext(
// 		context.TODO(),
// 		"",
// 		judgeSendQueue.Name,
// 		false,
// 		false,
// 		amqp.Publishing{
// 			DeliveryMode: amqp.Persistent,
// 			ContentType:  "application/json",
// 			Body:         msg,
// 		}); err != nil {
// 		return err
// 	}
// 	return nil
// }

var (
	judgeResultRecvConn  *amqp.Connection
	judgeResultRecvCh    *amqp.Channel
	judgeResultRecvQueue amqp.Queue
)

func InitJudgeResultRecvier(cfg configs.RabbitmqConfig) (<-chan amqp.Delivery, error) {
	username := cfg.Username
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
	log.Logger.Debug(url)

	var err error
	if judgeResultRecvConn, err = amqp.Dial(url); err != nil {
		return nil, err
	}
	if judgeResultRecvCh, err = judgeResultRecvConn.Channel(); err != nil {
		return nil, err
	}
	if judgeResultRecvQueue, err = judgeResultRecvCh.QueueDeclare(
		"JudgeResponseQueue",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	msgs, err := judgeResultRecvCh.Consume(
		judgeResultRecvQueue.Name,
		"",
		false, // auto ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func InitMQ() (<-chan amqp.Delivery, error) {
	cfg := configs.GetRBTConfig()
	log.LoggerSugar.Debugf("cfg = %+v", cfg)

	if err := InitJudgeSender(cfg); err != nil {
		return nil, err
	}
	return InitJudgeResultRecvier(cfg)
}
