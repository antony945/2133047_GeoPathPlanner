# Student Technical Documentation: GeoPathPlanner

## 1. Introduction

This document provides a technical overview of the Distributed Geographical Route Planning System, "GeoPathPlanner." The system is a web-based application designed to calculate optimal routes for aerial and maritime vehicles, incorporating obstacle avoidance and various pathfinding algorithms. It is built upon a scalable microservice architecture to ensure distributed processing and efficient resource utilization.

## 2. System Architecture Overview

GeoPathPlanner employs a distributed microservice architecture consisting of four main components:
*   **Frontend:** A user-facing React application for map interaction and route visualization.
*   **API Gateway:** The central entry point for all frontend requests, handling authentication and routing.
*   **User Management Service:** Manages user accounts, authentication, and user-specific data like route history.
*   **Routing Manager Service:** Responsible for executing pathfinding algorithms and computing optimal routes.

This architecture promotes modularity, scalability, and independent development of each component. Communication between services is primarily via REST APIs for synchronous interactions and Kafka for asynchronous messaging, particularly for routing requests.

## 3. Microservices Details

### 3.1. Frontend Service

*   **Technology Stack:** React (Vite, JavaScript), Bootstrap 5 (UI), React-Leaflet (Mapping), Leaflet-Draw (Map Interaction), Axios (API calls).
*   **Purpose:** Provides an interactive single-page application where users can define waypoints and obstacles, configure routing parameters, and visualize computed routes. It acts as the primary interface for user interaction with the GeoPathPlanner system.
*   **Key Features:**
    *   Interactive map for drawing waypoints and obstacles.
    *   Sidebar for input configuration (waypoints, obstacles, search, geolocalization).
    *   Display of computed routes.
    *   User authentication interface.

### 3.2. API Gateway Service

*   **Technology Stack:** [Specify technologies used, e.g., Node.js/Express, Python/FastAPI, Go/Gin]
*   **Purpose:** Serves as the single entry point for all client requests, abstracting the underlying microservices. It handles request routing, authentication (by validating JWTs with the User Management Service), and orchestrates complex workflows, such as forwarding routing requests to the Routing Manager via Kafka.
*   **Key Responsibilities:**
    *   Request routing to appropriate backend services.
    *   Authentication and authorization.
    *   Load balancing (if implemented).
    *   Protocol translation (e.g., HTTP to Kafka).

### 3.3. User Management Service

*   **Technology Stack:** [Specify technologies used, e.g., Node.js/NestJS, Python/Django, Java/Spring Boot]
*   **Purpose:** Manages all user-related functionalities, including user registration, login, profile management, and storage of user-specific data such as route history and bookmarks. It issues and validates JSON Web Tokens (JWT) for secure authentication.
*   **Key Features:**
    *   User registration and login (JWT-based).
    *   User profile management.
    *   Storage and retrieval of user route history.
    *   Password management.

### 3.4. Routing Manager Service

*   **Technology Stack:** [Specify technologies used, e.g., Go, Python, Java]
*   **Purpose:** The core computational engine of the system. It listens for routing requests via Kafka, applies various pathfinding algorithms (e.g., RRT, RRT\*, Bug Path), and publishes the computed route results back to Kafka.
*   **Key Features:**
    *   Consumption of routing requests from Kafka.
    *   Execution of multiple pathfinding algorithms.
    *   Consideration of various routing parameters (e.g., turning rate, altitude units).
    *   Publication of computed routes to Kafka.
    *   Scalability to handle parallel routing computations.

## 4. Key API Endpoints

### 4.1. API Gateway Endpoints

*   `POST /api/auth/login`: User login.
*   `POST /api/auth/register`: User registration.
*   `POST /api/route/compute`: Initiates a route computation request.
*   `GET /api/user/profile`: Retrieves user profile information.
*   `GET /api/history`: Retrieves user's route history.
*   `DELETE /api/history/{id}`: Deletes a specific route from history.

### 4.2. User Management Service Endpoints (Internal/Accessed via Gateway)

*   `POST /users`: Create a new user.
*   `POST /auth/login`: Authenticate user and issue JWT.
*   `GET /users/{id}`: Retrieve user details.
*   `PUT /users/{id}`: Update user details.
*   `GET /history/{userId}`: Retrieve route history for a user.
*   `DELETE /history/{userId}/{routeId}`: Delete a specific route from a user's history.

### 4.3. Routing Manager Service (Kafka Topics)

*   `routing-requests`: Topic for incoming route computation requests.
*   `routing-results`: Topic for outgoing computed route results.

## 5. Core Interaction Flows

### 5.1. Route Computation Flow

1.  **Frontend:** User defines waypoints and obstacles, sets parameters, and clicks "Compute."
2.  **Frontend:** Sends a `POST /api/route/compute` request to the API Gateway with the routing payload.
3.  **API Gateway:** Validates the user's authentication token.
4.  **API Gateway:** Produces a message to the `routing-requests` Kafka topic, containing the routing request details.
5.  **Routing Manager:** Consumes the message from `routing-requests`.
6.  **Routing Manager:** Executes the specified pathfinding algorithm, considering obstacles and parameters.
7.  **Routing Manager:** Publishes the computed route (or an error message) to the `routing-results` Kafka topic.
8.  **API Gateway:** Consumes the result message from `routing-results`.
9.  **API Gateway:** Responds to the Frontend with the computed route or an error.
10. **Frontend:** Displays the route on the map or shows an error message.

## 6. Deployment (Infrastructure as Code)

The entire GeoPathPlanner system is designed for easy deployment and reproducibility using Docker and Docker Compose.

*   **`docker-compose.yml`:** A single `docker-compose.yml` file (located in `source/`) orchestrates the deployment of all microservices (Frontend, API Gateway, User Management, Routing Manager) and their dependencies (e.g., Kafka, Redis, database).
*   **`Dockerfile`s:** Each microservice has its own `Dockerfile` (e.g., `source/frontend/Dockerfile`, `source/api/Dockerfile`, etc.) defining its build process and runtime environment.
*   **Deployment Process:** The system can be built and deployed with a single command: `docker-compose up --build`. This command will build all necessary Docker images, create and start the containers, and set up the network.

## 7. Limitations and Future Work

*   **Current Limitations:**
    *   [Placeholder: Describe any known limitations of the current implementation, e.g., limited algorithm choices, performance bottlenecks, lack of real-time data integration, specific browser compatibility issues.]
    *   [Example: The current obstacle generation is random and lacks real-world data integration.]
    *   [Example: Error handling for specific edge cases in routing algorithms might be basic.]
*   **Future Enhancements:**
    *   [Placeholder: Suggest potential improvements or new features, e.g., integration with real-time weather data, support for more vehicle types, advanced user analytics, improved UI/UX features.]
    *   [Example: Implement a more sophisticated geocoding service for waypoint search.]
    *   [Example: Develop a user dashboard for monitoring past route computations and system performance.]
