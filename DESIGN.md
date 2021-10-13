### Go Fetch Points API

## REQUIREMENTS
 - Serve HTTP Endpoints:
   - Adding Points per Payer
   - Using Points accumulated and returning amount deducted per Payer
   - Checking balance of Points per Payer
 - Points must be used in Timestamp order
 - User may not withdraw more Points than available
 - Payer points may not go negative on spend (No requirement for add?)
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
   - request-body: json array of keyed object with Timestamp, Payer, Points 
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
 - On add, insert each (Timestamp, Payer, Points) tuple into btree ordered by timestamp
   - If negative, add to unused deductions to realized during next spend
   - Add to aggregate balance of payer
 - On usage, 
   1. itterate btree in ascending order to realize unused deductions
     - update transactions to use deductions, removing if points are zero
   2. itterate btree in ascending order tabulating usage from each Payer and Timestamp
     - If total requested amount is met, deduct Points per Timestamp, return deducted Points per Payer, update balances
     - If total requested amount is unmet, return error stating insufficient Points to complete transaction
 - On check, return cached balances

## TESTS
 - Random insertion and inorder traversal of Point tuples 
 - Deduction algorithm
 - Balance algorithm
 - Integration tests for each endpoint
