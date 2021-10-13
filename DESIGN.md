### Go Fetch Points API

## REQUIREMENTS
 - Serve HTTP Endpoints:
   - Adding Points per Payer
   - Using Points accumulated and returning amount deducted per Payer
   - Checking balance of Points per Payer
 - Points must be used in Timestamp order
 - User may not withdraw more Points than available
 - Points may be added out of Timestamp order
 - Handle errors in a meaningful way

## ASSUMPTIONS
 - API is for a single User
 - Added Points will never occur simultaneously
 - Payers are case sensitive 

## ENDPOINTS
 - /add-points
   - method: POST
   - content-type: application/json
   - request-body: json keyed object with Timestamp, Payer, Points 
   - response: json keyed object with success boolean and possible error message
 - /use-points
   - method: POST
   - content-type: application/json
   - request-body: json keyed object with amount of Points to use
   - response: json keyed object with success booelan, error message on failure or data array of keyed objects containing Payer and deducted Points on success
 - /balance
   - method: GET
   - content-type: empty
   - request-body: empty
   - response: json keyed object with success boolean, error message on failure or data array of keyed objects containing Payer and cummulative Points on success

## WORKFLOWS
 - On add, insert (Timestamp, Payer, Points) tuple into btree ordered by timestamp
 - On usage, itterate btree in ascending order tabulating usage from each Payer and Timestamp
   - If total requested amount is met, deduct Points per Timestamp, return deducted Points per Payer
   - If total requested amount is unmet, return error stating insufficient Points to complete transaction
 - On check, iterate btree and sum all points per payer

## TESTS
 - Random insertion and inorder traversal of Point tuples 
 - Deduction algorithm
 - Balance algorithm
 - Integration tests for each endpoint
