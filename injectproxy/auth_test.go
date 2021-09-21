// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package injectproxy

import (
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/auth"
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	sampleKey = `eyJhbGciOiJSUzI1NiIsImtpZCI6ImIyYzVjOTkxYWMwMDJlZmYwMGZhNTMzYmM2Y2U0YWFhODc4MDBkMTcifQ.eyJpc3MiOiJodHRwOi8vZGV4LWxkYXAtdW1icmVsbGE6NTU1NiIsInN1YiI6IkNnWmtZV2x6ZVdRU0JHeGtZWEEiLCJhdWQiOiJhZXRoZXItcm9jLWd1aSIsImV4cCI6MTYzMjI5OTc1NywiaWF0IjoxNjMyMjEzMzU3LCJub25jZSI6IlNUUTNXSEZ1Y1doSE5FOUdMa3huVVVZelIydFlSRmhPV0UxRlgyVkdkVVZEZW0wMVZEQkJRbk5QWlZsVSIsImF0X2hhc2giOiI1UHRtWW93MEI3RVJnemJ6cENzUG1nIiwiY19oYXNoIjoiYjB5VndIdHJXcDVnRWdXeUZ3TFhIZyIsImVtYWlsIjoiZGFpc3lkQG9wZW5uZXR3b3JraW5nLm9yZyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsiY2hhcmFjdGVyc0dyb3VwIiwiRW50ZXJwcmlzZUFkbWluIiwic3RhcmJ1Y2tzIl0sIm5hbWUiOiJEYWlzeSBEdWtlIn0.aUH7MChqylA53r-D7tHY34XbAYbWH2zzhqbVtgYka9H1-OTw0QtjKRSLTL-1dBQjfr6oinO7lq0p-FgKQIwP7H5aYrrjjn1wZMMdCNWS2qX_JH7RJ-Eqo2FFXKiqQZNKGj5o5IUIQk2n_f1DYITuen7WBIPRIH-s-fo82pczzPaqixs32nagbK9H6btkl-G-io3mlpn8N8_F-nrTcaLvM4IK4zvrrvP9Qt0t9x7bLf76T660JrWXOFnIhOKmIpURi0I_07vf2tyCqrncleUqQ6DVOUZy68fIVhY0sKwznyRWpdr3Sd51EffcyIg1Ewp0XNtyu2m-NlMc1Q7DnD-l2w`
	sampleOIDCserver   = "http://testme"
)

const sampleWellKnownBody = `{
  "issuer": "http://testme",
  "authorization_endpoint": "http://testme/auth",
  "token_endpoint": "http://testme/token",
  "jwks_uri": "http://testme/keys",
  "userinfo_endpoint": "http://testme/userinfo",
  "device_authorization_endpoint": "http://testme/device/code",
  "grant_types_supported": [
    "authorization_code",
    "refresh_token",
    "urn:ietf:params:oauth:grant-type:device_code"
  ],
  "response_types_supported": [
    "code"
  ],
  "subject_types_supported": [
    "public"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256"
  ],
  "code_challenge_methods_supported": [
    "S256",
    "plain"
  ],
  "scopes_supported": [
    "openid",
    "email",
    "groups",
    "profile",
    "offline_access"
  ],
  "token_endpoint_auth_methods_supported": [
    "client_secret_basic",
    "client_secret_post"
  ],
  "claims_supported": [
    "iss",
    "sub",
    "aud",
    "iat",
    "exp",
    "email",
    "email_verified",
    "locale",
    "name",
    "preferred_username",
    "at_hash"
  ]
}`

const sampleKeysBody = `{
  "keys": [
    {
      "use": "sig",
      "kty": "RSA",
      "kid": "63f315fb92c7d43dd1d39dcd3d25c03d1744f449",
      "alg": "RS256",
      "n": "uBBCqtYmO91_CybA3UH8BIcekU8NRIZaH94KV8Hgw06_NZjCE_V2wvbNsf79dJ6huPo8GUvGP0bO-Vso0YeAOwSQ9hvXEcu80Tjd3TG83tihTsywF_8J_F74fjhhzwCo1k1X9AVFzFap_-byphwKDvVAQlXC5mxLPtZEYpTjfuOonCE54UFtmMNbn0nSt6XVe912-Btod8wJYu_1x5UZBkmDJWRppsLhr9DxqrP7x1LJJL-K2exatrBpBZLknh3vTDPAV97R6Aac11HLMJZjAATzSGKYh6daEM2KJBIisLBw46ZeqYQp3aYQqK-V3DwIkM2zDx9N0nJtTD73lc8ryQ",
      "e": "AQAB"
    },
    {
      "use": "sig",
      "kty": "RSA",
      "kid": "b2c5c991ac002eff00fa533bc6ce4aaa87800d17",
      "alg": "RS256",
      "n": "3Klr8U6ODE7t9Vsuq0TglbnPbrtooEEfUhB8YYv04dxsDnd4EJ1798IYZlO1ITbZiDMXUddDE_gMZuTqFoOdLwVcTKY8MJ1h65wdBL4P6LC3QlLdLyyDoR7vxYKD8-FWkes00jjNkYyw8QaRVAdQDlRXb2wDS0mOk00cEgw662S5508NnbNhTXRkaxBX-N9gxaw-ZtopXAuUFsAv6qW77C1sUrQoA637dMHAtR9nBhMZaGOfOWSTgFCpgpRi40gl6Q6ZYhBI7acNgKI45h4GyAAZ8TW2BM7aFR45WRGtujBJK87NBm-3X2qd5xuvNepqV3BoA97hvcGLvIrhx1z8kw",
      "e": "AQAB"
    },
    {
      "use": "sig",
      "kty": "RSA",
      "kid": "d8e480505f388fd4626d5af71dea0ebb7795381c",
      "alg": "RS256",
      "n": "3fC5CD1L8Pju6IEKQWGhIrL6PiSZ0pRzmMY0pfC6Xrpy1qfmxyg80dJ3Qqs5-3vohCb1W42aT6HeMXqy89TXOEua9s70KtVAuuLtFuMp1C0KiDi6Du7Ff_0kbE66zx7t9jXhRGxqu2kHwJlVcBSUODArzDNiG4s59zsMM_jEyTTNuPPwRH1PB0fxqTzBStSyrbRfIO7ARE3IZQCIx-6Vdo-nrlg2SIq2qAPE7h7lCxYLz5wTHoFciBebfGN3NoyGOrDWEKnoEhSTPmH_qCTFa_zSxBJNLT4rmHb5byAtxT5kYnFfBOF5vkMLa9uFt6APHlEYL5xhGI0PBh7E7c_KEw",
      "e": "AQAB"
    }
  ]
}`

func Test_enforceAuth(t *testing.T) {

	tsKeys := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sampleKeysBody)
	}))
	defer tsKeys.Close()

	tsWellKnown := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, strings.ReplaceAll(sampleWellKnownBody, sampleOIDCserver + "/keys", tsKeys.URL))
	}))
	defer tsWellKnown.Close()

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sampleKey))

	os.Setenv(auth.OIDCServerURL, tsWellKnown.URL)
	defer os.Setenv(auth.OIDCServerURL, "")

	w := httptest.NewRecorder()

	groups := enforceAuth(w, req)
	assert.Equal(t, 3, len(groups))
	for _, g := range groups {
		switch g {
		case "starbucks", "charactersGroup", "EnterpriseAdmin":
		default:
			t.Fatalf("unexpected group %s", g)
		}
	}
}

