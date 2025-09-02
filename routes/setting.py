from fastapi import APIRouter, Request, UploadFile, File
from fastapi.responses import FileResponse
from fastapi.responses import JSONResponse
import os, httpx, csv, psutil, struct
import ujson as json
import shutil, subprocess, logging
from datetime import datetime
from states.global_state import influx_state, redis_state, aesState,os_spec
from collections import defaultdict
from routes.auth import checkLoginAPI,get_mac_address
from routes.api import parameter_options
from .RedisBinary import Command, CmdType, ItemType

# Path 객체 절대경로
from pathlib import Path
base_dir = Path(__file__).resolve().parent
SETTING_FOLDER = base_dir.parent.parent / "config"  # ⬅️ 두 단계 상위로

setting_timeout = httpx.Timeout(
    connect=2.0,  # 연결에는 5초
    read=30.0,     # 응답 읽기는 2초
    write=2.0,    # 요청 전송은 5초
    pool=5.0      # 연결 풀은 5초
)

router = APIRouter()

def command_to_binary(command: Command) -> bytes:
    """Command 객체를 바이너리로 변환 (C 구조체와 호환)"""
    # C의 int는 보통 4바이트, enum도 int로 처리
    # 구조: int(4) + int(4) + int(4) = 12 bytes
    return struct.pack('iii', command.type, command.cmd, command.item)


def binary_to_command(binary_data: bytes) -> Command:
    """바이너리 데이터를 Command 객체로 변환"""
    if len(binary_data) != 12:
        raise ValueError("Invalid binary data size")

    type_val, cmd_val, item_val = struct.unpack('iii', binary_data)
    return Command(type=type_val, cmd=CmdType(cmd_val), item=ItemType(item_val))


def read_mac_plain(filepath):
    """파일에서 MAC 주소를 구분자 없이 읽기"""
    with open(filepath, 'rb') as f:
        mac_bytes = f.read(6)

        if len(mac_bytes) != 6:
            raise ValueError(f"잘못된 MAC 주소 길이: {len(mac_bytes)} bytes")

        # 구분자 없이 반환
        return mac_bytes.hex().lower()  # 더 간단한 방법

@router.get('/initDB')
async def initInflux():
    file_path = os.path.join(SETTING_FOLDER, 'influx.json')
    data = {
          "username": "admin",
          "password": "ntek9135",
          "org": "ntek",
          "bucket": "ntek",
          "retentionPeriodSeconds": 0  # 90 * 24 * 60 * 60
        }
    try:
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.post(f"http://127.0.0.1:8086/api/v2/setup", json=data)
            resData = response.json()
            cipher = aesState.encrypt(resData.get("auth").get("token"))
            influxdata = {
                "token" : cipher,
                "org" : resData.get("org").get("name"),
                "retention":data.get("retentionPeriodSeconds")
            }
        with open(file_path, "w", encoding="utf-8") as f:
            json.dump(influxdata, f, indent=4)

        # init_influx()  # ✅ 초기화 수행 (json 생성 + client 전역 등록)
        if influx_state.client is None:
            return {"result": False}

        if influx_state.error:
            return {"success": False, "message": influx_state.error}

        return {"success": True, "message": "InfluxDB initialized successfully"}
    except Exception as e:
        logging.error(f"❌ Influxdb Init Error: {e}")
        influx_state.client = None
        influx_state.error = f"Exception during init: {str(e)}"
        return {"success": False, "message": influx_state.error}


def parse_settings(setting):
    """설정을 파싱하여 결과 딕셔너리 생성"""
    # 기본값 설정
    result = {
        "result": "1",
        "mode": setting.get("mode", ""),
        "lang": setting.get("lang", ""),
        "location": "",
        "Diag_main": False,
        "Diag_sub": False,
        "enable_main": False,
        "enable_sub": False,
        "pq_main": False,
        "pq_sub": False,
        "assetType_main": "",
        "assetName_main": "",
        "assetType_sub": "",
        "assetName_sub": "",
        "assetNickname_main": "",
        "assetNickname_sub": "",
        "main_kva": -1,
        "sub_kva": -1,
        "pf_sign": -1,
        "va_type": -1,
        "unbalance":-1
    }

    # channel 검증
    if "channel" not in setting:
        result["result"] = "0"
        return result

    # device 모드인 경우만 상세 정보 파싱
    if 'device' in result["mode"]:
        # General 정보 파싱
        general_data = setting.get("General", {})
        use_function = general_data.get("useFuction", {})  # 오타 수정 필요
        dev_info = general_data.get("deviceInfo", {})
        result["pf_sign"] = general_data.get("pf_sign", 0)
        result["va_type"] = general_data.get("va_type", 0)
        result["unbalance"] = general_data.get("unbalance", 0)
        result["location"] = dev_info.get("location", "")
        result["Diag_main"] = bool(use_function.get("diagnosis_main", 0))
        result["Diag_sub"] = bool(use_function.get("diagnosis_sub", 0))

        # 채널별 정보 파싱
        for channel_data in setting["channel"]:
            channel_name = channel_data.get("channel", "")

            if channel_name == "Main":
                parse_channel_data(channel_data, result, "main", result["Diag_main"])
            elif channel_name == "Sub":
                parse_channel_data(channel_data, result, "sub", result["Diag_sub"])

    return result


def parse_channel_data(channel_data, result, prefix, diag_enabled):
    """채널 데이터를 파싱하여 결과에 추가"""
    result[f"enable_{prefix}"] = bool(channel_data.get("Enable", 0))
    result[f"pq_{prefix}"] = bool(channel_data.get("PowerQuality", 0))

    if diag_enabled:
        asset_info = channel_data.get("assetInfo", {})
        result[f"assetType_{prefix}"] = asset_info.get("type", "")
        result[f"assetName_{prefix}"] = asset_info.get("name", "")
        result[f"assetNickname_{prefix}"] = asset_info.get("nickname", "")

        if asset_info.get("type") == 'Transformer':
            result[f"{prefix}_kva"] = int(channel_data.get("n_kva", 0))

@router.get('/checkSettingFile') #check setup.json
def check_setupfile(request: Request):
    """
    우선순위:
    1. setup.json 파일 확인 (최우선)
    2. 파일이 없으면 default.json으로 생성
    3. Redis에 없으면 파일에서 로드하여 Redis에 저장
    """
    file_path = os.path.join(SETTING_FOLDER, 'setup.json')
    default_file_path = os.path.join(SETTING_FOLDER, 'default.json')
    deviceMac = get_mac_address()
    if os_spec.os != 'Windows':
        mac_file_path = os.path.join(SETTING_FOLDER, 'serial_num_do_not_modify.txt')
        if os.path.exists(mac_file_path):
            ser = read_mac_plain(mac_file_path)

            if ser != deviceMac:
                deviceMac = ser

    if redis_state.client is None:
        return {"result": "0", "error": "Redis not available"}

    try:
        redis_state.client.select(0)

        # 1. 먼저 setup.json 파일 확인 및 생성
        if not os.path.exists(file_path):
            # setup.json이 없으면 default.json으로 생성
            if os.path.exists(default_file_path):
                shutil.copy(default_file_path, file_path)
            else:
                return {"result": "0", "error": "No default.json found"}

        # 2. Redis에 setup이 있는지 확인
        if redis_state.client.hexists("System", "setup"):
            # Redis에 있으면 Redis 데이터 사용
            setting = json.loads(redis_state.client.hget("System", "setup"))
            # deviceMac = get_mac_address()
            if deviceMac != setting["General"]["deviceInfo"]["mac_address"]:
                setting["General"]["deviceInfo"]["mac_address"] = deviceMac
                setting["General"]["deviceInfo"]["serial_number"] = deviceMac
        else:
            # Redis에 없으면 파일에서 읽어서 Redis에 저장
            try:
                with open(file_path, "r", encoding="utf-8") as f:
                    setting = json.load(f)
            except json.JSONDecodeError:
                # setup.json이 손상된 경우 default.json으로 복구
                shutil.copy(default_file_path, file_path)
                with open(file_path, "r", encoding="utf-8") as f:
                    setting = json.load(f)

            # Redis에 저장
            # deviceMac = get_mac_address()
            if deviceMac != setting["General"]["deviceInfo"]["mac_address"]:
                setting["General"]["deviceInfo"]["mac_address"] = deviceMac
                setting["General"]["deviceInfo"]["serial_number"] = deviceMac

            redis_state.client.hset("System", "setup", json.dumps(setting))
            if "mode" in setting:
                redis_state.client.hset("System", "mode", setting["mode"])

        # 3. 설정 파싱 및 반환
        return parse_settings(setting)

    except Exception as e:
        return {"result": "0", "error": str(e)}


