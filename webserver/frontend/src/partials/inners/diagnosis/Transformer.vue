<template>
  <div class="transformer-diagnosis">
    <div class="header">
      <h2>변압기 진단 모니터링</h2>
      <div class="status-legend">
        <div class="legend-item">
          <div class="color-box" style="background: rgba(100, 116, 139, 0.8);"></div>
          <span>정지</span>
        </div>
        <div class="legend-item">
          <div class="color-box" style="background: rgba(16, 185, 129, 0.8);"></div>
          <span>정상</span>
        </div>
        <div class="legend-item">
          <div class="color-box" style="background: rgba(245, 158, 11, 0.8);"></div>
          <span>주의</span>
        </div>
        <div class="legend-item">
          <div class="color-box" style="background: rgba(249, 115, 22, 0.8);"></div>
          <span>경고</span>
        </div>
        <div class="legend-item">
          <div class="color-box" style="background: rgba(239, 68, 68, 0.8);"></div>
          <span>위험</span>
        </div>
      </div>
    </div>

    <div class="diagram-container">
      <svg width="700" height="500" viewBox="0 0 700 500" class="transformer-svg">
        <!-- 변압기 탱크 (외부 케이스) - 실제 탱크 형태 -->
        <path d="M120 160 L120 180 Q120 190 130 190 L570 190 Q580 190 580 180 L580 160 Q580 150 570 150 L130 150 Q120 150 120 160 Z" 
              fill="#e5e7eb" stroke="#374151" stroke-width="2"/>
        <rect x="130" y="190" width="440" height="200" 
              fill="#e5e7eb" stroke="#374151" stroke-width="2"/>
        <path d="M130 390 Q130 400 140 400 L560 400 Q570 400 570 390 L570 390 L130 390 Z" 
              fill="#e5e7eb" stroke="#374151" stroke-width="2"/>
        
        <!-- 냉각 라디에이터 (좌측) -->
        <rect x="80" y="200" width="50" height="180" 
              fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2" class="component-rect"/>
        <line x1="90" y1="210" x2="90" y2="370" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        <line x1="100" y1="210" x2="100" y2="370" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        <line x1="110" y1="210" x2="110" y2="370" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        <line x1="120" y1="210" x2="120" y2="370" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        
        <!-- 냉각 라디에이터 (우측) -->
        <rect x="570" y="200" width="50" height="180" 
              fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2" class="component-rect"/>
        <line x1="580" y1="210" x2="580" y2="370" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        <line x1="590" y1="210" x2="590" y2="370" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        <line x1="600" y1="210" x2="600" y2="370" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        <line x1="610" y1="210" x2="610" y2="370" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        
        <!-- 코어 (철심) - 실제 E자형 구조 -->
        <g class="component-group">
          <!-- 코어 메인 백 -->
          <rect x="320" y="200" width="60" height="180" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <!-- 코어 상단 연결부 -->
          <rect x="280" y="200" width="140" height="20" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <!-- 코어 하단 연결부 -->
          <rect x="280" y="360" width="140" height="20" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <!-- 코어 중앙 연결부 -->
          <rect x="280" y="285" width="140" height="30" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
        </g>
        
        <!-- 1차 권선 (고압측) - 코일 형태 -->
        <g class="component-group">
          <ellipse cx="270" cy="250" rx="25" ry="50" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="4"/>
          <ellipse cx="270" cy="250" rx="20" ry="45" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="3"/>
          <ellipse cx="270" cy="250" rx="15" ry="40" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="2"/>
          <ellipse cx="270" cy="330" rx="25" ry="50" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="4"/>
          <ellipse cx="270" cy="330" rx="20" ry="45" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="3"/>
          <ellipse cx="270" cy="330" rx="15" ry="40" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="2"/>
        </g>
        
        <!-- 2차 권선 (저압측) - 코일 형태 -->
        <g class="component-group">
          <ellipse cx="430" cy="250" rx="25" ry="50" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="4"/>
          <ellipse cx="430" cy="250" rx="20" ry="45" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="3"/>
          <ellipse cx="430" cy="250" rx="15" ry="40" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="2"/>
          <ellipse cx="430" cy="330" rx="25" ry="50" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="4"/>
          <ellipse cx="430" cy="330" rx="20" ry="45" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="3"/>
          <ellipse cx="430" cy="330" rx="15" ry="40" 
                   fill="none" stroke="rgba(245, 158, 11, 1)" stroke-width="2"/>
        </g>
        
        <!-- 고압 부싱 (HV Bushing) - 실제 부싱 형태 -->
        <g class="component-group">
          <rect x="245" y="80" width="20" height="70" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <ellipse cx="255" cy="70" rx="15" ry="10" 
                   fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <rect x="240" y="65" width="30" height="5" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        </g>
        
        <!-- 저압 부싱 (LV Bushing) -->
        <g class="component-group">
          <rect x="405" y="80" width="20" height="70" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <ellipse cx="415" cy="70" rx="15" ry="10" 
                   fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <rect x="400" y="65" width="30" height="5" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        </g>
        
        <!-- 탭 체인저 (OLTC) - 실제 탭체인저 형태 -->
        <g class="component-group">
          <circle cx="500" cy="280" r="35" 
                  fill="rgba(239, 68, 68, 0.8)" stroke="rgba(239, 68, 68, 1)" stroke-width="3"/>
          <circle cx="500" cy="280" r="25" 
                  fill="none" stroke="rgba(239, 68, 68, 1)" stroke-width="2"/>
          <circle cx="500" cy="280" r="15" 
                  fill="none" stroke="rgba(239, 68, 68, 1)" stroke-width="2"/>
          <!-- 탭 연결선들 -->
          <line x1="485" y1="265" x2="515" y2="295" stroke="rgba(239, 68, 68, 1)" stroke-width="2"/>
          <line x1="515" y1="265" x2="485" y2="295" stroke="rgba(239, 68, 68, 1)" stroke-width="2"/>
          <line x1="500" y1="250" x2="500" y2="310" stroke="rgba(239, 68, 68, 1)" stroke-width="2"/>
          <line x1="470" y1="280" x2="530" y2="280" stroke="rgba(239, 68, 68, 1)" stroke-width="2"/>
        </g>
        
        <!-- 커패시터 뱅크 (좌하단) -->
        <g class="component-group">
          <rect x="180" y="420" width="80" height="40" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <!-- 커패시터 단위들 -->
          <rect x="185" y="425" width="15" height="30" fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
          <rect x="205" y="425" width="15" height="30" fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
          <rect x="225" y="425" width="15" height="30" fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
          <rect x="245" y="425" width="15" height="30" fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        </g>
        
        <!-- 부하 연결부 (우하단) -->
        <g class="component-group">
          <polygon points="440,420 480,420 460,460" 
                   fill="rgba(249, 115, 22, 0.8)" stroke="rgba(249, 115, 22, 1)" stroke-width="2"/>
          <line x1="450" y1="430" x2="470" y2="430" stroke="rgba(249, 115, 22, 1)" stroke-width="2"/>
          <line x1="455" y1="440" x2="465" y2="440" stroke="rgba(249, 115, 22, 1)" stroke-width="2"/>
        </g>
        
        <!-- 스트레스 센서 (상단 우측) -->
        <g class="component-group">
          <circle cx="550" cy="120" r="15" 
                  fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <line x1="540" y1="110" x2="560" y2="130" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <line x1="560" y1="110" x2="540" y2="130" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <circle cx="550" cy="120" r="5" fill="rgba(16, 185, 129, 1)"/>
        </g>
        
        <!-- 온도 센서 (중앙 하단) -->
        <g class="component-group">
          <rect x="345" y="390" width="10" height="15" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
          <circle cx="350" cy="385" r="5" 
                  fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        </g>
        
        <!-- 오일 레벨 게이지 -->
        <g class="component-group">
          <rect x="530" y="180" width="25" height="50" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <rect x="535" y="185" width="15" height="40" fill="rgba(59, 130, 246, 0.6)"/>
          <circle cx="542" cy="175" r="3" fill="rgba(16, 185, 129, 1)"/>
          <line x1="538" y1="200" x2="548" y2="200" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
          <line x1="538" y1="210" x2="548" y2="210" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        </g>
        
        <!-- 접지 시스템 -->
        <g>
          <line x1="350" y1="400" x2="350" y2="430" stroke="#374151" stroke-width="3"/>
          <line x1="340" y1="430" x2="360" y2="430" stroke="#374151" stroke-width="3"/>
          <line x1="343" y1="435" x2="357" y2="435" stroke="#374151" stroke-width="2"/>
          <line x1="346" y1="440" x2="354" y2="440" stroke="#374151" stroke-width="1"/>
        </g>
        
        <!-- 전력선 연결 -->
        <line x1="255" y1="150" x2="255" y2="200" stroke="#374151" stroke-width="4"/>
        <line x1="415" y1="150" x2="415" y2="200" stroke="#374151" stroke-width="4"/>
        <line x1="465" y1="280" x2="455" y2="280" stroke="#374151" stroke-width="3"/>
        
        <!-- 절연 가스 모니터링 (상단 좌측) -->
        <g class="component-group">
          <rect x="150" y="110" width="20" height="20" 
                fill="rgba(16, 185, 129, 0.8)" stroke="rgba(16, 185, 129, 1)" stroke-width="2"/>
          <circle cx="160" cy="120" r="3" fill="rgba(16, 185, 129, 1)"/>
          <line x1="155" y1="115" x2="165" y2="125" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
          <line x1="165" y1="115" x2="155" y2="125" stroke="rgba(16, 185, 129, 1)" stroke-width="1"/>
        </g>
      </svg>
    </div>

    <!-- 상태 패널 -->
    <div class="status-panel">
      <h3>진단 상태</h3>
      <div class="status-grid">
        <div class="status-item">
          <div class="status-indicator" style="background-color: rgba(16, 185, 129, 0.8);"></div>
          <div class="status-info">
            <div class="component-name">코어</div>
            <div class="status-text">정상</div>
            <div class="status-value">85°C</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-indicator" style="background-color: rgba(245, 158, 11, 0.8);"></div>
          <div class="status-info">
            <div class="component-name">권선</div>
            <div class="status-text">주의</div>
            <div class="status-value">92°C</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-indicator" style="background-color: rgba(16, 185, 129, 0.8);"></div>
          <div class="status-info">
            <div class="component-name">부싱</div>
            <div class="status-text">정상</div>
            <div class="status-value">78°C</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-indicator" style="background-color: rgba(239, 68, 68, 0.8);"></div>
          <div class="status-info">
            <div class="component-name">탭 체인저</div>
            <div class="status-text">위험</div>
            <div class="status-value">105°C</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-indicator" style="background-color: rgba(16, 185, 129, 0.8);"></div>
          <div class="status-info">
            <div class="component-name">커패시터</div>
            <div class="status-text">정상</div>
            <div class="status-value">45°C</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-indicator" style="background-color: rgba(249, 115, 22, 0.8);"></div>
          <div class="status-info">
            <div class="component-name">부하</div>
            <div class="status-text">경고</div>
            <div class="status-value">88%</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-indicator" style="background-color: rgba(16, 185, 129, 0.8);"></div>
          <div class="status-info">
            <div class="component-name">스트레스</div>
            <div class="status-text">정상</div>
            <div class="status-value">0.8kV/mm</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-indicator" style="background-color: rgba(16, 185, 129, 0.8);"></div>
          <div class="status-info">
            <div class="component-name">지락 고장</div>
            <div class="status-text">정상</div>
            <div class="status-value">0.1mA</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-indicator" style="background-color: rgba(245, 158, 11, 0.8);"></div>
          <div class="status-info">
            <div class="component-name">불평형</div>
            <div class="status-text">주의</div>
            <div class="status-value">15%</div>
          </div>
        </div>
        <div class="status-item">
          <div class="status-indicator" style="background-color: rgba(16, 185, 129, 0.8);"></div>
          <div class="status-info">
            <div class="component-name">소음 및 진동</div>
            <div class="status-text">정상</div>
            <div class="status-value">65dB</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Transformer',
  mounted() {
    //console.log('실제 변압기 도면 렌더링 완료!');
  }
}
</script>

