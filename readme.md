# execs

This is a go package that offers a way to do things and handle them in a concurrent way.


```go
Executor(do(interface{}) interface{}, handle(interface{})) (fire(interface{}), kill())
```
 - `func do(interface{}) interface{}` executed concurrently on `fire(interface{})`
 - `func handle(interface{})` the result `do(interface{}) interface{}`
 - `func fire(interface{})` a message, argument of `do(interface{}) interface{}`
 - `func kill()` blocks until every message has been fired and handled

 Example:
 ```go
 // Create a Sub

// Handle a new execution type
fireSaveData, killSaveData := execs.Executor(db.SaveDB, analytic.RecordSaveDB)
defer killSaveData()

// this function can be called from anywhere
fireSaveData(db.User{0, "john"})
 ```