# def check_setupfile(request: Request):
#     file_path = os.path.join(SETTING_FOLDER, 'setup.json')
#     opmode = ''
#     if not redis_state.client is None:
#         redis_state.client.select(0)
#         if os.path.exists(file_path):
#             if redis_state.client.hexists("System", "setup"):
#                 opmode = redis_state.client.hget("System", "mode")
#                 redisContext = redis_state.client.hget("System", "setup")
#                 setting = json.loads(redisContext)
#             else:
#                 try:
#                     with open(file_path, "r", encoding="utf-8") as f:
#                         setting = json.load(f)
#                         redis_state.client.hset("System", "mode", setting["mode"])
#                         redis_state.client.hset("System", "setup",json.dumps(setting))
#                 except Exception as e:
#                     default_file_path = os.path.join(SETTING_FOLDER, 'default.json')
#                     shutil.copy(default_file_path, file_path)
#                     with open(file_path, "r", encoding="utf-8") as f:
#                         setting = json.load(f)
#                         redis_state.client.hset("System", "mode", setting["mode"])
#                         redis_state.client.hset("System", "setup", json.dumps(setting))
#                     return {"result": "0"}
#         else:
#             default_file_path = os.path.join(SETTING_FOLDER, 'default.json')
#             try:
#                 shutil.copy(default_file_path, file_path)
#                 with open(file_path, "r", encoding="utf-8") as f:
#                     setting = json.load(f)
#                     redis_state.client.hset("System", "mode", setting["mode"])
#                     redis_state.client.hset("System", "setup", json.dumps(setting))
#             except Exception as e:
#                 return {"result": "0"}
#     Diag_main = False
#     Diag_sub = False
#     Enable_main = False
#     Enable_sub = False
#     PQ_main = False
#     PQ_sub = False
#     assetType_main = ''
#     assetName_main = ''
#     assetType_sub = ''
#     assetName_sub = ''
#     assetNickName_main = ''
#     assetNickName_sub = ''
#     lang = ''
#     main_kva = -1
#     sub_kva = -1
#     if not "channel" in setting:
#         return {"result": "0"}
#     if opmode == '':
#         opmode = setting.get("mode", "")
#     lang = setting.get("lang", "")
#     if 'device' in opmode:
#         general_data = setting.get("General", {})
#         use_fuction = general_data.get("useFuction", {})
#         devInfo = general_data.get("deviceInfo",{})
#         devLocation = devInfo.get("location","")
#         Diag_main = bool(use_fuction.get("diagnosis_main", 0))
#         Diag_sub = bool(use_fuction.get("diagnosis_sub", 0))
#
#         main_channel_data = next((ch for ch in setting["channel"] if ch.get("channel") == "Main"), None)
#         if main_channel_data:
#             Enable_main = bool(main_channel_data.get("Enable", 0))
#             PQ_main = bool(main_channel_data.get("PowerQuality", 0))
#             if Diag_main:
#                 asset_info_main = main_channel_data.get("assetInfo", {})
#                 assetType_main = asset_info_main.get("type", "")
#                 assetName_main = asset_info_main.get("name", "")
#                 assetNickName_main = asset_info_main.get("nickname", "")
#                 if assetType_main == 'Transformer':
#                     main_kva = int(main_channel_data.get("n_kva", 0))
#
#         sub_channel_data = next((ch for ch in setting["channel"] if ch.get("channel") == "Sub"), None)
#         if sub_channel_data:
#             Enable_sub = bool(sub_channel_data.get("Enable", 0))
#             PQ_sub = bool(sub_channel_data.get("PowerQuality", 0))
#
#             if Diag_sub:
#                 asset_info_sub = sub_channel_data.get("assetInfo", {})
#                 assetType_sub = asset_info_sub.get("type", "")
#                 assetName_sub = asset_info_sub.get("name", "")
#                 assetNickName_sub = asset_info_sub.get("nickname", "")
#                 if assetType_sub == 'Transformer':
#                     sub_kva = int(sub_channel_data.get("n_kva", 0))
#
#         resultDict = {
#             "result": "1",
#             "mode": opmode,
#             "lang": lang,
#             "location":devLocation,
#             "Diag_main": Diag_main,
#             "Diag_sub": Diag_sub,
#             "enable_main": Enable_main,
#             "enable_sub": Enable_sub,
#             "pq_main": PQ_main,
#             "pq_sub": PQ_sub,
#             "assetType_main": assetType_main,
#             "assetName_main": assetName_main,
#             "assetType_sub": assetType_sub,
#             "assetName_sub": assetName_sub,
#             "assetNickname_main": assetNickName_main,
#             "assetNickname_sub": assetNickName_sub,
#             "main_kva": main_kva,
#             "sub_kva": sub_kva,
#         }
#     else:
#         resultDict = {
#             "result": "1",
#             "mode": opmode,
#             "lang": lang,
#             "location": '',
#             "Diag_main": False,
#             "Diag_sub": False,
#             "enable_main": False,
#             "enable_sub": False,
#             "pq_main": False,
#             "pq_sub": False,
#             "assetType_main": '',
#             "assetName_main": '',
#             "assetType_sub": '',
#             "assetName_sub": '',
#             "assetNickname_main": assetNickName_main,
#             "assetNickname_sub": assetNickName_sub,
#             "main_kva": main_kva,
#             "sub_kva": sub_kva,
#         }
#     # request.session["devMode"] = opmode
#     return resultDict

@router.get('/getSetting')
def get_setting():
    redis_state.client.execute_command("SELECT", 0)
    if redis_state.client.hexists("System", "setup"):
        redisContext = redis_state.client.hget("System", "setup")
        setting = json.loads(redisContext)
    else:
        file_path = os.path.join(SETTING_FOLDER, 'setup.json')
        if not os.path.exists(file_path):
            return {"passOK": 0, "error": "setting file not found"}
        try:
            with open(file_path, "r", encoding="utf-8") as f:
                setting = json.load(f)
        except Exception as e:
            return {"passOK": 0}

    setup_dict = {
        "mode": setting.get("mode", ""),
        "General": setting.get("General", {}),
        "main": {},
        "sub": {}
    }

    for ch in setting.get("channel", []):
        name = ch.get("channel", "")
        if name == "Main":
            setup_dict["main"] = ch
        elif name == "Sub":
            setup_dict["sub"] = ch

    return {"passOK": 1, "data":setup_dict}

