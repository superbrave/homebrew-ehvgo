package util

type InfisicalClientAuth struct {
	ClientId     string
	ClientSecret string
	ProjectId    string
}

const (
	DOK_INFISICAL_CLIENT_ID     = "42265328-925d-4987-b7de-500c4748926a"
	DOK_INFISICAL_CLIENT_SECRET = "cf84b5e39ccc32ecf6d0d202a8e6e89cb71d795b89faa204d7714663fd8b8a95"
	DOK_INFISICAL_PROJECT_ID    = "15ee529b-6602-4d88-ab67-c755b498b4a3"

	DOKTERONLINE_SLUG = "dok"
	SEEMENOPAUSE_SLUG = "smnp"

	INFISICAL_CLI_IMAGE = "ghcr.io/superbrave/infisical-cli"
)
