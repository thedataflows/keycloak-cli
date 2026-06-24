package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"time"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
)

const (
	PKG_CMD  = "cmd"
	APP_NAME = "keycloak-cli"
)

type Globals struct {
	LogLevel        string        `help:"Log level (trace,debug,info,warn,error)" default:"info"`
	LogFormat       string        `help:"Log format (console,json)" default:"console"`
	KeycloakBaseURL string        `short:"u" help:"Keycloak base URL" default:"http://localhost:8080"`
	Timeout         time.Duration `short:"t" help:"Request timeout duration" default:"5s"`
	SpecPath        string        `help:"Path to the Keycloak OpenAPI specification file" default:"keycloak-oapi/26.6.2.spec.json"`
	// StateDir  string `help:"Directory for state storage. Not yet used." default:".state/"`
}

// CLI represents the main CLI structure
type CLI struct {
	Globals    `kong:"embed"`
	Version    VersionCmd    `cmd:"" help:"Show version information"`
	Fetch      FetchCmd      `cmd:"" help:"Fetch object"`
	Compare    CompareCmd    `cmd:"" help:"Compare local manifests with fetched realm state"`
	Upload     UploadCmd     `cmd:"" help:"Upload objects"`
	Generate   GenerateCmd   `cmd:"" help:"Generate test data"`
	AdminToken AdminTokenCmd `cmd:"" help:"Get administrative access token from current instance"`
}

// AfterApply is called after Kong parses the CLI but before the command runs
func (cli *CLI) AfterApply(ctx *kong.Context) error {
	// Skip initialization for version command or help flag.
	if ctx.Command() == "version" || hasHelpFlag(ctx.Args) {
		return nil
	}

	// Configure log level.
	if err := setGlobalLoggerLogLevel(cli.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	// Configure log format.
	if err := setGlobalLoggerLogFormat(cli.LogFormat, nil); err != nil {
		return fmt.Errorf("set log format: %w", err)
	}

	return nil
}

// CreateContextWithTimeout creates a context with the configured timeout
func (g *Globals) CreateContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), g.effectiveTimeout())
}

func hasHelpFlag(args []string) bool {
	return slices.Contains(args, "--help") || slices.Contains(args, "-h")
}

func (g *Globals) effectiveTimeout() time.Duration {
	if g.Timeout <= 0 {
		return 5 * time.Second
	}
	if g.Timeout < time.Second {
		return g.Timeout * time.Second
	}
	return g.Timeout
}

func (cli *CLI) adminClient() (admin.Service, error) {
	return admin.New(admin.Config{
		BaseURL:  cli.KeycloakBaseURL,
		SpecPath: cli.SpecPath,
		Timeout:  cli.effectiveTimeout(),
	})
}

// Run executes the CLI with the given version
func Run(version string, args []string) error {
	// Load .env file; use Overload so file values take precedence over inherited
	// environment variables (e.g. after admin-token updates .env).
	_ = godotenv.Overload(
		".env",             // Current directory
		".local.env",       // Local overrides (common in web development)
		".development.env", // Development environment
	)

	var cli CLI

	parser, err := kong.New(&cli,
		kong.Name(APP_NAME),
		kong.Description("A CLI client for Keycloak"),
		kong.Configuration(kongyaml.Loader),
		kong.UsageOnError(),
		kong.DefaultEnvars(""),
	)
	if err != nil {
		return fmt.Errorf("create CLI parser: %w", err)
	}

	ctx, err := parser.Parse(args)
	if hasHelpFlag(args) {
		return nil
	}
	if err != nil {
		parser.FatalIfErrorf(err)
		return err
	}

	// Check if this is the version command - handle it specially without logging/config
	if ctx.Command() == "version" {
		return ctx.Run(version)
	}

	log.Logger.Info().Str("pkg", PKG_CMD).Str("app", ctx.Model.Name).Str("version", version).Msg("Starting application")

	return ctx.Run(ctx, &cli)
}

// setGlobalLoggerLogLevel configures the zerolog global level.
func setGlobalLoggerLogLevel(levelStr string) error {
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		return fmt.Errorf("parse log level %q: %w", levelStr, err)
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = log.Logger.Level(level)
	return nil
}

// setGlobalLoggerLogFormat configures the zerolog output format.
// If w is nil, output falls back to os.Stderr.
func setGlobalLoggerLogFormat(format string, w io.Writer) error {
	if w == nil {
		w = os.Stderr
	}
	switch format {
	case "console":
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{
			Out:        w,
			TimeFormat: time.RFC3339,
		})
	case "json":
		log.Logger = log.Logger.Output(w)
	default:
		return fmt.Errorf("invalid log format %q", format)
	}
	return nil
}
