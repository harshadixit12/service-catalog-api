# Service Catalog API
This repository contains an implementation of an API for storing and retrieving services and version - a catalog of services.  
It makes use of Gin HTTP framework (https://gin-gonic.com/docs/), and GORM (https://gorm.io/docs/) as the ORM, and makes use of SQLite (https://pkg.go.dev/github.com/mattn/go-sqlite3) for storage.  

Why Gin?
The framework is performant, and has inbuilt features like middlewares which would take additional effort to implement otherwise.

Why GORM?
It abstracts out the SQL, and makes writing queries easier, and also has drivers to support multiple databases.

Why SQLite?
Easy to set up and use for development, and to mock in tests as well.

# Terminology
1. Organization  
A company that is a customer of our API, and has multiple users belonging to it.

2. User  
A human user, who belongs to an organization, and can access the services in the organization.

3. Service  
Services belongs to an organisation which is our customer. Each service has a unique ID, name, description, and versions.
4. Version  
A service can have multiple versions, one or more being active at the same time. Each version has a name, and belongs to one service.

## Assumptions
This has been built on the following assumptions - 
1. The service will be read heavy
2. Service for managing users and organizations (ie customers that use the service catalog API) will be built and integrated later. Therefore, we are okay with mocking the user and organization for our implementation
3. Middleware for authorization / authentication will be implemented and integrated later
4. Access control (ie read / write access for certain users) will be built and implemented later
5. We have not set up a scalable RDS, and are okay with a lightweight DB like SQLite for this implementation

## Project Structure
```
.
├── controllers
│   ├── serviceController.go
│   └── versionController.go
├── main.go
├── middleware
│   └── authMiddleware.go
├── repository
│   ├── organization.go
│   ├── repository.go
│   ├── service.go
│   ├── user.go
│   └── version.go
└── resources
    ├── outputFormatter.go
    ├── response.go
    ├── service.go
    └── version.go
```

We have 5 modules, each with particular responsibilities:
1. main  
The module `main` initializes the service, as well as the database, and maps the handlers for each endpoint.
2. middleware  
Responsible for flows such as authentication, and in this case, mocking authentication and populating the customer identity - userID and organisationId of the user into the request context.
3. controllers  
The controllers in `controller` module are responsible for accepting requests, parsing, validating the user input, loading required data using `repository module` and then returning the response to users. This also includes parsing, processing and returning metadata related to pagination.
3. resources  
The elements in `resources` module are responsible for defining IO schema for the API - so that the responses have standardized schema, and the request bodies get parsed and validated.
4. repository  
This is the data storage layer, and has functions to initialise the database and to load data from the database.


## API Reference
| Endpoint               | HTTP Method | Request Body                                                 | Query params and values supported                                                                                                                                                                                                                                                | Description                                                                                                                       |
|------------------------|-------------|--------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------|
| /ping                  | GET         |                                                              |                                                                                                                                                                                                                                                                                  | Returns HTTP 200 OK if application has booted up.                                                                                 |
| /services              | GET         |                                                              | 1. page_size_limit: Integer in range [0-100]. <br>2. page_number: Integer > 0. <br>3. sort_field: ["id", "name","created_at","updated_at", "version_count"]. <br>4. sort_order: ["asc", "desc"]. <br>5. filter_field: ["name", "description"]. <br>6. filter_value: any string.  | Loads all Services in user's organisation.  <br>Supports filtering, sorting and pagination.<br>Default page size supported is 25. |
|                        | POST        | ```{"Name": "srv-name", "Description": "srv-description"}``` |                                                                                                                                                                                                                                                                                  | Creates a Service and returns it                                                                                                  |
| /services/:id          | GET         |                                                              |                                                                                                                                                                                                                                                                                  | Loads and returns a service based on given ID                                                                                     |
| /services/:id/versions | GET         |                                                              | 1. page_size_limit: Integer in range [0-100]. <br>2. page_number: Integer > 0.                                                                                                                                                                                                   | Returns all the versions associated with the given service ID.<br>This endpoint is paginated, and has default page size of 25.    |
|                        | POST        | ```{"Name": "v1.0.0"}```                                     |                                                                                                                                                                                                                                                                                  |                                                                                                                                   |


## Implementation details
### Database models and relationships
<img width="1072" alt="Screenshot 2024-10-16 at 12 43 45" src="https://github.com/user-attachments/assets/9efba4bf-d3d5-4581-8862-5cce075fae8c">

We want to store Services, and their Versions in our database. By extension, we also want to store the organisations these services belong to, and the users that belong to those organisations.
These have relationships - 1:N between services and versions, as well as 1:N between the organisations<>users, organisations<>services, etc as well.  

Besides, since we would be storing business critical systems, we'd want high consistency, transaction support, and support for high read throuput. 

So, we will use Relational Databases.  
Ideally, MySQL or PostGres - but for simplicity, I have chosen SQLite.

We have 4 Tables:  
1. organizations  
2. users  
3. services  
4. versions  

There are foreign key relationships defined to ensure data consistency.

There is a 1:N relationship between services and versions. However, I have made use of denormalisation so version count can be stored on services table, to avoid frequent join operations.

We want to store Services, and their Versions in our database. By extension, we also want to store the organisations these services belong to, and the users that belong to those organisations.
These have relationships - 1:N between services and versions, as well as 1:N between the organisations<>users, organisations<>services, etc as well.  

Besides, since we would be storing business critical systems, we'd want high consistency, transaction support, and support for high read throuput. 

So, we will use Relational Databases.  
Ideally, MySQL or PostGres - but for simplicity, I have chosen SQLite.

We have 4 Tables:  
1. organizations  
2. users  
3. services  
4. versions  

There are foreign key relationships defined to ensure data consistency.

There is a 1:N relationship between services and versions. However, I have made use of denormalisation so version count can be stored on services table, to avoid frequent join operations.

Services and Versions are identified by a Unique ID - generated using ulid package (https://pkg.go.dev/github.com/oklog/ulid/v2) - which is URL safe. We are using a column size of 36, even though ulid is of 26 characters to have a two way door supporting uuids in future.

The entities support soft deletion, by marking the `deleted_at` field.

### Validations
All input users give us, is validated in the controller layer, for example, the query parameters for pagination, sorting, etc.


## How to use
To start the server, run the following command  
```go run main.go```

This sets up the database, inserts relevant mock entries, and starts the service on `http://localhost:8080/`

To verify the service started successfully, we can make a GET request to `http://localhost:8080/ping`, and it should return a HTTP 200 OK response with `pong` in the response body.

## Testing


## Potential improvements in design, and otherwise
We should use a Database Management System - like MySQL or Postgres.

### Error handling
I have tried to hide internal details of implementation bubbling up to users in error messages by returning generic error messages, unless it is a 4xx - in which case the error message communicates to user what is wrong and how they can fix it.  
However, this can be further improved upon by using something like HTTP problem details (https://datatracker.ietf.org/doc/html/rfc7807#section-3) - where additional context about the error can be sent to the user.

### API design and implementation
- The `GET /services` endpoint supports filtering, however, for a user, a "search" operation could be more favorable - to search for services using a part of their name, description etc.
- More routes could be added to support Update and Delete operations on Service and Version resources.
- Rate limiting could be added.  
- We should also have some form of authentication for the API.

#### Pagination
For larger data sets, offset based pagination runs into performance issues with high offset values (https://www.pingcap.com/article/limit-offset-pagination-vs-cursor-pagination-in-mysql/), and cursor based pagination would be preffered - although slightly more complex to implement.
