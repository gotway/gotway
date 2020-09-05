#!/usr/bin/env bash

set -e

tmux new-session -d -s grpc
tab=0

function new_tab() {
    name="$1"
    cmd="$2"
    tab=$(($tab + 1))
    tmux new-window -t grpc:"$tab" -n "$name"
    tmux send-keys -t grpc:"$tab" "$cmd" enter
}

new_tab "client" "make cli"
new_tab "server" "make srv"
tmux rename-window -t grpc:0 'workspace'
tmux select-window -t grpc:0

tmux attach -t grpc
