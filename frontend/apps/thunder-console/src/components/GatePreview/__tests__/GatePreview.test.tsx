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

import {describe, it, expect, vi} from 'vitest';
import {render, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {OxygenUIThemeProvider} from '@wso2/oxygen-ui';
import type {Theme} from '@thunder/shared-design';
import GatePreview from '../GatePreview';

// Augment the global @asgardeo/react mock (from setup.ts) to include BaseSignIn,
// which GatePreview renders inside the canvas.
vi.mock('@asgardeo/react', () => ({
  AsgardeoProvider: ({children}: {children: React.ReactNode}) => children,
  BaseSignIn: () => <div data-testid="base-sign-in" />,
}));

// A minimal valid theme for rendering the preview (cast to avoid full type scaffolding)
const mockTheme = {
  colorSchemes: {
    light: {palette: {background: {default: '#ffffff'}}},
    dark: {palette: {background: {default: '#121212'}}},
  },
} as unknown as Theme;

function renderWithThemeProvider(ui: React.ReactElement) {
  return render(<OxygenUIThemeProvider>{ui}</OxygenUIThemeProvider>);
}

describe('GatePreview', () => {
  describe('Loading state', () => {
    it('should render a CircularProgress spinner when theme is null', () => {
      renderWithThemeProvider(<GatePreview theme={null} />);

      expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });

    it('should not render the preview canvas when theme is null', () => {
      renderWithThemeProvider(<GatePreview theme={null} />);

      expect(screen.queryByTestId('base-sign-in')).not.toBeInTheDocument();
    });
  });

  describe('Rendering with a valid theme', () => {
    it('should render the preview canvas (BaseSignIn mock) when a theme is provided', () => {
      renderWithThemeProvider(<GatePreview theme={mockTheme} />);

      expect(screen.getByTestId('base-sign-in')).toBeInTheDocument();
    });

    it('should not show a progress spinner when theme is provided', () => {
      renderWithThemeProvider(<GatePreview theme={mockTheme} />);

      expect(screen.queryByRole('progressbar')).not.toBeInTheDocument();
    });
  });

  describe('Toolbar visibility', () => {
    it('should render toolbar viewport controls by default (showToolbar=true)', () => {
      renderWithThemeProvider(<GatePreview theme={mockTheme} />);

      // PreviewToolbar contains icon buttons for mobile/tablet/desktop viewports
      const buttons = screen.getAllByRole('button');
      expect(buttons.length).toBeGreaterThan(0);
    });

    it('should not render toolbar buttons when showToolbar is false', () => {
      renderWithThemeProvider(<GatePreview theme={mockTheme} showToolbar={false} />);

      // Only the base-sign-in mock is rendered; no toolbar icon buttons
      const buttons = screen.queryAllByRole('button');
      expect(buttons).toHaveLength(0);
    });
  });

  describe('Display name', () => {
    it('should show "Preview" in the browser chrome when displayName is not set', () => {
      renderWithThemeProvider(<GatePreview theme={mockTheme} displayName="" />);

      expect(screen.getByText('Preview')).toBeInTheDocument();
    });

    it('should include displayName and "Preview" in the browser chrome', () => {
      renderWithThemeProvider(<GatePreview theme={mockTheme} displayName="My App" />);

      expect(screen.getByText('My App — Preview')).toBeInTheDocument();
    });
  });

  describe('Color scheme', () => {
    it('should render without errors when colorScheme is explicitly set to light', () => {
      renderWithThemeProvider(<GatePreview theme={mockTheme} colorScheme="light" />);

      expect(screen.getByTestId('base-sign-in')).toBeInTheDocument();
    });

    it('should render without errors when colorScheme is explicitly set to dark', () => {
      renderWithThemeProvider(<GatePreview theme={mockTheme} colorScheme="dark" />);

      expect(screen.getByTestId('base-sign-in')).toBeInTheDocument();
    });

    it('should render without errors when syncColorSchemeWithSystem is true', () => {
      renderWithThemeProvider(<GatePreview theme={mockTheme} syncColorSchemeWithSystem />);

      expect(screen.getByTestId('base-sign-in')).toBeInTheDocument();
    });
  });

  describe('Toolbar interactions', () => {
    it('should not crash when a viewport toggle button is clicked', async () => {
      renderWithThemeProvider(<GatePreview theme={mockTheme} />);

      const buttons = screen.getAllByRole('button');
      // Click each toolbar button to exercise viewport and zoom controls
      await Promise.all(
        buttons.map((button) =>
          userEvent.click(button).catch(() => {
            // Some buttons may be disabled; ignore errors
          }),
        ),
      );

      // Preview is still rendered after interactions
      expect(screen.getByTestId('base-sign-in')).toBeInTheDocument();
    });
  });
});
