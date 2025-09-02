from fastapi import APIRouter
import redis, httpx
from pydantic import BaseModel
import sqlite3
import ujson as json

router = APIRouter()

# Path 객체 절대경로
from pathlib import Path
base_dir = Path(__file__).resolve().parent
SETTING_FOLDER = base_dir.parent.parent / "config"  # ⬅️ 두 단계 상위로
router = APIRouter()
# os.makedirs(SETTING_FOLDER, exist_ok=True)  # 폴더가 없으면 생성

DB_PATH = "device.db"

def get_db_connection():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row  # → dict처럼 row 접근 가능
    return conn

# Pydantic 모델
class Device(BaseModel):
    id:str
    ipaddr: str
    location: str
    enable: str

@router.get('/getSetup/{ip}')
async def getSetup(ip):
    async with httpx.AsyncClient() as client:
        response = await client.get(f"http://{ip}:4000/setting/checkSettingFile")
        try:
            data = response.json()  # JSON 데이터 파싱
            return {"result" : 1, "data" : data}
        except ValueError:
            return {"result": 0}

@router.post('/saveDevice')
def save_device(device: Device):
    try:
        conn = get_db_connection()
        cursor = conn.cursor()
        cursor.execute('''
        CREATE TABLE IF NOT EXISTS device (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            ipaddr TEXT NOT NULL,
            location TEXT NOT NULL,
            enable TEXT NOT NULL
        )
        ''')

        conn.commit()
        if device.id == "0":
            cursor.execute(
                "INSERT INTO device (ipaddr, location, enable) VALUES (?, ?, ?)",
                (device.ipaddr, device.location, device.enable)
            )
        else:
            cursor.execute(
                "update device set ipaddr=?, location=?, enable=? where id=?",
                (device.ipaddr, device.location, device.enable, int(device.id))
            )
        conn.commit()
        conn.close()
        return {"status": "1"}
    except Exception as e:
        print(str(e))
        return {"status":"0"}

# @router.post('/saveDevice')
# async def save_device(request:Request):
#     try:
#         FILE_PATH = os.path.join(SETTING_FOLDER, "device.json")
#
#         # 클라이언트로부터 JSON 데이터 받기
#         data = await request.json()
#         print(data)
#         if not data:
#             return {"status": "0", "error": "No data provided"}
#         if os.path.exists(FILE_PATH):
#             with open(FILE_PATH, "r", encoding="utf-8") as f:
#                 setting = json.load(f)
#             idx = -1
#             for i in range(0, len(setting["device"])):
#                 if int(data["id"]) == int(setting["device"][i]["id"]):
#                     idx = i
#                     break
#             if idx >= 0:
#                 setting["device"][idx]["ipaddr"] = data["ipaddr"]
#                 setting["device"][idx]["location"] = data["location"]
#                 setting["device"][idx]["enable"] = data["enable"]
#             else:
#                 setting["device"].append({"id":data["id"] ,"ip": data["ipaddr"],"location":data["location"], "enable":data["enable"]})
#         else:
#             setting = {"device": []}
#             setting["device"].append({"id":int(data["id"]) ,"ip": data["ipaddr"],"location":data["location"], "enable":data["enable"]})
#         with open(FILE_PATH, "w", encoding="utf-8") as f:
#             json.dump(setting, f, indent=4)
#         return {"status": "1"}
#
#     except Exception as e:
#         return {"status":"0"}

@router.get('/deleteDevice/{id}')
def delete_device(id):
#     try:
#         FILE_PATH = os.path.join(SETTING_FOLDER, "device.json")
#         if os.path.exists(FILE_PATH):
#             with open(FILE_PATH, "r", encoding="utf-8") as f:
#                 setting = json.load(f)
#             idx = -1
#             for i in range(0, len(setting["device"])):
#                 if int(id) == int(setting["device"][i]["id"]):
#                     idx = i
#                     break
#             if idx >= 0:
#                 setting["device"].pop(idx)
#                 with open(FILE_PATH, "w", encoding="utf-8") as f:
#                     json.dump(setting, f, indent=4)
#                 return {"status": "1"}
#             else:
#                 return {"status":"0"}
#     except Exception as e:
#         return {"status":"0"}
#
# @router.post('/saveDevice')
# def save_device(device: Device):
    try:
        conn = get_db_connection()
        cursor = conn.cursor()
        cursor.execute(
            "delete from device where id=?",
            (int(id),)
        )
        conn.commit()
        conn.close()
        return {"status": "1"}
    except Exception as e:
        print(str(e))
        return {"status":"0"}

