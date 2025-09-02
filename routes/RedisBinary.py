from pydantic import BaseModel, Field, create_model,validator
from typing import List, Dict, Any, Optional, Type
from datetime import datetime
import struct
from enum import IntEnum
import logging
from abc import ABC, abstractmethod
from enum import Enum

# # 로거 설정
# logging = logging.getlogging(__name__)


class RedisDataType(str, Enum):
    """Redis 데이터 타입"""
    STRING = "string"
    HASH = "hash"
    LIST = "list"
    SET = "set"
    ZSET = "zset"


class BinaryStructConfig:
    """바이너리 구조체 설정"""

    def __init__(
            self,
            name: str,
            struct_format: str,
            field_names: List[str],
            field_descriptions: Optional[Dict[str, str]] = None,
            timestamp_fields: Optional[List[str]] = None,
            scale_factors: Optional[Dict[str, float]] = None
    ):
        self.name = name
        self.struct_format = struct_format
        self.field_names = field_names
        self.field_descriptions = field_descriptions or {}
        self.timestamp_fields = timestamp_fields or []
        self.scale_factors = scale_factors or {}
        self.struct_size = struct.calcsize(struct_format)

        # Pydantic 모델 동적 생성
        self.model = self._create_pydantic_model()

    def _create_pydantic_model(self) -> Type[BaseModel]:
        """구조체에 맞는 Pydantic 모델 동적 생성"""
        field_definitions = {}

        for field_name in self.field_names:
            field_type = float  # 기본 타입
            description = self.field_descriptions.get(field_name, field_name)

            # 타임스탬프 필드는 datetime으로
            if field_name in self.timestamp_fields:
                field_type = datetime

            field_definitions[field_name] = (field_type, Field(..., description=description))

        return create_model(
            f"{self.name}Model",
            **field_definitions
        )


class BinaryParser(ABC):
    """바이너리 파서 추상 클래스"""

    @abstractmethod
    def parse(self, binary_data: bytes, config: BinaryStructConfig) -> Dict[str, Any]:
        """바이너리 데이터 파싱"""
        pass


class StandardBinaryParser(BinaryParser):
    """표준 struct 기반 바이너리 파서"""

    def parse(self, binary_data: bytes, config: BinaryStructConfig) -> Dict[str, Any]:
        """바이너리 데이터를 딕셔너리로 파싱"""
        if len(binary_data) < config.struct_size:
            raise ValueError(f"데이터 크기 부족: {len(binary_data)} < {config.struct_size}")

        try:
            # 구조체 언팩
            unpacked = struct.unpack(config.struct_format, binary_data[:config.struct_size])

            # 필드명과 값 매핑
            result = {}
            for i, field_name in enumerate(config.field_names):
                if i < len(unpacked):
                    value = unpacked[i]

                    # 타임스탬프 변환
                    if field_name in config.timestamp_fields and isinstance(value, (int, float)):
                        value = datetime.fromtimestamp(value)

                    # 스케일 팩터 적용
                    elif field_name in config.scale_factors and isinstance(value, (int, float)):
                        value = value * config.scale_factors[field_name]

                    result[field_name] = value

            return result

        except struct.error as e:
            raise ValueError(f"구조체 파싱 오류: {e}")


