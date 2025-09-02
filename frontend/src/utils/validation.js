// utils/validation.js
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

  // 파일 경로 형식 검증
  validateFilePath(path, fieldName) {
    if (!path || typeof path !== 'string') {
      this.errors.push(`${fieldName}: File path is required`);
      return false;
    }
    
    // 기본적인 파일 경로 형식 검증
    const pathRegex = /^[a-zA-Z0-9._\-\/\\]+$/;
    if (!pathRegex.test(path)) {
      this.errors.push(`${fieldName}: Invalid file path format. Only alphanumeric characters, dots, hyphens, underscores, and slashes are allowed`);
      return false;
    }
    
    // 경로가 /로 시작하는지 확인 (절대 경로)
    if (!path.startsWith('/')) {
      this.warnings.push(`${fieldName}: File path should start with '/' for absolute path`);
    }
    
    return true;
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
    //this.validateRequired(deviceInfo.location, 'Device Location');
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
    //this.validateNumberRange(modbus.tcp_timeout, 1000, 30000, 'TCP Timeout (ms)');
    this.validateNumberRange(modbus.modbus_id, 1, 247, 'Modbus RTU ID');
    
    const validBaudRates = [0,1,2,3,4];
    if (!validBaudRates.includes(modbus.baud_rate)) {
      this.errors.push(`Baud Rate: Invalid baud rate (${modbus.baud_rate})`);
    }
    const validParity = [0,1,2];
    if (!validParity.includes(modbus.parity)) {
      this.errors.push(`Parity: Invalid parity (${modbus.parity})`);
    }
    const validDataBits = [7,8];
    if (!validDataBits.includes(modbus.data_bits)) {
      this.errors.push(`Data Bits: Invalid data bits (${modbus.data_bits})`);
    }
    const validStopBits = [1,2];
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
    
    if (![50, 60].includes(Number(ptInfo.linefrequency))) {
      this.errors.push('Line Frequency: Must be 50Hz or 60Hz');
    }
    
    if (![0, 1].includes(Number(ptInfo.wiringmode))) {
      this.errors.push('Wiring Mode: Must be 0 (3P4W) or 1 (3P3W)');
    }
    
    return this.errors.length === 0;
  }

  // Sampling 설정 검증
  validateSamplingSettings(sampling) {
    if (!sampling) return false;
    const validRates = [8000, 4000, 2000, 1000];
    const validDurations = [5, 10, 15, 20, 25, 30];
    const validPeriods = [1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60];
    
    if (!validRates.includes(Number(sampling.rate))) {
      this.errors.push(`Sampling Rate: Invalid rate (${sampling.rate})`);
    }
    
    if (!validDurations.includes(Number(sampling.duration))) {
      this.errors.push(`Sample Duration: Invalid duration (${sampling.duration})`);
    }
    
    if (!validPeriods.includes(Number(sampling.period))) {
      this.errors.push(`Sampling Interval: Invalid period (${sampling.period})`);
    }
    
    return this.errors.length === 0;
  }

  // Demand 설정 검증
  validateDemandSettings(demand) {
    if (!demand) return false;
    
    // Demand Target 검증
    this.validateNumberRange(demand.target, 0, 999999999, 'Demand Target');
    
    // Demand Interval 검증
    const validIntervals = [1, 5, 15, 30, 60];
    if (!validIntervals.includes(Number(demand.demand_interval))) {
      this.errors.push(`Demand Interval: Invalid interval (${demand.demand_interval}). Must be one of: ${validIntervals.join(', ')}`);
    }
    
    return this.errors.length === 0;
  }

   // FTP 설정 검증
  validateFTPSettings(inputDict) {
    if (inputDict.ftpInfo) {
      this.validateIPAddress(inputDict.ftpInfo.host, 'FTP HOST');
      this.validateRequired(inputDict.ftpInfo.id, 'FTP Account');
      this.validateRequired(inputDict.ftpInfo.pass, 'FTP Password');
      if (inputDict.ftpInfo.port) {
        this.validatePort(inputDict.ftpInfo.port, 'FTP Port');
      }
      
      // Manual 업로드가 활성화된 경우 파일 경로 검증
      if (inputDict.ftpInfo.isManual === true || inputDict.ftpInfo.isManual === 1) {
        if (inputDict.ftpInfo.upload_main) {
          this.validateFilePath(inputDict.ftpInfo.upload_main, 'FTP Main Channel Upload Path');
        }
        if (inputDict.ftpInfo.upload_sub) {
          this.validateFilePath(inputDict.ftpInfo.upload_sub, 'FTP Sub Channel Upload Path');
        }
      }
    } else {
      this.warnings.push('FTP is enabled but FTP settings are not configured');
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

  // Asset 설정 검증
  validateAssetSettings(assetInfo) {
    this.validateRequired(assetInfo.name, 'Asset Name');
    
    const namePattern = /^[a-zA-Z0-9_-]+$/;
    if (assetInfo.name && !namePattern.test(assetInfo.name)) {
      this.errors.push('Asset Name: Only alphanumeric characters, underscore, and hyphen are allowed');
    }
    
    return this.errors.length === 0;
  }

  // Power Quality 설정 검증
  validatePowerQualitySettings(inputDict) {
    if (inputDict.pqThresholds) {
      if (inputDict.pqThresholds.sagThreshold) {
        this.validateNumberRange(inputDict.pqThresholds.sagThreshold, 50, 90, 'Sag Threshold (%)');
      }
      if (inputDict.pqThresholds.swellThreshold) {
        this.validateNumberRange(inputDict.pqThresholds.swellThreshold, 110, 150, 'Swell Threshold (%)');
      }
      if (inputDict.pqThresholds.thdThreshold) {
        this.validateNumberRange(inputDict.pqThresholds.thdThreshold, 1, 50, 'THD Threshold (%)');
      }
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
      const validParams = [
        "None", "Temperature", "Frequency", "Line Voltage", "Phase Voltage",
        "Current", "Unbalance", "PF", "THD", "TDD", "Power", "Energy", "Demand"
      ];
      
      trendInfo.params.forEach((param, index) => {
        if (!validParams.includes(param)) {
          this.errors.push(`Trend Parameter ${index + 1}: Invalid parameter (${param})`);
        }
      });
      
      const nonNoneParams = trendInfo.params.filter(param => param !== "None");
      const uniqueParams = [...new Set(nonNoneParams)];
      if (nonNoneParams.length !== uniqueParams.length) {
        this.warnings.push('Trend Parameters: Duplicate parameters detected');
      }
    }
    
    return this.errors.length === 0;
  }

  // Event 설정 검증
  validateEventSettings(eventInfo) {
    const eventTypes = ['voltage', 'current', 'frequency', 'power'];
    
    eventTypes.forEach(type => {
      if (eventInfo[type]) {
        const typeInfo = eventInfo[type];
        
        if (typeInfo.upperLimit !== undefined) {
          this.validateNumberRange(typeInfo.upperLimit, 0, 999999, `${type} Upper Limit`);
        }
        if (typeInfo.lowerLimit !== undefined) {
          this.validateNumberRange(typeInfo.lowerLimit, 0, 999999, `${type} Lower Limit`);
        }
        
        if (typeInfo.upperLimit !== undefined && typeInfo.lowerLimit !== undefined) {
          if (Number(typeInfo.upperLimit) <= Number(typeInfo.lowerLimit)) {
            this.errors.push(`${type}: Upper limit must be greater than lower limit`);
          }
        }
        
        if (typeInfo.duration !== undefined) {
          this.validateNumberRange(typeInfo.duration, 1, 3600, `${type} Duration (seconds)`);
        }
      }
    });
    
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

  // 전체 설정 검증 (General + 모든 Channel)
  validateAllSettings(generalDict, mainChannelDict, subChannelDict) {
    this.reset();
    
    // 1. General 설정 검증
    this.validateGeneralSettingsInternal(generalDict);
    
    // 2. Main Channel 검증 (Enable되어 있을 때만)
    if (mainChannelDict && mainChannelDict.Enable === 1) {
      const mainErrors = this.errors.length;
      const mainWarnings = this.warnings.length;
      
      this.validateChannelSettingsInternal(mainChannelDict);
      
      // Main Channel 에러/경고에 접두사 추가
      for (let i = mainErrors; i < this.errors.length; i++) {
        this.errors[i] = `[Main Channel] ${this.errors[i]}`;
      }
      for (let i = mainWarnings; i < this.warnings.length; i++) {
        this.warnings[i] = `[Main Channel] ${this.warnings[i]}`;
      }
    }
    
    // 3. Sub Channel 검증 (Enable되어 있을 때만)
    if (subChannelDict && subChannelDict.Enable === 1) {
      const subErrors = this.errors.length;
      const subWarnings = this.warnings.length;
      
      this.validateChannelSettingsInternal(subChannelDict);
      
      // Sub Channel 에러/경고에 접두사 추가
      for (let i = subErrors; i < this.errors.length; i++) {
        this.errors[i] = `[Sub Channel] ${this.errors[i]}`;
      }
      for (let i = subWarnings; i < this.warnings.length; i++) {
        this.warnings[i] = `[Sub Channel] ${this.warnings[i]}`;
      }
    }
    
    // 4. 크로스 체크 (General과 Channel 간의 충돌 검사)
    this.validateCrossSettings(generalDict, mainChannelDict, subChannelDict);
    
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
    

    //ETC
    // PF Sign 검증
    if (inputDict.pf_sign !== undefined && ![0, 1].includes(Number(inputDict.pf_sign))) {
      this.errors.push('PF Sign: Must be 0 (IEC) or 1 (IEEE)');
    }
    
    // VA Type 검증
    if (inputDict.va_type !== undefined && ![0, 1].includes(Number(inputDict.va_type))) {
      this.errors.push('VA Type: Must be 0 (RMS) or 1 (vector)');
    }
    
    // FTP 사용 시에만 FTP 관련 설정 검증
    if (inputDict.useFuction?.ftp === 1) {
      this.validateFTPSettings(inputDict);
    }
    
    // SNTP 사용 시에만 SNTP 관련 설정 검증
    if (inputDict.useFuction?.sntp === 1) {
      this.validateSNTPSettings(inputDict);
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
  validateChannelSettingsInternal(inputDict) {
    if (!inputDict) return false;
    
    const currentErrors = this.errors.length;
    
    // 채널이 Enable되어 있을 때만 세부 설정 검증
    if (inputDict.Enable === 1) {
      this.validateCTSettings(inputDict.ctInfo);
      this.validatePTSettings(inputDict.ptInfo);
      this.validateSamplingSettings(inputDict.sampling);
      this.validateDemandSettings(inputDict.demand);
      
      // Asset 정보가 있을 때만 검증
      if (inputDict.assetInfo && inputDict.assetInfo.name) {
        this.validateAssetSettings(inputDict.assetInfo);
      }
      
      // Power Quality가 활성화되어 있을 때 추가 검증
      if (inputDict.PowerQuality === 1) {
        this.validatePowerQualitySettings(inputDict);
      }
      
      // Trend 설정 검증
      if (inputDict.trendInfo) {
        this.validateTrendSettings(inputDict.trendInfo);
      }
      
      // Event 설정 검증
      if (inputDict.eventInfo) {
        this.validateEventSettings(inputDict.eventInfo);
      }
    }
    
    return this.errors.length === currentErrors;
  }

  // General과 Channel 간의 크로스 검증
  validateCrossSettings(generalDict, mainChannelDict, subChannelDict) {
    // FTP와 Diagnosis 동시 사용 체크
    const isFTPEnabled = generalDict?.useFuction?.ftp === 1;
    const isMainDiagnosisEnabled = generalDict?.useFuction?.diagnosis_main === true || generalDict?.useFuction?.diagnosis_main === 1;
    const isSubDiagnosisEnabled = generalDict?.useFuction?.diagnosis_sub === true || generalDict?.useFuction?.diagnosis_sub === 1;
    const isMainChannelEnabled = mainChannelDict?.Enable === 1;
    const isSubChannelEnabled = subChannelDict?.Enable === 1;
    
    const isMainDiagnosisActive = isMainDiagnosisEnabled && isMainChannelEnabled;
    const isSubDiagnosisActive = isSubDiagnosisEnabled && isSubChannelEnabled;
    const isAnyDiagnosisEnabled = isMainDiagnosisActive || isSubDiagnosisActive;
    
    if (isFTPEnabled && isAnyDiagnosisEnabled) {
      this.errors.push('Cannot use Waveform FTP and Diagnosis simultaneously. Please disable one of them.');
    }
    
    
    // 주파수 설정 일치 검사
    if (isMainChannelEnabled && isSubChannelEnabled) {
      const mainFreq = mainChannelDict?.ptInfo?.linefrequency;
      const subFreq = subChannelDict?.ptInfo?.linefrequency;
      
      if (mainFreq && subFreq && mainFreq !== subFreq) {
        this.warnings.push(`Line frequency mismatch: Main Channel (${mainFreq}Hz) vs Sub Channel (${subFreq}Hz)`);
      }
    }
    
    return this.errors.length === 0;
  }

  // 검증 결과 반환 (운영 모드)
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
    
    let message = '⚠️ Warnings found:\n';
    message += this.warnings.map(warning => `• ${warning}`).join('\n');
    message += '\n\nDo you want to continue saving?';
    
    return message;
  }

  // 에러 전용 메시지 포맷팅  
  formatErrorOnlyMessage() {
    if (this.errors.length === 0) return '';
    
    let message = '❌ Cannot save due to validation errors:\n';
    message += this.errors.map(error => `• ${error}`).join('\n');
    message += '\n\nPlease fix these errors before saving.';
    
    return message;
  }
}

// 인스턴스 생성 및 export
export const settingValidator = new SettingValidator();