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
import {useState} from 'react';
import {describe, it, expect, vi} from 'vitest';
import type {AttestationConfig} from '../../../../models/oauth';
import AttestationSection from '../AttestationSection';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({t: (key: string) => key}),
}));

// Feeds back whatever the section emits as its next prop, mimicking the edit page's
// config round-trip. Guards against the field becoming un-typeable if the round-trip stalls.
function Harness() {
  const [attestation, setAttestation] = useState<AttestationConfig | null | undefined>(undefined);
  return <AttestationSection attestation={attestation ?? undefined} onAttestationChange={setAttestation} />;
}

describe('AttestationSection round-trip', () => {
  it('lets the user type the package name and reflects it back', async () => {
    const user = userEvent.setup();
    render(<Harness />);

    const input = screen.getByLabelText('applications:edit.advanced.attestation.labels.packageName');
    await user.type(input, 'com.example.app');

    expect(input).toHaveValue('com.example.app');
  });

  it('lets the user type the service account credentials', async () => {
    const user = userEvent.setup();
    render(<Harness />);

    const creds = screen.getByLabelText('applications:edit.advanced.attestation.labels.serviceAccountCredentials');
    await user.type(creds, 'abc123');

    expect(creds).toHaveValue('abc123');
  });
});
