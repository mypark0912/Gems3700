from fastapi import APIRouter, Request, UploadFile, File
from fastapi.responses import FileResponse
from fastapi.responses import JSONResponse
import json, os, httpx
import shutil, random

# @router.get('/checkSetting') #check setup_channel at influxdb
# def check_setting():
#     try:
#         # ✅ Flux 쿼리 작성 (Main과 Sub의 Diagnosis 및 PowerQuality 값을 가져옴)
#         flux_query = f'''
#         from(bucket: "ssdb")
#           |> range(start: -30d)
#           |> filter(fn: (r) => r["_measurement"] == "setup_channel")
#           |> filter(fn: (r) => r["_field"] == "Diagnosis" or r["_field"] == "PowerQuality" or r["_field"] == "Enable")
#           |> filter(fn: (r) => r["channel"] == "Main" or r["channel"] == "Sub")
#           |> last()
#         '''
#
#         # ✅ 쿼리 실행
#         result = query_api.query(org='ntek', query=flux_query)
#
#         # ✅ 결과 데이터 저장 변수 (기본값 False)
#         Diag_main = False
#         Diag_sub = False
#         Enable_main = False
#         Enable_sub = False
#
#         # ✅ 결과 데이터 파싱 및 변환 (0 → False, 1 → True)
#         for table in result:
#             for record in table.records:
#                 value = bool(record["_value"])  # ✅ 0 → False, 1 → True 변환
#
#                 if record["channel"] == "Main":
#                     if record["_field"] == "Diagnosis":
#                         Diag_main = value
#                     elif record["_field"] == "Enable":
#                         Enable_main = value
#                 elif record["channel"] == "Sub":
#                     if record["_field"] == "Diagnosis":
#                         Diag_sub = value
#                     elif record["_field"] == "Enable":
#                         Enable_sub = value
#
#         # ✅ 결과 반환 (Boolean 값)
#         resultDict = {
#             "result" : "1",
#             "Diag_main": Diag_main,
#             "Diag_sub": Diag_sub,
#             "enable_main" : Enable_main,
#             "enable_sub" : Enable_sub
#         }
#         return resultDict
#
#     except Exception as e:
#         print("Error:", e)
#         return {"result":"0"}
#
# @router.get('/getSetting/{channel}') #get setup_channel at influxdb
# def getSetting(channel, request:Request):
#     tb_name = ''
#     query=''
#     if 'user' not in request.session:
#         return {"error": "Unauthorized"}
#
#     if channel == 'General':
#         tb_name = 'setup_general'
#         query = (
#             f'from(bucket: "ssdb") |> range(start: -30d) |> filter(fn: (r) => r["_measurement"] == "{tb_name}") |> group(columns: ["_time", "device_model"])'
#             '|> map(fn: (r) => ({ r with _value: string(v: r._value) }))'
#             f'|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")')
#     else:
#         tb_name = 'setup_channel'
#         query = (f'from(bucket: "ssdb") |> range(start: -30d) '
#                  f'|> filter(fn: (r) => r["_measurement"] == "{tb_name}")'
#                  f'|> filter(fn: (r) => r["channel"] == "{channel}")'  # channel 필터 추가
#                  '|> group(columns: ["_time", "channel"])'
#                  '|> map(fn: (r) => ({ r with _value: string(v: r._value) }))'
#                  '|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")')
#     try:
#         tables = query_api.query(org='ntek', query=query)
#
#         results = []
#         for table in tables:
#             for record in table.records:
#                 record_data = record.values  # ✅ 모든 필드 데이터를 포함
#                 results.append(record_data)
#         print(results)
#         return results
#
#     except Exception as e:
#         return {"error": str(e)}
#
# @router.post('/save/{channel}') #save setting at influxdb
# async def saveSetting(channel, request:Request):
#     try:
#         # JSON 데이터 받기
#         data = await request.json()
#         if not data:
#             return {"status": "0"}
#         print("Received Data:", data)  # 콘솔에서 데이터 확인
#
#         meas = "setup_general" if channel == "General" else "setup_channel"
#
#         # ✅ `channel`이 General이면 `device_model`, 아니면 `channel`을 태그 키로 설정
#         tagKey = "device_model" if channel == "General" else "channel"
#         tagValue = str(data.get(tagKey, "unknown"))  # 태그 값 설정 (기본값 "unknown")
#
#         point = Point(meas).tag(tagKey, tagValue)
#
#         for key, value in data.items():
#             if key in ["channel", "device_model"]:  # 이미 태그로 저장된 값 제외
#                 continue
#             try:
#                 value = int(value) if isinstance(value, (int, bool)) or (isinstance(value, str) and value.isdigit()) else str(value)
#             except:
#                 value = str(value) if value is not None else ""
#             point = point.field(key, value)
#
#         write_api.write(bucket='ssdb', org='ntek', record=point)
#         return {"status": "1"}
#
#     except Exception as e:
#         print("Error:", e)
#         return {"status": "0"
#         }