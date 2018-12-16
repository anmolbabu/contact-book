# contact-book

## Overview

This is an application to store contact details identified by unique email id
The contact CRUD operations are exposed via Rest-end points authenticated by a basic auth.
This application is also ready with required configuration for deployment on heroku

## Installation

### Prerequisites
Install go. Ref: `https://golang.org/doc/install`

### Standalone
1. Create dir `$GOPATH/src/github.com/anmolbabu`
2. Clone the repo into the above dir `git clone https://github.com/anmolbabu/contact-book.git`
3. Build the binary `make`
4. Run the application `make run`

### Heroku
1. Create dir `$GOPATH/src/github.com/anmolbabu`
2. Clone the repo into the above dir `git clone https://github.com/anmolbabu/contact-book.git`
3. Create a Heroku account at `https://id.heroku.com/login`
4. Install heroku cli. Ref: `https://devcenter.heroku.com/articles/heroku-cli`
5. Create heroku app `heroku create`
6. Push the source to the app `git push heroku master`
7. Optionally, watch logs using `heroku logs --tail`

## Features

* Create a contact

```
ex:
URL: https://mysterious-springs-19562.herokuapp.com/contacts
Req Type: POST
Headers:
    Content-Type: application/json
    Accept: application/json
Body:
    {
        "name": "def",
        "emailid": "def@email.com"
    }
```

* Get Contacts

```
ex:
URL: https://mysterious-springs-19562.herokuapp.com/contacts
Req Type: GET
```

* Get contacts filtered by name and/or email id

```
ex:
URL: https://mysterious-springs-19562.herokuapp.com/contacts?name=def&emailid=def@email.com
Req Type: GET
```

* Delete contact

```
ex:
URL: https://mysterious-springs-19562.herokuapp.com/contacts/def@email.com
Req Type: DELETE
```

* Update contact
```
ex:
URL: https://mysterious-springs-19562.herokuapp.com/contacts/def@email.com
Req Type: PUT
Headers:
    Content-Type: application/json
    Accept: application/json
Body:
    {
        "name": "DEF"
    }
```
* Pagination
```
URL: http://127.0.0.1:8080/contacts?page=2&pagelimit=3&name=jkl
```

Note:

Basic authentication is in use and uses ENV names `PLIVO_AUTH_USER` and `PLIVO_AUTH_PASS` for user id and password respectively

```
Defaults:
    User name: `heroku`
    Password: `plivo`
```

Sample request with auth header:
curl --user heroku:plivo https://vast-scrubland-53840.herokuapp.com/contacts

## Running Tests

To run tests, use `make test`