package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
	"github.com/thedataflows/keycloak-cli/pkg/output"
)

// InvokeCmd performs one generic, spec-driven API call against the loaded
// Keycloak OpenAPI specification, or lists the spec's operations.
type InvokeCmd struct {
	OpID     string            `arg:"" optional:"" help:"OpenAPI operationId to invoke. Only usable with specs that define operationIds; the bundled spec has none, so prefer --resource with --verb"`
	Realm    string            `short:"r" help:"Realm scope; fills the {realm} path placeholder when the operation's template has one"`
	Param    map[string]string `help:"Path or query parameters as k=v (repeatable); keys matching path placeholders are substituted, the rest become query parameters"`
	Body     string            `help:"Request body: inline JSON or @file"`
	Resource string            `help:"Resource path segment to resolve (e.g. users); primary resolution mode, combined with --verb instead of an operationId"`
	Verb     string            `help:"HTTP verb when resolving via --resource" default:"GET" enum:"GET,POST,PUT,DELETE,PATCH"`
	List     string            `help:"List the spec's operations whose id, path or summary contain this substring (case-insensitive), print them as JSON and exit"`
}

func (c *InvokeCmd) Run(ctx *kong.Context, cli *CLI) error {
	log.Logger.Info().Str("pkg", PKG_CMD).Msg("Invoke API operation")
	log.Logger.Debug().Str("pkg", PKG_CMD).Msgf("Invoke command options: %+v; context: %+v", cli, ctx.Args)

	client, err := cli.Kcapi()
	if err != nil {
		return err
	}

	if c.List != "" {
		ops, err := client.Operations(kcapi.OpFilter{Search: c.List})
		if err != nil {
			return err
		}
		if ops == nil {
			ops = []kcapi.Operation{}
		}
		return output.WriteJSON(os.Stdout, ops)
	}

	var body json.RawMessage
	if c.Body != "" {
		body, err = readBodyArg(c.Body)
		if err != nil {
			return err
		}
	}

	call, err := invokeCall(c.OpID, c.Realm, c.Param, body, c.Resource, c.Verb)
	if err != nil {
		return err
	}

	commandCtx, cancel := cli.CreateContextWithTimeout()
	defer cancel()

	out, err := client.Invoke(commandCtx, call)
	if err != nil {
		return err
	}
	return output.WriteJSON(os.Stdout, out)
}

// invokeCall maps the parsed command arguments onto a kcapi.Call. Exactly one
// resolution mode may be used: an operationId, or --resource together with
// --verb.
func invokeCall(opID, realm string, params map[string]string, body json.RawMessage, resource, verb string) (kcapi.Call, error) {
	// Body stays a nil interface when there is no body: a typed-nil
	// json.RawMessage inside interface{} would make kcapi treat the call as
	// carrying a body and reject operations without a request schema.
	call := kcapi.Call{Realm: realm, Params: kcapi.P(params)}
	if len(body) > 0 {
		call.Body = body
	}
	switch {
	case opID != "" && resource != "":
		return kcapi.Call{}, fmt.Errorf("use either an operationId or --resource with --verb, not both")
	case resource != "":
		if verb == "" {
			return kcapi.Call{}, fmt.Errorf("--resource %q requires a --verb", resource)
		}
		call.Resource = resource
		call.Verb = kcapi.Verb(verb)
	case opID != "":
		call.Op = opID
	default:
		return kcapi.Call{}, fmt.Errorf("provide an operationId or --resource (with --verb) to select an operation")
	}
	return call, nil
}

// readBodyArg resolves the --body argument: a leading "@" reads the named
// file, anything else is inline JSON. The result must be valid JSON either way.
func readBodyArg(arg string) (json.RawMessage, error) {
	raw := []byte(arg)
	if len(arg) > 0 && arg[0] == '@' {
		var err error
		raw, err = os.ReadFile(arg[1:])
		if err != nil {
			return nil, fmt.Errorf("read body file: %w", err)
		}
	}
	if !json.Valid(raw) {
		return nil, fmt.Errorf("body is not valid JSON")
	}
	return json.RawMessage(raw), nil
}
