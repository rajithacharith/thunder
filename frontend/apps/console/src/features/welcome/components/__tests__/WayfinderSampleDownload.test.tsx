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

import {render, screen, userEvent, fireEvent} from '@thunderid/test-utils';
import {afterEach, describe, expect, it, vi} from 'vitest';

vi.mock('framer-motion', () => ({
  motion: {
    create: (Component: React.ElementType) => Component,
  },
}));

vi.mock('@wso2/oxygen-ui-icons-react', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@wso2/oxygen-ui-icons-react')>();
  return {
    ...actual,
    ChevronDown: () => <span data-testid="icon-chevron-down" />,
    ChevronUp: () => <span data-testid="icon-chevron-up" />,
    Download: () => <span data-testid="icon-download" />,
  };
});

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, opts?: Record<string, unknown>) => {
      if (opts) {
        return `${key}:${JSON.stringify(opts)}`;
      }
      return key;
    },
  }),
}));

const {mockDetectPlatform, mockUseWayfinderReleases} = vi.hoisted(() => ({
  mockDetectPlatform: vi.fn(),
  mockUseWayfinderReleases: vi.fn(),
}));

vi.mock('../../utils/downloadAssets', async (importOriginal) => {
  const actual = await importOriginal<typeof import('../../utils/downloadAssets')>();
  return {
    ...actual,
    detectPlatform: mockDetectPlatform,
  };
});

vi.mock('../../api/useWayfinderReleases', () => ({
  default: (...args: unknown[]): unknown => mockUseWayfinderReleases(...args),
}));

import WayfinderSampleDownload from '../WayfinderSampleDownload';

const mockAsset = {
  name: 'sample-app-wayfinder-1.0.0-linux-x64.zip',
  downloadUrl: 'https://example.com/sample-app-wayfinder-1.0.0-linux-x64.zip',
  sizeLabel: '10 MB',
};

const mockAssetMacos = {
  name: 'sample-app-wayfinder-1.0.0-macos-arm64.zip',
  downloadUrl: 'https://example.com/sample-app-wayfinder-1.0.0-macos-arm64.zip',
  sizeLabel: '12 MB',
};

const mockReleasesData = {
  latestRelease: {
    tagName: 'v1.0.0',
    assets: [mockAsset],
  },
  releases: [],
};

