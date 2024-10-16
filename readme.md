# Service Catalog API
This repository contains an implementation of an API for storing and retrieving services and version - a catalog of services.  
It makes use of Gin HTTP framework (https://gin-gonic.com/docs/), and GORM (https://gorm.io/docs/) as the ORM, and makes use of SQLite (https://pkg.go.dev/github.com/mattn/go-sqlite3) for storage.  

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
![alt text](<Screenshot 2024-10-16 at 10.02.06.png>)

### Validations



## How to use
To start the server, run the following command  
```go run main.go```

This sets up the database, inserts relevant mock entries, and starts the service on `http://localhost:8080/`

To verify the service started successfully, we can make a GET request to `http://localhost:8080/ping`, and it should return a HTTP 200 OK response with `pong` in the response body.

## Testing


## Potential improvements in design, and otherwise
