#!/bin/env bash
#MISE description="Download Keycloak OpenAPI spec"
#MISE alias="go"

set -euo pipefail

FORCE=0
OUTPUT_DIR=""
for arg in "$@"; do
  case $arg in
    --force)
      FORCE=1
      ;;
    *)
      if [[ -z "$KEYCLOAK_VERSION" ]]; then
        KEYCLOAK_VERSION="$arg"
      elif [[ -z "$OUTPUT_DIR" ]]; then
        OUTPUT_DIR="$arg"
      else
        echo "Too many arguments"
        exit 1
      fi
      ;;
  esac
done

if [[ -z "$KEYCLOAK_VERSION" ]]; then
    echo "Usage: $0 <keycloak-version> [output-dir] [--force]"
    exit 1
fi

OUTPUT_DIR="${OUTPUT_DIR:-./keycloak-oapi/}"

[[ -d "$OUTPUT_DIR" ]] || mkdir -p "$OUTPUT_DIR"

if [[ "$FORCE" -eq 0 && -f "${OUTPUT_DIR}/${KEYCLOAK_VERSION}.spec.json" ]]; then
    echo "[WARN] ${OUTPUT_DIR}/${KEYCLOAK_VERSION}.spec.json already exists. Use --force to overwrite."
else
  set -x
  curl "https://www.keycloak.org/docs-api/${KEYCLOAK_VERSION}/rest-api/openapi.json" | \
    jq --slurpfile patch "${MISE_TASK_DIR:-./mise-tasks/generate}/oapi-patch.json" -f "${MISE_TASK_DIR:-./mise-tasks/generate}/oapi-patch.jq" > "$OUTPUT_DIR/${KEYCLOAK_VERSION}.spec.json"
  { set +x; } 2>/dev/null
fi
