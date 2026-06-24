package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/auth"
)

// AdminTokenCmd gets Keycloak administrator access token
type AdminTokenCmd struct {
	Username        string `help:"Admin username" default:"admin" env:"KEYCLOAK_USERNAME,KC_BOOTSTRAP_ADMIN_USERNAME"`
	Password        string `help:"Admin password" default:"admin" env:"KEYCLOAK_PASSWORD,KC_BOOTSTRAP_ADMIN_PASSWORD"`
	Realm           string `help:"Keycloak realm" default:"master" env:"KEYCLOAK_REALM"`
	SetEnv          bool   `help:"Set tokens to environment variables" default:"true"`
	AccessTokenEnv  string `help:"Set access token to this environment variable" default:"KEYCLOAK_ACCESS_TOKEN"`
	RefreshTokenEnv string `help:"Set refresh token to this environment variable" default:"KEYCLOAK_REFRESH_TOKEN"`
}

func (c *AdminTokenCmd) Run(ctx *kong.Context, cli *CLI) error {
	log.Logger.Info().Str("pkg", PKG_CMD).Msg("Getting token")
	log.Logger.Debug().Str("pkg", PKG_CMD).Msgf("Token command options: %+v; context: %+v", cli, ctx.Args)

	authService := auth.New()
	tokenCtx, cancel := cli.CreateContextWithTimeout()
	defer cancel()
	tr, err := authService.PasswordToken(tokenCtx, cli.KeycloakBaseURL, c.Realm, c.Username, c.Password)
	if err != nil {
		return err
	}

	// set tokens to env vars if requested
	if c.SetEnv {
		tokens := []struct{ key, value string }{
			{c.AccessTokenEnv, tr.AccessToken},
			{c.RefreshTokenEnv, tr.RefreshToken},
		}
		for _, token := range tokens {
			if err := authService.SetEnvToken(token.key, token.value, ".env"); err != nil {
				return err
			}
		}
		return nil
	}

	fmt.Println(tr.AccessToken)

	return nil
}
