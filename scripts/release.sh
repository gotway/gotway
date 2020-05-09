#!/usr/bin/env bash

set -e

docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD

tag=$(git describe --abbrev=0 --tags)
project="microgateway"
image="gosmogolang/$project:$tag"
platform="linux/amd64,linux/arm64"

echo "👷   Creating builder $project ..."
docker buildx create --name "$project"
docker buildx use "$project"

echo "🏗    Building $image ..."
docker buildx build --platform "$platform" -t "$image" --push .
docker buildx imagetools inspect "$image"
