// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

import {render, screen} from '@thunderid/test-utils';
import {describe, expect, it} from 'vitest';
import AreaChart from '../AreaChart';

const POINTS = [
  {label: 'Mon', value: 10},
  {label: 'Tue', value: 30},
  {label: 'Wed', value: 20},
];

describe('AreaChart', () => {
  it('exposes an accessible description and shows every x label', () => {
    render(<AreaChart points={POINTS} ariaLabel="Sign-ins per day" />);

    expect(screen.getByRole('img', {name: 'Sign-ins per day'})).toBeInTheDocument();
    POINTS.forEach((p) => expect(screen.getByText(p.label)).toBeInTheDocument());
  });

  it('draws one line segment per point and closes the area', () => {
    const {container} = render(<AreaChart points={POINTS} ariaLabel="chart" />);

    const paths = container.querySelectorAll('path');
    expect(paths).toHaveLength(2);
    expect(paths[0].getAttribute('d')).toMatch(/^M[\d.]+,[\d.]+ L[\d.]+,[\d.]+ L[\d.]+,[\d.]+ L.* Z$/);
    expect(paths[1].getAttribute('d')).toMatch(/^M[\d.]+,[\d.]+ L[\d.]+,[\d.]+ L[\d.]+,[\d.]+$/);
  });

  it('renders without drawing a path when there is no data', () => {
    const {container} = render(<AreaChart points={[]} ariaLabel="empty chart" />);

    expect(container.querySelectorAll('path')).toHaveLength(0);
    expect(screen.getByRole('img', {name: 'empty chart'})).toBeInTheDocument();
  });
});
