import { Route, Router } from '@solidjs/router'
import { createContext, createSignal } from 'solid-js'

import './App.css'
import { Layout } from './components/Layout'
import { Home } from './pages/Home'
import { Signin } from './pages/Signin'
import { Signup } from './pages/Signup'

export const AppContext = createContext()

function App() {
  const [logged, setLogged] = createSignal(false)

  const context = [
    logged,
    {
      login: () => setLogged(true),
      logout: () => setLogged(false),
    },
  ]

  return (
    <AppContext.Provider value={context}>
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
    </AppContext.Provider>
  )
}

export default App
