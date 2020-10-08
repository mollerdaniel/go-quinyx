package main

import (
	"context"
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
		ClientID:       os.Getenv("CLIENTID"),
		ClientSecret:   os.Getenv("CLIENTSECRET"),
		TokenURL:       "https://api-test.quinyx.com/v2/oauth/token?grant_type=client_credentials",
		AuthStyle:      oauth2.AuthStyleInHeader,
		EndpointParams: urlValues,
	}
	client := conf.Client(ctx)

	q, err := quinyx.NewClient(client, quinyx.String("https://api-test.quinyx.com/"))
	if err != nil {
		log.Fatalf("Error creating a Client: %v", err)
	}

	_, res, err := q.Tags.GetAllTags(ctx, "myexternalid")
	if err != nil {
		if res != nil {
			log.Fatalf("Error: %v RequestUID: %s", err, res.QuinyxUID)
		}
		log.Fatalf("Error: %v", err)
	}
}
