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

import {fireEvent, render, screen} from '@testing-library/react';
import {describe, expect, it, vi} from 'vitest';
import ConnectionForm from '../ConnectionForm';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({t: (key: string) => key}),
  Trans: ({i18nKey}: {i18nKey: string}) => i18nKey,
}));
vi.mock('@thunderid/contexts', () => ({useToast: () => ({showToast: vi.fn()})}));

describe('ConnectionForm', () => {
  const baseProps = {
    type: 'google' as const,
    mode: 'create' as const,
    values: {
      name: '',
      clientId: '',
      clientSecret: '',
      redirectUri: 'https://id.acme.io/oauth/callback/google',
      scopes: '',
    },
    secretReplacing: false,
    hasStoredSecret: false,
    vendorDisplayName: 'Google',
    onFieldChange: vi.fn(),
    onSecretReplacingChange: vi.fn(),
  };

  it('shows field hints by default and replaces them with validation errors after blur', () => {
    render(<ConnectionForm {...baseProps} />);

    expect(screen.getByText('connections:form.fields.clientId.hint')).toBeInTheDocument();
    expect(screen.getByText('connections:form.fields.clientSecret.hint')).toBeInTheDocument();
    expect(screen.getByText('connections:form.fields.scopes.hint')).toBeInTheDocument();

    fireEvent.blur(screen.getByPlaceholderText('1234567890-abc.apps.googleusercontent.com'));

    expect(screen.getByText('connections:validation.required')).toBeInTheDocument();
    expect(screen.queryByText('connections:form.fields.clientId.hint')).not.toBeInTheDocument();
  });

  it('reports field edits through onFieldChange', () => {
    const onFieldChange = vi.fn();
    render(<ConnectionForm {...baseProps} onFieldChange={onFieldChange} />);

    fireEvent.change(screen.getByPlaceholderText('1234567890-abc.apps.googleusercontent.com'), {
      target: {value: 'my-client-id'},
    });

    expect(onFieldChange).toHaveBeenCalledWith('clientId', 'my-client-id');
  });

  it('renders the redirect URI as an editable field that reports edits', () => {
    const onFieldChange = vi.fn();
    render(<ConnectionForm {...baseProps} onFieldChange={onFieldChange} />);

    const input = screen.getByPlaceholderText('https://your-gate-host/gate/callback');
    fireEvent.change(input, {target: {value: 'https://gate.example.com/gate/callback'}});

    expect(onFieldChange).toHaveBeenCalledWith('redirectUri', 'https://gate.example.com/gate/callback');
  });
});
