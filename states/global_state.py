import base64
import logging
import platform, asyncio
from influxdb_client import InfluxDBClient
from influxdb_client.client.write_api import ASYNCHRONOUS
import ujson as json
import os, redis, httpx
from pathlib import Path
from typing import Optional
from Crypto.Util.Padding import unpad, pad
from Crypto.Cipher import AES
from hashlib import sha256


base_dir = Path(__file__).resolve().parent
SETTING_FOLDER = base_dir.parent.parent / "config"  # ⬅️ 두 단계 상위로
INIT_PATH = file_path = os.path.join(SETTING_FOLDER, 'influx.json')
ADMIN_PATH = os.path.join(SETTING_FOLDER, 'admin_secret.enc')

AES_KEY = b'ntekSystem_20250721_mypark_caner'  # 16바이트 (AES-128)

class OsSpec:
    def __init__(self, mode):
        self.os = self.get_os_info()
        if self.os == 'Windows':
            self.influxip = '192.168.1.91'
            self.restip = '192.168.1.24'
            self.logpath = '.'
            self.redisip = '192.168.1.92'
        else:
            self.redisip = '127.0.0.1'
            if mode > 0:
                self.influxip = '192.168.1.91'
                self.restip = '192.168.1.24'
                if mode == 1:
                    self.logpath = "/home/mypark/logs"
                else:
                    self.logpath = "/home/ntek/logs"
            else:
                self.influxip = '127.0.0.1'
                self.restip = '127.0.0.1'
                self.logpath = '/home/root/logs'


    def get_os_info(self):
        os_type = platform.system()

        if os_type == 'Windows':
            return "Windows"
        elif os_type == 'Linux':
            return "Linux"
        elif os_type == 'Darwin':
            return "macOS"
        else:
            return f"{os_type}"

    def get_path_separator(self):
        """OS별 경로 구분자 반환"""
        if self.os == "Windows":
            return "\\"
        else:
            return "/"

os_spec = OsSpec(0)

class AESEncDec:
    def __init__(self):
        self.iv = sha256(AES_KEY).digest()[:16]
        try:
            with open(ADMIN_PATH, 'rb') as f:
                raw = base64.b64decode(f.read())
        except Exception as e:
            logging.info(f"❌ {ADMIN_PATH} Admin Password Encryption File Error: {e}")
        self.adminIv = raw[:16]
        self.encrypted = raw[16:]  # ← 암호문은 따로 저장해둬야 함

    def checkAdmin(self, context):
        cipher = AES.new(AES_KEY, AES.MODE_CBC, self.adminIv)
        decrypted = unpad(cipher.decrypt(self.encrypted), AES.block_size).decode('utf-8')
        if decrypted == context:
            return True
        else:
            return False


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
        if os.path.exists(INIT_PATH):
            try:
                with open(INIT_PATH, "r", encoding="utf-8") as f:
                    influxSetting = json.load(f)
                #token = self.decrypt(influxSetting["token"])
                # print(token)
                #encoded_ciphertext = self.encrypt(token)
                resultDict = {"result":True, "cipher": influxSetting["token"], "org":influxSetting["org"]}
            except Exception as e:
                logging.error(f"❌ {INIT_PATH} Influxdb Setup File Error: {e}")
                resultDict = {"result":False}
        else:
            resultDict = {"result": False}
        return resultDict

aesState = AESEncDec()

class InfluxGlobalState:
    def __init__(self):
        self._client: Optional[InfluxDBClient] = None
        self._write_api = None
        self._query_api = None
        self._write_async_api = None
        self._error: Optional[str] = None

    def _connect(self):
        try:
            file_path = os.path.join(SETTING_FOLDER, 'influx.json')
            with open(file_path) as f:
                cfg = json.load(f)

            self._client = InfluxDBClient(
                url=cfg.get("url", f"http://{os_spec.influxip}:8086"),
                token=aesState.decrypt(cfg["token"]),
                org=cfg["org"],
                timeout=30_000,  # 30초로 증가 (밀리초 단위)
                enable_gzip=True  # 대용량 데이터 압축 전송
            )
            self._write_api = self._client.write_api()
            self._query_api = self._client.query_api()
            self._write_async_api = self._client.write_api(write_options=ASYNCHRONOUS)
            self._error = None

        except Exception as e:
            self._client = None
            self._error = f"[Influx Init Error] {str(e)}"
            logging.error(f"❌ Influx Connect Error: {e}")

    @property
    def client(self):
        if not self._client:
            self._connect()
        return self._client

    @property
    def write_api(self):
        if not self._write_api:
            self._connect()
        return self._write_api

    @property
    def query_api(self):
        if not self._query_api:
            self._connect()
        return self._query_api

    @property
    def write_async_api(self):
        if not self._write_async_api:
            self._connect()
        return self._write_async_api

    @property
    def error(self):
        return self._error

    def close(self):
        if self._client:
            self._client.close()

