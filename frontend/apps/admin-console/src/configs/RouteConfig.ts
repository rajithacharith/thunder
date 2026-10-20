// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

import type {AgentTypeRoutePaths} from '@thunderid/configure-agent-types';
import type {AgentRoutePaths} from '@thunderid/configure-agents';
import type {GroupRoutePaths} from '@thunderid/configure-groups';
import type {RoleRoutePaths} from '@thunderid/configure-roles';
import type {UserTypeRoutePaths} from '@thunderid/configure-user-types';
import type {UserRoutePaths} from '@thunderid/configure-users';

/**
 * Route paths the Admin Portal implements itself, rather than through a `@thunderid/configure-*`
 * package. Mirrors Console's own `ConsoleRoutePaths` for the same reason: nothing outside this app
 * depends on this shape, but centralizing it still keeps `App.tsx` and every navigate() call built
 * from the same source.
 *
 * @public
 */
export interface AdminPortalRoutePaths {
  dashboard: {
    list: () => string;
  };
}

/**
 * Route configuration for the Admin Portal app: the dashboard it implements itself, plus the two
 * domains it exposes via `@thunderid/configure-*` packages, composed from the route-shape each
 * package publishes.
 *
 * Each package reads its own slice through a package-level hook built on
 * `@thunderid/contexts`'s `useRoutes`, which lets a host application override any of those paths
 * without touching the package.
 *
 * @public
 */
export type RouteConfig = UserRoutePaths &
  UserTypeRoutePaths &
  AgentTypeRoutePaths &
  AgentRoutePaths &
  GroupRoutePaths &
  RoleRoutePaths &
  AdminPortalRoutePaths;

/**
 * Static path segment for each domain, with no leading slash.
 *
 * `App.tsx`'s `<Route path>` declarations are built from these same segments, so the mounted
 * routes and the destinations packages navigate to via `RouteConfig` below can never drift apart.
 *
 * @public
 */
export const ROUTE_SEGMENTS = {
  dashboard: 'dashboard',
  users: 'users',
  userTypes: 'user-types',
  agentTypes: 'agent-types',
  agents: 'agents',
  groups: 'groups',
  roles: 'roles',
} as const;

const RouteConfig: RouteConfig = {
  dashboard: {
    list: () => `/${ROUTE_SEGMENTS.dashboard}`,
  },
  users: {
    list: () => `/${ROUTE_SEGMENTS.users}`,
    detail: (userId) => `/${ROUTE_SEGMENTS.users}/${userId}`,
    add: () => `/${ROUTE_SEGMENTS.users}/add`,
  },
  userTypes: {
    list: () => `/${ROUTE_SEGMENTS.userTypes}`,
    detail: (id) => `/${ROUTE_SEGMENTS.userTypes}/${id}`,
    create: () => `/${ROUTE_SEGMENTS.userTypes}/create`,
  },
  agentTypes: {
    detail: (id) => `/${ROUTE_SEGMENTS.agentTypes}/${id}`,
  },
  agents: {
    list: () => `/${ROUTE_SEGMENTS.agents}`,
    detail: (id) => `/${ROUTE_SEGMENTS.agents}/${id}`,
    create: () => `/${ROUTE_SEGMENTS.agents}/create`,
  },
  groups: {
    list: () => `/${ROUTE_SEGMENTS.groups}`,
    detail: (id) => `/${ROUTE_SEGMENTS.groups}/${id}`,
    create: () => `/${ROUTE_SEGMENTS.groups}/create`,
  },
  roles: {
    list: () => `/${ROUTE_SEGMENTS.roles}`,
    detail: (id) => `/${ROUTE_SEGMENTS.roles}/${id}`,
    create: () => `/${ROUTE_SEGMENTS.roles}/create`,
  },
  // The agent pages link to an onboarding flow this app does not have, so both destinations are
  // mapped to the nearest equivalents here.
  welcome: {
    getStarted: () => `/${ROUTE_SEGMENTS.dashboard}`,
    getStartedAgentsCreate: () => `/${ROUTE_SEGMENTS.agents}/create`,
  },
};

export default RouteConfig;
