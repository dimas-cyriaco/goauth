import { useContext } from 'solid-js'
import { Show } from 'solid-js'

import { AppContext } from '../App'

export const Layout = (props) => {
  const [auth] = useContext(AppContext)

  return (
    <>
      <header>
        <nav>
          <ul>
            <li>
              <a href="/">
                <strong>GOAuth</strong>
              </a>
            </li>
          </ul>

          <ul>
            <Show when={!auth()}>
              <li>
                <a href="/signup">Signup</a>
              </li>

              <li>
                <a href="/signin">Signin</a>
              </li>
            </Show>

            <Show when={auth()}>
              <li>
                <a href="/signout">Signout</a>
              </li>
            </Show>
          </ul>
        </nav>
      </header>

      <main>{props.children}</main>
    </>
  )
}
