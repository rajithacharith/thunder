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
import React from 'react';
import './B2CIdentityJourney.css';
import {
  UseCaseBuildingBlockDetail,
  UseCaseBuildingBlockPanel,
  UseCaseBuildingBlocksExplorer,
} from './UseCaseBuildingBlocksExplorer';
import { UseCaseCapabilityMap, UseCaseMapGroup, UseCaseMapNode } from './UseCaseCapabilityMap';

export { UseCaseBuildingBlockPanel };

interface RoadmapNode {
  href: string;
  label: string;
  icon: React.ReactNode;
}

const roadmapNodes: RoadmapNode[] = [
  {
    href: '#b2c-identity-journey',
    label: 'Sign In',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="M14 3h6v18h-6" />
        <path d="M10 12h10" />
        <path d="m7 9 3 3-3 3" />
        <path d="M4 4h8v16H4" />
      </svg>
    ),
  },
  {
    href: '#enable-self-sign-up',
    label: 'Self Sign-Up',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
        <circle cx="9" cy="7" r="4" />
        <path d="M19 8v6" />
        <path d="M16 11h6" />
      </svg>
    ),
  },
  {
    href: '#add-self-service-profile-management',
    label: 'Manage Profile',
    icon: (
      <svg viewBox="0 0 24 24">
        <circle cx="12" cy="8" r="4" />
        <path d="M4 21a8 8 0 0 1 16 0" />
      </svg>
    ),
  },
  {
    href: '#add-account-recovery',
    label: 'Recover Access',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="M7 11V9a5 5 0 0 1 10 0v2" />
        <rect x="5" y="11" width="14" height="9" rx="2" />
        <path d="M12 15v2" />
      </svg>
    ),
  },
  {
    href: '#onboard-internal-users',
    label: 'Internal Users',
    icon: (
      <svg viewBox="0 0 24 24">
        <rect x="3" y="7" width="18" height="13" rx="2" />
        <path d="M9 7V5a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2" />
        <path d="M3 13h18" />
      </svg>
    ),
  },
  {
    href: '#handle-account-closure',
    label: 'Close Accounts',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
        <circle cx="9" cy="7" r="4" />
        <path d="m17 11 4 4m0-4-4 4" />
      </svg>
    ),
  },
  {
    href: '#defend-against-abuse-and-risk',
    label: 'Defend Against Abuse',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="M12 3 5 6v6c0 4.5 3 7.5 7 9 4-1.5 7-4.5 7-9V6l-7-3Z" />
        <path d="m9 12 2 2 4-4" />
      </svg>
    ),
  },
  {
    href: '#gain-identity-insights',
    label: 'Identity Insights',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="M3 3v18h18" />
        <path d="m7 16 4-4 4 4 4-8" />
      </svg>
    ),
  },
];

const roadmapIcons = roadmapNodes.map((node) => node.icon);

const journeyDetails: UseCaseBuildingBlockDetail[] = [
  {
    id: 'sign-in',
    label: 'Sign in',
    title: 'Add Sign-In to Your Application',
    icon: roadmapIcons[0],
    why:
      'As your most visible identity surface, sign-in needs to feel effortless. Consumers expect to choose the method they already prefer, such as password, social sign-in, passkey, or passwordless sign-in.',
    example: (
      <>
        A user installs your mobile application, taps <strong>Sign in with Google</strong>, and reaches their dashboard within seconds. Later, when they try to change their email address, the application asks them to confirm a one-time code as a step-up check. A power user enables a passkey on their phone and from then on signs in with a single tap of their passkey - no password required.
      </>
    ),
    capabilityGroups: [
      {
        title: 'Authentication methods',
        items: [
          'Password sign-in',
          'Email or SMS one-time code',
          'Magic link',
          'Passkey',
          'Social sign-in',
          'Enterprise identity provider sign-in',
        ],
      },
      {
        title: 'Security controls',
        items: [
          'Multi-factor authentication',
          'Step-up authentication',
          'Persistent sign-in / remember me',
        ],
      },
    ],
  },
  {
    id: 'self-sign-up',
    label: 'Self sign-up',
    title: 'Enable Self Sign-Up',
    icon: roadmapIcons[1],
    why:
      'New users decide whether your product is worth their time in the first minute. Use self sign-up when users should be able to create accounts without administrator involvement.',
    example: (
      <>
        A new user lands on your home page, taps <strong>Sign up with Google</strong>, and arrives signed in with a basic profile already filled in. Their consent decisions are recorded during sign-up and can be revisited later from settings.
      </>
    ),
    capabilityGroups: [
      {
        title: 'Registration methods',
        items: [
          'Email and password sign-up',
          'Passwordless sign-up',
          'Social sign-up',
          'Passkey-first registration',
          'Just-in-time account creation',
        ],
      },
      {
        title: 'Profile and trust',
        items: [
          'Progressive profile collection',
          'Email or phone verification',
          'Terms and marketing consent capture',
        ],
      },
    ],
  },
  {
    id: 'manage-profile',
    label: 'Manage profile',
    title: 'Manage Customer Profile',
    icon: roadmapIcons[2],
    why:
      'Once a consumer signs in, they expect a self-service area where they can view and change their own identity without contacting support. Profile management reduces support load and improves user trust.',
    example: (
      <>
        A user opens the account page and switches from password to passkey. They also enable two-factor authentication - the application asks for their phone number and verifies it via a one-time code before activating it. They remove an old linked Google account they no longer use and sign out a session on a device they sold last month. When they change their email address, the new address is verified via a magic link before the change takes effect.
      </>
    ),
    capabilityGroups: [
      {
        title: 'Account details',
        items: [
          'View and edit profile attributes',
          'Update verified email or phone with re-verification',
          'Manage linked social and enterprise identities',
        ],
      },
      {
        title: 'Security and privacy',
        items: [
          'Change password',
          'Add or remove a passkey',
          'Manage second factors',
          'View active sessions and revoke a specific session',
          'View and withdraw stored consents',
          'Account deletion or data export',
        ],
      },
    ],
  },
  {
    id: 'recover-access',
    label: 'Recover access',
    title: 'Recover Customer Access',
    icon: roadmapIcons[3],
    why:
      'When users lose access to their account, it tests whether they stay or leave. Recovery paths should be quick for legitimate users and resistant to account takeover.',
    example: (
      <>
        A user who forgot their password requests a magic link delivered to their email, clicks it, sets a new password, and is signed back in. Another user who lost their phone uses an email one-time code for recovery because their SMS channel is no longer verified. A third, whose account was locked after too many failed attempts, is automatically unlocked after the lockout window expires.
      </>
    ),
    capabilityGroups: [
      {
        title: 'Recovery methods',
        items: [
          'Forgotten-password reset',
          'Email magic link',
          'Email one-time code',
          'SMS one-time code',
        ],
      },
      {
        title: 'Recovery controls',
        items: [
          'Recovery channel verification',
          'Account unlock after lockout',
        ],
      },
    ],
  },
  {
    id: 'internal-users',
    label: 'Internal users',
    title: 'Manage Internal Users',
    icon: roadmapIcons[4],
    why:
      'Behind every consumer product is a team keeping it running. Support agents, administrators, and operations staff need a separate onboarding path with identity and role decided before they arrive.',
    example: (
      <>
        A new support agent receives an invitation email and follows the link. They set a password, accept the support terms, and land in the admin console with the support role pre-assigned. Separately, an operations admin creates ten new staff accounts directly for a regional support team. The initial passwords are auto-generated and distributed via a secure channel; the staff members rotate them on first sign-in.
      </>
    ),
    capabilityGroups: [
      {
        title: 'User creation',
        items: [
          'Invite a user by email',
          'Create a user account directly with initial credentials',
          'Bulk invite or bulk create',
          'Configurable invitation expiry and resend',
          'Revoke a pending invitation',
        ],
      },
      {
        title: 'Access setup',
        items: [
          'Onboarding after invitation acceptance',
          'Credential setup on first sign-in',
          'Initial role or permission assignment',
        ],
      },
    ],
  },
  {
    id: 'account-closure',
    label: 'Close accounts',
    title: 'Handle Account Closure',
    icon: roadmapIcons[5],
    why:
      'Accounts have an end as well as a beginning. Users decide to leave and expect a clear way to close their account. You need to remove accounts that violate your terms. Inactive accounts pile up over time and need to age out.',
    example: (
      <>
        A user closes their account from settings, and the account and its data are removed in line with your retention policy. A separate account is suspended by an admin after a fraud flag, with the reason captured against the record. An inactive account that has not been touched for two years receives a warning email, then is expired when the user does not return.
      </>
    ),
    capabilityGroups: [
      {
        title: 'User-initiated',
        items: [
          'Self-service account closure',
          'Audit record of the closure event',
        ],
      },
      {
        title: 'Admin-initiated',
        items: [
          'Account suspension with reason capture',
          'Inactive-account detection and expiry',
          'Prior notification before expiry',
        ],
      },
    ],
  },
  {
    id: 'defend-against-abuse',
    label: 'Defend against abuse',
    title: 'Defend Against Abuse and Risk',
    icon: roadmapIcons[6],
    why:
      'Identity flows attract abuse from day one. Bots try to mass-create accounts, attackers run credential stuffing, and even legitimate users can land in risky situations. The level of friction should adapt to the risk in the moment.',
    example: (
      <>
        A sign-up wave from a single IP range hits a bot challenge and is throttled before any accounts are created. A user signing in from a new country is asked for a one-time code as a step-up, then completes sign-in normally. A repeated wrong-password pattern triggers a temporary lockout, while genuine sign-ins continue to succeed.
      </>
    ),
    capabilityGroups: [
      {
        title: 'Detection',
        items: [
          'Bot detection on sign-up and sign-in',
          'CAPTCHA or invisible challenge integration',
          'Rate limiting per user, IP address, and device',
          'Credential stuffing and brute-force detection',
          'Risk signals: new device, new location, impossible travel',
        ],
      },
      {
        title: 'Response',
        items: [
          'Adaptive step-up authentication based on risk',
          'Account lockout with automatic unlock',
        ],
      },
    ],
  },
  {
    id: 'identity-insights',
    label: 'Identity insights',
    title: 'Gain Identity Insights',
    icon: roadmapIcons[7],
    why:
      'Identity is one of the highest-signal touch points your application has with each user. Without visibility, you optimize sign-up in the dark, miss security incidents until they escalate, and scramble when compliance asks.',
    example: (
      <>
        A product manager notices the sign-up completion rate dropped after a new terms screen went live. A security lead receives an alert about a spike in failed sign-ins from one IP range. A support agent reviews a user&apos;s recent sign-in attempts to diagnose an access issue. A compliance officer exports the audit log for an annual review.
      </>
    ),
    capabilityGroups: [
      {
        title: 'Analytics',
        items: [
          'Sign-up and sign-in funnel analytics with drop-off points',
          'Adoption metrics by authentication method',
          'Active user trends and registration over time',
        ],
      },
      {
        title: 'Audit and security',
        items: [
          'Audit log of identity events',
          'Per-user activity history for support investigations',
          'Security signals and real-time alerts',
          'Stream identity events to external analytics and SIEM tools',
          'Export audit data for compliance reporting',
        ],
      },
    ],
  },
];

