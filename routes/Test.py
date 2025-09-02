from fastapi import APIRouter, Request
import json, os
from pathlib import Path
base_dir = Path(__file__).resolve().parent
SETTING_FOLDER = base_dir.parent.parent / "config"  # ⬅️ 두 단계 상위로
router = APIRouter()

def findNode(supers, name, path): #check node level
    result = ''
    for i in range(0, len(supers)):
        if supers[i]["Name"] == name:
            if supers[i]["Path"] == path:
               break
            else:
                parts = [part.strip() for part in path.split("|")]
                result = "_".join(parts[1:])
                break
    return result

@router.get("/getDiagnosis_trans/{asset}") #Trans Temporary
async def get_diagnosis_trans(asset, request:Request):
    try:
        with open("./trans_st.json", "r", encoding="utf-8") as f:
            datas = json.load(f)
            # print("Raw JSON:", datas)  # JSON이 잘 읽히는지 확인
    except Exception as e:
        return {"passOK": "0", "error": f"An unexpected error occurred: {str(e)}"}

    if len(datas) > 0:
        return {"success": True, "data": datas["BarGraph"]}
    else:
        return {"success": False, "error": "No Data"}

@router.get("/getTrendParameter_trans/{asset}") #Temporary Router
async def get_trendParams(asset, request:Request):
    with open("./trans_pqTrend.json", "r", encoding="utf-8") as f:
        try:
            datas = json.load(f)
        except Exception as e:
            print(str(e))
            return {"success":False, "error": str(e)}
    if len(datas) > 0:
        superlist = []
        assetId = -1
        for a in range(0, len(datas)):
            if datas[a]["Name"] == asset:
                assetId = datas[a]["ID"]
                break
        for a in range(0, len(datas)):
            if datas[a]["ParentID"] == assetId:
                superdict = dict()
                superdict["ID"] = datas[a]["ID"]
                superdict["Name"] = datas[a]["Name"]
                superdict["Title"] = datas[a]["Title"]
                superdict["Path"] = datas[a]["Path"]
                superdict["isParent"] = True
                Titles = findNode(superlist, datas[a]["Name"], datas[a]["Path"])
                if Titles != "":
                    superdict["Title"] = superdict["Title"] +"_"+Titles
                    superlist.append(superdict)
                else:
                    superlist.append(superdict)
            else:
                continue
        for j in range(0, len(datas)):
            isSuper = True
            for i in range(0, len(superlist)):
                if superlist[i]["ID"] == datas[j]["ParentID"]:
                    subdict = dict()
                    subdict["ID"] = datas[j]["ID"]
                    subdict["Name"] = datas[j]["Name"]
                    subdict["Title"] = datas[j]["Title"]
                    if "children" in superlist[i]:
                        superlist[i]["children"].append(subdict)
                    else:
                        superlist[i]["children"] = []
                        superlist[i]["children"].append(subdict)
                    isSuper = False
                else:
                    continue
            if isSuper and datas[j]["NodeType"] == 11:
                superdict = dict()
                superdict["ID"] = datas[j]["ID"]
                superdict["Name"] = datas[j]["Name"]
                superdict["Title"] = datas[j]["Title"]
                superdict["Path"] = datas[j]["Path"]
                superdict["isParent"] = True
                Titles = findNode(superlist, datas[j]["Name"], datas[j]["Path"])
                if Titles != "":
                    superdict["Title"] = superdict["Title"] + "_" + Titles
                    superlist.append(superdict)
                else:
                    superlist.append(superdict)
    if len(superlist) > 0:
        return {"success":True, "superlist" : superlist}
    else:
        return {"success":False, "error": "No Data"}