influx_state = InfluxGlobalState()


class RedisGlobalState:
    client: Optional[redis.Redis] = None  # 텍스트용 (기존)
    binary_client: Optional[redis.Redis] = None  # 바이너리용 (추가)
    processor: Optional['BinaryDataProcessor'] = None
    handler: Optional['MaxMinDataHandler'] = None  # handler도 추가
    error: Optional[str] = None


redis_state = RedisGlobalState()


def init_redis():
    try:
        # 텍스트 데이터용 클라이언트 (기존)
        pool = redis.ConnectionPool(
            host=f"{os_spec.redisip}",
            port=6379,
            db=0,
            decode_responses=True  # 텍스트 자동 디코딩
        )
        redis_state.client = redis.Redis(connection_pool=pool)

        # 바이너리 데이터용 클라이언트 (추가)
        binary_pool = redis.ConnectionPool(
            host='127.0.0.1',
            port=6379,
            db=1,
            decode_responses=False  # 바이너리 모드
        )
        redis_state.binary_client = redis.Redis(connection_pool=binary_pool)

        redis_state.error = None

    except Exception as e:
        redis_state.client = None
        redis_state.binary_client = None
        redis_state.error = f"[Redis Init Error] {str(e)}"
        logging.error(f"❌ Redis Connect Error: {e}")


# class HTTPGlobalState:
#     def __init__(self):
#         self._client: Optional[httpx.AsyncClient] = None
#         self._lock: Optional[asyncio.Lock] = None
#         self._error: Optional[str] = None
#         self._config = {
#             "base_url": f"http://{os_spec.restip}:5000/api/",
#             "timeout": 30.0,
#             "connect_timeout": 5.0,
#             "max_connections": 5,  # ← 줄임
#             "max_keepalive_connections": 2  # ← 줄임
#         }
#
#     async def initialize(self):
#         """비동기 초기화"""
#         if self._lock is None:
#             self._lock = asyncio.Lock()
#
#         if not self._client:
#             async with self._lock:
#                 if not self._client:
#                     try:
#                         limits = httpx.Limits(
#                             max_connections=self._config["max_connections"],
#                             max_keepalive_connections=self._config["max_keepalive_connections"],
#                             keepalive_expiry=30.0
#                         )
#
#                         timeout = httpx.Timeout(
#                             timeout=self._config["timeout"],
#                             connect=self._config["connect_timeout"],
#                             read=30.0,
#                             write=5.0
#                         )
#
#                         self._client = httpx.AsyncClient(
#                             base_url=self._config["base_url"],
#                             limits=limits,
#                             timeout=timeout,
#                             http2=False
#                         )
#                         self._error = None
#                         logging.info(f"✅ Async HTTP Client initialized: {self._config['base_url']}")
#                     except Exception as e:
#                         self._error = f"[HTTP Init Error] {str(e)}"
#                         logging.error(f"❌ HTTP Client Init Error: {e}")
#                         raise
#
#     @property
#     def client(self):
#         """기존 코드 호환용 - 직접 client 접근"""
#         if not self._client:
#             raise RuntimeError("HTTP Client not initialized! Call 'await http_state.initialize()' first")
#         return self._client
#
#     async def post(self, endpoint: str, data: dict):
#         if not self._client:
#             await self.initialize()  # 자동 초기화
#         try:
#             response = await self._client.post(endpoint, json=data)
#             response.raise_for_status()
#             return response.json()
#         except httpx.HTTPError as e:
#             logging.error(f"HTTP POST Error: {e}")
#             raise
#
#     async def get(self, endpoint: str, params: dict = None):
#         if not self._client:
#             await self.initialize()  # 자동 초기화
#         try:
#             response = await self._client.get(endpoint, params=params)
#             response.raise_for_status()
#             return response.json()
#         except httpx.HTTPError as e:
#             logging.error(f"HTTP GET Error: {e}")
#             raise
#
#     async def close(self):
#         if self._client:
#             await self._client.aclose()
#             self._client = None
#             self._lock = None
#             logging.info("✅ Async HTTP Client closed")
#
#
# http_state = HTTPGlobalState()


async def cleanup_global_resources():
    """모든 전역 리소스 정리"""
    if influx_state.client:
        influx_state.client.close()
    if redis_state.client:
        redis_state.client.close()
    # await http_state.close()  # ← 비동기로 변경