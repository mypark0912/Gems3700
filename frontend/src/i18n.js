import { createI18n } from 'vue-i18n'
import ko from './locales/lang_ko.json'
import en from './locales/lang_en.json'
import ja from './locales/lang_jp.json'

const savedLang = localStorage.getItem('lang') || 'ko'  // ✅ localStorage에 없으면 ko

const i18n = createI18n({
  legacy: false,
  locale: savedLang,
  fallbackLocale: 'en',
  messages: {
    ko,
    en,
    ja,
  },
})

export default i18n
