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

import {render, screen, userEvent, waitFor} from '@thunderid/test-utils';
import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';

vi.mock('framer-motion', () => ({
  motion: {
    create: (Component: React.ElementType) => Component,
  },
}));

vi.mock('@wso2/oxygen-ui-icons-react', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@wso2/oxygen-ui-icons-react')>();
  return {
    ...actual,
    CheckCircle: () => <span data-testid="icon-check-circle" />,
    FileCode: () => <span data-testid="icon-file-code" />,
    FileText: () => <span data-testid="icon-file-text" />,
    FolderOpen: () => <span data-testid="icon-folder-open" />,
    RefreshCw: () => <span data-testid="icon-refresh-cw" />,
    XCircle: () => <span data-testid="icon-x-circle" />,
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

const {mockUseImportConfiguration, mockMutateAsync} = vi.hoisted(() => ({
  mockUseImportConfiguration: vi.fn(),
  mockMutateAsync: vi.fn(),
}));

vi.mock('../../../import-export/api/useImportConfiguration', () => ({
  default: mockUseImportConfiguration,
}));

vi.mock('@thunderid/contexts', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@thunderid/contexts')>();
  return {
    ...actual,
    useConfig: () => ({config: {brand: {product_name: 'ThunderID'}}}),
  };
});

function buildDirectoryPickerMock(
  folderName = 'my-project',
  {includeEnv = true}: {includeEnv?: boolean} = {},
): FileSystemDirectoryHandle {
  // Use plain async functions (not vi.fn()) so restoreMocks does not wipe them.
  const yamlMockFile = {
    name: 'thunderid-config.yaml',
    text: (): Promise<string> => Promise.resolve('yaml content'),
  };
  const envMockFile = {
    name: 'thunderid.env',
    text: (): Promise<string> => Promise.resolve('KEY=value'),
  };
  const yamlHandle = {getFile: () => Promise.resolve(yamlMockFile)};
  const envHandle = {getFile: () => Promise.resolve(envMockFile)};

  const configDir = {
    getFileHandle: (name: string) => {
      if (name === 'thunderid-config.yaml') return Promise.resolve(yamlHandle);
      if (name === 'thunderid-config.yml') return Promise.reject(new Error('not found'));
      if (name === 'thunderid.env' && includeEnv) return Promise.resolve(envHandle);
      return Promise.reject(new Error(`not found: ${name}`));
    },
  };

  return {
    name: folderName,
    getDirectoryHandle: (dirName: string) => {
      if (dirName === 'thunderid-config') return Promise.resolve(configDir);
      return Promise.reject(new Error(`not found: ${dirName}`));
    },
  } as unknown as FileSystemDirectoryHandle;
}

import WayfinderFolderImport from '../WayfinderFolderImport';

describe('WayfinderFolderImport', () => {
  const mockLocalStorageGetItem = vi.fn();
  const mockLocalStorageSetItem = vi.fn();

  beforeEach(() => {
    mockLocalStorageGetItem.mockReturnValue(null);
    vi.stubGlobal('localStorage', {
      getItem: mockLocalStorageGetItem,
      setItem: mockLocalStorageSetItem,
      removeItem: vi.fn(),
      clear: vi.fn(),
    });
    vi.stubGlobal('showDirectoryPicker', vi.fn().mockResolvedValue(buildDirectoryPickerMock()));
    mockUseImportConfiguration.mockReturnValue({mutateAsync: mockMutateAsync});
  });

  afterEach(() => {
    vi.clearAllMocks();
    vi.unstubAllGlobals();
  });

  it('renders without crashing', () => {
    const {container} = render(<WayfinderFolderImport />);
    expect(container).toBeInTheDocument();
  });

  it('shows "Select Wayfinder Sample Folder" button when idle', () => {
    render(<WayfinderFolderImport />);
    expect(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder')).toBeInTheDocument();
  });

  it('shows alreadyDone state when localStorage has the key', () => {
    mockLocalStorageGetItem.mockReturnValue(Date.now().toString());
    render(<WayfinderFolderImport />);
    expect(
      screen.getByText('common:welcome.wayfinderFolderImport.status.alreadyDone:{"productName":"ThunderID"}'),
    ).toBeInTheDocument();
    expect(screen.getByText('common:welcome.wayfinderFolderImport.actions.reImport')).toBeInTheDocument();
  });

  it('re-import button in alreadyDone state resets to idle', async () => {
    mockLocalStorageGetItem.mockReturnValue(Date.now().toString());
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.reImport'));

    expect(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder')).toBeInTheDocument();
  });

  it('stays idle when directory picker is aborted', async () => {
    const abortError = new DOMException('User aborted', 'AbortError');
    vi.stubGlobal('showDirectoryPicker', vi.fn().mockRejectedValue(abortError));
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));

    expect(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder')).toBeInTheDocument();
    expect(screen.queryByTestId('icon-x-circle')).not.toBeInTheDocument();
  });

  it('shows error state when picker fails with non-abort error', async () => {
    vi.stubGlobal('showDirectoryPicker', vi.fn().mockRejectedValue(new Error('permission denied')));
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));

    expect(await screen.findByText('common:welcome.wayfinderFolderImport.errors.cannotReadFolder')).toBeInTheDocument();
  });

  it('shows folder name and file names after folder is selected', async () => {
    vi.stubGlobal('showDirectoryPicker', vi.fn().mockResolvedValue(buildDirectoryPickerMock('my-project')));
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));

    expect(await screen.findByText('my-project/')).toBeInTheDocument();
    expect(screen.getByText(/thunderid-config\/thunderid-config\.yaml/)).toBeInTheDocument();
    expect(screen.getByText(/thunderid-config\/thunderid\.env/)).toBeInTheDocument();
  });

  it('shows "Configure in ThunderID" button after folder is selected', async () => {
    vi.stubGlobal('showDirectoryPicker', vi.fn().mockResolvedValue(buildDirectoryPickerMock()));
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));

    expect(await screen.findByText('common:welcome.wayfinderFolderImport.actions.importConfig')).toBeInTheDocument();
  });

  it('clicking import calls mutateAsync', async () => {
    mockMutateAsync.mockResolvedValue({
      summary: {totalDocuments: 2, imported: 2, failed: 0, importedAt: ''},
      results: [],
    });
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));
    await user.click(await screen.findByText('common:welcome.wayfinderFolderImport.actions.importConfig'));

    await waitFor(() => {
      expect(mockMutateAsync).toHaveBeenCalledWith(
        expect.objectContaining({
          content: 'yaml content',
          options: {upsert: true},
        }),
      );
    });
  });

  it('shows success message on import success', async () => {
    mockMutateAsync.mockResolvedValue({
      summary: {totalDocuments: 2, imported: 2, failed: 0, importedAt: ''},
      results: [],
    });
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));
    await user.click(await screen.findByText('common:welcome.wayfinderFolderImport.actions.importConfig'));

    expect(
      await screen.findByText('common:welcome.wayfinderFolderImport.status.success:{"productName":"ThunderID"}'),
    ).toBeInTheDocument();
  });

  it('shows error with count on partial failure', async () => {
    mockMutateAsync.mockResolvedValue({
      summary: {totalDocuments: 3, imported: 2, failed: 1, importedAt: ''},
      results: [{resourceType: 'application', status: 'failed', message: 'conflict'}],
    });
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));
    await user.click(await screen.findByText('common:welcome.wayfinderFolderImport.actions.importConfig'));

    expect(
      await screen.findByText('common:welcome.wayfinderFolderImport.errors.partialFailure:{"count":1}'),
    ).toBeInTheDocument();
  });

  it('shows error message on import error', async () => {
    mockMutateAsync.mockRejectedValue(new Error('server error'));
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));
    await user.click(await screen.findByText('common:welcome.wayfinderFolderImport.actions.importConfig'));

    expect(await screen.findByText('common:welcome.wayfinderFolderImport.errors.importFailed')).toBeInTheDocument();
  });

  it('shows error when yaml file is not found in selected folder', async () => {
    const noYamlDir = {
      getFileHandle: () => Promise.reject(new Error('not found')),
    };
    vi.stubGlobal(
      'showDirectoryPicker',
      vi.fn().mockResolvedValue({
        name: 'no-yaml-project',
        getDirectoryHandle: (dirName: string) => {
          if (dirName === 'thunderid-config') return Promise.resolve(noYamlDir);
          return Promise.reject(new Error('not found'));
        },
      }),
    );
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));

    expect(await screen.findByText('common:welcome.wayfinderFolderImport.errors.cannotReadFolder')).toBeInTheDocument();
  });

  it('change button in found state re-opens folder picker', async () => {
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));
    await screen.findByText('common:welcome.wayfinderFolderImport.actions.importConfig');

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.change'));

    expect(
      vi.mocked(
        (window as unknown as {showDirectoryPicker: () => Promise<FileSystemDirectoryHandle>}).showDirectoryPicker,
      ),
    ).toHaveBeenCalledTimes(2);
  });

  it('re-select folder button in error state resets and re-opens picker', async () => {
    mockMutateAsync.mockRejectedValue(new Error('server error'));
    const user = userEvent.setup();
    render(<WayfinderFolderImport />);

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.selectFolder'));
    await user.click(await screen.findByText('common:welcome.wayfinderFolderImport.actions.importConfig'));
    await screen.findByText('common:welcome.wayfinderFolderImport.errors.importFailed');

    await user.click(screen.getByText('common:welcome.wayfinderFolderImport.actions.reSelectFolder'));

    expect(
      vi.mocked(
        (window as unknown as {showDirectoryPicker: () => Promise<FileSystemDirectoryHandle>}).showDirectoryPicker,
      ),
    ).toHaveBeenCalledTimes(2);
  });
});
