# SYSTEM DESCRIPTION:

GeoPathPlanner is a distributed web-application used to find the shortest path between a set of geographical locations for an aerial (or sea) vehicle avoiding obstacles in a particular searching area. This will be achieved implementing path planning algorithms such that RRT, RRT*, and Bug Path. The project architecture is designed to be distributed as a series of microservices.

The system employs a distributed microservice architecture consisting of four main components:
*   **Frontend:** A React application for map interaction and route visualization.
*   **API:** The central entry point for all frontend routing requests
*   **User Management:** Manages user accounts, authentication, and user-profile data.
*   **Routing:** Responsible for executing pathfinding algorithms and computing optimal routes.

Communication between services is primarily via REST APIs for synchronous interactions and Kafka for asynchronous messaging, particularly for routing requests.

# USER STORIES:

1)  As a visitor, I want to try the app without creating an account so that I can quickly see how it works.
2)  As a user, I want to set points on a map where my drone should go so that I can get a possible route.
3)  As a user, I want to mark areas that the drone should avoid so that I can plan a safe path.
4)  As a user, I want to have options to customize which route to compute so to choose tradeoff between computation speed and quality of resulted route.
5)  As a user, I want to see the proposed drone route displayed on the map so that I understand how the drone would fly.
6)  As a user, I want to download or share the proposed route so that I can reuse it outside the app.
7)  As a visitor, I want to sign up into the application so that I can create my personal account to manage my data and routes.
8)  As a registered user, I want to log in to my personal space so that I can manage my routes.
9)  As a registered user, I want to log out to my personal account so that anyone can use my account.
10)  As a registered user, I want to edit my profile information.
11)  As a user, I want to see a clear explanation when no route can be found so that I understand what went wrong.
12)  As a user, I want to easily repeat my last request so that I can test it again.
13)  As a user, I want the option to adjust my inputs and try again after completing or failing request si I don’t have to start from scratch.
14)  As a user, I want the app to generate random no-fly areas so that I can quickly test different scenarios without having to draw them manually.
15)  As a user, I want the app to provide a easy way to clean the window so that I can easily restart over.
16)  As a registered user, I want to see a history of my past routes so that I can review or reuse them later.
17)  As a registered user, I want to delete past routes from my history so that I can no longer see routes that I don't want anymore.
18)  As a user, I want to load waypoints and constraints from geojson file, so that I can easily perform request without having to manually define it.

# CONTAINERS:

## CONTAINER_NAME: Frontend

### DESCRIPTION:
Provides an interactive single-page application where users can define waypoints and obstacles, configure routing parameters, and visualize computed routes. It acts as the primary interface for user interaction with the GeoPathPlanner system.

### USER STORIES:
1) As a visitor, I want to try the app without creating an account so that I can quickly see how it works.
2) As a user, I want to set points on a map where my drone should go so that I can get a possible route.
3) As a user, I want to mark areas that the drone should avoid so that I can plan a safe path.
5) As a user, I want to see the proposed drone route displayed on the map so that I understand how the drone would fly.
6) As a user, I want to download or share the proposed route so that I can reuse it outside the app.
7) As a visitor, I want to sign up into the application so that I can create my personal account to manage my data and routes.
8) As a registered user, I want to log in to my personal space so that I can manage my routes.
9) As a registered user, I want to log out to my personal account so that anyone can use my account.
10) As a registered user, I want to edit my profile information.
11) As a user, I want to see a clear explanation when no route can be found so that I understand what went wrong.
12) As a user, I want to easily repeat my last request so that I can test it again.
13) As a user, I want the option to adjust my inputs and try again after completing or failing request si I don’t have to start from scratch.
14) As a user, I want the app to generate random no-fly areas so that I can quickly test different scenarios without having to draw them manually.
15) As a user, I want the app to provide a easy way to clean the window so that I can easily restart over.
16) As a registered user, I want to see a history of my past routes so that I can review or reuse them later.
18) As a user, I want to load waypoints and constraints from geojson file, so that I can easily perform request without having to manually define it.

### PORTS:
8081:80

### PERSISTANCE EVALUATION
The Frontend container does not require data persistence.

### EXTERNAL SERVICES CONNECTIONS
- Connects to the *API* to send user requests and receive route data.
- Connects to the *User Management* microservice to manage user sessions and registrations.

### MICROSERVICES:

#### MICROSERVICE: frontend
- TYPE: frontend
- DESCRIPTION: Serves the main user interface for the application, built as a single-page application.
- PORTS: 8081:80
- TECHNOLOGICAL SPECIFICATION:
  - React (with Vite)
  - JavaScript
  - Bootstrap 5 for UI components
  - React-Leaflet for map rendering
  - Leaflet-Draw for map interactions
  - Axios for API calls
- SERVICE ARCHITECTURE:
The frontend is a React-based single-page application. The codebase is structured into several directories:
    - `pages`: Contains the main pages of the application (Homepage, Login, Profile, History).
    - `components`: Contains reusable UI components like the Map, Sidebar, and Modals.
    - `services`: Handles API communication with the backend.
    - `context`: Manages application-wide state, such as authentication.
- PAGES:

| Name | Description | Related Microservice | User Stories |
| ---- | ----------- | -------------------- | ------------ |
| Homepage.jsx | Main page with map and sidebar for route planning. | api, routing | 1, 2, 3, 4, 5, 11, 12, 13, 14, 15, 18 |
| Login.jsx | Handles user login. | user-management | 6, 7, 8 |
| Profile.jsx | Displays and allows editing of user profile. | user-management | 10 |
| History.jsx | Displays the user's past routes. | api | 16, 17 |

