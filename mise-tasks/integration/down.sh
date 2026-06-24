#!/bin/env bash
#MISE description="Stop local development environment for integration tests"
#MISE alias="id"

set -euo pipefail

set -x
docker compose -f "${MISE_TASK_DIR:-./mise-tasks/integration}/docker-compose.yaml" down
