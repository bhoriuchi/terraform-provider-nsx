#!/bin/sh

export TF_PROVIDER_PKG=nsx
export TF_PROVIDER_NAME="terraform-provider-$TF_PROVIDER_PKG"
export GOCACHE="$PWD/dist/gosrc_cache"

if [ -z "$GOX_OS_ARCH" ]; then
    export GOX_OS_ARCH="darwin/amd64 linux/amd64 windows/amd64"
fi

mkdir -p $GOCACHE/{src,bin}

if [ ! -d "$GOPATH" ]; then
    export GOPATH="$GOCACHE"
fi

# fetch godeps
echo "Fetching Dependencies... (this may take a while the first run)"
echo

docker run --rm \
  -w//build \
  -v/$PWD://build \
  -v/$GOCACHE/src://go/src \
  -v/$GOCACHE/bin://go/bin \
  golang:1.8 \
  go get -v github.com/mitchellh/gox

docker run --rm \
  -w//build \
  -v/$PWD://build \
  -v/$GOPATH/src://go/src \
  -v/$GOCACHE/bin://go/bin \
  golang:1.8 \
  go get -v -d ./...

# build our provider inside a go container
echo
echo "Building Provider --> build/dist/$TF_PROVIDER_NAME"
echo

docker run --rm \
  -w//build \
  -v/$PWD://build \
  -v/$PWD/$TF_PROVIDER_PKG://build/$TF_PROVIDER_PKG \
  -v/$PWD/dist://build/dist \
  -v/$GOPATH/src://go/src \
  -v/$GOCACHE/bin://go/bin \
  golang:1.8 \
  gox -osarch="$GOX_OS_ARCH" -output "dist/provider_{{.OS}}_{{.Arch}}"

cd "$PWD/dist"

for filename in $PWD/provider_*; do
  chmod +x $filename
  mv $filename $TF_PROVIDER_NAME
  tar -czvf "${filename##*/provider_}.tgz" $TF_PROVIDER_NAME
  rm -f $TF_PROVIDER_NAME
done

echo
