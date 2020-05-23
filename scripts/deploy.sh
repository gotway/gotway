#!/usr/bin/env bash

set -e

env=$1
if [ -z "$env" ]; then
    echo "❌ Environment argument is mandatory"
    exit 1
fi

configmap="manifests/01_configmap_$env.yml"
if [ ! -f "$configmap" ]; then
    echo "❌ Environment is invalid: $env"
    exit 1
fi

echo "🚀 Deploying to $env ..."
kubectl apply -f manifests/00_namespace.yml
kubectl apply -f "$configmap"
kubectl apply -f manifests/microgateway.yml