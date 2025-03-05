import Client, { Environment, Local, user } from './client'

type createAPIClientParams = {
  auth?: user.AuthData
  fetcher?: typeof fetch
}

export function createAPIClient(params?: createAPIClientParams): Client {
  const clientTarget =
    import.meta.env.VITE_ENCORE_ENVIRONMENT === 'local' ?
      Local
    : Environment('staging')

  const { auth, fetcher } = params || {}

  return new Client(clientTarget, { auth, fetcher })
}
