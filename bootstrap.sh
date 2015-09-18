#!/usr/bin/env bash

# Update aptitude
sudo apt-get update
sudo apt-get upgrade -y

# Install opencv & vim.
sudo apt-get install -y libopencv-dev
sudo apt-get install -y vim

#Ensure webcam module is loaded.
sudo modprobe uvcvideo

if [ ! -f go1.4.2.linux-arm~multiarch-armv7-1.tar.gz ]; then
	wget http://dave.cheney.net/paste/go1.4.2.linux-arm~multiarch-armv7-1.tar.gz
fi

sudo tar -C /usr/local -xzf go1.4.2.linux-arm~multiarch-armv7-1.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> /home/pi/.bashrc
source /home/pi/.bashrc
