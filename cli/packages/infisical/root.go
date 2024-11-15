package infisical

import (
	"context"

	ifc "github.com/infisical/go-sdk"
	"github.com/spf13/cobra"
)

var infisicalCommand = &cobra.Command {
  Use: "infisical",
  Short: "Infisical Configuration Management Tool",
  Run: func(cmd *cobra.Command, args []string) {},
}

func Execute(rootCmd *cobra.Command) {
  infisicalCommand.AddCommand(
    NewSetEnvironmentCommand(),
  )
  rootCmd.AddCommand(
    infisicalCommand,
  )
}

func GetClient() (ifc.InfisicalClientInterface, error) {
  client := ifc.NewInfisicalClient(context.Background(), ifc.Config{
    SiteUrl: INFISICAL_SITE_URL,
  })

  _, err := client.Auth().UniversalAuthLogin(DOK_INFISICAL_CLIENT_ID, DOK_INFISICAL_CLIENT_SECRET)
  if err != nil {
    return nil, err
  }

  return client, nil
}

func CompletionListSites() []string {
  return []string{"dokteronline", "dok", "seemenopause", "smnp"}
}