# Go REST API for User Management

This REST API is built using Go with no external web frameworks and provides endpoints for managing users. It uses the advanced routing features from the 1.22 update for the net/http package and sqlc to interact with a PostgreSQL database.

This project was developed as part of my university coursework on service-oriented architecture and microservice design, initially implemented in Python using Flask, but subsequently rewritten in Go to showcase the versatility and language-agnostic nature of RESTful APIs.

## Endpoints

`[GET] /users`

Retrieves a list of all users.

`[GET] /users/{id}`

Retrieves a user by ID or email.

`[PUT] /users/{id}`

Updates a user by ID.

`[DELETE] /users/{id}`

Deletes a user by ID.

`[POST] /users`

Creates a new user.

## Request and Response Bodies

### Request Body
The request body for PUT /users/{id} and POST /users should be a JSON object with the following properties:

- name: string
- email: string
- phone_number: string
- user_type: string
- address: string
- Response Body

The response body for `[GET] /users` and `[GET] /users/{id}` will be a JSON array or object with the following properties:

- id: int
- name: string
- email: string
- phone_number: string
- user_type: string
- address: string

## Error Handling
If an error occurs, the API will return a JSON error response with a 400 or 500 status code.
The error response will contain a message with the error details.

## Database
The API uses a PostgreSQL database to store user data.
The database connection is established using the DATABASE_URI environment variable.

## References
- [YouTube - TutorialEdge | Building REST APIs in Go 1.22 - New Features](https://www.youtube.com/watch?v=tgLvIghsJFo)
- [YouTube - BugBytes | SQLC in Go - Auto-Generating Database Code in Golang](https://www.youtube.com/watch?v=x_N2VjGQKr4)

## License
This API is licensed under the MIT License. See LICENSE for details.
