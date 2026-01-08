# Design Decisions

This document would contain every design decision I take for this project. It's
goal to show my thought process and rationale behind my decisions.

# Day 1: Project Structure

The project focus on building an FAQ API, with that in mind I want to focus on
separation of concerns. I need to have different layers.

## Project layers

Since this is an API, I am going with four layers:

- API Layer
- Services
- Repositories
- Models

### Layer 1: API Layer

In this layer I will define the API contract that I Would build my project
based on. This ensures that the project will always have a clear goal to strive
for.

The API Layer would be responsible for defining and handling all HTTP requests
and responses provided by the API. It will handle session validations using JWT

### Layer 2: Services

The services layer focus on the business logic of the application.

### Layer 3: Repositories

Repositories will handle all communication with the database. In a larger
environment, the services might define a repository interface and an
implementation is provided for a specific database much like a driver. However
in our case this might be a premature abstraction that will make things
complicated with no actual benefits

### Layer 4: Models

This where data models will live. Data models represent the main entities in
our application, it also represents the relations between the entities and any
business logic related to entities themselves. In a domain driven design
codebase, the domain would have been more complex with root aggregates. However
For the sake of simplicity the models would only represents the main entities
in our system
