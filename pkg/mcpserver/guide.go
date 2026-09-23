package mcpserver

import "fmt"

// Guide returns the operator quick start printed to stderr at startup: what
// the server exposes, the safety rule, and how to wire it into the common
// agent harnesses — so a fresh `keycloak-cli mcp` needs no docs to get
// running. transport is "stdio" or "http"; addr is the listen address in http
// mode (empty otherwise).
func Guide(transport, addr string) string {
	httpLine := "  Streamable HTTP:  keycloak-cli --spec-path /path/to/keycloak.spec.json mcp --transport http --http-addr 127.0.0.1:8081\n"
	if addr != "" {
		httpLine = fmt.Sprintf("  Streamable HTTP:  keycloak-cli --spec-path /path/to/keycloak.spec.json mcp --transport http --http-addr %s\n", addr)
	}
	endpointLine := "    then point the harness at \"http://127.0.0.1:8081/mcp\" (any path works; no authentication — keep it on loopback)\n"
	if addr != "" {
		endpointLine = fmt.Sprintf("    then point the harness at \"http://%s/mcp\" (any path works; no authentication — keep it on loopback)\n", addr)
	}
	return fmt.Sprintf(`keycloak-cli MCP server ready — transport: %s

Tools
  kc_operations  list spec operations (filters: resource, method, tag, search)
  kc_invoke      call one operation (resource+verb or op); non-GET needs "confirm": true
  kc_list        fetch ALL pages of a GET collection in one call
  kc_resolve     resolve one resource by type + name/id
  kc_neighbors   list the child collections hanging off a resource
  kc_edges       show the parent→child relationship vocabulary
  kc_fetch       realm-scoped export {resources, relationships, failures}
  kc_apply       apply those resources (create/update/delete); needs "confirm": true
  kc_reload      reload the OpenAPI spec from disk

Wire it into an agent harness
  Claude Code:  claude mcp add keycloak -- keycloak-cli --spec-path /path/to/keycloak.spec.json mcp
  mcp.json (Cursor, Claude Desktop, Windsurf, ...):
    {
      "mcpServers": {
        "keycloak": { "command": "keycloak-cli", "args": ["--spec-path", "/path/to/keycloak.spec.json", "mcp"] }
      }
    }
%s%s`, transport, httpLine, endpointLine)
}
