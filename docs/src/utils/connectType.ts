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

export const CONNECT_TYPE_STORAGE_KEY = 'thunder-connect-type';

const VALID_TYPES = new Set(['app', 'agent', 'mcp']);

export function toConnectType(raw: string | null): 'app' | 'agent' | 'mcp' {
  return raw !== null && VALID_TYPES.has(raw) ? (raw as 'app' | 'agent' | 'mcp') : 'app';
}

export function applyConnectType(type: 'app' | 'agent' | 'mcp'): void {
  localStorage.setItem(CONNECT_TYPE_STORAGE_KEY, type);
  document.documentElement.dataset.connectType = type;
}
