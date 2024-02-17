#!/bin/bash

set -eou pipefail

VAGRANTFILE=$(cat <<EOF
Vagrant.configure("2") do |config|
    config.vm.provider "docker" do |d|
        d.image = "alpine"

        # Keep the container running
        d.cmd = ["tail", "-f", "/dev/null"]
    end
end
EOF
)

# Get the absolute path of the script's directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Change the working directory to the script's directory
cd "$SCRIPT_DIR"

for i in {1..5}; do
    TEST_DIR="many-env/env$i"
    if [ -d "$TEST_DIR" ]; then
        pushd "$TEST_DIR" 2>/dev/null
            vagrant destroy -f &>/dev/null || true
        popd
        rm -rf "$TEST_DIR"
    fi

    mkdir -p "$TEST_DIR"
    pushd "$TEST_DIR" 2>/dev/null
        echo "$VAGRANTFILE" > Vagrantfile
        vagrant up
    popd
done
