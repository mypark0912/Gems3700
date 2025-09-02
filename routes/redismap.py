class RedisMapDetail:
    freq_keys = [
        {"id": 0, "label": "-", "key": "Frequency"}
    ]
    p_voltage_keys = [
        {"id": 0, "label": "L1", "key": "U1"},
        {"id": 1, "label": "L2", "key": "U2"},
        {"id": 2, "label": "L3", "key": "U3"},
        {"id": 3, "label": "Average", "key": "U~"},
    ]
    p_volt_angle_keys = [
        {"id": 0, "label": "L1", "key": "Angle U1"},
        {"id": 1, "label": "L2", "key": "Angle U2"},
        {"id": 2, "label": "L3", "key": "Angle U3"},
    ]
    l_voltage_keys = [
        {"id": 0, "label": "L1-L2", "key": "U12"},
        {"id": 1, "label": "L2-L3", "key": "U23"},
        {"id": 2, "label": "L3-L1", "key": "U31"},
    ]
    current_keys = [
        {"id": 0, "label": "L1", "key": "I1"},
        {"id": 1, "label": "L2", "key": "I2"},
        {"id": 2, "label": "L3", "key": "I3"},
        {"id": 3, "label": "Average", "key": "I~"},
        {"id": 4, "label": "Total", "key": "I total (I1+I2+I3)"},
    ]
    current_angle_keys = [
        {"id": 0, "label": "L1", "key": "Angle I1"},
        {"id": 1, "label": "L2", "key": "Angle I2"},
        {"id": 2, "label": "L3", "key": "Angle I3"},
    ]
    pf_keys = [
        {"id": 0, "label": "L1", "key": "PF1"},
        {"id": 1, "label": "L2", "key": "PF2"},
        {"id": 2, "label": "L3", "key": "PF3"},
        {"id": 3, "label": "Total", "key": "PF total"},
    ]
    a_powers_keys = [
        {"id": 0, "label": "L1", "key": "P1"},
        {"id": 1, "label": "L2", "key": "P2"},
        {"id": 2, "label": "L3", "key": "P3"},
        {"id": 3, "label": "Total", "key": "P total"},
    ]
    r_powers_keys = [
        {"id": 0, "label": "L1", "key": "Q1"},
        {"id": 1, "label": "L2", "key": "Q2"},
        {"id": 2, "label": "L3", "key": "Q3"},
        {"id": 3, "label": "Total", "key": "Q total"},
    ]
    ap_powers_keys = [
        {"id": 0, "label": "L1", "key": "S1"},
        {"id": 1, "label": "L2", "key": "S2"},
        {"id": 2, "label": "L3", "key": "S3"},
        {"id": 3, "label": "Total", "key": "S total"},
    ]

    thdu_keys = [
        {"id": 0, "label": "L1", "key": "THD V1"},
        {"id": 1, "label": "L2", "key": "THD V2"},
        {"id": 2, "label": "L3", "key": "THD V3"},
    ]

    thdi_keys = [
        {"id": 0, "label": "L1", "key": "THD I1"},
        {"id": 1, "label": "L2", "key": "THD I2"},
        {"id": 2, "label": "L3", "key": "THD I3"},
    ]

    tddi_keys = [
        {"id": 0, "label": "L1", "key": "TDD I1"},
        {"id": 1, "label": "L2", "key": "TDD I2"},
        {"id": 2, "label": "L3", "key": "TDD I3"},
    ]
    kwh_keys = [
        {"id": 0, "label": "Import", "key": "import kwh"},
        {"id": 1, "label": "Export", "key": "export kwh"},
    ]

    unbal_keys = [
        {"id": 0, "label": "Voltage", "key": "Uunbalance2"},
        {"id": 1, "label": "Current", "key": "Iunbalance2"},
    ]

    @classmethod
    def get_Datadict(cls,meters, maxmin, keys, unit):
        p_voltage_data = []
        for item in keys:
            key = item["key"]
            p_voltage_data.append({
                "id": item["id"],
                "subTitle": item["label"],
                "value": round(meters.get(key, 0), 2),
                "max": round(maxmin.get(f"{key}_max", 0), 2),
                # "maxTime": maxmin.get(f"{key}_maxTime", ""),
                "min": round(maxmin.get(f"{key}_min", 0), 2),
                # "minTime": maxmin.get(f"{key}_minTime", ""),
                "unit": unit
            })
        return p_voltage_data

