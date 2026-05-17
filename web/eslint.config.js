import { defineConfig, globalIgnores } from 'eslint/config'
import globals from 'globals'
import js from '@eslint/js'
import pluginVue from 'eslint-plugin-vue'
import pluginOxlint from 'eslint-plugin-oxlint'
import tseslint from 'typescript-eslint'

export default defineConfig([
  {
    name: 'app/files-to-lint',
    files: ['**/*.{vue,js,mjs,jsx,ts,tsx}'],
  },

  globalIgnores(['**/dist/**', '**/dist-ssr/**', '**/coverage/**']),

  {
    languageOptions: {
      globals: {
        ...globals.browser,
        __APP_VERSION__: 'readonly',
      },
    },
  },

  js.configs.recommended,
  ...tseslint.configs.recommended,
  ...pluginVue.configs['flat/essential'],

  {
    files: ['**/*.vue'],
    languageOptions: {
      parserOptions: {
        parser: tseslint.parser,
      },
    },
  },

  {
    rules: {
      'vue/multi-word-component-names': 'off',
      'vue/valid-v-slot': ['error', { allowModifiers: true }],
      '@typescript-eslint/no-explicit-any': 'off', // Disable for now to allow incremental migration
    }
  },

  ...pluginOxlint.buildFromOxlintConfigFile('.oxlintrc.json'),
])
