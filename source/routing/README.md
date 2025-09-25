# Routing Microservice
This folder contains the **routing microservice** for the GeoPathPlanner project. The microservice is responsible for implementing path-planning algorithms to find feasible paths between 3D waypoints (latitude, longitude, altitude) in environments with obstacles defined in GeoJSON format.

This microservice is implemented in **Go**, utilizing **PostgreSQL** and **Redis** to store temporary geospatial information required by the path-planning algorithms.

## Features

- Computes optimal or feasible paths between 3D waypoints.
- Considers obstacles described as GeoJSON polygons or multipolygons.
- Designed for integration with backend component.

## Supported Algorithms

- RRT - Coming soon
- RRT* - Coming soon
- PRM - Coming soon
- PRM* - Coming soon
- AntPath - Coming soon

docker build -t routing-test -f Dockerfile .
docker build -t routing-test-dev -f Dockerfile-dev .