class NestedBinaryParser(BinaryParser):
    """중첩 구조체 파서 (MaxMin 같은 경우)"""

    def __init__(self, base_format: str, base_count: int, field_mapping: List[tuple]):
        """
        Args:
            base_format: 기본 구조체 포맷 (예: '<fIfI' for MAXMIN_DATA)
            base_count: 전체 기본 구조체 개수
            field_mapping: [(field_name, count), ...] 필드명과 개수 매핑
        """
        self.base_format = base_format
        self.base_size = struct.calcsize(base_format)
        self.base_count = base_count
        self.field_mapping = field_mapping
        self.total_size = self.base_size * base_count

    def parse(self, binary_data: bytes, config: BinaryStructConfig) -> Dict[str, Any]:
        """중첩 구조체 파싱"""
        if len(binary_data) < self.total_size:
            raise ValueError(f"데이터 크기 부족: {len(binary_data)} < {self.total_size}")

        result = {}
        offset = 0

        for field_name, count in self.field_mapping:
            if count == 1:
                # 단일 값
                values = struct.unpack_from(self.base_format, binary_data, offset)
                result[field_name] = self._parse_base_struct(values)
                offset += self.base_size
            else:
                # 배열 값
                array_values = []
                for _ in range(count):
                    values = struct.unpack_from(self.base_format, binary_data, offset)
                    array_values.append(self._parse_base_struct(values))
                    offset += self.base_size
                result[field_name] = array_values

        return result

    def _parse_base_struct(self, values: tuple) -> Dict[str, Any]:
        """기본 구조체 파싱 (오버라이드 가능)"""
        # 기본 구현: 튜플 그대로 반환
        return values


class BinaryDataProcessor:
    """범용 바이너리 데이터 처리기"""

    def __init__(self, redis_client, db_num: Optional[int] = None):
        """
        초기화

        Args:
            redis_client: Redis 클라이언트 (decode_responses=False)
        """
        self.redis_client = redis_client
        self.db_num = db_num
        self.configs: Dict[str, BinaryStructConfig] = {}
        self.parsers: Dict[str, BinaryParser] = {}
        self.redis_key_patterns: Dict[str, str] = {}  # config_name -> redis_key_pattern

        # 기본 파서 등록
        self.default_parser = StandardBinaryParser()

    def _select_db(self):
        """설정된 DB 선택"""
        if self.db_num is not None:
            self.redis_client.execute_command("SELECT", self.db_num)

    def register_config(
            self,
            config: BinaryStructConfig,
            redis_key_pattern: str,
            parser: Optional[BinaryParser] = None
    ):
        """
        구조체 설정 등록

        Args:
            config: 바이너리 구조체 설정
            redis_key_pattern: Redis 키 패턴 (예: "maxmin_{period}_ch{channel}")
            parser: 커스텀 파서 (없으면 기본 파서 사용)
        """
        self.configs[config.name] = config
        self.redis_key_patterns[config.name] = redis_key_pattern
        self.parsers[config.name] = parser or self.default_parser

        logging.info(f"바이너리 설정 등록: {config.name}")

    def get_redis_data(
            self,
            key: str,
            data_type: RedisDataType = RedisDataType.STRING,
            field: Optional[str] = None,
            index: Optional[int] = None
    ) -> Optional[bytes]:
        """
        Redis에서 바이너리 데이터 가져오기

        Args:
            key: Redis 키
            data_type: Redis 데이터 타입
            field: Hash 필드명 (Hash 타입인 경우)
            index: List 인덱스 (List 타입인 경우)

        Returns:
            바이너리 데이터 또는 None
        """
        try:
            #            self._select_db()

            if data_type == RedisDataType.STRING:
                return self.redis_client.get(key)
            elif data_type == RedisDataType.HASH and field:
                return self.redis_client.hget(key, field)
            elif data_type == RedisDataType.LIST and index is not None:
                return self.redis_client.lindex(key, index)
            else:
                logging.warning(f"지원하지 않는 조회 방식: {data_type}")
                return None
        except Exception as e:
            logging.error(f"Redis 데이터 조회 오류: {e}")
            return None

    def parse_data(
            self,
            config_name: str,
            binary_data: bytes
    ) -> Dict[str, Any]:
        """
        바이너리 데이터 파싱

        Args:
            config_name: 설정 이름
            binary_data: 바이너리 데이터

        Returns:
            파싱된 데이터 딕셔너리
        """
        if config_name not in self.configs:
            raise ValueError(f"설정을 찾을 수 없습니다: {config_name}")

        config = self.configs[config_name]
        parser = self.parsers[config_name]

        return parser.parse(binary_data, config)

    def get_and_parse(
            self,
            config_name: str,
            key: str,
            data_type: RedisDataType = RedisDataType.STRING,
            field: Optional[str] = None,
            index: Optional[int] = None
    ) -> Optional[Dict[str, Any]]:
        """
        Redis에서 데이터를 가져와서 파싱

        Args:
            config_name: 설정 이름
            key: Redis 키
            data_type: Redis 데이터 타입
            field: Hash 필드명
            index: List 인덱스

        Returns:
            파싱된 데이터 또는 None
        """
        binary_data = self.get_redis_data(key, data_type, field, index)

        if not binary_data:
            return None

        try:
            return self.parse_data(config_name, binary_data)
        except Exception as e:
            logging.error(f"데이터 파싱 오류: {e}")
            return None

    def get_config(self, config_name: str) -> Optional[BinaryStructConfig]:
        """설정 조회"""
        return self.configs.get(config_name)

    def list_configs(self) -> List[str]:
        """등록된 설정 이름 목록"""
        return list(self.configs.keys())

    def parse_custom(
            self,
            struct_format: str,
            binary_data: bytes,
            field_names: Optional[List[str]] = None
    ) -> Dict[str, Any]:
        """
        커스텀 구조체 파싱 (디버깅/테스트용)

        Args:
            struct_format: 구조체 포맷
            binary_data: 바이너리 데이터
            field_names: 필드 이름 목록

        Returns:
            파싱된 데이터
        """
        try:
            unpacked = struct.unpack(struct_format, binary_data)

            if field_names and len(field_names) == len(unpacked):
                return dict(zip(field_names, unpacked))
            else:
                return {"values": list(unpacked)}

        except struct.error as e:
            raise ValueError(f"구조체 파싱 오류: {e}")


