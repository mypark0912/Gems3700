from fastapi import APIRouter, UploadFile, File
from states.global_state import redis_state
import os, csv, sqlite3
import ujson as json
from pathlib import Path
from pydantic import BaseModel
from datetime import date

router = APIRouter()

# Path 객체 절대경로

base_dir = Path(__file__).resolve().parent
SETTING_FOLDER = base_dir.parent.parent / "config"  # ⬅️ 두 단계 상위로
DB_PATH = os.path.join(SETTING_FOLDER, "maintenance.db")

def get_db_connection():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row  # → dict처럼 row 접근 가능
    return conn

class Post(BaseModel):
    title: str
    context: str
    mtype: int
    utype: str
    f_version: str
    a_version: str
    w_version: str

class CaliSet(BaseModel):
    cmd: str
    ref : str

def parse_csv(file_path):
    data = []
    with open(file_path, mode='r', encoding='utf-8') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            data.append(row)
    return data

def parse_csv_filelike(file_obj):
    data = []
    file_obj.seek(0)  # 위치 초기화
    reader = csv.DictReader(line.decode('utf-8') for line in file_obj)
    for row in reader:
        data.append(row)
    return data

@router.post('/upload')
def upload_file(file: UploadFile = File(...)):
    if file.filename == '':
        print('No selected file')
        return {'passOK': '0', 'error': 'No selected file'}
    else:
        try:
            calibration_file_path = os.path.join(SETTING_FOLDER, file.filename)
            with open(calibration_file_path, "wb") as f:
                content = file.file.read()
                f.write(content)
            data = parse_csv_filelike(file.file)
            redis_state.client.hset('calibration','setup', json.dumps(data))
            return {'passOK': '1', 'data': data, 'file_path': file.filename}
            # original_file_path = os.path.join(SETTING_FOLDER, 'setup.json')
            # if os.path.exists(original_file_path):
            #     with open(original_file_path, "r", encoding="utf-8") as f:
            #         setting = json.load(f)
            #     if len(data) > 0:
            #         for idx, ch in enumerate(setting["channel"]):
            #             ctData = setting["channel"][idx]["ctInfo"]
            #             for key, value in ctData:
            #                 ctData[key] = data[key]
            #             ptData = setting["channel"][idx]["ptInfo"]
            #             for key, value in ptData:
            #                 ptData[key] = data[key]
            #     if setting:
            #         with open(original_file_path, "w", encoding="utf-8") as f:
            #             json.dump(setting, f, indent=4)
            #         return {'passOK': '1', 'data': data, 'file_path':file.filename}
            #     else:
            #         return {'passOK': '0', 'error': 'Apply Failed'}
            # else:
            #     return {'passOK': '0', 'error': 'No exist setup file'}
        except Exception as e:
            return {'passOK': '0', 'error': str(e)}

@router.get('/checkSetup')
def check_setup():
    if redis_state.client.hexists('calibration','setup'):
        data = redis_state.client.hget('calibration','setup')
        datalist = json.loads(data)
        return {'passOK': '1', 'data':datalist}
    else:
        return {'passOK': '0', 'error': 'No exist calibration setup'}

@router.get('/calibrate/start')
def cali_start():
    redis_state.client.execute_command("SELECT", 1)
    if redis_state.client.hexists('calibration','setup'):
        data = redis_state.client.hget('calibration','setup')
        datalist = json.loads(data)
        print(datalist)
        redis_state.client.hset('Service', 'apply', 1)
        return {'passOK': '1'}
    else:
        return {'passOK': '0', 'error': 'No exist calibration setup'}

@router.get('/calibrate/end')
def cali_end():
    file_path = os.path.join(SETTING_FOLDER, 'setup.json')
    try:
        with open(file_path, "r", encoding="utf-8") as f:
            setting = json.load(f)
            redis_state.client.execute_command("SELECT", 1)
            redis_state.client.hset("System", "setup", json.dumps(setting))
            return {'passOK': '1'}
    except Exception as e:
        return {'passOK': '0', 'error': 'No exist setup file'}

@router.post('/calibrate/cmd/{types}')
def set_cmd(data:CaliSet, types):
    try:
        if types == 'SET':
            if data.ref != 'None':
                msg =  f"{data.cmd} {data.ref}"
            else:
                msg = data.cmd
        else:
            msg = data.cmd
        redis_state.client.execute_command("SELECT", 1)
        redis_state.client.hset('calibration','command',msg)
        redis_state.client.hset('calibration', 'cflag', 1)
        return {'passOK': '1'}
    except Exception as e:
        print(str(e))
        return {'passOK': '0'}

@router.post('/savePost/{mode}/{idx}')
def save_post(data: Post, mode: int, idx:int):
    title = data.title
    context = data.context
    mtype = data.mtype
    utype = data.utype
    f_version = data.f_version
    a_version = data.a_version
    w_version = data.w_version
    today = date.today()
    formatted = today.strftime("%Y-%m-%d")
    try:
        conn = get_db_connection()
        cursor = conn.cursor()
        if mode == 0:
            cursor.execute(
                "INSERT INTO `maintenance` (title,context, mtype, utype, f_version, a_version, w_version, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
                (title, context, mtype, utype, f_version,a_version,w_version, formatted)
            )
        else:
            cursor.execute(
                "Update `maintenance`set title=?,context=?, mtype=?, utype=?, f_version=?, a_version=?, w_version=?, date=? where id=?",
                (title, context, mtype, utype, f_version, a_version, w_version, formatted, idx )
            )
        conn.commit()
        conn.close()
        return {"passOK": "1"}
    except Exception as e:
        print(str(e))
        return {"passOK": "0", "msg": str(e)}

@router.get('/getPost')
def get_post():
    result = []
    try:
        db_exists = os.path.exists(DB_PATH)
        conn = get_db_connection()
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        flag = False
        if not db_exists:
            cursor.execute('''
                        CREATE TABLE IF NOT EXISTS maintenance (
                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                            title TEXT NOT NULL,
                            context TEXT NOT NULL,
                            mtype int NOT NULL,
                            utype TEXT NOT NULL,
                            f_version TEXT,
                            a_version TEXT,
                            w_version TEXT,
                            date TEXT
                        )
                    ''')
            cursor.execute("SELECT * FROM maintenance")
            rows = cursor.fetchall()
            if len(rows) > 0:
                result = [dict(row) for row in rows]
                flag = True
            else:
                flag = False
        else:
            cursor.execute("SELECT * FROM maintenance")
            rows = cursor.fetchall()
            if len(rows) > 0:
                result = [dict(row) for row in rows]
                flag = True
            else:
                flag = False
        conn.commit()
        conn.close()
    except Exception as e:
        print(str(e))
        flag = False
    if flag:
        return {"result": 1, "data": result}
    else:
        return {"result": 0}

@router.get('/deletePost/{idx}')
def delete_post(idx):
    try:
        db_exists = os.path.exists(DB_PATH)
        if db_exists:
            conn = get_db_connection()
            conn.row_factory = sqlite3.Row
            cursor = conn.cursor()
            cursor.execute("delete FROM maintenance where id=?",(idx,))
            conn.commit()
            conn.close()
            return {"result": 1}
        else:
            return {"result": 0}
    except Exception as e:
        print(str(e))
        return {"result": 0}
