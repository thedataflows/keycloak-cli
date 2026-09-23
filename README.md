# keycloak-cli

A Go CLI for declarative Keycloak administration. It speaks the Keycloak Admin REST API through an OpenAPI contract, so every operation is validated against the spec before it hits the wire. Endpoints are not hard-coded; they are deduced from the supplied OpenAPI spec, so passing a different Keycloak spec lets the CLI use the appropriate endpoints for any resource the spec describes.

The tool loads manifest files (JSON or YAML), validates them against the Keycloak OpenAPI spec (passed with the required `--spec-path` flag; a bundled copy lives at `keycloak-oapi/${KEYCLOAK_VERSION}.spec.json`), and applies them via the Admin REST API. It can also fetch live state, compare it with local manifests, generate test fixtures, and manage admin tokens.

## How it works

1. **OpenAPI contract** — At startup the CLI loads the Keycloak Admin REST OpenAPI spec (`--spec-path`). Every resource type, operation path, request schema, and path parameter is derived from the spec, so the tool adapts to whichever Keycloak OpenAPI document you provide.
2. **Resource mapping** — The CLI discovers CRUD operations per resource type (`realm`, `user`, `client`, `group`, `role`, etc.) by scanning the spec. It builds an internal routing table so `POST /admin/realms/{realm}/users` is used for creating users, `PUT /admin/realms/{realm}/users/{user-id}` for updates, and so on.
3. **Validation** — Before any network call, resources and relationships are validated against the spec schemas and operation contracts.
4. **Dependency ordering** — Resources are sorted by dependency (realms first, then users, groups, roles, clients, etc.) so uploads succeed without manual ordering.
5. **Relationship extraction** — The CLI can fetch and export relationships (user-group memberships, role composites, client scopes, etc.) by walking the live API state.
6. **Round-trip comparison** — The `compare` command canonicalizes both local manifests and fetched state (stripping server-managed fields) and performs subset-based comparison.

## Features

- **Spec-aware validation** — Resources and relationships are validated against the Keycloak OpenAPI spec before apply.
- **Dynamic resource discovery** — New resource types available in the spec are automatically supported for fetch and apply without code changes.
- **Relationship validation** — Relationship operations are matched against the admin API contract using path template inference.
- **Fetch with relationship export** — Pull live state including user-group memberships, role composites, client scope mappings, federated identities, and organization links.
- **Canonicalized fetch** — Strip server-managed fields to produce a clean manifest suitable for re-apply.
- **Manifest generation** — Generate test realms with configurable counts of users, clients, roles, groups, organizations, identity providers, client scopes, authentication flows, password policies, and security defenses.
- **Round-trip comparison** — Compare local manifests with live Keycloak state through normalized resource and relationship state.
- **Dry-run uploads** — Validate and preview changes without sending them.
- **Continue-on-error** — Apply as many resources as possible even when individual items fail.
- **Admin token retrieval** — Password-grant flow with automatic `.env` storage.
- **Integration stack** — Docker Compose stack for local Keycloak testing via `mise`.

## Build

```bash
go build .
```

This produces a `keycloak-cli` binary in the current directory. Once built (or installed on your `PATH`), use it directly:

```bash
keycloak-cli --help
```

## Commands

### Global flags

Every command accepts the same global flags:

| Flag                  | Short | Default                          | Description                               |
| --------------------- | ----- | -------------------------------- | ----------------------------------------- |
| `--keycloak-base-url` | `-u`  | `http://localhost:8080`          | Keycloak base URL                         |
| `--timeout`           | `-t`  | `5s`                             | Request timeout                           |
| `--spec-path`         |       | *(required)*                     | Path to the OpenAPI spec                  |
| `--log-level`         |       | `info`                           | `trace`, `debug`, `info`, `warn`, `error` |
| `--log-format`        |       | `console`                        | `console` or `json`                       |

### `fetch`

Fetch resources from Keycloak.

