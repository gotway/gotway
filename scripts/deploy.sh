#!/usr/bin/env bash

set -e

function deploy() {
    name="$1"
    path="$2"
    manifests="$path/manifests/"

    if [ ! -d "$manifests" ]; then
        echo "❌    Manifests not found: '$manifests'"
        exit 1
    fi
    echo "🚀    Deploying '$name'. Context: '$path'"
    kubectl apply -f "$manifests"
}

deploy "gotway" .

for ms in $(ls -d microservices/*); do
    name=$(basename "$ms")
    path="$ms"
    deploy "$name" "$path"
done
