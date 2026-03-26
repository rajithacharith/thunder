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

import {render, screen} from '@testing-library/react';
import {describe, it, expect, vi, beforeEach} from 'vitest';
import type {FlowComponent} from '../../../../models/flow';
import ImageAdapter from '../ImageAdapter';

/* eslint-disable @typescript-eslint/no-unused-vars */
const mockIsEmojiUri = vi.fn((_uri: string): boolean => false);
const mockExtractEmojiFromUri = vi.fn((_uri: string): string => '');
/* eslint-enable @typescript-eslint/no-unused-vars */

vi.mock('@asgardeo/react', () => ({
  isEmojiUri: (uri: string) => mockIsEmojiUri(uri),
  extractEmojiFromUri: (uri: string) => mockExtractEmojiFromUri(uri),
}));

vi.mock('@wso2/oxygen-ui', () => ({
  Box: ({component: Comp, src, alt, sx}: any) => {
    if (Comp === 'img') {
      return (
        <img
          data-testid="flow-image"
          src={src}
          alt={alt}
          style={sx}
        />
      );
    }
    return <div style={sx} />;
  },
}));

const baseComponent: FlowComponent = {
  id: 'img-1',
  type: 'IMAGE',
  src: 'https://example.com/logo.png',
  alt: 'Company Logo',
};

const emojiComponent: FlowComponent = {
  id: 'img-emoji',
  type: 'IMAGE',
  src: 'emoji://1f600',
  alt: 'Grinning Face',
};

describe('ImageAdapter', () => {
  beforeEach(() => {
    mockIsEmojiUri.mockReturnValue(false);
    mockExtractEmojiFromUri.mockReturnValue('😀');
  });

  it('renders an img element', () => {
    render(<ImageAdapter component={baseComponent} resolve={(s) => s} />);
    expect(screen.getByTestId('flow-image')).toBeInTheDocument();
  });

  it('uses the resolved src', () => {
    render(<ImageAdapter component={baseComponent} resolve={() => 'https://cdn.example.com/logo.png'} />);
    expect(screen.getByTestId('flow-image')).toHaveAttribute('src', 'https://cdn.example.com/logo.png');
  });

  it('falls back to component.src when resolve returns undefined', () => {
    render(<ImageAdapter component={baseComponent} resolve={() => undefined} />);
    expect(screen.getByTestId('flow-image')).toHaveAttribute('src', 'https://example.com/logo.png');
  });

  it('renders the alt text', () => {
    render(<ImageAdapter component={baseComponent} resolve={(s) => s} />);
    expect(screen.getByAltText('Company Logo')).toBeInTheDocument();
  });

  it('defaults alt to empty string when not provided', () => {
    const component = {...baseComponent, alt: undefined};
    render(<ImageAdapter component={component} resolve={(s) => s} />);
    expect(screen.getByTestId('flow-image')).toHaveAttribute('alt', '');
  });

  it('applies component width and height as px values via sx', () => {
    const component = {...baseComponent, width: '200', height: '100'};
    const {getByTestId} = render(<ImageAdapter component={component} resolve={(s) => s} />);
    const img = getByTestId('flow-image');
    expect(img.style.width).toBe('200px');
    expect(img.style.height).toBe('100px');
  });

  it('defaults to auto width and height when not provided', () => {
    const {getByTestId} = render(<ImageAdapter component={baseComponent} resolve={(s) => s} />);
    const img = getByTestId('flow-image');
    expect(img.style.width).toBe('auto');
    expect(img.style.height).toBe('auto');
  });

  it('applies maxWidth and maxHeight props via sx', () => {
    const {getByTestId} = render(
      <ImageAdapter component={baseComponent} resolve={(s) => s} maxWidth={80} maxHeight={80} />,
    );
    const img = getByTestId('flow-image');
    expect(img.style.maxWidth).toBe('80px');
    expect(img.style.maxHeight).toBe('80px');
  });

  it('defaults maxWidth and maxHeight to 100% when not provided', () => {
    const {getByTestId} = render(<ImageAdapter component={baseComponent} resolve={(s) => s} />);
    const img = getByTestId('flow-image');
    expect(img.style.maxWidth).toBe('100%');
    expect(img.style.maxHeight).toBe('100%');
  });

  describe('emoji rendering', () => {
    beforeEach(() => {
      mockIsEmojiUri.mockReturnValue(true);
      mockExtractEmojiFromUri.mockReturnValue('😀');
    });

    it('renders an emoji span instead of an img when src is an emoji URI', () => {
      render(<ImageAdapter component={emojiComponent} resolve={(s) => s} />);
      expect(screen.getByRole('img', {name: 'Grinning Face'})).toBeInTheDocument();
      expect(screen.queryByTestId('flow-image')).not.toBeInTheDocument();
    });

    it('uses concrete height as containerHeight when height is set', () => {
      const component = {...emojiComponent, height: '80', width: '80'};
      const {container} = render(<ImageAdapter component={component} resolve={(s) => s} />);
      const outer = container.firstChild as HTMLElement;
      expect(outer.style.height).toBe('80px');
    });

    it('uses concrete width as containerHeight when height is auto and width is set', () => {
      const component = {...emojiComponent, height: undefined, width: '60'};
      const {container} = render(<ImageAdapter component={component} resolve={(s) => s} />);
      const outer = container.firstChild as HTMLElement;
      expect(outer.style.height).toBe('60px');
    });

    it('falls back to default emoji container height when both width and height are non-concrete', () => {
      const component = {...emojiComponent, height: undefined, width: undefined};
      const {container} = render(<ImageAdapter component={component} resolve={(s) => s} />);
      const outer = container.firstChild as HTMLElement;
      expect(outer.style.height).toBe('4em');
    });
  });
});
