package es

import (
	"context"
	"time"

	"github.com/olivere/elastic/v7"
)

type Client struct {
	*elastic.Client
	*elastic.BulkProcessor
}

func NewClient(opt ...elastic.ClientOptionFunc) (*Client, error) {
	cli, err := elastic.NewClient(opt...)
	if err != nil {
		return nil, err
	}

	proc, err := cli.BulkProcessor().
		BulkSize(2 << 20).
		Workers(1000).
		BulkActions(1000).
		FlushInterval(100 * time.Millisecond).
		Do(context.Background())

	if err != nil {
		return nil, err

	}

	return &Client{cli, proc}, nil
}

func NewClientSimple(servers []string, user, pwd string) (*Client, error) {
	return NewClient(elastic.SetURL(servers...), elastic.SetBasicAuth(user, pwd))
}
