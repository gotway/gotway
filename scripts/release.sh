#!/usr/bin/env bash

set -e

docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD

tag=$(git describe --abbrev=0 --tags)

function release() {
    name="$1"
    tag="$2"
    path="$3"
    image="gotwaygateway/$name:$tag"
    platform="linux/amd64,linux/arm64,linux/arm"

    echo "🏗    Building '$image'. Context: '$path'"
    docker buildx create --name "$name" --use --append
    docker buildx build --platform "$platform" -t "$image" --push "$path"
    docker buildx imagetools inspect "$image"
}

release "gotway" "$tag" .

for ms in $(ls -d microservices/*); do
    name=$(basename "$ms")
    path="$ms"
    release "$name" "$tag" "$path"
done
