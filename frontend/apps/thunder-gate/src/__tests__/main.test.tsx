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

import {describe, it, expect, vi, beforeEach} from 'vitest';

// Mock ReactDOM
const mockRender = vi.fn();
vi.mock('react-dom/client', () => ({
  createRoot: vi.fn(() => ({
    render: mockRender,
  })),
}));

// Mock AppWithConfig
vi.mock('../AppWithConfig', () => ({
  default: () => <div>App</div>,
}));

// Mock i18next
vi.mock('i18next', () => ({
  default: {
    use: vi.fn().mockReturnThis(),
    init: vi.fn().mockResolvedValue(undefined),
  },
}));

// Mock react-i18next
vi.mock('react-i18next', () => ({
  initReactI18next: {},
}));

// Mock @thunder/i18n
vi.mock('@thunder/i18n/locales/en-US', () => ({
  default: {
    common: {},
    navigation: {},
    users: {},
    userTypes: {},
    integrations: {},
    applications: {},
    dashboard: {},
    auth: {},
    mfa: {},
    social: {},
    consent: {},
    errors: {},
  },
}));

// Mock @tanstack/react-query
vi.mock('@tanstack/react-query', () => ({
  QueryClient: vi.fn().mockImplementation(() => ({})),
  QueryClientProvider: ({children}: {children: React.ReactNode}) => children,
}));

// Mock @tanstack/react-query-devtools
vi.mock('@tanstack/react-query-devtools', () => ({
  ReactQueryDevtools: () => null,
}));

// Mock CSS import
vi.mock('../index.css', () => ({}));

describe('main', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Create a mock root element
    const root = document.createElement('div');
    root.id = 'root';
    document.body.appendChild(root);
  });

  it('should have a root element in the document', () => {
    const rootElement = document.getElementById('root');
    expect(rootElement).toBeInTheDocument();
  });
});
