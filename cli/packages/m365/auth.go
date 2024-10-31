package m365

import (
	"context"
	"ehvg/packages/util"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var credentials *azidentity.InteractiveBrowserCredential

var loginCommand = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Microsoft 365",
	Run:   Authenticate,
}

var logoutCommand = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Microsoft 365",
	Run:   DestroySession,
}

var whoamiCommand = &cobra.Command{
	Use:   "whoami",
	Short: "Shows the current logged in user",
	Run:   Whoami,
}

func Whoami(cmd *cobra.Command, args []string) {
	record, err := RetrieveAuthenticationRecord()
	if err == nil {
		fmt.Printf("You are currently logged in as: %v, on tenant %v", record.Username, record.TenantID)
	}
}

func DestroySession(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(GetAuthFile()); os.IsNotExist(err) {
		color.New(color.FgHiYellow).Println("No Microsoft 356 session found.")

		return
	} else {
		if err := os.Remove(GetAuthFile()); err != nil {
			util.HandleError(err, true)
		}
	}

	color.New(color.FgHiGreen).Println("Successfully logged out from Microsoft 365")
}

func Authenticate(cmd *cobra.Command, args []string) {
	r, err := RetrieveAuthenticationRecord()
	if err != nil {
		color.New(color.FgHiRed).Printf("Unable to retrieve authentication record: %v", err)
		return
	}

	cache, err := cache.New(nil)
	if err != nil {
		color.New(color.FgHiRed).Printf("Unable instantiate persistent cache: %v", err)
	}

	credentials, err = azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		Cache:                cache,
		TenantID:             util.AZURE_TENANT_ID,
		ClientID:             util.AZURE_APPLICATION_ID,
		RedirectURL:          util.AZURE_REDIRECT_URL,
		AuthenticationRecord: r,
	})
	if err != nil {
		color.New(color.FgHiRed).Printf("Unable to start authorization: %v", err)
		return
	}

	if r == (azidentity.AuthenticationRecord{}) {
		color.New(color.FgHiWhite).Println("Starting authentication")
		r, err = credentials.Authenticate(context.TODO(), nil)
		if err != nil {
			color.New(color.FgHiRed).Printf("Unable to authenticate: %v", err)
			return
		}

		err = StoreAuthenticationRecord(r)
		if err != nil {
			color.New(color.FgHiRed).Printf("Unable store authentication: %v", err)
			return
		}
	}
}

func GetAuthFile() string {
	filePath := strings.Join([]string{util.GetConfigDir(), "m365"}, string(os.PathSeparator))

	return filePath
}

func GetAndCreateAuthFile() string {
	if _, err := os.Stat(util.GetConfigDir()); os.IsNotExist(err) {
		if err := os.Mkdir(util.GetConfigDir(), 0700); err != nil {
			util.HandleError(err, true)
		}
	}

	file, err := os.OpenFile(GetAuthFile(), os.O_CREATE|os.O_RDWR, 0700)
	util.HandleError(err, true)

	file.Close()

	return GetAuthFile()
}

func RetrieveAuthenticationRecord() (azidentity.AuthenticationRecord, error) {
	record := azidentity.AuthenticationRecord{}

	b, err := os.ReadFile(GetAndCreateAuthFile())
	util.HandleError(err, false)

	if len(b) > 0 {
		dc, err := util.Decrypt(b)

		if err == nil && json.Valid(dc) {
			_ = json.Unmarshal(dc, &record)
		}
	}

	return record, err
}

func StoreAuthenticationRecord(record azidentity.AuthenticationRecord) error {
	b, err := json.Marshal(record)
	util.HandleError(err, true)

	e, err := util.Encrypt(b)
	util.HandleError(err, true)

	err = os.WriteFile(GetAndCreateAuthFile(), e, 0700)

	return err
}

// func NewGraphClient(scopes []string) (*msgraphsdk.GraphServiceClient, error) {
// 	creds, err := RetrieveRecord()
// 	if err != nil {
// 		color.New(color.FgHiRed).Println("Invalid credentials provided")
// 	}

// 	creds.GetToken()

// 	authProvider, err := azure.NewAzureIdentityAuthenticationProviderWithScopes(creds, scopes)
// }

func init() {
	m365Command.AddCommand(loginCommand)
	m365Command.AddCommand(logoutCommand)
	m365Command.AddCommand(whoamiCommand)
}
