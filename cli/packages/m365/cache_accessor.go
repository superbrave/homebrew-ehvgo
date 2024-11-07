package m365

import (
	"context"
	"ehvg/packages/util"
	"encoding/json"
	"os"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/fatih/color"
)

type TokenCache struct {
	file string
}

func (t *TokenCache) Replace(ctx context.Context, cache cache.Unmarshaler, hints cache.ReplaceHints) error {
	b, err := os.ReadFile(t.file)
	if err != nil {
		color.New(color.FgHiRed).Printf("Unable open authfile: %v", err)
		return cache.Unmarshal([]byte{})
	}

	if len(b) > 0 {
		dc, err := util.Decrypt(b)
		if err != nil {
			color.New(color.FgHiRed).Printf("Decryption failed: %v", err)
			return err
		}
		if !json.Valid(dc) {
			color.New(color.FgHiRed).Printf("Invalid JSON object: %v", err)
			return err
		}

		return cache.Unmarshal(dc)
	}

	return cache.Unmarshal([]byte{})
}

func (t *TokenCache) Export(ctx context.Context, cache cache.Marshaler, hints cache.ExportHints) error {
	data, err := cache.Marshal()
	if err != nil {
		color.New(color.FgHiRed).Printf("Unable marshal session info: %v", err)
		return err
	}

	edata, err := util.Encrypt(data)
	if err != nil {
		color.New(color.FgHiRed).Printf("Unable to encrypt session info: %v", err)
		return err
	}

	return os.WriteFile(t.file, edata, 0600)
}
