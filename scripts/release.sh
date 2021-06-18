#!/usr/bin/env bash

set -e

docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD

git fetch --all
tag=$(git describe --abbrev=0 --tags)
name="gotway"
image="gotwaygateway/$name"
platform="linux/amd64,linux/arm64,linux/arm"

echo "🏗    Building image '$image:$tag'..."
docker buildx create --name "$name" --use --append
docker buildx build --platform "$platform" -t "$image:$tag" -t "$image:latest" --push .
docker buildx imagetools inspect "$image:latest"