```bash
# Fetch default resources (realm, user, client, group, role)
keycloak-cli fetch

# Fetch specific resource types
keycloak-cli fetch user,group,role --realm demo

# Fetch with relationships
keycloak-cli fetch user --realm demo --relationships -f json

# Canonicalize output for re-apply
keycloak-cli fetch --canonicalize -o demo-manifest.json --realm demo

# Search and limit results
keycloak-cli fetch user --realm demo --search alice --max 10

# Filter fetched resources by exact name or id (case-insensitive)
keycloak-cli fetch client --realm demo myapp

# Output each resource to a separate file in a directory (requires structured format)
keycloak-cli fetch user --realm demo -f yaml -o ./out/

# Fetch nested child resources (e.g. authentication executions)
keycloak-cli fetch authenticationexecution --realm demo --parent "browser"

# Fetch up to N levels of child resources
keycloak-cli fetch client --realm demo --depth 2
```

Flags:
| Flag               | Short | Default                        | Description                                                       |
| ------------------ | ----- | ------------------------------ | ----------------------------------------------------------------- |
| `--realm`          | `-r`  |                                | Scope fetch to one realm                                          |
| `--resources`      | (arg) | `realm,user,client,group,role` | Comma-separated resource types                                    |
| `--filter`         | (arg) |                                | Filter fetched resources by exact name or id (case-insensitive)   |
| `--search`         | `-s`  |                                | Search filter for collection endpoints                            |
| `--max`            |       | `0`                            | Maximum results to return                                         |
| `--depth`          |       | `1`                            | Fetch child resources up to N levels deep                         |
| `--parent`         | `-p`  |                                | Parent resource identifier for nested resources                   |
| `--relationships`  |       | `false`                        | Fetch supported relationship state                                |
| `--canonicalize`   |       | `false`                        | Strip server-managed fields for re-apply                          |
| `--format`         | `-f`  | `table`                        | `table`, `json`, `yaml`, `toml`                                   |
| `--output`         | `-o`  |                                | Write output to a file or directory (ends with `/` for directory) |
| `--force`          |       | `false`                        | Overwrite existing output file                                    |
| `--exclude-fields` | `-e`  | `containerId`                  | Comma-separated fields to exclude                                 |
| `--long`           | `-l`  | `false`                        | Show detailed information                                         |

### `upload`

Load manifest files and apply them to Keycloak.

```bash
# Upload a directory of manifests
keycloak-cli upload generated/

# Dry run
keycloak-cli upload generated/ --dry-run

# Continue on error
keycloak-cli upload generated/ --continue-on-error

# Force delete behavior
keycloak-cli upload generated/ --delete
```

Flags:

| Flag                  | Default | Description                                  |
| --------------------- | ------- | -------------------------------------------- |
| `--dry-run`           | `false` | Validate without sending changes             |
| `--continue-on-error` | `false` | Keep applying after individual failures      |
| `--delete`            | `false` | Force delete behavior for uploaded resources |
| `--format`            |         | `table`                                      | `table`, `json`, `yaml`                                            |
| `--output`            |         |                                              | Write results to a file or directory (ends with `/` for directory) |
| `--force`             | `false` | Overwrite existing output file               |

### `generate`

Generate test realm manifests.

```bash
# Generate a basic realm
keycloak-cli generate --realm demo

# Generate with specific counts
keycloak-cli generate --realm demo --with-users 5 --with-clients 2 --with-roles 3

# Overwrite existing files
keycloak-cli generate --realm demo -o generated/ --overwrite
```

Flags:

