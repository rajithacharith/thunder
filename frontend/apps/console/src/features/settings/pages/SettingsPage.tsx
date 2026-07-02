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

import {Box, PageContent, PageTitle, Tab, Tabs} from '@wso2/oxygen-ui';
import {useState, type JSX, type SyntheticEvent} from 'react';
import {useTranslation} from 'react-i18next';
import CorsSection from '../components/cors/CorsSection';

export default function SettingsPage(): JSX.Element {
  const {t} = useTranslation();
  // Controlled Tabs; only the CORS tab exists for now.
  const [activeTab, setActiveTab] = useState(0);

  const handleTabChange = (_event: SyntheticEvent, newValue: number): void => {
    setActiveTab(newValue);
  };

  return (
    <PageContent>
      <PageTitle>
        <PageTitle.Header>{t('settings:page.title')}</PageTitle.Header>
        <PageTitle.SubHeader>{t('settings:page.subtitle')}</PageTitle.SubHeader>
      </PageTitle>

      <Tabs value={activeTab} onChange={handleTabChange} aria-label={t('settings:tabs.ariaLabel')}>
        <Tab
          label={t('settings:tabs.cors')}
          id="settings-tab-0"
          aria-controls="settings-tabpanel-0"
          sx={{textTransform: 'none'}}
        />
      </Tabs>

      <Box role="tabpanel" id="settings-tabpanel-0" aria-labelledby="settings-tab-0" sx={{py: 3}}>
        <CorsSection />
      </Box>
    </PageContent>
  );
}
