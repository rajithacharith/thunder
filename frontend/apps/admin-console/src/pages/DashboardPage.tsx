// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

import {Box, Card, Divider, Grid, PageContent, Stack, Typography} from '@wso2/oxygen-ui';
import {Activity, UserPlus, UserRoundCog, UsersRound} from '@wso2/oxygen-ui-icons-react';
import type {JSX} from 'react';
import {useTranslation} from 'react-i18next';
import AreaChart from '../components/AreaChart';
import StatTile from '../components/StatTile';

// TODO: replace with real counts once the Data Plane exposes a stats endpoint. Placeholder values
// only, so the dashboard's layout can be reviewed before any backend wiring exists.
const MOCK_STATS = {
  totalUsers: 128,
  userTypes: 4,
  newUsersThisWeek: 12,
  activeSessions: 37,
};

// TODO: replace with real per-day sign-in counts once the Data Plane records them. Placeholder
// values only, like the stats above.
const MOCK_SIGN_INS = [
  {label: 'Mon', value: 42},
  {label: 'Tue', value: 58},
  {label: 'Wed', value: 51},
  {label: 'Thu', value: 76},
  {label: 'Fri', value: 69},
  {label: 'Sat', value: 31},
  {label: 'Sun', value: 37},
];

// TODO: replace with the actual most-recently-created users once wired to the real API.
const MOCK_RECENT_USERS: {name: string; email: string; createdAgo: string}[] = [
  {name: 'Amara Chen', email: 'amara.chen@example.com', createdAgo: '2 hours ago'},
  {name: 'Diego Fontaine', email: 'diego.fontaine@example.com', createdAgo: 'Yesterday'},
  {name: 'Priya Raman', email: 'priya.raman@example.com', createdAgo: '2 days ago'},
  {name: 'Jonas Weber', email: 'jonas.weber@example.com', createdAgo: '3 days ago'},
  {name: 'Lindiwe Dube', email: 'lindiwe.dube@example.com', createdAgo: '5 days ago'},
];

export default function DashboardPage(): JSX.Element {
  const {t} = useTranslation('common');

  return (
    <PageContent>
      <Stack spacing={4}>
        <Box>
          <Typography variant="h4" fontWeight={600} gutterBottom>
            {t('dashboard.title', 'Dashboard')}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {t('dashboard.subtitle', "An overview of this deployment's users.")}
          </Typography>
        </Box>

        <Grid container spacing={2}>
          <Grid size={{xs: 12, sm: 6, md: 3}}>
            <StatTile
              icon={<UsersRound size={20} />}
              label={t('dashboard.stats.totalUsers', 'Total users')}
              value={MOCK_STATS.totalUsers}
            />
          </Grid>
          <Grid size={{xs: 12, sm: 6, md: 3}}>
            <StatTile
              icon={<UserRoundCog size={20} />}
              label={t('dashboard.stats.userTypes', 'User types')}
              value={MOCK_STATS.userTypes}
            />
          </Grid>
          <Grid size={{xs: 12, sm: 6, md: 3}}>
            <StatTile
              icon={<UserPlus size={20} />}
              label={t('dashboard.stats.newUsers', 'New users')}
              value={MOCK_STATS.newUsersThisWeek}
              note={t('dashboard.stats.newUsersNote', 'Last 7 days')}
            />
          </Grid>
          <Grid size={{xs: 12, sm: 6, md: 3}}>
            <StatTile
              icon={<Activity size={20} />}
              label={t('dashboard.stats.activeSessions', 'Active sessions')}
              value={MOCK_STATS.activeSessions}
            />
          </Grid>
        </Grid>

        <Card variant="outlined" sx={{p: 3}}>
          <Typography variant="subtitle1" fontWeight={600}>
            {t('dashboard.signIns.title', 'Sign-ins')}
          </Typography>
          <Typography variant="caption" color="text.secondary" sx={{display: 'block', mb: 2}}>
            {t('dashboard.signIns.subtitle', 'Last 7 days')}
          </Typography>
          <AreaChart
            points={MOCK_SIGN_INS}
            ariaLabel={t('dashboard.signIns.ariaLabel', 'Sign-ins per day over the last 7 days')}
          />
        </Card>

        <Card variant="outlined" sx={{p: 3}}>
          <Typography variant="subtitle1" fontWeight={600} gutterBottom>
            {t('dashboard.recentUsers.title', 'Recently added users')}
          </Typography>
          <Stack divider={<Divider />} sx={{mt: 1}}>
            {MOCK_RECENT_USERS.map((user) => (
              <Stack
                key={user.email}
                direction="row"
                justifyContent="space-between"
                alignItems="center"
                sx={{py: 1.25}}
              >
                <Box>
                  <Typography variant="body2" fontWeight={500}>
                    {user.name}
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    {user.email}
                  </Typography>
                </Box>
                <Typography variant="caption" color="text.secondary">
                  {user.createdAgo}
                </Typography>
              </Stack>
            ))}
          </Stack>
        </Card>
      </Stack>
    </PageContent>
  );
}
