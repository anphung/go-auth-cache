# go-auth-cache
## Description

This middleware assigns value "true" to header "go-auth-cache-Authenticated" to http.ResponseWriter if http.Request's header "Authorization" in this middleware's cache (cache hit). Assigns "false" otherwise.

One efficient usage is for Basic authentication. We use this middleware to know if the http.Request already authenticated in the past, avoid overhead of authenticating the same user over and over again.

Before exists, this middleware cache the request if http.ResponseWriter's header "go-auth-cache-Authenticated" is set to "cache".

## Usage:

```go
authCacheMiddleware := cache.NewDefault()
```

We then plug ```authCacheMiddleware``` in front of any middleware. Check [this](https://github.com/anphung/negroni-auth-dynamodb/commit/41ea28a8a2ac40db369ab90009c38ecede3cba2b) for example.

_Note:_ Any of the following middleware of this middleware has to __explicitly__ set http.ResponseWriter's header "go-auth-cache-Authenticated" to:
* "cache", if it want the request to be cached.
* don't set, do nothing.
