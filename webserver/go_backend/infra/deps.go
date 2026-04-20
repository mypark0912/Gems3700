package infra

import (
	"serverGO/config"
	"serverGO/crypto"
)

type Dependencies struct {
	Config *config.AppConfig
	Redis  *RedisState
	Influx *InfluxState
	Crypto *crypto.AESCipher
}
