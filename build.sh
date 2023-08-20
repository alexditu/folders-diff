#!/bin/bash

WORKSPACE="$(dirname "$(realpath "$0")")"
BUILDDIR="${WORKSPACE}/builddir"

# {path, output name} pairs
declare -A BINARIES=(\
["cmd/cli"]="fdiff"
)

build_go_binaries () {
	local OS="$1"
	shift
	local DEST_DIR="$1"
	shift
	local GCFLAGS="$1"
	shift

	local CUSTOM_FLAGS=""
	local OUT_BIN=""

	for k in "${!BINARIES[@]}"; do
		OUT_BIN="${DEST_DIR}/${BINARIES[$k]}"
		if [ "${OS}" == "windows" ]; then
			OUT_BIN="${OUT_BIN}.exe"
		fi

		echo "Building ${k} to ${OUT_BIN}"

		# print build command
		set -x

		go build -o "${OUT_BIN}" -gcflags="${GCFLAGS}" ${CUSTOM_FLAGS} -trimpath "../${k}"

		set +x

	done
}

go version

go mod tidy
go generate ./...

mkdir -p "${BUILDDIR}"
cd "${BUILDDIR}" || exit 1

# Build for Linux
echo -e "\n"

echo "Building for Linux x64"

echo -e "\nBuilding debug binaries"
build_go_binaries "linux" "${BUILDDIR}/linux/bin/amd64/debug" 'all=-N -l'

echo -e "\nBuilding release binaries"
build_go_binaries "linux" "${BUILDDIR}/linux/bin/amd64/release"


echo "Building for Windows x64"
# for full list run command: go tool dist list
export GOOS=windows GOARCH=amd64

echo -e "\nBuilding debug binaries"
build_go_binaries "windows" "${BUILDDIR}/windows/bin/amd64/debug" 'all=-N -l'

echo -e "\nBuilding release binaries"
build_go_binaries "windows" "${BUILDDIR}/windows/bin/amd64/release"