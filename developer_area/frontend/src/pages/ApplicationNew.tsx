import { useNavigate } from '@solidjs/router'
import { Show, createSignal } from 'solid-js'

import { isAPIError } from '../lib/client'
import { createAPIClient } from '../lib/clientUtils'

export const ApplicationNew = () => {
  const [name, setName] = createSignal('')
  const [errors, setErrors] = createSignal<Record<string, string[]>>()

  const navigate = useNavigate()

  async function onSubmit(event: Event) {
    event.preventDefault()

    const client = createAPIClient()

    try {
      await client.application.Create({ name: name() })
      setErrors({})
      navigate('/applications')
    } catch (error) {
      if (!isAPIError(error)) {
        return
      }

      setErrors(error.details)
    }
  }

  return (
    <>
      <h1>Create Application</h1>

      <Show when={errors()}>
        <small data-testid="login-error">{JSON.stringify(errors())}</small>
      </Show>

      <form onSubmit={onSubmit}>
        <fieldset>
          <label>
            Name
            <input
              aria-label="Name"
              autocomplete="name"
              data-testid="name"
              name="name"
              onChange={(e) => setName(e.target.value)}
              type="name"
              value={name()}
            />
          </label>
        </fieldset>

        <input data-testid="submit" type="submit" value="Signin" />
      </form>
    </>
  )
}
