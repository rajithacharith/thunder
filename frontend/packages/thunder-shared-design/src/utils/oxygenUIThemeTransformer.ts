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

import {OxygenTheme, extendTheme} from '@wso2/oxygen-ui';
import type {Theme} from '@wso2/oxygen-ui';
import {ThemeConfig} from '../models/theme';

/* eslint-disable @typescript-eslint/no-unsafe-assignment */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-explicit-any */

/**
 * Safely parses a border radius value to a number
 * Falls back to OxygenTheme default if parsing fails
 *
 * @param value - Border radius value (number or string like "8px")
 * @returns Parsed number or OxygenTheme default
 */
function parseBorderRadius(value: number | string | undefined): number {
  // If already a number, return it
  if (typeof value === 'number') {
    return value;
  }

  // If not a string or empty, use default
  if (typeof value !== 'string' || value.trim() === '') {
    const defaultRadius = OxygenTheme.shape?.borderRadius;
    return typeof defaultRadius === 'number' ? defaultRadius : 8;
  }

  // Remove 'px' suffix and parse
  const parsed = parseInt(value.replace('px', '').trim(), 10);

  // If parsing resulted in NaN or invalid number, use default
  if (Number.isNaN(parsed) || parsed < 0) {
    const defaultRadius = OxygenTheme.shape?.borderRadius;
    return typeof defaultRadius === 'number' ? defaultRadius : 8;
  }

  return parsed;
}

/**
 * Transforms a Thunder Theme object into an OxygenUI theme configuration
 *
 * @param theme - Thunder Theme configuration to transform
 * @returns Extended OxygenUI theme with design colors and styles
 */
export default function oxygenUIThemeTransformer(theme?: ThemeConfig): Theme {
  if (!theme) {
    return OxygenTheme;
  }

  const {colorSchemes} = theme;

  // Extract colors from light and dark schemes
  const lightColors = colorSchemes?.light?.colors;
  const darkColors = colorSchemes?.dark?.colors;

  // Build color schemes for the theme
  const themeColorSchemes: any = {
    ...(OxygenTheme as any).colorSchemes,
  };

  if (lightColors) {
    themeColorSchemes.light = {
      palette: {
        ...((OxygenTheme as any).colorSchemes?.light?.palette ?? {}),
        primary: {
          main: lightColors.primary.main,
          contrastText: lightColors.primary.contrastText,
          ...(lightColors.primary.dark && {dark: lightColors.primary.dark}),
        },
        ...(lightColors.secondary && {
          secondary: {
            main: lightColors.secondary.main,
            contrastText: lightColors.secondary.contrastText,
            ...(lightColors.secondary.dark && {dark: lightColors.secondary.dark}),
          },
        }),
        ...(lightColors.background && {
          background: {
            default: lightColors.background.default,
            paper: lightColors.background.paper,
          },
        }),
        ...(lightColors.text && {
          text: {
            primary: lightColors.text.primary,
            secondary: lightColors.text.secondary,
          },
        }),
      },
    };
  }

  if (darkColors) {
    themeColorSchemes.dark = {
      palette: {
        ...((OxygenTheme as any).colorSchemes?.dark?.palette ?? {}),
        primary: {
          main: darkColors.primary.main,
          contrastText: darkColors.primary.contrastText,
          ...(darkColors.primary.dark && {dark: darkColors.primary.dark}),
        },
        ...(darkColors.secondary && {
          secondary: {
            main: darkColors.secondary.main,
            contrastText: darkColors.secondary.contrastText,
            ...(darkColors.secondary.dark && {dark: darkColors.secondary.dark}),
          },
        }),
        ...(darkColors.background && {
          background: {
            default: darkColors.background.default,
            paper: darkColors.background.paper,
          },
        }),
        ...(darkColors.text && {
          text: {
            primary: darkColors.text.primary,
            secondary: darkColors.text.secondary,
          },
        }),
      },
    };
  }

  // Create static colors for primary theme
  const primaryColor = lightColors?.primary?.main ?? darkColors?.primary?.main;

  return extendTheme({
    ...OxygenTheme,
    ...(theme?.defaultColorScheme && {defaultColorScheme: theme.defaultColorScheme as 'light' | 'dark'}),
    ...(theme?.direction && {direction: theme.direction as 'ltr' | 'rtl'}),
    colorSchemes: themeColorSchemes,
    ...(theme?.shape && {
      shape: {
        borderRadius: parseBorderRadius(theme.shape.borderRadius),
      },
    }),
    ...(theme?.typography && {
      typography: {
        fontFamily: theme.typography.fontFamily,
      },
    }),
    components: {
      ...OxygenTheme.components,
      ...(primaryColor && {
        MuiButton: {
          styleOverrides: {
            ...(OxygenTheme.components?.MuiButton?.styleOverrides ?? {}),
            containedPrimary: {
              '&:not(:disabled)': {
                backgroundColor: primaryColor,
                color: lightColors?.primary?.contrastText ?? darkColors?.primary?.contrastText ?? '#fff',
                '&:hover': {
                  backgroundColor: lightColors?.primary?.dark ?? darkColors?.primary?.dark ?? primaryColor,
                  color: lightColors?.primary?.contrastText ?? darkColors?.primary?.contrastText ?? '#fff',
                },
              },
            },
          },
        },
      }),
    },
  });
}
