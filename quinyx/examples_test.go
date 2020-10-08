package quinyx_test

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/mollerdaniel/go-quinyx/quinyx"
	"golang.org/x/oauth2/clientcredentials"
)

func ExampleClient() {
	ctx := context.Background()

	urlValues := url.Values{}
	urlValues.Set("grant_type", "client_credentials")

	conf := clientcredentials.Config{
		ClientID:       "CLIENTID",
		ClientSecret:   "CLIENTSECRET", // Quinyx API does not accept URLEncoded secrets https://tools.ietf.org/html/rfc6749#section-2.3.1
		TokenURL:       "https://api.quinyx.com/v2/oauth2/token",
		EndpointParams: urlValues,
	}
	client := conf.Client(ctx)

	q, err := quinyx.NewClient(client, nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
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
