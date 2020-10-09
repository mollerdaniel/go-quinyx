# go-quinyx
[![GoDoc](https://godoc.org/github.com/mollerdaniel/go-quinyx/quinyx?status.svg)](https://godoc.org/github.com/mollerdaniel/go-quinyx/quinyx)

go-quinyx is a Go client library for accessing the [Quinyx REST API](https://api.quinyx.com/v2/docs/swagger-ui.html) 

## Example

Quickstart:

```go
 ctx := context.Background()

// Oauth2 Config
conf := clientcredentials.Config{
	ClientID:     os.Getenv("CLIENTID"),
	ClientSecret: os.Getenv("CLIENTSECRET"),
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
	log.Fatalf("Error: %v RequestUID: %s", err, res.GetQuinyxUID())
}

// Print each category
for _, category := range categories {
	fmt.Println(category)
}
  ```

## client_id and client_secret
Get the client_id and client_secret from [Quinyx](https://www.quinyx.com).

>! Note only client_secrets generated after *2020-10-09* is compatible with the [Oauth2](https://github.com/golang/oauth2) module