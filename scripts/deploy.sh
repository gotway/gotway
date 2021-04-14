#!/usr/bin/env bash

RELEASE="gotway"
REPO="gotway"

helm repo add "$REPO" https://charts.gotway.duckdns.org
helm repo update

echo "🚀 Deploying '${RELEASE}'..."
helm upgrade --install "$RELEASE" "$REPO/$RELEASE"