| Flag                          | Short | Default      | Description                                            |
| ----------------------------- | ----- | ------------ | ------------------------------------------------------ |
| `--realm`                     | `-r`  | `test-realm` | Target realm name                                      |
| `--output`                    | `-o`  |              | Output file or directory (ends with `/` for directory) |
| `--format`                    | `-f`  | `json`       | `json`, `yaml`, `toml`                                 |
| `--overwrite`                 |       | `false`      | Replace existing files                                 |
| `--summary`                   |       | `false`      | Print generation summary                               |
| `--with-users`                |       | `1`          | Number of users                                        |
| `--with-clients`              |       | `1`          | Number of clients                                      |
| `--with-roles`                |       | `1`          | Number of roles                                        |
| `--with-groups`               |       | `1`          | Number of groups                                       |
| `--with-organizations`        |       | `1`          | Number of organizations                                |
| `--with-identity-providers`   |       | `1`          | Number of identity providers                           |
| `--with-client-scopes`        |       | `1`          | Number of client scopes                                |
| `--with-authentication-flows` |       | `1`          | Number of authentication flows                         |
| `--with-password-policies`    |       | `1`          | Number of password policies                            |
| `--with-security-defenses`    |       | `1`          | Generate security defenses                             |

Output: `realm.<format>` and, when needed, `relationships.<format>` in the output directory. When output is a single file, only the realm resources are written; stdout includes both resources and relationships when relationships are generated.

### `compare`

Compare local manifests with fetched Keycloak state. Exits with an error if states differ.

```bash
# Compare a directory against a realm
keycloak-cli compare generated/ --realm demo

# JSON report
keycloak-cli compare generated/ --realm demo -f json

# Compare relationships only
keycloak-cli compare generated/relationships.json --realm demo
```

Flags:

| Flag       | Short | Default | Description                                                       |
| ---------- | ----- | ------- | ----------------------------------------------------------------- |
| `--realm`  | `-r`  |         | Target realm (inferred from manifests if omitted)                 |
| `--format` | `-f`  | `table` | `table`, `json`, `yaml`                                           |
| `--output` | `-o`  |         | Write report to a file or directory (ends with `/` for directory) |
| `--force`  |       | `false` | Overwrite existing output file                                    |

### `invoke`

Perform one generic, spec-driven API call against the loaded OpenAPI spec, or list the spec's operations.

```bash
# List operations whose id, path or summary contain a substring (case-insensitive)
keycloak-cli invoke --list "client role"

# Get one user. Resource+verb is the primary resolution mode: the bundled
# spec defines no operationIds. The resource is the plural path segment.
keycloak-cli invoke --resource users --verb GET --realm demo --param user-id=<uuid>

# Params that are not path placeholders become query parameters
keycloak-cli invoke --resource users --verb GET --realm demo --param max=5 --param username=alice

# Create a user with a JSON body (inline or @file)
keycloak-cli invoke --resource users --verb POST --realm demo --body @new-user.json

# Specs that do define operationIds can address operations directly
keycloak-cli invoke getUser --realm demo --param user-id=<uuid>
```

Flags:

| Flag         | Short | Default | Description                                                                                  |
| ------------ | ----- | ------- | -------------------------------------------------------------------------------------------- |
| (op-id)      | (arg) |         | OpenAPI operationId to invoke (optional; only usable with specs that define operationIds)    |
| `--resource` |       |         | Resource path segment to resolve (e.g. `users`), combined with `--verb`                      |
| `--verb`     |       | `GET`   | HTTP verb when resolving via `--resource`: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`           |
| `--realm`    | `-r`  |         | Fills the `{realm}` path placeholder when the operation's template has one                   |
| `--param`    |       |         | Path or query parameter as `k=v` (repeatable); keys matching path placeholders are substituted |
| `--body`     |       |         | Request body: inline JSON or `@file`                                                        |
| `--list`     |       |         | List the spec's operations whose id, path or summary contain this substring, print JSON and exit |

### `graph`

Inspect the resource graph the loaded spec implies: the relationship edges its path structure contains, single objects resolved by type and name, and the neighbours hanging off a resolved object. Resource types use the plural path-segment vocabulary of API paths (e.g. `users`), the same vocabulary `invoke --resource` uses.

```bash
# List every relationship edge implied by the spec (offline; no server contact)
keycloak-cli graph edges

