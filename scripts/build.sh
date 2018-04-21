#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

mkdir -p $DIR/../bin
NAME=csvnorm

echo "******************"
echo "Formatting $DIR/../cmd/$NAME"
cd $DIR/../cmd/$NAME
go fmt
echo "Done formatting."
echo "******************"
echo "Building program for $NAME"
cd $DIR/../bin
for GOOS in darwin linux windows; do
  GOOS=${GOOS} GOARCH=amd64 go build -o $NAME"_${GOOS}_amd64" github.com/pjdufour/$NAME/cmd/$NAME
done
if [[ "$?" != 0 ]] ; then
    echo "Error building program for $NAME"
    exit 1
fi
echo "Executable built at $(realpath $DIR/../bin/$NAME)"
