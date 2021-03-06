# Data Download

Measure The Future allows you to download all the measurements it has made about the usage of your physical space. This download comes as a zipFile containing:

* scouts.json
* A collection of JPG files (one for each scout).
* scout_summaries.json
* scout_interactions.json
* scout_healths.json

## scouts.json

Contains an array of scouts, one for each connected to the mothership. Each scout has the following format:

```
 {
  "id": 1,
  "uuid": "c91ff28c-f583-43be-adb8-d5c060080441",
  "ip_address": "0.0.0.0",
  "port": 8080,
  "authorised": true,
  "name": "Location 1",
  "state": "measuring",
  "summary": null
 }
```

* **id** Is a database identifier for the scout. It is automatically generated by the database.
* **uuid** Is a unique identifier for the scout, this is generated by the scout and is used as a key for uploading content to the mothership.
* **ip_address** Is the current IP address the scout is listening on.
* **port** Is the current port the scout is listening on.
* **authorised** Has a person authorised the scout (as identified by the UUID above) to send data to this instance of the mothership?
* **state** The current state of the mothership, the available options are 'idle', 'calibrating', 'calibrated', 'measuring'.
* **summary** Unused field.

## scout1.jpg (JPG file collection)

Each scout listed in scouts.json will also have a corresponding JPG file in the zip download. This is the callibration frame as displayed in the User Interface. The database identifier **id** above is used to match the the calibration frame with the scout in question. The calibration frame for a scout with an id '1' will be scout1.jpg.

## scout_summaries.json

Contains an array of interaction summaries, one for each scout. Each summary has the following format:

```
{
  "ScoutId": 1,
  "VisitorCount": 931,
  "VisitTimeBuckets": [[1,2,...,20],[1,2,...,20],...,[1,2,...,20]]
}
```

* **ScoutId** is used to match the summary with the database identifier **id** in scouts.json above.
* **VisitorCount** is the raw visitor count, the total number of interactions recorded within the space.
* **VisitTimeBuckets** Is the data used to generate the heatmap. This is 20 arrays of 20 elements. Each value in this 20x20 grid is the total accumulated interaction time at that place on the calibration frame. The first array is the first (top) row in the calibration image, and the first element within that array is the left most edge of the calibration image. The total accumulated interaction time is the total amount of time (in seconds) spent by all visitors at that part of the physical space.

## scout_interactions.json

Contains an array of interactions, one for each visitor interaction detected by the system. Each interaction has the following format:

```
{
  "Id": 2,
  "ScoutId": 1,
  "Duration": 3.1874993,
  "Waypoints": [[976,375],[1120,208],[1626,301]],
  "WaypointWidths": [[178,323],[176,207],[162,300]],
  "WaypointTimes": [7.642e-06,1.5644069,3.1874993],
  "Processed": true,
  "EnteredAt": "2016-09-16T20:15:00Z"
 }
```

* **Id** Is a database identifier for the interaction. This is automatically generated by the database.
* **ScoutId** Is used to match the interaction with the source scout. The corresponding scout in scouts.json will have the same **Id**.
* **Duration** This is the total amount of time (in seconds) that the scout observed this interaction occuring.
* **Waypoints** Is the path (in pixels) that the interaction took through the calibration frame. The address [0,0] corresponds with the top-left corner of the image.
* **WaypointWidths** This is the matching size (in pixels) of the interaction at each step along the path 'Waypoints'. The size of WaypointWidths and Waypoints will always be the same.
* **WaypointTimes** The offset time (in seconds) from 'EnteredAt' that each step along the path in waypoint occured.
* **Processed** Has this interaction been 'processed' and included as part of the summary as defined in scout_summaries.json?
* **EnteredAt** The time the interaction begun. This date/time is in UTC and deliberately rounded to the nearest 15 minutes. The rounding is an additional privacy protection measure, clumping multiple interactions into occuring at the same time. This to make it more difficult to cross-reference interaction data with other sources of metadata.

## scout_healths.json

Contains an array of scout healths, one for each scout at about a 15 minute interval. These healths give an approximation of the health of the measurement system:

```
 {
  "ScoutId": 1,
  "CPU": 5.09,
  "Memory": 0.34978658,
  "TotalMemory": 1.0066616e+09,
  "Storage": 0.5071577,
  "CreatedAt": "2016-09-16T20:21:45.149642Z"
 }
 ```

* **ScoutId** Is used to match the health with the source scout. The corresponding scout in scouts.json will have the same **Id**.
* **CPU** The five minute load average for the scout at that point in time. In this particular example 5.09 means that on average 4.09 processes were waiting on the CPU. A value of 0 would mean that the system was completely idle.
* **Memory** The percentage of the TotalMemory currently being used on the scout system.
* **TotalMemory** The total memory available on the scout in bytes.
* **Storage** The percentage of the total available storage being used on the scout system.
* **CreatedAt** When the health report was created.