@router.get('/getSettingData/{channel}')  # get setting with setup.json
def getDatafromSetting(channel):
    redis_state.client.execute_command("SELECT", 0)
    if redis_state.client.hexists("System", "setup"):
        redisContext = redis_state.client.hget("System", "setup")
        setting = json.loads(redisContext)
    else:
        file_path = os.path.join(SETTING_FOLDER, 'setup.json')
        if not os.path.exists(file_path):
            return {"passOK": "0", "error": "setting file not found"}
        try:
            with open(file_path, "r", encoding="utf-8") as f:
                    setting = json.load(f)
        except Exception as e:
            return {"passOK": "0"}

    if channel in setting:
        return {"passOK": "1", "data": setting[channel]}
    else:
        main_channel = [ch for ch in setting["channel"] if ch["channel"] == channel]
        if not main_channel:
            return {"passOK": "0"}
        else:
            return {"passOK": "1", "data": main_channel}

# @router.get('/SysService/{cmd}/{item}')
# async def get_SysService(cmd, item, request:Request):
#     if item == 'System':
#         async with httpx.AsyncClient(timeout=setting_timeout) as client:
#             response = await client.get(f"http://{os_spec.restip}:5000/api/serviceOp?cmd={cmd}&name=app")
#             data = response.json()
#             if response.status_code in [400, 401, 500]:
#                 flag = await checkLoginAPI(request)
#                 if not flag:
#                     return {"success": False, "error": "Restful API Login Failed"}
#                 else:
#                     response =  await client.get(f"http://{os_spec.restip}:5000/api/serviceOp?cmd={cmd}&name=app")
#                     data = response.json()
#         async with httpx.AsyncClient(timeout=setting_timeout) as client:
#             response = await client.get(f"http://{os_spec.restip}:5000/api/serviceOp?cmd={cmd}&name=db")
#             data2 = response.json()
#             if 'status' in data2 and data2['status'] == 500:
#                 flag = await checkLoginAPI(request)
#                 if not flag:
#                     return {"success": False, "error": "Restful API Login Failed"}
#                 else:
#                     response = await client.get(f"http://{os_spec.restip}:5000/api/serviceOp?cmd={cmd}&name=db")
#                     data2 = response.json()
#         if data and data2:
#             if cmd == 'restart':
#                 redis_state.client.hset('Service', 'reboot', 1)
#             return {"success": True}
#         else:
#             return {"success": False}
#     else:
#         if item == 'DBMS':
#             names = 'db'
#         else:
#             names = 'app'
#         try:
#             async with httpx.AsyncClient(timeout=setting_timeout) as client:
#                 response = await client.get(f"http://{os_spec.restip}:5000/api/serviceOp?cmd={cmd}&name={names}")
#                 data = response.json()
#                 if response.status_code in [400, 401, 500]:
#                     flag = await checkLoginAPI(request)
#                     if not flag:
#                         return {"success": False, "error": "Restful API Login Failed"}
#                     else:
#                         response = await client.get(f"http://{os_spec.restip}:5000/api/serviceOp?cmd={cmd}&name={names}")
#                         data = response.json()
#             if len(data) > 0:
#                 return {"success": True}
#             else:
#                 return {"success": False}
#         except Exception as e:
#             logging.error(f"❌ RestfulAPI Service Error: {e}")
#             return {"success": False}

@router.get('/Reset')
def reset():
    msg = ''

    if not redis_state.client is None:
        redis_state.client.execute_command("SELECT", 0)
        if redis_state.client.hexists("Service","setting"):
            checkflag = redis_state.client.hget("Service","setting")
            if int(checkflag) == 1:
                return {"success": False, "msg": "Modbus setting is activated"}

        setting_path = os.path.join(SETTING_FOLDER, 'setup.json')
        backup_file_path = os.path.join(SETTING_FOLDER, 'setup_backup.json')
        try:
            if os.path.exists(setting_path):
                shutil.copy(setting_path, backup_file_path)  # setup.json을 setting_backup.json으로 백업
                os.remove(setting_path)
            else:
                msg = "No exist setup.json"

            if redis_state.client.hexists("System", "setup"):
                nowSetup = json.loads(redis_state.client.hget("System", "setup"))
                default_file_path = os.path.join(SETTING_FOLDER, 'default.json')

                with open(default_file_path, "r", encoding="utf-8") as f:
                    defaults = json.load(f)

                defaults["mode"] = nowSetup["mode"]
                defaults["General"]["deviceInfo"]["mac_address"] = get_mac_address()
                defaults["General"]["deviceInfo"]["serial_number"] = get_mac_address()

                with open(default_file_path, "w", encoding="utf-8") as f:
                    json.dump(defaults, f, indent=2)  # indent 추가로 가독성 향상

                redis_state.client.hdel("System", "setup")
        except Exception as e:
            print(str(e))
            return {"success": False, "msg": str(e)}
    else:
        msg += "No saved setup in Redis"
    if msg == '':
        return {"success": True, "msg": "Reset success"}
    else:
        return {"success": True, "msg": msg}

@router.get('/ResetAll')
def resetAll():
    msg = ''
    if not redis_state.client is None:
        redis_state.client.execute_command("SELECT", 0)
        if redis_state.client.hexists("Service","setting"):
            checkflag = redis_state.client.hget("Service","setting")
            if int(checkflag) == 1:
                return {"success": False, "msg": "Modbus setting is activated"}
        setting_path = os.path.join(SETTING_FOLDER, 'setup.json')
        db_path = os.path.join(SETTING_FOLDER, 'user.db')
        backup_file_path = os.path.join(SETTING_FOLDER, 'setup_backup.json')
        try:
            if os.path.exists(setting_path):
                shutil.copy(setting_path, backup_file_path)  # setup.json을 setting_backup.json으로 백업
                os.remove(setting_path)
            else:
                msg = "No exist setup.json"
            if os.path.exists(db_path):
                os.remove(db_path)
            else:
                msg += "No exist user.db"

            default_file_path = os.path.join(SETTING_FOLDER, 'default.json')

            with open(default_file_path, "r", encoding="utf-8") as f:
                defaults = json.load(f)

            defaults["mode"] = 'device0'
            defaults["General"]["deviceInfo"]["mac_address"] = ""
            defaults["General"]["deviceInfo"]["serial_number"] = ""

            with open(default_file_path, "w", encoding="utf-8") as f:
                json.dump(defaults, f, indent=2)  # indent 추가로 가독성 향상

            if redis_state.client.hexists("System", "setup"):
                redis_state.client.hdel("System","setup")
            if redis_state.client.hexists("System", "mode"):
                redis_state.client.hdel("System","mode")
        except Exception as e:
            print(str(e))
            return {"success": False, "msg": str(e)}
    else:
        msg += "No saved setup in Redis"
    if msg == '':
        return {"success": True, "msg": "Reset success"}
    else:
        return {"success": True, "msg": msg}

