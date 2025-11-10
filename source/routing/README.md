# 🛰️Routing Microservice (Go)
This folder contains the **routing microservice** for the GeoPathPlanner project. The microservice is responsible for implementing path-planning algorithms to find feasible paths between 3D waypoints (latitude, longitude, altitude) in environments with obstacles defined in GeoJSON format.

This microservice is implemented in **Go**, utilizing **RTree** data structure to store temporary geospatial information required by the path-planning algorithms.

It gets routing requests through a **Kafka client** that:
- Consumes `RoutingRequest` messages from a Kafka topic (`routing_requests`)
- Validates the request (e.g., waypoints, constraints)
- Produces a `RoutingResponse` to another Kafka topic (`routing_responses`)

It uses the **[franz-go](https://github.com/twmb/franz-go)** library for efficient Kafka communication and is fully containerized using Docker Compose.

## Features

- Computes optimal or feasible paths between 3D waypoints.
- Considers obstacles described as GeoJSON polygons or multipolygons.
- Designed for integration with backend component.

## Supported Algorithms

- RRT - ✅
- RRT* - ✅
- AntPath - ✅

## ⚙️ Prerequisites

Before starting, make sure you have installed:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

---

## 🧩 Environment Configuration

All configuration values are stored in a `.env` file at the project root.  
Example:

```bash
KAFKA_BROKERS=kafka:9093
KAFKA_REQUEST_TOPIC=routing-requests
KAFKA_RESPONSE_TOPIC=routing-responses
KAFKA_ROUTING_CONSUMER_GROUP_ID=routing-group
REPLICAS=3
```

> ⚠️ Ensure this `.env` file exists before running `docker compose up`.

---

## 🚀 Running the Project

### 1️⃣ Build and start all containers

From the project root:

```bash
docker compose up --build
```

To run in detached mode:

```bash
docker compose up -d --build
```

Or if you want to see just the `app` logs:

```bash
docker compose up --build --attach app
```

This command will:
- Build the Go application (`app`)
- Start Kafka (in KRaft mode)
- Load environment variables from `.env`
- Wait for Kafka to become healthy before launching the app

---

### 2️⃣ Check that Kafka is running

You can verify that Kafka is up by listing topics:

```bash
docker exec -ti kafka /opt/kafka/bin/kafka-topics.sh \
  --list --bootstrap-server localhost:9092
```

If your topics are not yet created, you can manually create them:

```bash
docker exec -ti kafka /opt/kafka/bin/kafka-topics.sh \
  --create \
  --topic routing_requests \
  --bootstrap-server localhost:9092 \
  --partitions 3 \
  --replication-factor 1

docker exec -ti kafka /opt/kafka/bin/kafka-topics.sh \
  --create \
  --topic routing_responses \
  --bootstrap-server localhost:9092 \
  --partitions 3 \
  --replication-factor 1
```

---

### 3️⃣ (Optional) Test Kafka manually

**Produce a message to the requests topic:**
```bash
docker exec -ti kafka /opt/kafka/bin/kafka-console-producer.sh \
  --bootstrap-server localhost:9092 \
  --topic routing_requests
```
Paste a test JSON message:
```json
{
  "request_id": "abc123",
  "waypoints": [
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [
          4.433724687935722,
          50.872778105839274
        ]
      },
      "properties": {
        "altitudeUnit": "mt",
        "altitudeValue": 200
      }
    },
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [
          4.46992531620532,
          50.884400404439646
        ]
      },
      "properties": {
        "altitudeUnit": "mt",
        "altitudeValue": 300
      }
    },
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [
          4.45503208121508,
          50.890383059561145
        ]
      },
      "properties": {
        "altitudeUnit": "mt",
        "altitudeValue": 400
      }
    }
  ],
  "constraints": [
    {
      "id": 0,
      "type": "Feature",
      "geometry": {
        "type": "Polygon",
        "coordinates": [
          [
            [
              4.452166432497052,
              50.87565312347229
            ],
            [
              4.45389986799907,
              50.874291381256256
            ],
            [
              4.4498672299481825,
              50.88096516631941
            ],
            [
              4.451848154983281,
              50.8813222592305
            ],
            [
              4.45478419822345,
              50.87636729027514
            ],
            [
              4.4560577382486315,
              50.87672438629892
            ],
            [
              4.451812846246014,
              50.883420200613955
            ],
            [
              4.447815509645977,
              50.88214806230394
            ],
            [
              4.452166432497052,
              50.87565312347229
            ]
          ]
        ]
      },
      "properties": {
        "altitudeUnit": "mt",
        "maxAltitudeValue": 500,
        "minAltitudeValue": 400
      }
    },
    {
      "id": 1,
      "type": "Feature",
      "geometry": {
        "type": "Polygon",
        "coordinates": [
          [
            [
              4.435823054794525,
              50.87917754349243
            ],
            [
              4.435999901698551,
              50.876186530052024
            ],
            [
              4.443605337154025,
              50.878195458959425
            ],
            [
              4.439678842862065,
              50.88446720530686
            ],
            [
              4.435823054794525,
              50.87917754349243
            ]
          ]
        ]
      },
      "properties": {
        "altitudeUnit": "ft",
        "maxAltitudeValue": 1000,
        "minAltitudeValue": 100
      }
    },
    {
      "id": 2,
      "type": "Feature",
      "geometry": {
        "type": "Polygon",
        "coordinates": [
          [
            [
              4.443180802952469,
              50.87710169486769
            ],
            [
              4.445940351819587,
              50.874043594523414
            ],
            [
              4.449867003560769,
              50.875472226326934
            ],
            [
              4.447425915923816,
              50.88058383172759
            ],
            [
              4.443180802952469,
              50.87710169486769
            ]
          ]
        ]
      },
      "properties": {
        "altitudeUnit": "mt",
        "maxAltitudeValue": 999999,
        "minAltitudeValue": -999999
      }
    }
  ],
  "search_volume": {
    "type": "Feature",
    "geometry": {
      "type": "Polygon",
      "coordinates": [
        [
          [
            4.424480223895728,
            50.89367115387381
          ],
          [
            4.424480223895728,
            50.867778999101745
          ],
          [
            4.480653371960557,
            50.867778999101745
          ],
          [
            4.480653371960557,
            50.89367115387381
          ],
          [
            4.424480223895728,
            50.89367115387381
          ]
        ]
      ]
    },
    "properties": {
      "altitudeUnit": "mt",
      "maxAltitudeValue": 999999,
      "minAltitudeValue": -999999
    }
  },
  "parameters": {
    "algorithm": "rrtstar",
    "storage": "rtree",
    "goal_bias": 0.1,
    "max_iterations": 10000,
    "max_step_size_mt": 20,
    "sampler": "uniform",
    "seed": 10
  },
  "received_at": "2025-11-01T10:40:13Z"
}
```

**Consume from the responses topic:**
```bash
docker exec -ti kafka /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic routing_responses \
  --from-beginning
```

---

## 🧠 Kafka Configuration Notes

- **Brokers:** The app connects to internal broker addresses defined in `.env` (`kafka:9093`).
- **Topics:**  
  - Input: `routing_requests`  
  - Output: `routing_responses`
- **Group ID:** Each consumer group (`routing_group`) ensures only one instance of the app processes a given message.

Kafka automatically distributes messages across partitions (load balancing).  

---

## 🧹 Stopping and Cleaning Up
To stop the containers:

```bash
docker compose down
```

To remove all volumes (⚠️ this deletes Kafka data):

```bash
docker compose down -v
```

---

## 📘 Development Tips

- Use `docker compose logs -f app` to follow app logs live.
- You can rebuild only the Go app (without restarting Kafka) using:
  ```bash
  docker compose up -d --build app
  ```
- If you’re debugging code inside the container:
  ```bash
  docker exec -it app sh
  ```

---

## 🧩 License

This project is distributed under the MIT License.  
Feel free to use, modify, and distribute it.

---

## 👨‍💻 Author

Developed with ❤️ using Go and Kafka.  
If you find this helpful, give it a ⭐ on GitHub!