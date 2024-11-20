build-cli:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -C cli -o ../dist/ehvgo-${VERSION}-${GOOS}-${GOARCH} -v -ldflags "\
		-X ehvg/packages/infisical.DokClientID=${DOK_INFISICAL_CLIENT_ID} \
		-X ehvg/packages/infisical.DokClientSecret=${DOK_INFISICAL_CLIENT_SECRET} \
		-X ehvg/packages/infisical.SmnpClientID=${SMNP_INFISICAL_CLIENT_ID} \
		-X ehvg/packages/infisical.SmnpClientSecret=${SMNP_INFISICAL_CLIENT_SECRET} \
		-X ehvg/packages/infisical.EhvgClientID=${EHVG_INFISICAL_CLIENT_ID} \
		-X ehvg/packages/infisical.EhvgClientSecret=${EHVG_INFISICAL_CLIENT_SECRET} \
		-X ehvg/packages/util.EhvgoVersion=${VERSION} \
	"