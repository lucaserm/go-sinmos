package auth

import (
	"testing"

	"github.com/lucaserm/go-sinmos/internal/json"
)

// A malformed refresh token previously reached uuid.MustParse in the service
// and panicked. Validation must now reject it before the service is called.
func TestRefreshTokenPayloadRejectsMalformedUUID(t *testing.T) {
	cases := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{name: "empty", token: "", wantErr: true},
		{name: "garbage", token: "not-a-uuid", wantErr: true},
		{name: "valid", token: "0192f9c2-3a1b-7c4d-8e5f-6a7b8c9d0e1f", wantErr: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := json.Validate.Struct(RefreshTokenPayload{RefreshToken: tc.token})
			if (err != nil) != tc.wantErr {
				t.Fatalf("token %q: got err=%v, wantErr=%v", tc.token, err, tc.wantErr)
			}
		})
	}
}
