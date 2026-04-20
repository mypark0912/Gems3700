// utils/validation.js — idpm300 전용
export class SettingValidator {
  constructor() {
    this.errors = [];
    this.warnings = [];
  }

  // 에러와 경고 초기화
  reset() {
    this.errors = [];
    this.warnings = [];
  }

  // IP 주소 형식 검증
  validateIPAddress(ip, fieldName) {
    // SNTP Host인 경우 도메인 형식도 허용
    if (fieldName === 'SNTP Host') {
      if (ip === '') {
        return true;
      }
      const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
      const domainRegex = /^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$/;
      if (!ipRegex.test(ip) && !domainRegex.test(ip)) {
        this.errors.push(`${fieldName}: Invalid format. Must be IP address or domain name (${ip})`);
        return false;
      }
      return true;
    }
    const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
    if (!ipRegex.test(ip)) {
      this.errors.push(`${fieldName}: Invalid IP address format (${ip})`);
      return false;
    }
    return true;
  }

  // 숫자 범위 검증
  validateNumberRange(value, min, max, fieldName) {
    const num = parseFloat(value);
    if (isNaN(num)) {
      this.errors.push(`${fieldName}: Must be a valid number (${value})`);
      return false;
    }
    if (num < min || num > max) {
      this.errors.push(`${fieldName}: Must be between ${min} and ${max} (current: ${num})`);
      return false;
    }
    return true;
  }

  // 정수 검증
  validateInteger(value, fieldName) {
    const num = parseInt(value);
    if (isNaN(num) || num.toString() !== value.toString()) {
      this.errors.push(`${fieldName}: Must be a valid integer (${value})`);
      return false;
    }
    return true;
  }

  // 필수 필드 검증
  validateRequired(value, fieldName) {
    if (value === null || value === undefined || value === '') {
      this.errors.push(`${fieldName}: This field is required`);
      return false;
    }
    return true;
  }

  // 포트 번호 검증
  validatePort(port, fieldName) {
    return this.validateNumberRange(port, 1, 65535, fieldName);
  }

  // CT/PT 비율 검증
  validateRatio(primary, secondary, fieldName) {
    const primaryNum = parseFloat(primary);
    const secondaryNum = parseFloat(secondary);

    if (isNaN(primaryNum) || isNaN(secondaryNum)) {
      this.errors.push(`${fieldName}: Primary and Secondary must be valid numbers`);
      return false;
    }

    if (primaryNum <= 0 || secondaryNum <= 0) {
      this.errors.push(`${fieldName}: Primary and Secondary must be greater than 0`);
      return false;
    }

    const ratio = primaryNum / secondaryNum;
    if (ratio < 1) {
      this.warnings.push(`${fieldName}: Primary/Secondary ratio is less than 1 (${ratio.toFixed(2)})`);
    }

    return true;
  }

  // Device Information
  validateDeviceInfo(deviceInfo) {
    if (!deviceInfo) return false;
    this.validateRequired(deviceInfo.name, 'Device Name');
    return this.errors.length === 0;
  }

  // Communication
  validateTcpipSettings(tcpip) {
    if (!tcpip) return false;
    this.validateIPAddress(tcpip.ip_address, 'IP Address');
    this.validateIPAddress(tcpip.subnet_mask, 'Subnet Mask');
    this.validateIPAddress(tcpip.gateway, 'Gateway');
    this.validateIPAddress(tcpip.dnsserver, 'DNS Server');
    return this.errors.length === 0;
  }

