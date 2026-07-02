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

import {screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {renderWithProviders} from '@thunderid/test-utils';
import {describe, it, expect, vi} from 'vitest';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({t: (key: string) => key}),
}));

vi.mock('../../components/cors/CorsSection', () => ({
  default: () => <div data-testid="cors-section" />,
}));

const {default: SettingsPage} = await import('../SettingsPage');

describe('SettingsPage', () => {
  it('renders the title, subtitle, the CORS tab, and the CORS panel', () => {
    renderWithProviders(<SettingsPage />);

    expect(screen.getByText('settings:page.title')).toBeInTheDocument();
    expect(screen.getByText('settings:page.subtitle')).toBeInTheDocument();
    expect(screen.getByRole('tab', {name: 'settings:tabs.cors'})).toBeInTheDocument();
    expect(screen.getByTestId('cors-section')).toBeInTheDocument();
  });

  it('keeps the CORS panel active when its tab is clicked', async () => {
    const user = userEvent.setup();
    renderWithProviders(<SettingsPage />);
    await user.click(screen.getByRole('tab', {name: 'settings:tabs.cors'}));
    expect(screen.getByTestId('cors-section')).toBeInTheDocument();
  });
});
