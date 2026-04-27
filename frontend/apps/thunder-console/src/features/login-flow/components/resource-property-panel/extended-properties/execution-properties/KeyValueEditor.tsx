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

import {IconButton, Stack, TextField} from '@wso2/oxygen-ui';
import type {ReactNode} from 'react';

interface KeyValueEditorProps {
  entries: [string, string][];
  onAdd: () => void;
  onRemove: (index: number) => void;
  onKeyChange: (index: number, newKey: string) => void;
  onValueChange: (index: number, newValue: string) => void;
  keyPlaceholder: string;
  valuePlaceholder: string;
}

function KeyValueEditor({
  entries,
  onAdd,
  onRemove,
  onKeyChange,
  onValueChange,
  keyPlaceholder,
  valuePlaceholder,
}: KeyValueEditorProps): ReactNode {
  return (
    <Stack gap={1}>
      {entries.map(([key, value], index) => (
        // Using index as React key since entry keys can be duplicated or empty.
        // The list is short and only mutated via add/remove at known positions.
        // eslint-disable-next-line react/no-array-index-key
        <Stack key={index} direction="row" gap={1} alignItems="center">
          <TextField
            value={key}
            onChange={(e) => onKeyChange(index, e.target.value)}
            placeholder={keyPlaceholder}
            size="small"
            sx={{flex: 1}}
          />
          <TextField
            value={value}
            onChange={(e) => onValueChange(index, e.target.value)}
            placeholder={valuePlaceholder}
            size="small"
            sx={{flex: 1}}
          />
          <IconButton size="small" onClick={() => onRemove(index)} aria-label="Remove entry">
            &times;
          </IconButton>
        </Stack>
      ))}
      <IconButton size="small" onClick={onAdd} sx={{alignSelf: 'flex-start'}} aria-label="Add entry">
        +
      </IconButton>
    </Stack>
  );
}

export default KeyValueEditor;
