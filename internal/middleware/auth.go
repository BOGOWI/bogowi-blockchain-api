package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
)

type AuthMiddleware struct {
	firebaseProjectID string
	jwksURL           string
	tokenCache        *cache.Cache
	jwksCache         *cache.Cache
}

type FirebaseClaims struct {
	jwt.RegisteredClaims
	WalletAddress string                 `json:"walletAddress"`
	Firebase      map[string]interface{} `json:"firebase"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func NewAuthMiddleware(projectID string) *AuthMiddleware {
	return &AuthMiddleware{
		firebaseProjectID: projectID,
		jwksURL:           "https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com",
		tokenCache:        cache.New(5*time.Minute, 10*time.Minute),
		jwksCache:         cache.New(1*time.Hour, 2*time.Hour),
	}
}

func (a *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Check cache first
		if cachedClaims, found := a.tokenCache.Get(tokenString); found {
			claims := cachedClaims.(*FirebaseClaims)
			ctx := context.WithValue(r.Context(), "wallet", claims.WalletAddress)
			ctx = context.WithValue(ctx, "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Verify token
		claims, err := a.verifyToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// Cache valid token
		a.tokenCache.Set(tokenString, claims, cache.DefaultExpiration)

		// Add to context
		ctx := context.WithValue(r.Context(), "wallet", claims.WalletAddress)
		ctx = context.WithValue(ctx, "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (a *AuthMiddleware) verifyToken(tokenString string) (*FirebaseClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &FirebaseClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get key ID
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid in token header")
		}

		// Get public key
		return a.getPublicKey(kid)
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*FirebaseClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify issuer
	expectedIssuer := fmt.Sprintf("https://securetoken.google.com/%s", a.firebaseProjectID)
	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	// Verify audience
	if claims.Audience == nil || len(claims.Audience) == 0 || claims.Audience[0] != a.firebaseProjectID {
		return nil, fmt.Errorf("invalid audience")
	}

	// Check wallet address
	if claims.WalletAddress == "" {
		return nil, fmt.Errorf("missing wallet address in claims")
	}

	return claims, nil
}

func (a *AuthMiddleware) getPublicKey(kid string) (interface{}, error) {
	// Check cache
	cacheKey := "jwks"
	var jwks *JWKS

	if cachedJWKS, found := a.jwksCache.Get(cacheKey); found {
		jwks = cachedJWKS.(*JWKS)
	} else {
		// Fetch JWKS
		resp, err := http.Get(a.jwksURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
		}
		defer resp.Body.Close()

		jwks = &JWKS{}
		if err := json.NewDecoder(resp.Body).Decode(jwks); err != nil {
			return nil, fmt.Errorf("failed to decode JWKS: %v", err)
		}

		// Cache JWKS
		a.jwksCache.Set(cacheKey, jwks, cache.DefaultExpiration)
	}

	// Find key by kid
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			// Convert JWK to RSA public key
			return jwt.ParseRSAPublicKeyFromPEM([]byte(fmt.Sprintf(
				"-----BEGIN RSA PUBLIC KEY-----\n%s\n-----END RSA PUBLIC KEY-----",
				key.N,
			)))
		}
	}

	return nil, fmt.Errorf("key not found: %s", kid)
}

func GetWalletFromContext(ctx context.Context) (string, error) {
	wallet, ok := ctx.Value("wallet").(string)
	if !ok {
		return "", fmt.Errorf("wallet not found in context")
	}
	return wallet, nil
}

func GetClaimsFromContext(ctx context.Context) (*FirebaseClaims, error) {
	claims, ok := ctx.Value("claims").(*FirebaseClaims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}
