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

/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */

import {describe, it, expect, vi} from 'vitest';
import {render, screen} from '@testing-library/react';
import DividerAdapter from '../DividerAdapter';

let capturedDividerProps: Record<string, any>;

vi.mock('@wso2/oxygen-ui', () => ({
  Divider: ({children, orientation, sx}: any) => {
    capturedDividerProps = {orientation, sx};
    return (
      <div data-testid="divider" data-orientation={orientation} role="separator">
        {children}
      </div>
    );
  },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({t: (key: string) => key}),
}));

const mockResolve = (s: string | undefined) => s;

describe('DividerAdapter', () => {
  it('renders a horizontal divider by default', () => {
    const component = {id: 'd1', type: 'DIVIDER'} as any;
    render(<DividerAdapter component={component} resolve={mockResolve} />);

    expect(screen.getByTestId('divider')).toBeInTheDocument();
    expect(capturedDividerProps.orientation).toBe('horizontal');
  });

  it('renders a vertical divider when variant is VERTICAL', () => {
    const component = {id: 'd2', type: 'DIVIDER', variant: 'VERTICAL'} as any;
    render(<DividerAdapter component={component} resolve={mockResolve} />);

    expect(capturedDividerProps.orientation).toBe('vertical');
  });

  it('renders a horizontal divider when variant is HORIZONTAL', () => {
    const component = {id: 'd3', type: 'DIVIDER', variant: 'HORIZONTAL'} as any;
    render(<DividerAdapter component={component} resolve={mockResolve} />);

    expect(capturedDividerProps.orientation).toBe('horizontal');
  });

  it('renders label text when component has a label', () => {
    const component = {id: 'd4', type: 'DIVIDER', label: 'OR'} as any;
    render(<DividerAdapter component={component} resolve={mockResolve} />);

    expect(screen.getByText('OR')).toBeInTheDocument();
  });

  it('does not render label text when component has no label', () => {
    const component = {id: 'd5', type: 'DIVIDER'} as any;
    const {container} = render(<DividerAdapter component={component} resolve={mockResolve} />);

    const divider = container.querySelector('[data-testid="divider"]');
    expect(divider?.textContent).toBe('');
  });

  it('resolves label through the resolve function', () => {
    const resolveFn = vi.fn((s: string | undefined) => (s === '{{or}}' ? 'or_label' : s));
    const component = {id: 'd6', type: 'DIVIDER', label: '{{or}}'} as any;
    render(<DividerAdapter component={component} resolve={resolveFn} />);

    expect(resolveFn).toHaveBeenCalledWith('{{or}}');
    expect(screen.getByText('or_label')).toBeInTheDocument();
  });

  it('applies vertical margin styling', () => {
    const component = {id: 'd7', type: 'DIVIDER'} as any;
    render(<DividerAdapter component={component} resolve={mockResolve} />);

    expect(capturedDividerProps.sx).toEqual({my: 2});
  });
});
