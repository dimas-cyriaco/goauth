import Client, { Environment, Local } from './client'

export function createAPIClient(): Client {
  const clientTarget =
    import.meta.env.VITE_ENCORE_ENVIRONMENT === 'local' ?
      Local
    : Environment('staging')

  return new Client(clientTarget)
}
