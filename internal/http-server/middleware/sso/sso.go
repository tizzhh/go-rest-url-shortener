package sso

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	ssogrpc "url-shortener/internal/clients/sso/grpc"
)

func IsRequestAdmin(realm string, ssoClient *ssogrpc.Client, timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			auth = strings.Split(auth, "Basic ")[1]
			email, err := getEmailFromAuth(auth)
			if err != nil {
				respondWithError(w, realm, http.StatusInternalServerError)
				return
			}
			if email == "" {
				respondWithError(w, realm, http.StatusUnauthorized)
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			isAdmin, err := ssoClient.IsAdmin(ctx, email)
			if err != nil {
				respondWithError(w, realm, http.StatusInternalServerError)
				return
			}
			if !isAdmin {
				respondWithError(w, realm, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getEmailFromAuth(input string) (string, error) {
	const caller = "middleware.sso.getEmailFromAuth"

	auth, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	email, _, ok := strings.Cut(string(auth), ":")
	if !ok {
		return "", fmt.Errorf("%s: %w", caller, err)
	}
	return string(email), nil
}

func respondWithError(w http.ResponseWriter, realm string, code int) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(code)
}
