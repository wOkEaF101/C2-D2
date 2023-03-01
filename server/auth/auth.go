package auth

import (
	"C2-D2/server/models"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Super Secret
var jwtKey = []byte("BEEPBOOPBEEPBOOP")

// Generate Signed JWT for the given agent
func Login(agent models.Agent) models.Response {
	var err error

	// Claims to appear in the token - only base64 encoded so no secrets!
	claims := &jwt.MapClaims{
		"authorized": true,
		"agent":      agent.UUID,
	}

	// Generate the gosh darn thing
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		response := models.Response{Type: "ERROR", Message: "Server Error!", Data: err.Error()}
		return response
	}

	response := models.Response{Type: "SUCCESS", Message: agent.UUID, Data: tokenString}
	return response
}

// Validate JWT and return boolean if it is valid - don't think we need robust errors here... Yet.
func Auth(agent models.Agent) bool {
	tokenString := agent.Token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	// Error and not validated
	if err != nil {
		return false
	}

	// Token not valid - so it's not valid
	if !token.Valid {
		return false
	}

	// Token is valid and we're good to go
	return true
}
