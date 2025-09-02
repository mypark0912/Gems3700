from fastapi import APIRouter, Request
import bcrypt, os, logging, shutil
from pydantic import BaseModel
import sqlite3, httpx,uuid, psutil, subprocess
import ujson as json

from states.global_state import aesState, INIT_PATH, redis_state, os_spec

router = APIRouter()

LANG_FILES = {
    "ko": "lang_kor.json",
    "en": "lang_eng.json"
}

# Path 객체 절대경로
from pathlib import Path
base_dir = Path(__file__).resolve().parent
SETTING_FOLDER = base_dir.parent.parent / "config"  # ⬅️ 두 단계 상위로

router = APIRouter()
# os.makedirs(SETTING_FOLDER, exist_ok=True)  # 폴더가 없으면 생성

class Login(BaseModel):
    account: str
    password: str
    lang: str

class DefaultUser(BaseModel):
    role: str
    password: str
    newPassword: str

class Resetpassword(BaseModel):
    username: str
    password: str
    email: str

class SignupUser(BaseModel):
    username: str
    account: str
    password: str
    email: str
    role: str

class SignupAdmin(BaseModel):
    devType: int
    username: str
    account: str
    password: str
    email: str
    role: str
    adminPass: str
    lang: str

class User(BaseModel):
    username: str
    password: str
    email: str
    newPass: str
    role: str
    api : str

DB_PATH = os.path.join(SETTING_FOLDER, "user.db")
CAL_PATH = os.path.join(SETTING_FOLDER, "calibration.csv")
INIT_PATH = os.path.join(SETTING_FOLDER, 'influx.json')
SETUP_PATH = os.path.join(SETTING_FOLDER, 'setup.json')

#def get_mac_address():
#    """현재 시스템의 MAC 주소 가져오기"""
#    try:
#         uuid.getnode()로 MAC 주소 가져오기
#        mac_int = uuid.getnode()
#
#         정수를 16진수 문자열로 변환 (12자리)
#        mac_address = f'{mac_int:012x}'
#
#        return mac_address
#    except Exception as e:
#        return f"MAC 주소 가져오기 실패: {e}"

def is_service_active(service_name):
    try:
        result = subprocess.run(['systemctl', 'is-active', service_name],
                                stdout=subprocess.PIPE,
                                stderr=subprocess.PIPE,
                                text=True,
                                timeout=2)
        return result.stdout.strip() == 'active'
    except Exception as e:
        logging.info(f"❌ 서비스 상태 확인 실패: {service_name} - {e}")
        return False

def get_mac_address():
    """지정된 네트워크 카드들의 MAC 주소 가져오기"""

    # 사용할 네트워크 인터페이스 이름들 지정 (우선순위 순으로)
    TARGET_INTERFACES = ['sw0ep', 'end1']  # 필요에 따라 변경 가능

    try:
        # 지정된 인터페이스들을 순서대로 확인
        network_interfaces = psutil.net_if_addrs()

        for interface_name in TARGET_INTERFACES:
            if interface_name in network_interfaces:
                for addr in network_interfaces[interface_name]:
                    if addr.family == psutil.AF_LINK:  # MAC 주소
                        mac = addr.address.replace('-', '').replace(':', '').lower()
                        if mac and mac != '000000000000':
                            return mac

        # 지정된 인터페이스들이 모두 없으면 uuid.getnode() 사용
        mac_int = uuid.getnode()
        mac_address = f'{mac_int:012x}'
        return mac_address

    except Exception as e:
        print(f"MAC 주소 가져오기 실패: {e}")
        # 예외 발생 시에도 uuid.getnode() 사용
        mac_int = uuid.getnode()
        mac_address = f'{mac_int:012x}'
        return mac_address


@router.get('/getInfluxdb')
def get_influxkey():
    enJson = aesState.getInflux()
    if enJson["result"]:
        # plaintext_dec = aesState.decrypt(enJson["cipher"])
        # print("복호화 결과:", plaintext_dec)
        influxdata = {"token_cipher":enJson["cipher"], "org":enJson["org"]}
        return {"passOK":True, "data":influxdata}
    else:
        return {"passOK":False, "msg" : "No Configurations"}

