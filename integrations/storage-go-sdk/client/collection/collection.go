package collection

import "github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/client"

type collection struct {
	database string
	name     string

	client *client.Client
}
