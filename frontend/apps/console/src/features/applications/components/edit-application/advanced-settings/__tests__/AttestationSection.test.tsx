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

import {render, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {describe, it, expect, vi, beforeEach} from 'vitest';
import type {AttestationConfig} from '../../../../models/oauth';
import AttestationSection from '../AttestationSection';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

describe('AttestationSection', () => {
  const mockOnAttestationChange = vi.fn();

  beforeEach(() => {
    mockOnAttestationChange.mockClear();
  });

  describe('Rendering', () => {
    it('should render the attestation section', () => {
      render(<AttestationSection onAttestationChange={mockOnAttestationChange} />);

      expect(screen.getByText('applications:edit.advanced.labels.attestation')).toBeInTheDocument();
      expect(screen.getByText('applications:edit.advanced.attestation.intro')).toBeInTheDocument();
    });

    it('should render the package name and credentials fields', () => {
      render(<AttestationSection onAttestationChange={mockOnAttestationChange} />);

      expect(screen.getByLabelText('applications:edit.advanced.attestation.labels.packageName')).toBeInTheDocument();
      expect(
        screen.getByLabelText('applications:edit.advanced.attestation.labels.serviceAccountCredentials'),
      ).toBeInTheDocument();
    });

    it('should render the configured package name and digests', () => {
      render(
        <AttestationSection
          attestation={{android: {packageName: 'com.example.app', certificateSha256Digests: ['AA:BB', 'CC:DD']}}}
          onAttestationChange={mockOnAttestationChange}
        />,
      );

      expect(screen.getByDisplayValue('com.example.app')).toBeInTheDocument();
      expect(screen.getByDisplayValue('AA:BB')).toBeInTheDocument();
      expect(screen.getByDisplayValue('CC:DD')).toBeInTheDocument();
    });

    it('should not render the service account credentials value even when configured', () => {
      // The credentials field is write-only; the component never displays a stored value.
      render(
        <AttestationSection
          attestation={{android: {packageName: 'com.example.app', serviceAccountCredentials: 'secret-json'}}}
          onAttestationChange={mockOnAttestationChange}
        />,
      );

      expect(screen.queryByDisplayValue('secret-json')).not.toBeInTheDocument();
    });
  });

  describe('Editing', () => {
    it('should emit an attestation config when the package name is set', async () => {
      const user = userEvent.setup({delay: null});
      render(<AttestationSection onAttestationChange={mockOnAttestationChange} />);

      // The field is controlled by (unchanging) props in this test, so type a single character
      // and assert the emitted config carries it.
      const input = screen.getByLabelText('applications:edit.advanced.attestation.labels.packageName');
      await user.type(input, 'x');

      expect(mockOnAttestationChange).toHaveBeenLastCalledWith({android: {packageName: 'x'}});
    });

    it('should emit null when the only configured value is cleared', async () => {
      const user = userEvent.setup({delay: null});
      render(
        <AttestationSection
          attestation={{android: {packageName: 'x'}}}
          onAttestationChange={mockOnAttestationChange}
        />,
      );

      const input = screen.getByLabelText('applications:edit.advanced.attestation.labels.packageName');
      await user.clear(input);

      expect(mockOnAttestationChange).toHaveBeenLastCalledWith(null);
    });

    it('should add a digest row when Add Digest is clicked', async () => {
      const user = userEvent.setup({delay: null});
      render(<AttestationSection onAttestationChange={mockOnAttestationChange} />);

      await user.click(screen.getByText('applications:edit.advanced.attestation.addDigest'));

      expect(
        screen.getByPlaceholderText('applications:edit.advanced.attestation.placeholder.certificateSha256Digest'),
      ).toBeInTheDocument();
    });

    it('should emit the service account credentials when entered', async () => {
      const user = userEvent.setup({delay: null});
      render(
        <AttestationSection
          attestation={{android: {packageName: 'com.example.app'}}}
          onAttestationChange={mockOnAttestationChange}
        />,
      );

      const creds = screen.getByLabelText('applications:edit.advanced.attestation.labels.serviceAccountCredentials');
      await user.type(creds, '{{"type":"service_account"}');

      const calls = mockOnAttestationChange.mock.calls as [AttestationConfig | null][];
      const lastArg = calls[calls.length - 1][0];
      expect(lastArg?.android?.serviceAccountCredentials).toBe('{"type":"service_account"}');
    });
  });
});
