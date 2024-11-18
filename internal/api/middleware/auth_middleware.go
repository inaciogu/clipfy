package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
)

type CognitoConfig struct {
	UserPoolID string
	Region     string
	JWKSURL    string
}

// Configurações específicas do Cognito
var cognitoConfig = CognitoConfig{
	UserPoolID: os.Getenv("USER_POOL_ID"),
	Region:     "us-east-1",
	JWKSURL:    fmt.Sprintf("https://cognito-idp.us-east-1.amazonaws.com/%s/.well-known/jwks.json", os.Getenv("USER_POOL_ID")),
}

// ContextKey define uma chave para armazenar dados no contexto
type ContextKey string

const UserContextKey ContextKey = "user"

func parseRSAPublicKey(modulus string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(modulus)
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar modulus: %v", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString("AQAB") // Exponent is usually fixed as 65537 (0x10001)
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar exponent: %v", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes).Int64()

	pubKey := &rsa.PublicKey{
		N: n,
		E: int(e),
	}

	pubKeyPEM, err := encodeRSAPublicKeyToPEM(pubKey)
	if err != nil {
		return nil, fmt.Errorf("falha ao codificar chave pública: %v", err)
	}

	return jwt.ParseRSAPublicKeyFromPEM(pubKeyPEM)
}

func encodeRSAPublicKeyToPEM(pubKey *rsa.PublicKey) ([]byte, error) {
	pubKeyBytes := x509.MarshalPKCS1PublicKey(pubKey)
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})
	return pubKeyPEM, nil
}

// Função para validar o JWT (como no exemplo anterior)
func validateAndDecodeJWT(tokenString, jwksURL string) (*jwt.Token, jwt.MapClaims, error) {
	jwks, err := fetchJWKS(jwksURL)
	if err != nil {
		return nil, nil, err
	}

	// Função para buscar a chave pública com base no kid
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if kid, ok := token.Header["kid"].(string); ok {
			if modulus, exists := jwks[kid]; exists {
				return parseRSAPublicKey(modulus)
			}
		}
		return nil, errors.New("chave pública não encontrada para o kid")
	}

	// Validar o token
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao validar JWT: %v", err)
	}

	// Validar claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, errors.New("JWT inválido ou claims inválidos")
}

type JWKSKey struct {
	Keys []struct {
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
		Kty string `json:"kty"`
		Alg string `json:"alg"`
		Use string `json:"use"`
	} `json:"keys"`
}

func fetchJWKS(jwksURL string) (map[string]string, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar JWKS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status inesperado: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler JWKS: %v", err)
	}

	var jwks JWKSKey
	err = json.Unmarshal(body, &jwks)
	if err != nil {
		return nil, fmt.Errorf("falha ao fazer parse do JWKS: %v", err)
	}

	keys := make(map[string]string)
	for _, key := range jwks.Keys {
		keys[key.Kid] = key.N // Guardamos o Modulus (n) associado ao kid
	}
	return keys, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair o token do cabeçalho Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		// O formato esperado é: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token mal formado"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar o token JWT
		_, claims, err := validateAndDecodeJWT(tokenString, cognitoConfig.JWKSURL)
		if err != nil {
			fmt.Printf("falha ao validar token JWT: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalido"})
			c.Abort()
			return
		}

		c.Set(string(UserContextKey), claims)
		c.Next()
	}
}
