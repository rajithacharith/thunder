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

/**
 * Route paths this package needs from the host application.
 *
 * The host supplies these via `@thunderid/contexts`'s `RoutesProvider`. When absent (e.g. this
 * package rendered standalone in Storybook or a unit test), `useConnectionRoutes` falls back to
 * `defaultConnectionRoutePaths` below.
 *
 * Includes `trustedIssuers` alongside `connections`: a trusted-issuer instance is a connection
 * under the hood, but it opens a distinct detail page owned by the host application rather than
 * this package's own `ConnectionDetailPage`.
 *
 * @public
 */
export interface ConnectionRoutePaths {
  connections: {
    list: () => string;
    byType: (type: string) => string;
    detail: (type: string, id: string) => string;
    configure: (type: string) => string;
    create: () => string;
  };
  trustedIssuers: {
    detail: (id: string) => string;
  };
}

/**
 * Default connection (and trusted issuer) paths, used when no host-supplied override is present.
 *
 * @public
 */
export const defaultConnectionRoutePaths: ConnectionRoutePaths = {
  connections: {
    list: () => '/connections',
    byType: (type) => `/connections/${type}`,
    detail: (type, id) => `/connections/${type}/${id}`,
    configure: (type) => `/connections/${type}/configure`,
    create: () => '/connections/create',
  },
  trustedIssuers: {
    detail: (id) => `/trusted-issuers/${id}`,
  },
};
