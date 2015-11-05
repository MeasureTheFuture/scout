# Scout

This software powers measure the future 'scouts'. These are web cam based devices that monitor activity in and around various spaces.

## Compilation/Installation

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
