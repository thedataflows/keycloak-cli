package cmd

import (
	"context"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/mcpserver"
)

// McpCmd serves the kcapi library to LLM agents as an MCP server speaking
// JSON-RPC over stdio: stdout carries protocol frames, logging stays on
// stderr. State-changing kc_invoke calls require confirm:true.
type McpCmd struct{}

func (c *McpCmd) Run(ctx *kong.Context, cli *CLI) error {
	log.Logger.Info().Str("pkg", PKG_CMD).Msg("Serving MCP over stdio")
	client, err := cli.Kcapi()
	if err != nil {
		return err
	}
	return mcpserver.Run(context.Background(), client)
}
