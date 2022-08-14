package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	clientID     = "myclient"
	clientSecret = "JwEqBbJyYNaQLe4NSKXMs3JEKZ017YdH"
)

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/realms/myrealm")

	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8081/oauth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := "123"

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/oauth/callback", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Query().Get("state") != state {
			http.Error(writer, "state did not match", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(ctx, request.URL.Query().Get("code"))

		if err != nil {
			http.Error(writer, "Falha ao trocar o token", http.StatusBadRequest)
			return
		}

		resp := struct {
			AccessToken *oauth2.Token
		}{
			AccessToken: token,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Write(data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
