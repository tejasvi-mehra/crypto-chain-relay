package token

import (
	"builder-router/internal/framework/logger"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/caarlos0/env/v6"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	jwkjwt "github.com/lestrrat-go/jwx/v2/jwt"
)

var pkJson = `{
		"kid": "1",
		"d": "WmwFAu0rBzhD0KXt81VLgiNt1d6HpCk-YkM_nAxHnAUnCVmV6DzV8_iGnQJeM_KgSMBP3jhxqOTJBjazAitLx8Zp41uKipRrdy730SEkBKE-0XHs09TLylGiupJzNKSxHvf0KfKhfHO-HOLxf_KJkHS2WwBeUZruCZqxxMXOkiKYHDOZFOmQzib_fGC_C52flQhKXEr5M0kV_IqZr9ICPPQN5rjl6jMCj-WakC5QxIERsFl6VRIGXejclHqHayB_qj_gyLtT1lIAsNMUIKMjxkd5nd0_vI7fLhpXbyl-3Sqi7zTGRMocCizgBgJZmUnJMELPyVo02CE9QkHA9How-Q",
		"dp": "NlzZFADI3_-3GLKM0GbOlI77wmavDfdfZofI87BB41i8htjkGjNl8V9YQREU4r29fxV-_EKAcjOJZMm-QYKkKkQnXHjeEk2-NEss-906AwrJDpB2UnQXVoPVsrPnK9l8zHnjubfs9X7qGkUrF_fmI1jH5ox75vUPP12LgAWW69U",
		"dq": "jzgIy5zCd9KxWtuZZky6yQU3KIUyIIMyOihjwhDFgnLx5I3QVgqyxJ1XRn65t4PCJ_MY9ff0-7zCR4-UldBKv7ReZrctLYX5QJAO4rZ07rAI0t_RxuZaOK0bHjnO04hVNmzvFzGXqZr8m-ar_IAHy7gCcTnhgVlks51MiCF1tPE",
		"e": "AQAB",
		"kty": "RSA",
		"n": "pjB_13eFWfOlB8B3pSlx7oDpXaKfUF8hPJc9D-xEzwc3MrvG3aF6f89Bs6XZ0pqdidUv4CBG_6iZMwr-PXiywh2ekzcFg7ujjml4qAC55kdFs6Vp0kPbLZlQif6xF8NGkfpAFPR3SO-xUbj4e4NgwW1VGf6Lfm_LMVsZ2ycYBYdDaSq5P0W46mMqA5nKTqV3c22FxCCHrzxkC9hzDXXpDWJ3jkLCI5Uvib3gSQ4M5gp0hqoFrWbMDPNwgB1l1pE7LhENRYtAsRsoRMe9-1ojcLvk2GouTBfQ3nDEhZbyqZevr4kGC8wv8RRvmQxBZQB3n5f_vl5LjMNEhnuLHaaYVQ",
		"p": "zwRsqQIvUHugzkCjyeiZ6vXM6Oukjsb9QvSZFh9SEyN9cz5PmQ6cBW8FoP92uu3x-F51LqP-SPJ1KTDsvZTOLpftrueNsAditiaieyOUTvf_qb9G4OMy4btW_1btlLsUHE147bLO9pE79tN7_FsB3PAxBK0uq29Ewke9j705lnM",
		"q": "zYMG4ZVyJbId8Lygiwr-A61QNW48NZIQZauYvmwcRKekS97BSdN6GXh3TMHlTziCdvL0y6zV6iCLNlA-oQIcVvzOPkVvufmrpR0ZWdjw6tFctCG_MmcNT_5ROgH03Um5O7IvpRC4Lo9wGPNDrv50YniGCcE8u8R8A8ckTqQEnBc",
		"qi": "y5-JyiRLnrGlky9KRwAWLCoBpKnejXcpDuKxowj2RPn5Msily4EC8RDaat4HLa7ZSQU9s3JQ_LbWB8SPDJGMqE697KGd_fnmxIl7udZdzMbLHv9abnjFmW1FFuksLIQAEuqRtITT3nBplfUd2DUBnNZzbaH7kyrCO48PubIXP2g"
	  }`

