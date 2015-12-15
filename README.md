# Scout

This software powers measure the future 'scouts'. These are web cam based devices that monitor activity in and around various spaces.

![alpha](https://img.shields.io/badge/stability-alpha-orange.svg?style=flat "Alpha")&nbsp;
 ![GPLv3 License](https://img.shields.io/badge/license-GPLv3-blue.svg?style=flat "GPLv3 License")

## Compilation/Installation (Edison)

1. [Upgrade the Firmware on your Intel Edison to Yocto 2.1*](http://reprage.com/post/bootstrapping-the-intel-edison/).
2. ssh into your Edison.
3. Download and run the scout bootstrap script to configure and install all the development tooling:
```
	$ wget https://raw.githubusercontent.com/MeasureTheFuture/scout/master/bootstrap.sh
	$ chmod +x bootstrap.sh
	$ ./bootstrap.sh
```
4. After the bootstrap script has installed go, 3rd-party dependencies and downloaded the scout source code, it can be built with the following:
```
	$ source /etc/profile
	$ go build scout
```

## Compilation/Installation (OSX)

1. [Download & Install Go 1.5.1](https://golang.org/dl/)
2. Install OpenCV-3.0 via [Brew](http://brew.sh/):
```
	$ brew install opencv3
```
3. Create Project structure, build CV Bindings. Download and build the scout:
```
	$ mkdir mtf
	$ cd mtf
	$ git clone https://github.com/MeasureTheFuture/CVBindings.git
	$ cd CVBindings
	$ cmake .
	$ make
	$ cp CVBindings.h /usr/local/opt/opencv3/include
	$ cp libCVBindings.a /usr/local/opt/opencv3/lib
	$ mkdir scout
	$ cd scout
	$ export GOPATH=`pwd`
	$ go get github.com/onsi/ginkgo
	$ go get github.com/onsi/gomega
	$ go get github.com/MeasureTheFuture/scout
```

## API:

* To calibrate the scout (takes a new reference image):
```
	GET http://sco.ut.ip/calibrate
	returns: 200 OK on success.
```

* Once calibrated, the scout will make the following request to the mothership:
```
	POST http://moth.er.sh.ip/scout/<UUID>/calibrated
	This is a MIME multipart message with an attached file (file:calibrationFrame.jpg)
```

* To start measuring:
```
	GET http://sco.ut.ip/measure/start
	returns: 200 OK on success.
```

* During the measurement phase, the scout will make the following requests to the mothership:
```
	POST http://moth.er.sh.ip/scout/<UUID>/interaction
	This is a MIME multipart message with an attached file (file:interaction.json)
```

* To stop measuring:
```
	GET http://sco.ut.ip/measure/stop
	returns: 200 OK on success.
```



## TODO:
* Clean up existing code:
	* ~~Add configuration file.~~
	* ~~Command line option for overriding location of config file.~~
	* ~~Command line options to specify video file or live feed from webcam.~~
	* ~~Update monitor to loop while frames are available from the camera source.~~
	* ~~Make sure all the metadata fields are populated in the scene when tracking people (frame times).~~
	* ~~Make sure the calibration frame is always the one first pushed into the foreground subtractor.~~
	* Look at using the calibration frame to 'refresh' the foreground subtractor.
	* ~~Remove debug code from monitor, or add an optional flag for including it.~~
* ~~Cleanup up the OpenCV bindings, and bundle them with the other third party-dependencies.~~
* ~~Update compilation / installation instructions to suit.~~
* Build a couple more test datasets that are a bit more complicated (multiple people popping in and out of the frame).
* Do some more testing on the Edison. I have been just developing locally on my laptop.
* Start implementing the communication protocol with the mothership.
	* ~~Need to store the mothership ip address/endpoint in the configuration.~~
	* Need to write tests for inbound API.
	* Expand calibrate GET request to accept parameters (that can be adjusted on the scout).
	* ~~Need to implement calibrate response -> transmit calibration frame to the mothership.~~
	* ~~Need to implement broadcasting of interactions to mothership.~~
* ~~Generate UUID on initial startup, store as part of configuration.~~

## License

Copyright (C) 2015, Clinton Freeman

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
