// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

import {useConfig} from '@thunderid/contexts';
import {useLogger} from '@thunderid/logger/react';
import {SignOutButton, User, useThunderID} from '@thunderid/react';
import {
  AppShell,
  ColorSchemeImage,
  ColorSchemeToggle,
  Divider,
  Footer,
  Header,
  Sidebar,
  UserMenu,
} from '@wso2/oxygen-ui';
import {Bot, Group, LayoutDashboard, ShieldCheck, UserRoundCog, UsersRound} from '@wso2/oxygen-ui-icons-react';
import {useMemo, type JSX, type ReactNode} from 'react';
import {useTranslation} from 'react-i18next';
import {Link as NavigateLink, Outlet, useLocation} from 'react-router';
import RouteConfig from '../configs/RouteConfig';

/** A sidebar navigation entry. */
interface NavRoute {
  id: string;
  text: string;
  icon: JSX.Element;
  path: string;
}

interface NavCategory {
  category?: string;
  routes: NavRoute[];
}

export default function DashboardLayout(): ReactNode {
  const {clearSession, discovery} = useThunderID();
  const {isTrustedIssuerGenericOidc, getTrustedIssuerClientId, getClientUrl} = useConfig();
  const {t} = useTranslation();
  const logger = useLogger();

  const handleSignOut = (signOut: () => Promise<void>): void => {
    if (isTrustedIssuerGenericOidc()) {
      try {
        clearSession();
      } catch (error: unknown) {
        logger.error('Failed to clear local session before IdP sign out', {error});
      }

      const endSessionEndpoint = discovery?.wellKnown?.end_session_endpoint;
      if (!endSessionEndpoint) {
        logger.warn('end_session_endpoint missing from IdP discovery document; ending local session only');
        // eslint-disable-next-line react-hooks/immutability
        window.location.href = getClientUrl();
        return;
      }

      const logoutUrl = new URL(endSessionEndpoint);
      logoutUrl.searchParams.set('client_id', getTrustedIssuerClientId());
      logoutUrl.searchParams.set('post_logout_redirect_uri', getClientUrl());
      // eslint-disable-next-line react-hooks/immutability
      window.location.href = logoutUrl.toString();
      return;
    }

    // Native ThunderID session: signOut() performs OIDC RP-Initiated Logout, the SDK's default
    // behavior (no client config required). It clears the local session and redirects to ThunderID's
    // end_session_endpoint, which confirms the sign-out, terminates the SSO session server-side, and
    // returns to the app. It falls back to a local-only sign out when no end_session_endpoint is
    // advertised.
    signOut().catch((error: unknown) => {
      logger.error('Sign out failed', {error});
    });
  };

  const appRoutes: NavCategory[] = useMemo(
    () => [
      {
        routes: [
          {
            id: 'dashboard',
            text: t('navigation:pages.dashboard', 'Dashboard'),
            icon: <LayoutDashboard />,
            path: RouteConfig.dashboard.list(),
          },
          {
            id: 'users',
            text: t('navigation:pages.users'),
            icon: <UsersRound />,
            path: RouteConfig.users.list(),
          },
          {
            id: 'agents',
            text: t('navigation:pages.agents', 'Agents'),
            icon: <Bot />,
            path: RouteConfig.agents.list(),
          },
          {
            id: 'groups',
            text: t('navigation:pages.groups'),
            icon: <Group />,
            path: RouteConfig.groups.list(),
          },
          {
            id: 'roles',
            text: t('navigation:pages.roles'),
            icon: <ShieldCheck />,
            path: RouteConfig.roles.list(),
          },
          {
            id: 'user-types',
            text: t('navigation:pages.userTypes'),
            icon: <UserRoundCog />,
            path: RouteConfig.userTypes.list(),
          },
        ],
      },
    ],
    [t],
  );

  const {pathname} = useLocation();

  // Resolve the active leaf item by the longest matching path prefix, so nested
  // routes (e.g. /users/:userId) still highlight their menu item.
  const activeItem = useMemo((): string | undefined => {
    const leaves = appRoutes.flatMap((group) => group.routes);
    const match = leaves
      .filter((route) => pathname === route.path || pathname.startsWith(`${route.path}/`))
      .sort((a, b) => b.path.length - a.path.length)[0];
    return match?.id;
  }, [appRoutes, pathname]);

  return (
    <AppShell>
      <AppShell.Navbar>
        <Header>
          <Header.Toggle />
          <Header.Brand>
            <Header.BrandLogo>
              <ColorSchemeImage
                src={{
                  light: `${import.meta.env.BASE_URL}/assets/images/logo.svg`,
                  dark: `${import.meta.env.BASE_URL}/assets/images/logo-inverted.svg`,
                }}
                alt={{light: 'Logo (Light)', dark: 'Logo (Dark)'}}
                height={27}
                width="auto"
                alignItems="center"
                marginBottom="3px"
              />
            </Header.BrandLogo>
            <Header.BrandTitle>Admin Console</Header.BrandTitle>
          </Header.Brand>
          <Header.Spacer />
          <Header.Actions>
            <ColorSchemeToggle />
            <Divider orientation="vertical" flexItem sx={{mx: 1, display: {xs: 'none', sm: 'block'}}} />
            <User>
              {(user) => {
                const avatar = typeof user?.picture === 'string' ? user.picture : undefined;

                return (
                  <UserMenu>
                    <UserMenu.Trigger name={String(user?.name ?? '')} avatar={avatar} showName />
                    <UserMenu.Header
                      name={String(user?.name ?? '')}
                      email={String(user?.email ?? '')}
                      avatar={avatar}
                    />
                    <UserMenu.Divider />
                    <SignOutButton>
                      {({signOut}) => (
                        <UserMenu.Logout label={t('common:userMenu.signOut')} onClick={() => handleSignOut(signOut)} />
                      )}
                    </SignOutButton>
                  </UserMenu>
                );
              }}
            </User>
          </Header.Actions>
        </Header>
      </AppShell.Navbar>

      <AppShell.Sidebar>
        <Sidebar activeItem={activeItem}>
          <Sidebar.Nav>
            {appRoutes.map((categoryGroup) => (
              <Sidebar.Category key={categoryGroup.category}>
                {categoryGroup.category && <Sidebar.CategoryLabel>{categoryGroup.category}</Sidebar.CategoryLabel>}
                {categoryGroup.routes.map((route) => (
                  <Sidebar.Item key={route.id} id={route.id} link={<NavigateLink to={route.path} />}>
                    <Sidebar.ItemIcon>{route.icon}</Sidebar.ItemIcon>
                    <Sidebar.ItemLabel>{route.text}</Sidebar.ItemLabel>
                  </Sidebar.Item>
                ))}
              </Sidebar.Category>
            ))}
          </Sidebar.Nav>
        </Sidebar>
      </AppShell.Sidebar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>

      <AppShell.Footer>
        <Footer>
          <Footer.Copyright>© {new Date().getFullYear()} The ThunderID Authors</Footer.Copyright>
        </Footer>
      </AppShell.Footer>
    </AppShell>
  );
}
