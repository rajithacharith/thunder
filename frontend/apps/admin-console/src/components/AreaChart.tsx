// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

import {Box, Typography, useTheme} from '@wso2/oxygen-ui';
import {useId, type JSX} from 'react';

export interface AreaChartPoint {
  label: string;
  value: number;
}

export interface AreaChartProps {
  points: AreaChartPoint[];
  /** Accessible description of what the chart shows. */
  ariaLabel: string;
  height?: number;
}

const WIDTH = 600;
const PADDING = {top: 12, right: 12, bottom: 4, left: 12};

/**
 * A small dependency-free line-and-area chart drawn in SVG, filled from the theme's primary colour.
 * The x axis labels sit below the plot as plain text so they stay readable at any width.
 */
export default function AreaChart({points, ariaLabel, height = 180}: AreaChartProps): JSX.Element {
  const theme = useTheme();
  const gradientId = useId();
  const color = theme.vars?.palette.primary.main ?? theme.palette.primary.main;

  const plotWidth = WIDTH - PADDING.left - PADDING.right;
  const plotHeight = height - PADDING.top - PADDING.bottom;
  const max = Math.max(...points.map((p) => p.value), 1);
  const step = points.length > 1 ? plotWidth / (points.length - 1) : 0;

  const coords = points.map((p, i) => ({
    x: PADDING.left + i * step,
    y: PADDING.top + plotHeight - (p.value / max) * plotHeight,
  }));
  const line = coords.map((c, i) => `${i === 0 ? 'M' : 'L'}${c.x.toFixed(1)},${c.y.toFixed(1)}`).join(' ');
  const baseline = PADDING.top + plotHeight;
  const area =
    coords.length > 0
      ? `${line} L${coords[coords.length - 1].x.toFixed(1)},${baseline} L${coords[0].x.toFixed(1)},${baseline} Z`
      : '';

  return (
    <Box>
      <svg
        role="img"
        aria-label={ariaLabel}
        viewBox={`0 0 ${WIDTH} ${height}`}
        preserveAspectRatio="none"
        style={{width: '100%', height, display: 'block'}}
      >
        <defs>
          <linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor={color} stopOpacity={0.35} />
            <stop offset="100%" stopColor={color} stopOpacity={0} />
          </linearGradient>
        </defs>
        {[0.25, 0.5, 0.75].map((fraction) => (
          <line
            key={fraction}
            x1={PADDING.left}
            x2={WIDTH - PADDING.right}
            y1={PADDING.top + plotHeight * fraction}
            y2={PADDING.top + plotHeight * fraction}
            stroke="currentColor"
            strokeOpacity={0.1}
          />
        ))}
        {area && <path d={area} fill={`url(#${gradientId})`} />}
        {line && (
          <path
            d={line}
            fill="none"
            stroke={color}
            strokeWidth={2.5}
            strokeLinejoin="round"
            strokeLinecap="round"
            vectorEffect="non-scaling-stroke"
          />
        )}
      </svg>
      <Box sx={{display: 'flex', justifyContent: 'space-between', mt: 1}}>
        {points.map((p) => (
          <Typography key={p.label} variant="caption" color="text.secondary">
            {p.label}
          </Typography>
        ))}
      </Box>
    </Box>
  );
}
