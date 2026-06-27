// Token storage for the stateless-JWT API. Access token is a JWT (~7d);
// refresh token is a UUID — the API panics (500) on a non-UUID refresh value,
// so we validate the shape before ever sending it.
const ACCESS_KEY = 'sinmos.accessToken'
const REFRESH_KEY = 'sinmos.refreshToken'

const UUID_RE =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

export const authStorage = {
  getAccess: () => localStorage.getItem(ACCESS_KEY),
  getRefresh: () => localStorage.getItem(REFRESH_KEY),

  set(accessToken: string, refreshToken: string) {
    localStorage.setItem(ACCESS_KEY, accessToken)
    localStorage.setItem(REFRESH_KEY, refreshToken)
  },

  clear() {
    localStorage.removeItem(ACCESS_KEY)
    localStorage.removeItem(REFRESH_KEY)
  },

  hasValidRefresh() {
    const t = localStorage.getItem(REFRESH_KEY)
    return !!t && UUID_RE.test(t)
  },
}
