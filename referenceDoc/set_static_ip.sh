#!/bin/bash

# === 설정 값 ===
IFACE="ens33"           # 네트워크 인터페이스 이름 (예: eth0, ens33 등)
STATIC_IP="192.168.1.93/24"
GATEWAY="192.168.1.1"
DNS1="168.126.63.1"
DNS2="1.1.1.1"
NETPLAN_FILE="/etc/netplan/01-netcfg.yaml"

# === 백업 및 설정 ===
echo "🔧 기존 netplan 설정 백업..."
cp $NETPLAN_FILE ${NETPLAN_FILE}.bak

echo "🔧 고정 IP 설정 적용 중..."
cat <<EOF > $NETPLAN_FILE
network:
  version: 2
  ethernets:
    ${IFACE}:
      dhcp4: no
      addresses: [${STATIC_IP}]
      gateway4: ${GATEWAY}
      nameservers:
        addresses: [${DNS1}, ${DNS2}]
EOF

echo "✅ 설정 완료. netplan 적용 중..."
netplan apply

echo "📡 현재 IP:"
ip addr show $IFACE
