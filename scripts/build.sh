#!/usr/bin/env bash

set -e

tag=$(git describe --abbrev=0 --tags)

function build() {
    name="$1"
    tag="$2"
    path="$3"
    image="$name:$tag"

    echo "🏗    Building '$image'. Context: '$path'"
    docker build -t "$image" "$path"
}

build "microgateway" "$tag" .

for ms in $(ls -d microservices/*); do
    name=$(basename "$ms")
    path="$ms"
    build "$name" "$tag" "$path"
done