# @router.get('/DownloadBackup/{item}')
# async def get_backup(item):
#     try:
#         async with httpx.AsyncClient(timeout=setting_timeout) as client:
#             response = await client.get(f"http://{os_spec.restip}:5000/api/getFolder?name={item}")
#             response.raise_for_status()  # 에러 발생 시 예외 발생
#             file_path = os.path.join(DOWNLOAD_FOLDER, f'{item}.zip')
#             with open(file_path, "wb") as f:
#                 f.write(response.content)
#         if os.path.exists(file_path) and os.path.getsize(file_path) > 0:
#             response = FileResponse(file_path, media_type="application/json", filename="setup.json")
#             response.headers["Content-Disposition"] = "attachment; filename=setup.json"  # 다운로드 강제
#             response.headers["X-Message"] = "OK"
#             return response
#         else:
#             return JSONResponse(status_code=404, content={"success": False, "message": "파일이 존재하지 않습니다."},headers={"X-Message": "Fail"})
#     except Exception as e:
#         print(str(e))
#         return JSONResponse(status_code=404, content={"success": False, "message": "Restful API Service Error"},headers={"X-Message": "Fail"})

@router.get('/getDiagnosisSetting')
async def get_diagnosissetting(request:Request):
    # response = await  http_state.client.get(f"/getSettings")
    # data = response.json()
    async with httpx.AsyncClient(timeout=setting_timeout) as client:
        response = await client.get(f"http://{os_spec.restip}:5000/api/getSettings")
        data = response.json()
    #     if response.status_code in [400, 401, 500]:
    #         flag = await checkLoginAPI(request)
    #         if not flag:
    #             return {"success": False, "error": "Restful API Login Failed"}
    #         else:
    #             response = await client.get(f"http://{os_spec.restip}:5000/api/getSettings")
    #             data = response.json()
    if len(data) > 0:
        # print(data)
        return {"success": True, "data": data}
    else:
        return {"success": False, "error": "No Data"}

@router.get('/getDiagnosisProfile')
async def get_diagnosisprofile(request:Request):
    # response = await  http_state.client.get(f"/getProfile")
    # data = response.json()
    async with httpx.AsyncClient(timeout=setting_timeout) as client:
        response = await client.get(f"http://{os_spec.restip}:5000/api/getProfile")
        data = response.json()
    #     if response.status_code in [400, 401, 500]:
    #         flag = await checkLoginAPI(request)
    #         if not flag:
    #             return {"success": False, "error": "Restful API Login Failed"}
    #         else:
    #             response = await client.get(f"http://{os_spec.restip}:5000/api/getProfile")
    #             data = response.json()
    if len(data) > 0:
        return {"success": True, "data": data}
    else:
        return {"success": False, "error": "No Data"}

@router.post('/setDiagnosisSetting')
async def set_diagnosissetting(request: Request):
    data = await request.json()
    if not data:
        return {"status": "0", "error": "No data provided"}
    try:
        # response = await  http_state.client.post(f"/setSettings", json=data)
        # data = response.json()
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.post(f"http://{os_spec.restip}:5000/api/setSettings", json=data)
            data = response.json()
        #     if response.status_code in [400, 401, 500]:
        #         flag = await checkLoginAPI(request)
        #         if not flag:
        #             return {"success": False, "error": "Restful API Login Failed"}
        #         else:
        #             response = await client.post(f"http://{os_spec.restip}:5000/api/setSettings", json=data)
        #             data = response.json()
    except Exception as e:
        print("Error:", e)
        return {"status": "0", "error": str(e)}

    if len(data) > 0:
        return {"success": True }
    else:
        return {"success": False, "error": "Save Failed"}

@router.post('/setDiagnosisProfile')
async def set_diagnosisprofile(request: Request):
    data = await request.json()
    if not data:
        return {"status": "0", "error": "No data provided"}
    try:
        # response = await  http_state.client.post(f"/setProfile", json=data)
        # data = response.json()
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.post(f"http://{os_spec.restip}:5000/api/setProfile", json=data)
            data = response.json()
        #     if response.status_code in [400, 401, 500]:
        #         flag = await checkLoginAPI(request)
        #         if not flag:
        #             return {"success": False, "error": "Restful API Login Failed"}
        #         else:
        #             response = await client.post(f"http://{os_spec.restip}:5000/api/setProfile",json=data)
        #             data = response.json()
    except Exception as e:
        print("Error:", e)
        return {"status": "0", "error": str(e)}

    if len(data) > 0:
        return {"success": True }
    else:
        return {"success": False, "error": "Save Failed"}

@router.get("/getAssetTypes")  # Diagnosis, Report Vue : get Asset info
async def get_assetTypes(request: Request):
    # response = await  http_state.client.get(f"/getAssetTypes")
    # data = response.json()
    async with httpx.AsyncClient(timeout=setting_timeout) as client:
        response = await client.get(f"http://{os_spec.restip}:5000/api/getAssetTypes")
        data = response.json()
    #     if response.status_code in [400, 401, 500]:
    #         flag = await checkLoginAPI(request)
    #         if not flag:
    #             return {"success": False, "error": "Restful API Login Failed"}
    #         else:
    #             response = await client.get(f"http://{os_spec.restip}:5000/api/getAssetTypes")
    #             data = response.json()
    if data:
        return {"success": True, "data": data}
    else:
        return {"success": False, "error": "No Data"}

@router.get('/getAssetList')
async def get_assetlist(request: Request):
    # response = await  http_state.client.get(f"/getAssetHierarchy")
    # datas = response.json()
    async with httpx.AsyncClient(timeout=setting_timeout) as client:
        response = await client.get(f"http://{os_spec.restip}:5000/api/getAssetHierarchy")
        datas = response.json()
    #     if response.status_code in [400, 401, 500]:
    #         flag = await checkLoginAPI(request)
    #         if not flag:
    #             return {"success": False, "error": "Restful API Login Failed"}
    #         else:
    #             response = await client.get(f"http://{os_spec.restip}:5000/api/getAssetHierarchy")
    #             datas = response.json()
    if len(datas) > 0:
        grouped = defaultdict(list)
        for item in datas:
            aid = item["AssemblyID"]
            if aid is not None and aid != "null":
                grouped[aid].append({
                    "Name": item["Name"],
                    "Type": aid
                })
        return {'success': True, "data":grouped}
    else:
        return {'success': False, "msg": 'No exist assets'}

@router.post('/upload')  # upload setup.json
def upload_file(file: UploadFile = File(...)):
    if file.filename == '':
        return {'passOK': '0', 'error': 'No selected file'}

    # 기존 setup.json 파일 경로
    original_file_path = os.path.join(SETTING_FOLDER, 'setup.json')
    backup_file_path = os.path.join(SETTING_FOLDER, 'setup_backup.json')

    try:

        
        # 기존 setup.json 파일이 존재하면 백업
        if os.path.exists(original_file_path):
            shutil.copy(original_file_path, backup_file_path)  # setup.json을 setting_backup.json으로 백업

        # 업로드된 파일을 setup.json으로 저장
        with open(original_file_path, "wb") as buffer:
            shutil.copyfileobj(file.file, buffer)

        file.file.seek(0)  # 다시 읽기 위해 커서 초기화
        file_content = file.file.read().decode("utf-8")

        # Redis에 그대로 문자열로 저장
        redis_state.client.execute_command("SELECT", 0)
        redis_state.client.hset("System", "setup", file_content)
        redis_state.client.hset("Service", "save", 1)
        redis_state.client.hset("Service", "restart", 1)
        return {'passOK': '1', 'file_path': original_file_path}

    except Exception as e:
        return {'passOK': '0', 'error': str(e)}

