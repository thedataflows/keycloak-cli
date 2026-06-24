#!/bin/env bash
#MISE description="Convert Keycloak documentation to markdown format for use in the docs directory"
#MISE alias="kd"

set -euo pipefail

set -x
curl --no-progress-meter https://www.keycloak.org/docs/latest/server_admin/index.html | html2markdown --output docs/keycloak-admin.md $@
