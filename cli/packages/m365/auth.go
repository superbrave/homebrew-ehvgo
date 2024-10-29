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
	record, err := RetrieveRecord()
	if err == nil {
		fmt.Printf("You are currently logged in as: %v, on tenant %v", record.Username, record.TenantID)
	}
}

func DestroySession(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(GetAuthFile()); os.IsNotExist(err) {
		color.New(color.FgYellow).Println("No Microsoft 356 session found.")

		return
	} else {
		if err := os.Remove(GetAuthFile()); err != nil {
			util.HandleError(err, true)
		}
	}

	color.New(color.FgGreen).Println("Successfully logged out from Microsoft 365")
}

func Authenticate(cmd *cobra.Command, args []string) {
	record, err := RetrieveRecord()
	util.HandleError(err, true)

	cache, err := cache.New(nil)
	util.HandleError(err, false)

	creds, err := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		AuthenticationRecord: record,
		Cache:                cache,
	})
	util.HandleError(err, true)

	if (record == azidentity.AuthenticationRecord{}) {
		record, err = creds.Authenticate(context.Background(), nil)
		util.HandleError(err, true)

		if err := StoreRecord(record); err != nil {
			util.HandleError(err, true)
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

func RetrieveRecord() (azidentity.AuthenticationRecord, error) {
	record := azidentity.AuthenticationRecord{}

	b, err := os.ReadFile(GetAndCreateAuthFile())

	if err == nil && json.Valid(b) {
		err = json.Unmarshal(b, &record)
	}

	return record, err
}

func StoreRecord(record azidentity.AuthenticationRecord) error {
	b, err := json.Marshal(record)
	util.HandleError(err, true)

	err = os.WriteFile(GetAndCreateAuthFile(), b, 0700)

	return err
}

func init() {
	m365Command.AddCommand(loginCommand)
	m365Command.AddCommand(logoutCommand)
	m365Command.AddCommand(whoamiCommand)
}
