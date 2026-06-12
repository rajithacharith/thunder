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

import Link from '@docusaurus/Link';
import {Box} from '@wso2/oxygen-ui';
import React, {JSX} from 'react';

export interface InfographicStepProps {
  icon: string;
  title: string;
  href?: string;
  subtitle?: string;
  children?: React.ReactNode;
}

export const InfographicStep: React.FC<InfographicStepProps> = () => null;
InfographicStep.displayName = 'InfographicStep';

export function InfographicTimeline({children}: {children: React.ReactNode}): JSX.Element {
  const steps = React.Children.toArray(children).filter(
    (c): c is React.ReactElement<InfographicStepProps> =>
      React.isValidElement(c) && (c.type as React.FC).displayName === 'InfographicStep',
  );

  const darkMode = `html[data-theme='dark'] &`;

  return (
    <Box sx={{display: 'flex', flexDirection: 'column', gap: '2.5rem', margin: '2rem 0', position: 'relative'}}>
      {steps.map((step) => {
        const isLast = steps.indexOf(step) === steps.length - 1;
        return (
          <Box key={step.props.title} sx={{alignItems: 'flex-start', display: 'flex', position: 'relative'}}>
            {/* Marker: icon circle + connecting line */}
            <Box
              sx={{
                alignItems: 'center',
                display: 'flex',
                flexDirection: 'column',
                marginRight: '1.5rem',
                minWidth: '2.5rem',
                position: 'relative',
              }}
            >
              <Box
                sx={{
                  alignItems: 'center',
                  background: 'var(--ifm-color-primary, #0050b3)',
                  borderRadius: '50%',
                  boxShadow: '0 2px 8px rgba(0,0,0,0.08)',
                  color: '#fff',
                  display: 'flex',
                  fontSize: '1.3rem',
                  fontWeight: 'bold',
                  height: '2.5rem',
                  justifyContent: 'center',
                  width: '2.5rem',
                  zIndex: 2,
                  [darkMode]: {
                    background: 'var(--ifm-color-primary-dark, #7dc4fa)',
                    color: '#18191a',
                  },
                }}
              >
                {step.props.icon}
              </Box>
              {!isLast && (
                <Box
                  sx={{
                    background: 'var(--ifm-color-primary-light, #b3e0ff)',
                    borderRadius: '2px',
                    flex: '1 1 auto',
                    marginBottom: '-0.25rem',
                    marginTop: '0.25rem',
                    minHeight: '2.5rem',
                    width: '4px',
                    zIndex: 1,
                    [darkMode]: {background: 'var(--ifm-color-primary, #0050b3)'},
                  }}
                />
              )}
            </Box>

            {/* Content card */}
            <Box
              sx={{
                background: 'var(--ifm-card-background, #fff)',
                borderRadius: '12px',
                boxShadow: 'var(--ifm-global-shadow-lw, 0 2px 8px rgba(0,0,0,0.06))',
                display: 'flex',
                flex: 1,
                flexDirection: 'column',
                gap: '0.5rem',
                padding: '1.5rem 1.25rem',
                [darkMode]: {
                  background: 'var(--ifm-background-color, #18191a)',
                  border: '1px solid var(--ifm-color-emphasis-200, #23272f)',
                  boxShadow: 'var(--ifm-global-shadow-lw, 0 2px 8px rgba(0,0,0,0.18))',
                },
              }}
            >
              <Box
                component={step.props.href ? Link : 'span'}
                {...(step.props.href ? {to: step.props.href} : {})}
                sx={{
                  color: 'var(--ifm-color-primary, #0050b3)',
                  fontSize: '1.2rem',
                  fontWeight: 700,
                  marginBottom: '0.15rem',
                  textDecoration: 'none',
                  '&:hover': {textDecoration: 'underline'},
                  [darkMode]: {color: 'var(--ifm-color-primary-dark, #7dc4fa)'},
                }}
              >
                {step.props.title}
              </Box>
              {step.props.subtitle && (
                <Box
                  sx={{
                    color: 'var(--ifm-color-primary-dark, #003366)',
                    fontSize: '1rem',
                    marginBottom: '0.25rem',
                    [darkMode]: {color: 'var(--ifm-color-primary-light, #b3e0ff)'},
                  }}
                >
                  {step.props.subtitle}
                </Box>
              )}
              {step.props.children && (
                <Box
                  sx={{
                    color: 'var(--ifm-font-color-base, #222)',
                    fontSize: '1.05rem',
                    lineHeight: 1.5,
                    [darkMode]: {color: 'var(--ifm-font-color-base, #e6eaf3)'},
                  }}
                >
                  {step.props.children}
                </Box>
              )}
            </Box>
          </Box>
        );
      })}
    </Box>
  );
}
