# Wooble Service
AYY LMAO!

## Installation

```
docker-compose build && docker-compose up -d
```
## Configuration
```
docker exec -it gowooble /bin/bash
cd $CONFPATH

#dev.yml for dev environment
#test.yml for tests

cd $HOME/.aws

#credentials for amazon s3 creds
```

## API Resources

### Users

`POST /v1/users`
```js
Content-Type: application/json

{
  "email": <user email address>
  "name": <username>
  "secret": <password>
  "plan": <selected plan>
  "cardToken"?: <card token created by stripe>
  "isCreator"?: <is the use a creator>
}
```
```js
HTTP/1.1 201 Created
Content-Type: application/json
Location: /tokens

{
  "data": {
    "email": <user email address>
    "name": <username>
    "secret": <password>
    "plan": <selected plan>
    "isCreator"?: <is the use a creator>
  }
}
```

`DELETE /v1/users`
```js
Authorization: <user token>
```
```js
HTTP/1.1 204 NoContent
Authorization: <refreshed token>
```

`POST /v1/funds/bank`
`NOT YET IMPLEMENTED`

`POST /v1/funds/withdraw`
`NOT YET IMPLEMENTED`

### Token

`POST /v1/tokens`
```js
Content-Type: application/json

{
  "email": <user email address>
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

`PUT /v1/token`
```js
Authorization: <user token>
```
```js
HTTP/1.1 201 Created
Authorization: <refreshed token>

{
  "data": {
    "token": <access token>
  }
}
```

### Creations

`GET /v1/creations`
```js
HTTP/1.1 200 OK
Content-Type: application/json

{
  "data": [
    {
      "id": <creation id>
      "title": <creation title>
      "description"?: <creation description>
      "creator": {
        "name": <creator name>
      }
    	"versions": <creation versions>
      "createdAt": <creation date created>
      "updatedAt"?: <creation last updated date
    }
    {
      ...
    }
  ]
}
```

`GET /v1/creations/:creaID`
```js
Content-Type: application/json

{
  "data": {
    "id": <creation id>
    "title": <creation title>
    "description"?: <creation description>
    "creator": {
      "name": <creator name>
    }
    "versions": <creation versions>
    "createdAt": <creation date created>
    "updatedAt"?: <creation last updated date
  }
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
Location: /creations/:<id of the new creation>
Authorization: <refreshed token if expired>

{
  "data": {
    "id": <creation id>
    "title": <creation title>
    "description"?: <creation description>
    "creator": {
      "name": <creator name>
    }
    "versions": <creation versions>
    "createdAt": <creation date created>
    "updatedAt"?: <creation last updated date
  }
}
```

`PUT /v1/creations/:encid`
`It only updates the description and price`
```js
Content-Type: application/json
Authorization: <user token>

{
  "engine": <creation engine>
  "title": <creation title>
  "description"?: <creation description>
  "price"?: <creation price>
}
```
```js
HTTP/1.1 204 NoContent
Location: /creations/:<id of the creation>
Authorization: <refreshed token if expired>
```

`DELETE /v1/creations/:encid`
`Delete the creation if nobody use it, else it will make the creation unlisted by putting its state to 'delete'`
```js
Authorization: <user token>
```
```js
HTTP/1.1 204 NoContent
Authorization: <refreshed token if expired>
```

`PATCH /v1/creations/:encid/publish`
`Set creation state to 'public'`
```js
Authorization: <user token>
```
```js
HTTP/1.1 204 NoContent
Authorization: <refreshed token if expired>
```

`GET /v1/creations/:encid/code`
```js
Authorization: <user token>
```
```js
HTTP/1.1 200 OK
Content-Type: application/json
Authorization: <refreshed token if expired>

{
  "data": {
    "script": <script code>
    "document"?: <document code>
    "style"?: <style code>
  }
}
```

`POST /v1/creations/:encid/versions`
```js
Content-Type: application/json
Authorization: <user token>

{
  "version": <version to create>
}
```
```js
HTTP/1.1 204 NoContent
Authorization: <refreshed token if expired>
Location: /creations/<creation id>/code
```

`PUT /v1/creations/:encid/versions`
`Save the last version (only if in draft)`
```js
Content-Type: application/json
Authorization: <user token>

{
  "script": <script code>
  "document"?: <document code>
  "style"?: <style code>
}
```
```js
HTTP/1.1 204 NoContent
Authorization: <refreshed token if expired>
Location: /creations/<creation id>/versions
```

### Packages

`GET /v1/packages`
```js
Authorization: <user token>
```
```js
HTTP/1.1 200 OK
Content-Type: application/json
Authorization: <refreshed token if expired>

{
  "data": [
    {
      "id": <package id>
      "title": <package title>
      "domains":[<domains associated to the package>]
      "createdAt": <package creation date>
      "updatedAt"?: <package last update date>
      "creations"?: [...]
    },
    {
      ...
    }
  ]
}
```

`GET /v1/packages/:pkgID`
```js
Content-Type: application/json
Authorization: <user token>
```
```js
HTTP/1.1 200 OK
Content-Type: application/json
Authorization: <refreshed token if expired>

{
  "data": {
    "id": <package id>
    "title": <package title>
    "domains":[<domains associated to the package>]
    "createdAt": <package creation date>
    "updatedAt"?: <package last update date>
    "creations"?: [...]
  }
}
```

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
Content-Type: application/json
Location: /packages/:<id of the new package>
Authorization: <refreshed token if expired>
```

`PUT /v1/packages/:encid`
```js
Content-Type: application/json
Authorization: <user token>

{
  "title": <package title>
  "domains"?: [<domains with which the package will work>]
}
```
```js
HTTP/1.1 204 NoContent
Location: /packages/:<id of the new package>
Authorization: <refreshed token if expired>
```

`DELETE /v1/packages/:encid`
```js
Authorization: <user token>
```
```js
HTTP/1.1 204 NoContent
Authorization: <refreshed token if expired>
```

`POST /v1/packages/:encid/creations`
```js
Content-Type: application/json
Authorization: <user token>

{
  "creation": <IDs of creations to push in package>
}
```
```js
HTTP/1.1 204 NoContent
Location: /v1/packages/:<id of the package>
Authorization: <refreshed token if expired>
```

`DELETE /v1/packages/:encid/creations`
```js
Authorization: <user token>
```
```js
HTTP/1.1 204 NoContent
Authorization: <refreshed token if expired>
```

`PUT /v1/packages/:packageID/build`
```js
Authorization: <user token>
```
```js
HTTP/1.1 200 OK
Content-Type: application/json
Authorization: <refreshed token if expired>

{
  "data": {
    "source": <CDN url of the builded package>
  }
}
```