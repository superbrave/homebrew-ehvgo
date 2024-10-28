package util

func GetInfisicalClientAuth(projectContext string) (InfisicalClientAuth, error) {
	var infisicalClientAuth InfisicalClientAuth
	var err error

	switch pc := projectContext; pc {
	case DOKTERONLINE_SLUG:
		infisicalClientAuth.ClientId = DOK_INFISICAL_CLIENT_ID
		infisicalClientAuth.ClientSecret = DOK_INFISICAL_CLIENT_SECRET
		infisicalClientAuth.ProjectId = DOK_INFISICAL_PROJECT_ID

		err = nil
	}

	return infisicalClientAuth, err
}

func GetInfisicalmage(version string) string {
	return INFISICAL_CLI_IMAGE + ":" + version
}
