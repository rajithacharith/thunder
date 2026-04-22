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

import {FormControl, FormLabel, MenuItem, Select} from '@wso2/oxygen-ui';
import isEmpty from 'lodash-es/isEmpty';
import {memo, useCallback, useMemo, type ReactElement} from 'react';
import {useTranslation} from 'react-i18next';
import ButtonExtendedProperties from './extended-properties/ButtonExtendedProperties';
import ExecutionExtendedProperties from './extended-properties/ExecutionExtendedProperties';
import FieldExtendedProperties from './extended-properties/FieldExtendedProperties';
import RulesProperties from './nodes/RulesProperties';
import ResourcePropertyFactory from './ResourcePropertyFactory';
import TextPropertyField from '@/features/flows/components/resource-property-panel/TextPropertyField';
import type {ResourcePropertiesProps} from '@/features/flows/context/FlowBuilderCoreProvider';
import type {FieldKey, FieldValue} from '@/features/flows/models/base';
import {ElementCategories, ElementTypes, type Element} from '@/features/flows/models/elements';
import type {Resource} from '@/features/flows/models/resources';
import {ExecutionTypes, StepCategories, StepTypes} from '@/features/flows/models/steps';

/**
 * Factory to generate the property configurator for the given password recovery flow resource.
 *
 * @param props - Props injected to the component.
 * @returns The ResourceProperties component.
 */
