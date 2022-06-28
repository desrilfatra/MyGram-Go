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
                    "email": "andrikuwito2@gmail.com",
                    "password": "$2a$10$cvpW1zR8RXkG5VBosoBJ/./kXKaO7pKXmzaLfUgsE6rU61TxqEJvi",
                    "username": "desril",
                    "date": "2022-06-27T13:06:37.558+07:00"
                }
            }
            ```
