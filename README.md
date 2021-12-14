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
curl --location --request GET 'http://localhost:8080/todos' \
--data-raw ''
```