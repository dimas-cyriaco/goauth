import { useNavigate } from '@solidjs/router'
import { Show, createSignal } from 'solid-js'

import { isAPIError } from '../lib/client'
import { createAPIClient } from '../lib/clientUtils'

export const Signin = () => {
  const [email, setEmail] = createSignal('')
  const [password, setPassword] = createSignal('')
  const [error, setError] = createSignal<string>()

  const navigate = useNavigate()

  async function onSubmit(event: Event) {
    event.preventDefault()

    const client = createAPIClient()

    const body = new FormData()
    body.set('email', email())
    body.set('password', password())

    try {
      await client.user.Login('POST', body)
      setError(undefined)
      navigate('/')
    } catch (error) {
      if (isAPIError(error) && error.status === 401) {
        setError('Wrong email or password.')
        return
      }

      setError('An error occured. Try again later or contact support.')
    }
  }

  return (
    <>
      <h1>Login</h1>

      <Show when={error()}>
        <small data-testid="login-error">{error()}</small>
      </Show>

      <form onSubmit={onSubmit}>
        <fieldset>
          <label>
            Email
            <input
              aria-label="Email"
              autocomplete="email"
              data-testid="email"
              name="email"
              onChange={(e) => setEmail(e.target.value)}
              type="email"
              value={email()}
            />
          </label>
        </fieldset>

        <fieldset>
          <label>
            Password
            <input
              aria-label="Password"
              data-testid="password"
              name="password"
              onChange={(e) => setPassword(e.target.value)}
              type="password"
              value={password()}
            />
          </label>
        </fieldset>

        <input
          data-testid="submit"
          type="submit"
          value="Signin"
        />
      </form>
    </>
  )
}
