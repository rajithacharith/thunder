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

import {render, screen} from '@thunderid/test-utils';
import {describe, expect, it} from 'vitest';
import CodeInline from '../CodeInline';

describe('CodeInline', () => {
  it('renders children', () => {
    render(<CodeInline>my-token</CodeInline>);
    expect(screen.getByText('my-token')).toBeInTheDocument();
  });

  it('renders as a code element', () => {
    const {container} = render(<CodeInline>value</CodeInline>);
    expect(container.querySelector('code')).toBeInTheDocument();
  });

  it('renders without children', () => {
    const {container} = render(<CodeInline />);
    expect(container).toBeInTheDocument();
  });
});
