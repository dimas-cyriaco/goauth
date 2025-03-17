import Client, { Environment, Local, user } from './client'

type createAPIClientParams = {
  auth?: user.AuthData
  fetcher?: typeof fetch
}

export function createAPIClient(params?: createAPIClientParams): Client {
  const clientTarget = import.meta.env.VITE_ENCORE_ENVIRONMENT === 'local' ? Local : Environment('staging')

  const csrfToken = getCookie('csrf_token')

  const fetcher = async (url: RequestInfo | URL, params: RequestInit | undefined) => {
    const { headers, ...rest } = params || {}

    const requestInit: RequestInit = {
      ...rest,
      headers: { ...headers, 'X-CSRF-Token': csrfToken ? csrfToken : '' },
      credentials: 'include',
    }

    return await fetch(url, requestInit)
  }

  const { auth } = params || {}

  return new Client(clientTarget, { auth, fetcher })
}

function getCookie(name: string) {
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) return parts.pop()?.split(';').shift()
  return null
}