class RedisMapped:
    mapdict = {
        'meter': ['meter', 'thd', 'power', 'energy', 'maxmin'],
        'pq': ['harmonics', 'i_harmonics', 'waveform'],
        'alarm': ['alarm_status', 'alarm_log', 'pq_event', 'event_log', 'event_list'],
        'demand': ['demand', 'demand_log']
    }

    dashboard_meter ={
        "meter":{"U1","U2","U3","U4","I1","I2","I3","Itot","Ubal1","Ibal1","Freq","Ig","Temp", "Ubal_nema","Ibal_nema"},
        "power":{"P1","P2","P3","P4","PF4","PF1","PF2","PF3","S4", "Q4"},
        "energy":{"total_kwh_import","thismonth_kwh_import"},
        "thd":{"THD_U1","THD_U2","THD_U3","THD_I1","THD_I2","THD_I3","TDD_I1","TDD_I2","TDD_I3"}
    }

    dashboard_trans = {
        "meter": {"Freq", "Ig", "Temp"},
    }

    energy_report = [
        "import kwh","export kwh","import kwh this month","export kwh this month","import kwh last month","export kwh last month"
    ]

class RedisMapDetail2:
    freq_keys = [
        {"id": 0, "label": "-", "key": "Freq"}
    ]

    p_voltage_keys = [
        {"id": 0, "label": "L1", "key": "U1"},
        {"id": 1, "label": "L2", "key": "U2"},
        {"id": 2, "label": "L3", "key": "U3"},
    ]
    p_volt_angle_keys = [
        {"id": 0, "label": "L1", "key": "Uangle1"},
        {"id": 1, "label": "L2", "key": "Uangle2"},
        {"id": 2, "label": "L3", "key": "Uangle3"},
    ]
    l_voltage_keys = [
        {"id": 0, "label": "L1-L2", "key": "Upp1"},
        {"id": 1, "label": "L2-L3", "key": "Upp2"},
        {"id": 2, "label": "L3-L1", "key": "Upp3"},
    ]
    current_keys = [
        {"id": 0, "label": "L1", "key": "I1"},
        {"id": 1, "label": "L2", "key": "I2"},
        {"id": 2, "label": "L3", "key": "I3"},
        {"id": 3, "label": "Ground", "key": "Ig"},
    ]
    current_angle_keys = [
        {"id": 0, "label": "L1", "key": "Iangle1"},
        {"id": 1, "label": "L2", "key": "Iangle2"},
        {"id": 2, "label": "L3", "key": "Iangle3"},
    ]
    pf_keys = [
        {"id": 0, "label": "Total", "key": "PF4"},
    ]
    a_powers_keys = [
        {"id": 0, "label": "Total", "key": "P4"},
    ]
    r_powers_keys = [
        {"id": 0, "label": "Total", "key": "Q4"},
    ]
    ap_powers_keys = [
        {"id": 0, "label": "Total", "key": "S4"},
    ]

    thdu_keys = [
        {"id": 0, "label": "L1", "key": "THD_U1"},
        {"id": 1, "label": "L2", "key": "THD_U2"},
        {"id": 2, "label": "L3", "key": "THD_U3"},
    ]

    thdi_keys = [
        {"id": 0, "label": "L1", "key": "THD_I1"},
        {"id": 1, "label": "L2", "key": "THD_I2"},
        {"id": 2, "label": "L3", "key": "THD_I3"},
    ]

    tddi_keys = [
        {"id": 0, "label": "L1", "key": "TDD_I1"},
        {"id": 1, "label": "L2", "key": "TDD_I2"},
        {"id": 2, "label": "L3", "key": "TDD_I3"},
    ]
    kwh_keys = [
        {"id": 0, "label": "Import", "key": "total_kwh_import"},
        {"id": 1, "label": "Export", "key": "total_kwh_export"},
        {"id": 2, "label":"Import", "key":"thismonth_kwh_import"},
        {"id": 3, "label": "Export", "key": "thismonth_kwh_export"},
    ]

    unbal_keys = [
        {"id": 0, "label": "Voltage", "key": "Ubal1"},
        {"id": 1, "label": "Current", "key": "Ibal1"},
    ]

    unbal_keys2 = [
        {"id": 0, "label": "Voltage", "key": "Ubal_nema"},
        {"id": 1, "label": "Current", "key": "Ibal_nema"}
    ]

    @classmethod
    def get_Datadict(cls,meters, maxmin, keys, unit):
        p_voltage_data = []
        for item in keys:
            key = item["key"]
            p_voltage_data.append({
                "id": item["id"],
                "subTitle": item["label"],
                "value": round(meters.get(key, 0), 2),
                "max": round(maxmin.get(f"{key}_max", 0), 2),
                "maxTime": maxmin.get(f"{key}_maxTime", ""),
                "min": round(maxmin.get(f"{key}_min", 0), 2),
                "minTime": maxmin.get(f"{key}_minTime", ""),
                "unit": unit
            })
        return p_voltage_data

