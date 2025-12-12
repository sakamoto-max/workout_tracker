package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"workout_tracker/auth"
)


type contexKey string

var ClaimsKey contexKey = "claims"


// what can admin do?
// add exercises to the exercises table

// user cannot add to the exercise table

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		cookie, err := r.Cookie("jwtToken")
		if err != nil {
			repsonse := map[string]string{
				"message" : "error getting cookie from request",
			}

			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(repsonse)
		}
		
		tokenFormCookie := cookie.Value

		if tokenFormCookie == "" {
			repsonse := map[string]string{
				"message" : "token is empty",
			}
		
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(repsonse)
		}else {
			claims, err := auth.VerifyJwtToken(tokenFormCookie)

			if err != nil{
				repsonse := map[string]string{
					"message" : "token is not valid",
				}
		
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(repsonse)
			}else {

				ctx := context.WithValue(context.Background(), ClaimsKey, claims)
	
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}			
	})
}

func GetValuesFromJWT() {

}

