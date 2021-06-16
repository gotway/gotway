#!/usr/bin/env bash

release="gotway"
repo="gotway/gotway"

git fetch --all
tag=$(git describe --abbrev=0 --tags)

helm repo add "$repo" https://charts.gotway.duckdns.org
helm repo update

echo "🚀 Deploying '$repo' with image version '$tag'..."
helm upgrade --install "$release" "$repo" \
  --set image.tag=$tag \
  --set catalog.image.tag=$tag \
  --set stock.image.tag=$tag \
  --set route.image.tag=$tag
