import pluginJs from '@eslint/js'
import eslintConfigPrettier from 'eslint-config-prettier'
import solid from 'eslint-plugin-solid/configs/recommended'
import globals from 'globals'
import tseslint from 'typescript-eslint'

/** @type {import('eslint').Linter.Config[]} */
export default [
  { files: ['**/*.{js,mjs,cjs,ts}'] },
  { languageOptions: { globals: globals.browser } },
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
  solid,
  eslintConfigPrettier,
]
