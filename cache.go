// Package cache implements middleware caching Authorization header to authenticate user.
package cache

import (
	"net/http"
	"time"

	gcache "github.com/pmylund/go-cache"
)

// Type Cache stores cache for basic auth, and record.
type Cache struct {
	gcache.Cache
}

const (
	DefaultAuthExpireTime = 10 * time.Minute
	DefaultAuthPurseTime  = 60 * time.Second
	AuthenticatedHeader   = "go-auth-cache-Authenticated"
)

// NewDefault returns a Cache with default value for AuthCache, RecordCache
func NewDefault() *Cache {
	return &Cache{
		// Cache for basic authentication:
		//      Default expiration time: 10 minutes.
		//      Purges expired items time: every 60 seconds.
		*gcache.New(DefaultAuthExpireTime, DefaultAuthPurseTime),
	}
}

// New returns a Cache with given expire time, and purse time.
func New(authExpireTime, authPurseTime time.Duration) *Cache {
	return &Cache{
		// Cache for basic authentication:
		//      Default expiration time: 10 minutes.
		//      Purges expired items time: every 60 seconds.
		*gcache.New(authExpireTime, authPurseTime),
	}
}

// ServeHTTP inform next middleware this request already authenticated by looking into the cache.
func (c *Cache) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	found := false

	if (req.Method != "GET") && (req.Method != "POST") {
		next(w, req)
		return
	}
	// Get credential from request header.
	credential := req.Header.Get("Authorization")
	// Get authentication status by credential.
	authenticated, found := c.Get(credential)

	// Cache hit
	if found && (authenticated == "true") {
		// Inform next middleware this request is authenticated.
		w.Header().Set(AuthenticatedHeader, "true")
		next(w, req)
	} else {
		// Cache miss. Unauthenticated.
		w.Header().Set(AuthenticatedHeader, "false")

		// Call next middleware
		next(w, req)

		// Set cache for this credential.
		authenticated = w.Header().Get(AuthenticatedHeader)
		if authenticated == "true" {
			c.Set(credential, authenticated, gcache.DefaultExpiration)
		}

		// Remove AuthenticatedHeader out of response header.
		w.Header().Del(AuthenticatedHeader)
	}
}