const crossCuttingIcons = {
  federation: (
    <svg viewBox="0 0 24 24">
      <circle cx="7" cy="12" r="3" />
      <circle cx="17" cy="7" r="3" />
      <circle cx="17" cy="17" r="3" />
      <path d="M9.5 10.5 14.5 8.5" />
      <path d="M9.5 13.5 14.5 15.5" />
    </svg>
  ),
  authorization: (
    <svg viewBox="0 0 24 24">
      <path d="M4 12h16" />
      <path d="M12 4v16" />
      <path d="M7 7h10v10H7z" />
    </svg>
  ),
  consent: (
    <svg viewBox="0 0 24 24">
      <path d="M9 11 12 14 20 6" />
      <path d="M20 12v6a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h9" />
    </svg>
  ),
  branding: (
    <svg viewBox="0 0 24 24">
      <path d="M12 4a8 8 0 1 0 0 16c1.4 0 2.5-1.1 2.5-2.5 0-1.1.9-2 2-2h1A4.5 4.5 0 0 0 22 11 7 7 0 0 0 12 4Z" />
      <circle cx="7.5" cy="11" r="1" />
      <circle cx="10.5" cy="8" r="1" />
      <circle cx="14" cy="8" r="1" />
    </svg>
  ),
  localization: (
    <svg viewBox="0 0 24 24">
      <circle cx="12" cy="12" r="10" />
      <path d="M12 2a14.5 14.5 0 0 0 0 20M12 2a14.5 14.5 0 0 1 0 20" />
      <path d="M2 12h20" />
    </svg>
  ),
  privacy: (
    <svg viewBox="0 0 24 24">
      <path d="M12 3 5 6v6c0 4.5 3 7.5 7 9 4-1.5 7-4.5 7-9V6l-7-3Z" />
      <path d="M9 12h6" />
      <path d="M12 9v6" />
    </svg>
  ),
};