@router.get('/checkBearing')
def load_and_merge_bearing_data():
    filenames = ["Bearing.csv", "Bearing.json", "Bearings.csv", "Bearings.json"]
    merged_list = []
    seen_names = set()
    nonExistCount = 0
    for fname in filenames:
        file_path = os.path.join(SETTING_FOLDER, fname)
        if not os.path.isfile(file_path):
            nonExistCount += 1
            continue

        if fname.endswith(".csv"):
            try:
                with open(file_path, newline='', encoding='utf-8') as f:
                    reader = csv.DictReader(f)
                    for row in reader:
                        name = row.get("Name")
                        if name and name not in seen_names:
                            seen_names.add(name)
                            merged_list.append(row)
            except Exception as e:
                return {'passOK': '0', 'msg': f"[CSV Error] {fname}: {e}"}

        elif fname.endswith(".json"):
            try:
                with open(file_path, "r", encoding="utf-8") as f:
                    data = json.load(f)
                    if not isinstance(data, list):
                        print(f"[JSON Warning] {fname}는 리스트 형식이 아님, 무시됨")
                        continue
                    for item in data:
                        name = item.get("Name")
                        if name and name not in seen_names:
                            seen_names.add(name)
                            merged_list.append(item)
            except Exception as e:
                return {'passOK': '0', 'msg': f"[CSV JSON] {fname}: {e}"}

    if nonExistCount == len(filenames):
        return {'passOK': '0', 'msg': f"No exist files"}
    else:
        return {'passOK': '1', 'data':merged_list}

@router.post('/uploadBearing')  # upload setup.json
def upload_file(file: UploadFile = File(...)):
    if file.filename == '':
        return {'passOK': '0', 'error': 'No selected file'}

    # 기존 setup.json 파일 경로
    original_file_path = os.path.join(SETTING_FOLDER, file.filename)

    try:

        with open(original_file_path, "wb") as buffer:
            shutil.copyfileobj(file.file, buffer)

        data = get_Bearing(original_file_path)

        return {'passOK': '1', 'file_path': original_file_path, 'data':data}

    except Exception as e:
        return {'passOK': '0', 'error': str(e)}


@router.get('/download')  # download setup.json
def download_setting():
    file_path = os.path.join(SETTING_FOLDER, 'setup.json')

    # 파일이 존재하지 않으면 에러 처리
    if not os.path.exists(file_path):
        return {"passOK": "0", "error": "setting file not found"}

    try:
        # 파일을 클라이언트로 다운로드 응답
        response = FileResponse(file_path, media_type="application/json", filename="setup.json")
        response.headers["Content-Disposition"] = "attachment; filename=setup.json"  # 다운로드 강제
        return response
    except Exception as e:
        return {"passOK": "0", "error": f"An unexpected error occurred: {str(e)}"}


@router.post('/createAsset')
async def createAsset(request: Request):
    data = await request.json()

    if not data:
        return {"status": "0", "error": "No data provided"}
    assetType = data['assetType']
    assetName = data['assetName']
    print("assetType",assetType,"assetName",assetName)
    try:
        # response = await  http_state.client.get(f"/createAsset?type={assetType}&name={assetName}")
        # datas = response.json()
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.get(f"http://{os_spec.restip}:5000/api/createAsset?type={assetType}&name={assetName}")
            datas = response.json()
        #     if response.status_code in [400, 401, 500]:
        #         flag = await checkLoginAPI(request)
        #         if not flag:
        #             return {"success": False, "error": "Restful API Login Failed"}
        #         else:
        #             response = await client.get(f"http://{os_spec.restip}:5000/api/createAsset?type={assetType}&name={assetName}")
        #             datas = response.json()
    except Exception as e:
        print(str(e))
        return {"status": "0", "error": "Created Failed"}

    if len(datas) > 0:
        if assetName in datas:
            return {"status": "1"}
        else:
            return {"status": "0", "error": "Create Failed"} #성공해도 리스트에는 추가안된상태로 리턴되서 여기로빠짐
    else:
        return {"status": "0", "error": "Create Failed"}

@router.post('/modifyAsset')
async def modifyAsset(request: Request):
    data = await request.json()
    if not data:
        return {"status": "0", "error": "No data provided"}
    assetName = data.get("assetName") 
    newName = data.get("newName")
    try:
        # response = await  http_state.client.get(f"/renameAsset?name={assetName}&newname={newName}")
        # datas = response.json()
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.get(f"http://{os_spec.restip}:5000/api/renameAsset?name={assetName}&newname={newName}")
            datas = response.json()
        #     if response.status_code in [400, 401, 500]:
        #         flag = await checkLoginAPI(request)
        #         if not flag:
        #             return {"success": False, "error": "Restful API Login Failed"}
        #         else:
        #             response = await client.get(f"http://{os_spec.restip}:5000/api/renameAsset?name={assetName}&newname={newName}")
        #             datas = response.json()
    except Exception as e:
        return {"status": "0", "error": "Modify Failed"}

    if len(datas) > 0:     
        if newName in datas:
            return {"status": "1"}
        else:
            return {"status": "0", "error": "Modify Failed"}
    else:
        return {"status": "0", "error": "Modify Failed"}

@router.get('/deleteAsset/{asset}')
async def deleteAsset(asset, request: Request):
    try:
        # response = await  http_state.client.get(f"/deleteAsset?name={asset}")
        # datas = response.json()
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.get(f"http://{os_spec.restip}:5000/api/deleteAsset?name={asset}")
            datas = response.json()
        #     if response.status_code in [400, 401, 500]:
        #         flag = await checkLoginAPI(request)
        #         if not flag:
        #             return {"success": False, "error": "Restful API Login Failed"}
        #         else:
        #             response = await client.get(
        #                 f"http://{os_spec.restip}:5000/api/deleteAsset?name={asset}")
        #             datas = response.json()
    except Exception as e:
        print(str(e))
        return {"status": "0", "error": "Delete Failed"}

    if len(datas) > 0:
        if not asset in datas:
            return {"status": "1"}
        else:
            return {"status": "0", "error": "Delete Failed"}
    else:
        return {"status": "1"}

@router.get('/unregisterAsset/{channel}/{asset}')
async def unreg_Asset(channel, asset):
    # finflag = False
    # response = await  http_state.client.get(f"/unregisterAsset?name={asset}")
    # data = response.json()
    async with httpx.AsyncClient(timeout=setting_timeout) as client:
        response = await client.get(f"http://{os_spec.restip}:5000/api/unregisterAsset?name={asset}")
        data = response.json()
    if len(data) > 0:
        finflag = True
    else:
        finflag = False
    redis_state.client.execute_command("SELECT", 0)
    if redis_state.client.hexists("System","setup"):
        datStr = redis_state.client.hget("System","setup")
        setting = json.loads(datStr)
        for chInfo in setting["channel"]:
            if channel == chInfo["channel"]:
                chInfo["assetInfo"]["name"] =''
                chInfo["assetInfo"]["type"] = ''
                finflag = True
                break
        redis_state.client.hset("System","setup", json.dumps(setting))
    if finflag:
        return {"success":True}
    else:
        return {"success":False}

@router.get('/registerAsset/{channel}/{assetName}/{assetType}')
async def reg_Asset(channel, assetName, assetType):
    # finflag = False
    # response = await  http_state.client.get(f"/registerAsset?name={assetName}")
    # data = response.json()
    async with httpx.AsyncClient(timeout=setting_timeout) as client:
        response = await client.get(f"http://{os_spec.restip}:5000/api/registerAsset?name={assetName}")
        data = response.json()
    if len(data) > 0:
        finflag = True
    else:
        finflag = False
    redis_state.client.execute_command("SELECT", 0)
    if redis_state.client.hexists("System","setup"):
        datStr = redis_state.client.hget("System","setup")
        setting = json.loads(datStr)
        for chInfo in setting["channel"]:
            if channel == chInfo["channel"]:
                chInfo["assetInfo"]["name"] = assetName
                chInfo["assetInfo"]["type"] = assetType
                finflag = True
                break
        redis_state.client.hset("System","setup", json.dumps(setting))
    if finflag:
        return {"success":True}
    else:
        return {"success":False}

