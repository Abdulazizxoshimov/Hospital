package token

import (
	"net/http"
	"strings"

	"github.com/Abdulazizxoshimov/Hospital/config"
	"github.com/spf13/cast"
)

func GetIdFromToken(r *http.Request, cfg *config.Config) (string, int) {
	var softToken string
	token := r.Header.Get("Authorization")

	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		softToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		softToken = token
	}

	claims, err := ExtractClaim(softToken, []byte(cfg.Token.SignInKey))
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}

	return cast.ToString(claims["sub"]), 0
}

func GetRoleFromToken(r *http.Request, cfg *config.Config) (string, int) {
	var softToken string
	token := r.Header.Get("Authorization")

	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		softToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		softToken = token
	}

	claims, err := ExtractClaim(softToken, []byte(cfg.Token.SignInKey))
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}

	return cast.ToString(claims["role"]), 0
}