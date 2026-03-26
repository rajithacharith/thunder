/**
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import type {Linter} from 'eslint';
import jsxA11y from 'eslint-plugin-jsx-a11y';
import reactPlugin from 'eslint-plugin-react';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import globals from 'globals';
import createParserOptions from '../utils/tsconfig-resolver';

const reactConfig: Linter.Config[] = [
  reactPlugin.configs.flat['recommended'],
  reactPlugin.configs.flat['jsx-runtime'],
  reactHooks.configs.flat.recommended,
  jsxA11y.flatConfigs.recommended as Linter.Config,
  reactRefresh.configs.recommended,
  {
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
      parserOptions: createParserOptions(),
    },
  },
  {
    name: 'thunder/react-settings',
    settings: {
      react: {
        version: 'detect',
      },
      'import/resolver': {
        typescript: {
          alwaysTryTypes: true,
        },
        alias: {
          map: [
            ['@', './src'],
            ['@/components', './src/components'],
            ['@/layouts', './src/layouts'],
            ['@/theme', './src/theme'],
            ['@/contexts', './src/contexts'],
            ['@/lib', './src/lib'],
            ['@/hooks', './src/hooks'],
            ['@/types', './src/types'],
            ['@/test', './src/test'],
          ],
          extensions: ['.js', '.jsx', '.ts', '.tsx', '.json'],
        },
      },
    },
  },
  {
    name: 'thunder/react-overrides',
    rules: {
      // Turn off the requirement to have React in scope for JSX.
      // https://github.com/jsx-eslint/eslint-plugin-react/blob/c9f5eb264e881f7de66188cbb20904fa8edf3985/docs/rules/jsx-use-react.md
      'react/jsx-use-react': 'off',
      // Turn off the requirement to have React in scope for JSX.
      // https://github.com/jsx-eslint/eslint-plugin-react/blob/c9f5eb264e881f7de66188cbb20904fa8edf3985/docs/rules/react-in-jsx-scope.md
      'react/react-in-jsx-scope': 'off',
      // Override the default `airbnb` rule to allow prop spreading in JSX.
      // https://github.com/jsx-eslint/eslint-plugin-react/blob/958954de7422c5c78e8758fa02fc8b6aa2db67ec/docs/rules/jsx-props-no-spreading.md
      'react/jsx-props-no-spreading': 'off',
      // Override the default `airbnb` rule to avoid the deprecated `defaultProps` usage.
      // https://github.com/jsx-eslint/eslint-plugin-react/blob/958954de7422c5c78e8758fa02fc8b6aa2db67ec/docs/rules/require-default-props.md
      'react/require-default-props': [
        'error',
        {
          forbidDefaultForRequired: true,
          classes: 'ignore',
          functions: 'defaultArguments',
        },
      ],
      // Allow imports without file extensions for TypeScript/JavaScript files
      // This is especially useful for path aliases like @/ that resolve to TypeScript files
      'import/extensions': [
        'error',
        'ignorePackages',
        {
          js: 'never',
          jsx: 'never',
          ts: 'never',
          tsx: 'never',
        },
      ],
    },
  },
];

export default reactConfig;
