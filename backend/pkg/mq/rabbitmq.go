// Package mq RabbitMQ 消息队列基础封装
package mq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const defaultConfirmTimeout = 5 * time.Second

// MessageRecord 消息发送记录
type MessageRecord struct {
	ID         string
	Exchange   string
	RoutingKey string
	Body       []byte
	Status     int8 // 0-待确认 1-确认成功 2-确认失败
	CreatedAt  time.Time
}

// RabbitMQ 连接管理器
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewRabbitMQ 创建 RabbitMQ 连接
func NewRabbitMQ(host string, port int, user, password, vhost string) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", user, password, host, port, vhost)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("连接 RabbitMQ 失败: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建 RabbitMQ Channel 失败: %w", err)
	}

	// 启用发布确认模式
	if err := ch.Confirm(false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("启用 Publisher Confirm 失败: %w", err)
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

// GetChannel 获取 Channel
func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.ch
}

// DeclareExchange 声明交换机（持久化）
func (r *RabbitMQ) DeclareExchange(name, kind string) error {
	return r.ch.ExchangeDeclare(
		name,
		kind,
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	)
}

// DeclareQueue 声明队列（持久化）
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.ch.QueueDeclare(
		name,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
}

// BindQueue 绑定队列到交换机
func (r *RabbitMQ) BindQueue(queueName, routingKey, exchangeName string) error {
	return r.ch.QueueBind(queueName, routingKey, exchangeName, false, nil)
}

// Publish 发布持久化消息（不等待确认，快速返回）
func (r *RabbitMQ) Publish(exchange, routingKey string, body []byte) error {
	return r.ch.Publish(
		exchange,
		routingKey,
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}

// PublishWithConfirm 发布消息并等待 Broker 确认（ack/nack 或超时）
func (r *RabbitMQ) PublishWithConfirm(ctx context.Context, exchange, routingKey string, body []byte) error {
	confirms := r.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	if err := r.Publish(exchange, routingKey, body); err != nil {
		return fmt.Errorf("发布消息失败: %w", err)
	}

	// 移除监听以避免 channel 泄漏
	defer func() {
		// 清空 channel 中尚未处理的通知
		for range confirms {
		}
	}()

	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			return fmt.Errorf("消息被 RabbitMQ 拒绝 (Nack)")
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("等待 RabbitMQ 确认超时: %w", ctx.Err())
	case <-time.After(defaultConfirmTimeout):
		return fmt.Errorf("等待 RabbitMQ 确认超时 (%v)", defaultConfirmTimeout)
	}
}

// Consume 消费消息（手动 ACK）
func (r *RabbitMQ) Consume(queueName, consumerTag string) (<-chan amqp.Delivery, error) {
	return r.ch.Consume(
		queueName,
		consumerTag,
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
}

// NotifyConfirm 获取 Publisher Confirm 通知 Channel
func (r *RabbitMQ) NotifyConfirm() chan amqp.Confirmation {
	return r.ch.NotifyPublish(make(chan amqp.Confirmation, 1))
}

// Close 关闭连接
func (r *RabbitMQ) Close() {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
