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

# Day 2: Authentication and Authorization

## Users

In our app, users have different types: customers, merchants, and admins.

Each one has to be treated differently at sign up.

- Customers: Customers are the simplest of three. It requires no permissions,
  and there is no additional tasks required rather than adding the new user.

- merchants: There are no permissions needed for merchants, however we should
  also add a store for the merchant. This had an implication on how we design the
  services for each layer.

- admins: Registering an admin can only be done by an existing admin. That's
  why we have seeded an admin in the migrations.

The difference between the creation of each user type have resulted in
separation between them in both the API layer and service layer, while being
united in the models layer with the User model.

### Unit of Work

The merchant users raises a question on the repository layer. Should the user
repository only deal with the user table. This means the creation would happen
on the user repository and the creation of store is in the store repository.
However this raises a problem in case of failures, what will happen if the
merchant user has been created while the store failed for any reason such as
network failures?

One answer is to use the Unit of Work, which mean passing transactions to the
repositories and committing transactions at the end.this however might add
unnecessary complexity for a small scale application like this.

So I have chosen to have a method in the user repository that starts a
transaction inside and add both a merchant and store. This is also closer to
the ideas of Domain Driven Design where there should be a root aggregate
entity, where entities attached to it can be manipulated by its repository.

## Session Management

This application uses JWT as the main session management option. Other options
might be using of stateful tokens. While I might prefer stateful tokens because
it provide me more control over the sessions since I can easily delete them
from my caching system, it will require me to setup a caching system like Redis
or implement my own simple cache (which I implemented in my Last Project
available on my GitHub called Showtime).

# Day 3: FAQ Design

## Database and Model design

### Language

The next phase is to implement the FAQ, FAQ categories, and the translations.
The first design decisions is where to handle supported languages, should we
make a database table for languages and reference it in the FAQs or
translations tables? or should we use go Enums? I opted for the second option
since it was not mentioned in the design document that there would be a
language table.

### FAQ Categories

Although based on database relations between tables, the FAQ category has many
FAQs and each FAQ has only one FAQ category. I have chosen to include the
Category inside the FAQ model so the model is encapsulated with categories,
languages, and translation. We can consider the model as the root aggregate in
this scenario.

### Store FAQs vs Global FAQs

For each FAQ there is two possible outcomes. Local and set to a store or
Global. Based on this I've decided that there would be is_global boolean field
and nullable store_id field. This would allow the same table to contain both
store-specific and global FAQs.
