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

import {describe, expect, it, beforeEach} from 'vitest';
import cn, {setCnPrefix, getCnPrefix} from '../cn';

describe('cn', () => {
  beforeEach(() => {
    setCnPrefix('Thunder');
  });

  it('should prefix a single class name with the default prefix', () => {
    expect(cn('SignIn--root')).toBe('ThunderSignIn--root');
  });

  it('should join multiple class names with the prefix', () => {
    expect(cn('SignInBox--root', 'SignInBox--paper')).toBe(
      'ThunderSignInBox--root ThunderSignInBox--paper',
    );
  });

  it('should filter out falsy values', () => {
    expect(cn('SignIn--root', false && 'SignIn--primary')).toBe('ThunderSignIn--root');
    expect(cn('SignIn--root', null)).toBe('ThunderSignIn--root');
    expect(cn('SignIn--root', undefined)).toBe('ThunderSignIn--root');
    expect(cn('SignIn--root', 0)).toBe('ThunderSignIn--root');
  });

  it('should include truthy conditional classes', () => {
    expect(cn('SignIn--root', true && 'SignIn--primary')).toBe(
      'ThunderSignIn--root ThunderSignIn--primary',
    );
  });

  it('should return an empty string when no truthy classes are provided', () => {
    expect(cn(false, null, undefined, 0)).toBe('');
  });

  it('should return an empty string when called with no arguments', () => {
    expect(cn()).toBe('');
  });
});

describe('setCnPrefix / getCnPrefix', () => {
  beforeEach(() => {
    setCnPrefix('Thunder');
  });

  it('should change the prefix used by cn()', () => {
    setCnPrefix('Asgardeo');
    expect(cn('SignIn--root')).toBe('AsgardeoSignIn--root');
  });

  it('should return the current prefix via getCnPrefix()', () => {
    expect(getCnPrefix()).toBe('Thunder');
    setCnPrefix('Custom');
    expect(getCnPrefix()).toBe('Custom');
  });
});
