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

import {render, screen} from '@thunderid/test-utils';
import {describe, expect, it, vi} from 'vitest';
import {MCP_INSPECTOR_URL} from '../../constants/sample-urls';

vi.mock('@wso2/oxygen-ui-icons-react', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@wso2/oxygen-ui-icons-react')>();
  return {
    ...actual,
    ExternalLink: () => <span data-testid="icon-external-link" />,
  };
});

import ExternalLink from '../ExternalLink';

describe('ExternalLink', () => {
  it('renders an anchor with the given href', () => {
    render(<ExternalLink href={MCP_INSPECTOR_URL}>Open Inspector</ExternalLink>);
    const link = screen.getByRole('link', {name: /Open Inspector/});
    expect(link).toHaveAttribute('href', 'http://localhost:6274');
  });

  it('opens in a new tab with noopener noreferrer', () => {
    render(<ExternalLink href={MCP_INSPECTOR_URL}>Link</ExternalLink>);
    const link = screen.getByRole('link', {name: /Link/});
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveAttribute('rel', 'noopener noreferrer');
  });

  it('renders the external link icon', () => {
    render(<ExternalLink href={MCP_INSPECTOR_URL}>Link</ExternalLink>);
    expect(screen.getByTestId('icon-external-link')).toBeInTheDocument();
  });

  it('renders children text', () => {
    render(<ExternalLink href={MCP_INSPECTOR_URL}>Click here</ExternalLink>);
    expect(screen.getByText('Click here')).toBeInTheDocument();
  });

  it('renders without children', () => {
    render(<ExternalLink href={MCP_INSPECTOR_URL} />);
    expect(screen.getByRole('link')).toBeInTheDocument();
  });
});
