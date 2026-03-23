import { gql } from './useGraphQL'
import { useAuthStore } from '@/stores/useAuthStore'
import { useRouter } from 'vue-router'

const REQUEST_MAGIC_LINK = `
  mutation RequestMagicLink($email: String!) {
    requestMagicLink(email: $email)
  }
`

const VERIFY_MAGIC_LINK = `
  mutation VerifyMagicLink($token: String!) {
    verifyMagicLink(token: $token) {
      sessionId
      user { id email }
    }
  }
`

const LOGOUT = `
  mutation Logout {
    logout
  }
`

const ME = `
  query Me {
    me { id email }
  }
`

export function useAuthAPI() {
  const auth = useAuthStore()
  const router = useRouter()

  async function requestMagicLink(email: string): Promise<void> {
    await gql<{ requestMagicLink: boolean }>(REQUEST_MAGIC_LINK, { email })
    // Sin sesión todavía — el usuario debe hacer click en el enlace
  }

  async function verifyMagicLink(token: string): Promise<void> {
    const data = await gql<{
      verifyMagicLink: { sessionId: string; user: { id: string; email: string } }
    }>(VERIFY_MAGIC_LINK, { token })

    const { sessionId, user } = data.verifyMagicLink
    auth.login(sessionId, user)
    await router.push('/')
  }

  async function logout(): Promise<void> {
    try {
      await gql(LOGOUT, {}, auth.sessionId)
    } finally {
      auth.logout()
      await router.push('/auth')
    }
  }

  async function restoreSession(): Promise<void> {
    if (!auth.sessionId) return
    try {
      const data = await gql<{ me: { id: string; email: string } }>(ME, {}, auth.sessionId)
      if (data.me) {
        auth.user = data.me
      } else {
        auth.logout()
      }
    } catch {
      auth.logout()
    }
  }

  return { requestMagicLink, verifyMagicLink, logout, restoreSession }
}
