package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nihar-hegde/valtro-backend/internal/repositories/user"
	"github.com/nihar-hegde/valtro-backend/internal/utils/response"
	"gorm.io/gorm"
)

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// Thread-safe caches with proper synchronization
type JWTCache struct {
	jwks        *JWKSResponse
	jwksTime    time.Time
	userIDCache map[string]string // clerkID -> internalID
	userIDTime  map[string]time.Time
	mu          sync.RWMutex
}

var jwtCache = &JWTCache{
	userIDCache: make(map[string]string),
	userIDTime:  make(map[string]time.Time),
}

// ClerkJWTMiddleware validates Clerk JWT tokens using their JWKS endpoint
func ClerkJWTMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.SendUnauthorized(w, "Authorization header required")
			return
		}

		// Extract Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.SendUnauthorized(w, "Invalid authorization format. Expected 'Bearer <token>'")
			return
		}

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Get the kid from token header
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("no kid found in token header")
			}

			// Get public key from Clerk's JWKS endpoint
			publicKey, err := getPublicKeyFromJWKS(kid)
			if err != nil {
				return nil, fmt.Errorf("failed to get public key: %v", err)
			}

			return publicKey, nil
		})

		if err != nil {
			response.SendUnauthorized(w, "Invalid token: "+err.Error())
			return
		}

		// Validate token and extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract Clerk user ID from 'sub' claim
			clerkUserID, ok := claims["sub"].(string)
			if !ok || clerkUserID == "" {
				response.SendUnauthorized(w, "Invalid token: user ID not found")
				return
			}

			// Convert Clerk user ID to internal UUID (with caching)
			internalUserID, err := getCachedInternalUserID(db, clerkUserID)
			if err != nil {
				response.SendUnauthorized(w, "User not found: "+err.Error())
				return
			}

			// Set X-User-ID header with internal UUID for organization handlers
			r.Header.Set("X-User-ID", internalUserID)

			// Add both IDs to context
			ctx := context.WithValue(r.Context(), "clerkUserID", clerkUserID)
			ctx = context.WithValue(ctx, "internalUserID", internalUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			response.SendUnauthorized(w, "Invalid token claims")
			return
		}
	})
	}
}

// getPublicKeyFromJWKS fetches public key from Clerk's JWKS endpoint with thread-safe caching
func getPublicKeyFromJWKS(kid string) (*rsa.PublicKey, error) {
	// Check cached JWKS with read lock
	jwtCache.mu.RLock()
	if jwtCache.jwks != nil && time.Since(jwtCache.jwksTime) < 5*time.Minute {
		cachedJWKS := jwtCache.jwks
		jwtCache.mu.RUnlock()
		return findKeyInJWKS(cachedJWKS, kid)
	}
	jwtCache.mu.RUnlock()

	// Get Clerk domain from environment
	clerkDomain := os.Getenv("CLERK_FRONTEND_API")
	if clerkDomain == "" {
		return nil, fmt.Errorf("CLERK_FRONTEND_API not configured")
	}

	// Fetch JWKS from Clerk with timeout
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", clerkDomain)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JWKS endpoint returned status: %d", resp.StatusCode)
	}

	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %v", err)
	}

	// Cache the JWKS with write lock
	jwtCache.mu.Lock()
	jwtCache.jwks = &jwks
	jwtCache.jwksTime = time.Now()
	jwtCache.mu.Unlock()

	return findKeyInJWKS(&jwks, kid)
}

// findKeyInJWKS finds the public key with matching kid
func findKeyInJWKS(jwks *JWKSResponse, kid string) (*rsa.PublicKey, error) {
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return jwkToRSAPublicKey(key)
		}
	}
	return nil, fmt.Errorf("key with kid '%s' not found", kid)
}

// jwkToRSAPublicKey converts JWK to RSA public key
func jwkToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decode the modulus
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %v", err)
	}

	// Decode the exponent
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %v", err)
	}

	// Convert to big integers
	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

// getCachedInternalUserID converts Clerk user ID to internal UUID with caching
func getCachedInternalUserID(db *gorm.DB, clerkUserID string) (string, error) {
	const userCacheTTL = 10 * time.Minute

	// Check cache first with read lock
	jwtCache.mu.RLock()
	if internalID, exists := jwtCache.userIDCache[clerkUserID]; exists {
		if cacheTime, timeExists := jwtCache.userIDTime[clerkUserID]; timeExists {
			if time.Since(cacheTime) < userCacheTTL {
				jwtCache.mu.RUnlock()
				return internalID, nil
			}
		}
	}
	jwtCache.mu.RUnlock()

	// Cache miss or expired - fetch from database
	userRepo := user.NewRepository(db)
	userModel, err := userRepo.GetByClerkUserID(clerkUserID)
	if err != nil {
		return "", fmt.Errorf("user with Clerk ID %s not found", clerkUserID)
	}

	internalID := userModel.ID.String()

	// Update cache with write lock
	jwtCache.mu.Lock()
	jwtCache.userIDCache[clerkUserID] = internalID
	jwtCache.userIDTime[clerkUserID] = time.Now()
	jwtCache.mu.Unlock()

	return internalID, nil
}

// getInternalUserID is kept for backward compatibility but uses the cached version
func getInternalUserID(db *gorm.DB, clerkUserID string) (string, error) {
	return getCachedInternalUserID(db, clerkUserID)
}
