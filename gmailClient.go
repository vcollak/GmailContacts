package main

import (
	"golang.org/x/net/context"
	"google.golang.org/api/gmail/v1"
)

func getGmailClient() (*gmail.Service, error) {
	ctx := context.Background()

	config, err := getGmailConfig()
	if err != nil {
		return nil, err
	}

	client := getClient(ctx, config)
	svc, err := gmail.New(client)
	if err != nil {
		return nil, err
	}

	return svc, nil

}
