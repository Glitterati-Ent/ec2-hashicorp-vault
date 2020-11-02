package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"log"
	"os"
)

func main() {

	text := `path "secret/data/*" {
  capabilities = ["create", "update"]
}

path "secret/data/policy" {
  capabilities = ["read", "write"]
}`

	err := os.Setenv(api.EnvVaultToken, "education")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Setenv(api.EnvVaultAddress, "http://127.0.0.1:8200")
	if err != nil {
		log.Fatal(err)
	}
	token := os.Getenv(api.EnvVaultToken

	addr := os.Getenv(api.EnvVaultAddress)

	config := &api.Config{
		Address: addr,
	}
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	client.SetToken(token)

	policies, err := client.Sys().ListPolicies()
	if err != nil {
		log.Fatal(err)
	}

	client.Sys().PutPolicy("my_policy", text)

	policy, err := client.Sys().GetPolicy("my_policy")
	if err != nil {
		log.Fatal(err)
	}

	userPolicy := []string{"my_policy"}

	user, err := client.Auth().Token().Create(&api.TokenCreateRequest{DisplayName: "test", Policies: userPolicy })
	if err != nil {
		log.Fatal(err)
	}

	client.ClearToken()
	client.SetToken(user.Auth.ClientToken)

	ud := map[string]interface{}{"user": user}
	_, err = client.Logical().Write("secret/data/policy", ud)
	if err != nil {
		log.Fatal(err)
	}

	secret, err := client.Logical().Read("secret/data/policy")
	if err != nil {
		log.Fatal(err)
	}
	m, ok := secret.Data["user"].(map[string]interface{})
	if !ok {
		fmt.Printf("%T %#v\n", secret.Data["user"], secret.Data["user"])
		return
	}
	fmt.Println("MAP: ", m)

}
