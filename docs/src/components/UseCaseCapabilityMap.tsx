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

import React from 'react';
import './UseCaseCapabilityMap.css';

export interface UseCaseMapNode {
  id: string;
  href: string;
  label: string;
  icon: React.ReactNode;
}

export interface UseCaseMapGroup {
  id: string;
  label: string;
  nodes: UseCaseMapNode[];
}

interface UseCaseCapabilityMapProps {
  ariaLabel: string;
  root: UseCaseMapNode;
  groups: UseCaseMapGroup[];
}

export function UseCaseCapabilityMap({ ariaLabel, root, groups }: UseCaseCapabilityMapProps) {
  return (
    <nav className="uc-capability-map" aria-label={ariaLabel}>
      <div className="uc-capability-map__canvas">
        <svg className="uc-capability-map__path" viewBox="0 0 1000 700" preserveAspectRatio="none" aria-hidden="true">
          <path d="M500 92 V112" />
          <path d="M500 112 L116 112" />
          <path d="M500 112 L372 112" />
          <path d="M500 112 L628 112" />
          <path d="M500 112 L884 112" />
          <path d="M116 112 V124" />
          <path d="M372 112 V124" />
          <path d="M628 112 V124" />
          <path d="M884 112 V124" />
        </svg>

        <a
          href={root.href}
          className="uc-capability-map__node uc-capability-map__node--root"
        >
          <span className="uc-capability-map__icon" aria-hidden>{root.icon}</span>
          <span className="uc-capability-map__label">{root.label}</span>
        </a>

        {groups.map((group, groupIndex) => (
          <div
            key={group.id}
            className="uc-capability-map__group"
            style={{ gridColumn: groupIndex + 1, gridRow: 2 }}
          >
            <div
              className="uc-capability-map__category"
              style={{ gridColumn: groupIndex + 1, gridRow: 2 }}
            >
              <span className="uc-capability-map__category-label">
                {group.label}
              </span>
            </div>
            {group.nodes.map((node, nodeIndex) => (
              <a
                key={node.id}
                href={node.href}
                className="uc-capability-map__node"
                style={{ gridColumn: groupIndex + 1, gridRow: nodeIndex + 3 }}
              >
                <span className="uc-capability-map__icon" aria-hidden>{node.icon}</span>
                <span className="uc-capability-map__label">{node.label}</span>
              </a>
            ))}
          </div>
        ))}
      </div>
    </nav>
  );
}
