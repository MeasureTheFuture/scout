#!/usr/bin/env bash

# Update opkg.
opkg update
opkg upgrade
opkg install git


# Install golang.
if [ ! -f go1.5.1.linux-386.tar.gz ]; then
	wget https://storage.googleapis.com/golang/go1.5.1.linux-386.tar.gz
fi

mkdir /usr/local
tar -C /usr/local -xzf go1.5.1.linux-386.tar.gz
echo "PATH=$PATH:/usr/local/go/bin" >> /etc/profile


# Install pre-compiled binaries of our third party dependencies.
if [ ! -f opencv-3-edison.tgz ]; then
	wget https://github.com/MeasureTheFuture/scout-dependencies/releases/download/v0.2-alpha/opencv-3-edison.tgz
fi

tar -C /usr/local -xzf opencv-3-edison.tgz


# Configure the go project and get the scout source code
mkdir mtf
mkdir mtf/src
echo "GOPATH=`pwd`/mtf" >> /etc/profile
echo "export GOPATH PATH" >> /etc/profile
source /etc/profile
