import { ParentComponent, Show } from 'solid-js'

import { useAuthContext } from '../contexts/auth-context'

export const Layout: ParentComponent = (props) => {
  const [isLogged] = useAuthContext()

  return (
    <>
      <header>
        <nav>
          <ul>
            <li>
              <a
                data-testid="link-to-home"
                href="/"
              >
                <strong>GOAuth</strong>
              </a>
            </li>
          </ul>

          <ul>
            <Show when={!isLogged()}>
              <li>
                <a
                  data-testid="link-to-signup"
                  href="/signup"
                >
                  Signup
                </a>
              </li>

              <li>
                <a
                  data-testid="link-to-signin"
                  href="/signin"
                >
                  Signin
                </a>
              </li>
            </Show>

            <Show when={isLogged()}>
              <li>
                <a
                  data-testid="link-to-signout"
                  href="/signout"
                >
                  Signout
                </a>
              </li>
            </Show>
          </ul>
        </nav>
      </header>

      <main>{props.children}</main>
    </>
  )
}
