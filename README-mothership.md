# Mothership

Experimental software for powering the measure the future 'motherships'. These are devices that collect from the 'scouts' and generate web reports.

![alpha](https://img.shields.io/badge/stability-alpha-orange.svg?style=flat "Alpha")&nbsp;
 ![GPLv3 License](https://img.shields.io/badge/license-GPLv3-blue.svg?style=flat "GPLv3 License")

## Copmilation/installation (Edison)

Installation Instructions for Measure The Future can be found [here](https://github.com/MeasureTheFuture/installer).

## Compilation/installation (OSX)

The following steps will configure an OSX machine as a development environment for the mothership:

1. [Download and Install Go 1.6.3](https://golang.org/dl/)
2. [Download and Install PostgreSQL](http://postgresapp.com/)
3. Create SSL keys for PostgreSQL
```
	$ cd ~/Library/Application\ Support/Postgres/var-9.4/
	$ openssl req -new -text -out server.req
	$ openssl rsa -in privkey.pem -out server.key
	$ rm privkey.pem
	$ openssl req -x509 -in server.req -text -key server.key -out server.crt
	$ chmod og-rwx server.key

```
4. Edit postgres config by uncommenting the line #sss = on (Remove the #)
```
	$ vim postgresql.conf
```
5. Restart Postgres.app
6. Create Databases
```
	$ psql
		postgres=# create database mothership;
		postgres=# create database mothership_test;
		postgres=# create user mothership_user with password 'password';
		postgres=# ALTER ROLE mothership_user SET client_encoding TO 'utf8';
		postgres=# ALTER ROLE mothership_user SET default_transaction_isolation TO 'read committed';
		postgres=# ALTER ROLE mothership_user SET timezone TO 'UTC';
		postgres=# GRANT ALL PRIVILEGES ON DATABASE mothership TO mothership_user;
		postgres=# GRANT ALL PRIVILEGES ON DATABASE mothershop_test TO mothership_user;
		postgres=# \q
```
7. [Download and Install Node v4.4.7](https://nodejs.org/en/)
8. Update NPM
```
	$ sudo npm install npm -g
```
9. Create project folder and get source from Github
```
	$ mkdir mtf
	$ cd mtf
	$ mkdir src
	$ cd src
	$ git clone git@github.com:MeasureTheFuture/mothership.git
	$ cd ..
	$ export GOPATH=`pwd`
	$ go get github.com/labstack/echo
	$ go get github.com/lib/pq
	$ go get github.com/onsi/ginkgo
	$ go get github.com/onsi/gomega
	$ go get -u github.com/mattes/migrate
```
10. Build the backend, create config file and migrate databases.
```
	$ go build mothership
	$ cp src/mothership/mothership.json_example mothership.json
	$ ./bin/migrate -url postgres://mothership_test@localhost:5432/mothership -path ./src/mothership/migrations up
	$ ./bin/migrate -url postgres://mothership_test@localhost:5432/mothership_test -path ./src/mothership/migrations up
```
11. Run the front end. (ctrl-c stops the mothership)
```
	$./mothership
```
12. In a new terminal, build the frontend.
```
	$ cd src/mothership/frontend
	$ npm install
	$ npm run build
```
13. Visit localhost:1323 in your browser.

## Testing

Testing for the backend:
```
	$ go test -p 1 github.com/MeasureTheFuture/mothership/...
```

Testing for the frontend:
```
	$ npm run test
```

## TODO

- [x] Implement summarise proccess - update scout summaries with unprocessed interaction data.
- [x] Add control endpoint for getting scout summary data for the interface.
- [x] Add 'measure' button to the interface.
- [x] Update the interface to display visitor count from backend.
- [x] Mothership generates internal server error (500) when no scouts are present.
- [x] Ability to stop measuring from the UI
- [x] Heatmap rendering
- [x] Heatmap number crunching
- [x] Finish pending scout interaction tests
- [x] Need to render heatmap scale.
- [ ] Better identification that the UI is in measuring mode / working
- [ ] Need to generate histogram / timeline.
- [ ] Write tests for the frontend components
- [ ] Tidy up responsive design (for small screensizes)
- [ ] Need to calculate and store heatmap / statistics overtime.
- [ ] Write tests for summarise and heatmap generation. summarise_test
- [ ] Write tests for all the other cases when an AABB intersects a shaft (shaft_test)
- [ ] An easier way to tweak the OpenCV settings for best results

## License

Copyright (C) 2016, Clinton Freeman

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
