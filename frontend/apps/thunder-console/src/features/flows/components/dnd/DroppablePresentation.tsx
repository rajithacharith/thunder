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

import {Box, type BoxProps} from '@wso2/oxygen-ui';
import {memo, type ReactElement, type ReactNode} from 'react';

/**
 * Props interface for DroppablePresentation
 */
export interface DroppablePresentationProps {
  children: ReactNode;
  className?: string;
  sx?: BoxProps['sx'];
  isDragActive?: boolean;
}

/**
 * Memoized presentation component for Droppable content.
 * PERFORMANCE FIX: Based on dnd-kit issue #389 - separate presentation from hook
 * This prevents children from re-rendering when useDroppable causes parent re-renders.
 * @see https://github.com/clauderic/dnd-kit/issues/389
 *
 * @param props - Props injected to the component.
 * @returns DroppablePresentation component.
 */
function DroppablePresentation({
  children,
  className = undefined,
  sx = {},
  isDragActive = false,
}: DroppablePresentationProps): ReactElement {
  return (
    <Box
      className={className}
      sx={{
        display: 'inline-flex',
        flexDirection: 'column',
        height: '100%',
        width: '100%',
        ...sx,
        // Applied after sx spread so drag-active padding isn't overridden
        // by caller shorthands like `p`.
        paddingTop: isDragActive ? '20px' : '0px',
        paddingBottom: isDragActive ? '20px' : '0px',
        transition: 'padding 0.2s ease',
      }}
    >
      {children}
    </Box>
  );
}

export default memo(DroppablePresentation);
