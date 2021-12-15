# Simple api TODO for studenty

## Routes HTTP

#### Create new todo
```bash
curl --location --request POST 'http://localhost:8080/todos' \
--header 'Content-Type: application/json' \
--data-raw '{
	"title": "title",
	"description": "description",
	"statusId": 1
}'
```


#### Get all todos 
```bash
curl --location --request GET 'http://localhost:8080/todos' 
```

#### Create Status Todo
```bash
curl --location --request POST 'http://localhost:8080/todos/status' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "doing"
}'
```

#### Get Status Todo
```bash
curl --location --request GET 'http://localhost:8080/todos/status/6'
```

#### Update Status Todo
```bash
curl --location --request PUT 'http://localhost:8080/todos/9' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "bar other",
    "description": "foo other",
    "statusId": 2
}'
```