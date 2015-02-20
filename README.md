# DashDB
The fast, simple, scalable, and safe key-value DB in Go. 
This database is designed to be used with [conductor](https://github.com/Vluxe/conductor) for scalability.

## Examples

```go
  d, err := dash.New()
  if err != nil {
    fmt.Println(err)
    return
  }
  defer d.Cleanup()

  loaded := d.Get("name")
  fmt.Println("disk value is:", loaded)

  d.Set("name", "Dalton")
  val := d.Get("name")
  fmt.Println("value is:", val)

  d.Remove("name")
  fmt.Println("remove value")

  d.Set("name", "Austin")
  v := d.Get("name")
  fmt.Println("value is:", v)

```

Also show conductor integration!
