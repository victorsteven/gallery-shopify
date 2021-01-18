
### Photo Gallery Application

#### Installation
Clone the application using:
```
git clone git@github.com:victorsteven/gallery-shopify.git
```


### The Server
Change to the ``server`` directory:
```
cd server
```
- Update the ``.env`` file with correct details. Both for the real database and the testing database.

- Run unit and integration test from the root of the server directory using:
```
go test -v ./...
```
- In the root of the server directory, start the server using:
```
go run main.go
```
The server will start running on port 7070(the default port from the .env file)

### Testing Endpoints.

#### Login with any of the seeded user email and password and copy the token provided
#### Use this token as `Bearer Token`

#### Delete one image using:
```DELETE http://localhost:7070/images/1```


#### Delete bulk images using:
```DELETE http://localhost:7070/bulk_delete/images```

