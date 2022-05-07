# ASYNC ARCHITECTURE

## AUTH SERVICE

### API

POST /auth/register

```
curl -X POST localhost:8080/auth/register -d '{"username":"test","email":"test@example.com","password":"secret"}'
```

POST /auth/login

```
curl -X POST localhost:8080/auth/login -d '{"email":"test@example.com","password":"secret"}'
```

POST /auth/logout

```
curl -X POST localhost:8080/auth/logout -d '{"token":"auth-token"}'
```

POST /auth/check

```
curl -X POST localhost:8080/auth/check -d '{"token":"auth-token"}'
```

PACH /user/{id}

```
curl -X PATCH localhost:8080/auth/user/123 -d '{"role":"manager"}'
```

POST /users/delele

```
curl -X DELETE localhost:8080/auth/user/123
```

GET /users

```
curl localhost:8080/auth/user'
```

### Events

UserCreated
UserDeleted
UserRoleChanged

## CLIENT SERVICE

### API

GET /client/ping

curl -H "Authorization: Bearer $ACCESS_TOKEN" localhost:8081/client/ping
