package main

import (
	"github.com/pborges/gbridge/oauth"
	"io/ioutil"
	"log"
	"net/http"
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("HTTP ", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	authProvider := oauth.MapBasedAuthProvider{}
	authProvider.RegisterClient("123456", "654321")
	oathServer := oauth.Server{AuthProvider: &authProvider}

	mux.HandleFunc("/oauth", oathServer.HandleAuth())
	mux.HandleFunc("/token", oathServer.HandleToken())

	mux.HandleFunc("/smarthome", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		for name, headers := range r.Header {
			for _, h := range headers {
				log.Println("\tHEADER:", name, "->", h)
			}
		}

		//var client Client
		//accessToken := r.Header.Get("Authorization")
		//clientDBLock.Lock()
		//for _, c := range clientDB {
		//	if "Bearer "+c.AccessToken == accessToken {
		//		client = c
		//		break
		//	}
		//}
		//clientDBLock.Unlock()
		//
		//if client.ID == "" {
		//	http.Error(w, "{success: false, error: \"failed auth\"}", http.StatusForbidden)
		//	return
		//}

		body, _ := ioutil.ReadAll(r.Body)
		log.Println(string(body))
	})

	log.Fatal(http.ListenAndServe(":8085", logger(mux)))
}
