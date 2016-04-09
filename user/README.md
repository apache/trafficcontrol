* `POST /` 

  Create a new user
  
  ```
  curl -Lkvs -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' 'http://localhost:8080/' -d'{"username":"jvd", "firstName":"vandoorn", "password": "secret"}'
  ```
  
* `GET /{id}`

  Read user with id
* `PUT /{id}` 

  Update user with id
* `DELETE /{id}`

  Delete user with id
