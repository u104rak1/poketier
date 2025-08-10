package cors_test

import (
	"poketier/pkg/cors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCORSConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName             string
		appEnv               string
		allowOrigins         string
		wantAllowAll         bool
		wantAllowOrigins     []string
		wantAllowCredentials bool
	}{
		{
			caseName:             "正常系: production環境ではAllowOriginsが設定されること",
			appEnv:               "production",
			allowOrigins:         "https://example.com,https://foo.com",
			wantAllowAll:         false,
			wantAllowOrigins:     []string{"https://example.com", "https://foo.com"},
			wantAllowCredentials: true,
		},
		{
			caseName:             "正常系: local環境ではAllowAllOriginsがtrueになること",
			appEnv:               "local",
			allowOrigins:         "",
			wantAllowAll:         true,
			wantAllowOrigins:     nil,
			wantAllowCredentials: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			got := cors.GetCORSConfig(tt.allowOrigins, tt.appEnv)

			// Assert
			assert.Equal(t, tt.wantAllowAll, got.AllowAllOrigins, "AllowAllOrigins should match")
			assert.Equal(t, tt.wantAllowOrigins, got.AllowOrigins, "AllowOrigins should match")
			assert.Equal(t, tt.wantAllowCredentials, got.AllowCredentials, "AllowCredentials should match")
		})
	}
}
