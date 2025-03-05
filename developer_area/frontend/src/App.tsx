import { Route, Router } from '@solidjs/router'
import { Show, createEffect, createResource } from 'solid-js'

import './App.css'
import { Layout } from './components/Layout'
import { AuthContext, makeAuthContext } from './contexts/auth-context'
import { createAPIClient } from './lib/clientUtils'
import { Home } from './pages/Home'
import { Signin } from './pages/Signin'
import { Signup } from './pages/Signup'

function getCookie(name: string) {
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) return parts.pop()?.split(';').shift()
  return null
}

const fetchMe = async () => {
  try {
    const csrfToken = getCookie('csrf_token')

    const fetcher = async (url: RequestInfo | URL, _: unknown) => {
      return await fetch(url, {
        method: 'GET',
        headers: {
          'X-CSRF-Token': csrfToken ? csrfToken : '',
        },
        credentials: 'include',
      })
    }

    const client = createAPIClient({
      fetcher,
    })

    const me = await client.user.Me()
    return !!me
  } catch (error) {
    return false
  }
}

function App() {
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
          <Route
            path="/"
            component={Home}
          />
          <Route
            path="/signup"
            component={Signup}
          />
          <Route
            path="/signin"
            component={Signin}
          />
        </Router>
      </AuthContext.Provider>
    </Show>
  )
}

export default App
