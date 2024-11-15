package util

const (
	DOKTERONLINE_SLUG = "dok"
	SEEMENOPAUSE_SLUG = "smnp"

	EHVGO_VERSION = "0.0.1"

	AZURE_APPLICATION_ID = "e388aa60-cb81-437f-a14e-bad2974ea418"
	AZURE_TENANT_ID      = "be8e47c5-fe9b-49b6-a09b-050ee2a44ec0"
	AZURE_REDIRECT_URL   = "http://localhost:8000"
	AZURE_TENANT_URI     = "https://login.microsoftonline.com/be8e47c5-fe9b-49b6-a09b-050ee2a44ec0"
)

func GetAzureScopes() []string {
  return []string{"User.Read.All"}
}