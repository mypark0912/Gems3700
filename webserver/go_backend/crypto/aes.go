package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var AESKey = []byte("ntekSystem_20250721_mypark_caner") // 32 bytes

type AESCipher struct {
	iv        []byte
	adminIV   []byte
	adminData []byte
	configDir string
}

func NewAESCipher(configDir string) (*AESCipher, error) {
	hash := sha256.Sum256(AESKey)
	iv := hash[:16]

	c := &AESCipher{
		iv:        iv,
		configDir: configDir,
	}

	adminPath := filepath.Join(configDir, "admin_secret.enc")
	raw, err := os.ReadFile(adminPath)
	if err == nil {
		decoded, err := base64.StdEncoding.DecodeString(string(raw))
		if err == nil && len(decoded) > 16 {
			c.adminIV = decoded[:16]
			c.adminData = decoded[16:]
		}
	}

	return c, nil
}

func (c *AESCipher) CheckAdmin(password string) bool {
	if c.adminIV == nil || c.adminData == nil {
		return false
	}
	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return false
	}
	mode := cipher.NewCBCDecrypter(block, c.adminIV)
	plaintext := make([]byte, len(c.adminData))
	mode.CryptBlocks(plaintext, c.adminData)
	plaintext = pkcs7Unpad(plaintext)
	return string(plaintext) == password
}

func (c *AESCipher) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return "", err
	}
	padded := pkcs7Pad([]byte(plaintext), aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, c.iv)
	mode.CryptBlocks(ciphertext, padded)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *AESCipher) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return "", err
	}
	if len(data)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of block size")
	}
	mode := cipher.NewCBCDecrypter(block, c.iv)
	plaintext := make([]byte, len(data))
	mode.CryptBlocks(plaintext, data)
	plaintext = pkcs7Unpad(plaintext)
	return string(plaintext), nil
}

type InfluxConfig struct {
	URL   string `json:"url"`
	Token string `json:"token"`
	Org   string `json:"org"`
	OrgID string `json:"org_id"`
}

func (c *AESCipher) GetInflux() (*InfluxConfig, error) {
	path := filepath.Join(c.configDir, "influx.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg InfluxConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padBytes := make([]byte, padding)
	for i := range padBytes {
		padBytes[i] = byte(padding)
	}
	return append(data, padBytes...)
}

func pkcs7Unpad(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding > aes.BlockSize {
		return data
	}
	return data[:len(data)-padding]
}
