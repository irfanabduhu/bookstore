# Database Setup
```
docker run -d \
-e POSTGRES_DB=bookstore \
-e POSTGRES_USER=irfan \
-e POSTGRES_PASSWORD=123456 \
-p 5432:5432 postgres
```

Run the [queries](/src/config/model.sql) to setup the database.

# Start the development server
```
cd src
go mod tidy
go run main.go
```

Use the [Client](/src/client.http) to interact with the server.
