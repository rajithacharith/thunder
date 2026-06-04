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

import {render, screen, userEvent, waitFor, fireEvent} from '@thunderid/test-utils';
import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';

const {mockDetectPlatform} = vi.hoisted(() => ({mockDetectPlatform: vi.fn()}));
const mockLocalStorageGetItem = vi.fn();
const mockLocalStorageSetItem = vi.fn();
const mockSessionStorageGetItem = vi.fn();
const mockSessionStorageSetItem = vi.fn();

vi.mock('@thunderid/contexts', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@thunderid/contexts')>();
  return {
    ...actual,
    useConfig: () => ({
      config: {
        brand: {
          product_name: 'ThunderID',
          docs_url: 'https://docs.example.com/',
        },
      },
    }),
  };
});

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, opts?: Record<string, unknown>) => {
      if (opts?.productName) return `${key}:${opts.productName as string}`;
      return key;
    },
  }),
}));

vi.mock('../../utils/downloadAssets', async (importOriginal) => {
  const actual = await importOriginal<typeof import('../../utils/downloadAssets')>();
  return {...actual, detectPlatform: mockDetectPlatform};
});

vi.mock('../TerminalBlock', () => ({
  default: ({command, tabs}: {command: string; tabs?: React.ReactNode}) => (
    <div>
      {tabs}
      <pre data-testid="terminal-block">{command}</pre>
    </div>
  ),
}));

vi.mock('../WayfinderFolderImport', () => ({
  default: ({onSuccess}: {onSuccess?: () => void}) => (
    <div data-testid="wayfinder-folder-import">
      <button onClick={onSuccess} type="button">
        mock-import-success
      </button>
    </div>
  ),
}));

vi.mock('../WayfinderSampleDownload', () => ({
  default: () => <div data-testid="wayfinder-sample-download" />,
}));

vi.mock('@wso2/oxygen-ui-icons-react', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@wso2/oxygen-ui-icons-react')>();
  return {
    ...actual,
    CheckCircle: () => <span data-testid="icon-check-circle" />,
    ChevronRight: () => <span data-testid="icon-chevron-right" />,
    Database: () => <span data-testid="icon-database" />,
    Download: () => <span data-testid="icon-download" />,
    Play: () => <span data-testid="icon-play" />,
    Settings: () => <span data-testid="icon-settings" />,
  };
});

import WayfinderSampleSetup from '../WayfinderSampleSetup';

