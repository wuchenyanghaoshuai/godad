module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true
  },
  extends: [
    'eslint:recommended',
    'plugin:vue/vue3-essential',
    'plugin:@typescript-eslint/recommended'
  ],
  parser: 'vue-eslint-parser',
  parserOptions: {
    ecmaVersion: 'latest',
    parser: '@typescript-eslint/parser',
    sourceType: 'module'
  },
  plugins: [
    'vue',
    '@typescript-eslint'
  ],
  rules: {
    'vue/multi-word-component-names': 'off',
    'vue/no-mutating-props': 'off',
    '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],
    '@typescript-eslint/no-explicit-any': 'off',
    '@typescript-eslint/ban-ts-comment': 'warn',
    'no-undef': 'off',
    'no-console': 'off',
    'no-dupe-keys': 'error'
  },
  overrides: [
    {
      files: ['tests/**/*'],
      rules: {
        '@typescript-eslint/no-unused-vars': 'off',
      }
    }
  ],
  globals: {
    NodeJS: true
  }
}