@router.get("/getDiagTransPQ")
def get_diag_transpq():
    try:
        with open("./trans_pq.json", "r", encoding="utf-8") as f:
            datas = json.load(f)
            # print("Raw JSON:", datas)  # JSON이 잘 읽히는지 확인
    except Exception as e:
        return {"passOK": "0", "error": f"An unexpected error occurred: {str(e)}"}
    superNodes = []
    if len(datas) > 0:
        statuslist = datas["BarGraph"]
        for a in range(0, len(datas["BarGraph"])):
            superNodes.append(datas["BarGraph"][a]["NodeType"])
        superlist = []
        for a in range(0, len(datas["TreeList"])):
            if datas["TreeList"][a]["NodeType"] in superNodes:
                superdict = dict()
                superdict["ID"] = datas["TreeList"][a]["ID"]
                superdict["Name"] = datas["TreeList"][a]["Name"]
                superdict["Title"] = datas["TreeList"][a]["Title"]
                superdict["Path"] = datas["TreeList"][a]["Path"]
                superdict["Status"] = datas["TreeList"][a]["Status"]
                superdict["Value"] = datas["TreeList"][a]["Value"]
                superdict["isParent"] = True
                Titles = findNode(superlist, datas["TreeList"][a]["Name"], datas["TreeList"][a]["Path"])
                if Titles != "":
                    superdict["Title"] = superdict["Title"] + "_" + Titles
                    superlist.append(superdict)
                else:
                    superlist.append(superdict)
            else:
                continue
        for i in range(0, len(superlist)):
            for j in range(0, len(datas["TreeList"])):
                if superlist[i]["ID"] == datas["TreeList"][j]["ParentID"]:
                    subdict = dict()
                    subdict["ID"] = datas["TreeList"][j]["ID"]
                    subdict["Name"] = datas["TreeList"][j]["Name"]
                    subdict["Title"] = datas["TreeList"][j]["Name"]
                    subdict["Status"] = datas["TreeList"][j]["Status"]
                    subdict["Value"] = datas["TreeList"][j]["Value"]
                    subdict["NodeType"] = datas["TreeList"][j]["NodeType"]
                    if datas["TreeList"][j]["NodeType"] == 10:
                        subdict["isSub"] = True
                    else:
                        subdict["isSub"] = False
                    if "children" in superlist[i]:
                        superlist[i]["children"].append(subdict)
                    else:
                        superlist[i]["children"] = []
                        superlist[i]["children"].append(subdict)
                else:
                    continue
        for i in range(0, len(superlist)):
            if "children" in superlist[i]:
                for k in range(0, len(superlist[i]["children"])):
                    if superlist[i]["children"][k]["isSub"]:
                        for j in range(0, len(datas["TreeList"])):
                            if datas["TreeList"][j]["NodeType"] == 11 and superlist[i]["children"][k]["ID"] == \
                                    datas["TreeList"][j]["ParentID"]:
                                subdict = dict()
                                subdict["ID"] = datas["TreeList"][j]["ID"]
                                subdict["Name"] = datas["TreeList"][j]["Name"]
                                subdict["Title"] = datas["TreeList"][j]["Title"]
                                subdict["Status"] = datas["TreeList"][j]["Status"]
                                subdict["Value"] = datas["TreeList"][j]["Value"]
                                if "children" in superlist[i]["children"][k]:
                                    superlist[i]["children"][k]["children"].append(subdict)
                                else:
                                    superlist[i]["children"][k]["children"] = []
                                    superlist[i]["children"][k]["children"].append(subdict)
                            else:
                                continue
                    else:
                        continue
            else:
                continue
    # print(superlist)
    if len(superlist) > 0:
        return {"success": True, "data_status": statuslist, "data_tree": superlist}
    else:
        return {"success": False, "error": "No Data"}

@router.get("/getDiagTrans")  #Temporary JSON : Transformer
def get_diag_trans():
    superNodes = []
    try:
        with open("./trans_st.json", "r", encoding="utf-8") as f:
            datas = json.load(f)
            # print("Raw JSON:", datas)  # JSON이 잘 읽히는지 확인
    except Exception as e:
        return {"passOK": "0", "error": f"An unexpected error occurred: {str(e)}"}
    if len(datas) > 0:
        statuslist = datas["BarGraph"]
        for a in range(0, len(datas["BarGraph"])):
            superNodes.append(datas["BarGraph"][a]["NodeType"])
        superlist = []
        for a in range(0, len(datas["TreeList"])):
            if datas["TreeList"][a]["NodeType"] in superNodes:
                superdict = dict()
                superdict["ID"] = datas["TreeList"][a]["ID"]
                superdict["Name"] = datas["TreeList"][a]["Name"]
                superdict["Title"] = datas["TreeList"][a]["Title"]
                superdict["Path"] = datas["TreeList"][a]["Path"]
                superdict["Status"] = datas["TreeList"][a]["Status"]
                superdict["Value"] = datas["TreeList"][a]["Value"]
                superdict["isParent"] = True
                Titles = findNode(superlist, datas["TreeList"][a]["Name"], datas["TreeList"][a]["Path"])
                if Titles != "":
                    superdict["Title"] = superdict["Title"] +"_"+Titles
                    superlist.append(superdict)
                else:
                    superlist.append(superdict)
            else:
                continue
        for i in range(0, len(superlist)):
            for j in range(0, len(datas["TreeList"])):
                if superlist[i]["ID"] == datas["TreeList"][j]["ParentID"]:
                    subdict = dict()
                    subdict["ID"] = datas["TreeList"][j]["ID"]
                    subdict["Name"] = datas["TreeList"][j]["Name"]
                    subdict["Title"] = datas["TreeList"][j]["Title"]
                    subdict["Status"] = datas["TreeList"][j]["Status"]
                    subdict["Value"] = datas["TreeList"][j]["Value"]
                    if "children" in superlist[i]:
                        superlist[i]["children"].append(subdict)
                    else:
                        superlist[i]["children"] = []
                        superlist[i]["children"].append(subdict)
                else:
                    continue
    if len(superlist) > 0:
        return {"success":True, "data_status" : statuslist, "data_tree" : superlist}
    else:
        return {"success":False, "error": "No Data"}