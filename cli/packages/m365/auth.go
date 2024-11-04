package m365

import (
	"context"
	"ehvg/packages/util"
	"errors"
	"os"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	msal "github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cacheAccessor = &TokenCache{file: GetAndCreateAuthFile()}

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

var authStorage = &TokenCache{
	file: GetAndCreateAuthFile(),
}

func Whoami(cmd *cobra.Command, args []string) {
	a, err := Account()
	if err != nil {
		color.New(color.FgHiWhite).Printf("Unable to retrieve account: %v", err)
		return
	}

	color.New().Printf("You are currently logged in as %v on tenant ID %v", a.PreferredUsername, a.Realm)
}

func Account() (msal.Account, error) {
	ctx := context.Background()

	client, err := msal.New(util.AZURE_APPLICATION_ID, public.WithAuthority(util.AZURE_TENANT_URI), public.WithCache(authStorage))
	if err != nil {
		return msal.Account{}, err
	}

	accounts, err := client.Accounts(ctx)
	if err != nil {
		return msal.Account{}, err
	}

	if len(accounts) == 0 {
		err := errors.New("no active session found")
		return msal.Account{}, err
	}

	return accounts[0], err
}

func DestroySession(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	client, err := msal.New(util.AZURE_APPLICATION_ID, public.WithAuthority(util.AZURE_TENANT_URI), public.WithCache(authStorage))
	if err != nil {
		color.New(color.FgHiRed).Printf("Unable to initalize MSAL: %v", err)
		return
	}

	a, err := Account()
	if err != nil {
		color.New(color.FgHiRed).Printf("Unable to retrieve account: %v", err)
		return
	}

	if err := client.RemoveAccount(ctx, a); err != nil {
		color.New(color.FgHiRed).Printf("Unable to logout: %v", err)
		return
	}

	color.New(color.FgHiGreen).Println("Successfully logged out from Microsoft 365")
}

func Authenticate(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	isForced, _ := cmd.Flags().GetBool("force")

	client, err := msal.New(util.AZURE_APPLICATION_ID, public.WithAuthority(util.AZURE_TENANT_URI), public.WithCache(authStorage))
	if err != nil {
		color.New(color.FgHiRed).Printf("Unable to initalize MSAL: %v", err)
	}

	accounts, err := client.Accounts(ctx)
	if err != nil {
		color.New(color.FgHiRed).Printf("Something went wrong while retrieving session: %v", err)
		return
	}

	if len(accounts) == 0 || isForced {
		_, err := client.AcquireTokenInteractive(ctx, []string{"User.Read"})
		if err != nil {
			color.New(color.FgHiRed).Printf("Unable to initalize authentication: %v", err)
		}
	} else {
		_, err := client.AcquireTokenSilent(ctx, []string{"User.Read"}, public.WithSilentAccount(accounts[0]))
		if err != nil {
			color.New(color.FgHiRed).Printf("Unable to initalize authentication: %v", err)
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

func init() {
	loginCommand.Flags().BoolP("force", "", false, "Force a new login, overwriting any existing sessions")
	m365Command.AddCommand(loginCommand)
	m365Command.AddCommand(logoutCommand)
	m365Command.AddCommand(whoamiCommand)
}
