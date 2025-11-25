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

import type {JSX} from 'react';
import {
  Box,
  Button,
  Checkbox,
  Divider,
  FormLabel,
  FormControl,
  Alert,
  FormControlLabel,
  TextField,
  Typography,
  styled,
  AlertTitle,
  Paper,
  ColorSchemeImage,
  Stack,
  IconButton,
  InputAdornment,
} from '@wso2/oxygen-ui';
import {useState} from 'react';
import {SignIn, SignUp} from '@asgardeo/react';
import {Smartphone, Google, Facebook, GitHub, Eye, EyeClosed} from '@wso2/oxygen-ui-icons-react';
import {useNavigate, useSearchParams} from 'react-router';
import ROUTES from '../../constants/routes';

const StyledPaper = styled(Paper)(({theme}) => ({
  display: 'flex',
  flexDirection: 'column',
  alignSelf: 'center',
  width: '100%',
  padding: theme.spacing(4),
  gap: theme.spacing(2),
  [theme.breakpoints.up('sm')]: {
    width: '450px',
  },
}));

export default function SignInBox(): JSX.Element {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [usernameError, setUsernameError] = useState(false);
  const [usernameErrorMessage, setUsernameErrorMessage] = useState('');
  const [passwordError, setPasswordError] = useState(false);
  const [passwordErrorMessage, setPasswordErrorMessage] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showPasswordMap, setShowPasswordMap] = useState<Record<string, boolean>>({});

  const togglePasswordVisibility = () => {
    setShowPassword((prev) => !prev);
  };

  const togglePasswordVisibilityForField = (identifier: string): void => {
    setShowPasswordMap((prev: Record<string, boolean>) => ({
      ...prev,
      [identifier]: !prev[identifier],
    }));
  };

  const validateInputs = () => {
    let isValid = true;

    if (!username || username.trim().length === 0) {
      setUsernameError(true);
      setUsernameErrorMessage('Username is required.');
      isValid = false;
    } else {
      setUsernameError(false);
      setUsernameErrorMessage('');
    }

    if (!password) {
      setPasswordError(true);
      setPasswordErrorMessage('Password is required.');
      isValid = false;
    } else {
      setPasswordError(false);
      setPasswordErrorMessage('');
    }

    return isValid;
  };

  return (
    <Stack gap={5}>
      <ColorSchemeImage
        src={{
          light: `${import.meta.env.BASE_URL}/assets/images/logo.svg`,
          dark: `${import.meta.env.BASE_URL}/assets/images/logo-inverted.svg`,
        }}
        alt={{light: 'Logo (Light)', dark: 'Logo (Dark)'}}
        height={30}
        width="auto"
        sx={{
          display: {xs: 'flex', md: 'none'},
        }}
      />
      <StyledPaper variant="outlined">
        <SignIn
          onError={(error) => {
            setUsernameError(true);
            setUsernameErrorMessage(error?.message || 'Authentication failed');
          }}
        >
          {({onSubmit, isLoading, components, error, isInitialized}) => (
            <>
              <Typography variant="h2" sx={{width: '100%', mb: 2}}>
                Sign in
              </Typography>

              {(isLoading || !isInitialized) ? (
                <Box sx={{display: 'flex', justifyContent: 'center', p: 3}}>
                  <Typography>Loading authentication...</Typography>
                </Box>
              ) : (
                <>
                  {error && (
                    <Alert severity="error" sx={{mb: 2}}>
                      <AlertTitle>Error</AlertTitle>
                      {error.message || error.toString()}
                    </Alert>
                  )}

                  {/* Handle different flow types */}
                  {(() => {
                    // Check if we have button components (multi-option flow)
                    if (components && components.length > 0 && components.some((c) => c.type === 'BUTTON')) {
                      return (
                        <>
                          {/* Always render BasicAuth form first if basic_auth is available */}
                          {components.some((c) => c.config?.actionId === 'basic_auth') && (
                            <Box
                              component="form"
                              onSubmit={(event) => {
                                event.preventDefault();
                                if (validateInputs()) {
                                  // Tracker: https://github.com/asgardeo/javascript/issues/222
                                  // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-explicit-any
                                  onSubmit({
                                    actionId: 'basic_auth',
                                    inputs: {
                                      username,
                                      password,
                                    },
                                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                                  } as any).catch(() => {
                                    // Error handled by onError callback
                                  });
                                }
                              }}
                              noValidate
                              sx={{display: 'flex', flexDirection: 'column', width: '100%', gap: 2}}
                            >
                              <FormControl>
                                <FormLabel htmlFor="username">Username</FormLabel>
                                <TextField
                                  error={usernameError}
                                  helperText={usernameErrorMessage}
                                  id="username"
                                  type="text"
                                  name="username"
                                  placeholder="Enter your username"
                                  autoComplete="username"
                                  autoFocus
                                  required
                                  fullWidth
                                  variant="outlined"
                                  color={usernameError ? 'error' : 'primary'}
                                  disabled={isLoading}
                                  value={username}
                                  onChange={(e) => setUsername(e.target.value)}
                                />
                              </FormControl>
                              <FormControl>
                                <FormLabel htmlFor="password">Password</FormLabel>
                                <TextField
                                  error={passwordError}
                                  helperText={passwordErrorMessage}
                                  name="password"
                                  placeholder="••••••"
                                  type={showPassword ? 'text' : 'password'}
                                  id="password"
                                  autoComplete="current-password"
                                  required
                                  fullWidth
                                  variant="outlined"
                                  color={passwordError ? 'error' : 'primary'}
                                  disabled={isLoading}
                                  value={password}
                                  onChange={(e) => setPassword(e.target.value)}
                                  slotProps={{
                                    input: {
                                      endAdornment: (
                                        <InputAdornment position="end">
                                          <IconButton
                                            aria-label="toggle password visibility"
                                            onClick={togglePasswordVisibility}
                                            edge="end"
                                            disabled={isLoading}
                                          >
                                            {showPassword ? <EyeClosed /> : <Eye />}
                                          </IconButton>
                                        </InputAdornment>
                                      ),
                                    },
                                  }}
                                />
                              </FormControl>
                              <FormControlLabel
                                control={<Checkbox value="remember" color="primary" disabled={isLoading} />}
                                label="Remember me"
                              />
                              <Button type="submit" fullWidth variant="contained" disabled={isLoading} sx={{mt: 2}}>
                                {isLoading ? 'Signing in...' : 'Sign in'}
                              </Button>
                            </Box>
                          )}

                          {/* Show divider if there are other auth options besides basic_auth */}
                          {components.some((c) => c.type === 'BUTTON' && c.config?.actionId !== 'basic_auth') && (
                            <Divider sx={{my: 2}}>OR</Divider>
                          )}

                          {/* Show other authentication options as buttons */}
                          <Box sx={{display: 'flex', flexDirection: 'column', gap: 2}}>
                            {components
                              .filter((c) => c.type === 'BUTTON')
                              .filter((c) => c.config?.actionId !== 'basic_auth')
                              .map((button) => {
                                const actionId = (button.config?.actionId as string) ?? '';

                                const getIcon = () => {
                                  if (actionId.includes('google')) return <Google />;
                                  if (actionId.includes('facebook')) return <Facebook />;
                                  if (actionId.includes('github')) return <GitHub />;
                                  if (actionId.includes('mobile')) return <Smartphone />;
                                  return null;
                                };

                                const getLabel = () => {
                                  if (button.config?.text && typeof button.config.text === 'string') {
                                    return button.config.text;
                                  }
                                  if (actionId.includes('google')) return 'Sign in with Google';
                                  if (actionId.includes('github')) return 'Sign in with GitHub';
                                  if (actionId.includes('facebook')) return 'Sign in with Facebook';
                                  if (actionId.includes('mobile')) return 'Sign in with SMS OTP';
                                  return `Sign in with ${actionId.replace('_auth', '').replace('_', ' ')}`;
                                };

                                return (
                                  <Button
                                    key={button.id}
                                    fullWidth
                                    variant="outlined"
                                    onClick={() => {
                                      // Tracker: https://github.com/asgardeo/javascript/issues/222
                                      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-explicit-any
                                      onSubmit({actionId} as any).catch(() => {
                                        // Error handled by onError callback
                                      });
                                    }}
                                    disabled={isLoading}
                                    startIcon={getIcon()}
                                  >
                                    {getLabel()}
                                  </Button>
                                );
                              })}
                          </Box>
                        </>
                      );
                    }

                    // Handle FORM components
                    if (components && components.length > 0 && components.some((c) => c.type === 'FORM')) {
                      return (
                        <Box sx={{display: 'flex', flexDirection: 'column', gap: 2}}>
                          {components
                            .filter((c) => c.type === 'FORM')
                            .map((form) => (
                              <Box
                                key={form.id}
                                component="form"
                                onSubmit={(event) => {
                                  event.preventDefault();
                                  const data = new FormData(event.currentTarget);
                                  const inputs: Record<string, string> = {};

                                  // Extract inputs from form components
                                  form.components
                                    ?.filter((c) => c.type === 'INPUT')
                                    .forEach((input) => {
                                      if (input.config?.identifier) {
                                        inputs[input.config.identifier as string] = data.get(
                                          input.config.identifier as string,
                                        ) as string;
                                      }
                                    });

                                  // Tracker: https://github.com/asgardeo/javascript/issues/222
                                  // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-explicit-any
                                  onSubmit({inputs} as any).catch(() => {
                                    // Error handled by onError callback
                                  });
                                }}
                                noValidate
                                sx={{display: 'flex', flexDirection: 'column', width: '100%', gap: 2}}
                              >
                                {form.components?.map((component) => {
                                  if (component.type === 'INPUT') {
                                    const isPasswordField: boolean = (component.config?.type as string) === 'password';
                                    const identifier: string = component.config?.identifier as string;
                                    const showPasswordForField: boolean = showPasswordMap[identifier] || false;

                                    let inputType: string =
                                      (component.config?.type as string) === 'text'
                                        ? 'text'
                                        : (component.config?.type as string);

                                    if (isPasswordField) {
                                      inputType = showPasswordForField ? 'text' : 'password';
                                    }

                                    return (
                                      <FormControl key={component.id}>
                                        <FormLabel htmlFor={identifier}>{component.config?.label as string}</FormLabel>
                                        <TextField
                                          id={identifier}
                                          name={identifier}
                                          type={inputType}
                                          placeholder={component.config?.placeholder as string}
                                          required={component.config?.required as boolean}
                                          fullWidth
                                          variant="outlined"
                                          disabled={isLoading}
                                          slotProps={
                                            isPasswordField
                                              ? {
                                                  input: {
                                                    endAdornment: (
                                                      <InputAdornment position="end">
                                                        <IconButton
                                                          aria-label="toggle password visibility"
                                                          onClick={(): void =>
                                                            togglePasswordVisibilityForField(identifier)
                                                          }
                                                          edge="end"
                                                          disabled={isLoading}
                                                        >
                                                          {showPasswordForField ? <EyeClosed /> : <Eye />}
                                                        </IconButton>
                                                      </InputAdornment>
                                                    ),
                                                  },
                                                }
                                              : undefined
                                          }
                                        />
                                        {(component.config?.hint as string) && (
                                          <Typography variant="caption" color="text.secondary">
                                            {component.config.hint as string}
                                          </Typography>
                                        )}
                                      </FormControl>
                                    );
                                  }

                                  if (component.type === 'BUTTON' && component.config?.type === 'submit') {
                                    return (
                                      <Button
                                        key={component.id}
                                        type="submit"
                                        fullWidth
                                        variant={component.variant === 'PRIMARY' ? 'contained' : 'outlined'}
                                        disabled={isLoading}
                                        sx={{mt: 2}}
                                      >
                                        {isLoading ? 'Submitting...' : (component.config?.text as string)}
                                      </Button>
                                    );
                                  }
                                  return null;
                                })}
                              </Box>
                            ))}
                        </Box>
                      );
                    }

                    // If no components send empty
                    return (
                      <Box sx={{display: 'flex', justifyContent: 'center', p: 3}}>
                        <Typography>Loading authentication...</Typography>
                      </Box>
                    );
                  })()}
                </>
              )}
            </>
          )}
        </SignIn>
          
        <SignUp>
          {({components}) => {
            if (components && components.length > 0) {
              return (
                <Typography sx={{textAlign: 'center', mt: 2}}>
                  Don&apos;t have an account?{' '}
                  <Button
                    variant="text"
                    onClick={() => {
                      const currentParams = searchParams.toString();
                      const createUrl = currentParams ? `${ROUTES.AUTH.SIGN_UP}?${currentParams}` : ROUTES.AUTH.SIGN_UP;
                      // eslint-disable-next-line @typescript-eslint/no-floating-promises
                      navigate(createUrl);
                    }}
                    sx={{
                      p: 0,
                      minWidth: 'auto',
                      textTransform: 'none',
                      color: 'primary.main',
                      textDecoration: 'underline',
                      '&:hover': {
                        textDecoration: 'underline',
                        backgroundColor: 'transparent',
                      },
                    }}
                  >
                    Sign up
                  </Button>
                </Typography>
              );
            }

            return null;
          }}
        </SignUp>
      </StyledPaper>
    </Stack>
  );
}