def initialize_alarm_configs(channel, alams):
    rediskey = f"alarm_status:{channel}"
    alarmdict = {}
    for j in range(1, 33):
        if alams[str(j)][0] != 0:
            alarmdict[str(j)] = alams[str(j)]
    if len(alarmdict.keys()) > 0:
        redis_state.client.execute_command("SELECT", 1)
        redis_state.client.delete(rediskey)
        for key in alarmdict:
            if alarmdict[key][1] == 0:
                condstr = 'UNDER'
            else:
                condstr = 'OVER'
            init_data = {
                "status": 0,
                "count": 0,
                "condition": condstr,
                "value": 0,
                "chan": alarmdict[key][0],
                "cond": alarmdict[key][1],
                "level": alarmdict[key][3],
                "chan_text": parameter_options[alarmdict[key][0]],
                "last_update": int(datetime.now().timestamp()),
                "status_text": "None"
            }
            redis_state.client.hset(rediskey, key, json.dumps(init_data))

@router.post('/savefile/{channel}')  # save setup.json
async def saveSetting(channel: str, request: Request):
    deviceMac = get_mac_address()
    if os_spec.os != 'Windows':
        mac_file_path = os.path.join(SETTING_FOLDER, 'serial_num_do_not_modify.txt')
        if os.path.exists(mac_file_path):
            ser = read_mac_plain(mac_file_path)
            if ser != deviceMac:
                deviceMac = ser
    redis_state.client.execute_command("SELECT", 0)
    if redis_state.client.hexists("Service","setting"):
        checkflag = redis_state.client.hget("Service","setting")
        if int(checkflag) == 1:
            return {"status": "0", "error": "Modbus setting is activated"}
    try:
        FILE_PATH = os.path.join(SETTING_FOLDER, "setup.json")
        #FTP_PATH = os.path.join(SETTING_FOLDER, "WaveformFTP.json")
        # 클라이언트로부터 JSON 데이터 받기
        data = await request.json()
        if not data:
            return {"status": "0", "error": "No data provided"}

        # 파일이 존재하면 읽기, 없으면 기본 구조 생성
        if os.path.exists(FILE_PATH):
            with open(FILE_PATH, "r", encoding="utf-8") as f:
                setting = json.load(f)
        else:
            # 기본 JSON 구조 (필요에 따라 수정)
            setting = {"mode": "device", "General": {}, "channel": []}

        # 채널 파라미터가 "general"이면 "General" 부분 업데이트
        if channel.lower() == "general":        
            setting["General"] = data
         
        else:
            # 그 외의 경우, "channel" 배열 내에서 해당 채널 업데이트 (없으면 추가)
            if "channel" not in setting or not isinstance(setting["channel"], list):
                setting["channel"] = []
            updated = False
            for idx, ch in enumerate(setting["channel"]):
                if ch.get("channel", "").lower() == channel.lower():
                    setting["channel"][idx] = data
                    updated = True
                    break
            if not updated:
                setting["channel"].append(data)

        # 업데이트된 설정을 파일에 저장 (덮어쓰기)
        if deviceMac != setting["General"]["deviceInfo"]["mac_address"]:
            setting["General"]["deviceInfo"]["mac_address"] = deviceMac
            setting["General"]["deviceInfo"]["serial_number"] = deviceMac

        with open(FILE_PATH, "w", encoding="utf-8") as f:
            json.dump(setting, f, indent=4)

        if len(setting["channel"]) > 0:
            for i in range(0,len(setting["channel"])):
                initialize_alarm_configs(setting["channel"][i]["channel"], setting["channel"][i]["alarm"])

        redis_state.client.select(0)
        redis_state.client.hset("System", "setup", json.dumps(setting))
        # redis_state.client.hset("Service","save",1)
        # redis_state.client.hset("Service", "restart", 1)

        return {"status": "1", "data": setting}
    except Exception as e:
        print("Error:", e)
        return {"status": "0", "error": str(e)}
    
@router.get('/restartasset/{flag}')  # save setup.json
async def restartasset(request:Request, flag):
    redis_state.client.select(0)
    if redis_state.client.hexists("Service","setting"):
        checkflag = redis_state.client.hget("Service","setting")
        if int(checkflag) == 1:
            return {"status": "0","success": False, "error": "Modbus setting is activated"}
    try:
        redis_state.client.hset("Service","save",1)
        redis_state.client.hset("Service", "restart", 1)
        if flag:
            # response = await  http_state.client.get(f"/isServiceRestartNeeded")
            # data = response.json()
            if is_service_active("smartsystemsservice") and is_service_active("smartsystemsrestapiservice"):
                async with httpx.AsyncClient(timeout=setting_timeout) as client:
                    response = await client.get(f"http://{os_spec.restip}:5000/api/isServiceRestartNeeded")
                    data = response.json()
                    # if response.status_code in [400, 401, 500]:
                    #     flag = await checkLoginAPI(request)
                    #     if not flag:
                    #         return {"success": False, "error": "Restful API Login Failed"}
                    #     else:
                    #         response = await client.get(f"http://{os_spec.restip}:5000/api/isServiceRestartNeeded")
                    #         data = response.json()
                if data['ServiceRestartRequired']:
                    return {"status": "1", "success": True}
                else:
                    return {"status": "1", "success": False}
            else:
                return {"status": "1", "success": False}
        else:
            return {"status": "1", "success": False}
    except Exception as e:
        print("Error:", e)
        return {"status": "0", "success": False}

@router.post('/restoreSetting') #restore setup.json
def restore_setting():
    if redis_state.client.hexists("Service","setting"):
        checkflag = redis_state.client.hget("Service","setting")
        if int(checkflag) == 1:
            return {"status": "0", "error": "Modbus setting is activated"}

    setting_file = os.path.join(SETTING_FOLDER, 'setup.json')
    backup_file = os.path.join(SETTING_FOLDER, 'setup_backup.json')

    # 백업 파일이 존재하는지 확인
    if not os.path.exists(backup_file):
        return JSONResponse(status_code=400, content={"passOK": "0", "error": "setting_backup.json not found"})

    try:
        # 백업 파일을 setup.json으로 덮어쓰기
        shutil.copy(backup_file, setting_file)
        with open(backup_file, "r", encoding="utf-8") as f:
            setting = json.load(f)
        redis_state.client.execute_command("SELECT", 0)
        redis_state.client.hset("System", "setup", json.dumps(setting))
        redis_state.client.hset("Service", "save", 1)
        redis_state.client.hset("Service", "restart", 1)
        return {"passOK": "1", "message": "Restore successful"}
    except Exception as e:
        return JSONResponse(status_code=500, content={"passOK": "0", "error": f"Error occurred while restoring: {str(e)}"})


# def findNode(node, superlist): #check node level
#     exist =False
#     if node["Assembly"] in superlist:
#     for i in range(0, len(superlist)):
#         if superlist[i]["Assembly"] == node["Assembly"]:
#             exist = True
#             break
#     if not exist:
#         superlist.append(node)
#     return superlist

