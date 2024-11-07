package util

const (
	INFISICAL_SITE_URL = "https://infisical.ehealthsystems.nl"

	DOK_INFISICAL_CLIENT_ID           = "42265328-925d-4987-b7de-500c4748926a"
	DOK_INFISICAL_CLIENT_SECRET       = "cf84b5e39ccc32ecf6d0d202a8e6e89cb71d795b89faa204d7714663fd8b8a95"
	DOK_CHECKOUT_INFISICAL_PROJECT_ID = "15ee529b-6602-4d88-ab67-c755b498b4a3"

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