const crossCuttingDetails: UseCaseBuildingBlockDetail[] = [
  {
    id: 'federated-identity',
    label: 'Federation',
    title: 'Federated Identity',
    icon: crossCuttingIcons.federation,
    why:
      'External identity providers, both social and enterprise, let users bring an identity they already have to your app. Done well, federation creates one user record per real person regardless of how many sign-in methods they use.',
    capabilityGroups: [
      {
        title: 'Identity providers',
        items: [
          'Social identity provider sign-in',
          'Enterprise OIDC identity provider sign-in',
          'Connected identity sign-out behavior',
        ],
      },
      {
        title: 'Account mapping',
        items: [
          'Just-in-time account creation',
          'Account linking',
          'Federated profile mapping',
        ],
      },
    ],
  },
  {
    id: 'authorization',
    label: 'Authorization',
    title: 'Authorization',
    icon: crossCuttingIcons.authorization,
    why:
      'When your app calls APIs on behalf of the user, it needs the right level of access. Scopes describe what the app may do and audiences describe which API the token is valid for.',
    capabilityGroups: [
      {
        title: 'Token controls',
        items: [
          'OAuth2 scopes',
          'Audience-restricted tokens',
          'Claims for application decisions',
        ],
      },
      {
        title: 'Access decisions',
        items: [
          'Role-aware access where needed',
          'API authorization',
          'Least-privilege access requests',
        ],
      },
    ],
  },
  {
    id: 'consent',
    label: 'Consent',
    title: 'Consent',
    icon: crossCuttingIcons.consent,
    why:
      'Where authorization describes what the app requests, consent is where the user agrees to it. Consent decisions should be recorded so users can review or revoke them later.',
    capabilityGroups: [
      {
        title: 'Consent capture',
        items: [
          'Profile-sharing consent',
          'Permission consent',
          'Terms of service acceptance',
          'Privacy policy acceptance',
          'Marketing preference capture',
        ],
      },
      {
        title: 'Consent lifecycle',
        items: [
          'Consent review and revocation',
          'Consent records for audit',
        ],
      },
    ],
  },
  {
    id: 'branding',
    label: 'Branding',
    title: 'Branding',
    icon: crossCuttingIcons.branding,
    why:
      'Your sign-in, sign-up, and recovery surfaces should match your brand whether they live on hosted pages or inside your own app screens.',
    capabilityGroups: [
      {
        title: 'Visual identity',
        items: [
          'Hosted page branding',
          'App-native screen consistency',
          'Logo and color customization',
          'Branded copy',
        ],
      },
      {
        title: 'Experience consistency',
        items: [
          'Localized sign-in experience',
          'Recovery flow branding',
        ],
      },
    ],
  },
  {
    id: 'localization',
    label: 'Localization',
    title: 'Localization',
    icon: crossCuttingIcons.localization,
    why:
      "Identity surfaces should speak the user's language. Sign-in, sign-up, recovery, and profile screens render in the user's locale, with right-to-left layouts where the language needs them.",
    capabilityGroups: [
      {
        title: 'Language and layout',
        items: [
          'Locale-aware identity screens',
          'Right-to-left layout support',
          'Localized notification emails and SMS',
        ],
      },
      {
        title: 'Accessibility',
        items: [
          'Keyboard navigation',
          'Screen reader support',
          'Sufficient contrast and accessibility standards',
        ],
      },
    ],
  },
  {
    id: 'privacy',
    label: 'Privacy',
    title: 'Privacy',
    icon: crossCuttingIcons.privacy,
    why:
      'Consent capture and policy alignment should be built into registration and profile interactions, not bolted on later.',
    capabilityGroups: [
      {
        title: 'Data visibility',
        items: [
          'Stored data visibility',
          'Shared data visibility',
          'Consent history',
        ],
      },
      {
        title: 'User controls',
        items: [
          'Account deletion workflows',
          'Data export workflows',
          'Privacy preference management',
        ],
      },
    ],
  },
];

const b2cRootNode: UseCaseMapNode = {
  id: 'add-login',
  href: '#add-login-to-your-application',
  label: 'Add Sign-In',
  icon: roadmapIcons[0],
};

const b2cUseCaseGroups: UseCaseMapGroup[] = [
  {
    id: 'identity-access',
    label: 'Identity & Access',
    nodes: [
      {
        id: 'self-sign-up',
        href: '#enable-self-sign-up',
        label: 'Enable Self Sign-Up',
        icon: roadmapIcons[1],
      },
      {
        id: 'recover-access',
        href: '#add-account-recovery',
        label: 'Configure Account Recovery',
        icon: roadmapIcons[3],
      },
      {
        id: 'federated-identity',
        href: '#federated-identity',
        label: 'Add Federated Sign-In',
        icon: crossCuttingIcons.federation,
      },
      {
        id: 'authorization',
        href: '#authorization',
        label: 'Authorize Access',
        icon: crossCuttingIcons.authorization,
      },
    ],
  },
  {
    id: 'administration',
    label: 'Administration',
    nodes: [
      {
        id: 'manage-profile',
        href: '#add-self-service-profile-management',
        label: 'Enable Profile Management',
        icon: roadmapIcons[2],
      },
      {
        id: 'internal-users',
        href: '#onboard-internal-users',
        label: 'Manage Internal Team Access',
        icon: roadmapIcons[4],
      },
      {
        id: 'account-closure',
        href: '#handle-account-closure',
        label: 'Handle Account Closure',
        icon: roadmapIcons[5],
      },
    ],
  },
  {
    id: 'configuration',
    label: 'Configuration',
    nodes: [
      {
        id: 'consent',
        href: '#consent',
        label: 'Manage Consent',
        icon: crossCuttingIcons.consent,
      },
      {
        id: 'branding',
        href: '#branding',
        label: 'Customize Branding',
        icon: crossCuttingIcons.branding,
      },
      {
        id: 'localization',
        href: '#localization',
        label: 'Localize Identity Surfaces',
        icon: crossCuttingIcons.localization,
      },
    ],
  },
  {
    id: 'operations',
    label: 'Operations',
    nodes: [
      {
        id: 'defend-against-abuse',
        href: '#defend-against-abuse-and-risk',
        label: 'Defend Against Abuse',
        icon: roadmapIcons[6],
      },
      {
        id: 'identity-insights',
        href: '#gain-identity-insights',
        label: 'Gain Identity Insights',
        icon: roadmapIcons[7],
      },
      {
        id: 'privacy',
        href: '#privacy',
        label: 'Protect Customer Data',
        icon: crossCuttingIcons.privacy,
      },
    ],
  },
];

export function B2CIdentityUseCaseMap() {
  return (
    <UseCaseCapabilityMap
      ariaLabel="B2C identity use case capability map"
      root={b2cRootNode}
      groups={b2cUseCaseGroups}
    />
  );
}

export function B2CIdentityJourneyExplorer() {
  return (
    <>
      <h3>Primary B2C Journeys</h3>
      <p>Each block below represents a distinct identity use case. Select one to see what it covers and which capabilities are involved.</p>
      <UseCaseBuildingBlocksExplorer
        ariaLabel="Primary B2C identity journeys"
        detailPanelId="b2c-journey-detail"
        groups={[
          {
            id: 'primary-journeys',
            nodes: journeyDetails,
          },
        ]}
      />
      <h3>Cross-Cutting Capabilities</h3>
      <p>These capabilities are not tied to a single journey. They apply across the identity system and are relevant to most B2C applications.</p>
      <UseCaseBuildingBlocksExplorer
        ariaLabel="Cross-cutting B2C identity capabilities"
        detailPanelId="b2c-cross-cutting-detail"
        groups={[
          {
            id: 'cross-cutting-capabilities',
            nodes: crossCuttingDetails,
            variant: 'secondary',
          },
        ]}
      />
    </>
  );
}

const solutionPatternNodes: RoadmapNode[] = [
  {
    href: '#redirect-based',
    label: 'Redirect-Based',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="M4 12h12" />
        <path d="m12 6 6 6-6 6" />
        <circle cx="20" cy="12" r="2" />
      </svg>
    ),
  },
  {
    href: '#app-native',
    label: 'App-Native',
    icon: (
      <svg viewBox="0 0 24 24">
        <rect x="3" y="4" width="7" height="7" rx="1" />
        <rect x="14" y="4" width="7" height="7" rx="1" />
        <rect x="3" y="13" width="7" height="7" rx="1" />
        <rect x="14" y="13" width="7" height="7" rx="1" />
      </svg>
    ),
  },
  {
    href: '#direct-api',
    label: 'Direct API',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="m8 4-6 8 6 8" />
        <path d="m16 4 6 8-6 8" />
        <path d="M14 4 10 20" />
      </svg>
    ),
  },
];

