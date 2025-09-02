from .RedisBinary import BinaryStructConfig, BinaryDataProcessor
from .demand_parser import DemandParser
from .RedisBinary import register_waveform_config

import logging


def setup_demand_processor(processor: BinaryDataProcessor):
    """
    Demand 데이터 처리기 설정

    Args:
        processor: BinaryDataProcessor 인스턴스
    """
    # Demand 설정 생성
    demand_config = BinaryStructConfig(
        name="demand",
        struct_format='<' + 'B' * 124,  # 108 바이트 (실제로는 파서가 처리)
        field_names=["demand_data"],  # 파서가 처리하므로 단순화
        field_descriptions={
            "demand_data": "Demand data including Max, Dynamic, Present, and Predict values"
        }
    )

    # Demand 파서 등록
    demand_parser = DemandParser()

    # 프로세서에 등록
    processor.register_config(
        config=demand_config,
        redis_key_pattern="Demand",  # Redis 키 패턴
        parser=demand_parser
    )

    return demand_config


def setup_all_processors(processor: BinaryDataProcessor):
    """
    모든 바이너리 구조체 설정을 등록
    """
    # Demand 설정
    setup_demand_processor(processor)
    # handler = MaxMinDataHandler(processor)

    # Waveform 설정
    register_waveform_config(processor)


    # 기타 설정들...

    logging.info(f"Registered {len(processor.configs)} binary configurations")