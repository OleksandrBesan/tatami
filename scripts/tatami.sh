#!/bin/bash
# Tatami shell wrapper
# Add this to your ~/.zshrc or ~/.bashrc:
#   source /path/to/tatami/scripts/tatami.sh

tatami() {
    local output
    output=$(TATAMI_WRAPPER=1 command tatami "$@")
    local exit_code=$?

    if [[ $exit_code -eq 0 && -d "$output" ]]; then
        cd "$output"
    elif [[ -n "$output" ]]; then
        echo "$output"
    fi
    return $exit_code
}