const solutionPatternDetails: UseCaseBuildingBlockDetail[] = [
  {
    id: 'redirect-based',
    label: 'Redirect-based',
    title: 'Redirect-Based',
    icon: solutionPatternNodes[0].icon,
    why:
      'ThunderID hosts the identity screens. Your app redirects users there and gets them back signed in.',
  },
  {
    id: 'app-native',
    label: 'App-native',
    title: 'App-Native',
    icon: solutionPatternNodes[1].icon,
    why:
      'Your app renders every screen, but ThunderID owns the journey — step ordering, branching, and policy stay on the server.',
  },
  {
    id: 'direct-api',
    label: 'Direct API',
    title: 'Direct API',
    icon: solutionPatternNodes[2].icon,
    why:
      "Your app calls ThunderID's primitive APIs directly — low-level, single-purpose operations with no hosted pages and no journey to configure. You decide what to call, when to call it, and what to do with the result.",
  },
];

export function B2CIntegrationApproachesCards() {
  return (
    <div className="uc-approach-cards">
      {solutionPatternDetails.map((pattern) => (
        <UseCaseBuildingBlockPanel
          key={pattern.id}
          icon={pattern.icon}
          title={pattern.title}
          why={pattern.why}
          capabilityGroups={pattern.capabilityGroups}
        />
      ))}
    </div>
  );
}

const tokensAndApisDetails: UseCaseBuildingBlockDetail[] = [
  {
    id: 'session-token-strategy',
    label: 'Session and Token Strategy',
    title: 'Session and Token Strategy',
    icon: (
      <svg viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="9" />
        <path d="M12 7v5l3 3" />
      </svg>
    ),
    why: 'Once a user signs in, the application holds a token or session that represents them. The shape of that credential decides how long the user stays signed in, how quickly access can be revoked when needed, and how much load lands on the identity product. Most B2C apps make this decision once and live with it for years, so it pays to pick deliberately.',
    capabilityGroups: [
      {
        title: 'Patterns',
        items: [
          'Stateless tokens: short-lived access tokens backed by longer refresh tokens, with no server-side session record. Scales easily; revocation waits for the token to expire.',
          'Server-backed sessions: every refresh is backed by a server-side record the identity product can revoke instantly, enabling true sign-out-everywhere.',
          'Sliding-expiry sessions: sessions extend on each use so returning users rarely have to sign in again, in return for longer-lived credentials.',
        ],
      },
      {
        title: 'Capabilities',
        items: [
          'Stateless JWT access tokens with refresh tokens',
          'Server-side session record backing each refresh',
          'Instant token and session revocation',
          'Refresh-token rotation on use',
          'Sliding-expiry sessions for stay-signed-in',
          'Configurable access, refresh, and session lifetimes',
          'Single logout across application surfaces',
        ],
      },
    ],
  },
  {
    id: 'protect-apis',
    label: 'Protect APIs',
    title: 'Protect APIs the App Calls',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="M12 2 3 7v6c0 5.5 3.8 10.7 9 12 5.2-1.3 9-6.5 9-12V7L12 2z" />
        <path d="m9 12 2 2 4-4" />
      </svg>
    ),
    why: "Sign-in is one half of identity; the other is protecting the APIs your app calls after sign-in. The identity product issues tokens during sign-in, and the same tokens carry the user's permissions to your APIs. The API does not need to know who the user is; it only needs to validate the token and check the permissions inside.",
    capabilityGroups: [
      {
        title: 'Capabilities',
        items: [
          'Issue an OAuth2 access token to the app on sign-in',
          'Validate the token at the API edge via JWT signature verification or introspection',
          'Scope-based authorization checks',
          'Bind a token to a specific API via audience or resource indicator',
          'Group related APIs into a resource server so they share permission rules',
        ],
      },
    ],
  },
];

/** @deprecated use B2CIntegrationApproachesCards */
export const B2CIntegrationApproachesRoadmap = B2CIntegrationApproachesCards;

export function B2CTokensAndApisCards() {
  return (
    <div className="uc-approach-cards">
      {tokensAndApisDetails.map((item) => (
        <UseCaseBuildingBlockPanel
          key={item.id}
          icon={item.icon}
          title={item.title}
          why={item.why}
          capabilityGroups={item.capabilityGroups}
        />
      ))}
    </div>
  );
}

const operationsDetails: UseCaseBuildingBlockDetail[] = [
  {
    id: 'identity-as-code',
    label: 'Identity-as-Code',
    title: 'Identity-as-Code',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
        <path d="M14 2v6h6" />
        <path d="M9 15 7 13l2-2" />
        <path d="m13 11 2 2-2 2" />
      </svg>
    ),
    why: 'Identity configuration spans many resources such as user types, applications, roles, flows, federated providers, and branding, and it grows over time. Managing it by hand breaks down as soon as you have separate dev, staging, and production environments. Identity-as-Code treats configuration as versioned source files the team reviews, tests, and promotes.',
    capabilityGroups: [
      {
        title: 'Capabilities',
        items: [
          'Declarative configuration files for identity resources',
          "Environment-specific values supplied separately from the configuration",
          "Version control, review, and rollback through the application's existing workflow",
          'Promotion of changes across dev, staging, and production tenants',
          'Drift detection between configuration files and the live tenant',
        ],
      },
    ],
  },
  {
    id: 'resilience',
    label: 'Resilience',
    title: 'Resilience and Multi-Region Deployment',
    icon: (
      <svg viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="9" />
        <path d="M3.6 9h16.8M3.6 15h16.8" />
        <path d="M12 3a14.5 14.5 0 0 0 0 18M12 3a14.5 14.5 0 0 1 0 18" />
      </svg>
    ),
    why: 'The identity product is on the critical path for every sign-in. If it is down, no one gets in. The deployment shape needs to match the availability you have promised your users. The right pattern depends on your reliability target, latency budget, and where your users are.',
    capabilityGroups: [
      {
        title: 'Deployment patterns',
        items: [
          'Single-region: simplest deployment, fits regional audiences and modest availability targets',
          'Active-passive: replicates user data to a second region that takes over on failure',
          'Active-active: runs every region live, routing each user to the nearest healthy region for lowest latency and smallest blast radius',
          'Regional sharding: pins specific users to specific regions, often to satisfy data-residency rules',
        ],
      },
      {
        title: 'Capabilities',
        items: [
          'Single-region deployment',
          'Active-passive deployment with failover',
          'Active-active deployment with regional routing',
          'Regional sharding pinned to data residency',
          'User-data replication across regions',
          'Health-checked routing and failover',
          'Independent regional configuration and secrets',
        ],
      },
    ],
  },
  {
    id: 'monitoring',
    label: 'Activity Monitoring',
    title: 'Activity Monitoring and Audit',
    icon: (
      <svg viewBox="0 0 24 24">
        <path d="M3 3v18h18" />
        <path d="m7 16 4-8 4 4 3-6" />
      </svg>
    ),
    why: 'Identity events are some of the highest-signal data in your stack. Every sign-up, sign-in, recovery, and consent change is a row worth keeping. Product teams want funnel and adoption metrics. Security teams want anomalies. Compliance teams want a durable audit log they can hand to a regulator.',
    capabilityGroups: [
      {
        title: 'Capabilities',
        items: [
          'Structured audit log of every identity event with a consistent schema',
          'Built-in dashboards for sign-up funnels, sign-in success rates, and method adoption',
          'Per-user activity timelines for support investigation',
          'Stream events to external analytics, data warehouse, or SIEM platforms',
          'Anomaly detection: failed sign-in spikes, brute force, impossible travel',
          'Real-time alerts for security-relevant events',
          'Export audit logs for compliance reporting',
        ],
      },
    ],
  },
  {
    id: 'connect-to-systems',
    label: 'Connect to Systems',
    title: 'Connect to Other Systems',
    icon: (
      <svg viewBox="0 0 24 24">
        <circle cx="18" cy="5" r="2" />
        <circle cx="6" cy="12" r="2" />
        <circle cx="18" cy="19" r="2" />
        <path d="m8 11.5 8.1-5M8 12.5l8.1 5.1" />
      </svg>
    ),
    why: 'Identity is not an isolated system; it belongs to the business. New sign-ups should land in your CRM; password resets should trigger security event logs; failed sign-ins should surface in your monitoring. The identity product is the natural place to emit these signals because it sees every identity event.',
    capabilityGroups: [
      {
        title: 'Capabilities',
        items: [
          'Emit identity events (sign-up, password change, and so on) to subscribers',
          'Push new-user data to a CRM or marketing tool on sign-up',
          'Send SMS and email notifications through providers you choose',
          'Call out to your own systems mid-journey for validation or data lookup',
          'Run checks before or after key identity actions',
        ],
      },
    ],
  },
];

