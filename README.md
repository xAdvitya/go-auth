# go-auth
A simple authentication service written in Go.

## Features

- User registration
- User login
- Password hashing
- JWT token generation and validation

## Installation

```sh
git clone https://github.com/xAdvitya/go-auth.git
cd go-auth
go mod tidy
```

## Usage

1. Start the server:

    ```sh
    go run main.go
    ```

2. The server will be running at `http://localhost:8080`.

## API Endpoints

- `POST /register` - Register a new user
- `POST /login` - Login a user and get a JWT token

## Configuration

You can configure the application using environment variables:

- `PORT` - The port on which the server will run (default: 8080)
- `JWT_SECRET` - The secret key used for signing JWT tokens
- `MONGO_URI` - The URI of your MongoDB instance
