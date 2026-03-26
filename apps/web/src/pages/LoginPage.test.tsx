import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it, vi, beforeEach } from 'vitest'

import { LoginPage } from './LoginPage'
import { useAuth } from '../context/AuthContext'
import { useI18n } from '../i18n/I18nContext'

vi.mock('../context/AuthContext', () => ({
  useAuth: vi.fn(),
}))

vi.mock('../i18n/I18nContext', () => ({
  useI18n: vi.fn(),
}))

describe('LoginPage 冒烟测试', () => {
  const loginMock = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()

    vi.mocked(useI18n).mockReturnValue({
      locale: 'zh-CN',
      setLocale: vi.fn(),
      t: (key: string) => {
        const dict: Record<string, string> = {
          'common.language': '语言',
          'common.email': '邮箱',
          'common.password': '密码',
          'common.login': '登录',
          'common.signingIn': '登录中...',
          'login.welcomeBack': '欢迎回来',
          'login.signInToSystem': '登录通用管理系统',
          'login.failed': '登录失败',
        }
        return dict[key] ?? key
      },
    })

    vi.mocked(useAuth).mockReturnValue({
      token: null,
      user: null,
      isReady: true,
      login: loginMock,
      logout: vi.fn(),
    })
  })

  it('应能正常渲染登录页核心元素', () => {
    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>,
    )

    expect(screen.getByText('欢迎回来')).toBeInTheDocument()
    expect(screen.getByText('登录通用管理系统')).toBeInTheDocument()
    expect(screen.getByLabelText('邮箱')).toBeInTheDocument()
    expect(screen.getByLabelText('密码')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '登录' })).toBeInTheDocument()
  })

  it('提交表单时应调用 login', async () => {
    loginMock.mockResolvedValue(undefined)

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>,
    )

    fireEvent.change(screen.getByLabelText('邮箱'), { target: { value: 'tester@gms.local' } })
    fireEvent.change(screen.getByLabelText('密码'), { target: { value: 'secret123' } })
    fireEvent.click(screen.getByRole('button', { name: '登录' }))

    await waitFor(() => {
      expect(loginMock).toHaveBeenCalledWith('tester@gms.local', 'secret123')
    })
  })

  it('登录失败时应显示错误文案', async () => {
    loginMock.mockRejectedValue(new Error('登录失败'))

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>,
    )

    fireEvent.click(screen.getByRole('button', { name: '登录' }))

    expect(await screen.findByText('登录失败')).toBeInTheDocument()
  })
})