def get_db_connection():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row  # → dict처럼 row 접근 가능
    return conn

def check_useDiagnosis():
    if not redis_state.client is None:
        redis_state.client.execute_command("SELECT", 0)
        if redis_state.client.hexists("System","setup"):
            redisContext = redis_state.client.hget("System","setup")
            setting = json.loads(redisContext)
            usefile = False
        else:
            usefile = True
    else:
        usefile = True
    if usefile :
        # file_path = os.path.join(SETTING_FOLDER, 'setup.json')
        try:
            with open(SETUP_PATH, "r", encoding="utf-8") as f:
                setting = json.load(f)
        except Exception as e:
            return False

    if not "channel" in setting:
        return False
    opmode = setting.get("mode", "")
    if opmode == 'device1':
        general_data = setting.get("General", {})
        use_fuction = general_data.get("useFuction", {})
        diag_main = bool(use_fuction.get("diagnosis_main", False))
        diag_sub = bool(use_fuction.get("diagnosis_sub", False))
        return diag_main or diag_sub
    else:
        return False

@router.get('/checkInstall')
async def checkInstall():
    db_exists = os.path.exists(DB_PATH)
    conn = get_db_connection()
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    if not db_exists:
        cursor.execute('''
                    CREATE TABLE IF NOT EXISTS user (
                        account TEXT PRIMARY KEY,
                        username TEXT NOT NULL,
                        password TEXT NOT NULL,
                        email TEXT NOT NULL,
                        role TEXT NOT NULL,
                        api TEXT NOT NULL
                    )
                ''')
        cursor.execute("SELECT * FROM user")
        rows = cursor.fetchall()
        if len(rows) > 0:
            flag = True
        else:
            flag = False
    else:
        table_name = "user"
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name=?", (table_name,))
        exists = cursor.fetchone()
        if exists:
            cursor.execute("SELECT * FROM user")
            rows = cursor.fetchall()
            if len(rows) > 0:
                flag = True
            else:
                flag = False
        else:
            flag = True
    conn.commit()
    conn.close()
    if os.path.exists(CAL_PATH):
        calibated = True
    else:
        calibated = False

    if os.path.exists(INIT_PATH):
        installed = True
    else:
        installed = False
    # user.db 없거나 admin 없음 / calibration 없음 / db init 안됨

    if flag and installed:
        if calibated:
            return {"result": 2, "calibration" : True}
        else:
            return {"result": 2, "calibration" : False}
    elif flag and not installed:
        if calibated:
            return {"result": 1, "calibration" : True}
        else:
            return {"result": 1, "calibration" : False}
    else:
        return {"result": 0}

@router.get("/checkDBMS")
def check_Db():
    if os.path.exists(INIT_PATH):
        return {"result": 1}
    else:
        return {"result": 0}

