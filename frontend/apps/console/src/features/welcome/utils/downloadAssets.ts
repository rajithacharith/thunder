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

export type {
  ArchKey,
  DetectedPlatform,
  DistributionAsset,
  OsKey,
  ReleaseAssetInput,
  ReleaseEntry,
  ReleasesData,
} from '../models/download-assets';
import OS_LABELS from '../constants/download-assets';
import type {ArchKey, DetectedPlatform, DistributionAsset, OsKey, ReleaseAssetInput} from '../models/download-assets';

const VALID_OS = new Set<string>(Object.keys(OS_LABELS));
const VALID_ARCH = new Set<string>(['arm64', 'x64']);

function archLabel(os: OsKey, arch: ArchKey): string {
  if (os === 'macos') return arch === 'arm64' ? 'ARM64 (Apple Silicon)' : 'x64 (Intel)';
  return arch === 'arm64' ? 'ARM64' : 'x64';
}

export function parseDistributionAssets(assets: ReleaseAssetInput[], pattern: RegExp): DistributionAsset[] {
  const result: DistributionAsset[] = [];
  for (const asset of assets) {
    pattern.lastIndex = 0;
    const match = pattern.exec(asset.name);
    if (!match) continue;
    if (!VALID_OS.has(match[1]) || !VALID_ARCH.has(match[2])) continue;
    const os = match[1] as OsKey;
    const arch = match[2] as ArchKey;
    result.push({
      arch,
      archLabel: archLabel(os, arch),
      downloadUrl: asset.downloadUrl,
      name: asset.name,
      os,
      osLabel: OS_LABELS[os],
      sizeLabel: asset.sizeLabel,
    });
  }
  return result;
}

const PREFERRED_ARCHS: ArchKey[] = ['x64', 'arm64'];

export function pickAssetForPlatform(
  assets: DistributionAsset[],
  platform: DetectedPlatform | null,
): DistributionAsset | null {
  if (assets.length === 0) return null;
  if (platform?.arch) {
    return (
      assets.find((a) => a.os === platform.os && a.arch === platform.arch) ??
      assets.find((a) => a.os === platform.os) ??
      assets[0]
    );
  }
  const osAssets = assets.filter((a) => a.os === platform?.os);
  if (osAssets.length > 0) {
    return PREFERRED_ARCHS.map((arch) => osAssets.find((a) => a.arch === arch)).find(Boolean) ?? osAssets[0];
  }
  return assets[0];
}

export function groupAssetsByOs(
  assets: DistributionAsset[],
  preferredOs?: OsKey | null,
): {os: OsKey; assets: DistributionAsset[]}[] {
  const defaultOrder: OsKey[] = ['linux', 'win', 'macos'];
  const order = preferredOs ? [preferredOs, ...defaultOrder.filter((o) => o !== preferredOs)] : defaultOrder;
  return order.map((os) => ({os, assets: assets.filter((a) => a.os === os)})).filter((g) => g.assets.length > 0);
}

interface NavigatorWithUserAgentData extends Navigator {
  userAgentData?: {
    getHighEntropyValues?: (
      hints: ('architecture' | 'bitness' | 'platform')[],
    ) => Promise<{architecture?: string; bitness?: string; platform?: string}>;
  };
}

function detectOperatingSystem(userAgent: string, platform: string): OsKey | null {
  const ua = userAgent.toLowerCase();
  const pf = platform.toLowerCase();
  if (pf.includes('mac') || /(mac os x|macintosh)/.test(ua)) return 'macos';
  if (pf.includes('win') || ua.includes('windows')) return 'win';
  if (pf.includes('linux') || ua.includes('linux')) return 'linux';
  return null;
}

function detectArchitecture(userAgent: string, os: OsKey | null): ArchKey | null {
  const ua = userAgent.toLowerCase();
  if (/(arm64|aarch64|armv8|apple silicon|silicon)/.test(ua)) return 'arm64';
  if (/\b(wow64|win64|x64|x86_64|amd64|intel)\b/.test(ua)) {
    if (os === 'macos') return null;
    return 'x64';
  }
  return null;
}

export async function detectPlatform(): Promise<DetectedPlatform> {
  if (typeof navigator === 'undefined') return {arch: null, os: null};
  const {platform, userAgent} = navigator;
  const fallbackOs = detectOperatingSystem(userAgent, platform);
  const fallback: DetectedPlatform = {arch: detectArchitecture(userAgent, fallbackOs), os: fallbackOs};
  const {userAgentData} = navigator as NavigatorWithUserAgentData;
  if (!userAgentData?.getHighEntropyValues) return fallback;
  try {
    const v = await userAgentData.getHighEntropyValues(['architecture', 'bitness', 'platform']);
    const detPf = v.platform?.toLowerCase() ?? '';
    const detArch = v.architecture?.toLowerCase() ?? '';
    const detBits = v.bitness?.toLowerCase() ?? '';
    const os: OsKey | null =
      detPf === 'macos' ? 'macos' : detPf === 'windows' ? 'win' : detPf === 'linux' ? 'linux' : fallback.os;
    const arch: ArchKey | null =
      detArch === 'arm' ? 'arm64' : detArch === 'x86' && detBits === '64' ? 'x64' : fallback.arch;
    return {arch, os};
  } catch {
    return fallback;
  }
}
