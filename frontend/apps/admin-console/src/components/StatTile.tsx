// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

import {Box, Card, Typography} from '@wso2/oxygen-ui';
import type {JSX, ReactNode} from 'react';

export interface StatTileProps {
  icon: ReactNode;
  label: string;
  value: string | number;
  /** A short trailing note, e.g. "+12 this week" — omitted when there's nothing to add. */
  note?: string;
}

/**
 * A single at-a-glance number on the dashboard: an icon, a label, and the value, matching the
 * icon-in-a-soft-box card idiom Console's own home cards already use
 * (see apps/console/src/components/home/cards/HomeNextStepCard.tsx).
 */
export default function StatTile({icon, label, value, note = undefined}: StatTileProps): JSX.Element {
  return (
    <Card variant="outlined" sx={{p: 3, display: 'flex', flexDirection: 'column', gap: 1.5}}>
      <Box
        sx={{
          width: 40,
          height: 40,
          borderRadius: 1,
          bgcolor: 'action.hover',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: 'text.secondary',
        }}
      >
        {icon}
      </Box>
      <Typography variant="h3" fontWeight={600} sx={{fontSize: '1.9rem'}}>
        {value}
      </Typography>
      <Box>
        <Typography variant="body2" color="text.secondary">
          {label}
        </Typography>
        {note && (
          <Typography variant="caption" color="text.secondary" sx={{display: 'block', mt: 0.25}}>
            {note}
          </Typography>
        )}
      </Box>
    </Card>
  );
}