@router.get("/getStatus/{ip}/{asset}") #Master Dashboard Status
async def getStatus(ip, asset):
    async with httpx.AsyncClient() as client:
        response = await client.get(f"http://{ip}:5000/api/getDiagnostic?name={asset}")

        try:
            data = response.json()  # JSON 데이터 파싱
        except ValueError:
            return {"status": -1}

        # success가 True인지 확인
        if len(data) == 0:
            return {"status": -1}

        # "data" 항목이 존재하는지 확인
        status_list = [item["Status"] for item in data["BarGraph"] if "Status" in item]

        # 모든 값이 1이면 "OK"
        if all(status == 1 for status in status_list):
            return {"status": 1}

        # 모든 값이 0이면 "No Data"
        if all(status == 0 for status in status_list):
            return {"status": 0}

        # 우선순위: Repair(4) > Inspect(3) > Warning(2)
        if 4 in status_list:
            return {"status": 4}
        elif 3 in status_list:
            return {"status": 3}
        elif 2 in status_list:
            return {"status": 2}
        # 기본값 (예상하지 못한 경우)
        return {"status": -2}

@router.get("/getDeviceInfo") #Master Dashboard
async def getDeviceInfo():
    try:
        conn = get_db_connection()
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM device")
        rows = cursor.fetchall()
        conn.close()
        devicelist = [dict(row) for row in rows]  # JSON 배열 반환
    except Exception as e:
        return {"result": "0"}

    # file_path = os.path.join(SETTING_FOLDER, 'device.json')
    # try:
    #     with open(file_path, "r", encoding="utf-8") as f:
    #         devices = json.load(f)
    # except Exception as e:
    #     return {"result": "0"}
    # if devices:
    #     devicelist = devices["device"]
    # else:
    #     return {"result": "0"}

    deviceDataList = list()
    for i in range(0, len(devicelist)):  #get device setup (httpx), meter(redis), diagnosis(httpx)
        if devicelist[i]["enable"] == 'true':
            devsetting = await getSetup(devicelist[i]["ipaddr"]) #remote
            if devsetting["result"] == 1:
                dev_setup = devsetting["data"]
                tmpPool = redis.ConnectionPool(host=devicelist[i]["ipaddr"], port=6379, db=6)
                tmpclient = redis.Redis(connection_pool=tmpPool)
                #dev_setup = check_setupfile() #local
                main_st = -1
                sub_st = -1
                meters_sub_dict = dict()

                if dev_setup["result"] == "1":      #each device setup
                    meters_main = tmpclient.get("meter_main")
                    meters_main_dict = json.loads(meters_main)
                    if dev_setup["Diag_main"] and dev_setup["assetName_main"] != '':
                        main_st = await getStatus('192.168.1.24', dev_setup["assetName_main"])
                        meters_main_dict["status"] = main_st["status"]
                        meters_main_dict["st_type"] = dev_setup["assetType_main"]
                    else:
                        meters_main_dict["status"] = main_st
                        meters_main_dict["st_type"] = 'None'
                    if dev_setup["enable_sub"]:
                        meters_sub = tmpclient.get("meter_sub")
                        meters_sub_dict = json.loads(meters_sub)
                        if dev_setup["Diag_sub"] and dev_setup["assetName_sub"] != '':
                            sub_st = await getStatus('192.168.1.24', dev_setup["assetName_sub"])
                            meters_sub_dict["status"] = sub_st["status"]
                            meters_sub_dict["st_type"] = dev_setup["assetType_sub"]
                        else:
                            meters_sub_dict["status"] = sub_st
                            meters_sub_dict["st_type"] = 'None'
                    deviceDataList.append({
                        "result" :1,
                        "id": devicelist[i]["id"],
                        "ip" : devicelist[i]["ipaddr"],
                        "location" : devicelist[i]["location"],
                        "enable": devicelist[i]["enable"],
                        "main": meters_main_dict,
                        "sub":meters_sub_dict
                    })
                else:
                    deviceDataList.append({ "result":0 })
                tmpclient.close()
        else:
            deviceDataList.append({
                "result": 1,
                "id": devicelist[i]["id"],
                "ip": devicelist[i]["ipaddr"],
                "location": devicelist[i]["location"],
                "enable": devicelist[i]["enable"],
                "main": {},
                "sub": {}
            })
    if len(deviceDataList) > 0:
        return {"result" : 1, "data": deviceDataList}
    else:
        return { "result":0 }