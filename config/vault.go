package config

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/vault/api"
)

func LoadVaultSecrets(cfg *Config) {
	if cfg.Vault.Host == "" {
		log.Println("Vault host is empty, skipping Vault initialization")
		return
	}

	vaultCfg := api.DefaultConfig()
	vaultCfg.Address = cfg.Vault.Host

	client, err := api.NewClient(vaultCfg)
	if err != nil {
		log.Printf("Failed to create Vault client: %v", err)
		return
	}

	client.SetToken(cfg.Vault.Token)

	log.Printf("Vault client initialized: host=%s, path=%s", cfg.Vault.Host, cfg.Vault.Path)

	secret, err := client.Logical().Read(fmt.Sprintf("secret/data/%s", cfg.Vault.Path))
	if err != nil {
		log.Printf("Failed to read secret from Vault: %v", err)
		return
	}

	if secret == nil || secret.Data == nil {
		log.Printf("Secret not found at path: %s", cfg.Vault.Path)
		return
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Println("Invalid secret data format")
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal secret data: %v", err)
		return
	}

	var secretConfig SecretVaultConfig
	if err := json.Unmarshal(jsonData, &secretConfig); err != nil {
		log.Printf("Failed to unmarshal secret config: %v", err)
		return
	}

	// Config에 Vault에서 가져온 값 적용
	if secretConfig.DatabaseSecret.Password != "" {
		cfg.Database.Password = secretConfig.DatabaseSecret.Password
		log.Println("Database password loaded from Vault")
	}
	if secretConfig.RedisSecret.Password != "" {
		cfg.Redis.Password = secretConfig.RedisSecret.Password
		log.Println("Redis password loaded from Vault")
	}
	if secretConfig.JWTSecret != "" {
		cfg.Secret.JwtSecret = secretConfig.JWTSecret
		log.Println("JWT secret loaded from Vault")
	}
}