# Filter edges by parent and child resource type
keycloak-cli graph edges --parent realms --child roles

# Resolve a user by username
keycloak-cli graph resolve users alice --realm demo

# Resolve straight by id (the <name> argument is still required but ignored)
keycloak-cli graph resolve users "" --realm demo --id <uuid>

# Resolve a realm by its own name (the realm is the graph root)
keycloak-cli graph resolve realms demo --realm demo

# Walk the realm root: list its child collections (users, clients, roles, ...)
keycloak-cli graph neighbors realms demo --realm demo

# List the objects related to a resolved user, with the edges that produced them
keycloak-cli graph neighbors users alice --realm demo

# Only walk edges whose child is a given resource type
keycloak-cli graph neighbors users alice --realm demo --child groups
```

Subcommands:

| Subcommand  | Arguments       | Description                                                                                   |
| ----------- | --------------- | --------------------------------------------------------------------------------------------- |
| `edges`     |                 | List relationship edges implied by the spec, as JSON (`--parent`/`--child` filters)           |
| `resolve`   | `<type> <name>` | Resolve one object by plural type plus name (`--realm` required; `--id` skips the name lookup) |
| `neighbors` | `<type> <name>` | Fetch each reachable child collection once; prints `{"nodes": [...], "edges": [...]}`         |

### `admin-token`

Get an admin access token through the password grant flow.

```bash
# Default: admin/admin on master realm, store to .env
keycloak-cli admin-token

# Custom credentials
keycloak-cli admin-token --username admin --password admin --realm master

# Print token only, do not write .env
keycloak-cli admin-token --set-env=false
```

Flags:

| Flag                  | Default                  | Description                    |
| --------------------- | ------------------------ | ------------------------------ |
| `--username`          | `admin`                  | Admin username                 |
| `--password`          | `admin`                  | Admin password                 |
| `--realm`             | `master`                 | Keycloak realm                 |
| `--set-env`           | `true`                   | Write tokens to `.env`         |
| `--access-token-env`  | `KEYCLOAK_ACCESS_TOKEN`  | Env var name for access token  |
| `--refresh-token-env` | `KEYCLOAK_REFRESH_TOKEN` | Env var name for refresh token |

### `version`

Print the version.

```bash
keycloak-cli version
```

### `mcp`

Serve the same capabilities to LLM agents as an MCP (Model Context Protocol) server — nine tools: `kc_operations`, `kc_invoke`, `kc_list`, `kc_resolve`, `kc_neighbors`, `kc_edges`, `kc_fetch`, `kc_apply`, `kc_reload` — over stdio (for locally-launched agents) or streamable HTTP (for remote ones). The HTTP transport implements the [2026-07-28](https://blog.modelcontextprotocol.io/posts/2026-07-28/) stateless core: no initialize handshake, no sessions — point a harness at the URL and POST. `kc_list` is GET-only by construction; state-changing work (`kc_invoke` non-GET verbs, every `kc_apply`) always requires an explicit `confirm: true` tool argument. On startup the server prints a quick guide to stderr: the tool list, the safety rule, and wiring snippets for Claude Code and generic `mcp.json` harnesses. See [docs/design/mcp-server.md](docs/design/mcp-server.md) for the full design.

```bash
# stdio (default; stdout carries protocol frames, the startup guide and logging go to stderr)
keycloak-cli --spec-path keycloak-oapi/${KEYCLOAK_VERSION}.spec.json mcp

