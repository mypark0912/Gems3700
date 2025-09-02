import base64
import logging

from influxdb_client import InfluxDBClient
from influxdb_client.client.write_api import ASYNCHRONOUS
import json
import os, redis
from pathlib import Path
from typing import Optional
from Crypto.Util.Padding import unpad, pad
from Crypto.Cipher import AES
from hashlib import sha256

# base_dir = Path(__file__).resolve().parent
# SETTING_FOLDER = base_dir.parent.parent / "config"  # ⬅️ 두 단계 상위로
# INIT_PATH = file_path = os.path.join(SETTING_FOLDER, 'influx.json')
# ADMIN_PATH = os.path.join(SETTING_FOLDER, 'admin_secret.enc')

AES_KEY = b'ntekSystem_202504_mypark'  # 16바이트 (AES-128)

class AESEncDec:
    def __init__(self, path):
        self.iv = sha256(AES_KEY).digest()[:16]
        self.path = path
        # try:
        #     with open(ADMIN_PATH, 'rb') as f:
        #         raw = base64.b64decode(f.read())
        # except Exception as e:
        #     logging.info(f"❌ {ADMIN_PATH} Admin Password Encryption File Error: {e}")
        # self.adminIv = raw[:16]
        # self.encrypted = raw[16:]  # ← 암호문은 따로 저장해둬야 함


    # def checkAdmin(self, context):
    #     cipher = AES.new(AES_KEY, AES.MODE_CBC, self.adminIv)
    #     decrypted = unpad(cipher.decrypt(self.encrypted), AES.block_size).decode('utf-8')
    #     if decrypted == context:
    #         return True
    #     else:
    #         return False


    def encrypt(self, context):
        cipher = AES.new(AES_KEY, AES.MODE_CBC, self.iv)
        ciphertext = cipher.encrypt(pad(context.encode('utf-8'), AES.block_size))
        encoded_ciphertext = base64.b64encode(ciphertext).decode('utf-8')
        return encoded_ciphertext


    def decrypt(self, context):
        ciphertext_dec = base64.b64decode(context)
        cipher_dec = AES.new(AES_KEY, AES.MODE_CBC, self.iv)
        plaintext_dec = unpad(cipher_dec.decrypt(ciphertext_dec), AES.block_size).decode('utf-8')
        return plaintext_dec


    def getInflux(self):
        if os.path.exists(self.path):
            try:
                file_path = os.path.join(self.path, 'influx.json')
                with open(file_path, "r", encoding="utf-8") as f:
                    influxSetting = json.load(f)
                token = self.decrypt(influxSetting["token"])
                print(token)
                # encoded_ciphertext = self.encrypt(token)
                resultDict = {"result": True, "cipher": influxSetting["token"], "org": influxSetting["org"]}
            except Exception as e:
                logging.info(f"❌ {file_path} Influxdb Setup File Error: {e}")
                resultDict = {"result": False}
        else:
            resultDict = {"result": False}
        return resultDict

aes = AESEncDec(r"C:\config")
aes.getInflux()