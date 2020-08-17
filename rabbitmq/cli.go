package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type RabbitOption func(op *option)

func WithAddr(addr string) RabbitOption {
	return func(op *option) {
		op.Addr = addr
	}
}

type option struct {
	Addr     string
	Name     string
	Exchange string
}

type Client struct {
	option
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewClient(ops ...RabbitOption) (*Client, error) {
	opt := option{}

	for _, op := range ops {
		op(&opt)
	}

	conn, err := amqp.Dial(opt.Addr)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare("", false, true, false, false, nil)
	if err != nil {
		return nil, err
	}

	opt.Name = q.Name

	res := &Client{option: opt, conn: conn, channel: ch}

	return res, nil
}

func (c *Client) Bind(exchange string) error {
	err := c.channel.QueueBind(c.Name, "", exchange, false, nil)
	c.Exchange = exchange

	return err
}

func (c *Client) Send(queue string, body interface{}) error {
	str, err := json.Marshal(body)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{ReplyTo: c.Name, Body: str}
	err = c.channel.Publish("", queue, false, false, msg)

	return err
}

func (c *Client) Publish(exchange string, body interface{}) error {
	str, err := json.Marshal(body)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{ReplyTo: c.Name, Body: str}
	err = c.channel.Publish(exchange, "", false, false, msg)

	return err
}

func (c *Client) Consume() (<-chan amqp.Delivery, error) {
	ch, err := c.channel.Consume(c.Name, "", true, false, false, false, nil)

	return ch, err
}

func (c *Client) Close() {
	c.conn.Channel()
}
