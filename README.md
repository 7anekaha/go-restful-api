# TotalCoder - Assignment 17 - Rest API with MongoDB and InMemoryDB
The assignment is to create a REST API with three endpoints using MongoDNB and InMemoryDB.
Endoints:
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

# Steps to run the project
1. Clone the repository
2. Compile and Run the project
```bash
make run
```
