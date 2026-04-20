import resolveConfig from 'tailwindcss/resolveConfig';
import tailwindConfigFile from '@tailwindConfig'

export const tailwindConfig = () => {
  return resolveConfig(tailwindConfigFile)
}

export const hexToRGB = (h) => {
  let r = 0;
  let g = 0;
  let b = 0;
  if (h.length === 4) {
    r = `0x${h[1]}${h[1]}`;
    g = `0x${h[2]}${h[2]}`;
    b = `0x${h[3]}${h[3]}`;
  } else if (h.length === 7) {
    r = `0x${h[1]}${h[2]}`;
    g = `0x${h[3]}${h[4]}`;
    b = `0x${h[5]}${h[6]}`;
  }
  return `${+r},${+g},${+b}`;
};

export const formatValue = (value) => Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
  maximumSignificantDigits: 3,
  notation: 'compact',
}).format(value);

export const formatThousands = (value) => Intl.NumberFormat('en-US', {
  maximumSignificantDigits: 3,
  notation: 'compact',
}).format(value);

/**
 * IANA 타임존 이름 → UTC 오프셋 (시간 단위)
 * ex) "Asia/Seoul" → 9, "America/New_York" → -5 or -4 (DST)
 * @param {string} timezone - IANA timezone string
 * @returns {number} UTC offset in hours (e.g. 9, -5, 5.5, 5.75)
 */
export const getUtcOffset = (timezone) => {
  const now = new Date();
  const fmt = new Intl.DateTimeFormat('en-US', {
    timeZone: timezone,
    timeZoneName: 'shortOffset',
  });
  const parts = fmt.formatToParts(now);
  const tzPart = parts.find(p => p.type === 'timeZoneName');
  if (!tzPart) return 0;

  const match = tzPart.value.match(/^GMT([+-]?)(\d{1,2})(?::(\d{2}))?$/);
  if (!match) return 0;

  const sign = match[1] === '-' ? -1 : 1;
  const hours = parseInt(match[2], 10);
  const minutes = match[3] ? parseInt(match[3], 10) : 0;
  return sign * (hours + minutes / 60);
};

/**
 * IANA 타임존 이름 → UTC 오프셋 문자열
 * ex) "Asia/Seoul" → "UTC+9", "Asia/Kathmandu" → "UTC+5:45"
 * @param {string} timezone - IANA timezone string
 * @returns {string} formatted offset string
 */
export const getUtcOffsetString = (timezone) => {
  const offset = getUtcOffset(timezone);
  if (offset === 0) return 'UTC';
  const sign = offset >= 0 ? '+' : '-';
  const abs = Math.abs(offset);
  const h = Math.floor(abs);
  const m = Math.round((abs - h) * 60);
  return m === 0 ? `UTC${sign}${h}` : `UTC${sign}${h}:${String(m).padStart(2, '0')}`;
};