# MaxMin 데이터용 커스텀 파서
class MaxMinDataParser(NestedBinaryParser):
    """MaxMin 전용 파서"""

    def __init__(self):
        # MAXMIN_DATA 구조체: float max, uint32 max_ts, float min, uint32 min_ts
        super().__init__(
            base_format='<fIfI',
            base_count=21,  # MaxMin1Sec와 MaxMin15Min 모두 21개
            field_mapping=[]  # 서브클래스에서 설정
        )

    def _parse_base_struct(self, values: tuple) -> Dict[str, Any]:
        """MAXMIN_DATA 구조체 파싱"""
        return {
            "max": values[0],
            "max_timestamp": datetime.fromtimestamp(values[1]),
            "min": values[2],
            "min_timestamp": datetime.fromtimestamp(values[3])
        }


class MaxMin1SecParser(MaxMinDataParser):
    """1초 주기 MaxMin 파서"""

    def __init__(self):
        super().__init__()
        self.field_mapping = [
            ("freq", 1),
            ("temp", 1),
            ("u", 4),
            ("upp", 4),
            ("i", 4),
            ("itot", 1),
            ("in", 1),
            ("isum", 1),
            ("pf", 4)
        ]


class MaxMin15MinParser(MaxMinDataParser):
    """15분 주기 MaxMin 파서"""

    def __init__(self):
        super().__init__()
        self.field_mapping = [
            ("p", 4),
            ("q", 4),
            ("s", 4),
            ("thd_u", 3),
            ("thd_upp", 3),
            ("thd_i", 3)
        ]


class CmdType(IntEnum):
    CMD_CLR = 0
    CMD_REBOOT = 1


class ItemType(IntEnum):
    ITEM_DEMAND = 0
    ITEM_MAXMIN = 1
    ITEM_ENERGY = 2
    ITEM_ALARM = 3
    ITEM_EVENT = 4
    ITEM_REBOOT = 5
    ITEM_RUNHOUR = 6

# Pydantic 모델
class Command(BaseModel):
    type: int
    cmd: CmdType
    item: ItemType

    @validator('type')
    def validate_type(cls, v):
        if v < 0:
            raise ValueError('type must be non-negative')
        return v

    class Config:
        use_enum_values = True  # enum을 정수값으로 직렬화


