# Example Security Check

Attackers will sometimes leave a malicious process running but delete the original binary 
on disk.  This example will use OSQuery to run a query that returns processes in which the 
original binary has been deleted from disk.  The example converts the OSQuery results  into something 
analyzer can consume and then uses analyzer to interpret the results. 

The check is defined in the config.json file included with the example.
```json
{
  "collectors": [
    {
      "id": "process-check",
      "expression": "SELECT name, path, pid FROM processes WHERE on_disk = 0;"
    }
  ],
  "checks": [
    {
      "collector-id": "process-check",
      "Description": "processes running without binary on disk",
      "conditions": [
        {
          "type": "fail",
          "predicate": "@len( $ ) > 2",
          "message": "lots of possible malicious processes exist"
        },
        {
          "type": "warning",
          "predicate": "@len( $ ) > 0 && @len( $ ) <= 2",
          "message": "possible malicious processes exist"
        },
        {
          "type": "pass",
          "predicate": "@len( $ ) == 0",
          "message": "no suspect processes found"
        }
      ]
    }
  ]
}
```
