#!/usr/bin/env python3
"""
Redis 데이터 마이그레이션 스크립트
From: 192.168.2.92 => To: 127.0.0.1 (로컬)

마이그레이션 대상:
- DB 1: 해시 meter_main, pq_main의 모든 데이터
- DB 8: 해시 tblist1, tblist2, tblist3, tblist4의 모든 데이터
"""

import redis
import sys
from datetime import datetime
import logging

# 로깅 설정
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('redis_migration.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

class RedisMigration:
    def __init__(self, source_host='192.168.2.92', source_port=6379, 
                 dest_host='127.0.0.1', dest_port=6379,
                 source_password=None, dest_password=None):
        """
        Redis 마이그레이션 초기화
        
        Args:
            source_host: 소스 Redis 호스트
            source_port: 소스 Redis 포트
            dest_host: 대상 Redis 호스트
            dest_port: 대상 Redis 포트
            source_password: 소스 Redis 비밀번호 (옵션)
            dest_password: 대상 Redis 비밀번호 (옵션)
        """
        self.source_host = source_host
        self.source_port = source_port
        self.dest_host = dest_host
        self.dest_port = dest_port
        
        # Redis 연결 설정
        try:
            self.source_redis = redis.Redis(
                host=source_host,
                port=source_port,
                password=source_password,
                decode_responses=False  # 바이너리 데이터 처리를 위해 False
            )
            
            self.dest_redis = redis.Redis(
                host=dest_host,
                port=dest_port,
                password=dest_password,
                decode_responses=False
            )
            
            # 연결 테스트
            self.source_redis.ping()
            self.dest_redis.ping()
            
            logger.info(f"소스 Redis 연결 성공: {source_host}:{source_port}")
            logger.info(f"대상 Redis 연결 성공: {dest_host}:{dest_port}")
            
        except redis.ConnectionError as e:
            logger.error(f"Redis 연결 실패: {e}")
            sys.exit(1)
    
    def migrate_hash(self, source_db, dest_db, hash_key):
        """
        특정 해시 키의 모든 데이터를 마이그레이션
        
        Args:
            source_db: 소스 DB 번호
            dest_db: 대상 DB 번호
            hash_key: 마이그레이션할 해시 키
        
        Returns:
            마이그레이션된 필드 수
        """
        try:
            # DB 선택
            self.source_redis.select(source_db)
            self.dest_redis.select(dest_db)
            
            # 해시 키가 존재하는지 확인
            if not self.source_redis.exists(hash_key):
                logger.warning(f"DB {source_db}에 해시 키 '{hash_key.decode() if isinstance(hash_key, bytes) else hash_key}'가 존재하지 않습니다.")
                return 0
            
            # 모든 해시 필드와 값 가져오기
            hash_data = self.source_redis.hgetall(hash_key)
            
            if not hash_data:
                logger.warning(f"해시 키 '{hash_key.decode() if isinstance(hash_key, bytes) else hash_key}'에 데이터가 없습니다.")
                return 0
            
            # 대상 Redis에 데이터 저장
            # 파이프라인을 사용하여 성능 향상
            pipeline = self.dest_redis.pipeline()
            
            # 기존 데이터 삭제 (선택사항)
            # pipeline.delete(hash_key)
            
            # 해시 데이터 설정
            for field, value in hash_data.items():
                pipeline.hset(hash_key, field, value)
            
            pipeline.execute()
            
            field_count = len(hash_data)
            logger.info(f"DB {source_db} -> DB {dest_db}: 해시 '{hash_key.decode() if isinstance(hash_key, bytes) else hash_key}' - {field_count}개 필드 마이그레이션 완료")
            
            return field_count
            
        except Exception as e:
            logger.error(f"해시 마이그레이션 중 오류 발생: {e}")
            raise
    
    def verify_migration(self, db_num, hash_key):
        """
        마이그레이션 검증 - 소스와 대상의 데이터 비교
        
        Args:
            db_num: DB 번호
            hash_key: 검증할 해시 키
        
        Returns:
            검증 성공 여부
        """
        try:
            self.source_redis.select(db_num)
            self.dest_redis.select(db_num)
            
            source_data = self.source_redis.hgetall(hash_key)
            dest_data = self.dest_redis.hgetall(hash_key)
            
            if source_data == dest_data:
                logger.info(f"검증 성공: DB {db_num}, 해시 '{hash_key.decode() if isinstance(hash_key, bytes) else hash_key}'")
                return True
            else:
                logger.error(f"검증 실패: DB {db_num}, 해시 '{hash_key.decode() if isinstance(hash_key, bytes) else hash_key}' - 데이터 불일치")
                return False
                
        except Exception as e:
            logger.error(f"검증 중 오류 발생: {e}")
            return False
    
    def run_migration(self):
        """
        전체 마이그레이션 실행
        """
        logger.info("=" * 50)
        logger.info("Redis 데이터 마이그레이션 시작")
        logger.info(f"소스: {self.source_host}:{self.source_port}")
        logger.info(f"대상: {self.dest_host}:{self.dest_port}")
        logger.info("=" * 50)
        
        start_time = datetime.now()
        total_fields = 0
        success_count = 0
        
        # 마이그레이션 작업 정의
        migration_tasks = [
            # DB 1: meter_main, pq_main
            {'db': 1, 'hashes': [b'meter_main', b'pq_main']},
            {'db': 5, 'hashes': [b'device:192.168.9.100:3:main']},
            # DB 8: tblist1, tblist2, tblist3, tblist4
            {'db': 8, 'hashes': [b'tblist1', b'tblist2', b'tblist3', b'tblist4']}
        ]
        
        # 마이그레이션 실행
        for task in migration_tasks:
            db_num = task['db']
            logger.info(f"\nDB {db_num} 마이그레이션 시작...")
            
            for hash_key in task['hashes']:
                try:
                    # 데이터 마이그레이션
                    field_count = self.migrate_hash(db_num, db_num, hash_key)
                    total_fields += field_count
                    
                    # 검증
                    if field_count > 0:
                        if self.verify_migration(db_num, hash_key):
                            success_count += 1
                        
                except Exception as e:
                    logger.error(f"마이그레이션 실패 - DB {db_num}, 해시 '{hash_key.decode()}': {e}")
        
        # 완료 시간 계산
        end_time = datetime.now()
        duration = end_time - start_time
        
        # 결과 요약
        logger.info("=" * 50)
        logger.info("마이그레이션 완료")
        logger.info(f"소요 시간: {duration}")
        logger.info(f"총 마이그레이션된 필드 수: {total_fields}")
        logger.info(f"성공한 해시 수: {success_count}")
        logger.info("=" * 50)
    
    def close_connections(self):
        """
        Redis 연결 종료
        """
        try:
            self.source_redis.close()
            self.dest_redis.close()
            logger.info("Redis 연결 종료")
        except Exception as e:
            logger.error(f"연결 종료 중 오류: {e}")


def main():
    """
    메인 실행 함수
    """
    # Redis 연결 설정
    # 필요시 비밀번호 추가
    SOURCE_PASSWORD = None  # 소스 Redis 비밀번호 (필요시 설정)
    DEST_PASSWORD = None    # 대상 Redis 비밀번호 (필요시 설정)
    
    try:
        # 마이그레이션 객체 생성
        migration = RedisMigration(
            source_host='192.168.2.88',
            source_port=6379,
            dest_host='127.0.0.1',
            dest_port=6379,
            source_password=SOURCE_PASSWORD,
            dest_password=DEST_PASSWORD
        )
        
        # 마이그레이션 실행
        migration.run_migration()
        
        # 연결 종료
        migration.close_connections()
        
    except KeyboardInterrupt:
        logger.info("\n마이그레이션이 사용자에 의해 중단되었습니다.")
        sys.exit(1)
    except Exception as e:
        logger.error(f"예상치 못한 오류 발생: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()