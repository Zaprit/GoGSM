#!/bin/sh

OSES=("windows" "darwin" "linux")
ARCHES=("386" "amd64" "arm64")

for os in ${OSES[@]}
do
 for arch in ${ARCHES[@]}
 do
  echo "Building for $os/$arch..."
  GOOS=$os GOARCH=$arch go build -o GoGSM-$os-$arch .
  echo "Done!"
  echo
 done
done
