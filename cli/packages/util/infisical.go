package util

import (
	"context"
	"fmt"

	infisical "github.com/infisical/go-sdk"
)

func NewInfisicalClient() {
	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl: INFISICAL_SITE_URL,
	})
	fmt.Println(client)

	//_, err := client.Auth().UniversalAuthLogin()

	//HandleError(err, true)
}
