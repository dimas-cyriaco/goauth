import { useNavigate } from '@solidjs/router'
import { Show, createSignal } from 'solid-js'

import { isAPIError, user } from '../lib/client'
import { createAPIClient } from '../lib/clientUtils'

export const Signup = () => {
  const [email, setEmail] = createSignal('')
  const [password, setPassword] = createSignal('')
  const [passwordConfirmation, setPasswordConfirmation] = createSignal('')
  const [hasSubmited, setHasSubmited] = createSignal(false)
  const [errors, setErrors] = createSignal<Record<string, string[]>>({})

  const navigate = useNavigate()

  async function onSubmit(event: Event) {
    event.preventDefault()
    setHasSubmited(true)

    const client = createAPIClient()

    const registrationParams: user.RegistrationParams = {
      email: email(),
      password: password(),
      password_confirmation: passwordConfirmation(),
    }

    try {
      await client.user.Registration(registrationParams)
      setErrors({})
      navigate('/signin')
    } catch (error) {
      if (!isAPIError(error)) {
        return
      }

      setErrors(error.details)
    }
  }

  return (
    <>
      <h1>Create Account</h1>

      <form onSubmit={onSubmit}>
        <fieldset>
          <label>
            Email
            <input
              aria-invalid={hasSubmited() ? !!errors().email : undefined}
              aria-label="Email"
              autocomplete="email"
              data-testid="email"
              name="email"
              onChange={(e) => setEmail(e.target.value)}
              type="email"
              value={email()}
            />
            <Show when={errors().email}>
              <small
                data-testid="email-error"
                id="email-helper"
              >
                {errors().email?.join(' ')}
              </small>
            </Show>
          </label>
        </fieldset>

        <fieldset>
          <label>
            Password
            <input
              aria-invalid={hasSubmited() ? !!errors().password : undefined}
              aria-label="Password"
              data-testid="password"
              name="password"
              onChange={(e) => setPassword(e.target.value)}
              type="password"
              value={password()}
            />
            <Show when={errors().password}>
              <small
                data-testid="password-error"
                id="password-helper"
              >
                {errors().password?.join(' ')}
              </small>
            </Show>
          </label>
        </fieldset>

        <fieldset>
          <label>
            Password Confirmation
            <input
              aria-invalid={
                hasSubmited() ? !!errors().password_confirmation : undefined
              }
              aria-label="PasswordConfirmation"
              data-testid="password-confirmation"
              name="password_confirmation"
              onChange={(e) => setPasswordConfirmation(e.target.value)}
              type="password"
              value={passwordConfirmation()}
            />
            <Show when={errors().password_confirmation}>
              <small
                data-testid="password-confirmation-error"
                id="password_confirmation-helper"
              >
                {errors().password_confirmation?.join(' ')}
              </small>
            </Show>
          </label>
        </fieldset>

        <input
          data-testid="submit"
          type="submit"
          value="Create Account"
        />
      </form>
    </>
  )
}

// const MAP_ENCORE_ERROR_TO_MESSAGES: Record<string, Record<string, string>> = {
//   Email: {
//     taken: 'Email already taken',
//     email: 'Invalid email format',
//   },
//   Password: {
//     min: 'Password too short. Should be at least 6 characters',
//     max: 'Password too long. Should be at most 72 characters',
//   },
//   PasswordConfirmation: {
//     eqcsfield: 'Password confirmation does not match password',
//   },
// }