# @router.post('/joinAdmin')
# def join(data: SignupAdmin):
#     devType = data.devType
#     name = data.username
#     account = data.account
#     password = data.password
#     email = data.email
#     adminPass = data.adminPass
#     lang = data.lang
#
#     if devType < 3:
#         mode = f"device{devType}"
#         if devType == 0:
#             diag = 'No'
#         else:
#             diag = 'Yes'
#     else:
#         diag = 'No'
#         mode = "server"
#     # setting = {"mode": mode, "lang": lang ,"General": {"deviceInfo":{"mac_address":get_mac_address(), "serial_number":get_mac_address()}}, "channel": []}
#     default_file_path = os.path.join(SETTING_FOLDER, 'default.json')
#     # setting_path = os.path.join(SETTING_FOLDER, 'setup.json')
#     shutil.copy(default_file_path, SETUP_PATH)
#     with open(SETUP_PATH, "r", encoding="utf-8") as f:
#         setting = json.load(f)
#         setting["mode"] = mode
#         setting["General"]["deviceInfo"]["mac_address"] = get_mac_address()
#         setting["General"]["deviceInfo"]["serial_number"] = get_mac_address()
#     # FILE_PATH = os.path.join(SETTING_FOLDER, "setup.json")
#     with open(SETUP_PATH, "w", encoding="utf-8") as f:
#         json.dump(setting, f, indent=2)
#
#     redis_state.client.execute_command("SELECT", 0)
#     redis_state.client.hset("System", "setup", json.dumps(setting))
#     redis_state.client.hset("System", "mode",  mode)
#     redis_state.client.hset("Service", "save", 1)
#     redis_state.client.hset("Service", "restart", 1)
#     # print(adminPass)
#     # filename = os.path.join(SETTING_FOLDER, 'admin_secret.enc')
#     # with open(filename, 'rb') as f:
#     #     raw = base64.b64decode(f.read())
#     # iv = raw[:16]
#     # encrypted = raw[16:]
#     # cipher = AES.new(AES_KEY, AES.MODE_CBC, iv)
#     # decrypted = unpad(cipher.decrypt(encrypted), AES.block_size).decode('utf-8')
#     # print(decrypted)
#     # if  decrypted == adminPass:
#     if aesState.checkAdmin(adminPass):
#         try:
#             conn = get_db_connection()
#             conn.row_factory = sqlite3.Row
#             cursor = conn.cursor()
#             hashed_password = bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")
#             cursor.execute(
#                 "INSERT INTO `user` (account,username, password, email, role, api) VALUES (?, ?, ?, ?, ?, ?)",
#                 (account, name, hashed_password, email, '3', diag )
#             )
#             conn.commit()
#             conn.close()
#             return {"passOK": "1"}
#         except Exception as e:
#             print(str(e))
#             return {"passOK": "0", "msg" : str(e)}
#     else:
#         return {"passOK": "0", "msg" : 'Admin Password is Wrong'}
#

@router.post('/joinAdmin')
def join(data: SignupAdmin):
    devType = data.devType
    name = data.username
    account = data.account
    password = data.password
    email = data.email
    adminPass = data.adminPass
    lang = data.lang

    if devType < 3:
        mode = f"device{devType}"
        if devType == 0:
            diag = 'No'
        else:
            diag = 'Yes'
    else:
        diag = 'No'
        mode = "server"

    default_file_path = os.path.join(SETTING_FOLDER, 'default.json')
    shutil.copy(default_file_path, SETUP_PATH)
    with open(SETUP_PATH, "r", encoding="utf-8") as f:
        setting = json.load(f)
        setting["mode"] = mode
        setting["General"]["deviceInfo"]["mac_address"] = get_mac_address()
        setting["General"]["deviceInfo"]["serial_number"] = get_mac_address()

    with open(SETUP_PATH, "w", encoding="utf-8") as ef:
        json.dump(setting, ef, indent=2)

    redis_state.client.execute_command("SELECT", 0)
    redis_state.client.hset("System", "setup", json.dumps(setting))
    redis_state.client.hset("System", "mode", mode)
    redis_state.client.hset("Service", "save", 1)
    redis_state.client.hset("Service", "restart", 1)

    if aesState.checkAdmin(adminPass):
        try:
            conn = get_db_connection()
            conn.row_factory = sqlite3.Row
            cursor = conn.cursor()
            hashed_password = bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")
            cursor.execute(
                "INSERT INTO `user` (account,username, password, email, role, api) VALUES (?, ?, ?, ?, ?, ?)",
                (account, name, hashed_password, email, '3', diag)
            )
            conn.commit()
            conn.close()
            return {"passOK": "1"}
        except Exception as e:
            print(str(e))
            return {"passOK": "0", "msg": str(e)}
    else:
        return {"passOK": "0", "msg": 'Admin Password is Wrong'}

@router.get('/checkSession')
def check_session(request: Request):
    user = request.session.get("user")
    userRole = request.session.get("userRole")
    devMode = request.session.get("devMode")
    if user is None:
        return {"loggedIn": False, "username": None, "userRole":None, "mode":None}
    else:
        return {"loggedIn": True, "username": user, "userRole":userRole, "mode":devMode}

