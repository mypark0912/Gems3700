package setting

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"serverGO/infra"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type CMDManager struct {
	deps  *infra.Dependencies
	seqNo uint32
}

type RequestMessage struct {
	Seq   uint32 `json:"seq"`
	Cmd   string `json:"cmd"`
	CanID int    `json:"canid,omitempty"`
}

func NewCMDManager(deps *infra.Dependencies) *CMDManager {
	return &CMDManager{deps: deps}
}

func (m *CMDManager) nextSeq() uint32 {
	return atomic.AddUint32(&m.seqNo, 1)
}

func (m *CMDManager) pushRequest(req RequestMessage) (uint32, error) {
	seq := m.nextSeq()
	req.Seq = seq

	data, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	ctx := context.Background()
	err = m.deps.Redis.Client0.LPush(ctx, "Request", string(data)).Err()
	if err != nil {
		return 0, err
	}

	return seq, nil
}

// WebSocket 핸들러
func (m *CMDManager) handleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ws upgrade error:", err)
		return
	}
	defer conn.Close()
	log.Println("ws connected:", c.ClientIP())

	// 응답 대기 채널 관리
	var mu sync.Mutex
	var wsMu sync.Mutex // WebSocket 쓰기 보호
	pending := make(map[uint32]chan map[string]interface{})

	// Response 큐 감시 고루틴 (단일)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				result, err := m.deps.Redis.Client0.RPop(ctx, "Response").Result()
				if err != nil {
					time.Sleep(100 * time.Millisecond)
					continue
				}

				var resp map[string]interface{}
				if err := json.Unmarshal([]byte(result), &resp); err != nil {
					continue
				}

				respSeq, ok := resp["seq"].(float64)
				if !ok {
					continue
				}

				mu.Lock()
				ch, exists := pending[uint32(respSeq)]
				mu.Unlock()

				if exists {
					ch <- resp
				} else {
					// 매칭 안 되면 다시 넣음
					m.deps.Redis.Client0.LPush(ctx, "Response", result)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()

	// WebSocket 메시지 수신 루프
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("ws disconnected:", c.ClientIP())
			return
		}

		var req RequestMessage
		if err := json.Unmarshal(msg, &req); err != nil {
			wsMu.Lock()
			conn.WriteJSON(gin.H{"error": "invalid message"})
			wsMu.Unlock()
			continue
		}

		seq, err := m.pushRequest(req)
		if err != nil {
			wsMu.Lock()
			conn.WriteJSON(gin.H{"error": err.Error()})
			wsMu.Unlock()
			continue
		}

		// 응답 대기 채널 등록
		ch := make(chan map[string]interface{}, 1)
		mu.Lock()
		pending[seq] = ch
		mu.Unlock()

		// 응답 대기 고루틴
		go func(seq uint32, ch chan map[string]interface{}) {
			select {
			case resp := <-ch:
				wsMu.Lock()
				conn.WriteJSON(resp)
				wsMu.Unlock()
			case <-time.After(10 * time.Second):
				wsMu.Lock()
				conn.WriteJSON(gin.H{"seq": seq, "error": "timeout"})
				wsMu.Unlock()
			}
			mu.Lock()
			delete(pending, seq)
			mu.Unlock()
		}(seq, ch)
	}
}

// 테스트용: scan_resp
func (m *CMDManager) testScanResp(c *gin.Context) {
	ctx := context.Background()

	result, err := m.deps.Redis.Client0.RPop(ctx, "Request").Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "no request in queue"})
		return
	}

	var req RequestMessage
	json.Unmarshal([]byte(result), &req)

	resp := map[string]interface{}{
		"seq": req.Seq,
		"cmd": "scan_resp",
		"tapboxes": []map[string]interface{}{
			{"canid": 300001, "type": (1 << 8) | (1 << 4) | 2},
			{"canid": 300002, "type": (5 << 8) | (5 << 4) | 2},
			{"canid": 300003, "type": (1 << 8) | (1 << 4) | 3},
			{"canid": 300004, "type": (4 << 8) | (4 << 4) | 3},
		},
	}

	data, _ := json.Marshal(resp)
	m.deps.Redis.Client0.LPush(ctx, "Response", string(data))

	c.JSON(http.StatusOK, gin.H{"status": "ok", "seq": req.Seq})
}

// 테스트용: ping_resp (Request 큐의 ping_req를 모두 처리)
func (m *CMDManager) testPingResp(c *gin.Context) {
	ctx := context.Background()
	count := 0

	for {
		result, err := m.deps.Redis.Client0.RPop(ctx, "Request").Result()
		if err != nil {
			break
		}

		var req RequestMessage
		json.Unmarshal([]byte(result), &req)

		if req.Cmd != "ping_req" {
			m.deps.Redis.Client0.LPush(ctx, "Request", result)
			break
		}

		resp := map[string]interface{}{
			"seq":   req.Seq,
			"cmd":   "ping_resp",
			"canid": req.CanID,
			"res":   "ok",
		}

		data, _ := json.Marshal(resp)
		m.deps.Redis.Client0.LPush(ctx, "Response", string(data))
		count++
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "count": count})
}