<style scoped>
.transformer-diagnosis {
  font-family: 'Arial', sans-serif;
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  background: #f5f5f5;
  border-radius: 8px;
}

.header h2 {
  color: #333;
  margin: 0;
}

.status-legend {
  display: flex;
  gap: 20px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.color-box {
  width: 20px;
  height: 20px;
  border-radius: 4px;
}

.diagram-container {
  background: white;
  border: 2px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  text-align: center;
}

.transformer-svg {
  border: 1px solid #ccc;
  background: #fafafa;
}

.title-label {
  font-size: 14px;
  fill: #374151;
  font-weight: bold;
}

.component-label {
  font-size: 10px;
  fill: #1f2937;
  font-weight: bold;
}

.component-rect, .component-circle, .component-ellipse {
  cursor: pointer;
  transition: all 0.3s ease;
}

.component-rect:hover, .component-circle:hover, .component-ellipse:hover {
  stroke-width: 4;
  filter: brightness(1.1);
}

.status-panel {
  background: white;
  border: 2px solid #ddd;
  border-radius: 8px;
  padding: 20px;
}

.status-panel h3 {
  margin-top: 0;
  color: #333;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 10px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.status-item:hover {
  background: #f9f9f9;
}

.status-indicator {
  width: 16px;
  height: 16px;
  border-radius: 50%;
}

.status-info {
  flex: 1;
}

.component-name {
  font-weight: bold;
  color: #333;
}

.status-text {
  font-size: 12px;
  color: #666;
}

.status-value {
  font-size: 14px;
  font-weight: bold;
  color: #333;
}

@media (max-width: 768px) {
  .header {
    flex-direction: column;
    gap: 10px;
  }
  
  .status-legend {
    flex-direction: column;
    gap: 10px;
  }
  
  .status-grid {
    grid-template-columns: 1fr;
  }
}
</style>