# @router.get('/fetchlangset')
# def get_langset(request: Request):
#     lang = 'en'
#     langs = request.session.get("lang", "eng")
#     if langs :
#         if langs == 'kor':
#             lang = 'ko'
#         else:
#             lang = 'en'
#     else:
#         lang = 'en'
#     file_path = LANG_FILES.get(lang, LANG_FILES["en"])  # 기본값 영어
#     readpath = os.path.join(SETTING_FOLDER, file_path)
#     try:
#         with open(readpath, "r", encoding="utf-8") as f:
#             lang_data = json.load(f)
#         return {"success": True, "data": lang_data}
#     except FileNotFoundError:
#         return {"success": False, "error": "Language file not found"}

@router.get('/getUser')
def getUser(request: Request):
    if 'user' not in request.session:
        return {"passOK": "0"}

    user = request.session.get("user")
    # devMode = request.session.get("devMode")
    try:
        conn = get_db_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM user where account=?",(user,))
        rows = cursor.fetchone()
        conn.close()
        # userlist = [dict(row) for row in rows]
        if len(rows) >0:
            result = dict()
            result["username"] = rows["username"]
            result["email"] = rows["email"]
            result["api"] = rows["api"]
            # result["devMode"] = devMode
            return {"passOK": "1", "data": result}
        else:
            return {"passOK": "0"}
    except Exception as e:
        return {"passOK": "0"}

@router.get('/getUserList')
def get_Userlist(request: Request):
    if 'user' not in request.session:
        return {"passOK": "0"}

    # if request.session["API"] != 'Login':
    #     logout()

    user = request.session.get("user")
    try:
        conn = get_db_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM user")
        rows = cursor.fetchall()
        conn.close()
        userlist = [dict(row) for row in rows]
        if len(rows) > 0:
            return {"passOK": "1", "data": userlist}
        else:
            return {"passOK": "0"}
    except Exception as e:
        return {"passOK": "0"}

@router.post('/resetPassword')
def resetPassword(data: Resetpassword):
    #print(f"🔍 요청된 데이터: account={data.username}, email={data.email}")
    password = data.password
    account = data.username  # 프론트엔드에서 username으로 보내지만 실제로는 account
    email = data.email
    result = 0
    
    try:
        conn = get_db_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        
        # account와 email이 모두 일치하는 사용자를 찾음
        cursor.execute(
            "SELECT * FROM user WHERE account=? AND email=?", 
            (account, email)
        )
        rows = cursor.fetchone()
        
        if rows:
            print(f"✅ 찾은 사용자: account={rows['account']}, username={rows['username']}, email={rows['email']}")
            # account와 email이 모두 일치하는 경우에만 비밀번호 변경
            hashed_password = bcrypt.hashpw(
                password.encode("utf-8"), 
                bcrypt.gensalt()
            ).decode("utf-8")
            
            cursor.execute(
                "UPDATE user SET password=? WHERE account=? AND email=?",
                (hashed_password, account, email)
            )
            #print(f"🔄 비밀번호 변경 완료: account={account}")
            result = 1  # 성공
        else:
            #print(f"❌ 일치하는 계정을 찾을 수 없음: account={account}, email={email}")
            result = 0  # 실패
            
        conn.commit()
        conn.close()
        return {"passOK": str(result)}
        
    except Exception as e:
        #print(f"💥 에러 발생: {str(e)}")
        return {"passOK": str(result)}

async def getAPIUsers():
    async with httpx.AsyncClient() as client:
        try:
            response = await client.get(
                f"http://{os_spec.restip}:5000/api/getUsers")
            datas = response.json()
            if response.status_code in [400, 401, 500]:
                data = []
            else:
                data = datas
        except Exception as e:
            data = []
        return data

async def changeLevelAPI(username, role):
    async with httpx.AsyncClient() as client:
        try:
            # response = await client.get(
            #     f"http://{os_spec.restip}:5000/api/changeUserLevel?username={username}&level={role}")
            data = {
                "username":username,
                "level":int(role)
            }
            response = await client.post(
                f"http://{os_spec.restip}:5000/api/changeUserLevel", json=data)
            datas = response.json()
            if response.status_code in [400, 401, 500]:
                flag = 2
            else:
                flag = 1
        except Exception as e:
            flag = 3
    return flag

