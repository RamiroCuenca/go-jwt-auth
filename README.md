
# Authentication API - Implementation of JWT on Go

In this project i made an app that allow the user to:
- Register / SignUp.
- LogIn. 

And when the user authenticates with it credentials it receives a JWT (Json Web Token) with wich he will be able to access other endpoints such as:
- ReadAll: Fetch all users data.
- ReadById: Fetch the data from a specific user.
- UpdateById: Update "username" from a specific user.
- DeleteById: Delete a specific user.
Also, in this project i applied testing with native "Testing" library from Go, security protocols in order to store safely the password from the user and finally applied the necessary procedures to use JWT.

## API Reference

#### SignUp / Register a new user

Returns a fresh Json Web Token through the headers under the key "Token".

```http
  POST /api/v1/register
```

| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `username` | `string` | **Required** - *Unique* - Between 1 and 50 digits|
| `email` | `string` | **Required** - *Unique* - Should be valid|
| `password` | `string` |  **Required** - At least 6 digits|

---

#### Login with an existing user

Returns a fresh Json Web Token through the headers under the key "Token".

```http
  POST /api/v1/login
```

| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `email` | `string` | **Required** |
| `password` | `string` |  **Required** |

---

#### Fetch all users

Returns a json with data from all registered users.

```http
  GET /api/v1/readall
```

| Header Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `Authorization` | `string` | **Required** - JWT token - Should still be active|

---

#### Fetch a specific user

Returns a json with data from the fetched user.

```http
  GET /api/v1/readbyid
```

| URL Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required** |

| Header Parameter| Type     | Description                |
| :-------- | :------- | :------------------------- |
| `Authorization` | `string` | **Required** - JWT token - Should still be active|

---

#### Update a specific user

Returns a json with data from the updated user.

```http
  PUT /api/v1/updatebyid
```

| URL Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required** |

| Body Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `username` | `string` | **Required** - Must be unique|

| Header Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `Authorization` | `string` | **Required** - JWT token - Should still be active|

---

#### Delete a specific user

Returns a json with data from the deleted user.

```http
  DELETE /api/v1/deletebyid
```

| URL Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `int` | **Required** |

| Header Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `Authorization` | `string` | **Required** - JWT token - Should still be active|

## Database Reference

The database that this project use is a PostgresDB and consists in just 1 table called "users".

| Column | Data Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `SERIAL` | *PRIMARY KEY* |
| `username` | `VARCHAR(50)` | **NOT NULL** - *Unique* |
| `email` | `VARCHAR(80)` | **NOT NULL** - *Unique* |
| `hashed_password` | `VARCHAR(255)` |  **NOT NULL** |
| `created_at` | `TIMESTAMP` | **NOT NULL** - *DEFAULT NOW()* |
| `updated_at` | `TIMESTAMP` |  |


## Author

- [@ramirocuencasalinas](https://www.linkedin.com/in/ramiro-cuenca-salinas/)

  
