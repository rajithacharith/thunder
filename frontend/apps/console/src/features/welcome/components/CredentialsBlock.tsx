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

import {Box, IconButton, Stack, Typography} from '@wso2/oxygen-ui';
import {Check, Copy, Eye, EyeOff} from '@wso2/oxygen-ui-icons-react';
import type {JSX} from 'react';
import {useState} from 'react';

interface CredentialRowProps {
  field: 'username' | 'password';
  value: string;
  showPassword: boolean;
  isCopied: boolean;
  onToggleShow: () => void;
  onCopy: () => void;
}

function CredentialRow({field, value, showPassword, isCopied, onToggleShow, onCopy}: CredentialRowProps): JSX.Element {
  const isPassword = field === 'password';
  return (
    <Box
      sx={{
        border: '1px solid',
        borderColor: 'divider',
        borderRadius: 1.5,
        px: 2,
        py: 1,
        display: 'flex',
        alignItems: 'center',
        gap: 1,
      }}
    >
      <Box sx={{flex: 1, minWidth: 0}}>
        <Typography variant="caption" color="text.secondary" sx={{display: 'block', textTransform: 'capitalize'}}>
          {field}
        </Typography>
        <Typography variant="body2" fontFamily="monospace">
          {isPassword && !showPassword ? '••••••••' : value}
        </Typography>
      </Box>
      {isPassword && (
        <IconButton
          size="small"
          aria-label={showPassword ? 'Hide password' : 'Show password'}
          onClick={onToggleShow}
          sx={{color: 'text.secondary'}}
        >
          {showPassword ? <EyeOff size={14} /> : <Eye size={14} />}
        </IconButton>
      )}
      <IconButton
        size="small"
        aria-label={`Copy ${field}`}
        onClick={onCopy}
        sx={{color: isCopied ? 'success.main' : 'text.secondary'}}
      >
        {isCopied ? <Check size={14} /> : <Copy size={14} />}
      </IconButton>
    </Box>
  );
}

export interface CredentialsBlockProps {
  username: string;
  password: string;
}

export default function CredentialsBlock({username, password}: CredentialsBlockProps): JSX.Element {
  const [showPassword, setShowPassword] = useState(false);
  const [copiedField, setCopiedField] = useState<'username' | 'password' | null>(null);

  const handleCopy = (field: 'username' | 'password', value: string): void => {
    void navigator.clipboard.writeText(value).then(() => {
      setCopiedField(field);
      setTimeout(() => setCopiedField(null), 2000);
    });
  };

  return (
    <Stack spacing={1}>
      <CredentialRow
        field="username"
        value={username}
        showPassword={showPassword}
        isCopied={copiedField === 'username'}
        onToggleShow={() => setShowPassword((v) => !v)}
        onCopy={() => handleCopy('username', username)}
      />
      <CredentialRow
        field="password"
        value={password}
        showPassword={showPassword}
        isCopied={copiedField === 'password'}
        onToggleShow={() => setShowPassword((v) => !v)}
        onCopy={() => handleCopy('password', password)}
      />
    </Stack>
  );
}
