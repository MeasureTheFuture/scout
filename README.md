# Scout

This software powers measure the future 'scouts'. These are web cam based devices that monitor activity in and around various spaces.

![alpha](https://img.shields.io/badge/stability-alpha-orange.svg?style=flat "Alpha")&nbsp;
 ![GPLv3 License](https://img.shields.io/badge/license-GPLv3-blue.svg?style=flat "GPLv3 License")

## Compilation/Installation (Edison)

Installation instructions for Measure The Future can be found [here](https://github.com/MeasureTheFuture/installer).

## Compilation/Installation (OSX)

1. [Download & Install Go 1.6](https://golang.org/dl/)
2. Install OpenCV-3.1 via [Brew](http://brew.sh/):
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

### Connecting to a mothership

The wifi on the scout needs to be configured to connect to the Access Point running on the mothership:

```
	$ configure_edison --wifi
```
The name of the access point will be the same as the device name you configured for the mothership. The password will be the same as the root password you supplied when you ran `configure_edison` on the mothership.

### Subsequent uses

At the moment, the scout software doesn't automatically start on boot. Each time the scout is powered up you need to login and run:

```
	$ ./scout
```


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
	POST http://moth.er.sh.ip/scout_api/calibrated
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
	POST http://moth.er.sh.ip/scout_api/interaction
	This is a MIME multipart message with an attached file (file:interaction.json) containing:
	{
		"UUID":"59ef7180-f6b2-4129-99bf-970eb4312b4b",	// Unique identifier of scout.
		"Version":"0.1",								// Transmission protocol version.
		"Entered":"2015-03-07T11:00:00Z",				// When interaction began, rounded to nearest half hour.
		"Duration":2.3,									// The duration in seconds.
		"Path":[
			{
				"XPixels":4,							// x-coordinate of waypoint centroid in pixels.
				"YPixels":5,							// y-coordinate of waypoint centroid in pixels.
				"HalfWidthPixels":2,					// Half the width of the waypoint in pixels.
				"HalfHeightPixels":2,					// Half the height of the waypoint in pixels.
				"T":0.5									// The number of seconds since the interaction start.
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
	POST http://moth.er.sh.ip/scout_api/log
	This is a MIME multipart message with an attached file (file:scout.log).

	Within the request header is the following key "Mothership-Authorization", it
	contains the UUID of the scout.
```

* When the scout is running, it will send periodic health heart beats to the mothership:
```
	POST http://moth.er.sh.ip/scout_api/heartbeat
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

- [x] Filter out interactions that are 'noise', ones that last less than a second.
- [x] Build a couple of extra test datasets that are more complicated (multiple people popping in and out of the frame).
- [ ] Edison testing - Long running tests / memory leaks and other hardware issues.
- [ ] Integration testing.
- [ ] Look at using calibration frame to 'refresh' the foreground subtractor.
- [ ] Calibration frame could also be periodically updated when we have no people detected in the frame (to compenstate for subtle lighting changes).


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
