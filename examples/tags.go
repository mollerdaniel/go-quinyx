package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/mollerdaniel/go-quinyx/quinyx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	ctx := context.Background()

	// Oauth2 Config
	conf := clientcredentials.Config{
		ClientID:     os.Getenv("CLIENTID"),
		ClientSecret: os.Getenv("CLIENTSECRET"), // Quinyx API does not accept URLEncoded secrets https://tools.ietf.org/html/rfc6749#section-2.3.1
		TokenURL:     "https://api-rc.quinyx.com/v2/oauth/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
		EndpointParams: url.Values{
			"grant_type": {"client_credentials"},
		},
	}

	// HTTP Client
	client := conf.Client(ctx)

	// Quinyx API Client
	q, err := quinyx.NewClient(client, quinyx.String("https://api-rc.quinyx.com"))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Call Tags Service to get all categories
	categories, res, err := q.Tags.GetAllCategories(ctx)
	if err != nil {
		if res != nil {
			log.Fatalf("Error: %v RequestUID: %s", err, res.QuinyxUID)
		}
		log.Fatalf("Error: %v", err)
	}

	// Dump each category to console
	for _, tagCategory := range categories {
		fmt.Println(tagCategory)
	}
}