export function B2COperationsCards() {
  return (
    <div className="uc-approach-cards">
      {operationsDetails.map((item) => (
        <UseCaseBuildingBlockPanel
          key={item.id}
          icon={item.icon}
          title={item.title}
          why={item.why}
          capabilityGroups={item.capabilityGroups}
        />
      ))}
    </div>
  );
}

interface SolutionArchitectureOption {
  id: string;
  label: string;
  description: string;
  graphic: {
    left: string;
    center: string;
    right: string;
    notes: string[];
  };
}

interface SolutionArchitectureStage {
  id: string;
  title: string;
  question: string;
  description: string;
  options: SolutionArchitectureOption[];
}

const solutionArchitectureStages: SolutionArchitectureStage[] = [
  {
    id: 'integration',
    title: 'Choose an integration approach',
    question: 'Where should identity screens live, and who should drive the journey?',
    description: 'Decide who owns the screens and who drives the journey.',
    options: [
      {
        id: 'redirect',
        label: 'Redirect-based',
        description: 'ThunderID hosts the screens and controls the identity journey.',
        graphic: {
          left: 'Your app',
          center: 'Hosted ThunderID journey',
          right: 'Signed-in user',
          notes: ['Sign-in, sign-up, recovery, consent', 'Tokens return to the app', 'Fastest secure path'],
        },
      },
      {
        id: 'app-native',
        label: 'App-native',
        description: 'Your application renders screens while ThunderID controls journey policy.',
        graphic: {
          left: 'App screens',
          center: 'ThunderID journey state',
          right: 'Next step or tokens',
          notes: ['Custom UI', 'Server-controlled branching', 'SDK-driven flow calls'],
        },
      },
      {
        id: 'direct-api',
        label: 'Direct API',
        description: 'Your application owns the screens and composes identity primitives.',
        graphic: {
          left: 'App orchestration',
          center: 'Identity APIs',
          right: 'App-managed outcome',
          notes: ['Maximum control', 'No guided journey', 'More app-side responsibility'],
        },
      },
    ],
  },
  {
    id: 'identity-data',
    title: 'Choose identity sources and data',
    question: 'Where do users come from, and which system owns the user record?',
    description: 'Decide where consumer identities come from and where user records live.',
    options: [
      {
        id: 'managed-directory',
        label: 'ThunderID user store',
        description: 'ThunderID owns the canonical consumer user record.',
        graphic: {
          left: 'Consumers',
          center: 'ThunderID user store',
          right: 'Application profile',
          notes: ['Directly managed users', 'Recovery and profile features', 'Simple operating model'],
        },
      },
      {
        id: 'federation',
        label: 'Federation',
        description: 'Users bring identities from social or enterprise providers.',
        graphic: {
          left: 'External IdPs',
          center: 'Federated sign-in',
          right: 'Linked user',
          notes: ['Social and enterprise OIDC', 'Just-in-time provisioning', 'Home-realm discovery'],
        },
      },
      {
        id: 'mixed',
        label: 'Mixed model',
        description: 'Some users live in ThunderID while other users come from external sources.',
        graphic: {
          left: 'Local and external users',
          center: 'Account linking',
          right: 'One customer identity',
          notes: ['Migration-friendly', 'Segment-specific sources', 'Flexible source of truth'],
        },
      },
    ],
  },
  {
    id: 'tokens-apis',
    title: 'Design tokens, sessions, and APIs',
    question: 'How should the signed-in user be represented to the app and APIs?',
    description: 'Decide what the app receives after sign-in and how APIs trust it.',
    options: [
      {
        id: 'stateless',
        label: 'Stateless tokens',
        description: 'Use short-lived access tokens backed by refresh tokens.',
        graphic: {
          left: 'Application',
          center: 'JWT access token',
          right: 'Protected APIs',
          notes: ['Scales easily', 'Signature validation', 'Best-effort revocation'],
        },
      },
      {
        id: 'revocable',
        label: 'Server-backed sessions',
        description: 'Back refresh and session behavior with revocable server-side records.',
        graphic: {
          left: 'Application session',
          center: 'ThunderID session record',
          right: 'Revocable access',
          notes: ['Sign out everywhere', 'Fast revocation', 'Session lookup'],
        },
      },
      {
        id: 'api-protection',
        label: 'API protection',
        description: 'Shape scopes, audiences, and resource servers around the APIs the app calls.',
        graphic: {
          left: 'Access token',
          center: 'API gateway or middleware',
          right: 'Resource server',
          notes: ['Scopes and permissions', 'Audience validation', 'Shared policy boundary'],
        },
      },
    ],
  },
  {
    id: 'operations',
    title: 'Plan operations and integrations',
    question: 'How should identity configuration run, change, and emit signals?',
    description: 'Decide how identity configuration runs, scales, and connects to your stack.',
    options: [
      {
        id: 'identity-as-code',
        label: 'Identity-as-Code',
        description: 'Keep identity resources in versioned configuration and promote them across environments.',
        graphic: {
          left: 'Configuration files',
          center: 'Deployment pipeline',
          right: 'ThunderID tenants',
          notes: ['Review and rollback', 'Environment values', 'Drift detection'],
        },
      },
      {
        id: 'resilience',
        label: 'Resilience',
        description: 'Choose a deployment shape that matches availability, latency, and residency needs.',
        graphic: {
          left: 'Regional traffic',
          center: 'Healthy region routing',
          right: 'Available sign-in',
          notes: ['Single-region', 'Active-passive or active-active', 'Regional sharding'],
        },
      },
      {
        id: 'events',
        label: 'Audit and events',
        description: 'Send identity activity to dashboards, audit storage, and business systems.',
        graphic: {
          left: 'Identity events',
          center: 'Audit and event stream',
          right: 'Analytics, SIEM, CRM',
          notes: ['Structured audit log', 'Real-time alerts', 'Webhooks and enrichment'],
        },
      },
    ],
  },
];

