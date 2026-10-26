// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

import {createOxygenTheme, OxygenTheme} from '@wso2/oxygen-ui';

/**
 * AdminConsoleTheme - the Data Plane admin console's own theme, deliberately distinct from the
 * Control Plane console's blue glass look so the two are told apart at a glance.
 * Features: orange primary, deep rust accent, solid warm-neutral surfaces with no blur or ambient glow.
 */
export const AdminConsoleThemeConfig = {
  colorSchemes: {
    light: {
      palette: {
        primary: {
          main: '#f26b0f',
          mainChannel: '242 107 15',
          dark: '#c9540a',
          darkChannel: '201 84 10',
          light: '#f79a5a',
          lightChannel: '247 154 90',
          contrastText: '#ffffff',
        },
        secondary: {
          main: '#a13d0b',
          mainChannel: '161 61 11',
          dark: '#7d2f08',
          darkChannel: '125 47 8',
          light: '#c9672e',
          lightChannel: '201 103 46',
          contrastText: '#ffffff',
        },
        background: {
          default: '#faf6f2',
          paper: '#ffffff',
          acrylic: '#ffffff',
        },
        text: {
          primary: '#26180f',
          secondary: 'rgba(38, 24, 15, 0.72)',
        },
      },
    },
    dark: {
      palette: {
        primary: {
          main: '#ff8a3d',
          mainChannel: '255 138 61',
          dark: '#e06a1a',
          darkChannel: '224 106 26',
          light: '#ffb07a',
          lightChannel: '255 176 122',
          contrastText: '#2a1204',
        },
        secondary: {
          main: '#e5a06b',
          mainChannel: '229 160 107',
          dark: '#c4804a',
          darkChannel: '196 128 74',
          light: '#f0c29c',
          lightChannel: '240 194 156',
          contrastText: '#241204',
        },
        background: {
          default: '#16110d',
          paper: '#211913',
          acrylic: '#211913',
        },
        text: {
          primary: '#fbf4ee',
          secondary: 'rgba(251, 244, 238, 0.7)',
        },
      },
    },
  },
  shape: {
    borderRadius: 6,
  },
  components: {
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none' as const,
          boxShadow: 'none',
        },
      },
    },
    MuiTextField: {
      defaultProps: {
        size: 'small' as const,
      },
    },
    MuiSelect: {
      defaultProps: {
        size: 'small' as const,
      },
    },
    MuiAutocomplete: {
      defaultProps: {
        size: 'small' as const,
      },
    },
  },
};

const AdminConsoleTheme = createOxygenTheme(AdminConsoleThemeConfig, OxygenTheme);

export default AdminConsoleTheme;
