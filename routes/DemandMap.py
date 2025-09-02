import math
from typing import List, Dict, Any


class DemandDataFormatter:
    """Demand 데이터를 지정된 형식으로 변환하는 클래스"""

    @staticmethod
    def sanitize_float(value: float, default: float = 0.0) -> float:
        """INF/NaN을 체크하고 기본값으로 대체"""
        if value is None:
            return default
        try:
            if math.isfinite(value):
                return float(value)
            else:
                return default
        except (TypeError, ValueError):
            return default

    @classmethod
    def safe_round(cls, value: float, decimals: int = 2, default: float = 0.0) -> float:
        """안전한 반올림"""
        sanitized = cls.sanitize_float(value, default)
        return round(sanitized, decimals)

    @classmethod
    def get_Datadict(cls, meters: Dict[str, float], maxmin: Dict[str, Any],
                     keys: List[Dict[str, Any]], unit: str) -> List[Dict[str, Any]]:
        """기존 형식과 동일한 데이터 포맷"""
        p_voltage_data = []
        for item in keys:
            key = item["key"]
            p_voltage_data.append({
                "id": item["id"],
                "subTitle": item["label"],
                "value": cls.safe_round(meters.get(key, 0), 2),
                "max": cls.safe_round(maxmin.get(f"{key}_max", 0), 2),
                "maxTime": maxmin.get(f"{key}_maxTime", ""),
                "min": cls.safe_round(maxmin.get(f"{key}_min", 0), 2),
                "minTime": maxmin.get(f"{key}_minTime", ""),
                "unit": unit
            })
        return p_voltage_data

    @classmethod
    def format_demand_data(cls, parsed_demand: Dict[str, Any]) -> Dict[str, Any]:
        """파싱된 Demand 데이터를 응답 형식으로 변환"""

        # 1. Present Demand (현재값) - meters
        meters = {
            "p_import": parsed_demand["CD_P"][0],
            "p_export": parsed_demand["CD_P"][1],
            "q_import": parsed_demand["CD_Q"][0],
            "q_export": parsed_demand["CD_Q"][1],
            "s": parsed_demand["CD_S"],
            "i_a": parsed_demand["DD_I"][0],  # Dynamic Demand의 전류값 사용
            "i_b": parsed_demand["DD_I"][1],
            "i_c": parsed_demand["DD_I"][2]
        }

        # 2. Max Demand - maxmin
        maxmin = {
            # 유효전력
            "p_import_max": parsed_demand["MD_P"][0]["value"],
            "p_import_maxTime": parsed_demand["MD_P"][0]["time"].strftime("%Y-%m-%d %H:%M:%S"),
            "p_import_min": 0,  # Demand는 최소값이 없으므로 0
            "p_import_minTime": "",

            "p_export_max": parsed_demand["MD_P"][1]["value"],
            "p_export_maxTime": parsed_demand["MD_P"][1]["time"].strftime("%Y-%m-%d %H:%M:%S"),
            "p_export_min": 0,
            "p_export_minTime": "",

            # 무효전력
            "q_import_max": parsed_demand["MD_Q"][0]["value"],
            "q_import_maxTime": parsed_demand["MD_Q"][0]["time"].strftime("%Y-%m-%d %H:%M:%S"),
            "q_import_min": 0,
            "q_import_minTime": "",

            "q_export_max": parsed_demand["MD_Q"][1]["value"],
            "q_export_maxTime": parsed_demand["MD_Q"][1]["time"].strftime("%Y-%m-%d %H:%M:%S"),
            "q_export_min": 0,
            "q_export_minTime": "",

            # 피상전력
            "s_max": parsed_demand["MD_S"]["value"],
            "s_maxTime": parsed_demand["MD_S"]["time"].strftime("%Y-%m-%d %H:%M:%S"),
            "s_min": 0,
            "s_minTime": "",

            # 전류
            "i_a_max": parsed_demand["MD_I"][0]["value"],
            "i_a_maxTime": parsed_demand["MD_I"][0]["time"].strftime("%Y-%m-%d %H:%M:%S"),
            "i_a_min": 0,
            "i_a_minTime": "",

            "i_b_max": parsed_demand["MD_I"][1]["value"],
            "i_b_maxTime": parsed_demand["MD_I"][1]["time"].strftime("%Y-%m-%d %H:%M:%S"),
            "i_b_min": 0,
            "i_b_minTime": "",

            "i_c_max": parsed_demand["MD_I"][2]["value"],
            "i_c_maxTime": parsed_demand["MD_I"][2]["time"].strftime("%Y-%m-%d %H:%M:%S"),
            "i_c_min": 0,
            "i_c_minTime": ""
        }

        # 3. 키 정의
        power_keys = [
            {"id": 1, "key": "p_import", "label": "Import"},
            {"id": 2, "key": "p_export", "label": "Export"}
        ]

        reactive_keys = [
            {"id": 1, "key": "q_import", "label": "Import"},
            {"id": 2, "key": "q_export", "label": "Export"}
        ]

        apparent_keys = [
            {"id": 1, "key": "s", "label": "Total"}
        ]

        current_keys = [
            {"id": 1, "key": "i_a", "label": "L1"},
            {"id": 2, "key": "i_b", "label": "L2"},
            {"id": 3, "key": "i_c", "label": "L3"}
        ]

        # 4. 형식에 맞게 변환
        return {
            "power_demand": cls.get_Datadict(meters, maxmin, power_keys, "kW"),
            "reactive_demand": cls.get_Datadict(meters, maxmin, reactive_keys, "kVar"),
            "apparent_demand": cls.get_Datadict(meters, maxmin, apparent_keys, "kVA"),
            "current_demand": cls.get_Datadict(meters, maxmin, current_keys, "A"),
            "predict_demand": {
                "value": cls.safe_round(parsed_demand["PD_P"], 2),
                "unit": "kW"
            },
            "dynamic_demand": {
                "timestamp": parsed_demand["ddTime"].strftime("%Y-%m-%d %H:%M:%S"),
                "power": parsed_demand["DD_P"],
                "reactive": parsed_demand["DD_Q"],
                "apparent": parsed_demand["DD_S"],
                "current": parsed_demand["DD_I"]
            }
        }