const solutionCrossCuttingChoices = [
  'Application type',
  'Branding',
  'Session model',
  'Data residency',
  'Token lifetimes',
];

export function B2CSolutionArchitectureMap() {
  const [activeStageId, setActiveStageId] = React.useState(solutionArchitectureStages[0].id);
  const [selectedOptions, setSelectedOptions] = React.useState<Record<string, string>>(
    Object.fromEntries(solutionArchitectureStages.map((stage) => [stage.id, stage.options[0].id])),
  );

  const activeStage =
    solutionArchitectureStages.find((stage) => stage.id === activeStageId) ?? solutionArchitectureStages[0];
  const activeOption =
    activeStage.options.find((option) => option.id === selectedOptions[activeStage.id]) ?? activeStage.options[0];

  const selectOption = (stageId: string, optionId: string) => {
    setSelectedOptions((current) => ({
      ...current,
      [stageId]: optionId,
    }));
  };

  return (
    <section className="uc-solution-map" aria-label="B2C solution architecture decision map">
      <div className="uc-solution-map__rail" aria-hidden />
      <div className="uc-solution-map__stages" role="tablist" aria-label="Solution decision sections">
        {solutionArchitectureStages.map((stage, index) => (
          <button
            key={stage.id}
            type="button"
            role="tab"
            aria-selected={activeStage.id === stage.id}
            className={`uc-solution-map__stage${activeStage.id === stage.id ? ' uc-solution-map__stage--active' : ''}`}
            onClick={() => setActiveStageId(stage.id)}
          >
            <span className="uc-solution-map__index" aria-hidden>
              {index + 1}
            </span>
            <span className="uc-solution-map__content">
              <span className="uc-solution-map__content-title">{stage.title}</span>
              <span className="uc-solution-map__content-description">{stage.description}</span>
              <span className="uc-solution-map__selection">
                {stage.options.find((option) => option.id === selectedOptions[stage.id])?.label ?? stage.options[0].label}
              </span>
            </span>
          </button>
        ))}
      </div>
      <article className={`uc-solution-map__detail uc-solution-map__detail--${activeStage.id}`} role="tabpanel">
        <div className="uc-solution-map__question">
          <span className="uc-solution-chooser__eyebrow">Decision {solutionArchitectureStages.indexOf(activeStage) + 1}</span>
          <h3>{activeStage.question}</h3>
          <div className="uc-solution-map__options" role="group" aria-label={activeStage.question}>
            {activeStage.options.map((option) => (
              <button
                key={option.id}
                type="button"
                className={`uc-solution-map__option${activeOption.id === option.id ? ' uc-solution-map__option--active' : ''}`}
                onClick={() => selectOption(activeStage.id, option.id)}
              >
                <strong>{option.label}</strong>
                <span>{option.description}</span>
              </button>
            ))}
          </div>
        </div>
        <div className="uc-solution-map__graphic" aria-label={`${activeOption.label} architecture graphic`}>
          <div className="uc-solution-map__graphic-flow">
            <div className="uc-solution-map__graphic-node">{activeOption.graphic.left}</div>
            <div className="uc-solution-map__graphic-arrow" aria-hidden>
              →
            </div>
            <div className="uc-solution-map__graphic-node uc-solution-map__graphic-node--primary">
              {activeOption.graphic.center}
            </div>
            <div className="uc-solution-map__graphic-arrow" aria-hidden>
              →
            </div>
            <div className="uc-solution-map__graphic-node">{activeOption.graphic.right}</div>
          </div>
          <ul>
            {activeOption.graphic.notes.map((note) => (
              <li key={note}>{note}</li>
            ))}
          </ul>
        </div>
      </article>
      <aside className="uc-solution-map__cross-cutting" aria-label="Cross-cutting choices">
        <div>
          <h3>Apply cross-cutting choices</h3>
          <p>These choices affect every stage of the solution.</p>
        </div>
        <ul>
          {solutionCrossCuttingChoices.map((choice) => (
            <li key={choice}>{choice}</li>
          ))}
        </ul>
      </aside>
    </section>
  );
}

export function B2CIdentitySourcesDataGraphic() {
  return (
    <section className="uc-identity-sources-diagram" aria-label="Identity sources and data overview">
      <div className="uc-identity-sources-diagram__hero">
        <div className="uc-identity-sources-diagram__hero-icon">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <circle cx="12" cy="8" r="4" />
            <path d="M4 20c0-4 3.6-7 8-7s8 3 8 7" />
          </svg>
        </div>
        <strong>User identity</strong>
      </div>

      <div className="uc-identity-sources-diagram__fork" aria-hidden="true">
        <svg viewBox="0 0 200 32" preserveAspectRatio="none">
          <line x1="100" y1="0" x2="100" y2="16" />
          <line x1="30" y1="16" x2="170" y2="16" />
          <line x1="30" y1="16" x2="30" y2="32" />
          <line x1="170" y1="16" x2="170" y2="32" />
        </svg>
      </div>

      <div className="uc-identity-sources-diagram__questions">
        <div className="uc-identity-sources-diagram__question">
          <div className="uc-identity-sources-diagram__question-head">
            <span className="uc-identity-sources-diagram__eyebrow">Question 1</span>
            <h3><a href="#identity-federation">How does identity enter the app?</a></h3>
            <p>How consumer identities arrive at your application.</p>
          </div>
          <div className="uc-identity-sources-diagram__item-group">
            <h4>Identity providers</h4>
            <ul className="uc-identity-sources-diagram__items">
              <li>
                <span className="uc-identity-sources-diagram__item-icon">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <circle cx="12" cy="12" r="9" />
                    <path d="M12 3c-2.5 3-4 5.6-4 9s1.5 6 4 9M12 3c2.5 3 4 5.6 4 9s-1.5 6-4 9M3 12h18" />
                  </svg>
                </span>
                <span>
                  <strong>Social sign-in</strong>
                  <span>Google, GitHub, and other consumer providers</span>
                </span>
              </li>
              <li>
                <span className="uc-identity-sources-diagram__item-icon">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <rect x="3" y="6" width="18" height="13" rx="2" />
                    <path d="M8 6V4h8v2" />
                  </svg>
                </span>
                <span>
                  <strong>Enterprise OIDC</strong>
                  <span>Connect enterprise identity providers</span>
                </span>
              </li>
            </ul>
          </div>

          <div className="uc-identity-sources-diagram__item-group">
            <h4>Federation decisions</h4>
            <ul className="uc-identity-sources-diagram__items">
              <li>
                <span className="uc-identity-sources-diagram__item-icon">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <circle cx="12" cy="12" r="9" />
                    <path d="M12 8v8M8 12h8" />
                  </svg>
                </span>
                <span>
                  <strong>Provisioning model</strong>
                  <span>JIT on first sign-in, or invitation-only onboarding</span>
                </span>
              </li>
              <li>
                <span className="uc-identity-sources-diagram__item-icon">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M8 12h8" />
                    <circle cx="6" cy="12" r="2" />
                    <circle cx="18" cy="12" r="2" />
                  </svg>
                </span>
                <span>
                  <strong>Account linking policy</strong>
                  <span>Verified email, explicit user action, or both</span>
                </span>
              </li>
              <li>
                <span className="uc-identity-sources-diagram__item-icon">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M8 7H5a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3" />
                    <path d="M8 7l2.5-3h3L16 7" />
                  </svg>
                </span>
                <span>
                  <strong>Single logout</strong>
                  <span>Sign the user out of your app and connected provider</span>
                </span>
              </li>
              <li>
                <span className="uc-identity-sources-diagram__item-icon">
                  <svg viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M3 9l9-6 9 6v10a2 2 0 0 1-2 2h-3" />
                    <path d="M12 12h6M15 9l3 3-3 3" />
                  </svg>
                </span>
                <span>
                  <strong>Home-realm discovery</strong>
                  <span>Route users to the right provider by email domain</span>
                </span>
              </li>
            </ul>
          </div>
        </div>

        <div className="uc-identity-sources-diagram__question">
          <div className="uc-identity-sources-diagram__question-head">
            <span className="uc-identity-sources-diagram__eyebrow">Question 2</span>
            <h3><a href="#user-stores">Where is it stored?</a></h3>
            <p>Which system owns the canonical user record.</p>
          </div>
          <ul className="uc-identity-sources-diagram__items">
            <li>
              <span className="uc-identity-sources-diagram__item-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <ellipse cx="12" cy="7" rx="9" ry="3" />
                  <path d="M3 7v5c0 1.66 4.03 3 9 3s9-1.34 9-3V7M3 12v5c0 1.66 4.03 3 9 3s9-1.34 9-3v-5" />
                </svg>
              </span>
              <span>
                <strong>Product-managed directory</strong>
                <span>ThunderID owns the canonical user record</span>
              </span>
            </li>
            <li>
              <span className="uc-identity-sources-diagram__item-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <circle cx="12" cy="12" r="9" />
                  <path d="M9 12h6M12 9l3 3-3 3" />
                </svg>
              </span>
              <span>
                <strong>Federated-only</strong>
                <span>No local record; identity stays with the provider</span>
              </span>
            </li>
            <li>
              <span className="uc-identity-sources-diagram__item-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <rect x="2" y="3" width="8" height="5" rx="1" />
                  <rect x="14" y="3" width="8" height="5" rx="1" />
                  <rect x="2" y="16" width="8" height="5" rx="1" />
                  <path d="M6 8v4h12V8M18 16v-4" />
                </svg>
              </span>
              <span>
                <strong>External directory</strong>
                <span>LDAP or custom backing store you control</span>
              </span>
            </li>
            <li>
              <span className="uc-identity-sources-diagram__item-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <rect x="3" y="9" width="7" height="6" rx="1" />
                  <rect x="14" y="9" width="7" height="6" rx="1" />
                  <path d="M10 12h4" />
                </svg>
              </span>
              <span>
                <strong>Mixed</strong>
                <span>Some managed, others federated</span>
              </span>
            </li>
          </ul>
        </div>
      </div>
    </section>
  );
}

