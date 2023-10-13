# Mackenzie - Generic Function Caching Library in Go

Mackenzie is a lightweight, highly flexible caching library for Go applications, leveraging the power of generics and reflection to offer enhanced caching capabilities.

## Features

- Uses generics for flexibility across different data types.
- Cache item expiration based on configurable lifetime.
- Automatic cleanup of expired items at regular intervals.
- Strong type checking ensures that caching methods match the expected signatures.
- Force-get to bypass cache and fetch fresh data.
- Provides time until a cached item expires.
- Efficient locking for thread safety.

## Installation

```bash
go get github.com/nokusukun/mackenzie
```

## Usage

### 1. Creating a Cache

To create a cache instance, use the `Create` method:

```go
cache, err := mackenzie.Create[YourType](yourFunction, mackenzie.Config{Lifetime: time.Minute, CleanInterval: time.Minute * 10})
if err != nil {
    log.Fatal(err)
}
```

The passed function should return the data type you want to cache and may optionally return an error as its second return value.

### 2. Get Data

To get data from the cache:

```go
item, err := cache.Get(arguments...)
if err != nil {
    log.Fatal(err)
}
```

If the data is in the cache and hasn't expired, it will be returned. Otherwise, the cache will call the provided function to fetch and cache the data.

### 3. Force Get

To bypass the cache and fetch fresh data:

```go
item, err := cache.ForceGet(arguments...)
```

### 4. Clear Cache

Clear the entire cache:

```go
cache.Clear()
```

Clear a specific key:

```go
cache.ClearKey(arguments...)
```

### 5. Check Item Expiry

To check how long until a cached item expires:

```go
duration := cache.ExpiresIn(arguments...)
```

### 6. Cleanup

Before your application exits or if you're done with the cache instance:

```go
cache.Unload()
```

## Errors

You can check for specific Mackenzie-related errors using the helper functions:

- `IsErrCallInvalid(err)`: Returns true if the error is due to an invalid call.
- `IsErrMackenzie(err)`: Returns true for generic Mackenzie errors.

# Example

### Basic Square Function
This example illustrates the basic operation of a function that takes in a `Point` structure and returns the product of its `X` and `Y` fields.

```go
p := &Point{X: 5, Y: 5}
result := Square(p)
fmt.Println(result)  // Outputs: 25
```

### Creating a Cache
This example shows how to create a cache for the `Square` function:

```go
myCache, err := mackenzie.Create[int](Square, mackenzie.Config{Lifetime: 1 * time.Second})
if err != nil {
    log.Fatalf("Failed to create cache: %v", err)
}
```

### Fetching Data from Cache
This demonstrates how to get data from the cache. If the data is not present in the cache, the `Square` function will be called, and the result will be cached.

```go
result, err := myCache.Get(&Point{1, 2})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)  // Outputs the product of 1 and 2, i.e., 2
```

### Checking Cache Expiry
You can check how long it takes for a cached item to expire:

```go
expiry := myCache.ExpiresIn(&Point{1, 2})
fmt.Printf("Expiry time: %v\n", expiry)
```

### Forcing Data Fetch
This shows how to force a fetch, which bypasses the cache and retrieves fresh data:

```go
result, err := myCache.ForceGet(&Point{1, 2})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)  // Outputs the product of 1 and 2, i.e., 2
```

### Error Handling
In case you pass incorrect arguments to the `Get` function, `mackenzie` provides specific errors that you can check:

```go
_, err = myCache.Get(1)
if mackenzie.IsErrCallInvalid(err) {
    fmt.Println("Invalid function call!")
}
if mackenzie.IsErrMackenzie(err) {
    fmt.Println("A Mackenzie-specific error occurred!")
}
```
