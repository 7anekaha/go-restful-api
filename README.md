# TotalCoder - Assignment 17 - Rest API with MongoDB and InMemoryDB
The assignment is to create a REST API with three endpoints using MongoDB and InMemoryDB.
(In the assigment a mongo uri was provided but I used a docker container to run MongoDB instead.)

Doing this assignment I learned:
- How TestMain works in Go
- How to use [testcontainers](https://golang.testcontainers.org/) to test MongoDB
- How to overrride UnmarshalJSON to add a custom validation

## Endoints:
1. POST /mongodb - To fetch data from MongoDB
Request Body:
```json
{
    "startDate": "YYYY-MM-DD"
    "endDate": "YYYY-MM-DD"
    "minCount": 2800
    "maxCount": 2900
}
```
Response:
```json
{
    "code": 0,
    "msg": "Success",
    "records": [
        {
            "key": "TAKwGc6Jr4i8Z487",
            "createdAt": "2017-01-28T01:22:14.398Z",
            "totalCount": 2800
        },
        {
            "key": "NAeQ8eX7e5TEg7oH",
            "createdAt": "2017-01-27T08:19:14.135Z",
            "totalCount": 2900
        }
    ]
}
```
Code 0 means success and 1 means bad request.

The message will be "Success" for code 0 and "Structure of the request body is not correct" or "Start date should be before end date and minCount should be less than maxCount" for code 1.

2. GET /inmemory - To fetch data from InMemoryDB
Request: "http://localhost:3000/inmemory/?key=active-tabs"
Response:
```json
{
   "key": "active-tabs",
   "value": "getir"
}
```
3. POST /inmemory - To store data from InMemoryDB
Request Body:
```json
{
    "key": "active-tabs",
    "value": "getir"
}
```
Response:
```json
{
    "key": "active-tabs",
    "value": "getir"
}
```

## Requirements
- Go 1.22
- Docker
- Docker Compose

## Steps to run the project
1. Clone the repository
2. Compile and Run the project
```bash
make run
```

## Test the project
```bash
make test
```
