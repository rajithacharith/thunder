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

import {Link} from '@wso2/oxygen-ui';
import {ExternalLink as ExternalLinkIcon} from '@wso2/oxygen-ui-icons-react';
import type {JSX, ReactNode} from 'react';

// TODO: Move this to oxygen-ui and use.
export default function ExternalLink({href, children = null}: {href: string; children?: ReactNode}): JSX.Element {
  return (
    <Link
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      style={{color: 'inherit', fontWeight: 600, display: 'inline-flex', alignItems: 'center', gap: 2}}
    >
      {children}
      <ExternalLinkIcon size={12} style={{flexShrink: 0, opacity: 0.7}} />
    </Link>
  );
}
