/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/**
 * Settings Page Object Model (CORS allowed origins)
 *
 * Encapsulates the Settings > CORS panel: adding/removing custom allowed origins and saving.
 *
 * @example
 * const settingsPage = new SettingsPage(page, baseUrl);
 * await settingsPage.goto();
 * await settingsPage.addAllowedOrigin("https://app.example.com");
 */

import { Page, Locator, expect } from "@playwright/test";
import { ConsoleRoutes } from "../../configs/routes/console-routes";
import { BasePage } from "../base.page";
import { Timeouts } from "../../constants/timeouts";

// Matches the placeholder rendered on each editable (custom) origin input.
const ORIGIN_PLACEHOLDER = "https://app.example.com";

export class SettingsPage extends BasePage {
  readonly baseUrl: string;

  readonly corsTab: Locator;
  readonly addOriginButton: Locator;
  readonly saveButton: Locator;
  readonly discardButton: Locator;
  // Editable (custom) origin inputs and their delete buttons render in the same row order,
  // so the Nth input aligns with the Nth remove button. Read-only rows carry neither.
  readonly originInputs: Locator;
  readonly removeButtons: Locator;

  constructor(page: Page, baseUrl: string) {
    super(page);
    this.baseUrl = baseUrl;

    this.corsTab = page.getByRole("tab", { name: /cors/i });
    this.addOriginButton = page.getByRole("button", { name: /add origin/i });
    this.saveButton = page.getByRole("button", { name: /save changes/i });
    this.discardButton = page.getByRole("button", { name: /discard/i });
    this.originInputs = page.getByPlaceholder(ORIGIN_PLACEHOLDER);
    this.removeButtons = page.getByRole("button", { name: /remove origin/i });
  }

  /** Navigate to the Settings (CORS) page. */
  async goto() {
    await this.page.goto(`${this.baseUrl}${ConsoleRoutes.settings}`, {
      waitUntil: "networkidle",
      timeout: Timeouts.PAGE_LOAD,
    });
    await this.corsTab.first().waitFor({ state: "visible", timeout: Timeouts.ELEMENT_VISIBILITY });
  }

  /** Index of the editable row holding the given origin, or -1 if absent. */
  private async indexOfOrigin(origin: string): Promise<number> {
    const count = await this.originInputs.count();
    for (let i = 0; i < count; i++) {
      if ((await this.originInputs.nth(i).inputValue()) === origin) {
        return i;
      }
    }
    return -1;
  }

  /** Whether a custom (editable) origin with the given value is currently listed. */
  async hasCustomOrigin(origin: string): Promise<boolean> {
    return (await this.indexOfOrigin(origin)) !== -1;
  }

  /** Add a custom allowed origin and persist it. */
  async addAllowedOrigin(origin: string) {
    await this.addOriginButton.click();
    const input = this.originInputs.last();
    await input.fill(origin);
    await input.blur();
    await this.save();
  }

  /** Remove a custom allowed origin (no-op if absent) and persist. */
  async removeAllowedOrigin(origin: string) {
    const index = await this.indexOfOrigin(origin);
    if (index === -1) {
      return;
    }
    await this.removeButtons.nth(index).click();
    await this.save();
  }

  /** Click Save changes and wait for the unsaved-changes bar to clear (success). */
  private async save() {
    await expect(this.saveButton).toBeEnabled({ timeout: Timeouts.ELEMENT_VISIBILITY });
    await this.saveButton.click();
    await expect(this.saveButton).toBeHidden({ timeout: Timeouts.ELEMENT_VISIBILITY });
  }
}
