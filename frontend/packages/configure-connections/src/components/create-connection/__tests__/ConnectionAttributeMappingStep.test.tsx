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

import {fireEvent, render, screen} from '@thunderid/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import ConnectionAttributeMappingStep from '../ConnectionAttributeMappingStep';

// Stub the mapping control so the step can drive its onChange without pulling in user-type hooks.
vi.mock('../../AttributeMappingSection', () => ({
  default: ({onChange}: {onChange: (config: unknown, valid: boolean) => void}) => (
    <button
      type="button"
      data-testid="stub-add-mapping"
      onClick={() =>
        onChange(
          {
            userTypeResolution: {default: 'Person'},
            userTypeAttributeMappings: [
              {userType: 'Person', attributes: [{externalAttribute: 'a', localAttribute: 'b'}]},
            ],
          },
          true,
        )
      }
    >
      add
    </button>
  ),
}));

describe('ConnectionAttributeMappingStep', () => {
  const props = {
    vendorDisplayName: 'Google',
    onChange: vi.fn(),
    onBack: vi.fn(),
    onCreate: vi.fn(),
    isPending: false,
    createDisabled: false,
  };
  beforeEach(() => vi.clearAllMocks());

  it('shows the skip-and-create label until a mapping is added', () => {
    render(<ConnectionAttributeMappingStep {...props} />);
    expect(screen.getByText('Map provider attributes to your users')).toBeInTheDocument();
    // The stepper's short step label (`wizard.steps.attributeMapping`) must not be rendered here;
    // assert on the raw key since its real translation collides with the section title text.
    expect(screen.queryByText('wizard.steps.attributeMapping')).not.toBeInTheDocument();
    expect(screen.getByTestId('wizard-create')).toHaveTextContent('Skip and Create connection');

    fireEvent.click(screen.getByTestId('stub-add-mapping'));
    expect(screen.getByTestId('wizard-create')).toHaveTextContent('Create connection');
  });

  it('invokes onCreate and onBack from the button row', () => {
    render(<ConnectionAttributeMappingStep {...props} />);
    fireEvent.click(screen.getByTestId('wizard-create'));
    expect(props.onCreate).toHaveBeenCalledTimes(1);

    fireEvent.click(screen.getByText('Back'));
    expect(props.onBack).toHaveBeenCalledTimes(1);
  });
});
