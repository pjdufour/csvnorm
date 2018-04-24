#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "Testing sample.csv"
cat $DIR/../examples/sample.csv | $DIR/../bin/csvnorm_linux_amd64
echo "***********************************************"
echo ""
echo "Testing sample-with-broken-utf8.csv"
cat $DIR/../examples/sample-with-broken-utf8.csv | $DIR/../bin/csvnorm_linux_amd64
echo "***********************************************"
echo ""
