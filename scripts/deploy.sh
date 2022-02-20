#!/usr/bin/env bash

set -e

repo="gotway"
release="gotway"
chart="gotway/gotway"
namespace="gotway"

git fetch --all
tag=$(git describe --tags $(git rev-list --tags --max-count=1))

helm repo add "$repo" https://charts.gotway.duckdns.org
helm repo update

echo "🚀 Deploying '$chart' with image version '$tag'..."
helm upgrade --install "$release" "$chart" --set image.tag=$tag --namespace "$namespace"