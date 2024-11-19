package infisical

import (
	"context"
	"errors"
	"os"
	"strings"

	ifc "github.com/infisical/go-sdk"
)

var (
  DokClientID string 
  DokClientSecret string

  SmnpClientID string
  SmnpClientSecret string

  EhvgClientID string
  EhvgClientSecret string
)

func GetClient(a string) (ifc.InfisicalClientInterface, error) {
  client := ifc.NewInfisicalClient(context.Background(), ifc.Config{
    SiteUrl: INFISICAL_SITE_URL,
  })

  id, err := getInfisicalClientId(a)
  if err != nil {
    return nil, err
  }

  secret, err := getInfisicalClientSecret(a)
  if err != nil {
    return nil, err
  }
  _, err = client.Auth().UniversalAuthLogin(id, secret)
  if err != nil {
    return nil, err
  }

  return client, nil
}

func getInfisicalClientId(a string) (string, error) {
  if id := os.Getenv("INFISICAL_CLIENT_ID"); id != "" {
    return id, nil
  }

  switch strings.ToLower(a) {
  case "dok", "dokteronline":
    return DokClientID, nil
  case "smnp", "seemenopause":
    return SmnpClientID, nil
  case "ehvg", "ehealthventuresgroup":
    return EhvgClientID, nil
  default:
    return "", errors.New("no client ID found")
  }
}

func getInfisicalClientSecret(a string) (string, error) {
  if id := os.Getenv("INFISICAL_CLIENT_SECRET"); id != "" {
    return id, nil
  }
  
  switch strings.ToLower(a) {
  case "dok", "dokteronline":
    return DokClientSecret, nil
  case "smnp", "seemenopause":
    return SmnpClientSecret, nil
  case "ehvg", "ehealthventuresgroup":
    return EhvgClientSecret, nil
  default:
    return "", errors.New("no client secret found")
  }
}