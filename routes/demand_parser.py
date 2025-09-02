from datetime import datetime
from typing import Dict, Any
import struct
from .RedisBinary import NestedBinaryParser, BinaryStructConfig


class DemandParser(NestedBinaryParser):
    """Demand 데이터 전용 파서"""

    def __init__(self):
        # DemandRedis 구조체 전체 크기: 108 bytes
        # MAX_DEMAND 구조체: uint32 + float = 8 bytes
        # 총 8개의 MAX_DEMAND + 나머지 float 값들

        # 구조를 정확히 계산:
        # - 8개의 MAX_DEMAND (64 bytes)
        # - 1개의 uint32 (4 bytes)
        # - 10개의 float (40 bytes)
        # 총합: 108 bytes

        # 기본 구조체는 사용하지 않고 직접 파싱
        super().__init__(
            base_format='',  # 사용하지 않음
            base_count=0,  # 사용하지 않음
            field_mapping=[]  # 사용하지 않음
        )

        self.total_size = 124

    def parse(self, binary_data: bytes, config: BinaryStructConfig) -> Dict[str, Any]:
        """Demand 데이터 파싱"""
        if len(binary_data) != self.total_size:
            raise ValueError(f"데이터 크기 불일치: 예상 {self.total_size}바이트, 실제 {len(binary_data)}바이트")

        result = {}
        offset = 0

        try:
            # 1. Max Demand 파싱
            # MD_P[2]
            result["MD_P"] = []
            for i in range(2):
                mdTime, value = struct.unpack_from('<If', binary_data, offset)
                result["MD_P"].append({
                    "time": datetime.fromtimestamp(mdTime),
                    "timestamp": mdTime,
                    "value": value
                })
                offset += 8

            # MD_Q[2]
            result["MD_Q"] = []
            for i in range(2):
                mdTime, value = struct.unpack_from('<If', binary_data, offset)
                result["MD_Q"].append({
                    "time": datetime.fromtimestamp(mdTime),
                    "timestamp": mdTime,
                    "value": value
                })
                offset += 8

            # MD_S (단일)
            mdTime, value = struct.unpack_from('<If', binary_data, offset)
            result["MD_S"] = {
                "time": datetime.fromtimestamp(mdTime),
                "timestamp": mdTime,
                "value": value
            }
            offset += 8

            # MD_I[3]
            result["MD_I"] = []
            for i in range(3):
                mdTime, value = struct.unpack_from('<If', binary_data, offset)
                result["MD_I"].append({
                    "time": datetime.fromtimestamp(mdTime),
                    "timestamp": mdTime,
                    "value": value
                })
                offset += 8

            # 2. Dynamic Demand 파싱
            # ddTime
            ddTime = struct.unpack_from('<I', binary_data, offset)[0]
            result["ddTime"] = datetime.fromtimestamp(ddTime)
            result["ddTimestamp"] = ddTime
            offset += 4

            # DD_P[2]
            result["DD_P"] = list(struct.unpack_from('<2f', binary_data, offset))
            offset += 8

            # DD_Q[2]
            result["DD_Q"] = list(struct.unpack_from('<2f', binary_data, offset))
            offset += 8

            # DD_S
            result["DD_S"] = struct.unpack_from('<f', binary_data, offset)[0]
            offset += 4

            # DD_I[3]
            result["DD_I"] = list(struct.unpack_from('<3f', binary_data, offset))
            offset += 12

            # 3. Present Demand 파싱
            # CD_P[2]
            result["CD_P"] = list(struct.unpack_from('<2f', binary_data, offset))
            offset += 8

            # CD_Q[2]
            result["CD_Q"] = list(struct.unpack_from('<2f', binary_data, offset))
            offset += 8

            # CD_S
            result["CD_S"] = struct.unpack_from('<f', binary_data, offset)[0]
            offset += 4

            # 4. Predict Demand 파싱
            # PD_P
            result["PD_P"] = struct.unpack_from('<f', binary_data, offset)[0]
            offset += 4

            return result

        except struct.error as e:
            raise ValueError(f"Demand 데이터 파싱 오류: {e}")