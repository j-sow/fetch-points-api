# Points API
See DESIGN.md for assumptions and design plan

## Prerequisites
 - gnu make or equivalent: use either system package manager or Homebrew, if on windows see http://gnuwin32.sourceforge.net/packages/make.htm
 - go (tested with >1.14): See instructions at https://golang.org/doc/install
 - docker: See instructions at https://docs.docker.com/get-docker/

## Building
To build standalone with go
```
make build
```

## Running
The program can either be run with
```
make run
```

Or by simply executing
```
./RunRewardsAPI
```

Default port is 8080, but can be changed with `-host` flag
```
./RunRewardsAPI -port 8888
```

## Tests
To run tests
```
make test
```

## Dockerized
To run with docker (no build necessary)
```
make run-docker
```

To run tests with docker
```
make test-docker
```

## API
### **Add Points**
----
Add transactions to points store

* **URL**

  /add-points

* **Method:**
  
  `POST`
  
* **Data Params**

  _JSON array_ 

  **REQUIRED PER OBJECT** <br />
  `timestamp: [string]` - time of transaction <br />
  `points: [int64]` - points added or subtracted <br />
  `payer: [string]` - with whom transaction is being made

* **Success Response:**
  
  _On success you should see a 200 HTTP Status and json object_

  * **Code:** 200 <br />
    **Content:** _JSON encoded body_ <br />
    `success: [bool]` - if addition was succesful 

* **Example:**
```
curl -X POST http://localhost:8080/add-points \
   -H 'Content-Type: application/json' \
   -d '[{"timestamp": "2021-10-10T00:00:00Z", "points": 200, "payer": "GENERAL MILLS"}]'
```

### **Check Balance**
----
  Request balance for all payers in transaction history

* **URL**

  /check-balance

* **Method:**
  
  `GET`
  
* **Data Params**
  
  _None_

* **Success Response:**
  
  _On success you should see a 200 HTTP Status and json encoded body_

  * **Code:** 200 <br />
    **Content:** _JSON encoded body_ <br/>
    `success: [bool]` - if checking was successful <br/>
    `data: [Object]` - json object with payers as keys and points as values

* **Example:**
```
curl http://localhost:8080/check-balance
```

### **Use Points**
----
  Use points from points store

* **URL**

  /use-points

* **Method:**
  
  `POST`
  
* **Data Params**

  _JSON object_ 

  **REQUIRED** <br />
  `points: [int64]` - points added or subtracted <br />

* **Success Response:**
  
  _On success you should see a 200 HTTP Status and json encoded body_

  * **Code:** 200 <br />
    **Content:** _JSON encoded body_ <br />
    `success: [bool]` - if using points was successful <br />
    `data: [Array]` - array of objects containing payers and amounts deducted

* **Example:**
```
curl -X POST http://localhost:8080/use-points
   -H 'Content-Type: application/json'
   -d '{"points": 200}'

```


## Prompt
Our users have points in their accounts. Users only see a single balance in their accounts. But for reporting purposes we actually track their
points per payer/partner. In our system, each transaction record contains: payer (string), points (integer), timestamp (date).

For earning points it is easy to assign a payer, we know which actions earned the points. And thus which partner should be paying for the points.

When a user spends points, they don't know or care which payer the points come from. But, our accounting team does care how the points are
spent. There are two rules for determining what points to "spend" first:
 - We want the oldest points to be spent first (oldest based on transaction timestamp, not the order they’re received)
 - We want no payer's points to go negative.

Provide routes that:
- Add transactions for a specific payer and date.
- Spend points using the rules above and return a list of { "payer": <string>, "points": <integer> } for each call.
- Return all payer point balances.

Note:
- We are not defining specific requests/responses. Defining these is part of the exercise.
- We don’t expect you to use any durable data store. Storing transactions in memory is acceptable for the exercise.

### Example
Suppose you call your add transaction route with the following sequence of calls:
```
{ "payer": "DANNON", "points": 1000, "timestamp": "2020-11-02T14:00:00Z" }
{ "payer": "UNILEVER", "points": 200, "timestamp": "2020-10-31T11:00:00Z" }
{ "payer": "DANNON", "points": -200, "timestamp": "2020-10-31T15:00:00Z" }
{ "payer": "MILLER COORS", "points": 10000, "timestamp": "2020-11-01T14:00:00Z" }
{ "payer": "DANNON", "points": 300, "timestamp": "2020-10-31T10:00:00Z" }
```

Then you call your spend points route with the following request:
```
{ "points": 5000 }
```

The expected response from the spend call would be:
```
[
    { "payer": "DANNON", "points": -100 },
    { "payer": "UNILEVER", "points": -200 },
    { "payer": "MILLER COORS", "points": -4,700 }
]
```

A subsequent call to the points balance route, after the spend, should returns the following results:
```
{
    "DANNON": 1000,
    "UNILEVER": 0,
    "MILLER COORS": 5300
}
