package main

import (
	"context"
	"fmt"
	"os"

	client "github.com/zinclabs/sdk-go-zincsearch"
)

func main() {
	user := *client.NewMetaUser() // MetaUser | User data
	userName := "admin"
	password := "admin"
	user.Name = &userName
	user.Password = &password
	configuration := client.NewConfiguration()
	apiClient := client.NewAPIClient(configuration)

	data := *client.NewMetaIndexSimple() // MetaIndexSimple | Index data
	resp, r, err := apiClient.Index.Create(context.Background()).Data(data).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `Index.Create``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Create`: MetaHTTPResponseIndex
	fmt.Fprintf(os.Stdout, "Response from `Index.Create`: %v\n", resp)
}
