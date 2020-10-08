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
	// TODO
	ctx := context.Background()
	urlValues := url.Values{}
	urlValues.Set("grant_type", "client_credentials")
	conf := clientcredentials.Config{
		ClientID:     os.Getenv("CLIENTID"),
		ClientSecret: os.Getenv("CLIENTSECRET"), // Quinyx API does not accept URLEncoded secrets https://tools.ietf.org/html/rfc6749#section-2.3.1
		TokenURL:     "https://api-test.quinyx.com/v2/oauth/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
		EndpointParams: url.Values{
			"grant_type": {"client_credentials"},
		},
	}
	client := conf.Client(ctx)

	q, err := quinyx.NewClient(client, quinyx.String("https://api-test.quinyx.com"))
	if err != nil {
		log.Fatalf("Error creating a Client: %v", err)
	}

	categories, res, err := q.Tags.GetAllCategories(ctx)
	if err != nil {
		if res != nil {
			log.Fatalf("Error: %v RequestUID: %s", err, res.QuinyxUID)
		}
		log.Fatalf("Error: %v", err)
	}

	for _, tagCategory := range categories {
		fmt.Println(tagCategory)
	}
}
