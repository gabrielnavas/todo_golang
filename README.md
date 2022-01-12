# Simple api TODO for studenty

# Routes HTTP
<br>

# Routes User

#### Create new user
```bash
curl --location --request POST 'http://localhost:8080/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "navas",
	"username": "navas",
	"email": "navas@email.com",
    "password": "123456",
	"passwordConfirmation": "123456"
}'
```

#### Update user
```bash
curl --location --request PUT 'http://localhost:8080/users/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "navas10",
	"username": "navas10",
	"email": "navas10@email.com",
    "password": "654321",
	"passwordConfirmation": "654321",
    "levelAccess": 2
}'
```

#### Get All users
```bash
curl --location --request GET 'http://localhost:8080/users'
```

#### Get user
```bash
curl --location --request GET 'http://localhost:8080/users/1'
```

#### Delete user
```bash
curl --location --request DELETE 'http://localhost:8080/users/1'
```

#### Change Password user
```bash
curl --location --request POST 'http://localhost:8080/users/change_password/3' \
--header 'Content-Type: application/json' \
--data-raw '{
    "oldPassword": "123456",
    "newPassword": "112233",
    "newPasswordConfirmation": "112233"
}'
```

#### Patch Photo user
```bash
curl --location --request PATCH 'http://localhost:8080/users/photo/3' \
--form 'photo=@"/home/navas/Desktop/my_photo.jpg"'
```

#### Delete Photo user
```bash
curl --location --request DELETE 'http://localhost:8080/users/photo/3'
```

#### Get Photo user
```bash
curl --location --request GET 'http://localhost:8080/users/photo/2'
```

# Routes Login
```bash
curl --location --request POST 'http://localhost:8080/users/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "navas",
    "password": "123456"
}'
```

# Routes Todo
<br>

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

#### Get todo
```bash
curl --location --request GET 'http://localhost:8080/todos/2'
```

#### Get all todos 
```bash
curl --location --request GET 'http://localhost:8080/todos' 
```

#### Delete todo
```bash
curl --location --request DELETE 'http://localhost:8080/todos/2'
```

#### Update Image Todo
```bash
curl --location --request PATCH 'http://localhost:8080/todos/image/1' \
--form 'image=@"photo.jpg"'
```

#### Get Image Todo
```bash
curl --location --request GET 'http://localhost:8080/todos/image/1'
```

#### Delete Image Todo
```bash
curl --location --request DELETE 'http://localhost:8080/todos/image/1'
```

#### Create Status Todo
```bash
curl --location --request POST 'http://localhost:8080/todos/status' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "doing"
}'
```

#### Update Status Todo
```bash
curl --location --request PUT 'http://localhost:8080/todos/status/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "finish"
}'
```

#### Get All Status Todo
```bash
curl --location --request GET 'http://localhost:8080/todos/status'
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

#### Delete Status Todo
```bash
curl --location --request DELETE 'http://localhost:8080/todos/status/3'
```