export function B2CSolutionPatternsExplorer() {
  const [uiOwner, setUiOwner] = React.useState<'thunderid' | 'app'>('thunderid');
  const [journeyOwner, setJourneyOwner] = React.useState<'thunderid' | 'app'>('thunderid');

  const selectedPatternId =
    uiOwner === 'app' && journeyOwner === 'app'
      ? 'direct-api'
      : uiOwner === 'app'
        ? 'app-native'
        : 'redirect-based';
  const selectedPattern =
    solutionPatternDetails.find((p) => p.id === selectedPatternId) ?? solutionPatternDetails[0];

  const selectPattern = (patternId: string) => {
    if (patternId === 'direct-api') {
      setUiOwner('app');
      setJourneyOwner('app');
      return;
    }

    if (patternId === 'app-native') {
      setUiOwner('app');
      setJourneyOwner('thunderid');
      return;
    }

    setUiOwner('thunderid');
    setJourneyOwner('thunderid');
  };

  const selectUiOwner = (owner: 'thunderid' | 'app') => {
    setUiOwner(owner);
    if (owner === 'thunderid') setJourneyOwner('thunderid');
  };

  const isFirstRender = React.useRef(true);
  React.useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return;
    }
    if (typeof window !== 'undefined') {
      window.dispatchEvent(
        new CustomEvent('thunder:pattern-selected', {
          detail: { id: selectedPatternId, label: selectedPattern.label },
        })
      );
    }
  }, [selectedPatternId]);

  return (
    <section className="uc-solution-chooser" aria-label="B2C solution pattern chooser">
      <div className="uc-solution-chooser__questions" aria-label="Architecture decisions">
        <div className="uc-solution-chooser__question">
          <div>
            <span className="uc-solution-chooser__eyebrow">Decision 1</span>
            <h3>Who owns the identity screens?</h3>
            <p>Choose where users see sign-in, sign-up, recovery, and consent screens.</p>
          </div>
          <div className="uc-solution-chooser__options" role="group" aria-label="Choose who owns the identity screens">
            <button
              type="button"
              className={`uc-solution-chooser__option${uiOwner === 'thunderid' ? ' uc-solution-chooser__option--active' : ''}`}
              onClick={() => selectUiOwner('thunderid')}
            >
              ThunderID
            </button>
            <button
              type="button"
              className={`uc-solution-chooser__option${uiOwner === 'app' ? ' uc-solution-chooser__option--active' : ''}`}
              onClick={() => selectUiOwner('app')}
            >
              Your application
            </button>
          </div>
        </div>

        <div className="uc-solution-chooser__question">
          <div>
            <span className="uc-solution-chooser__eyebrow">Decision 2</span>
            <h3>Who owns the identity journey?</h3>
            <p>Choose who decides the next step, applies policy, and handles branching.</p>
          </div>
          <div className="uc-solution-chooser__options" role="group" aria-label="Choose who owns the identity journey">
            <button
              type="button"
              className={`uc-solution-chooser__option${journeyOwner === 'thunderid' ? ' uc-solution-chooser__option--active' : ''}`}
              onClick={() => setJourneyOwner('thunderid')}
            >
              ThunderID
            </button>
            <button
              type="button"
              className={`uc-solution-chooser__option${journeyOwner === 'app' ? ' uc-solution-chooser__option--active' : ''}`}
              disabled={uiOwner === 'thunderid'}
              onClick={() => setJourneyOwner('app')}
            >
              Your application
            </button>
          </div>
        </div>
      </div>

      <div className="uc-solution-chooser__patterns" role="tablist" aria-label="Solution patterns">
        {solutionPatternDetails.map((pattern) => (
          <button
            key={pattern.id}
            type="button"
            role="tab"
            aria-selected={selectedPattern.id === pattern.id}
            className={`uc-building-block-node${selectedPattern.id === pattern.id ? ' uc-building-block-node--active' : ''}`}
            onClick={() => selectPattern(pattern.id)}
          >
            <span className="uc-building-block-node__icon" aria-hidden>
              {pattern.icon}
            </span>
            {selectedPattern.id === pattern.id && <span className="uc-solution-chooser__recommended">Recommended</span>}
            <span className="uc-building-block-node__label">{pattern.label}</span>
          </button>
        ))}
      </div>

      <article className="uc-building-blocks__panel uc-solution-chooser__detail" role="tabpanel">
        <div className="uc-building-blocks__body">
          <p>{selectedPattern.why}</p>
          <a href={`#${selectedPattern.id}`} className="uc-solution-chooser__rec-link">
            Read the {selectedPattern.title} details
            <svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round" width={14} height={14}>
              <path d="M5 12h14M12 5l7 7-7 7" />
            </svg>
          </a>
        </div>
      </article>
    </section>
  );
}

