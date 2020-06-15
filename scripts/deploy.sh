#!/usr/bin/env bash

set -e

env=$1
if [ -z "$env" ]; then
    echo "❌    Environment argument is mandatory"
    exit 1
fi

function deploy() {
    name="$1"
    env="$2"
    path="$3"
    manifests="$path/manifests/$env"

    if [ ! -d "$manifests" ]; then
        echo "❌    Manifests not found: '$manifests'"
        exit 1
    fi
    echo "🚀    Deploying '$name' to '$env'. Context: '$path'"
    kubectl apply -f "$manifests"
}

deploy "microgateway" "$env" .

for ms in $(ls -d microservices/*); do
    name=$(basename "$ms")
    path="$ms"
    deploy "$name" "$env" "$path"
done
