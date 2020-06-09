#!/bin/bash

TUFIN_DOMAIN="$1"
TUFIN_PROJECT="$2"
TUFIN_AGENT_TOKEN="$3"
TUFIN_NAMESPACE="${4:-default}"

main() {
    crt="$(echo | openssl s_client -connect "guru.tufin.io:443" 2> /dev/null | awk '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/')"

    if [[ "$?" != "0" ]]; then
        echo "Error fetching Orca certificates using openssl - please make sure you're connected to the internet & that 'openssl' is installed."
        exit 1
    fi

    kubectl create secret generic orca-secrets \
    --from-literal=docker-repo-username="${TUFIN_DOMAIN}_${TUFIN_PROJECT}" \
    --from-literal=guru-api-key="$TUFIN_AGENT_TOKEN" \
    --from-literal=guru-crt="$crt" \
    --namespace="$TUFIN_NAMESPACE" \
    --dry-run -o yaml
}

main
