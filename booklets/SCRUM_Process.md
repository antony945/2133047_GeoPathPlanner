# Booklet: SCRUM Process

## 1. Introduction

This project was developed using the SCRUM agile methodology. We organized our work into two-week sprints, each with a specific goal, a planned backlog, and a review/retrospective session. This approach allowed us to adapt to challenges and deliver value incrementally.

## 2. Sprint 1: Foundation and Core Map Interaction

*   **Duration:** 2 Weeks
*   **Sprint Goal:** Establish the project foundation, including the basic frontend layout, and implement core map functionalities for defining waypoints and obstacles.

### Sprint 1 Backlog

| User Story ID | Description | Status |
| :------------ | :---------- | :------- |
| US-2 | Set waypoints on the map | Done |
| US-3 | Mark no-fly zones (obstacles) | Done |
| US-15 | Add a "clear map" button | Done |
| - | Setup React project with Vite | Done |
| - | Implement basic Homepage layout | Done |

### Sprint 1 Review & Retrospective

*   **Review:** Successfully demonstrated the ability to add and remove markers (waypoints) and polygons (obstacles) on the Leaflet map.
*   **Retrospective:** The initial setup of React-Leaflet with Vite was more complex than expected. We decided to allocate more time for library integration in future sprints.

## 3. Sprint 2: Route Computation Flow (Stubbed) & User Auth UI

*   **Duration:** 2 Weeks
*   **Sprint Goal:** Implement the UI for user authentication and connect the frontend to a stubbed backend to simulate the route computation flow.

### Sprint 2 Backlog

| User Story ID | Description | Status |
| :------------ | :---------- | :------- |
| US-1 | Allow guest access to the map | Done |
| US-7 | Create UI for user sign-up | Done |
| US-8 | Create UI for user login | Done |
| US-4 | Add options to customize route request | In Progress |
| US-5 | Display a fake computed route on the map | Done |

### Sprint 2 Review & Retrospective

*   **Review:** Showcased the login/signup pages and the end-to-end flow of sending a request and displaying a hardcoded route result on the map.
*   **Retrospective:** Managing application state (user auth, map data, results) is becoming complex. We will prioritize implementing a state management solution (like React Context) in the next sprint.

## 4. Burndown Chart

A burndown chart was used to track our progress against the plan for each sprint. The chart plots the remaining effort (in story points) against the days of the sprint.

*(You should include a screenshot of your burndown chart from a spreadsheet or project management tool here.)*
