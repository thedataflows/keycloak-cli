#!/bin/env bash
#MISE description="Start local development environment for integration tests"
#MISE alias="iu"

set -euo pipefail

set -x
docker compose -f "${MISE_TASK_DIR:-./mise-tasks/integration}/docker-compose.yaml" up -d
docker compose -f "${MISE_TASK_DIR:-./mise-tasks/integration}/docker-compose.yaml" ps
{ set +x; } 2>/dev/null

KC_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' integration-keycloak-1)
KEYCLOAK_BASE_URL="http://keycloak.${KC_IP}.nip.io:8080"
echo "Running on: $KEYCLOAK_BASE_URL"
echo "Management: ${KEYCLOAK_BASE_URL%:*}:9000"

# Update .env file with KEYCLOAK_BASE_URL
sed -i -E "s,^\s*(KEYCLOAK_BASE_URL)\s*=.*,\1=$KEYCLOAK_BASE_URL," .env || echo "KEYCLOAK_BASE_URL=$KEYCLOAK_BASE_URL" >> .env