async def deleteAPI(account):
    async with httpx.AsyncClient() as client:
        try:
            # response = await client.get(
            #     f"http://{os_spec.restip}:5000/api/deleteUser?username={account}")
            data = {
                "username":account
            }
            response = await client.post(
                f"http://{os_spec.restip}:5000/api/deleteUser", json=data)
            datas = response.json()
            if response.status_code in [400, 401, 500]:
                flag = 2
            else:
                flag = 1
        except Exception as e:
            flag = 3
    return flag

async def changePassAPI(account,password, newPass):
    async with httpx.AsyncClient() as client:
        try:
            data ={
                "username":account,
                "oldpassword":password,
                "password":newPass
            }
            response = await client.post(
                f"http://{os_spec.restip}:5000/api/changePassword",json=data)
            datas = response.json()
            if response.status_code in [400, 401, 500]:
                flag = 2
            else:
                flag = 1
        except Exception as e:
            flag = 3
    return flag

async def joinAPI(account, newPass, role):
    async with httpx.AsyncClient() as client:
        try:
            data = {
                "username":account,
                "password":newPass,
                "level":int(role)
            }
            # response = await client.get(
            #     f"http://{os_spec.restip}:5000/api/createUser?username={account}&password={newPass}&level={role}")
            response = await client.post(f"http://{os_spec.restip}:5000/api/createUser", json=data)
            datas = response.json()
            if response.status_code in [400, 401, 500]:
                flag = 2
            else:
                flag = 1
        except Exception as e:
            flag = 3
    return flag

async def loginAPI(account, newPass, role):
    async with httpx.AsyncClient(timeout=httpx.Timeout(2.0)) as client:
        try:
            if role == '3':
                data = {
                    "username" : account,
                    "password" : newPass,
                }
            else:
                data = {
                    "username": account,
                    "password": newPass,
                    "level": role
                }
            # response = await client.get(
            #     f"http://{os_spec.restip}:5000/api/login?username={account}&password={newPass}")
            response = await client.post(f"http://{os_spec.restip}:5000/api/login", json=data)
            datas = response.json()
            if response.status_code in [400, 401, 500]:
                flag = 2
            else:
                flag = 1
        except Exception as e:
            logging.error(f"❌ Smart System Login Error: {e}")
            flag = 3
    return flag

async def logoffAPI():
    async with httpx.AsyncClient(timeout=httpx.Timeout(2.0)) as client:
        try:
            # response = await client.get(f"http://{os_spec.restip}:5000/api/logoff")
            data = {}
            response = await client.post(
                f"http://{os_spec.restip}:5000/api/logff", json=data)
            datas = response.json()
            if response.status_code in [400, 401, 500]:
                flag = 2
            else:
                flag = 1
        except Exception as e:
            flag = 3
            print(str(e))
    return flag

async def checkLoginAPI(request:Request):
    userRole = request.session.get('userRole')

    try:
        conn = get_db_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM user")
        rows = cursor.fetchall()
        op = False
        gst = False
        for row in rows:
            if row["role"] == 0:
                op = True
            if row["role"] == 1:
                gst = True
        if int(userRole) == 3:
            user = 'ntek'
            upass = 'ntek9135!'
        elif int(userRole) == 2:
            user = 'admin'
            upass = 'admin'
        elif int(userRole) == 0:
            if op:
                user = 'operator'
                upass = 'operator'
            else:
                user = None
                upass = None
        else:
            if gst:
                user = 'guest'
                upass = 'guest'
            else:
                user = None
                upass = None
        if user and upass:
            await loginAPI(user,upass,userRole)
            return True
        else:
            return False
    except Exception as e:
        print(str(e))
        return False


@router.get('/getAPIUsers')
async def get_users():
    datas = await getAPIUsers()
    if len(datas) > 0:
        return {"success": True, "data": datas}
    else:
        return {"success": False}

@router.get('/saveUser/{account}/{role}')
def updateRole(account, role):
    flag = 0
    try:
        conn = get_db_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM user WHERE account=?", (account,))
        rows = cursor.fetchone()
        if rows:
            cursor.execute(
                "update user set role=? where account=?",
                (role, account)
            )
            flag = 1
        else:
            flag = 0
        conn.commit()
        conn.close()
        return {"passOK": str(flag)}
    except Exception as e:
        return {"passOK": str(flag)}

