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

import {describe, expect, it} from 'vitest';
import getConfigurationTypeFromTemplate from '../getConfigurationTypeFromTemplate';
import {ApplicationCreateFlowConfiguration} from '../../models/application-create-flow';
import type {ApplicationTemplate} from '../../models/application-templates';

describe('getConfigurationTypeFromTemplate', () => {
  it('returns NONE for null template config', () => {
    const result = getConfigurationTypeFromTemplate(null);
    expect(result).toBe(ApplicationCreateFlowConfiguration.NONE);
  });

  it('returns NONE when redirectUris is already populated', () => {
    const template: ApplicationTemplate = {
      name: 'Test App',
      description: 'Test application',
      inboundAuthConfig: [
        {
          type: 'oauth2',
          config: {
            redirectUris: ['https://example.com/callback'],
            grantTypes: ['authorization_code'],
            responseTypes: ['code'],
            publicClient: true,
          },
        },
      ],
    };

    const result = getConfigurationTypeFromTemplate(template);
    expect(result).toBe(ApplicationCreateFlowConfiguration.NONE);
  });

  it('returns DEEPLINK for mobile applications', () => {
    const template: ApplicationTemplate = {
      name: 'Mobile Application',
      description: 'Mobile app',
      inboundAuthConfig: [
        {
          type: 'oauth2',
          config: {
            redirectUris: [],
            grantTypes: ['authorization_code'],
            responseTypes: ['code'],
            publicClient: true,
          },
        },
      ],
    };

    const result = getConfigurationTypeFromTemplate(template);
    expect(result).toBe(ApplicationCreateFlowConfiguration.DEEPLINK);
  });

  it('returns URL for browser applications', () => {
    const template: ApplicationTemplate = {
      name: 'Browser Application',
      description: 'Browser app',
      inboundAuthConfig: [
        {
          type: 'oauth2',
          config: {
            redirectUris: [],
            grantTypes: ['authorization_code'],
            responseTypes: ['code'],
            publicClient: true,
          },
        },
      ],
    };

    const result = getConfigurationTypeFromTemplate(template);
    expect(result).toBe(ApplicationCreateFlowConfiguration.URL);
  });

  it('returns URL for server applications', () => {
    const template: ApplicationTemplate = {
      name: 'Server Application',
      description: 'Server app',
      inboundAuthConfig: [
        {
          type: 'oauth2',
          config: {
            redirectUris: [],
            grantTypes: ['authorization_code'],
            responseTypes: ['code'],
            publicClient: false,
          },
        },
      ],
    };

    const result = getConfigurationTypeFromTemplate(template);
    expect(result).toBe(ApplicationCreateFlowConfiguration.URL);
  });

  it('returns NONE for backend applications', () => {
    const template: ApplicationTemplate = {
      name: 'Backend Application',
      description: 'Backend service',
      inboundAuthConfig: [
        {
          type: 'oauth2',
          config: {
            redirectUris: [],
            grantTypes: ['client_credentials'],
            responseTypes: [],
            publicClient: false,
          },
        },
      ],
    };

    const result = getConfigurationTypeFromTemplate(template);
    expect(result).toBe(ApplicationCreateFlowConfiguration.NONE);
  });

  it('returns URL as default for unknown application types', () => {
    const template: ApplicationTemplate = {
      name: 'Unknown Application',
      description: 'Unknown app type',
      inboundAuthConfig: [
        {
          type: 'oauth2',
          config: {
            redirectUris: [],
            grantTypes: ['authorization_code'],
            responseTypes: ['code'],
            publicClient: true,
          },
        },
      ],
    };

    const result = getConfigurationTypeFromTemplate(template);
    expect(result).toBe(ApplicationCreateFlowConfiguration.URL);
  });

  it('handles template with no OAuth config', () => {
    const template: ApplicationTemplate = {
      name: 'No OAuth App',
      description: 'App without OAuth',
      inboundAuthConfig: [],
    };

    const result = getConfigurationTypeFromTemplate(template);
    expect(result).toBe(ApplicationCreateFlowConfiguration.URL);
  });

  it('handles template with undefined redirectUris', () => {
    const template: ApplicationTemplate = {
      name: 'Undefined URIs App',
      description: 'App with undefined redirectUris',
      inboundAuthConfig: [
        {
          type: 'oauth2',
          config: {
            // redirectUris not specified
            grantTypes: ['authorization_code'],
            responseTypes: ['code'],
            publicClient: true,
          },
        },
      ],
    };

    const result = getConfigurationTypeFromTemplate(template);
    expect(result).toBe(ApplicationCreateFlowConfiguration.URL);
  });

  it('handles case-insensitive template names', () => {
    const template: ApplicationTemplate = {
      name: 'MOBILE APPLICATION',
      description: 'Mobile app with uppercase name',
      inboundAuthConfig: [
        {
          type: 'oauth2',
          config: {
            redirectUris: [],
            grantTypes: ['authorization_code'],
            responseTypes: ['code'],
            publicClient: true,
          },
        },
      ],
    };

    const result = getConfigurationTypeFromTemplate(template);
    expect(result).toBe(ApplicationCreateFlowConfiguration.DEEPLINK);
  });
});
