
# Booklet: System Architecture

## 1. Introduction

This document provides a detailed overview of the system architecture for the **GeoPathPlanner** project. The primary goal of this project is to create a distributed, scalable, and efficient web application for planning optimal routes for aerial and maritime vehicles while avoiding designated obstacles.

The architecture is designed to support key functional requirements, including interactive map-based input, complex algorithmic processing, user data management, and a responsive user experience.

## 2. Architectural Style: Microservices

We have chosen a **microservices architectural style** to structure our system. This approach decomposes the application into a set of independently deployable services. Each service is organized around a specific business capability.

The key drivers for this decision include:

*   **Scalability:** Individual services (e.g., the `Routing`) can be scaled independently to handle high computational loads without affecting other parts of the system.
*   **Separation of Concerns:** Each microservice has a clearly defined responsibility, making the system easier to understand, develop, and maintain. For example, user authentication is entirely managed by the `User Management` service.
*   **Technological Flexibility:** Each service can be implemented using the most appropriate technology stack. We leverage React for the frontend, Go for high-performance computing in the `Routing`, and NestJS for the `User Management` service.
*   **Resilience:** The use of asynchronous communication patterns ensures that the computationally intensive routing process does not block the backend services, contributing to overall system resilience and scalability. While the frontend displays a loading screen during this waiting period, the backend remains responsive to other requests.

## 3. System Components (Microservices)

The system is composed of four main microservices, containerized using Docker.

### 3.1. Frontend

*   **Description:** A single-page application (SPA) that provides the user interface. It allows users to interact with a map to define waypoints and obstacles, configure routing parameters, and visualize the computed results.
*   **Technology Stack:** React, Vite, React-Leaflet (for map interaction), Bootstrap 5, Axios.
*   **Responsibilities:**
    *   Rendering the map and UI components.
    *   Capturing user input (drawing on the map, uploading files).
    *   Communicating with the backend via the API.
    *   Handling client-side state (e.g., authentication tokens, map data).

### 3.2. API

*   **Description:** The single entry point for all client requests related to route computation and history. It orchestrates the asynchronous pathfinding process and manages the persistence of results.
*   **Technology Stack:** Python (with FastAPI).
*   **Responsibilities:**
    *   **Route Handling:** Manages all `/routes` endpoints, including computation, history retrieval, and deletion.
    *   **Authentication:** Validates JWT tokens locally for all incoming requests to secure the endpoints.
    *   **API Composition:** Orchestrates the asynchronous route computation flow by publishing jobs to a Kafka topic and waiting for a response from another topic.
    *   **Data Persistence:** Saves the final routing results to a database and retrieves them for user history requests.

### 3.3. User Management

*   **Description:** A RESTful service responsible for all user-related concerns.
*   **Technology Stack:** NestJS (a Node.js framework), using JWT for authentication.
*   **Responsibilities:**
    *   User registration and login.
    *   Issuing and validating JSON Web Tokens (JWT).
    *   Managing user data, including route history and preferences.
    *   Providing a secure interface for other services to verify user identity.

### 3.4. Routing Manager

*   **Description:** A backend service dedicated to performing the computationally intensive pathfinding calculations. It operates asynchronously.
*   **Technology Stack:** Go (chosen for its performance and concurrency features), Kafka.
*   **Responsibilities:**
    *   Consuming route computation requests from a Kafka topic.
    *   Executing pathfinding algorithms (AntPath, RRT, RRT*).
    *   Considering waypoints, obstacles, and other constraints during calculation.
    *   Publishing the final computed route to another Kafka topic for the API to retrieve.

## 4. Communication Patterns

The system employs a hybrid communication model to ensure both responsiveness and scalability.

*   **Synchronous (REST):** For immediate, request-response interactions like user login, registration, or fetching route history, the Frontend communicates with the API using standard RESTful API calls over HTTP.
*   **Asynchronous (Kafka):** For the time-consuming route computation process, we use an event-driven approach with Apache Kafka. This decouples the backend services, allowing for scalable and resilient processing, even though the frontend displays a loading state while awaiting the result.

## 5. Data Flow: Route Computation

The asynchronous route computation flow is central to the system's design:

1.  **Request Submission:** The user defines waypoints and obstacles on the **Frontend** and clicks "Compute."
2.  **API Receives Request:** The Frontend sends a `POST` request to the `/route/compute` endpoint on the **API**.
3.  **Job Publishing:** The API validates the request and publishes a message containing the routing parameters (waypoints, obstacles) to a Kafka topic (e.g., `routing_requests`).
4.  **Asynchronous Processing:** The **Routing**, which is subscribed to the `routing_requests` topic, consumes the message and begins the pathfinding calculation. During this time, the **API** waits for the result from Kafka, and the **Frontend** displays a loading screen to the user.
5.  **Result Publishing:** Once the calculation is complete, the Routing Manager publishes the resulting route data to a different Kafka topic (e.g., `routing_results`).
6.  **Result Consumption:** The API listens to the `routing_results` topic. Upon receiving the result corresponding to the user's request, it forwards the data to the Frontend.
7.  **Visualization:** The Frontend receives the computed route and displays it on the map.

## 6. Deployment (Infrastructure as Code)

The entire system is designed to be deployed using **Docker** and **Docker Compose**. This complies to the principle of Infrastructure as Code (IaC).

*   Each microservice has its own `Dockerfile` that defines its environment and dependencies.
*   A master `docker-compose.yml` file in the `source/` directory defines and links all the services (including Kafka and any databases), allowing the entire application stack to be launched with a single command: `docker-compose up --build`.

This approach ensures a consistent, reproducible deployment environment for both development and production.

## 7. Architectural Diagram

**Diagram Description:**
*   A **User** interacts with the **Frontend** (React).
*   For session management the **Frontend** interacts with the **User Management** microservice.
*   The **Frontend** sends all routing requests to the **API** (Python/FastAPI).
*   The **API** validates JWT tokens for authenticated requests (tokens issued by the **User Management** service).
*   For routing requests, the **API** publishes a message to a **Kafka** topic.
*   The **Routing** (Go) consumes from this Kafka topic, computes the route, and publishes the result to another Kafka topic.
*   The **API** consumes the result from Kafka and sends it back to the **Frontend**.
*   The **API** also persists and retrieves route history from its own **Database**.
*   The **User Management** service is connected to a **Database** to persist user data.