  // Modbus
  validateModbusSettings(modbus) {
    if (!modbus) return false;
    this.validatePort(modbus.tcp_port, 'Modbus TCP Port');

    const validBaudRates = [0, 1, 2, 3, 4];
    if (!validBaudRates.includes(modbus.baud_rate)) {
      this.errors.push(`Baud Rate: Invalid baud rate (${modbus.baud_rate})`);
    }
    const validParity = [0, 1, 2];
    if (!validParity.includes(modbus.parity)) {
      this.errors.push(`Parity: Invalid parity (${modbus.parity})`);
    }
    const validDataBits = [7, 8];
    if (!validDataBits.includes(modbus.data_bits)) {
      this.errors.push(`Data Bits: Invalid data bits (${modbus.data_bits})`);
    }
    const validStopBits = [1, 2];
    if (!validStopBits.includes(modbus.stop_bits)) {
      this.errors.push(`Stop Bits: Invalid stop bits (${modbus.stop_bits})`);
    }
    return this.errors.length === 0;
  }

  // CT 설정 검증
  validateCTSettings(ctInfo) {
    if (!ctInfo) return false;
    this.validateNumberRange(ctInfo.startingcurrent, 0.1, 1000, 'Starting Current');
    this.validateNumberRange(ctInfo.inorminal, 1, 10000, 'Rated Current');
    this.validateRatio(ctInfo.ct1, this.getCTSecondaryValue(ctInfo.ct2), 'CT Ratio');
    this.validateNumberRange(ctInfo.zctscale, 0.1, 1000, 'ZCT Scale');

    if (!Array.isArray(ctInfo.direction) || ctInfo.direction.length !== 3) {
      this.errors.push('CT Direction: Must have 3 direction settings');
    } else {
      ctInfo.direction.forEach((dir, index) => {
        if (![0, 1].includes(Number(dir))) {
          this.errors.push(`CT${index + 1} Direction: Must be 0 (Positive) or 1 (Negative)`);
        }
      });
    }

    return this.errors.length === 0;
  }

  // PT 설정 검증
  validatePTSettings(ptInfo) {
    if (!ptInfo) return false;
    this.validateNumberRange(ptInfo.vnorminal, 1, 100000, 'Rated Voltage');
    this.validateRatio(ptInfo.pt1, ptInfo.pt2, 'PT Ratio');

    const validWiringModes = [0, 1, 2, 3, 4];
    if (!validWiringModes.includes(Number(ptInfo.wiringmode))) {
      this.errors.push(`Wiring Mode: Invalid wiring mode index (${ptInfo.wiringmode})`);
    }

    const validLineFreq = [0, 1];
    if (!validLineFreq.includes(Number(ptInfo.linefrequency))) {
      this.errors.push(`Line Frequency: Invalid line frequency index (${ptInfo.linefrequency})`);
    }

    return this.errors.length === 0;
  }

  // Demand 설정 검증
  validateDemandSettings(demand) {
    if (!demand) return false;
    this.validateNumberRange(demand.target, 0, 999999999, 'Demand Target');

    const validIntervalIndices = [0, 1, 2, 3, 4];
    if (!validIntervalIndices.includes(Number(demand.demand_interval))) {
      this.errors.push(`Demand Interval: Invalid interval index (${demand.demand_interval})`);
    }

    return this.errors.length === 0;
  }

  // SNTP 설정 검증
  validateSNTPSettings(inputDict) {
    if (inputDict.sntpInfo) {
      if (inputDict.sntpInfo.host) {
        this.validateIPAddress(inputDict.sntpInfo.host, 'SNTP Host');
      }
      if (inputDict.sntpInfo.timezone !== undefined) {
        this.validateRequired(inputDict.sntpInfo.timezone, 'Timezone');
      }
    } else {
      this.warnings.push('SNTP is enabled but SNTP settings are not configured');
    }

    return this.errors.length === 0;
  }

  // MQTT 설정 검증
  validateMQTTSettings(mqttInfo) {
    if (!mqttInfo) return true;

    if (mqttInfo.host) {
      this.validateIPAddress(mqttInfo.host, 'MQTT Host');
    }
    if (mqttInfo.port) {
      this.validatePort(mqttInfo.port, 'MQTT Port');
    }
    if (mqttInfo.externalport) {
      this.validatePort(mqttInfo.externalport, 'MQTT External Port');
    }

    return this.errors.length === 0;
  }

