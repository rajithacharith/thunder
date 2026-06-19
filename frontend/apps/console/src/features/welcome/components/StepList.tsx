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

import {Box, Stack, Typography} from '@wso2/oxygen-ui';
import type {JSX, ReactNode} from 'react';

export interface StepListProps {
  steps: ReactNode[];
  startFrom?: number;
}

export default function StepList({steps, startFrom = 1}: StepListProps): JSX.Element {
  return (
    <Stack spacing={1} component="ol" sx={{pl: 0, m: 0, listStyle: 'none'}}>
      {steps.map((step, i) => (
        <Stack
          // eslint-disable-next-line react/no-array-index-key
          key={i}
          component="li"
          direction="row"
          spacing={1.5}
          alignItems="flex-start"
        >
          <Box
            sx={{
              width: 20,
              height: 20,
              borderRadius: '50%',
              bgcolor: 'action.selected',
              color: 'text.secondary',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: '0.7rem',
              fontWeight: 700,
              flexShrink: 0,
              mt: 0.15,
            }}
          >
            {startFrom + i}
          </Box>
          <Typography variant="body2" color="text.secondary" sx={{flex: 1}}>
            {step}
          </Typography>
        </Stack>
      ))}
    </Stack>
  );
}
