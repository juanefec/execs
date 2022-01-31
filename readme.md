# execs

Concurrent action execution and result handling.

```go
Executor(action(interface{}) interface{}, handle(interface{})) (fire(interface{}), kill())
```
 - `func action(interface{}) interface{}` is the action executed concurrently on `fire(interface{})`.
 - `func handle(interface{})` takes the result of `action(interface{}) interface{}` and executes concurrently.
 - `func fire(interface{})` pushes the data to execute the action concurrently.
 - `func kill()` blocks until every action has been done and handled, fireing after kill will panic.

 Example:
 ```go
// Run an Executor
// Here you dont care enough to wait for a database save,
// so you just fire saves when needed and a handler takes care of the results
save, endSave := execs.Executor(db.SaveDB, bussines.AfterSaveDB)
defer endSave()

// this function can be called from anywhere
save(db.User{0, "john"})
 ```
 