// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

/**
 * Sample App Login Page Object
 *
 * Page Object Model for the Vanilla Sample App login functionality.
 * Provides methods to interact with login form and verify authentication.
 */

import { Page, expect } from "@playwright/test";
import { BasePage } from "../base.page";
import { Timeouts } from "../../constants/timeouts";

export class SampleAppLoginPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  /**
   * Navigate to the sample app
   * @param url - Sample app URL (default: https://localhost:3000)
   */
  async goto(url: string = "https://localhost:3000") {
    await this.page.goto(url, { waitUntil: "commit" });
  }

  /**
   * Verify the login page is loaded
   *
   * Unlike a redirect-based sample, this app renders its login form directly at "/" - there's no
   * separate landing page with a "Sign In" button, so this just waits for the form itself.
   */
  async verifyHomePageLoaded() {
    await this.verifyLoginPageLoaded();
  }

  /**
   * No-op: kept so existing call sites (goto -> verifyHomePageLoaded -> clickSignInButton ->
   * verifyLoginPageLoaded) don't need restructuring. This app has no separate landing page to
   * navigate away from - the login form is already visible after verifyHomePageLoaded.
   */
  async clickSignInButton() {
    // Intentionally empty.
  }

  async verifyLoginPageLoaded() {
    await this.page.waitForSelector('input[name="username"], input[placeholder*="username" i]', {
      timeout: Timeouts.NETWORK_IDLE,
      state: "visible",
    });
  }

  /**
   * Fill in the login form
   * @param username - Username to enter
   * @param password - Password to enter
   */
  async fillLoginForm(username: string, password: string) {
    // Fill username
    const usernameInput = this.page.locator('input[name="username"], input[placeholder*="username" i]').first();
    await usernameInput.waitFor({ state: "visible", timeout: Timeouts.DEFAULT_ACTION });
    await usernameInput.fill(username);

    // Fill password
    const passwordInput = this.page.locator('input[name="password"], input[placeholder*="password" i]').first();
    await passwordInput.waitFor({ state: "visible", timeout: Timeouts.DEFAULT_ACTION });
    await passwordInput.fill(password);
  }

  /**
   * Click the login/sign in button
   */
  async clickLogin() {
    // Try multiple selector strategies for the login button
    const loginButton = this.page
      .locator(
        'button[type="submit"], button:has-text("Sign In"), button:has-text("Login"), button:has-text("Sign in")'
      )
      .first();

    await loginButton.waitFor({ state: "visible", timeout: Timeouts.DEFAULT_ACTION });
    await loginButton.click();
  }

  /**
   * Perform complete login flow
   * @param username - Username to login with
   * @param password - Password to login with
   */
  async login(username: string, password: string) {
    await this.fillLoginForm(username, password);
    await this.clickLogin();

    // Wait for navigation/response after login
    await this.page.waitForLoadState("networkidle");
  }

  /**
   * Verify user is logged in successfully
   * Checks for common indicators like avatar, profile information, or welcome message
   */
  async verifyLoggedIn() {
    // Wait for login to complete - look for logged-in state indicators
    await this.page.waitForLoadState("networkidle");

    // Check for common logged-in indicators (adjust selectors based on your app)
    const loggedInIndicators = [
      this.page.locator('button[aria-haspopup="true"]'), // Avatar menu button
      this.page.locator('button:has(> div[class*="MuiAvatar"])'), // Avatar button
      this.page.locator('[role="menuitem"]:has-text("Sign Out")'), // May be visible if menu is open
      this.page.locator('[data-testid="user-profile"]'),
      this.page.getByText(/welcome|hello/i),
      this.page.locator(".user-profile"),
      this.page.locator(".logged-in"),
      this.page.locator(".token-container"), // Token display container
    ];

    // Require at least one indicator to appear. This must assert rather than probe: a sign-in that
    // fails client-side (for example a blocked cross-origin token read) still leaves the app on a
    // page with no login form, so tolerating a missing indicator would report success for a broken
    // sign-in and push the failure into whatever assertion happens to run next.
    const anyLoggedInIndicator = loggedInIndicators.reduce((combined, locator) => combined.or(locator));
    await expect(anyLoggedInIndicator.first()).toBeVisible({ timeout: Timeouts.DEFAULT_ACTION });

    // Take a screenshot for verification
    await this.screenshot("logged-in-state");
  }

  /**
   * Verify specific user information is displayed
   * @param userInfo - Expected user information (e.g., username, email)
   */
  async verifyUserInfo(userInfo: string) {
    await expect(this.page.locator(`text=${userInfo}`)).toBeVisible({ timeout: Timeouts.DEFAULT_ACTION });
  }

  /**
   * Click logout button
   * The logout option is in a dropdown menu accessed via Avatar button
   */
  async logout() {
    // First, look for the avatar/menu button to open the menu
    const avatarButton = this.page
      .locator(
        'button[aria-haspopup="true"], button[aria-controls="account-menu"], button:has(> div[class*="MuiAvatar"])'
      )
      .first();

    // Fallback: direct logout button (for other app implementations)
    const logoutButton = this.page
      .locator('button:has-text("Logout"), button:has-text("Sign Out"), [data-testid="logout-button"]')
      .first();

    // Wait (with auto-retry) for whichever logout indicator is visible first
    await avatarButton.or(logoutButton).first().waitFor({ state: "visible", timeout: Timeouts.DEFAULT_ACTION });

    if (await avatarButton.isVisible().catch(() => false)) {
      // Click avatar to open menu
      await avatarButton.click();

      // Wait for menu to appear
      await this.page.waitForSelector('#account-menu, [role="menu"]', {
        state: "visible",
        timeout: Timeouts.DEFAULT_ACTION,
      });

      // Click the Sign Out menu item
      const signOutMenuItem = this.page
        .locator('[role="menuitem"]:has-text("Sign Out"), [role="menuitem"]:has-text("Logout")')
        .first();

      await signOutMenuItem.waitFor({ state: "visible", timeout: Timeouts.DEFAULT_ACTION });
      await signOutMenuItem.dispatchEvent("click");
    } else {
      await logoutButton.click();
    }

    // RP-initiated logout redirects to the gate's sign-out confirmation flow before the
    // post-logout redirect back to the app; confirm it if it appears.
    const confirmSignOutButton = this.page.getByRole("button", { name: "Sign out" });
    const confirmed = await confirmSignOutButton
      .waitFor({ state: "visible", timeout: Timeouts.DEFAULT_ACTION })
      .then(() => true)
      .catch(() => false);
    if (confirmed) {
      await confirmSignOutButton.click();
    }

    await this.page.waitForLoadState("networkidle");
  }

  /**
   * Verify logout was successful
   * RP-initiated logout redirects back to the app's post-logout redirect URI (the home page),
   * not the login form.
   */
  async verifyLoggedOut() {
    await this.verifyHomePageLoaded();
  }

  /**
   * Verify OTP input page is loaded
   * Used for MFA authentication flows where OTP verification is required
   * Handles both single input and MUI separate digit inputs
   */
  async verifyOTPPageLoaded() {
    // Wait for either "Verify OTP" heading or OTP input fields
    await this.page.waitForSelector(
      'h3:has-text("Verify OTP"), input[aria-label*="OTP digit" i], input[name="otp"], input[placeholder*="otp" i]',
      {
        timeout: Timeouts.NETWORK_IDLE,
        state: "visible",
      }
    );
  }

  /**
   * Fill in the OTP input field
   * Handles both single input field and MUI separate digit inputs (6 boxes)
   * @param otp - OTP code to enter (e.g., "123456")
   */
  async fillOTP(otp: string) {
    // Check if MUI separate digit inputs exist (aria-label="OTP digit 1", etc.)
    const digitInputs = this.page.locator('input[aria-label*="OTP digit" i]');
    const digitCount = await digitInputs.count();

    if (digitCount > 0) {
      // MUI separate digit inputs - fill each digit individually
      console.log(`  Filling ${digitCount} separate OTP digit inputs...`);
      const digits = otp.split("");

      // Each digit is filled via its own locator, which focuses it as part of Playwright's own
      // actionability checks, so this doesn't depend on the app's auto-advance-on-input JS having
      // already run by the time the next iteration starts.
      for (let i = 0; i < Math.min(digits.length, digitCount); i++) {
        const input = digitInputs.nth(i);
        await input.waitFor({ state: "visible", timeout: Timeouts.DEFAULT_ACTION });
        await input.fill(digits[i]);
      }
    } else {
      // Single OTP input field
      const otpInput = this.page
        .locator('input[name="otp"], input[placeholder*="otp" i], input[placeholder*="code" i]')
        .first();

      await otpInput.waitFor({ state: "visible", timeout: Timeouts.DEFAULT_ACTION });
      await otpInput.fill(otp);
    }
  }

  /**
   * Click the submit/verify OTP button
   * Looks for "Verify" button specifically for MUI OTP form
   */
  async clickVerifyOTP() {
    // MUI OTP form has a "Verify" button (not type="submit")
    const verifyButton = this.page
      .locator(
        'button:has-text("Verify"):not(:has-text("Resend")), button:has-text("Submit"), button:has-text("Continue"), button[type="submit"]'
      )
      .first();

    await verifyButton.waitFor({ state: "visible", timeout: Timeouts.DEFAULT_ACTION });
    await verifyButton.click();
  }

  /**
   * Complete OTP verification step
   * @param otp - OTP code to verify
   */
  async verifyOTP(otp: string) {
    await this.fillOTP(otp);
    await this.clickVerifyOTP();

    // Wait for navigation/response after OTP verification
    await this.page.waitForLoadState("networkidle");
  }

  /**
   * Returns the visible OTP error message text, or an empty string if none is present.
   */
  async getOTPErrorMessage(): Promise<string> {
    const text = await this.page
      .locator(".MuiAlert-message, .MuiAlert-colorError .MuiAlertTitle-root")
      .first()
      .textContent()
      .catch(() => "");
    return text?.trim() ?? "";
  }
}
