#!/usr/bin/env bash

set -e

docker-compose -f docker-compose.redis.yml up -d

tmux new-session -d -s mgw
tab=0

function new_tab() {
    name="$1"
    path="$2"
    tab=$(($tab + 1))
    tmux new-window -t mgw:"$tab" -n "$name"
    tmux send-keys -t mgw:"$tab" "cd $path; make run" enter
}

new_tab "microgateway" .

for ms in $(ls -d microservices/*); do
    name=$(basename "$ms")
    path="$ms"
    new_tab "$name" "$path"
done

tmux rename-window -t mgw:0 'workspace'
tmux select-window -t mgw:0

tmux attach -t mgw
