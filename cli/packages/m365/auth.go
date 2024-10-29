package m365

import (
	"context"
	"ehvg/packages/util"
	"encoding/json"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/spf13/cobra"
)

var authCommand = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Microsoft 365",
	Run:   Authenticate,
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
	h, _ := os.UserHomeDir()
	dir := strings.Join([]string{h, ".ehvg"}, string(os.PathSeparator))
	filePath := strings.Join([]string{dir, "m365"}, string(os.PathSeparator))

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0700); err != nil {
			util.HandleError(err, true)
		}

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0700)
		util.HandleError(err, true)

		file.Close()
	}

	return filePath
}

func RetrieveRecord() (azidentity.AuthenticationRecord, error) {
	record := azidentity.AuthenticationRecord{}

	b, err := os.ReadFile(GetAuthFile())

	if err == nil {
		if json.Valid(b) {
			err = json.Unmarshal(b, &record)
		}
	}

	return record, err
}

func StoreRecord(record azidentity.AuthenticationRecord) error {
	b, err := json.Marshal(record)
	util.HandleError(err, true)

	err = os.WriteFile(GetAuthFile(), b, 0700)

	return err
}

func init() {
	m365Command.AddCommand(authCommand)
}