  // Trend 설정 검증
  validateTrendSettings(trendInfo) {
    const validPeriods = [1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60];
    if (!validPeriods.includes(Number(trendInfo.period))) {
      this.errors.push(`Trend Period: Invalid period (${trendInfo.period})`);
    }

    if (Array.isArray(trendInfo.params)) {
      trendInfo.params.forEach((param, index) => {
        const val = Number(param);
        if (isNaN(val) || val < 0 || val > 9) {
          this.errors.push(`Trend Parameter ${index + 1}: Invalid parameter value (${param})`);
        }
      });

      const nonNoneParams = trendInfo.params.filter(param => Number(param) !== 0);
      const uniqueParams = [...new Set(nonNoneParams.map(Number))];
      if (nonNoneParams.length !== uniqueParams.length) {
        this.warnings.push('Trend Parameters: Duplicate parameters detected');
      }
    }

    return this.errors.length === 0;
  }

  // Event 설정 검증
  validateEventSettings(eventInfo) {
    if (!eventInfo) return false;

    const eventTypes = ['sag', 'swell', 'inter', 'rvc'];
    eventTypes.forEach(type => {
      const action = eventInfo[`${type}_action`];
      if (action !== undefined && ![0, 1, 2].includes(Number(action))) {
        this.errors.push(`${type}: Invalid action value (${action})`);
      }
      const level = eventInfo[`${type}_level`];
      if (level !== undefined) {
        this.validateNumberRange(level, 0, 999999, `${type} Level`);
      }
      const holdoff = eventInfo[`${type}_holdofftime`];
      if (holdoff !== undefined) {
        this.validateNumberRange(holdoff, 0, 999999, `${type} HoldOff time`);
      }
    });

    return this.errors.length === 0;
  }

  // Alarm 설정 검증
  validateAlarmSettings(alarm) {
    if (!alarm) return false;

    for (let i = 1; i <= 32; i++) {
      const alarmConfig = alarm[i];
      if (!alarmConfig || typeof alarmConfig !== 'object') continue;

      const { chan, cond, threshold, hysteresis, action, delay } = alarmConfig;

      if (Number(chan) !== 0) {
        const thresholdNum = Number(threshold);
        if (isNaN(thresholdNum) || thresholdNum <= 0) {
          this.errors.push(`Alarm ${i}: Threshold must be greater than 0 when parameter is active (current: ${threshold})`);
        } else {
          this.validateNumberRange(thresholdNum, 0.01, 999999, `Alarm ${i} Threshold`);
        }
      }

      if (chan !== undefined && (!Number.isInteger(Number(chan)) || Number(chan) < 0)) {
        this.errors.push(`Alarm ${i}: Parameter must be a non-negative integer (current: ${chan})`);
      }

      if (delay !== undefined) {
        const delayNum = Number(delay);
        if (isNaN(delayNum) || delayNum < 0) {
          this.errors.push(`Alarm ${i}: Delay must be a non-negative number (current: ${delay})`);
        }
      }
    }

    return this.errors.length === 0;
  }

  // IP 주소인지 확인하는 헬퍼 함수
  isIPAddress(address) {
    const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
    return ipRegex.test(address);
  }

  // CT Secondary 값 변환 헬퍼
  getCTSecondaryValue(ctType) {
    const ctTypes = {
      '0': 5,      // 5A
      '1': 0.1,    // 100mA
      '2': 0.333,  // 333mV
      '3': 1       // Rogowski (기본값)
    };
    return ctTypes[ctType] || 5;
  }

