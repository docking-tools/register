#! /bin/bash

set -e

# Invoke ./cover for HTML output
COVER=${COVER:-"-cover"}

GO_BUILD_FLAGS=-a

TEST=( ./**/*_test.go )


echo ${TEST}

echo "Running tests..."

	MACHINE_TYPE=$(uname -m)
	if [ $MACHINE_TYPE != "armv7l" ]; then
		RACE="--race"
	fi
	go test -v -timeout 3m ${COVER} ${RACE} $(glide novendor)