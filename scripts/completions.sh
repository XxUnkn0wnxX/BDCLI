#!/bin/sh
set -eu

output_dir="${1:-completions}"

rm -rf "$output_dir"
mkdir -p "$output_dir"

for sh in bash zsh fish; do
	go run main.go completion "$sh" >"$output_dir/bdcli.$sh"
done
