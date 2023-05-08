# Database Setup
```
docker run -d \
-e POSTGRES_DB=bookstore \
-e POSTGRES_USER=irfan \
-e POSTGRES_PASSWORD=123456 \
-p 5432:5432 postgres
```

# Start the development server
```
go run src/main.go
```