import { Route, Router } from '@solidjs/router'
import { Show, createEffect, createResource } from 'solid-js'

import './App.css'
import { Layout } from './components/Layout'
import { AuthContext, makeAuthContext } from './contexts/auth-context'
import { createAPIClient } from './lib/clientUtils'
import { Home } from './pages/Home'
import { Signin } from './pages/Signin'
import { Signup } from './pages/Signup'

export default function App() {
  const [isLogged] = createResource(fetchMe)

  createEffect(() => {
    if (isLogged.loading) {
      return
    }

    if (isLogged.error) {
      return
    }
  })

  return (
    <Show when={!isLogged.loading && !isLogged.error}>
      <AuthContext.Provider value={makeAuthContext(!!isLogged())}>
        <Router root={Layout}>
          <Route path="/" component={Home} />
          <Route path="/signup" component={Signup} />
          <Route path="/signin" component={Signin} />
        </Router>
      </AuthContext.Provider>
    </Show>
  )
}

// function WithAuth(Component: Component) {
//   return () => {
//     const navigate = useNavigate()
//
//     const [isLogged] = useAuthContext()
//
//     createEffect(() => {
//       if (!isLogged()) {
//         navigate('/signin', { replace: true })
//       }
//     })
//
//     return <Component />
//   }
// }

const fetchMe = async () => {
  try {
    const client = createAPIClient()
    const me = await client.account.Me()
    return !!me
  } catch {
    return false
  }
}
