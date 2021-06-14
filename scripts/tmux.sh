#!/usr/bin/env bash

set -e

tmux new-session -d -s gotway
tab=0

function new_tab() {
    name="$1"
    path="$2"
    tab=$(($tab + 1))
    tmux new-window -t gotway:"$tab" -n "$name"
    tmux send-keys -t gotway:"$tab" "cd $path; make run" enter
}

new_tab "gotway" .

for ms in $(ls -d cmd/*); do
    name=$(basename "$ms")
    if [ "$name" = "gotway" ]; then
      continue
    fi
    path="$ms"
    new_tab "$name" "$path"
done

tmux select-window -t gotway:1

tmux attach -t gotway
