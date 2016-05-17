# Scout

This software powers measure the future 'scouts'. These are web cam based devices that monitor activity in and around various spaces.

![alpha](https://img.shields.io/badge/stability-alpha-orange.svg?style=flat "Alpha")&nbsp;
 ![GPLv3 License](https://img.shields.io/badge/license-GPLv3-blue.svg?style=flat "GPLv3 License")

## Compilation/Installation (Edison)

1. [Upgrade the Firmware on your Intel Edison to Yocto 3.0](http://reprage.com/post/bootstrapping-the-intel-edison/).
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
	$ go get github.com/onsi/ginkgo
	$ go get github.com/onsi/gomega
	$ go get github.com/shirou/gopsutil
	$ go get github.com/MeasureTheFuture/scout
	$ go build github.com/MeasureTheFuture/scout
```

## Connecting to a mothership and running (Edison)

```
	$ configure_edison --wifi
```
The name of the mothership access point will be the same as the device name you configured on the mothership. The password will be the same as the root password you supplied when you ran configure_edison on the mothership.

To run the scout:
```
	$ ./scout
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
	$ go get github.com/shirou/gopsutil
	$ go get github.com/MeasureTheFuture/scout
```

## Operating Instructions:

The scout is a command line application that broadcasts interaction data to 'mothership' or any other location for agregation/reporting.

```
	$ ./scout -help


	  Usage of ./scout:
      -configFile string
    	The path to the configuration file (default "scout.json")
      -debug
    	Should we run scout in debug mode, and render frames of detected materials
      -logFile string
    	The output path for log files. (default "scout.log")
      -videoFile string
    	The path to a video file to detect motion from instead of a webcam
```

If the configuration doesn't exist at the specified place, the scout will create one for you. The scout will fill it with default values that can be customised.


## API:

* To calibrate the scout (takes a new reference image):
```
	GET http://sco.ut.ip/calibrate
	returns: 200 OK on success.

	Calibrate accepts parameters for tweaking OpenCV settings:

	MinArea            float64 // The minimum area enclosed by a contour to be counted as an interaction.
	DilationIterations int     // The number of iterations to perform while dilating the foreground mask.
	ForegroundThresh   int     // A value between 0 and 255 to use when thresholding the foreground mask.
	GaussianSmooth     int     // The size of the filter to use when gaussian smoothing.
	MogHistoryLength   int     // The length of history to use for the MOG2 subtractor.
	MogThreshold       float64 // Threshold to use with the MOG2 subtractor.
	MogDetectShadows   int     // 1 if you want the MOG2 subtractor to detect shadows, 0 otherwise.

	For example:
	GET http://sco.ut.ip/calibrate?MinArea=1300&MogDetectShadows=1

	Sets the minimum detectable area of a 'person' to be 1300 pixels and enables shadow detection.
```

* Once calibrated, the scout will make the following request to the mothership:
```
	POST http://moth.er.sh.ip/scout/calibrated
	This is a MIME multipart message with an attached file (file:calibrationFrame.jpg)

	Within the request header is the following key "Mothership-Authorization", it
	contains the UUID of the scout.
```

* To start measuring:
```
	GET http://sco.ut.ip/measure/start
	returns: 200 OK on success.
```

* During the measurement phase, the scout will make the following requests to the mothership:
```
	POST http://moth.er.sh.ip/scout/interaction
	This is a MIME multipart message with an attached file (file:interaction.json) containing:
	{
		"UUID":"59ef7180-f6b2-4129-99bf-970eb4312b4b",	// Unique identifier of scout.
		"Version":"0.1",								// Transmission protocol version.
		"Entered":"2015-03-07 11:00:00 +0000 UTC",		// When interaction began, rounded to nearest half hour.
		"Duration":2.3,									// The duration in seconds.
		"Path":[
			{
				XPixels:4,								// x-coordinate of waypoint centroid in pixels.
				YPixels:5,								// y-coordinate of waypoint centroid in pixels.
				HalfWidthPixels:2,						// Half the width of the waypoint in pixels.
				HalfHeightPixels:2,						// Half the height of the waypoint in pixels.
				T:0.5									// The number of seconds since the interaction start.
			}
		]
	}

	Within the request header is the following key "Mothership-Authorization", it
	contains the UUID of the scout.
```

* To stop measuring:
```
	GET http://sco.ut.ip/measure/stop
	returns: 200 OK on success.
```

* On startup, the scout will transmit the log file from its previous run to the mothership:
```
	POST http://moth.er.sh.ip/scout/log
	This is a MIME multipart message with an attached file (file:scout.log).

	Within the request header is the following key "Mothership-Authorization", it
	contains the UUID of the scout.
```

* When the scout is running, it will send periodic health heart beats to the mothership:
```
	POST http://moth.er.sh.ip/scout/heartbeat
	This is a MIME multipart message with an attached file (file:heartbeat.json) containing:
	{
		"UUID":"59ef7180-f6b2-4129-99bf-970eb4312b4b",	// Unique identifier of scout.
		"Version":"0.1",								// Transmission protocol version.
		"Health":{
			"IpAddress":"10.1.1.1",						// Current IP address of the scout.
			"CPU":0.4,									// Current CPU load, 0.0 - no load.
			"Memory":0.1,								// Current Memory usage, 0.0 - not used, 1.0 all used.
			"TotalMemory":1233312.0,					// Total Memory available in bytes.
			"Storage":0.1								// Current Storage usage, 0.0 - not used, 1.0 all full.
		}
	}

	Within the request header is the following key "Mothership-Authorization", it
	contains the UUID of the scout.
```

## TODO:
* Build a couple more test datasets that are a bit more complicated (multiple people popping in and out of the frame).

* Do more testing on the Edison. I have been just developing locally on my laptop.
	* ~~Setup a test with mothership on edison, and latest code running on laptop.~~
		* ~~Depends on mothership supporting API.~~
	* Long running tests / memory leaks and any other hardware issues.
	* Multiple people testing.
		* Depends on mothership reporting data from the scout.

* Look at using the calibration frame to 'refresh' the foreground subtractor.
	* Calibration frame of foreground subtractor could also be periodically updated when we have no
	* people detected in the frame (to compensate for subtle lighting changes).

* Build dependencies into ipk so that they are easier to install on the Edison.
	* OpenCV
	* Nginx
	* Postgres
	* Setup a little spot on reprage to host them.

* Upgrade scout from OpenCV 3.0 to OpenCV 3.1

* Configure and setup monit to automatically start and monitor the scout application on the Edison.

## Completed:
* ~~Update bootstrap.sh with the latest dependencies used by the scout.~~
* ~~Fix bug in simplification code (demo-data.mov)~~
* ~~Remove UUID from URL structure used to communicate with mothership.~~
	* ~~Remove UUID from calibrate and log endpoints (shifting it inside form metadata).~~
* ~~Rename mothership URL endpoint from scout to scout_api.~~
* ~~Free C.free(unsafe.Pointer(file)) for calibration image.~~
* ~~Clean up existing code:~~
	* ~~Add configuration file.~~
	* ~~Command line option for overriding location of config file.~~
	* ~~Command line options to specify video file or live feed from webcam.~~
	* ~~Update monitor to loop while frames are available from the camera source.~~
	* ~~Make sure all the metadata fields are populated in the scene when tracking people (frame times).~~
	* ~~Make sure the calibration frame is always the one first pushed into the foreground subtractor.~~
	* ~~Remove debug code from monitor, or add an optional flag for including it.~~
* ~~Cleanup up the OpenCV bindings, and bundle them with the other third party-dependencies.~~
* ~~Update compilation / installation instructions to suit.~~
* ~~Start implementing the communication protocol with the mothership.~~
	* ~~Need to store the mothership ip address/endpoint in the configuration.~~
	* ~~Need to write tests for inbound API.~~
	* ~~Expand calibrate GET request to accept parameters (that can be adjusted on the scout).~~
	* ~~Need to implement calibrate response -> transmit calibration frame to the mothership.~~
	* ~~Go over protocol document and double check that I'm sending everything that needs to be transmitted.~~
		* ~~Need to send version identifiers with all communication to the motherhsip.~~
		* ~~Need to include UUID with interactions transmitted to mothership.~~
		* ~~Finish implementing new interaction utility in interaction.go~~
	* ~~Health heart beat:~~
		* ~~Send first heart beat on startup, and then every 15 minutes after that.~~
		* ~~Get IP address for transmission~~
		* ~~Get memory usage.~~
		* ~~Get Disk usage.~~
		* ~~Get CPU usage.~~
		* ~~Transmit any application error logs.~~
	* ~~Implement Douglas-Peucker to simplify interaction pathway before transmission (to reduce data size of transmitted interactions).~~
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
