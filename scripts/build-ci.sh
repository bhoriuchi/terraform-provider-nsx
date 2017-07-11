#!/bin/bash
#
# This script builds the application from source for multiple platforms.

# Get the parent directory of where this script is.

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd "$DIR"

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
XC_OS=${XC_OS:-linux darwin windows}

gox \
	    -os="${XC_OS}" \
	        -arch="${XC_ARCH}" \
		    -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}" \
		        ./...

# Done!
echo
echo "==> Results:"
ls -hl dist/*
