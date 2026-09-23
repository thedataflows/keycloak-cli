package cmd

import (
	"context"
	"net"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/mcpserver"
)

// McpCmd serves the kcapi library to LLM agents as an MCP server: JSON-RPC
// over stdio (default) or over the streamable HTTP transport. stdout carries
// protocol frames in stdio mode, logging stays on stderr. State-changing
// kc_invoke calls require confirm:true on both transports.
type McpCmd struct {
	Transport string `help:"MCP transport (stdio,http)" enum:"stdio,http" default:"stdio"`
	// HTTPAddr only applies to --transport=http. The HTTP surface has no
	// authentication, so the default stays on loopback.
	HTTPAddr string `help:"Listen address for --transport=http" default:"127.0.0.1:8081"`
}

func (c *McpCmd) Run(ctx *kong.Context, cli *CLI) error {
	client, err := cli.Kcapi()
	if err != nil {
		return err
	}
	if c.Transport == "http" {
		ln, err := net.Listen("tcp", c.HTTPAddr)
		if err != nil {
			return err
		}
		log.Logger.Info().Str("pkg", PKG_CMD).Str("addr", ln.Addr().String()).Msg("Serving MCP over streamable HTTP")
		return mcpserver.RunHTTP(context.Background(), client, ln)
	}
	log.Logger.Info().Str("pkg", PKG_CMD).Msg("Serving MCP over stdio")
	return mcpserver.Run(context.Background(), client)
}
