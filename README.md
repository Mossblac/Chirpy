# Chirpy

###

> a boot.dev project

###

![Welcome to Chirpy!](/github.com/Mossblac/Chirpy/assets/logo.png)


###

## Project Goals

###

- Understand web servers
- Build a production style HTTP server in Go, without using framework
- Use JSON headers and status codes to communicate with clients via RESTful API
- **Learn** what makes Go great for building servers
- Use type safe SQL to store and retrieve from a Postgres database
- Implement a secure authentication/ authorization system with well-tested cryptography libraries
- Build and understand webhooks and API keys
- Document the REST API with markdown  


###

## What Chirpy Does

> thankfully it is not called X-er

this is a mock version of a popular social media app.   
It allows users to create, store, and access messages stored within a local server.

I had a lot of fun with this project.

There is no front end other than a single image (above), a landing page, and a few GET commands.  
It is mostly command line driven,but great for learning and testing these (previously new to me) concepts.

###

## To Install 

in a Go module 

> go get https://github.com/Mossblac/Chirpy

###

## Setup Instructions

Follow these steps to get the project running locally:

1. **Install Dependencies**

   Make sure you have [Go](https://golang.org/dl/) installed. You will also need [Goose](https://github.com/pressly/goose) for managing database migrations.
   
   To install Goose:
   ```sh
   go install github.com/pressly/goose/v3/cmd/goose@latest

2. **Create The Database**

> touch ./your_database.db

3. **Run Database Migrations**

> goose -dir migrations sqlite3 ./your_database.db up

## How To Use 

> test.http 

###

**test.http** includes descriptions and instructions for all the supported commands of the application  
it is formatted for testing with the REST client within VScode.   




