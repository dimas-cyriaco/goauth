POST http://localhost:4000/users

{
  "name": "Dimas Cyriaco"
}

HTTP 200

[Captures]
user_id: jsonpath "$['id']"

GET http://localhost:4000/users/{{user_id}}

HTTP 200

POST http://localhost:4000/users

{}

HTTP 400
