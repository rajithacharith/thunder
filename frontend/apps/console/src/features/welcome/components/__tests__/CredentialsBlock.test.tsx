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

import {render, screen, fireEvent} from '@thunderid/test-utils';
import {afterEach, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';

vi.mock('@wso2/oxygen-ui-icons-react', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@wso2/oxygen-ui-icons-react')>();
  return {
    ...actual,
    Check: () => <span data-testid="icon-check" />,
    Copy: () => <span data-testid="icon-copy" />,
    Eye: () => <span data-testid="icon-eye" />,
    EyeOff: () => <span data-testid="icon-eye-off" />,
  };
});

import CredentialsBlock from '../CredentialsBlock';

describe('CredentialsBlock', () => {
  it('renders the username', () => {
    render(<CredentialsBlock username="john.doe" password="john.doe" />);
    expect(screen.getByText('john.doe')).toBeInTheDocument();
  });

  it('masks the password by default', () => {
    render(<CredentialsBlock username="john.doe" password="secret" />);
    expect(screen.getByText('••••••••')).toBeInTheDocument();
    expect(screen.queryByText('secret')).not.toBeInTheDocument();
  });

  it('shows the password when the eye button is clicked', () => {
    render(<CredentialsBlock username="john.doe" password="secret" />);
    fireEvent.click(screen.getByRole('button', {name: /show password/i}));
    expect(screen.getByText('secret')).toBeInTheDocument();
  });

  it('hides the password again when eye-off is clicked', () => {
    render(<CredentialsBlock username="john.doe" password="secret" />);
    fireEvent.click(screen.getByRole('button', {name: /show password/i}));
    fireEvent.click(screen.getByRole('button', {name: /hide password/i}));
    expect(screen.getByText('••••••••')).toBeInTheDocument();
  });

  describe('clipboard', () => {
    let writeTextSpy: ReturnType<typeof vi.fn>;

    beforeAll(() => {
      Object.defineProperty(navigator, 'clipboard', {
        value: {writeText: vi.fn()},
        writable: true,
        configurable: true,
      });
    });

    beforeEach(() => {
      writeTextSpy = vi.spyOn(navigator.clipboard, 'writeText').mockResolvedValue(undefined);
    });

    afterEach(() => {
      vi.restoreAllMocks();
    });

    it('copies username on copy button click', () => {
      render(<CredentialsBlock username="john.doe" password="secret" />);
      fireEvent.click(screen.getByRole('button', {name: /copy username/i}));
      expect(writeTextSpy).toHaveBeenCalledWith('john.doe');
    });

    it('copies password on copy button click', () => {
      render(<CredentialsBlock username="john.doe" password="secret" />);
      fireEvent.click(screen.getByRole('button', {name: /copy password/i}));
      expect(writeTextSpy).toHaveBeenCalledWith('secret');
    });
  });
});
