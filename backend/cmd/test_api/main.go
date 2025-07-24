package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

type UsersResponse struct {
	Users []struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	} `json:"users"`
}

func main() {
	baseURL := "http://localhost:8081"

	// 1. Se connecter pour obtenir un token
	fmt.Println("1. Connexion...")
	loginData := LoginRequest{
		Email:    "superadmin@eduqr.com",
		Password: "superadmin123",
	}

	loginJSON, _ := json.Marshal(loginData)
	resp, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		fmt.Printf("Erreur de connexion: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var loginResp LoginResponse
	json.Unmarshal(body, &loginResp)

	// La réponse de login ne contient pas de champ "success"
	if loginResp.Token == "" {
		fmt.Printf("Échec de la connexion: %s\n", string(body))
		return
	}

	token := loginResp.Token
	fmt.Printf("Token obtenu: %s...\n", token[:20])

	// 2. Récupérer tous les utilisateurs
	fmt.Println("\n2. Récupération des utilisateurs...")
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/all", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Erreur lors de la récupération des utilisateurs: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	var usersResp UsersResponse
	json.Unmarshal(body, &usersResp)

	fmt.Printf("Nombre total d'utilisateurs: %d\n", len(usersResp.Users))
	fmt.Println("\nListe des utilisateurs:")
	for _, user := range usersResp.Users {
		fmt.Printf("- ID: %d, Email: %s, Nom: %s %s, Rôle: %s\n",
			user.ID, user.Email, user.FirstName, user.LastName, user.Role)
	}

	// 3. Filtrer les professeurs
	fmt.Println("\n3. Professeurs uniquement:")
	professors := 0
	for _, user := range usersResp.Users {
		if user.Role == "professeur" {
			fmt.Printf("- %s %s (%s)\n", user.FirstName, user.LastName, user.Email)
			professors++
		}
	}
	fmt.Printf("Nombre de professeurs: %d\n", professors)
}