# streamable HTTP on a loopback address (no authentication — keep it local)
keycloak-cli --spec-path keycloak-oapi/${KEYCLOAK_VERSION}.spec.json mcp --transport http --http-addr 127.0.0.1:8081
```

Wire it into an agent harness (also printed at startup):

```bash
# Claude Code
claude mcp add keycloak -- keycloak-cli --spec-path keycloak-oapi/${KEYCLOAK_VERSION}.spec.json mcp
```

```json
{
  "mcpServers": {
    "keycloak": { "command": "keycloak-cli", "args": ["--spec-path", "keycloak-oapi/${KEYCLOAK_VERSION}.spec.json", "mcp"] }
  }
}
```

## Library usage (pkg/kcapi)

The CLI is a thin layer over [`pkg/kcapi`](pkg/kcapi/), a spec-driven Keycloak client library: every operation in the supplied OpenAPI spec is discoverable and invocable generically, plus graph queries over the object relationships the spec's path structure implies.

```go
client, err := kcapi.New(kcapi.Config{
    BaseURL: "https://kc.example.com",
    Spec:    kcapi.SpecSource{Path: "keycloak-oapi/${KEYCLOAK_VERSION}.spec.json"}, // Path, URL, or Raw bytes
})
```

Tokens resolve exactly as they do for the CLI: `KEYCLOAK_ACCESS_TOKEN` from the environment (refreshed via `KEYCLOAK_REFRESH_TOKEN`), so the same `.env` file works. `Config.Credentials` validates which grant shape you intend (password pair vs client secret), but the secret values themselves come from the environment; set `Config.Auth` to your own `auth.Service` to source tokens programmatically (kcapi only calls its `AccessToken` method).

```go
// Discovery: list the spec's operations matching a filter (zero fields match everything).
ops, err := client.Operations(kcapi.OpFilter{Resource: "users", Method: kcapi.Get})

// Invoke one operation. Exactly one resolution mode: Op, or Resource+Verb.
// The bundled Keycloak spec defines no operationIds, so Resource+Verb is the
// primary mode; the resource is the plural path segment (e.g. "users").
user, err := client.Invoke(ctx, kcapi.Call{
    Resource: "users",
    Verb:     kcapi.Get,
    Realm:    "demo",
    Params:   kcapi.P{"user-id": "59b0..."}, // path placeholders first; the rest go to the query string
})

// Page through a collection operation via first/max (kcapi.DefaultPageSize, 100, per page).
all, err := client.ListAll(ctx, kcapi.Call{Resource: "users", Verb: kcapi.Get, Realm: "demo"})

// Graph queries. Types speak the same plural path-segment vocabulary as Call.Resource.
node, err := client.Resolve(ctx, kcapi.Ref{Type: "users", Name: "alice", Realm: "demo"})
neighbors, walked, err := client.Neighbors(ctx, node, kcapi.EdgeFilter{Child: "groups"})
edges := client.Edges() // every parent→child edge the spec implies (230 edges / 35 pairs on the bundled spec)

