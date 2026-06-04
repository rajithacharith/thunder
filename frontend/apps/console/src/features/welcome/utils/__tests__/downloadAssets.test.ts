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

import {afterEach, describe, expect, it, vi} from 'vitest';
import type {DistributionAsset, ReleaseAssetInput} from '../../models/download-assets';
import {detectPlatform, groupAssetsByOs, parseDistributionAssets, pickAssetForPlatform} from '../downloadAssets';

const PATTERN = /^sample-app-wayfinder-[0-9A-Za-z.+-]+-(macos|linux|win)-(arm64|x64)\.zip$/i;

function makeAsset(name: string, extra: Partial<ReleaseAssetInput> = {}): ReleaseAssetInput {
  return {downloadUrl: `https://example.com/${name}`, name, sizeLabel: '10 MB', ...extra};
}

describe('parseDistributionAssets', () => {
  it('parses matching linux x64 asset', () => {
    const result = parseDistributionAssets([makeAsset('sample-app-wayfinder-1.0.0-linux-x64.zip')], PATTERN);
    expect(result).toHaveLength(1);
    expect(result[0].os).toBe('linux');
    expect(result[0].arch).toBe('x64');
    expect(result[0].osLabel).toBe('Linux');
    expect(result[0].archLabel).toBe('x64');
  });

  it('parses matching linux arm64 asset', () => {
    const result = parseDistributionAssets([makeAsset('sample-app-wayfinder-1.0.0-linux-arm64.zip')], PATTERN);
    expect(result).toHaveLength(1);
    expect(result[0].os).toBe('linux');
    expect(result[0].arch).toBe('arm64');
    expect(result[0].archLabel).toBe('ARM64');
  });

  it('parses matching macos arm64 asset with Apple Silicon label', () => {
    const result = parseDistributionAssets([makeAsset('sample-app-wayfinder-1.0.0-macos-arm64.zip')], PATTERN);
    expect(result).toHaveLength(1);
    expect(result[0].os).toBe('macos');
    expect(result[0].arch).toBe('arm64');
    expect(result[0].archLabel).toBe('ARM64 (Apple Silicon)');
    expect(result[0].osLabel).toBe('Mac OS');
  });

  it('parses matching macos x64 asset with Intel label', () => {
    const result = parseDistributionAssets([makeAsset('sample-app-wayfinder-1.0.0-macos-x64.zip')], PATTERN);
    expect(result).toHaveLength(1);
    expect(result[0].os).toBe('macos');
    expect(result[0].arch).toBe('x64');
    expect(result[0].archLabel).toBe('x64 (Intel)');
  });

  it('parses matching win x64 asset', () => {
    const result = parseDistributionAssets([makeAsset('sample-app-wayfinder-1.0.0-win-x64.zip')], PATTERN);
    expect(result).toHaveLength(1);
    expect(result[0].os).toBe('win');
    expect(result[0].arch).toBe('x64');
    expect(result[0].osLabel).toBe('Windows');
  });

  it('ignores non-matching asset names', () => {
    const result = parseDistributionAssets(
      [makeAsset('some-other-tool-1.0.0-linux-x64.zip'), makeAsset('README.md')],
      PATTERN,
    );
    expect(result).toHaveLength(0);
  });

  it('returns correct download URL and sizeLabel', () => {
    const asset = makeAsset('sample-app-wayfinder-1.0.0-linux-x64.zip', {
      downloadUrl: 'https://releases.example.com/dl/file.zip',
      sizeLabel: '22 MB',
    });
    const result = parseDistributionAssets([asset], PATTERN);
    expect(result[0].downloadUrl).toBe('https://releases.example.com/dl/file.zip');
    expect(result[0].sizeLabel).toBe('22 MB');
  });

  it('parses multiple assets', () => {
    const assets = [
      makeAsset('sample-app-wayfinder-1.0.0-linux-x64.zip'),
      makeAsset('sample-app-wayfinder-1.0.0-macos-arm64.zip'),
      makeAsset('sample-app-wayfinder-1.0.0-win-x64.zip'),
    ];
    const result = parseDistributionAssets(assets, PATTERN);
    expect(result).toHaveLength(3);
  });
});

describe('pickAssetForPlatform', () => {
  const linuxX64: DistributionAsset = {
    arch: 'x64',
    archLabel: 'x64',
    downloadUrl: 'https://example.com/linux-x64.zip',
    name: 'sample-app-wayfinder-1.0.0-linux-x64.zip',
    os: 'linux',
    osLabel: 'Linux',
    sizeLabel: '10 MB',
  };
  const linuxArm64: DistributionAsset = {
    arch: 'arm64',
    archLabel: 'ARM64',
    downloadUrl: 'https://example.com/linux-arm64.zip',
    name: 'sample-app-wayfinder-1.0.0-linux-arm64.zip',
    os: 'linux',
    osLabel: 'Linux',
    sizeLabel: '11 MB',
  };
  const macosArm64: DistributionAsset = {
    arch: 'arm64',
    archLabel: 'ARM64 (Apple Silicon)',
    downloadUrl: 'https://example.com/macos-arm64.zip',
    name: 'sample-app-wayfinder-1.0.0-macos-arm64.zip',
    os: 'macos',
    osLabel: 'Mac OS',
    sizeLabel: '12 MB',
  };

  it('returns null when assets array is empty', () => {
    expect(pickAssetForPlatform([], {os: 'linux', arch: 'x64'})).toBeNull();
  });

  it('returns exact match when os and arch match', () => {
    const result = pickAssetForPlatform([linuxX64, linuxArm64, macosArm64], {os: 'linux', arch: 'x64'});
    expect(result).toBe(linuxX64);
  });

  it('falls back to same OS when arch does not match', () => {
    const result = pickAssetForPlatform([linuxX64, macosArm64], {os: 'linux', arch: 'arm64'});
    expect(result).toBe(linuxX64);
  });

  it('falls back to first asset when no OS matches', () => {
    const result = pickAssetForPlatform([linuxX64, macosArm64], {os: 'win', arch: 'x64'});
    expect(result).toBe(linuxX64);
  });

  it('returns null only when empty regardless of platform', () => {
    expect(pickAssetForPlatform([], null)).toBeNull();
    expect(pickAssetForPlatform([], {os: 'macos', arch: 'arm64'})).toBeNull();
  });

  it('picks preferred arch when platform has os but no arch', () => {
    const result = pickAssetForPlatform([linuxX64, linuxArm64], {os: 'linux', arch: null});
    // PREFERRED_ARCHS = ['x64', 'arm64'], x64 comes first
    expect(result).toBe(linuxX64);
  });
});