export function B2CIdentityJourneyRoadmap() {
  return (
    <nav className="uc-b2c-roadmap" aria-label="B2C identity use case roadmap">
      {roadmapNodes.map((node) => (
        <a key={node.href} href={node.href} className="uc-b2c-roadmap__node">
          <span className="uc-b2c-roadmap__icon" aria-hidden>
            {node.icon}
          </span>
          <span className="uc-b2c-roadmap__label">{node.label}</span>
        </a>
      ))}
    </nav>
  );
}

export function B2CSolutionPatternsRoadmap() {
  return (
    <nav className="uc-b2c-roadmap" aria-label="B2C solution pattern roadmap">
      {solutionPatternNodes.map((node) => (
        <a key={node.href} href={node.href} className="uc-b2c-roadmap__node">
          <span className="uc-b2c-roadmap__icon" aria-hidden>
            {node.icon}
          </span>
          <span className="uc-b2c-roadmap__label">{node.label}</span>
        </a>
      ))}
    </nav>
  );
}

interface ArchDecisionCard {
  id: 'integration' | 'identity-sources' | 'tokens-and-apis' | 'operations';
  title: string;
  question: string;
  href: string;
  icon: React.ReactNode;
}

const b2cArchDecisions: ArchDecisionCard[] = [
  {
    id: 'integration',
    title: 'Integration Pattern',
    question: 'Where do identity screens live, and who controls the journey?',
    href: '../integration-patterns',
    icon: (
      <svg viewBox="0 0 24 24">
        <circle cx="12" cy="18" r="3" />
        <circle cx="6" cy="6" r="3" />
        <circle cx="18" cy="6" r="3" />
        <path d="M18 9v2c0 .6-.4 1-1 1H7c-.6 0-1-.4-1-1V9" />
        <path d="M12 12v3" />
      </svg>
    ),
  },
  {
    id: 'identity-sources',
    title: 'Identity Sources',
    question: 'Where do identities come from, and which system owns the record?',
    href: '../identity-sources',
    icon: (
      <svg viewBox="0 0 24 24">
        <ellipse cx="12" cy="5" rx="9" ry="3" />
        <path d="M3 5v14c0 1.66 4.03 3 9 3s9-1.34 9-3V5" />
        <path d="M3 12c0 1.66 4.03 3 9 3s9-1.34 9-3" />
      </svg>
    ),
  },
  {
    id: 'tokens-and-apis',
    title: 'Tokens & APIs',
    question: 'How are post-sign-in credentials shaped, and how do your APIs validate them?',
    href: '../tokens-and-apis',
    icon: (
      <svg viewBox="0 0 24 24">
        <circle cx="7.5" cy="15.5" r="5.5" />
        <path d="m21 2-9.6 9.6" />
        <path d="m15.5 7.5 3 3L22 7l-3-3" />
      </svg>
    ),
  },
  {
    id: 'operations',
    title: 'Run & Observe',
    question: 'How do you configure, deploy, monitor, and connect the identity system?',
    href: '../operations',
    icon: (
      <svg viewBox="0 0 24 24">
        <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
      </svg>
    ),
  },
];

function GlassCard({
  href = undefined,
  className = '',
  children = undefined,
}: {
  href?: string;
  className?: string;
  children?: React.ReactNode;
}) {
  const cls = ['uc-glass-card', className].filter(Boolean).join(' ');
  if (href) return <Link to={href} className={cls}>{children}</Link>;
  return <div className={cls}>{children}</div>;
}

export function B2CNextSteps({ href = './try-it-out' }: { href?: string } = {}) {
  const [patternLabel, setPatternLabel] = React.useState<string | null>(null);

  React.useEffect(() => {
    const handle = (e: Event) => {
      setPatternLabel((e as CustomEvent<{ label: string }>).detail.label);
    };
    window.addEventListener('thunder:pattern-selected', handle);
    return () => window.removeEventListener('thunder:pattern-selected', handle);
  }, []);

  return (
    <div className="uc-next-steps">
      <GlassCard href={href} className="uc-next-steps__try">
        <div className="uc-next-steps__try-eyebrow">Try It Out</div>
        <div className="uc-next-steps__try-title">
          {patternLabel ? `See ${patternLabel} working in practice` : 'See your pattern working in practice'}
        </div>
        <p className="uc-next-steps__try-desc">
          Walk through a working B2C setup and see how your selected integration pattern behaves end to end.
        </p>
        <span className="uc-next-steps__try-btn">Start the walkthrough &#8594;</span>
      </GlassCard>
    </div>
  );
}

export function B2CArchitectureDecisions({
  currentDecision,
  prioritizeIntegration = false,
}: {
  currentDecision?: ArchDecisionCard['id'];
  prioritizeIntegration?: boolean;
} = {}) {
  const cards = currentDecision
    ? b2cArchDecisions.filter((d) => d.id !== currentDecision)
    : b2cArchDecisions;

  if (prioritizeIntegration) {
    const integration = b2cArchDecisions.find((d) => d.id === 'integration') ?? b2cArchDecisions[0];
    const supporting = b2cArchDecisions.filter((d) => d.id !== 'integration');
    return (
      <div className="uc-arch-decisions uc-arch-decisions--prioritized">
        <div className="uc-arch-decisions__primary">
          <span className="uc-arch-decisions__step">Start here</span>
          <GlassCard href={integration.href} className="uc-arch-decision-card uc-arch-decision-card--primary">
            <div className="useCaseJourneyCardIcon">{integration.icon}</div>
            <div className="uc-arch-decision-card__body">
              <div className="uc-arch-decision-card__title">{integration.title}</div>
              <p className="uc-arch-decision-card__question">{integration.question}</p>
              <span className="uc-arch-decision-card__cta">Choose an integration pattern &#8594;</span>
            </div>
          </GlassCard>
        </div>
        <div className="uc-arch-decisions__supporting">
          <span className="uc-arch-decisions__step">Supporting decisions</span>
          <div className="uc-arch-decisions__supporting-grid">
            {supporting.map((d) => (
              <GlassCard key={d.id} href={d.href} className="uc-arch-decision-card uc-arch-decision-card--supporting">
                <div className="useCaseJourneyCardIcon">{d.icon}</div>
                <div className="uc-arch-decision-card__title">{d.title}</div>
                <p className="uc-arch-decision-card__question">{d.question}</p>
                <span className="uc-arch-decision-card__cta">Explore &#8594;</span>
              </GlassCard>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="uc-arch-decisions">
      <div className="uc-arch-decisions__grid">
        {cards.map((d) => (
          <GlassCard key={d.id} href={d.href} className="uc-arch-decision-card">
            <div className="useCaseJourneyCardIcon">{d.icon}</div>
            <div className="uc-arch-decision-card__title">{d.title}</div>
            <p className="uc-arch-decision-card__question">{d.question}</p>
            <span className="uc-arch-decision-card__cta">Explore &#8594;</span>
          </GlassCard>
        ))}
      </div>
    </div>
  );
}
