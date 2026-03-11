import { createContext, useContext, useMemo, useState } from 'react'

import { Locale, messages } from './messages'

type TranslateParams = Record<string, string | number>

type I18nContextValue = {
  locale: Locale
  setLocale: (locale: Locale) => void
  t: (key: string, params?: TranslateParams) => string
}

const localeStorageKey = 'gms-locale'
const locales: Locale[] = ['zh-CN', 'en-US']

const I18nContext = createContext<I18nContextValue | undefined>(undefined)

function isLocale(value: string | null): value is Locale {
  return value !== null && locales.includes(value as Locale)
}

function resolveInitialLocale(): Locale {
  const saved = typeof window !== 'undefined' ? window.localStorage.getItem(localeStorageKey) : null
  if (isLocale(saved)) {
    return saved
  }

  const browserLanguage = typeof window !== 'undefined' ? window.navigator.language : ''
  if (browserLanguage.toLowerCase().startsWith('zh')) {
    return 'zh-CN'
  }

  return 'en-US'
}

function interpolate(template: string, params?: TranslateParams): string {
  if (!params) {
    return template
  }

  return template.replace(/\{\{\s*(\w+)\s*\}\}/g, (_, key: string) => String(params[key] ?? `{{${key}}}`))
}

export function I18nProvider({ children }: { children: React.ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>(resolveInitialLocale)

  function setLocale(nextLocale: Locale) {
    setLocaleState(nextLocale)
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(localeStorageKey, nextLocale)
    }
  }

  const value = useMemo<I18nContextValue>(() => {
    const localePack = messages[locale]
    const fallbackPack = messages['en-US']

    function t(key: string, params?: TranslateParams): string {
      const template = localePack[key] ?? fallbackPack[key] ?? key
      return interpolate(template, params)
    }

    return {
      locale,
      setLocale,
      t,
    }
  }, [locale])

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>
}

export function useI18n() {
  const context = useContext(I18nContext)
  if (!context) {
    throw new Error('useI18n must be used within I18nProvider')
  }
  return context
}
