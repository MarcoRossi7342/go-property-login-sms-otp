package main

import (
	"encoding/json"
	"log"
	"net/http"
	"propertylogin/infrai"
	"propertylogin/property"
)

type verifyRequest struct { To, Code, TenantID string }

func main() {
	c, err := infrai.New(); if err != nil { log.Fatal(err) }
	http.HandleFunc("/login/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		var in verifyRequest; if json.NewDecoder(r.Body).Decode(&in) != nil { http.Error(w, "invalid json", http.StatusBadRequest); return }
		env, err := c.VerifyOTP(in.To, in.Code, "property-login-"+in.To); if err != nil { http.Error(w, err.Error(), http.StatusBadGateway); return }
		decision, err := property.VerifyLogin(env.OK, in.TenantID); if err != nil { http.Error(w, err.Error(), http.StatusUnauthorized); return }
		json.NewEncoder(w).Encode(decision)
	})
	log.Println("property login listening on :8080"); log.Fatal(http.ListenAndServe(":8080", nil))
}