describe('WayfinderSampleDownload', () => {
  afterEach(() => {
    vi.clearAllMocks();
    vi.unstubAllGlobals();
  });

  it('returns null when isError is true', () => {
    mockDetectPlatform.mockResolvedValue({os: 'linux', arch: 'x64'});
    mockUseWayfinderReleases.mockReturnValue({data: undefined, isError: true});
    const {container} = render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);
    expect(container).toBeEmptyDOMElement();
  });

  it('returns null when assets array is empty', () => {
    mockDetectPlatform.mockResolvedValue({os: 'linux', arch: 'x64'});
    mockUseWayfinderReleases.mockReturnValue({
      data: {latestRelease: {tagName: 'v1.0.0', assets: []}, releases: []},
      isError: false,
    });
    const {container} = render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);
    expect(container).toBeEmptyDOMElement();
  });

  it('returns null when data is not yet loaded', () => {
    mockDetectPlatform.mockResolvedValue({os: 'linux', arch: 'x64'});
    mockUseWayfinderReleases.mockReturnValue({data: undefined, isError: false});
    const {container} = render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);
    expect(container).toBeEmptyDOMElement();
  });

  it('shows download button with OS label', async () => {
    mockDetectPlatform.mockResolvedValue({os: 'linux', arch: 'x64'});
    mockUseWayfinderReleases.mockReturnValue({data: mockReleasesData, isError: false});
    render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);

    expect(await screen.findByText(/common:welcome\.wayfinderSampleDownload\.downloadButton/)).toBeInTheDocument();
  });

  it('shows "Recommended for this device" label when platform matches', async () => {
    mockDetectPlatform.mockResolvedValue({os: 'linux', arch: 'x64'});
    mockUseWayfinderReleases.mockReturnValue({data: mockReleasesData, isError: false});
    render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);

    expect(await screen.findByText('common:welcome.wayfinderSampleDownload.recommendedLabel')).toBeInTheDocument();
  });

  it('shows "Selected download" label when platform does not match exactly', async () => {
    mockDetectPlatform.mockResolvedValue({os: 'macos', arch: 'arm64'});
    mockUseWayfinderReleases.mockReturnValue({data: mockReleasesData, isError: false});
    render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);

    expect(await screen.findByText('common:welcome.wayfinderSampleDownload.selectedLabel')).toBeInTheDocument();
  });

  it('does not show "Other download options" toggle when only one OS group', async () => {
    mockDetectPlatform.mockResolvedValue({os: 'linux', arch: 'x64'});
    mockUseWayfinderReleases.mockReturnValue({data: mockReleasesData, isError: false});
    render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);

    await screen.findByText(/common:welcome\.wayfinderSampleDownload\.downloadButton/);
    expect(screen.queryByText('common:welcome.wayfinderSampleDownload.otherPlatforms')).not.toBeInTheDocument();
  });

  it('shows "Other download options" toggle when multiple OS groups exist', async () => {
    mockDetectPlatform.mockResolvedValue({os: 'linux', arch: 'x64'});
    const multiOsData = {
      latestRelease: {
        tagName: 'v1.0.0',
        assets: [mockAsset, mockAssetMacos],
      },
      releases: [],
    };
    mockUseWayfinderReleases.mockReturnValue({data: multiOsData, isError: false});
    render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);

    expect(await screen.findByText('common:welcome.wayfinderSampleDownload.otherPlatforms')).toBeInTheDocument();
  });

  it('toggling "Other platforms" button shows other download options', async () => {
    mockDetectPlatform.mockResolvedValue({os: 'linux', arch: 'x64'});
    const multiOsData = {
      latestRelease: {
        tagName: 'v1.0.0',
        assets: [mockAsset, mockAssetMacos],
      },
      releases: [],
    };
    mockUseWayfinderReleases.mockReturnValue({data: multiOsData, isError: false});
    const user = userEvent.setup();
    render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);

    const toggle = await screen.findByText('common:welcome.wayfinderSampleDownload.otherPlatforms');
    await user.click(toggle);

    expect(await screen.findByText('common:welcome.wayfinderSampleDownload.hidePlatforms')).toBeInTheDocument();
  });

  it('handles detectPlatform rejection gracefully', async () => {
    mockDetectPlatform.mockRejectedValue(new Error('platform detection failed'));
    mockUseWayfinderReleases.mockReturnValue({data: mockReleasesData, isError: false});
    render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);

    expect(await screen.findByText(/common:welcome\.wayfinderSampleDownload\.downloadButton/)).toBeInTheDocument();
  });

  it('toggles other platforms on Enter keypress', async () => {
    mockDetectPlatform.mockResolvedValue({os: 'linux', arch: 'x64'});
    const multiOsData = {
      latestRelease: {
        tagName: 'v1.0.0',
        assets: [
          {
            name: 'sample-app-wayfinder-1.0.0-linux-x64.zip',
            downloadUrl: 'https://example.com/linux-x64.zip',
            sizeLabel: '10 MB',
          },
          {
            name: 'sample-app-wayfinder-1.0.0-macos-arm64.zip',
            downloadUrl: 'https://example.com/macos-arm64.zip',
            sizeLabel: '12 MB',
          },
        ],
      },
      releases: [],
    };
    mockUseWayfinderReleases.mockReturnValue({data: multiOsData, isError: false});
    render(<WayfinderSampleDownload releasesUrl="https://example.com/releases.json" />);

    const toggle = await screen.findByText('common:welcome.wayfinderSampleDownload.otherPlatforms');
    fireEvent.keyDown(toggle.closest('[role="button"]') ?? toggle, {key: 'Enter'});

    expect(await screen.findByText('common:welcome.wayfinderSampleDownload.hidePlatforms')).toBeInTheDocument();
  });
});
