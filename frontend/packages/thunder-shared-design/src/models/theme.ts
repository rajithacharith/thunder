/**
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
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

/**
 * Color configuration for a specific color type (primary, secondary)
 */
export interface ThemeColor {
  /**
   * Main color value
   * @example "#1976d2"
   */
  main: string;

  /**
   * Dark variant of the color
   * @example "#0d47a1"
   */
  dark: string;

  /**
   * Text color that contrasts with the main color
   * @example "#ffffff"
   */
  contrastText: string;
}

/**
 * Background color configuration
 */
export interface BackgroundColors {
  /**
   * Default background color
   * @example "#FFFFFF"
   */
  default: string;

  /**
   * Paper/surface background color
   * @example "#FBFBFA"
   */
  paper: string;
}

/**
 * Text color configuration
 */
export interface TextColors {
  /**
   * Primary text color
   * @example "#2F3437"
   */
  primary: string;

  /**
   * Secondary text color
   * @example "#6B7280"
   */
  secondary: string;
}

/**
 * Color configuration for a color scheme
 */
export interface ColorSchemeColors {
  /**
   * Primary color configuration
   */
  primary: ThemeColor;

  /**
   * Secondary color configuration
   */
  secondary: ThemeColor;

  /**
   * Background color configuration
   */
  background?: BackgroundColors;

  /**
   * Text color configuration
   */
  text?: TextColors;
}

/**
 * Color scheme configuration for light or dark mode
 */
export interface ColorScheme {
  /**
   * Color configuration for this color scheme
   */
  colors: ColorSchemeColors;
}

/**
 * Shape configuration
 */
export interface ShapeConfig {
  /**
   * Border radius value
   * @example "8px" or 8
   */
  borderRadius: string | number;
}

/**
 * Typography configuration
 */
export interface TypographyConfig {
  /**
   * Font family
   * @example "'Roboto', sans-serif"
   */
  fontFamily: string;
}

/**
 * Theme configuration containing color schemes, shape, and typography
 */
export interface ThemeConfig {
  /**
   * Text direction
   * @example "ltr"
   */
  direction: string;

  /**
   * Default color scheme ("light" or "dark")
   * @example "light"
   */
  defaultColorScheme: string;

  /**
   * Available color schemes
   */
  colorSchemes: {
    light?: ColorScheme;
    dark?: ColorScheme;
  };

  /**
   * Shape configuration
   */
  shape?: ShapeConfig;

  /**
   * Typography configuration
   */
  typography?: TypographyConfig;
}
