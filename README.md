# Wooble Service
AYY LMAO!

## Installation

```
docker-compose build && docker-compose up -d
```
## Configuration
```
docker exec -it woobleservice_gowooble_1 /bin/bash
cd $CONFPATH

#dev.yml for dev environment
#test.yml for tests

cd $HOME/.aws

#credentials for amazon s3 creds
```

## API Resources

### Token

`POST /v1/token/generate`
```js
Content-Type: application/json

{
  "login": <username or user email address>
  "secret": <password>
}
```
```js
HTTP/1.1 201 Created
Content-Type: application/json

{
  "data": {
    "token": <access token>
  }
}
```

`POST /v1/token/refresh`
```js
Authorization: <user token>

{}
```
```js
HTTP/1.1 201 Created
Authorization: <refreshed token>

{}
```

### Creations

`GET /v1/creations/:id`
```js
Content-Type: application/json

{
  "data": {
    "id": <creation id>
    "title": <creation title>
    "creator": {
      "name": <creator name>
    }
    "version": <creation version>
    "createdAt": <creation date created>
    "updatedAt": <creation last updated date
  }
}
```

`GET /v1/creations`
```js
Content-Type: application/json

{
  "data": [
    {
      "id": <creation id>
      "title": <creation title>
      "creator": {
        "name": <creator name>
      }
    	"version": <creation version>
      "createdAt": <creation date created>
      "updatedAt": <creation last updated date
    },
    {
      ...
    }
  ]
}
```

`POST /v1/creations`
```js
Content-Type: application/json
Authorization: <user token>

{
  "engine"?: <engine, JSES5 default>
  "title": <creation title>
  "document"?: <creation HTML>
  "script"?: <creation Script>
  "style"?: <creation CSS>
}
```
```js
HTTP/1.1 201 Created
Location: /v1/creations/:<id of the new creation>
Authorization: <refreshed token if expired>

{}
```

### Packages

`POST /v1/packages`
```js
Content-Type: application/json
Authorization: <user token>

{
  "title": <package title>
  "domains"?: [<domains with which the package will work>]
}
```
```js
HTTP/1.1 201 Created
Location: /v1/users/:<name of package owner>/packages/:<id of the new package>
Authorization: <refreshed token if expired>

{}
```

`POST /v1/users/:username/packages/:packageID/push`
```js
Content-Type: application/json
Authorization: <user token>

{
  "creations": [<IDs of creations to push in package>]
}
```
```js
HTTP/1.1 201 Created
Location: /v1/users/:<name of package owner>/packages/:<id of the package>
Authorization: <refreshed token if expired>

{}
```

`GET /v1/users/:username/packages/:packageID/build`
```js
Content-Type: application/json
Authorization: <user token>
```
```js
HTTP/1.1 200 OK
Authorization: <refreshed token if expired>

{
  "data": {
    "source": <CDN url of the builded package>
  }
}
```