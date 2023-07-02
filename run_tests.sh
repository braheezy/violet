#!/bin/bash

set -euo pipefail

YELLOW='#f9e2af'
GREEN='#a6e3a1'
BLUE='#89b4fa'
DARK_BLUE='#1e66f5'
PURPLE='#cba6f7'
DARK_PURPLE='#8839ef'

# Get the absolute path of the script's directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Change the working directory to the script's directory
cd "$SCRIPT_DIR"

print_test_name() {
  gum style --foreground $BLUE --bold -- "$1"
}

print_delimiter() {
  gum style \
    --foreground $DARK_PURPLE \
    -- \
    '-------------------------------'
}

print_header() {
  gum style --background '#8839ef' --foreground '#f9e2af' --margin "1 1" --padding "1 4" --bold "Violet Integraton Test"
}

run_test_file() {
  local test=$(print_test_name "$(basename $1)")

  # Set up test header
  print_header
  print_delimiter
  gum style \
    --foreground $PURPLE \
    --bold \
    --padding "0 1" \
    --align "center" \
    "🚀 Running $test 🚀"

  # Run test file
  set +e
  bash "$file"
  retval=$?
  set -e

  # Handle test result
  if [ $retval -ne 0 ]; then
    gum style \
      --foreground $PURPLE \
      --padding "0 1" \
      --align "center" \
      "🚨 Failed: $test 🚨"
  else
    gum style \
      --foreground $PURPLE \
      --padding "0 1" \
      --align "center" \
      "✅ Success: $test ✅"
  fi
}

# Execute Go tests, defined in Makefile
make test

# Execute all the scripts in the tests directory
for file in test/test*.sh; do
  if [[ -f "$file" ]]; then
    run_test_file "$file"
  fi
done
