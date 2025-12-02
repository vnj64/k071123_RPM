#!/usr/bin/env bash
set -euo pipefail

if [ ! -f /usr/share/kibana/config/kibana.token ]; then
  echo "ERROR: /usr/share/kibana/config/kibana.token not found!"
  exit 1
fi

export ELASTICSEARCH_SERVICEACCOUNTTOKEN
ELASTICSEARCH_SERVICEACCOUNTTOKEN=$(cat /usr/share/kibana/config/kibana.token)
echo "Using service account token: ${ELASTICSEARCH_SERVICEACCOUNTTOKEN:0:6}…"

exec /usr/local/bin/kibana-docker
