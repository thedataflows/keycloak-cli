package cmd

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
	"github.com/thedataflows/keycloak-cli/pkg/output"
)

// GraphCmd inspects the resource graph the loaded spec implies: the
// relationship edges its path structure contains, a single object resolved by
// type and name, and the neighbours hanging off a resolved object.
type GraphCmd struct {
	Edges     GraphEdgesCmd     `cmd:"" help:"List relationship edges implied by the spec"`
	Resolve   GraphResolveCmd   `cmd:"" help:"Resolve a Keycloak object by type + name"`
	Neighbors GraphNeighborsCmd `cmd:"" help:"List objects related to a resolved object"`
}

// GraphEdgesCmd lists the spec's structural edges. Edges are inferred from the
// path templates alone, so this subcommand works without a live server.
type GraphEdgesCmd struct {
	Parent string `help:"Only edges with this parent resource type (as in API paths, e.g. realms)"`
	Child  string `help:"Only edges with this child resource type (as in API paths, e.g. roles)"`
}

func (c *GraphEdgesCmd) Run(_ *kong.Context, cli *CLI) error {
	log.Logger.Info().Str("pkg", PKG_CMD).Msg("List spec graph edges")

	client, err := cli.Kcapi()
	if err != nil {
		return err
	}
	// Edges are inferred from the path templates alone, so no server is
	// contacted; kcapi.Edges has no filter parameter, so the flags apply here.
	edges := filterGraphEdges(client.Edges(), graphFilter(c.Parent, c.Child))
	if edges == nil {
		edges = []kcapi.Edge{}
	}
	return output.WriteJSON(os.Stdout, edges)
}

// GraphResolveCmd resolves one object by plural resource type plus human name,
// or straight by id with --id.
type GraphResolveCmd struct {
	Type  string `arg:"" help:"Resource type as in API paths, e.g. users"`
	Name  string `arg:"" help:"Object name (username, client id, role name)"`
	ID    string `help:"Skip name lookup, use this id"`
	Realm string `required:"" help:"Realm to resolve in"`
}

func (c *GraphResolveCmd) Run(_ *kong.Context, cli *CLI) error {
	log.Logger.Info().Str("pkg", PKG_CMD).Msgf("Resolve %s %q in realm %q", c.Type, c.Name, c.Realm)

	client, err := cli.Kcapi()
	if err != nil {
		return err
	}

	commandCtx, cancel := cli.CreateContextWithTimeout()
	defer cancel()

	node, err := client.Resolve(commandCtx, graphRef(c.Type, c.Name, c.Realm, c.ID))
	if err != nil {
		return err
	}
	return output.WriteJSON(os.Stdout, node)
}

// GraphNeighborsCmd resolves the named object and walks every edge whose
// parent is its type, fetching each reachable child collection once.
type GraphNeighborsCmd struct {
	Type   string `arg:"" help:"Resource type as in API paths, e.g. users"`
	Name   string `arg:"" help:"Object name (username, client id, role name)"`
	Realm  string `required:"" help:"Realm to resolve in"`
	Child  string `help:"Only edges with this child resource type (as in API paths, e.g. roles)"`
	Parent string `help:"Only edges with this parent resource type (as in API paths, e.g. realms)"`
}

func (c *GraphNeighborsCmd) Run(_ *kong.Context, cli *CLI) error {
	log.Logger.Info().Str("pkg", PKG_CMD).Msgf("List neighbors of %s %q in realm %q", c.Type, c.Name, c.Realm)

	client, err := cli.Kcapi()
	if err != nil {
		return err
	}

	commandCtx, cancel := cli.CreateContextWithTimeout()
	defer cancel()

	node, err := client.Resolve(commandCtx, graphRef(c.Type, c.Name, c.Realm, ""))
	if err != nil {
		return err
	}
	nodes, edges, err := client.Neighbors(commandCtx, node, graphFilter(c.Parent, c.Child))
	if err != nil {
		return err
	}
	if nodes == nil {
		nodes = []kcapi.Node{}
	}
	if edges == nil {
		edges = []kcapi.Edge{}
	}
	out := struct {
		Nodes []kcapi.Node `json:"nodes"`
		Edges []kcapi.Edge `json:"edges"`
	}{Nodes: nodes, Edges: edges}
	return output.WriteJSON(os.Stdout, out)
}

// graphFilter maps the --parent/--child flags onto an EdgeFilter; empty flags
// stay empty, which the walk treats as match-everything.
func graphFilter(parent, child string) kcapi.EdgeFilter {
	return kcapi.EdgeFilter{Parent: parent, Child: child}
}

// filterGraphEdges keeps only the edges the filter admits; an empty filter
// field matches everything, so an empty filter returns the slice unchanged.
func filterGraphEdges(edges []kcapi.Edge, filter kcapi.EdgeFilter) []kcapi.Edge {
	if filter.Parent == "" && filter.Child == "" {
		return edges
	}
	filtered := make([]kcapi.Edge, 0, len(edges))
	for _, edge := range edges {
		if filter.Parent != "" && edge.Parent != filter.Parent {
			continue
		}
		if filter.Child != "" && edge.Child != filter.Child {
			continue
		}
		filtered = append(filtered, edge)
	}
	return filtered
}

// graphRef maps a resolve target onto a Ref; a non-empty id short-circuits the
// name lookup.
func graphRef(typ, name, realm, id string) kcapi.Ref {
	return kcapi.Ref{Type: typ, Name: name, Realm: realm, ID: id}
}