## CONTAINER_NAME: API

### DESCRIPTION:
Serves as the single entry point for all routing requests, abstracting the underlying microservice. It handles request routing, authentication of routing requests (by validating JWTs with the User Management Service Secret JWT key), and forwards routing requests to the Routing Manager via Kafka.

### USER STORIES:
1) As a visitor, I want to try the app without creating an account so that I can quickly see how it works.
11) As a user, I want to see a clear explanation when no route can be found so that I understand what went wrong.
16) As a registered user, I want to see a history of my past routes so that I can review or reuse them later.
17) As a registered user, I want to delete past routes from my history so that I can no longer see routes that I don't want anymore.

### PORTS:
8000:8000

### PERSISTANCE EVALUATION
The API-Gateway container requires persistent storage to maintain details about the user's route history. It connects to a PostGIS database.

### EXTERNAL SERVICES CONNECTIONS
- Connects to **Kafka** to send routing requests to the Routing Manager and receive results.
- Connects to a **PostGIS** database to store and retrieve route history.

### MICROSERVICES:

#### MICROSERVICE: api
- TYPE: backend
- DESCRIPTION: Manages API requests from the frontend, communicates with other backend services, and handles route history.
- PORTS: 8000:8000
- TECHNOLOGICAL SPECIFICATION:
  - Python
  - FastAPI
- SERVICE ARCHITECTURE:
The service is built with FastAPI and is structured into modules:
    - `main.py`: The main application file with endpoint definitions.
    - `kafka.py`: Handles producing and consuming messages from Kafka.
    - `db.py`: Manages the connection and queries to the PostgreSQL database.
    - `token.py`: Handles JWT token validation using the *User Management* JWT Secret Key.
- ENDPOINTS:

| HTTP METHOD | URL | Description | User Stories |
| ----------- | --- | ----------- | ------------ |
| POST | /api/route/compute | Initiates a route computation request. | 4, 5 |
| GET | /api/history | Retrieves user's route history. | 16 |
| DELETE | /api/history/{id} | Deletes a specific route from history. | 17 |

- DB STRUCTURE:

**_routes_** : | **_request_id_** | user_id | response | created_at | updated_at |

## CONTAINER_NAME: User-Management

### DESCRIPTION:
Manages all user-related functionalities, including user registration, login and profile management. It issues and validates JSON Web Tokens (JWT) for secure authentication.

### USER STORIES:
7) As a visitor, I want to sign up into the application so that I can create my personal account to manage my data and routes.
8) As a registered user, I want to log in to my personal space so that I can manage my routes.
9) As a registered user, I want to log out to my personal account so that anyone can use my account.
10) As a registered user, I want to edit my profile information.

### PORTS:
3000:3000

### PERSISTANCE EVALUATION
The User-Management container requires persistent storage to manage user credentials and profile information. It uses a PostgreSQL database.

### EXTERNAL SERVICES CONNECTIONS
Connects to its own PostgreSQL database.

### MICROSERVICES:

#### MICROSERVICE: user-management
- TYPE: backend
- DESCRIPTION: Handles user creation, authentication, and profile management.
- PORTS: 3000:3000
- TECHNOLOGICAL SPECIFICATION:
  - NestJS (Node.js)
  - TypeScript
  - PostgreSQL with TypeORM
- SERVICE ARCHITECTURE:
The service is a NestJS application organized into modules:
    - `AuthModule`: Handles JWT-based authentication, including login and token validation.
    - `UsersModule`: Manages user CRUD operations (Create, Read, Update, Delete).
- ENDPOINTS:

| HTTP METHOD | URL | Description | User Stories |
| ----------- | --- | ----------- | ------------ |
| POST | /users | Create a new user. | 7 |
| POST | /auth/login | Authenticate user and issue JWT. | 8 |
| GET | /users/{id} | Retrieve user details. | 10 |
| PUT | /users/{id} | Update user details. | 10 |
| POST | /auth/logout | Logout request | 9 |

- DB STRUCTURE:

**_users_** : | **_user_id_** | username | email | password | country | name | surname | created_at | updated_at |

## CONTAINER_NAME: Routing

### DESCRIPTION:
The core computational engine of the system. It listens for routing requests via Kafka, applies various pathfinding algorithms (e.g., RRT, RRT*, Bug Path), and publishes the computed route results back to Kafka.

### USER STORIES:
4)  As a user, I want to have options to customize which route to compute so to choose tradeoff between computation speed and quality of resulted route.
11)  As a user, I want to see a clear explanation when no route can be found so that I understand what went wrong.

### PORTS:
None

### PERSISTANCE EVALUATION
The Routing container does not require data persistence. It is a stateless computational service.

### EXTERNAL SERVICES CONNECTIONS
Connects to **Kafka** to consume routing requests and produce computed routes.

### MICROSERVICES:

#### MICROSERVICE: routing
- TYPE: backend
- DESCRIPTION: Executes pathfinding algorithms based on requests received from the API via Kafka.
- PORTS: None
- TECHNOLOGICAL SPECIFICATION:
  - Go
- SERVICE ARCHITECTURE:
The service is a Go application that runs as a background worker.
    - It initializes a Kafka consumer to listen for messages on the `routing-requests` topic.
    - Upon receiving a request, it deserializes the data and invokes the appropriate pathfinding algorithm from the `algorithm` package.
    - After computation, it serializes the result and produces a message to the `routing-results` Kafka topic.
- KAFKA TOPICS:

| Topic Name | Description |
| ----------- | ----------- |
| `routing-requests` | Topic for incoming route computation requests. |
| `routing-results` | Topic for outgoing computed route results. |