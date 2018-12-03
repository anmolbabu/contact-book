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
        "Name": "def",
        "EmailID": "def@email.com"
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
URL: https://mysterious-springs-19562.herokuapp.com/contacts?Name=def&EmailID=def@email.com
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
        "Name": "DEF"
    }
```