@router.get('/removeUser/{username}')
def removeUser(username):
    removed = 0
    if username == 'admin':
        return {"passOK": str(removed)}
    try:
        conn = get_db_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM user WHERE account=?", (username,))
        rows = cursor.fetchone()
        if rows:
            cursor.execute(
                "delete from user where account=?",
                (username,)
            )
        conn.commit()
        conn.close()
        removed = 1
        return {"passOK": str(removed)}
    except Exception as e:
        return {"passOK": str(removed)}


@router.post('/updateProfile')
async def updateProfile(request:Request, data: User):
    flag = 0
    if check_useDiagnosis():  #Use Diagnonis
        apis = 'Yes'
    else:
        apis = 'No'
    try:
        email = data.email
        username = data.username
        password = data.password
        newPass = data.newPass
        api = data.api
        account = request.session.get("user")
        conn = get_db_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM user WHERE account=?", (account,))
        rows = cursor.fetchone()
        if rows:
            if rows["role"] == '3':
                cursor.execute(
                    "update user set username=?,email=?,api=? where account=?",
                    (username, email, api, account)
                )
                flag = 1
            else:
                if bcrypt.checkpw(password.encode("utf-8"), rows["password"].encode("utf-8")):
                    hashed_password = bcrypt.hashpw(newPass.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")
                    if apis == 'Yes':
                        if rows["role"] == '2':
                            if rows["api"] == 'No' and api != rows["api"]:
                                flag = await joinAPI(account,newPass,rows["role"])
                                if flag == 1:
                                    flag = await loginAPI(account, newPass, rows["role"])
                            elif rows["api"] == 'Yes' and api != rows["api"]:
                                flag = await deleteAPI(account)
                            else:
                                flag = await changePassAPI(account,password, newPass)
                    
                    cursor.execute(
                        "update user set username=?,password=?,email=?,api=? where account=?",
                        (username, hashed_password, email, api, account)
                    )
               
                    flag = 1
                else:
                    flag = 0
        else:
            flag = 0
        conn.commit()
        conn.close()
        return {"passOK": str(flag)}
    except Exception as e:
        return {"passOK": str(flag)}


@router.get('/joinAPI/{role}')
async def join_default(role):
    if check_useDiagnosis():
        apis = 'Yes'
    else:
        apis = 'No'
    if role == '0':
        username = 'operator'
        password = 'operator'
    else:
        username = 'guest'
        password = 'guest'
    if apis:
        flag = await joinAPI(username,password,role)
    else:
        flag = 3
    if apis and flag == 1:
        try:
            conn = get_db_connection()
            conn.row_factory = sqlite3.Row
            cursor = conn.cursor()
            hashed_password = bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")
            cursor.execute(
                "INSERT INTO `user` (account, username, password, email, role, api) VALUES (?, ?, ?, ?, ?, ?)",
                (username, username, hashed_password, '', role, 'No' )
            )
            conn.commit()
            conn.close()
            flag = 1
        except Exception as e:
            print(str(e))
            flag = 3
    return {"passOK": str(flag)}

@router.get('/removeAPI/{role}')
async def remove_default(role):
    if check_useDiagnosis():
        apis = 'Yes'
    else:
        apis = 'No'
    if role == '0':
        username = 'operator'
    else:
        username = 'guest'
    if apis:
        flag = await deleteAPI(username)
    else:
        flag = 3
    if apis and flag == 1:
        try:
            conn = get_db_connection()
            conn.row_factory = sqlite3.Row
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM user WHERE account=?", (username,))
            rows = cursor.fetchone()
            if rows :
                cursor.execute(
                    "delete from user where account=?",
                    (username,)
                )
            conn.commit()
            conn.close()
            flag = 1
        except Exception as e:
            flag = 3
    return {"passOK": str(flag)}

@router.post('/join')
def join(data: SignupUser):
    name = data.username
    account = data.account
    password = data.password
    email = data.email
    role = data.role
    # if check_useDiagnosis() and role == '0':
    #     apis = 'Yes'
    # else:
    #     apis = 'No'
    if not account:
        return {"passOK": "0"}
    try:
        conn = get_db_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        hashed_password = bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")
        cursor.execute(
            "INSERT INTO `user` (account, username, password, email, role, api) VALUES (?, ?, ?, ?, ?, ?)",
            (account, name, hashed_password, email, role, 'No' )
        )
        conn.commit()
        conn.close()
        return {"passOK": "1"}
    except Exception as e:
        print(str(e))
        return {"passOK": "0", "msg" : str(e)}
    
def get_mode_from_redis(redis_client) -> str:  
    try:
        redis_client.execute_command("SELECT", 0)
        if redis_client.hexists("System", "setup"):
            setup_json = redis_client.hget("System", "setup")
            setting = json.loads(setup_json)
            return setting.get("mode", "")
    except Exception as e:
        print(f"[ERROR] Redis mode fetch 실패: {e}")
    return ""


# @router.post('/checkLogins')
# async def checkLogins(request: Request, data: Login):
#     # data = request.json()
#     name = data.account
#     password = data.password
#     lang = data.lang
#     flag = 0  # flag = 0 초기, 1 성공,2 API 로그인 실패, 3 API 로그인 계정 없음, 4 비번 오류, 5 계정 오류,
#     if not name:
#         return {"passOK": str(flag)}
#     try:
#         conn = get_db_connection()
#         conn.row_factory = sqlite3.Row
#         cursor = conn.cursor()
#         table_name = "user"
#         cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name=?", (table_name,))
#         exists = cursor.fetchone()
#         role = ''
#         if exists:
#             cursor.execute("SELECT * FROM user where account=?", (name,))
#             rows = cursor.fetchone()
#             if rows:
#                 if bcrypt.checkpw(password.encode("utf-8"), rows["password"].encode("utf-8")):
#                     request.session["user"] = name
#                     request.session["lang"] = lang
#                     request.session["userRole"] = rows["role"]
#                     request.session["api"] = rows["api"]
#                     role = rows["role"]
#
#                     flag = 1
#                     if role == '3':
#                         mode = get_mode_from_redis(redis_state.client)
#                         if mode and mode != 'device0':
#                             if rows["api"] == 'Yes':
#                                 ssflag = is_service_active('smartsystemsservice')
#                                 sraflag = is_service_active('smartsystemsrestapiservice')
#                                 if ssflag and sraflag:
#                                     tmpResult = await loginAPI(name, password, role)
#                                     if tmpResult == 1:
#                                         flag = 1
#                                     else:
#                                         flag = 2
#                                 else:
#                                     flag = 1
#                             else:
#                                 flag = 1
#                         else:
#                             flag = 1
#                     elif check_useDiagnosis():
#                         if role == '2':
#                             if rows["api"] == 'Yes':
#                                 ssflag = is_service_active('smartsystemsservice')
#                                 sraflag = is_service_active('smartsystemsrestapiservice')
#                                 if ssflag and sraflag:
#                                     tmpResult = await loginAPI(name, password, role)
#                                     if tmpResult == 1:
#                                         flag = 1
#                                     else:
#                                         flag = 2
#                                 else:
#                                     flag = 1
#                             else:
#                                 flag = 1
#                         else:
#                             if role == '0':
#                                 apiName = 'operator'
#                             else:
#                                 apiName = 'guest'
#                             cursor.execute("SELECT * FROM user WHERE account=?", (apiName,))
#                             rows = cursor.fetchone()
#                             if rows:
#                                 loginEnable = True
#                             else:
#                                 loginEnable = False
#                             if loginEnable:
#                                 tmpResult = await loginAPI(apiName, apiName, role)
#                                 if tmpResult == 1:
#                                     flag = 1
#                                 else:
#                                     flag = 2
#                             else:
#                                 flag = 3
#                 else:
#                     flag = 4
#             else:
#                 print(f"❌ 로그인 실패")
#                 flag = 5
#         else:
#             flag = 0
#         conn.commit()
#         conn.close()
#         redis_state.client.execute_command("SELECT", 0)
#         if redis_state.client.hexists("System", "mode"):
#             mode = redis_state.client.hget("System", "mode")
#             request.session["devMode"] = mode
#         else:
#             mode = 'device0'
#             request.session["devMode"] = 'device0'
#         if flag == 1:
#             # print(f"🔐 로그인 성공: {name}, 역할: {role}, 언어: {lang}")
#             return {"passOK": str(flag), "data": {"lang": lang, "userRole": role}, "mode": mode}
#         else:
#             # print(f"❌ 로그인 실패: {name}, 역할: {role}, 언어: {lang}")
#             return {"passOK": str(flag)}
#     except Exception as e:
#         print(str(e))
#         return {"passOK": "0"}

@router.post('/checkLogins')
async def checkLogins(request:Request, data: Login):
    # data = request.json()
    name = data.account
    password = data.password
    lang = data.lang
    flag = 0   #flag = 0 초기, 1 성공, 4 비번 오류, 5 계정 오류
    if not name:
        return {"passOK": str(flag)}
    try:
        conn = get_db_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        table_name = "user"
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name=?", (table_name,))
        exists = cursor.fetchone()
        role = ''
        if exists:
            cursor.execute("SELECT * FROM user where account=?", (name,))
            rows = cursor.fetchone()
            if rows :
                if bcrypt.checkpw(password.encode("utf-8"), rows["password"].encode("utf-8")):
                    request.session["user"] = name
                    request.session["lang"] = lang
                    request.session["userRole"] = rows["role"]
                    role = rows["role"]
                    flag = 1
                else:
                    flag = 4
            else:
                print(f"❌ 로그인 실패")
                flag = 5
        else:
            flag = 0
        conn.commit()
        conn.close()
        mode_setup = ''
        redis_state.client.execute_command("SELECT", 0)
        if redis_state.client.hexists("System","setup"):
            setups = redis_state.client.hget("System","setup")
            setup = json.loads(setups)
            mode_setup = setup["mode"]
        if redis_state.client.hexists("System","mode"):
            mode = redis_state.client.hget("System","mode")
            if mode != mode_setup:
                mode = mode_setup
            request.session["devMode"] = mode
            redis_state.client.hset("System", "mode", mode)
        else:
            if mode_setup != '':
                mode = mode_setup
                request.session["devMode"] = mode_setup
        if flag:
            # print(f"🔐 로그인 성공: {name}, 역할: {role}, 언어: {lang}")
            return {"passOK": str(flag), "data":{"lang": lang, "userRole" : role}, "mode":mode}
        else:
            # print(f"❌ 로그인 실패: {name}, 역할: {role}, 언어: {lang}")
            return {"passOK": str(flag)}
    except Exception as e:
        print(str(e))
        return {"passOK": "0"}

@router.get('/checkRemote/{user}')
async def check_remoteUser(user, request:Request):
    if not user:
        return {"success":False}
    try:
        redis_state.client.execute_command("SELECT", 0)
        if redis_state.client.hexists("System","mode"):
            mode = redis_state.client.hget("System","mode")
            request.session["devMode"] = mode
        else:
            mode = 'device0'
            request.session["devMode"] = 'device0'
        if user == 'ntek':
            data = {
                "account" : user,
                "userRole" : 3,
                "mode": mode
            }
            request.session["user"] = user
            request.session["userRole"] = 3
            return {"success": True, "data": data}
        else:
            conn = get_db_connection()
            conn.row_factory = sqlite3.Row
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM user where account=?", (user,))
            rows = cursor.fetchone()
            if rows:
                flag = 1
                data = {
                    "account": user,
                    "userRole": 2,
                    "mode": mode
                }
                request.session["user"] = user
                request.session["userRole"] = 2
            else:
                flag = 0
            conn.commit()
            conn.close()
        if flag == 1:
            return {"success":True, "mode":mode, "data":data }
        else:
            return {"success":False}
    except Exception as e:
        print(str(e))

@router.get('/logout')
async def logout(request:Request):
    request.session.clear()  # 모든 세션 삭제
    if check_useDiagnosis():
        flag = await logoffAPI()
    return {"success": True}
