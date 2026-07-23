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

import {useRoutes} from '@thunderid/contexts';
import {defaultOrganizationUnitRoutePaths, type OrganizationUnitRoutePaths} from './types';

/**
 * Resolves the organization unit route paths, preferring the host application's configuration
 * (supplied via `RoutesProvider`) and falling back to this package's own defaults.
 *
 * Components should never hardcode organization unit destination paths; they should call this
 * hook and build the destination from the returned functions instead.
 *
 * @public
 */
export default function useOrganizationUnitRoutes(): OrganizationUnitRoutePaths['organizationUnits'] {
  const routes = useRoutes<Partial<OrganizationUnitRoutePaths>>();
  return routes.organizationUnits ?? defaultOrganizationUnitRoutePaths.organizationUnits;
}
