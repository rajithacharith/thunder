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
import StepList from '../StepList';

describe('StepList', () => {
  it('renders each step', () => {
    render(<StepList steps={['First step', 'Second step']} />);
    expect(screen.getByText('First step')).toBeInTheDocument();
    expect(screen.getByText('Second step')).toBeInTheDocument();
  });

  it('numbers steps starting from 1 by default', () => {
    render(<StepList steps={['Step A', 'Step B']} />);
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
  });

  it('numbers steps from a custom startFrom value', () => {
    render(<StepList steps={['Step A', 'Step B']} startFrom={3} />);
    expect(screen.getByText('3')).toBeInTheDocument();
    expect(screen.getByText('4')).toBeInTheDocument();
    expect(screen.queryByText('1')).not.toBeInTheDocument();
  });

  it('renders JSX node steps', () => {
    render(
      <StepList
        steps={[
          <span key="a" data-testid="jsx-step">
            JSX content
          </span>,
        ]}
      />,
    );
    expect(screen.getByTestId('jsx-step')).toBeInTheDocument();
  });
});
