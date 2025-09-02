#!/bin/bash

IP_FILE="ip.txt"
KEY_FILE="key.pem"
CERT_FILE="cert.pem"
CRT_FILE="cert.crt"
CN="sv500.local"

# IP 주소 목록 확인
if [ ! -f "$IP_FILE" ]; then
    echo "❌ $IP_FILE 파일이 존재하지 않습니다."
    exit 1
fi

# subjectAltName 문자열 생성
IP_LIST=$(awk '{print "IP:" $0}' "$IP_FILE" | paste -sd, -)
SAN="DNS:$CN,$IP_LIST"

echo "📥 SAN 구성: $SAN"

# 기존 파일 삭제
rm -f "$KEY_FILE" "$CERT_FILE" "$CRT_FILE"

# 인증서 생성
openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout "$KEY_FILE" -out "$CERT_FILE" -days 365 \
  -subj "/C=KR/ST=Seoul/L=Gangnam/O=MyCompany/OU=Dev/CN=$CN" \
  -addext "subjectAltName=$SAN"

# cert.crt 복사본 생성
cp "$CERT_FILE" "$CRT_FILE"

echo "✅ 인증서 생성 완료:"
echo " - 개인키: $KEY_FILE"
echo " - 인증서: $CERT_FILE"
echo " - 복사본: $CRT_FILE"