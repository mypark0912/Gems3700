package infra

import (
	"fmt"
	"log"
	"sync"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"serverGO/config"
	"serverGO/crypto"
)

type InfluxState struct {
	mu       sync.Mutex
	client   influxdb2.Client
	queryAPI api.QueryAPI
	writeAPI api.WriteAPI
	deleteAPI api.DeleteAPI
	cfg      *config.AppConfig
	cipher   *crypto.AESCipher
	org      string
	err      error
}

func NewInfluxState(cfg *config.AppConfig, cipher *crypto.AESCipher) *InfluxState {
	return &InfluxState{cfg: cfg, cipher: cipher}
}

func (s *InfluxState) connect() {
	influxCfg, err := s.cipher.GetInflux()
	if err != nil {
		s.err = fmt.Errorf("influx config: %w", err)
		log.Println("InfluxDB config error:", err)
		return
	}

	token, err := s.cipher.Decrypt(influxCfg.Token)
	if err != nil {
		s.err = fmt.Errorf("influx token decrypt: %w", err)
		log.Println("InfluxDB token decrypt error:", err)
		return
	}

	url := influxCfg.URL
	if url == "" {
		url = fmt.Sprintf("http://%s:8086", s.cfg.InfluxIP)
	}

	s.client = influxdb2.NewClientWithOptions(url, token,
		influxdb2.DefaultOptions().
			SetHTTPRequestTimeout(30))

	s.org = influxCfg.Org
	s.queryAPI = s.client.QueryAPI(s.org)
	s.writeAPI = s.client.WriteAPI(s.org, "") // bucket set per-call
	s.deleteAPI = s.client.DeleteAPI()
	s.err = nil

	log.Println("InfluxDB connected:", url)
}

func (s *InfluxState) QueryAPI() api.QueryAPI {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client == nil {
		s.connect()
	}
	return s.queryAPI
}

func (s *InfluxState) DeleteAPI() api.DeleteAPI {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client == nil {
		s.connect()
	}
	return s.deleteAPI
}

func (s *InfluxState) Client() influxdb2.Client {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client == nil {
		s.connect()
	}
	return s.client
}

func (s *InfluxState) Org() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client == nil {
		s.connect()
	}
	return s.org
}

func (s *InfluxState) Error() error {
	return s.err
}

func (s *InfluxState) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client != nil {
		s.client.Close()
	}
}
