import { FormEvent, useState } from 'react'
import { Navigate, useLocation } from 'react-router-dom'

import { useAuth } from '../context/AuthContext'
import { useI18n } from '../i18n/I18nContext'
import { Locale } from '../i18n/messages'

export function LoginPage() {
  const { token, login, isReady } = useAuth()
  const { locale, setLocale, t } = useI18n()
  const location = useLocation()
  const [email, setEmail] = useState('admin@gms.local')
  const [password, setPassword] = useState('admin123')
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const from = (location.state as { from?: { pathname?: string } } | null)?.from?.pathname || '/app/pm'

  if (isReady && token) {
    return <Navigate to={from} replace />
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError(null)
    setIsSubmitting(true)

    try {
      await login(email, password)
    } catch (err) {
      setError(err instanceof Error ? err.message : t('login.failed'))
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="login-page">
      <div className="login-card stack-md">
        <div className="locale-picker row-between">
          <span className="muted">{t('common.language')}</span>
          <select value={locale} onChange={(event) => setLocale(event.target.value as Locale)}>
            <option value="zh-CN">中文</option>
            <option value="en-US">English</option>
          </select>
        </div>

        <div className="stack-sm">
          <h1>{t('login.welcomeBack')}</h1>
          <p className="muted">{t('login.signInToSystem')}</p>
        </div>

        <form onSubmit={onSubmit} className="stack-md">
          <label className="stack-sm">
            <span>{t('common.email')}</span>
            <input
              type="email"
              value={email}
              autoComplete="email"
              onChange={(event) => setEmail(event.target.value)}
              required
            />
          </label>
          <label className="stack-sm">
            <span>{t('common.password')}</span>
            <input
              type="password"
              value={password}
              autoComplete="current-password"
              onChange={(event) => setPassword(event.target.value)}
              required
            />
          </label>
          {error && <p className="error-text">{error}</p>}
          <button type="submit" className="btn-primary" disabled={isSubmitting}>
            {isSubmitting ? t('common.signingIn') : t('common.login')}
          </button>
        </form>
      </div>
    </div>
  )
}
