# MyGram-Go
Final Project Scalable Web Service with Golang Hacktiv8


* #### Endpoint List : 
    * ##### User : 
        * #### Register
            
            [POST]```http://localhost:8080/users/register```
            
            body :

            ```json
            {
                "age": 23,
                "email":"desril@gmail.com",
                "password":"desril",
                "username":"desril"
            }
            ```

            response
            ```json
            {
                "data": {
                    "id": 1,
                    "age": 21,
                    "email": "desril@gmail.com",
                    "password": "$2a$10$cvpW1zR8RXkG5VBosoBJ/./kXKaO7pKXmzaLfUgsE6rU61TxqEJvi",
                    "username": "desril",
                    "date": "2022-06-27T13:06:37.558+07:00"
                }
            }
            ```
            
        * #### Login
            [POST]```http://localhost:8080/users/login```
            
            body :

            ```json
            {
                "email": "desril@mailz.com",
                "password": "desril"
            }
            ```

            response
            ```json
            {
                "data": {
                   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImRlc3JpbEBtYWlsei5jb20iLCJpZCI6MX0._CTowajg3nwUKU4qn4qJzP02pOmQBzzYcaZTEinmQs8"
                }
            }
            ```
            
       * #### Update User
            [PUT]```http://localhost:8080/users/1```
            
            Authorization :

            ```
            {
                Authorization: "Bearer {{token}}"
            }
            ```

            contoh :
            ```
            {
                Authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZF91c2VyIjoyNH0.XpOW4v9hpneBw9gnsVAGli_zDqj7VmMLW6ZHL80MauQ"
            }
            ```

            body :

            ```json
            {
                "email" : "desril12@mailz.com",
                "username":"desril12"
            }  
            ```

            response
            ```json
            {
                "status": "ok",
                "data": {
                          "id": 1,
                          "email": "desril12@mailz.com",
                          "username": "desril12",
                          "age": 22,
                          "updated_at": "2022-06-29T22:01:39.6471884+07:00"
                        }
            }
            ```
        
             * #### Delete User
            [DELETE]```http://localhost:8080/users```
            
            Authorization :

            ```
            {
                Authorization: "Bearer {{token}}"
            }
            ```

            token didapatkan ketika melalui proses login

            contoh :
            ```
            {
                Authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZF91c2VyIjoyNH0.XpOW4v9hpneBw9gnsVAGli_zDqj7VmMLW6ZHL80MauQ"
            }
            ```            

            response
            ```json
            {
                  {
                          "message": "Your account has been successfully deleted"
                  }
            }
            ```
