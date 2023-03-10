# Simple RESTful API with GO
Building a simple RESTful API in the GO programming language

##### Example Input: 
```
{
    "title": "My Kazakhstan",
    "author": "Ali Khodzhaev",
    "rating": 4
}
```

##### Example Response: 
```
{
    "status": 200,
    "message": Create data successfully 
} 
```

### POST /

Creates new book


### GET /books

Returns all books

##### Example Response: 
```
[
	{
            "id": "1",
            "title": "Mother",
            "author": "Maksim Gorki",
            "publish_date" "2022-01-10",
            "rating": 5
        }        
]
```

### DELETE /books/1

Deletes books by ID:

### UPDATE /books/1

Update books by ID:

##### Example Body: 
```
{
	"title": "My Kazakhstan",
        "author": "Ali Ahmedov",
        "rating": 3            
}
```

## Run Project

Use ```go run main.go``` to run and install all dependencies in import