// Hot-swap the spec from its configured source; in-flight calls keep the spec they started with.
err = client.Reload(ctx)
```

Failed calls return `*kcapi.Error{Kind, Op, Status, Body, Err}`; match failures with `errors.Is` against the `kcapi.ErrNotFound`, `ErrConflict`, `ErrValidation`, and `ErrAuth` sentinels. The declarative `fetch`/`upload`/`compare` flows live in [`pkg/manifest`](pkg/manifest/), built on the same client.

## Supported resources

Resource types are discovered dynamically from the OpenAPI spec. The CLI does not maintain a hard-coded list of endpoints; instead it scans the spec, creates a contract for every resource type with recognizable CRUD operations, and uses the appropriate endpoint for each request. Common types discovered from the bundled Keycloak ${KEYCLOAK_VERSION} spec include:

| Resource type             | Typical operations     |
| ------------------------- | ---------------------- |
| `realm`                   | GET, POST, PUT, DELETE |
| `user`                    | GET, POST, PUT, DELETE |
| `client`                  | GET, POST, PUT, DELETE |
| `group`                   | GET, POST, PUT, DELETE |
| `role`                    | GET, POST, PUT, DELETE |
| `clientscope`             | GET, POST, PUT, DELETE |
| `component`               | GET, POST, PUT, DELETE |
| `identityprovider`        | GET, POST, PUT, DELETE |
| `organization`            | GET, POST, PUT, DELETE |
| `authenticationflow`      | GET, POST, PUT         |
| `authenticationexecution` | GET                    |
| `protocolmapper`          | GET, POST, PUT         |
| `identityprovidermapper`  | GET, POST, PUT         |
| `requiredactionprovider`  | GET, PUT               |
| `clientpolicies`          | GET, PUT               |
| `clientprofiles`          | GET, PUT               |
| `localization`            | GET, POST, PUT, DELETE |
| `workflow`                | GET, POST, PUT, DELETE |

Additional types—such as `resource`, `scope`, `policy`, `credential`, `event`, and any new types introduced in a future spec—are supported automatically through dynamic operation mapping.

## Supported relationships

Relationship families are inferred from the admin API contract by scanning relationship-like paths in the OpenAPI spec; a built-in registry supplies friendly names and payload shapes for the most common patterns. Both fetch and apply are supported for the following families:

| Relationship                       | Description                          |
| ---------------------------------- | ------------------------------------ |
| `user-group-membership`            | User group memberships               |
| `user-realm-role-mapping`          | User realm role mappings             |
| `user-client-role-mapping`         | User client role mappings            |
| `user-federated-identity`          | User federated identity links        |
| `group-realm-role-mapping`         | Group realm role mappings            |
| `group-client-role-mapping`        | Group client role mappings           |
| `role-composite-mapping`           | Realm role composite memberships     |
| `client-role-composite`            | Client role composite memberships    |
| `client-scope-realm-role-mapping`  | Client scope to realm role mappings  |
| `client-scope-client-role-mapping` | Client scope to client role mappings |
| `client-realm-scope-mapping`       | Client realm scope mappings          |
| `client-client-scope-mapping`      | Client to client scope mappings      |
| `default-group-membership`         | Realm default groups                 |
| `realm-default-client-scope`       | Realm default client scopes          |
| `realm-optional-client-scope`      | Realm optional client scopes         |
| `client-default-scope`             | Client default scope mappings        |
| `client-optional-scope`            | Client optional scope mappings       |
| `organization-member`              | Organization member assignments      |
| `organization-identity-provider`   | Organization identity provider links |

New relationship patterns that appear in a different Keycloak spec can be picked up automatically or customized with a `relationship-overrides.yaml` file next to the spec.

## Manifest format

### Resource manifest

A resource manifest is a JSON or YAML array of resource objects.

```json
[
  {
    "type": "realm",
    "realm": "demo",
    "data": {
      "realm": "demo",
      "enabled": true,
      "displayName": "Demo Realm"
    }
  },
  {
    "type": "user",
    "realm": "demo",
    "data": {
      "username": "alice",
      "enabled": true,
      "email": "alice@example.com"
    }
  },
  {
    "type": "client",
    "realm": "demo",
    "data": {
      "clientId": "my-app",
      "name": "My Application",
      "enabled": true
    }
  }
]
```

Fields:

- `type` — Resource type (must match a type discoverable from the spec).
- `realm` — Target realm name.
- `data` — The resource payload. Must conform to the schema for the resource type and operation.
- `parentType` — Optional. Disambiguates the parent resource type for nested resources (e.g. a `protocolmapper` under a `clientscope` vs a `client`).
- `delete` — If `true`, the resource is deleted instead of created/updated.

### Relationship manifest

Relationship manifests use a `relationships` envelope.

```json
{
  "relationships": [
    {
      "path": "demo/users/alice/groups/admin-group"
    },
    {
      "path": "demo/users/alice/role-mappings/realm",
      "data": [
        { "name": "admin-role" }
      ]
    },
    {
      "path": "demo/roles-by-id/admin-role/composites",
      "data": [
        { "name": "composite-role" }
      ]
    }
  ]
}
```

Relationship paths are resolved against the actual spec templates discovered from the OpenAPI document. The example paths above match templates such as `{realm}/users/{user-id}/groups/{groupId}`, `{realm}/users/{user-id}/role-mappings/realm`, and `{realm}/roles-by-id/{role-id}/composites`. The CLI validates each path and payload against the spec operation contract.

### Loading manifests

Manifests can be loaded from files or directories. The CLI accepts both JSON and YAML formats.

```bash
# Single file
keycloak-cli upload realm.json

