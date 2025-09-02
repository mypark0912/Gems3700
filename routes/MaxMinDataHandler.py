from typing import Dict, List, Any, Optional
from datetime import datetime
from .RedisBinary import (
    BinaryDataProcessor,
    BinaryStructConfig,
    MaxMin1SecParser,
    MaxMin15MinParser
)


class MaxMinDataHandler:
    """새로운 C 구조체 형식의 MaxMin 데이터 처리기"""

    def __init__(self, processor: BinaryDataProcessor):
        """
        Args:
            processor: BinaryDataProcessor 인스턴스
        """
        self.processor = processor

        # MaxMin 설정이 등록되어 있지 않으면 등록
        if "maxmin_1sec" not in processor.list_configs():
            self._register_configs()

    def _register_configs(self):
        """MaxMin 설정 등록"""
        # 1초 데이터 설정
        maxmin_1sec_config = BinaryStructConfig(
            name="maxmin_1sec",
            struct_format='<' + 'fIfI' * 21,  # MAXMIN_DATA 21개
            field_names=["data"]  # 파서가 처리
        )

        self.processor.register_config(
            config=maxmin_1sec_config,
            redis_key_pattern="maxmin_1sec_ch{channel_id}",
            parser=MaxMin1SecParser()
        )

        # 15분 데이터 설정
        maxmin_15min_config = BinaryStructConfig(
            name="maxmin_15min",
            struct_format='<' + 'fIfI' * 21,
            field_names=["data"]
        )

        self.processor.register_config(
            config=maxmin_15min_config,
            redis_key_pattern="maxmin_15min_ch{channel_id}",
            parser=MaxMin15MinParser()
        )

    def get_maxmin_data(self, channel: str, period: str = "1sec") -> Optional[Dict[str, Any]]:
        """
        MaxMin 데이터 가져오기

        Args:
            channel: 채널명 (예: "Main", "Sub")
            period: "1sec" 또는 "15min"

        Returns:
            파싱된 MaxMin 데이터
        """
        # 채널명을 ID로 변환 (예시)
        channel_mapping = {"Main": 1, "Sub": 2}
        channel_id = channel_mapping.get(channel, 1)

        # Redis 키 생성
        redis_key = f"maxmin_{period}_ch{channel_id}"

        # 데이터 파싱
        return self.processor.get_and_parse(
            config_name=f"maxmin_{period}",
            key=redis_key
        )

    def convert_to_legacy_format(self, parsed_data: Dict[str, Any]) -> Dict[str, Any]:
        """
        새 구조체 형식을 기존 형식으로 변환

        기존: U1_max, U1_min, U1_maxTime, U1_minTime
        새: u[0] = {max: ..., max_timestamp: ..., min: ..., min_timestamp: ...}
        """
        legacy_format = {}

        # 1초 데이터인 경우
        if "freq" in parsed_data:
            # 주파수
            if "freq" in parsed_data:
                legacy_format["Freq_max"] = parsed_data["freq"]["max"]
                legacy_format["Freq_min"] = parsed_data["freq"]["min"]
                legacy_format["Freq_maxTime"] = parsed_data["freq"]["max_timestamp"]
                legacy_format["Freq_minTime"] = parsed_data["freq"]["min_timestamp"]

            # 온도
            if "temp" in parsed_data:
                legacy_format["Temp_max"] = parsed_data["temp"]["max"]
                legacy_format["Temp_min"] = parsed_data["temp"]["min"]
                legacy_format["Temp_maxTime"] = parsed_data["temp"]["max_timestamp"]
                legacy_format["Temp_minTime"] = parsed_data["temp"]["min_timestamp"]

            # 선간전압 U
            if "u" in parsed_data:
                for i, u_data in enumerate(parsed_data["u"]):
                    idx = i + 1
                    legacy_format[f"U{idx}_max"] = u_data["max"]
                    legacy_format[f"U{idx}_min"] = u_data["min"]
                    legacy_format[f"U{idx}_maxTime"] = u_data["max_timestamp"]
                    legacy_format[f"U{idx}_minTime"] = u_data["min_timestamp"]

            # 상전압 Upp
            if "upp" in parsed_data:
                for i, upp_data in enumerate(parsed_data["upp"]):
                    idx = i + 1
                    legacy_format[f"Upp{idx}_max"] = upp_data["max"]
                    legacy_format[f"Upp{idx}_min"] = upp_data["min"]
                    legacy_format[f"Upp{idx}_maxTime"] = upp_data["max_timestamp"]
                    legacy_format[f"Upp{idx}_minTime"] = upp_data["min_timestamp"]

            # 전류 I
            if "i" in parsed_data:
                for i, i_data in enumerate(parsed_data["i"]):
                    idx = i + 1
                    legacy_format[f"I{idx}_max"] = i_data["max"]
                    legacy_format[f"I{idx}_min"] = i_data["min"]
                    legacy_format[f"I{idx}_maxTime"] = i_data["max_timestamp"]
                    legacy_format[f"I{idx}_minTime"] = i_data["min_timestamp"]

            # 전체 전류
            if "itot" in parsed_data:
                legacy_format["Itot_max"] = parsed_data["itot"]["max"]
                legacy_format["Itot_min"] = parsed_data["itot"]["min"]
                legacy_format["Itot_maxTime"] = parsed_data["itot"]["max_timestamp"]
                legacy_format["Itot_minTime"] = parsed_data["itot"]["min_timestamp"]

            # 역률 PF
            if "pf" in parsed_data:
                for i, pf_data in enumerate(parsed_data["pf"]):
                    idx = i + 1
                    legacy_format[f"PF{idx}_max"] = pf_data["max"]
                    legacy_format[f"PF{idx}_min"] = pf_data["min"]
                    legacy_format[f"PF{idx}_maxTime"] = pf_data["max_timestamp"]
                    legacy_format[f"PF{idx}_minTime"] = pf_data["min_timestamp"]

        # 15분 데이터인 경우
        elif "p" in parsed_data:
            # 유효전력 P
            if "p" in parsed_data:
                for i, p_data in enumerate(parsed_data["p"]):
                    idx = i + 1
                    legacy_format[f"P{idx}_max"] = p_data["max"]
                    legacy_format[f"P{idx}_min"] = p_data["min"]
                    legacy_format[f"P{idx}_maxTime"] = p_data["max_timestamp"]
                    legacy_format[f"P{idx}_minTime"] = p_data["min_timestamp"]

            # 무효전력 Q
            if "q" in parsed_data:
                for i, q_data in enumerate(parsed_data["q"]):
                    idx = i + 1
                    legacy_format[f"Q{idx}_max"] = q_data["max"]
                    legacy_format[f"Q{idx}_min"] = q_data["min"]
                    legacy_format[f"Q{idx}_maxTime"] = q_data["max_timestamp"]
                    legacy_format[f"Q{idx}_minTime"] = q_data["min_timestamp"]

            # THD
            if "thd_u" in parsed_data:
                for i, thd_data in enumerate(parsed_data["thd_u"]):
                    idx = i + 1
                    legacy_format[f"THD_U{idx}_max"] = thd_data["max"]
                    legacy_format[f"THD_U{idx}_min"] = thd_data["min"]
                    legacy_format[f"THD_U{idx}_maxTime"] = thd_data["max_timestamp"]
                    legacy_format[f"THD_U{idx}_minTime"] = thd_data["min_timestamp"]

        return legacy_format

    def get_data_dict(self, meters: Dict[str, Any], parsed_data: Dict[str, Any], keys: List[Dict[str, str]],
                      unit: str) -> List[Dict[str, Any]]:
        """
        특정 키들의 데이터를 추출 (기존 RedisMapDetail2.get_Datadict와 동일한 형식)

        Args:
            meters: 현재값 딕셔너리
            parsed_data: 파싱된 MaxMin 데이터
            keys: 키 정보 리스트 [{"id": "1", "key": "U1", "label": "전압1"}, ...]
            unit: 단위

        Returns:
            추출된 데이터 리스트
        """
        result = []

        if parsed_data is None:
            # 기본값으로 처리
            for item in keys:
                result.append({
                    "id": item["id"],
                    "subTitle": item["label"],
                    "value": round(meters.get(item["key"], 0), 2),
                    "max": 0,
                    "maxTime": "",
                    "min": 0,
                    "minTime": "",
                    "unit": unit
                })
            return result

        for item in keys:
            key = item["key"]

            # 새 구조체에서 데이터 추출
            max_value = 0
            min_value = 0
            max_time = ""
            min_time = ""

            # 키 파싱 (예: "U1" -> field="u", index=0)
            if key and len(key) > 1:
                field_base = key[:-1].lower() if key[-1].isdigit() else key.lower()
                index = int(key[-1]) - 1 if key[-1].isdigit() else None

                # 필드 매핑
                field_mapping = {
                    "u": "u",  # 선간전압
                    "upp": "upp",  # 상전압
                    "i": "i",  # 전류
                    "pf": "pf",  # 역률
                    "p": "p",  # 유효전력
                    "q": "q",  # 무효전력
                    "s": "s",  # 피상전력
                    "freq": "freq",  # 주파수
                    "temp": "temp",  # 온도
                    "itot": "itot",  # 전체 전류
                    "in": "in",  # 중성선 전류
                    "isum": "isum"  # 전류 합계
                }

                field_name = field_mapping.get(field_base, field_base)

                # 데이터 추출
                if field_name in parsed_data:
                    if index is not None and isinstance(parsed_data[field_name], list):
                        if 0 <= index < len(parsed_data[field_name]):
                            data = parsed_data[field_name][index]
                            max_value = data.get("max", 0)
                            min_value = data.get("min", 0)
                            max_time = data.get("max_timestamp", "")
                            min_time = data.get("min_timestamp", "")
                    else:
                        # 단일 값 (freq, temp 등)
                        data = parsed_data[field_name]
                        if isinstance(data, dict):
                            max_value = data.get("max", 0)
                            min_value = data.get("min", 0)
                            max_time = data.get("max_timestamp", "")
                            min_time = data.get("min_timestamp", "")

            # datetime을 문자열로 변환
            if isinstance(max_time, datetime):
                max_time = max_time.strftime("%Y-%m-%d %H:%M:%S")
            if isinstance(min_time, datetime):
                min_time = min_time.strftime("%Y-%m-%d %H:%M:%S")

            result.append({
                "id": item["id"],
                "subTitle": item["label"],
                "value": round(meters.get(key, 0), 2),
                "max": round(max_value, 2),
                "maxTime": max_time,
                "min": round(min_value, 2),
                "minTime": min_time,
                "unit": unit
            })

        return result

    def get_max_values(self, channel: str, parsed_data: Dict[str, Any]) -> List[float]:
        """
        최대값 리스트 추출 (기존 코드의 max 배열과 동일)

        Args:
            channel: 채널명
            parsed_data: 파싱된 데이터

        Returns:
            [U1_max, U2_max, U3_max, I1_max, I2_max, I3_max]
        """
        max_values = []

        # U1, U2, U3
        if "u" in parsed_data:
            for i in range(3):
                if i < len(parsed_data["u"]):
                    max_values.append(parsed_data["u"][i]["max"])
                else:
                    max_values.append(0.0)
        else:
            max_values.extend([0.0, 0.0, 0.0])

        # I1, I2, I3
        if "i" in parsed_data:
            for i in range(3):
                if i < len(parsed_data["i"]):
                    max_values.append(parsed_data["i"][i]["max"])
                else:
                    max_values.append(0.0)
        else:
            max_values.extend([0.0, 0.0, 0.0])

        return max_values


