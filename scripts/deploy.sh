#!/usr/bin/env bash

set -e

repo="gotway"
release="gotway"
chart="gotway/gotway"
namespace="gotway"

helm repo add "$repo" https://charts.gotway.duckdns.org
helm repo update

echo "🚀 Deploying '$chart' with image version '$tag'..."
helm upgrade --install "$release" "$chart" --namespace "$namespace"