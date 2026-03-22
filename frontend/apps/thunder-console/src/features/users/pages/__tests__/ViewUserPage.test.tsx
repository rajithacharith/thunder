/**
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
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

import {beforeEach, describe, expect, it, vi} from 'vitest';
import {render, screen, userEvent, waitFor} from '@thunder/test-utils';
import ViewUserPage from '../ViewUserPage';
import type {ApiError, ApiUser, ApiUserSchema, UserSchemaListResponse} from '../../types/users';

const mockNavigate = vi.fn();
const mockMutateAsync = vi.fn();
const mockUpdateReset = vi.fn();
const mockDeleteMutate = vi.fn();

vi.mock('react-router', async () => {
  const actual = await vi.importActual<typeof import('react-router')>('react-router');

  return {
    ...actual,
    useNavigate: () => mockNavigate,
    useParams: () => ({userId: 'user123'}),
  };
});

interface QueryResult<TData> {
  data: TData;
  isLoading: boolean;
  error: ApiError | null;
}

interface UpdateMutationResult {
  mutateAsync: typeof mockMutateAsync;
  isPending: boolean;
  error: ApiError | null;
  reset: typeof mockUpdateReset;
}

interface DeleteMutationResult {
  mutate: typeof mockDeleteMutate;
  isPending: boolean;
}

const mockUseGetUser = vi.fn<() => QueryResult<ApiUser | undefined>>();
const mockUseGetUserSchemas = vi.fn<() => QueryResult<UserSchemaListResponse | undefined>>();
const mockUseGetUserSchema = vi.fn<() => QueryResult<ApiUserSchema | undefined>>();
const mockUseUpdateUser = vi.fn<() => UpdateMutationResult>();
const mockUseDeleteUser = vi.fn<() => DeleteMutationResult>();

vi.mock('../../api/useGetUser', () => ({
  default: () => mockUseGetUser(),
}));

vi.mock('../../api/useGetUserSchemas', () => ({
  default: () => mockUseGetUserSchemas(),
}));

vi.mock('../../api/useGetUserSchema', () => ({
  default: () => mockUseGetUserSchema(),
}));

vi.mock('../../api/useUpdateUser', () => ({
  default: () => mockUseUpdateUser(),
}));

vi.mock('../../api/useDeleteUser', () => ({
  default: () => mockUseDeleteUser(),
}));

describe('ViewUserPage', () => {
  const user: ApiUser = {
    id: 'user123',
    display: 'John Doe',
    type: 'Employee',
    attributes: {
      username: 'john_doe',
      email: 'john@example.com',
    },
  };

  const schemas: UserSchemaListResponse = {
    totalResults: 1,
    startIndex: 1,
    count: 1,
    schemas: [{id: 'schema-employee', name: 'Employee', ouId: 'ou-1'}],
  };

  const schema: ApiUserSchema = {
    id: 'schema-employee',
    name: 'Employee',
    ouId: 'ou-1',
    schema: {
      username: {
        type: 'string',
        required: true,
      },
      email: {
        type: 'string',
        required: true,
      },
    },
  };

  beforeEach(() => {
    vi.clearAllMocks();

    mockUseGetUser.mockReturnValue({
      data: user,
      isLoading: false,
      error: null,
    });

    mockUseGetUserSchemas.mockReturnValue({
      data: schemas,
      isLoading: false,
      error: null,
    });

    mockUseGetUserSchema.mockReturnValue({
      data: schema,
      isLoading: false,
      error: null,
    });

    mockUseUpdateUser.mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: false,
      error: null,
      reset: mockUpdateReset,
    });

    mockUseDeleteUser.mockReturnValue({
      mutate: mockDeleteMutate,
      isPending: false,
    });

    mockNavigate.mockResolvedValue(undefined);
    mockMutateAsync.mockResolvedValue(user);
  });

  it('shows loading indicator while user is loading', () => {
    mockUseGetUser.mockReturnValue({
      data: undefined,
      isLoading: true,
      error: null,
    });

    render(<ViewUserPage />);

    expect(screen.getByRole('progressbar')).toBeInTheDocument();
  });

  it('shows error alert and back action when user query fails', () => {
    mockUseGetUser.mockReturnValue({
      data: undefined,
      isLoading: false,
      error: {
        code: 'USR-1000',
        message: 'Failed to load user',
        description: 'User API failed',
      },
    });

    render(<ViewUserPage />);

    expect(screen.getByRole('alert')).toHaveTextContent('Failed to load user');
    expect(screen.getByRole('button', {name: /back|users:manageuser.back/i})).toBeInTheDocument();
  });

  it('renders user header and attribute values', async () => {
    render(<ViewUserPage />);

    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('Employee')).toBeInTheDocument();
    expect(screen.getByText('john_doe')).toBeInTheDocument();
    expect(screen.getByText('john@example.com')).toBeInTheDocument();

    const copyControl = await screen.findByRole('button', {name: /copy user id|copy user/i});
    expect(copyControl).toBeInTheDocument();
  });

  it('submits updated values through mutateAsync', async () => {
    const userEventInstance = userEvent.setup();

    render(<ViewUserPage />);

    await userEventInstance.click(screen.getByRole('button', {name: /edit|common:actions.edit/i}));
    await userEventInstance.click(screen.getByRole('button', {name: /save changes|save|common:actions.save/i}));

    await waitFor(() => {
      expect(mockMutateAsync).toHaveBeenCalledTimes(1);
      expect(mockMutateAsync).toHaveBeenCalledWith({
        userId: 'user123',
        data: {
          ouId: 'ou-1',
          type: 'Employee',
          attributes: {
            username: 'john_doe',
            email: 'john@example.com',
          },
        },
      });
    });
  });

  it('resets edit state on cancel', async () => {
    const userEventInstance = userEvent.setup();

    render(<ViewUserPage />);

    await userEventInstance.click(screen.getByRole('button', {name: /edit|common:actions.edit/i}));
    await userEventInstance.click(screen.getByRole('button', {name: /cancel|common:actions.cancel/i}));

    expect(mockUpdateReset).toHaveBeenCalledTimes(1);
    expect(screen.getByRole('button', {name: /edit|common:actions.edit/i})).toBeInTheDocument();
  });

  it('opens delete dialog and deletes user successfully', async () => {
    const userEventInstance = userEvent.setup();

    mockDeleteMutate.mockImplementation(
      (_id: string, options?: {onSuccess?: () => void; onError?: (error: Error) => void}) => {
        options?.onSuccess?.();
      },
    );

    render(<ViewUserPage />);

    await userEventInstance.click(screen.getByRole('button', {name: /^delete$|common:actions.delete/i}));
    await userEventInstance.click(screen.getByRole('button', {name: /^delete$|common:actions.delete/i}));

    await waitFor(() => {
      expect(mockDeleteMutate).toHaveBeenCalledTimes(1);

      const [mutateUserId, callbacks] = mockDeleteMutate.mock.calls[0] as [
        string,
        {onSuccess?: () => void; onError?: (error: Error) => void} | undefined,
      ];

      expect(mutateUserId).toBe('user123');
      expect(typeof callbacks?.onSuccess).toBe('function');
      expect(typeof callbacks?.onError).toBe('function');
      expect(mockNavigate).toHaveBeenCalledWith('/users');
    });
  });

  it('shows delete error when mutation fails', async () => {
    const userEventInstance = userEvent.setup();

    mockDeleteMutate.mockImplementation(
      (_id: string, options?: {onSuccess?: () => void; onError?: (error: Error) => void}) => {
        options?.onError?.(new Error('Delete failed'));
      },
    );

    render(<ViewUserPage />);

    await userEventInstance.click(screen.getByRole('button', {name: /^delete$|common:actions.delete/i}));
    await userEventInstance.click(screen.getByRole('button', {name: /^delete$|common:actions.delete/i}));

    expect(screen.getByText('Delete failed')).toBeInTheDocument();
  });
});