  // 전체 설정 검증 (General + Main Channel)
  validateAllSettings(generalDict, mainChannelDict) {
    this.reset();

    // 1. General 설정 검증
    this.validateGeneralSettingsInternal(generalDict);

    // 2. Main Channel 검증 (Enable되어 있을 때만)
    if (mainChannelDict && (mainChannelDict.Enable == 1 || mainChannelDict.Enable == true)) {
      const mainErrors = this.errors.length;
      const mainWarnings = this.warnings.length;

      this.validateChannelSettingsInternal(mainChannelDict);

      for (let i = mainErrors; i < this.errors.length; i++) {
        this.errors[i] = `[Main Channel] ${this.errors[i]}`;
      }
      for (let i = mainWarnings; i < this.warnings.length; i++) {
        this.warnings[i] = `[Main Channel] ${this.warnings[i]}`;
      }
    }

    return this.getValidationResult();
  }

  // General 설정 검증 (내부용)
  validateGeneralSettingsInternal(inputDict) {
    if (!inputDict) return false;

    const currentErrors = this.errors.length;
    const currentWarnings = this.warnings.length;

    // 기본 필수 설정들은 항상 검증
    this.validateDeviceInfo(inputDict.deviceInfo);
    this.validateTcpipSettings(inputDict.tcpip);
    this.validateModbusSettings(inputDict.modbus);

    // etc 검증
    if (inputDict.etc) {
      if (inputDict.etc.pf_sign !== undefined && ![0, 1].includes(Number(inputDict.etc.pf_sign))) {
        this.errors.push('PF Sign: Must be 0 (IEC) or 1 (IEEE)');
      }
      if (inputDict.etc.va_type !== undefined && ![0, 1].includes(Number(inputDict.etc.va_type))) {
        this.errors.push('VA Type: Must be 0 (RMS) or 1 (vector)');
      }
    }

    // SNTP 사용 시에만 SNTP 관련 설정 검증
    if (inputDict.useFunction?.sntp === 1) {
      this.validateSNTPSettings(inputDict);
    }

    // MQTT 사용 시 MQTT 설정 검증
    if (inputDict.useFunction?.mqtt === 1) {
      this.validateMQTTSettings(inputDict.MQTT);
    }

    // General 관련 에러/경고에 접두사 추가
    for (let i = currentErrors; i < this.errors.length; i++) {
      this.errors[i] = `[General] ${this.errors[i]}`;
    }
    for (let i = currentWarnings; i < this.warnings.length; i++) {
      this.warnings[i] = `[General] ${this.warnings[i]}`;
    }

    return this.errors.length === currentErrors;
  }

  // Channel 설정 검증 (내부용)
  validateChannelSettingsInternal(channelDict) {
    if (!channelDict) return false;

    const currentErrors = this.errors.length;

    if (channelDict.Enable == 1 || channelDict.Enable == true) {
      this.validateCTSettings(channelDict.ctInfo);
      this.validatePTSettings(channelDict.ptInfo);
      this.validateDemandSettings(channelDict.demand);

      if (channelDict.trendInfo) {
        this.validateTrendSettings(channelDict.trendInfo);
      }

      if (channelDict.eventInfo) {
        this.validateEventSettings(channelDict.eventInfo);
      }

      if (channelDict.alarm) {
        this.validateAlarmSettings(channelDict.alarm);
      }
    }

    return this.errors.length === currentErrors;
  }

  // 검증 결과 반환
  getValidationResult() {
    return {
      isValid: this.errors.length === 0,
      errors: [...this.errors],
      warnings: [...this.warnings],
      hasWarnings: this.warnings.length > 0,
      hasErrors: this.errors.length > 0
    };
  }

  // 경고 전용 메시지 포맷팅
  formatWarningMessage() {
    if (this.warnings.length === 0) return '';
    let message = 'Warnings found:\n';
    message += this.warnings.map(warning => `- ${warning}`).join('\n');
    message += '\n\nDo you want to continue saving?';
    return message;
  }

  // 에러 전용 메시지 포맷팅
  formatErrorOnlyMessage() {
    if (this.errors.length === 0) return '';
    let message = 'Cannot save due to validation errors:\n';
    message += this.errors.map(error => `- ${error}`).join('\n');
    message += '\n\nPlease fix these errors before saving.';
    return message;
  }
}

// 인스턴스 생성 및 export
export const settingValidator = new SettingValidator();