function ResourceProperties({
  properties,
  resource,
  onChange,
  onVariantChange,
}: ResourcePropertiesProps): ReactElement | null {
  const {t} = useTranslation();
  const handleChange = useCallback(
    (propertyKey: string, newValue: unknown, changedResource: unknown): void => {
      let processedValue: string | boolean | object;
      if (typeof newValue === 'boolean') {
        processedValue = newValue;
      } else if (typeof newValue === 'object' && newValue !== null) {
        processedValue = newValue;
      } else if (typeof newValue === 'string' || typeof newValue === 'number') {
        processedValue = String(newValue);
      } else {
        processedValue = '';
      }

      onChange(propertyKey, processedValue, changedResource as Resource);
    },
    [onChange],
  );
  const selectedVariant = useMemo<Element | undefined>(() => {
    if (!resource?.variants || resource.variants.length === 0) {
      return undefined;
    }
    const currentVariant = resource.variants.find((v: Element) => v.variant === (resource as Element).variant) as
      | Element
      | undefined;
    return currentVariant ?? (resource.variants[0] as Element);
  }, [resource]);

  const renderElementId = (): ReactElement => (
    <ResourcePropertyFactory
      key={`${resource.id}-$id`}
      resource={resource}
      propertyKey="id"
      propertyValue={resource.id}
      onChange={handleChange}
    />
  );

  const renderElementPropertyFactory = () => {
    const hasVariants = !isEmpty(resource?.variants);

    return (
      <>
        {hasVariants && (
          <div>
            <FormLabel htmlFor="variant-select">Variant</FormLabel>
            <Select
              id="variant-select"
              value={selectedVariant?.variant ?? ''}
              onChange={(e) => {
                const newVariant = resource?.variants?.find((variant: Element) => variant.variant === e.target.value);
                onVariantChange?.((newVariant?.variant as string) ?? '');
              }}
              fullWidth
            >
              {resource?.variants?.map((variant: Element) => (
                <MenuItem key={variant.variant as string} value={variant.variant as string}>
                  {variant.variant as string}
                </MenuItem>
              ))}
            </Select>
          </div>
        )}
        {properties &&
          Object.entries(properties)?.map(([key, value]: [FieldKey, FieldValue]) => (
            <ResourcePropertyFactory
              key={`${resource.id}-${key}`}
              resource={resource}
              propertyKey={key}
              propertyValue={value}
              data-componentid={`${resource.id}-${key}`}
              onChange={handleChange}
            />
          ))}
      </>
    );
  };

  switch (resource.category) {
    case StepCategories.Interface:
      if (resource.type === StepTypes.End) {
        return (
          <>
            {renderElementId()}
            {/* <FlowCompletionProperties resource={resource} onChange={onChange} /> */}
          </>
        );
      }

      return null;
    case ElementCategories.Field:
      return (
        <>
          {renderElementId()}
          <FieldExtendedProperties resource={resource} onChange={handleChange} />
          {renderElementPropertyFactory()}
        </>
      );
    case ElementCategories.Action:
      return (
        <>
          {renderElementId()}
          {resource.type === ElementTypes.Action && (
            <ButtonExtendedProperties resource={resource} onChange={handleChange} onVariantChange={onVariantChange} />
          )}
          {renderElementPropertyFactory()}
        </>
      );
    case StepCategories.Decision:
      if (resource.type === StepTypes.Rule) {
        return (
          <>
            {renderElementId()}
            <RulesProperties />
          </>
        );
      }

      return null;
    case StepCategories.Workflow:
      if (
        resource.type === StepTypes.Execution &&
        (resource?.data as {action?: {executor?: {name?: string}}})?.action?.executor?.name ===
          ExecutionTypes.ConfirmationCode
      ) {
        return (
          <>
            {renderElementId()}
            {/* <ConfirmationCodeProperties resource={resource} onChange={onChange} /> */}
            {renderElementPropertyFactory()}
          </>
        );
      }
      return (
        <>
          {renderElementId()}
          <ExecutionExtendedProperties resource={resource} onChange={handleChange} />
          {renderElementPropertyFactory()}
        </>
      );
    case ElementCategories.Display:
      if (resource.type === ElementTypes.Text) {
        const hasVariants = !isEmpty(resource?.variants);

        return (
          <>
            {renderElementId()}
            {hasVariants && (
              <div>
                <FormLabel htmlFor="variant-select">Variant</FormLabel>
                <Select
                  id="variant-select"
                  value={selectedVariant?.variant ?? ''}
                  onChange={(e) => {
                    const newVariant = resource?.variants?.find(
                      (variant: Element) => variant.variant === e.target.value,
                    );
                    onVariantChange?.((newVariant?.variant as string) ?? '');
                  }}
                  fullWidth
                >
                  {resource?.variants?.map((variant: Element) => (
                    <MenuItem key={variant.variant as string} value={variant.variant as string}>
                      {variant.variant as string}
                    </MenuItem>
                  ))}
                </Select>
              </div>
            )}
            <TextPropertyField
              resource={resource}
              propertyKey="label"
              propertyValue={(resource as Element & {label?: string}).label ?? ''}
              onChange={(_key, value, res) => handleChange('label', value, res)}
            />
            <FormControl fullWidth size="small">
              <FormLabel htmlFor="align-select">{t('flows:core.elements.text.align.label')}</FormLabel>
              <Select
                id="align-select"
                value={(resource as Element & {align?: string}).align ?? 'left'}
                onChange={(e) => handleChange('align', e.target.value, resource)}
              >
                {(['left', 'center', 'right', 'justify', 'inherit'] as const).map((opt) => (
                  <MenuItem key={opt} value={opt}>
                    {t(`flows:core.elements.text.align.options.${opt}`)}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </>
        );
      }
      if (resource.type === ElementTypes.Image) {
        return (
          <>
            {renderElementId()}
            <TextPropertyField
              resource={resource}
              propertyKey="src"
              propertyValue={(resource as Element & {src?: string}).src ?? ''}
              onChange={(_key, value, res) => handleChange('src', value, res)}
            />
            <TextPropertyField
              resource={resource}
              propertyKey="alt"
              propertyValue={(resource as Element & {alt?: string}).alt ?? ''}
              onChange={(_key, value, res) => handleChange('alt', value, res)}
            />
            <TextPropertyField
              resource={resource}
              propertyKey="width"
              propertyValue={(resource as Element & {width?: string}).width ?? ''}
              onChange={(_key, value, res) => handleChange('width', value, res)}
            />
            <TextPropertyField
              resource={resource}
              propertyKey="height"
              propertyValue={(resource as Element & {height?: string}).height ?? ''}
              onChange={(_key, value, res) => handleChange('height', value, res)}
            />
          </>
        );
      }
      return (
        <>
          {renderElementId()}
          {renderElementPropertyFactory()}
        </>
      );

    default:
      return (
        <>
          {renderElementId()}
          {renderElementPropertyFactory()}
        </>
      );
  }
}

export default memo(ResourceProperties);