var pubkKeys = json.RawMessage(`{
		"keys": [
			{
				"kid": "1",
				"e": "AQAB",
				"kty": "RSA",
				"n": "pjB_13eFWfOlB8B3pSlx7oDpXaKfUF8hPJc9D-xEzwc3MrvG3aF6f89Bs6XZ0pqdidUv4CBG_6iZMwr-PXiywh2ekzcFg7ujjml4qAC55kdFs6Vp0kPbLZlQif6xF8NGkfpAFPR3SO-xUbj4e4NgwW1VGf6Lfm_LMVsZ2ycYBYdDaSq5P0W46mMqA5nKTqV3c22FxCCHrzxkC9hzDXXpDWJ3jkLCI5Uvib3gSQ4M5gp0hqoFrWbMDPNwgB1l1pE7LhENRYtAsRsoRMe9-1ojcLvk2GouTBfQ3nDEhZbyqZevr4kGC8wv8RRvmQxBZQB3n5f_vl5LjMNEhnuLHaaYVQ"
			}
		]
	  }`)

func GenerateJwt(tokenType string, accountID string, expAfter time.Duration) (string, error) {
	// this function should not be here but instead should be more secure
	pk, err := jwk.ParseKey([]byte(pkJson))
	if err != nil {
		return "", err
	}

	builder := jwkjwt.NewBuilder().
		Audience([]string{"testing"}).
		Issuer("issuer").
		JwtID("0987654321").
		Subject(accountID).
		Claim("typ", tokenType)

	if expAfter != 0 {
		builder = builder.Expiration(time.Now().Add(expAfter))
	}

	t, err := builder.Build()
	if err != nil {
		return "", err
	}

	signed, err := jwkjwt.Sign(t, jwkjwt.WithKey(jwa.RS256, pk))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

type JwkConfig struct {
	JwkURL     string `env:"JWK_URL,notEmpty" envDefault:""`
	CaCertPath string `env:"CA_CERT_PATH" envDefault:""`
}

func NewJwkConfig(options ...env.Options) (*JwkConfig, error) {
	c := &JwkConfig{}
	if err := env.Parse(c, options...); err != nil {
		return c, fmt.Errorf("fail to load jwk config: %s", err)
	}
	return c, nil
}

func Jwk(skipper middleware.Skipper, jwkConfig *JwkConfig) echo.MiddlewareFunc {
	log := logger.GetLogger("jwk")

	// var client *http.Client
	// jwks, err := keyfunc.Get(jwkConfig.JwkURL, keyfunc.Options{
	// 	// Refreshes only based on setting the RefreshInterval & RefreshRateLimit
	// 	RefreshErrorHandler: func(err error) {
	// 		log.Panicw("Failed with JWK KeyFunc during refresh", "error", err)
	// 	},
	// 	Client:            client,
	// 	RefreshUnknownKID: true,
	// })
	jwks, err := keyfunc.NewJSON(pubkKeys)
	if err != nil {
		log.Panicw("Failed to create JWKS from url", "error", err)
	}
	// Error handling for invalid JWK configuration with keys during startup
	if len(jwks.ReadOnlyKeys()) == 0 {
		log.Panicw("Failed to get the keys from JWK url", "jwk", jwks.KIDs())
	}

	return echojwt.WithConfig(echojwt.Config{
		Skipper: skipper,
		KeyFunc: func(t *jwt.Token) (interface{}, error) {
			return jwks.Keyfunc(t)
		},

		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			token, err := jwt.NewParser(
				jwt.WithExpirationRequired(),
			).Parse(auth, jwks.Keyfunc)
			if err != nil {
				log.Errorw("failed to parse token", "err", err)
				return nil, &echojwt.TokenError{Token: token, Err: err}
			}
			if !token.Valid {
				log.Errorw("invalid token", "token", token)

				return nil, &echojwt.TokenError{Token: token, Err: errors.New("invalid token")}
			}

			return token, nil
		},
	})
}
