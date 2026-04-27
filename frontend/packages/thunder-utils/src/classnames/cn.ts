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

const DEFAULT_PREFIX = 'Thunder';

let prefix: string = DEFAULT_PREFIX;

type ClassValue = string | false | null | undefined | 0;

/**
 * Sets the global class name prefix used by `cn()`.
 *
 * Should be called once at app bootstrap, typically using the product name
 * from `config.js`.
 *
 * @example
 * ```ts
 * setCnPrefix('<PRODUCT_NAME>');
 * ```
 *
 * @param newPrefix - The prefix to use for all class names
 */
export function setCnPrefix(newPrefix: string): void {
  // Strip characters that are invalid in CSS class names to prevent broken selectors.
  prefix = newPrefix.replace(/[^a-zA-Z0-9_-]/g, '');
}

/**
 * Returns the current class name prefix.
 */
export function getCnPrefix(): string {
  return prefix;
}

/**
 * Constructs a className string from conditional class values, automatically
 * prefixing each class with the configured product name (default: "Thunder").
 *
 * Follows MUI-style BEM convention: `{Prefix}{Component}--{slot}`.
 *
 * @example
 * ```tsx
 * cn("SignIn--root")
 * // => "ThunderSignIn--root"
 *
 * cn("SignIn--root", isPrimary && "SignIn--primary")
 * // => "ThunderSignIn--root ThunderSignIn--primary" (when isPrimary is true)
 * // => "ThunderSignIn--root" (when isPrimary is false)
 *
 * cn("SignInBox--root", "SignInBox--paper", isActive && "SignInBox--active")
 * // => "ThunderSignInBox--root ThunderSignInBox--paper ThunderSignInBox--active"
 * ```
 *
 * @param classes - Class name strings or falsy values for conditional classes
 * @returns The joined className string with the configured prefix applied
 */
export default function cn(...classes: ClassValue[]): string {
  return classes
    .filter(Boolean)
    .map((cls) => `${prefix}${cls as string}`)
    .join(' ');
}
