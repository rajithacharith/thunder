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

import FormFieldsBlock from '../FormFieldsBlock';

const basicFields = [
  {label: 'URL', value: 'http://localhost:8787/mcp'},
  {label: 'Client ID', value: 'MY-CLIENT'},
];

describe('FormFieldsBlock', () => {
  it('renders each field label', () => {
    render(<FormFieldsBlock fields={basicFields} />);
    expect(screen.getByText('URL')).toBeInTheDocument();
    expect(screen.getByText('Client ID')).toBeInTheDocument();
  });

  it('renders each field value', () => {
    render(<FormFieldsBlock fields={basicFields} />);
    expect(screen.getByText('http://localhost:8787/mcp')).toBeInTheDocument();
    expect(screen.getByText('MY-CLIENT')).toBeInTheDocument();
  });

  it('does not render a copy button for readOnly fields', () => {
    render(<FormFieldsBlock fields={[{label: 'Transport', value: 'Streamable HTTP', readOnly: true}]} />);
    expect(screen.queryByRole('button', {name: /copy transport/i})).not.toBeInTheDocument();
  });

  it('renders a copy button for non-readOnly fields', () => {
    render(<FormFieldsBlock fields={[{label: 'URL', value: 'http://localhost'}]} />);
    expect(screen.getByRole('button', {name: /copy url/i})).toBeInTheDocument();
  });

  it('masks password fields by default', () => {
    render(<FormFieldsBlock fields={[{label: 'Secret', value: 'mysecret', isPassword: true}]} />);
    expect(screen.getByText('••••••••')).toBeInTheDocument();
    expect(screen.queryByText('mysecret')).not.toBeInTheDocument();
  });

  it('shows password field value when eye button is clicked', () => {
    render(<FormFieldsBlock fields={[{label: 'Secret', value: 'mysecret', isPassword: true}]} />);
    fireEvent.click(screen.getByRole('button', {name: /show password/i}));
    expect(screen.getByText('mysecret')).toBeInTheDocument();
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

    it('copies the field value on copy button click', () => {
      render(<FormFieldsBlock fields={[{label: 'URL', value: 'http://localhost:8787/mcp'}]} />);
      fireEvent.click(screen.getByRole('button', {name: /copy url/i}));
      expect(writeTextSpy).toHaveBeenCalledWith('http://localhost:8787/mcp');
    });
  });
});
