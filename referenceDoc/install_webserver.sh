#!/bin/sh

# ===== 설정 =====
APP_DIR=/home/root/webserver
VENV_DIR=$APP_DIR/venv
MAIN_FILE=main_fast.py
INFLUXDB_VERSION=2.7.11
PACKAGE_JSON=$APP_DIR/package_pip.json

# 1. 디렉토리 생성
mkdir -p $APP_DIR
cd $APP_DIR

# 2. Python 가상환경 생성
python3 -m venv $VENV_DIR
. $VENV_DIR/bin/activate

# 3. pip 최신화
pip install --upgrade pip

# 4. packages.json 기반 의존성 설치
echo ">>> packages.json을 기반으로 의존성 설치 중..."
if [ -f "$PACKAGE_JSON" ]; then
    DEPS=$(python3 -c "
import json
with open('$PACKAGE_JSON') as f:
    data = json.load(f)
print(' '.join(data.get('dependencies', [])))
")
    pip install $DEPS
else
    echo "⚠️ packages.json 파일이 없습니다: $PACKAGE_JSON"
    exit 1
fi

# 5. FastAPI 앱 파일 복사 (수정 필요)
cp /path/to/your/$MAIN_FILE $APP_DIR

# 6. FastAPI systemd 등록
cat <<EOF > /etc/systemd/system/webserver.service
[Unit]
Description=FastAPI Web Server
After=network.target

[Service]
ExecStart=$VENV_DIR/bin/uvicorn $MAIN_FILE:app --host 0.0.0.0 --port 4000
WorkingDirectory=$APP_DIR
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

# 7. InfluxDB 2.7.11 설치
cd /tmp
wget https://dl.influxdata.com/influxdb/releases/influxdb2-${INFLUXDB_VERSION}-linux-armv7.tar.gz
tar -xzf influxdb2-${INFLUXDB_VERSION}-linux-armv7.tar.gz
cp influxdb2-${INFLUXDB_VERSION}-linux-armv7/influxd /usr/local/bin/
cp influxdb2-${INFLUXDB_VERSION}-linux-armv7/influx /usr/local/bin/

# 8. InfluxDB 디렉토리 구성
mkdir -p /etc/influxdb2
mkdir -p /var/lib/influxdb2
mkdir -p /var/log/influxdb

# 9. InfluxDB systemd 서비스 등록 (setup은 웹서버에서 진행)
cat <<EOF > /etc/systemd/system/influxdb.service
[Unit]
Description=InfluxDB 2.7.11
After=network.target

[Service]
ExecStart=/usr/local/bin/influxd --bolt-path /var/lib/influxdb2/influxd.bolt --engine-path /var/lib/influxdb2/engine
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

# 10. 서비스 등록 및 시작
systemctl daemon-reexec
systemctl daemon-reload
systemctl enable influxdb
systemctl enable webserver
systemctl start influxdb
systemctl start webserver

echo "✅ 설치 완료: FastAPI + InfluxDB 2.7.11. 의존성은 packages.json 기준 설치됨"