# Directory (loads all `.json` and `.yaml` files)
keycloak-cli upload generated/

# Multiple files
keycloak-cli upload realm.json relationships.json
```

## Environment setup

The CLI loads `.env`, `.local.env`, and `.development.env` if they exist.

Useful variables:

| Variable                      | Description                                  |
| ----------------------------- | -------------------------------------------- |
| `KEYCLOAK_BASE_URL`           | Keycloak server URL                          |
| `KEYCLOAK_ACCESS_TOKEN`       | Bearer token for authenticated requests      |
| `KEYCLOAK_REFRESH_TOKEN`      | Refresh token for token renewal              |
| `KEYCLOAK_USERNAME`           | Admin username for token command             |
| `KEYCLOAK_PASSWORD`           | Admin password for token command             |
| `KC_BOOTSTRAP_ADMIN_USERNAME` | Alternative admin username for token command |
| `KC_BOOTSTRAP_ADMIN_PASSWORD` | Alternative admin password for token command |
| `KEYCLOAK_REALM`              | Default realm for token command              |

Values in `.env` take precedence over inherited environment variables, so an `admin-token` run that updates `.env` is immediately picked up by the next command.

Example `.env`:

```bash
KEYCLOAK_BASE_URL=http://keycloak.172.21.0.3.nip.io:8080
KEYCLOAK_USERNAME=admin
KEYCLOAK_PASSWORD=admin
KEYCLOAK_REALM=master
```

## Local integration stack

A Docker Compose stack is available for local Keycloak testing. Tasks are managed with [mise-en-place](https://mise.jdx.dev/).

Start:

```bash
mise run integration:up
# or
bash mise-tasks/integration/up.sh
```

Test:

```bash
mise run integration:test
# or
bash mise-tasks/integration/test.sh
```

Stop:

```bash
mise run integration:down
# or
bash mise-tasks/integration/down.sh
```

The stack reads `KEYCLOAK_VERSION`, `KC_BOOTSTRAP_ADMIN_USERNAME`, and `KC_BOOTSTRAP_ADMIN_PASSWORD` from `.env`.

## Updating the OpenAPI spec

Download a new Keycloak OpenAPI spec:

```bash
KEYCLOAK_VERSION=26.7.4 mise run generate:oapi
```

Force overwrite an existing spec:

```bash
KEYCLOAK_VERSION=26.7.4 mise run generate:oapi -- --force
```

## Development

Run tests:

```bash
go test ./...
```

Run vet:

```bash
go vet ./...
```

Run linter:

```bash
golangci-lint run
```

Format code:

```bash
go fmt ./...
```

## Architecture

```properties
cmd/            — CLI commands (Kong)
pkg/kcapi/      — Spec-driven client library: discovery, invoke, graph, runtime, errors
pkg/manifest/   — Declarative fetch/apply/compare flows on top of kcapi
pkg/auth/       — Token acquisition and .env storage
pkg/output/     — Table/JSON/YAML/TOML formatting
pkg/realmgen/   — Test realm generation
keycloak-oapi/  — Bundled Keycloak OpenAPI specs
```

## Current scope

The CLI does not hard-code Keycloak endpoints. It deduces resource types, CRUD operations, and relationship patterns directly from the supplied OpenAPI spec. In principle, any Keycloak Admin REST OpenAPI document can be passed with `--spec-path`, and the CLI will dynamically discover the appropriate endpoints for the requested resources. The contract layer, validation path, relationship export, round-trip normalization, and compare flow are spec-driven and do not require code changes when Keycloak adds new resource types.