class WaveformParser(BinaryParser):
    """WAVEFORM_L16 구조체 전용 파서"""

    def __init__(self):
        # WAVEFORM_L16 구조체 크기 계산
        # int16_t U[3][160]: 3 * 160 * 2 = 960 bytes
        # int16_t Upp[3][160]: 3 * 160 * 2 = 960 bytes
        # int16_t I[3][160]: 3 * 160 * 2 = 960 bytes
        # float vscale: 4 bytes
        # float iscale: 4 bytes
        # int16_t r[56]: 56 * 2 = 112 bytes
        # Total: 960 + 960 + 960 + 4 + 4 + 112 = 3000 bytes
        self.struct_size = 3000
        self.waveform_size = 160
        self.phase_count = 3

    def parse(self, binary_data: bytes, config: BinaryStructConfig = None) -> Dict[str, Any]:
        """WAVEFORM_L16 바이너리 데이터 파싱"""
        if len(binary_data) < self.struct_size:
            raise ValueError(f"데이터 크기 부족: {len(binary_data)} < {self.struct_size}")

        try:
            offset = 0
            result = {}

            # U[3][160] 파싱
            u_data = []
            for phase in range(self.phase_count):
                phase_data = []
                for i in range(self.waveform_size):
                    value = struct.unpack_from('<h', binary_data, offset)[0]  # int16_t = 'h'
                    phase_data.append(value)
                    offset += 2
                u_data.append(phase_data)

            # Upp[3][160] 파싱
            upp_data = []
            for phase in range(self.phase_count):
                phase_data = []
                for i in range(self.waveform_size):
                    value = struct.unpack_from('<h', binary_data, offset)[0]
                    phase_data.append(value)
                    offset += 2
                upp_data.append(phase_data)

            # I[3][160] 파싱
            i_data = []
            for phase in range(self.phase_count):
                phase_data = []
                for i in range(self.waveform_size):
                    value = struct.unpack_from('<h', binary_data, offset)[0]
                    phase_data.append(value)
                    offset += 2
                i_data.append(phase_data)

            # vscale, iscale 파싱
            vscale = struct.unpack_from('<f', binary_data, offset)[0]  # float = 'f'
            offset += 4
            iscale = struct.unpack_from('<f', binary_data, offset)[0]
            offset += 4

            # r[56] 파싱 (예약 영역, 필요시 사용)
            reserved = []
            for i in range(56):
                value = struct.unpack_from('<h', binary_data, offset)[0]
                reserved.append(value)
                offset += 2

            # 결과 구성
            result = {
                # Phase Voltage
                "Wave Form V1": u_data[0],
                "Wave Form V2": u_data[1],
                "Wave Form V3": u_data[2],

                # Line Voltage
                "Wave Form V12": upp_data[0],
                "Wave Form V23": upp_data[1],
                "Wave Form V31": upp_data[2],

                # Current
                "Wave Form I1": i_data[0],
                "Wave Form I2": i_data[1],
                "Wave Form I3": i_data[2],

                # Scale factors
                "vscale": vscale,
                "iscale": iscale,

                # Reserved (옵션)
                # "reserved": reserved
            }

            return result

        except struct.error as e:
            raise ValueError(f"Waveform 구조체 파싱 오류: {e}")


# Waveform 설정 등록 함수
def register_waveform_config(processor: BinaryDataProcessor):
    """Waveform 구조체 설정 등록"""

    # Waveform 설정 생성 (구조체 포맷은 파서에서 직접 처리하므로 간단히)
    waveform_config = BinaryStructConfig(
        name="waveform",
        struct_format="<3000s",  # 전체를 바이트 배열로
        field_names=["raw_data"],
        field_descriptions={
            "raw_data": "Waveform raw binary data"
        }
    )

    # 커스텀 파서와 함께 등록
    processor.register_config(
        config=waveform_config,
        redis_key_pattern="waveform_ch{channel}",  # 예시 패턴
        parser=WaveformParser()
    )

