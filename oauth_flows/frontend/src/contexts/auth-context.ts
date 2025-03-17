import { createContext, createSignal, useContext } from 'solid-js'

export const makeAuthContext = (isLogged: boolean) => {
  const [logged, setLogged] = createSignal<boolean>(isLogged)

  return [
    logged,
    {
      login: () => {
        return setLogged(true)
      },
      logout: () => setLogged(false),
    },
  ] as const
}
export type AuthContextType = ReturnType<typeof makeAuthContext>
export const AuthContext = createContext<AuthContextType>()

export const useAuthContext = () => {
  const authContext = useContext(AuthContext)
  if (!authContext) {
    throw new Error(
      'useCountNameContext should be called inside its ContextProvider',
    )
  }
  return authContext
}
