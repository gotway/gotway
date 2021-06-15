#!/usr/bin/env bash

release="gotway"
repo="gotway"
tag=$(git describe --abbrev=0 --tags)

helm repo add "$repo" https://charts.gotway.duckdns.org
helm repo update

echo "🚀 Deploying '${release}'..."
helm upgrade --install "$release" "$repo/$release" --set image.tag=$tag