# def set_structNameplate(datas:list):
#     if len(datas) > 0:
#         superlist = []
#         superNodes = []
#         for a in range(0, len(datas)):
#             pName = f"{datas[a]['Assembly']}_{datas[a]['AssemblyID']}"
#             if not pName in superNodes:
#                 superNodes.append(pName)
#
#         for i in range(0, len(superNodes)):
#             superdict = dict()
#             tmpstrs = superNodes[i].split('_')
#             superdict["Assembly"] = tmpstrs[0]
#             superdict["Name"] = superNodes[i]
#             superdict["Title"] = superNodes[i]
#             if len(tmpstrs) > 1:
#                 superdict["AssemblyID"] = tmpstrs[1]
#             else:
#                 superdict["AssemblyID"] = tmpstrs[0]
#             superdict["ParentAssemly"] = superNodes[i]
#             superdict["Abstract"] = True
#             superlist.append(superdict)
#
#         for i in range(0,len(superlist)):
#             for j in range(0, len(datas)):
#                 pName = f"{datas[j]['Assembly']}_{datas[j]['AssemblyID']}"
#                 if superlist[i]["ParentAssemly"] == pName:
#                     if "children" in superlist[i]:
#                         superlist[i]["children"].append(datas[j])
#                     else:
#                         superlist[i]["children"] = []
#                         superlist[i]["children"].append(datas[j])
#     return superlist

def get_Bearing(filename):
    _, ext = os.path.splitext(filename)

    if ext.lower() != ".csv":
        return json.loads(filename)
    else:
        parsed_result = []
        with open(filename, "r", encoding="utf-8") as f:
            reader = csv.DictReader(f)
            for row in reader:
                # 문자열 → float 변환 (Name은 그대로)
                for key in row:
                    if key != "Name":
                        row[key] = float(row[key])
                parsed_result.append(row)
        return parsed_result

@router.get("/getAssetConfig/{asset}")  #Setting Vue : get Asset info
async def get_assetconfig(asset, request:Request):
    # response = await  http_state.client.get(f"/getNameplate?name="+asset)
    # data = response.json()
    async with httpx.AsyncClient(timeout=setting_timeout) as client:
        response = await client.get(f"http://{os_spec.restip}:5000/api/getNameplate?name="+asset)
        data = response.json()
    #     if response.status_code in [400, 401, 500]:
    #         flag = await checkLoginAPI(request)
    #         if not flag:
    #             return {"success": False, "error": "Restful API Login Failed"}
    #         else:
    #             response = await client.get(f"http://{os_spec.restip}:5000/api/getNameplate?name={asset}")
    #             data = response.json()
    #     superlist = set_structNameplate(data)
    if len(data) > 0:
        return {"success": True, "data": data}
    else:
        return {"success": False, "error": "No Data"}

@router.post("/checkAssetConfig/{asset}")
async def check_assetconfig(asset:str, request:Request):
    data = await request.json()
    if not data:
        return {"status": "0", "error": "No data provided"}
    try:
        # response = await  http_state.client.post(f"/checkNameplate?name={asset}", json=data)
        # result = response.json()
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.post(f"http://{os_spec.restip}:5000/api/checkNameplate?name={asset}", json=data)
            result = response.json()
        #     if 'status' in result and result['status'] == 500:
        #         flag = await checkLoginAPI(request)
        #         if not flag:
        #             return {"success": False, "error": "Restful API Login Failed"}
        #         else:
        #             response = await client.post(f"http://{os_spec.restip}:5000/api/setNameplate?name={asset}", json=data)
        #             result = response.json()
    except Exception as e:
        print("Error:", e)
        return {"status": "0", "error": str(e)}

    if len(result) > 0:
        return {"success":True, "result":result['resetRequired']}
    else:
        return {"success": False, "error": "No Data"}

@router.post("/setAssetConfig/{asset}")
async def set_assetconfig(asset:str, request:Request):
    data = await request.json()
    if not data:
        return {"status": "0", "error": "No data provided"}
    try:
        # response = await  http_state.client.post(f"/setNameplate?name={asset}", json=data)
        # result = response.json()
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.post(f"http://{os_spec.restip}:5000/api/setNameplate?name={asset}", json=data)
            result = response.json()
        #     if 'status' in result and result['status'] == 500:
        #         flag = await checkLoginAPI(request)
        #         if not flag:
        #             return {"success": False, "error": "Restful API Login Failed"}
        #         else:
        #             response = await client.post(f"http://{os_spec.restip}:5000/api/setNameplate?name={asset}", json=data)
        #             result = response.json()
    except Exception as e:
        print("Error:", e)
        return {"status": "0", "error": str(e)}

    if len(result) > 0:
        return {"success": True}
    else:
        return {"success": False, "error": "No Data"}

@router.get("/getAssetParams/{asset}")  #Setting Vue : get Asset info
async def get_assetParams(asset, request:Request):
    # response = await  http_state.client.get(f"/getParameters?name="+asset)
    # data = response.json()
    async with httpx.AsyncClient(timeout=setting_timeout) as client:
        response = await client.get(f"http://{os_spec.restip}:5000/api/getParameters?name="+asset)
        data = response.json()
    #     if response.status_code in [400, 401, 500]:
    #         flag = await checkLoginAPI(request)
    #         if not flag:
    #             return {"success": False, "error": "Restful API Login Failed"}
    #         else:
    #             response = await client.get(f"http://{os_spec.restip}:5000/api/getParameters?name={asset}")
    #             data = response.json()
    if len(data) > 0:
        return {"success": True, "data": data }
    else:
        return {"success": False, "error": "No Data"}

@router.post("/setAssetParams/{asset}")
async def set_assetParams(asset:str, request:Request):
    data = await request.json()
    if not data:
        return {"status": "0", "error": "No data provided"}
    try:
        # response = await  http_state.client.post(f"/setParameters?name={asset}", json=data)
        # result = response.json()
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.post(f"http://{os_spec.restip}:5000/api/setParameters?name={asset}", json=data)
            result = response.json()
        #     if 'status' in result and result['status'] == 500:
        #         flag = await checkLoginAPI(request)
        #         if not flag:
        #             return {"success": False, "error": "Restful API Login Failed"}
        #         else:
        #             response = await client.post(f"http://{os_spec.restip}:5000/api/setParameters?name={asset}", json=data)
        #             result = response.json()
    except Exception as e:
        print("Error:", e)
        return {"status": "0", "error": str(e)}

    if len(result) > 0:
        return {"success": True}
    else:
        return {"success": False, "error": "No Data"}

@router.get("/test/{asset}")
async def test_asset(asset, request:Request):
    try:
        # HTTP 클라이언트 상태 확인
        # if http_state.error:
        #     return {"success": False, "error": http_state.error}
        #
        # # 클라이언트 가져오기
        # client = http_state.client
        # if not client:
        #     return {"success": False, "error": "HTTP client not initialized"}

        # response = await http_state.client.get(f"/getComm?name={asset}")
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.get(f"http://{os_spec.restip}:5000/api/getComm?name={asset}")
            data = response.json()

        result = dict()
        result["AssetName"] = data["AssetName"]
        result["SerialNumber"] = data["SerialNumber"]
        result["Channel"] = data["Channel"]
        result["Commissions"] = data["Commissions"]

        if len(data) > 0:
            return {"success": True, "data": result}
        else:
            return {"success": False, "error": "No Data"}
    except Exception as e:
        print(f"Exception: {type(e).__name__}: {e}")
        return {"success": False, "error": "No Data"}