describe('WayfinderSampleSetup', () => {
  beforeEach(() => {
    mockDetectPlatform.mockResolvedValue({os: 'linux', arch: 'x64'});
    vi.stubGlobal('localStorage', {
      getItem: mockLocalStorageGetItem,
      setItem: mockLocalStorageSetItem,
      removeItem: vi.fn(),
      clear: vi.fn(),
    });
    vi.stubGlobal('sessionStorage', {
      getItem: mockSessionStorageGetItem,
      setItem: mockSessionStorageSetItem,
      removeItem: vi.fn(),
      clear: vi.fn(),
    });
    mockLocalStorageGetItem.mockReturnValue(null);
    mockSessionStorageGetItem.mockReturnValue(null);
  });

  afterEach(() => {
    vi.clearAllMocks();
    vi.unstubAllGlobals();
  });

  it('renders without crashing', () => {
    const {container} = render(<WayfinderSampleSetup />);
    expect(container).toBeInTheDocument();
  });

  it('renders the setup title', () => {
    render(<WayfinderSampleSetup />);
    expect(screen.getByText('common:welcome.wayfinderSampleSetup.title')).toBeInTheDocument();
  });

  it('renders the one-time setup badge', () => {
    render(<WayfinderSampleSetup />);
    expect(screen.getByText('common:welcome.wayfinderSampleSetup.oneTimeSetup')).toBeInTheDocument();
  });

  it('expands by default when not previously imported', () => {
    mockLocalStorageGetItem.mockReturnValue(null);
    render(<WayfinderSampleSetup />);
    expect(screen.getByTestId('wayfinder-folder-import')).toBeInTheDocument();
    expect(screen.getByTestId('terminal-block')).toBeInTheDocument();
  });

  it('collapses by default when previously imported', () => {
    mockLocalStorageGetItem.mockReturnValue('1234567890');
    render(<WayfinderSampleSetup />);
    expect(screen.queryByTestId('wayfinder-folder-import')).not.toBeInTheDocument();
  });

  it('shows setupComplete message when done and collapsed', () => {
    mockLocalStorageGetItem.mockReturnValue('1234567890');
    mockSessionStorageGetItem.mockReturnValue('false');
    render(<WayfinderSampleSetup />);
    expect(screen.getByText('common:welcome.wayfinderSampleSetup.setupComplete')).toBeInTheDocument();
  });

  it('toggles expand/collapse when header is clicked', async () => {
    const user = userEvent.setup();
    render(<WayfinderSampleSetup />);

    expect(screen.getByTestId('terminal-block')).toBeInTheDocument();

    const header = screen.getByRole('button', {name: /wayfinderSampleSetup\.title/i});
    await user.click(header);

    expect(screen.queryByTestId('terminal-block')).not.toBeInTheDocument();
    expect(mockSessionStorageSetItem).toHaveBeenCalledWith('thunderid-wayfinder-setup-expanded', 'false');
  });

  it('respects sessionStorage expanded=true even when already imported', () => {
    mockLocalStorageGetItem.mockReturnValue('1234567890');
    mockSessionStorageGetItem.mockReturnValue('true');
    render(<WayfinderSampleSetup />);
    expect(screen.getByTestId('terminal-block')).toBeInTheDocument();
  });

  it('shows WayfinderSampleDownload when expanded', () => {
    render(<WayfinderSampleSetup />);
    expect(screen.getByTestId('wayfinder-sample-download')).toBeInTheDocument();
  });

  it('shows step titles when expanded', () => {
    render(<WayfinderSampleSetup />);
    expect(screen.getByText('common:welcome.wayfinderSampleSetup.steps.getSample.title')).toBeInTheDocument();
    expect(screen.getByText(/wayfinderSampleSetup.steps.configure.title/)).toBeInTheDocument();
    expect(screen.getByText('common:welcome.wayfinderSampleSetup.steps.run.title')).toBeInTheDocument();
  });

  it('shows unix command by default', () => {
    render(<WayfinderSampleSetup />);
    expect(screen.getByTestId('terminal-block')).toHaveTextContent('./start.sh');
  });

  it('switches to windows command when OS tab is changed', async () => {
    const user = userEvent.setup();
    render(<WayfinderSampleSetup />);

    await user.click(screen.getByRole('tab', {name: 'Windows'}));

    expect(screen.getByTestId('terminal-block')).toHaveTextContent('.\\start.ps1');
  });

  it('auto-selects windows tab when platform is win', async () => {
    mockDetectPlatform.mockResolvedValue({os: 'win', arch: 'x64'});
    render(<WayfinderSampleSetup />);

    await waitFor(() => {
      expect(screen.getByTestId('terminal-block')).toHaveTextContent('.\\start.ps1');
    });
  });

  it('collapses and marks done after import success', async () => {
    const user = userEvent.setup();
    render(<WayfinderSampleSetup />);

    expect(screen.getByTestId('wayfinder-folder-import')).toBeInTheDocument();

    await user.click(screen.getByText('mock-import-success'));

    expect(screen.queryByTestId('wayfinder-folder-import')).not.toBeInTheDocument();
    expect(mockSessionStorageSetItem).toHaveBeenCalledWith('thunderid-wayfinder-setup-expanded', 'false');
  });

  it('toggles expand/collapse when header receives Enter keypress', () => {
    render(<WayfinderSampleSetup />);

    expect(screen.getByTestId('terminal-block')).toBeInTheDocument();

    const header = screen.getByRole('button', {name: /wayfinderSampleSetup\.title/i});
    fireEvent.keyDown(header, {key: 'Enter'});

    expect(screen.queryByTestId('terminal-block')).not.toBeInTheDocument();
  });
});
