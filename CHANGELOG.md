# Change Log

ping route

## Unreleased

### Added
 - Unittests (tests.py)


## 1.4.0 - 2016-03-18

### Added
 - Create customer/apikey route
 - Create apikey to database for datasource permissions
 - customer apikey required for reading and writing datasources
 - NewFeatureHandler EOF error


## 1.4.0 - 2016-03-16

### Added
 - Load datasource route
 - Unload datasource route
 - View loaded datasources route
 - Gracefull shutdown. Waits to shut down app until all websockets are disconnected
 - server profile route (uptime, runtime, server status)
 - Sphinx documentation templates

### Changed
 - get_requirements.sh checks for src packages before downloading
 - json response from delete layer fix


## 1.4.0 - 2016-03-09

### Changed
 - Cleaned up logging syntax
 - Better logging messages
 - Fixed database lock bug on no datasource found
 - Improved error handling and messaging for database
 - Improved error handling and messaging for GET routes
 - Fixed http status codes for GET routes


## 1.3.0 - 2016-02-28

### Added
 - Version flag (-v)
 - broadcastAllDsViewers for viewer count messaging

### Changed
 - use toGeoJSON to send feature payloads.
 - refactor of sendFeature in find.draw

### Removed
 - package feature and featuretypes removed from find.draw


## 1.2.0 - 2016-02-23

### Changed
 - Fixed logging messages for socket handlers

### Added
 - Viewer count to map.html template
 - send json through websocket with viewer count and instructions to update layer
 - Touch screen support for drawing features (https://github.com/michaelguild13/Leaflet.draw) 

### Removed
 - Redundant logging from socketHandlers.go


## 1.1.0 - 2016-02-22

### Added
 - Added Socket hub

Initial Release

### Added
 - Web interface for drawing.
 - Go based RESTful GeoJson server.
 - Bolt database
 - Cache layer for database