import devtools from 'solid-devtools/vite'
import { defineConfig } from 'vite'
// @ts-expect-error this package do not have typings
import solid from 'vite-plugin-solid'

export default defineConfig({
  plugins: [
    solid(),
    devtools({
      /* features options - all disabled by default */
      autoname: true, // e.g. enable autoname
    }),
  ],
})
