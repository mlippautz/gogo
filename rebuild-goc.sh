#!/bin/sh

export GOROOT=$HOME/go
export GOARCH=amd64
export GOOS=linux

echo -n "Do you want to rebuild the go compiler? [y/N]:"
read -n1 -s input
echo ""
if [ -z $input ] || test $input != "y" ; then
    echo "Negative. Exiting"
    exit -1
fi

# delete the old one
rm -rf $GOROOT

# rebuild
hg clone -r release https://go.googlecode.com/hg/ $GOROOT
cd $GOROOT/src
./all.bash

