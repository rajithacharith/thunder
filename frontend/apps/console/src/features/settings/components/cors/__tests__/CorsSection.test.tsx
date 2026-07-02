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

import {fireEvent, screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {renderWithProviders} from '@thunderid/test-utils';
import {describe, it, expect, beforeEach, afterEach, vi} from 'vitest';
import type {CorsConfigResponse} from '../../../models/responses';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({t: (key: string) => key}),
}));

const mockUseGetCorsConfig =
  vi.fn<() => {data: CorsConfigResponse | undefined; isLoading: boolean; error: Error | null}>();
vi.mock('../../../api/useGetCorsConfig', () => ({
  default: () => mockUseGetCorsConfig(),
}));

const mockMutate = vi.fn();
vi.mock('../../../api/useUpdateCorsConfig', () => ({
  default: () => ({mutate: mockMutate, isPending: false}),
}));

const {default: CorsSection} = await import('../CorsSection');

function makeData(overrides?: Partial<CorsConfigResponse>): CorsConfigResponse {
  return {
    readOnly: {allowedOrigins: ['https://console.example.com']},
    writable: {allowedOrigins: ['https://app.acme.com']},
    merged: {allowedOrigins: []},
    ...overrides,
  };
}

describe('CorsSection', () => {
  beforeEach(() => {
    mockUseGetCorsConfig.mockReset();
    mockMutate.mockReset();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('does not render the editor while loading', () => {
    mockUseGetCorsConfig.mockReturnValue({data: undefined, isLoading: true, error: null});
    renderWithProviders(<CorsSection />);
    expect(screen.queryByRole('button', {name: 'settings:cors.addOrigin'})).toBeNull();
  });

  it('shows an alert on load error', () => {
    mockUseGetCorsConfig.mockReturnValue({data: undefined, isLoading: false, error: new Error('load failed')});
    renderWithProviders(<CorsSection />);
    expect(screen.getByRole('alert')).toBeInTheDocument();
  });

  it('renders read-only origins (incl. regex patterns), editable origins, and the Add control', () => {
    mockUseGetCorsConfig.mockReturnValue({
      data: makeData({readOnly: {allowedOrigins: ['https://console.example.com', {regex: '^https://x$'}]}}),
      isLoading: false,
      error: null,
    });
    renderWithProviders(<CorsSection />);

    expect(screen.getByDisplayValue('https://console.example.com')).toHaveAttribute('readonly');
    expect(screen.getByDisplayValue('^https://x$')).toHaveAttribute('readonly');
    expect(screen.getByDisplayValue('https://app.acme.com')).not.toHaveAttribute('readonly');
    expect(screen.getByRole('button', {name: 'settings:cors.addOrigin'})).toBeInTheDocument();
    expect(screen.getByText('settings:cors.readOnlyHint')).toBeInTheDocument();
    expect(screen.getByRole('button', {name: 'settings:cors.removeOrigin'})).toBeInTheDocument();
  });

  it('removes an editable origin when its delete button is clicked', async () => {
    const user = userEvent.setup();
    mockUseGetCorsConfig.mockReturnValue({
      data: makeData({readOnly: {allowedOrigins: []}, writable: {allowedOrigins: ['https://remove.example.com']}}),
      isLoading: false,
      error: null,
    });
    renderWithProviders(<CorsSection />);

    expect(screen.getByDisplayValue('https://remove.example.com')).toBeInTheDocument();
    await user.click(screen.getByRole('button', {name: 'settings:cors.removeOrigin'}));
    expect(screen.queryByDisplayValue('https://remove.example.com')).toBeNull();
  });

  it('adds an editable row when Add origin is clicked', async () => {
    const user = userEvent.setup();
    mockUseGetCorsConfig.mockReturnValue({
      data: makeData({writable: {allowedOrigins: []}}),
      isLoading: false,
      error: null,
    });
    renderWithProviders(<CorsSection />);

    expect(screen.queryByPlaceholderText('settings:cors.originPlaceholder')).toBeNull();
    await user.click(screen.getByRole('button', {name: 'settings:cors.addOrigin'}));
    expect(screen.getByPlaceholderText('settings:cors.originPlaceholder')).toBeInTheDocument();
  });

  it('saves the edited origins via the update mutation and clears the unsaved bar on success', async () => {
    const user = userEvent.setup();
    mockMutate.mockImplementation((...args: unknown[]) => {
      const opts = args[1] as {onSuccess?: () => void} | undefined;
      opts?.onSuccess?.();
    });
    mockUseGetCorsConfig.mockReturnValue({
      data: makeData({readOnly: {allowedOrigins: []}, writable: {allowedOrigins: []}}),
      isLoading: false,
      error: null,
    });
    renderWithProviders(<CorsSection />);

    await user.click(screen.getByRole('button', {name: 'settings:cors.addOrigin'}));
    await user.type(screen.getByPlaceholderText('settings:cors.originPlaceholder'), 'https://new.example.com');

    const saveButton = await screen.findByRole('button', {name: 'settings:cors.save'});
    await user.click(saveButton);

    expect(mockMutate).toHaveBeenCalledWith(
      expect.objectContaining({data: {allowedOrigins: ['https://new.example.com']}}),
      expect.anything(),
    );
    // onSuccess → reset() clears the overlay, so the unsaved bar disappears.
    await waitFor(() => {
      expect(screen.queryByRole('button', {name: 'settings:cors.save'})).toBeNull();
    });
  });

  it('does not save when a row fails the submit-time validation guard', async () => {
    const user = userEvent.setup();
    mockUseGetCorsConfig.mockReturnValue({
      data: makeData({readOnly: {allowedOrigins: []}, writable: {allowedOrigins: []}}),
      isLoading: false,
      error: null,
    });
    renderWithProviders(<CorsSection />);

    await user.click(screen.getByRole('button', {name: 'settings:cors.addOrigin'}));
    // Change without blurring (fireEvent), so the row-level error isn't shown yet and Save is enabled.
    // '(bad' is neither a valid origin nor a compilable regex, so the submit-time guard must block it.
    fireEvent.change(screen.getByPlaceholderText('settings:cors.originPlaceholder'), {target: {value: '(bad'}});
    fireEvent.click(screen.getByRole('button', {name: 'settings:cors.save'}));

    expect(mockMutate).not.toHaveBeenCalled();
  });

  it('blocks Save when a row is a duplicate', async () => {
    const user = userEvent.setup();
    mockUseGetCorsConfig.mockReturnValue({
      data: makeData({readOnly: {allowedOrigins: []}, writable: {allowedOrigins: []}}),
      isLoading: false,
      error: null,
    });
    renderWithProviders(<CorsSection />);

    await user.click(screen.getByRole('button', {name: 'settings:cors.addOrigin'}));
    await user.click(screen.getByRole('button', {name: 'settings:cors.addOrigin'}));
    const inputs = screen.getAllByPlaceholderText('settings:cors.originPlaceholder');
    await user.type(inputs[0], 'https://dup.example.com');
    await user.type(inputs[1], 'https://dup.example.com');

    const saveButton = await screen.findByRole('button', {name: 'settings:cors.save'});
    expect(saveButton).toBeDisabled();
    expect(mockMutate).not.toHaveBeenCalled();
  });
});