@router.get("/testwave/{asset}")
async def test_asset(asset, request:Request):
    try:
        # HTTP 클라이언트 상태 확인
        # if http_state.error:
        #     return {"success": False, "error": http_state.error}
        #
        # # 클라이언트 가져오기
        # client = http_state.client
        # if not client:
        #     return {"success": False, "error": "HTTP client not initialized"}
        #
        # response = await http_state.client.get(f"/getWaveform?name={asset}")
        # data = response.json()
        async with httpx.AsyncClient(timeout=setting_timeout) as client:
            response = await client.get(f"http://{os_spec.restip}:5000/api/getWaveform?name={asset}")
            data = response.json()

        if len(data) > 0:
            waveTDict = dict()
            waveTDict["SR"] = data["SR"]
            waveTDict["Duration"] = data["Duration"]
            waveTDict["Vwave"] = data["Vwave"]
            waveTDict["Iwave"] = data["Iwave"]
            waveTDict["Vfft"] = data["Vfft"]
            waveTDict["Ifft"] = data["Ifft"]
            return {"success": True, "waveT": waveTDict}
        else:
            return {"success": False, "error": "No Data"}
    except Exception as e:
        print(f"Exception: {type(e).__name__}: {e}")
        return {"success": False, "error": "No Data"}

# @router.get("/healthCheck")
# async def check_health(request:Request):
#     async with httpx.AsyncClient(timeout=setting_timeout) as client:
#         response = await client.get(f"http://{os_spec.restip}:5000/api/health")
#         data = response.json()
#         if response.status_code in [400, 401, 500]:
#             flag = await checkLoginAPI(request)
#             if not flag:
#                 return {"success": False, "error": "Restful API Login Failed"}
#             else:
#                 response = await client.get(f"http://{os_spec.restip}:5000/api/health")
#                 data = response.json()
#     if len(data) > 0:
#         return {"success":True, "data":data}
#     else:
#         return {"success":False, "msg":"No Data"}
#
# @router.get("/checkAPI")
# async def check_health(request:Request):
#     async with httpx.AsyncClient(timeout=setting_timeout) as client:
#         response = await client.get(f"http://{os_spec.restip}:5000/api/alive")
#         data = response.json()
#         if response.status_code in [400, 401, 500]:
#             flag = await checkLoginAPI(request)
#             if not flag:
#                 return {"success": False, "error": "Restful API Login Failed"}
#             else:
#                 response = await client.get(f"http://{os_spec.restip}:5000/api/alive")
#                 data = response.json()
#     if len(data) > 0:
#         return {"success":True, "data":data}
#     else:
#         return {"success":False, "msg":"No Data"}

def is_service_active(name):
    try:
        result = subprocess.run(['systemctl', 'is-active', name],
                                stdout=subprocess.PIPE,
                                stderr=subprocess.PIPE,
                                text=True,
                                timeout=2)
        return result.stdout.strip() == "active"
    except Exception as e:
        print(f"❌ 서비스 상태 확인 실패: {name} - {e}")
        return False

@router.get("/SysCheck")
def check_sysStatus():
    data = {}

    if os_spec.os == 'Windows':
        if not redis_state.client is None:
            redis_state.client.execute_command("SELECT", 0)
            if redis_state.client.hexists("System", "Status"):
                data = json.loads(redis_state.client.hget("System", "Status"))
        data["core"] = 0

        usage = psutil.disk_usage('/')
        disk1 = {
            "drive": "/",
            "totalGB": round(usage.total / (1024 ** 3), 2),
            "freeGB": round(usage.free / (1024 ** 3), 2),
            "status": "ok" if usage.percent < 90 else "warning"
        }

        return {"success": True, "data": data, "disk": [disk1]}
    else:
        servicedict = {
            'smartsystem' : 'smartsystemsservice',
            'smartapi': 'smartsystemsrestapiservice',
            'redis': 'redis',
            'influxdb': 'influxdb',
            'core': 'core',
            'webserver':'webserver'
        }
        service_status = {}

        # 각 서비스의 상태 확인
        for key, service_name in servicedict.items():
            service_status[key] = is_service_active(service_name)

        usage = psutil.disk_usage('/')
        disk1 = {
            "drive": "/",
            "totalGB": round(usage.total / (1024 ** 3), 2),
            "freeGB": round(usage.free / (1024 ** 3), 2),
            "status": "ok" if usage.percent < 90 else "warning"
        }

        usage = psutil.disk_usage('/usr/local')
        disk2 = {
            "drive": "/usr/local",
            "totalGB": round(usage.total / (1024 ** 3), 2),
            "freeGB": round(usage.free / (1024 ** 3), 2),
            "status": "ok" if usage.percent < 90 else "warning"
        }

        return {"success": True, "data":service_status,  "disk":[disk1, disk2]}

@router.get('/SysService/{cmd}/{item}')
def sysService(cmd, item):
    itemdict = {
      "Redis":"redis",
      "InfluxDB":"influxdb",
      "SmartSystems":"smartsystemsservice",
      "SmartAPI":"smartsystemsrestapiservice",
      "Core":"core",
      "WebServer":"webserver",
    }

    try:
        # systemctl 명령 실행
        result = subprocess.run(
            ['sudo', 'systemctl', cmd, itemdict[item]],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            timeout=30
        )

        return {
            'success': True,
            'stdout': result.stdout.strip(),
            'stderr': result.stderr.strip(),
            'returncode': result.returncode,
            'service': itemdict[item],
            'action': cmd
        }

    except subprocess.TimeoutExpired:
        return {
            'success': False,
            'error': f'Timeout expired while {cmd}ing {itemdict[item]}',
            'service': itemdict[item],
            'action': cmd
        }
    except Exception as e:
        return {
            'success': False,
            'error': str(e),
            'service': itemdict[item],
            'action': cmd
        }

@router.post("/command")
async def push_command_left(command: Command):
    """새로운 command를 Redis 리스트의 왼쪽에 추가"""
    try:
        binary_data = command_to_binary(command)
        redis_state.client.select(1)
        redis_state.client.lpush('command', binary_data)

        return {
            "success": True,
            "message": "Command pushed to left of list",
            "data": command.dict()
        }
    except Exception as e:
        return {
            "success": False,
            "message": str(e)
        }

#@router.post("/command")
#async def push_command_left(request: Request):
#    """새로운 command를 Redis 리스트의 왼쪽에 추가"""
#    try:
#        # Raw body 먼저 확인
#        body_bytes = await request.body()
#        print(f"Raw body: {body_bytes}")
#
#        # JSON 파싱 시도
#        try:
#            body_json = json.loads(body_bytes)
#            print(f"Parsed JSON: {body_json}")
#            print(f"Types: type={type(body_json.get('type'))}, cmd={type(body_json.get('cmd'))}, item={type(body_json.get('item'))}")
#        except json.JSONDecodeError as e:
#            print(f"JSON decode error: {e}")
#            return {"success": False, "message": f"Invalid JSON: {e}"}
#
#        # 수동으로 Command 객체 생성
#        try:
#            command = Command(
#                type=body_json.get('type'),
#                cmd=body_json.get('cmd'),
#                item=body_json.get('item')
#            )
#        except ValidationError as e:
#            print(f"Validation error: {e.errors()}")
#            return {"success": False, "message": "Validation failed", "errors": e.errors()}
#
#        binary_data = command_to_binary(command)
#        redis_state.client.select(1)
#        redis_state.client.lpush('command', binary_data)
#
#        return {
#            "success": True,
#            "message": "Command pushed to left of list",
#            "data": command.dict()
#        }
#    except Exception as e:
#        print(f"Unexpected error: {e}")
#        import traceback
#        traceback.print_exc()
#        return {"success": False, "message": str(e)}
