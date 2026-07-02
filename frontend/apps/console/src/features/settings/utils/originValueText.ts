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

import isStringOrigin from './isStringOrigin';
import type {AllowedOrigin} from '../models/responses';

/**
 * Returns the editable text for an allowed origin entry.
 *
 * @param entry - The allowed origin entry
 * @returns The literal string, or the pattern of a regex entry
 */
export default function originValueText(entry: AllowedOrigin): string {
  return isStringOrigin(entry) ? entry : entry.regex;
}
