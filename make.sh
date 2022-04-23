#!/bin/sh

OSES=("windows" "darwin" "linux")
ARCHES=("386" "amd64" "arm64")

rm -rf build/
mkdir -p build

for os in ${OSES[@]}
do
 for arch in ${ARCHES[@]}
 do
  echo "Building for $os/$arch..."
  if [ "$os" == "windows" ]
  then
   GOOS=$os GOARCH=$arch go build -o build/GoGSM-$os-$arch.exe .
  else
   GOOS=$os GOARCH=$arch go build -o build/GoGSM-$os-$arch .
  fi
  echo "Done!"
  echo
 done
done