# 사용 예시
"""
# 1. 초기화
handler = MaxMinDataHandler(processor)

# 2. 현재값 데이터 (meters) - 기존처럼 Redis에서 가져오거나 다른 소스에서
meters = {
    "U1": 220.5,
    "U2": 221.0,
    "U3": 219.8,
    "I1": 10.5,
    "I2": 11.2,
    "I3": 10.8,
    "Freq": 60.0
}

# 3. MaxMin 데이터 가져오기 (새 구조체)
parsed_maxmin = handler.get_maxmin_data("Main", "1sec")

# 4. 기존 키 정의 형식 (딕셔너리 리스트)
p_voltage_keys = [
    {"id": "1", "key": "Upp1", "label": "R상 전압"},
    {"id": "2", "key": "Upp2", "label": "S상 전압"},
    {"id": "3", "key": "Upp3", "label": "T상 전압"}
]

l_voltage_keys = [
    {"id": "1", "key": "U1", "label": "RS 전압"},
    {"id": "2", "key": "U2", "label": "ST 전압"},
    {"id": "3", "key": "U3", "label": "TR 전압"}
]

current_keys = [
    {"id": "1", "key": "I1", "label": "R상 전류"},
    {"id": "2", "key": "I2", "label": "S상 전류"},
    {"id": "3", "key": "I3", "label": "T상 전류"}
]

freq_keys = [
    {"id": "1", "key": "Freq", "label": "주파수"}
]

# 5. 기존 방식과 동일하게 사용
p_voltage_data = handler.get_data_dict(meters, parsed_maxmin, p_voltage_keys, 'V')
l_voltage_data = handler.get_data_dict(meters, parsed_maxmin, l_voltage_keys, 'V')
current_data = handler.get_data_dict(meters, parsed_maxmin, current_keys, 'A')
freq_data = handler.get_data_dict(meters, parsed_maxmin, freq_keys, 'Hz')

# 결과 예시:
# p_voltage_data = [
#     {
#         "id": "1",
#         "subTitle": "R상 전압",
#         "value": 220.5,
#         "max": 230.2,
#         "maxTime": "2024-01-15 10:30:00",
#         "min": 210.8,
#         "minTime": "2024-01-15 09:45:00",
#         "unit": "V"
#     },
#     ...
# ]

# 6. 필요시 기존 형식으로 변환
legacy_maxmin = handler.convert_to_legacy_format(parsed_maxmin)
# 이제 legacy_maxmin은 {"U1_max": 230.2, "U1_min": 210.8, ...} 형식

# 7. 최대값 배열 추출
max_values = handler.get_max_values("Main", parsed_maxmin)
"""