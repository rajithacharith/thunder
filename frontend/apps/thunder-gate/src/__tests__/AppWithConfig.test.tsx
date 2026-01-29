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

import {describe, it, expect, vi} from 'vitest';
import {render} from '@thunder/test-utils';
import AppWithConfig from '../AppWithConfig';

// Mock the AppWithTheme component
vi.mock('../AppWithTheme', () => ({
  default: () => <div data-testid="app-with-theme">App With Theme</div>,
}));

// Mock the BrandingProvider
vi.mock('@thunder/shared-branding', () => ({
  BrandingProvider: ({children}: {children: React.ReactNode}) => <div data-testid="branding-provider">{children}</div>,
}));

describe('AppWithConfig', () => {
  it('renders without crashing', () => {
    const {container} = render(<AppWithConfig />);
    expect(container).toBeInTheDocument();
  });

  it('renders AppWithTheme component', () => {
    const {getByTestId} = render(<AppWithConfig />);
    expect(getByTestId('app-with-theme')).toBeInTheDocument();
  });

  it('wraps with BrandingProvider', () => {
    const {getByTestId} = render(<AppWithConfig />);
    expect(getByTestId('branding-provider')).toBeInTheDocument();
  });
});
