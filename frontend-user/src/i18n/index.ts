import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import ru from './locales/ru.json'
import tj from './locales/tj.json'

i18n.use(initReactI18next).init({
  resources: { ru: { translation: ru }, tj: { translation: tj } },
  lng: localStorage.getItem('mf-lang') ?? 'ru',
  fallbackLng: 'ru',
  interpolation: { escapeValue: false },
})

export default i18n
