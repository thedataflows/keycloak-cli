#!/bin/env bash
#MISE description="Test local development environment integration"
#MISE alias="it"

set -euo pipefail

set -x

## Get Keycloak Admin access token, willbe written to .env as KEYCLOAK_ACCESS_TOKEN (default)
go run . --keycloak-base-url="$KEYCLOAK_BASE_URL" admin-token

## Reload the env to get the new token
mise env >/dev/null

xh \
  --verbose --body --json \
  --bearer "$KEYCLOAK_ACCESS_TOKEN" \
  "$KEYCLOAK_BASE_URL/admin/realms" \
  "@${MISE_TASK_DIR:-./mise-tasks/integration}/test-realm.json"
