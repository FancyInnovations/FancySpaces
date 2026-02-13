package collection

import "github.com/fancyinnovations/fancyspaces/storage-sdk/client"

type collection struct {
	database string
	name     string

	client *client.Client
}