describe('groupAssetsByOs', () => {
  const linuxX64: DistributionAsset = {
    arch: 'x64',
    archLabel: 'x64',
    downloadUrl: 'https://example.com/linux-x64.zip',
    name: 'linux-x64.zip',
    os: 'linux',
    osLabel: 'Linux',
    sizeLabel: '10 MB',
  };
  const winX64: DistributionAsset = {
    arch: 'x64',
    archLabel: 'x64',
    downloadUrl: 'https://example.com/win-x64.zip',
    name: 'win-x64.zip',
    os: 'win',
    osLabel: 'Windows',
    sizeLabel: '11 MB',
  };
  const macosArm64: DistributionAsset = {
    arch: 'arm64',
    archLabel: 'ARM64 (Apple Silicon)',
    downloadUrl: 'https://example.com/macos-arm64.zip',
    name: 'macos-arm64.zip',
    os: 'macos',
    osLabel: 'Mac OS',
    sizeLabel: '12 MB',
  };

  it('groups assets by OS', () => {
    const groups = groupAssetsByOs([linuxX64, winX64, macosArm64]);
    const osList = groups.map((g) => g.os);
    expect(osList).toContain('linux');
    expect(osList).toContain('win');
    expect(osList).toContain('macos');
  });

  it('preferred OS appears first', () => {
    const groups = groupAssetsByOs([linuxX64, winX64, macosArm64], 'macos');
    expect(groups[0].os).toBe('macos');
  });

  it('excludes empty OS groups', () => {
    const groups = groupAssetsByOs([linuxX64, macosArm64]);
    expect(groups.find((g) => g.os === 'win')).toBeUndefined();
    expect(groups).toHaveLength(2);
  });

  it('returns default order when no preferred OS', () => {
    const groups = groupAssetsByOs([linuxX64, winX64, macosArm64]);
    // Default order is linux, win, macos
    expect(groups[0].os).toBe('linux');
    expect(groups[1].os).toBe('win');
    expect(groups[2].os).toBe('macos');
  });
});

describe('detectPlatform', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('returns {os: null, arch: null} when navigator is undefined', async () => {
    vi.stubGlobal('navigator', undefined);
    const result = await detectPlatform();
    expect(result).toEqual({os: null, arch: null});
  });

  it('detects macOS from userAgent', async () => {
    vi.stubGlobal('navigator', {
      platform: 'MacIntel',
      userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)',
    });
    const result = await detectPlatform();
    expect(result.os).toBe('macos');
  });

  it('detects Windows from userAgent', async () => {
    vi.stubGlobal('navigator', {
      platform: 'Win32',
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',
    });
    const result = await detectPlatform();
    expect(result.os).toBe('win');
  });

  it('detects Linux from userAgent', async () => {
    vi.stubGlobal('navigator', {
      platform: 'Linux x86_64',
      userAgent: 'Mozilla/5.0 (X11; Linux x86_64)',
    });
    const result = await detectPlatform();
    expect(result.os).toBe('linux');
  });

  it('uses userAgentData.getHighEntropyValues when available', async () => {
    vi.stubGlobal('navigator', {
      platform: 'Linux x86_64',
      userAgent: 'Mozilla/5.0',
      userAgentData: {
        getHighEntropyValues: vi.fn().mockResolvedValue({
          platform: 'macOS',
          architecture: 'arm',
          bitness: '64',
        }),
      },
    });
    const result = await detectPlatform();
    expect(result.os).toBe('macos');
    expect(result.arch).toBe('arm64');
  });

  it('returns null os when platform is not recognized', async () => {
    vi.stubGlobal('navigator', {
      platform: 'FreeBSD',
      userAgent: 'some unknown browser on freebsd',
    });
    const result = await detectPlatform();
    expect(result.os).toBeNull();
  });

  it('falls back to userAgent detection when getHighEntropyValues throws', async () => {
    vi.stubGlobal('navigator', {
      platform: 'Linux x86_64',
      userAgent: 'Mozilla/5.0 (X11; Linux x86_64)',
      userAgentData: {
        getHighEntropyValues: vi.fn().mockRejectedValue(new Error('not allowed')),
      },
    });
    const result = await detectPlatform();
    expect(result.os).toBe('linux');
  });
});
