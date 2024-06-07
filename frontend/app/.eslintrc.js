module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: [
    'plugin:react/recommended', // Uses the recommended rules from @eslint-plugin-react
    'standard-with-typescript',
    'plugin:import/typescript',
    'plugin:prettier/recommended', // Enables eslint-plugin-prettier and eslint-config-prettier. This will display prettier errors as ESLint errors. Make sure this is always the last configuration in the extends array.
    'plugin:@typescript-eslint/recommended', // Uses the recommended rules from the @typescript-eslint/eslint-plugin
    'plugin:react-hooks/recommended',
  ],
  overrides: [],
  parser: '@typescript-eslint/parser', // Specifies the ESLint parser
  parserOptions: {
    ecmaVersion: 2020, // Allows for the parsing of modern ECMAScript features
    sourceType: 'module', // Allows for the use of imports
    project: ['./tsconfig.eslint.json'],
    ecmaFeatures: {
      jsx: true, // Allows for the parsing of JSX
    },
  },
  settings: {
    react: {
      version: 'detect', // Tells eslint-plugin-react to automatically detect the version of React to use
    },
    'import/parsers': {
      '@typescript-eslint/parser': ['.ts', '.tsx'],
    },
    'import/resolver': {
      typescript: {
        project: './tsconfig.eslint.json',
        alwaysTryTypes: true,
      },
    },
  },
  plugins: ['react', '@typescript-eslint', 'prettier'],
  rules: {
    'react/react-in-jsx-scope': ['off'],
    'react/jsx-uses-react': ['off'],
    'react/jsx-props-no-spreading': ['off'],
    'react/no-unescaped-entities': ['off'],
    '@typescript-eslint/no-non-null-assertion': ['off'],
    'import/no-self-import': 'off',
    '@typescript-eslint/consistent-type-assertions': 'off',
    '@typescript-eslint/strict-boolean-expressions': [
      'off',
      {
        allowString: true,
        allowNumber: true,
        allowNullableObject: true,
        allowNullableBoolean: true,
        allowNullableNumber: false,
        allowNullableString: true,
        allowAny: false,
      },
    ],
    'import/order': [
      'error',
      {
        groups: ['builtin', 'external', 'internal', 'index', 'parent', 'sibling'],
      },
    ],
    'react/prop-types': 0,
    'import/no-named-as-default': 'off',
    'import/no-default-export': 'error',
    '@typescript-eslint/explicit-function-return-type': 'off',
    '@typescript-eslint/naming-convention': [
      'error',
      {
        selector: ['enum', 'enumMember'],
        format: ['PascalCase', 'UPPER_CASE'],
      },
      {
        selector: 'function',
        format: ['camelCase', 'PascalCase'],
      },
      {
        // Don't enforce format on some exceptions
        selector: 'variable',
        types: ['boolean'],
        filter: {
          regex: '^loading|disabled|checked$',
          match: true,
        },
        format: null,
      },
      {
        // Ignore rules for underscore-only function parameters like `_` or `__`
        selector: 'parameter',
        filter: {
          regex: '^_+$',
          match: true,
        },
        format: null,
      },
      {
        selector: 'variable',
        format: ['camelCase', 'PascalCase', 'UPPER_CASE'],
      },
      {
        selector: 'parameter',
        format: ['camelCase', 'PascalCase'],
      },
      {
        selector: 'typeLike',
        format: ['PascalCase'],
      },
      {
        selector: 'objectLiteralProperty',
        filter: {
          regex: '^__typename|__html|aria-.+$',
          match: true,
        },
        format: null,
      },
      {
        selector: 'objectLiteralProperty',
        filter: {
          regex: '^[0-9]+$',
          match: true,
        },
        format: null,
      },
      {
        selector: 'interface',
        format: ['PascalCase'],
      },
    ],
  },
};
