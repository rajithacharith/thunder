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

const productConfig = require('../../docusaurus.product.config');

const PRODUCT_NAME = productConfig.project.name;
const CONSOLE_URL = productConfig.local.consoleUrl;
const WAYFINDER_SAMPLE_URL = productConfig.local.samples.wayfinderUrl;
const WAYFINDER_MAIL_URL = productConfig.local.samples.wayfinderMailUrl;

/**
 * Clean an MDX source file into plain Markdown suitable for LLM consumption.
 * Removes imports, replaces known JSX components, and strips leftover JSX.
 */
async function processMarkdownFile(content, _constants, _sourceDir, linkContext) {
  let result = content;

  result = removeImports(result);
  const { frontmatter, body } = extractFrontmatter(result);
  result = body;
  result = replaceThunderIDComponents(result);
  result = processAdmonitions(result);
  result = stripRemainingJsx(result);
  result = rewriteRelativeDocLinks(result, linkContext);
  result = cleanupWhitespace(result);

  if (frontmatter.title) {
    result = `# ${frontmatter.title}\n\n${result}`;
  }

  return result;
}

function extractFrontmatter(content) {
  const match = content.match(/^---\n([\s\S]*?)\n---\n?/);
  if (!match) return { frontmatter: {}, body: content };

  const frontmatter = {};
  const titleMatch = match[1].match(/^title:\s*(.+)/m);
  if (titleMatch) frontmatter.title = titleMatch[1].trim().replace(/^['"]|['"]$/g, '');

  return { frontmatter, body: content.slice(match[0].length) };
}

function removeImports(content) {
  return content.replace(/^import\s+[\s\S]*?from\s+['"][^'"]+['"];?\s*$/gm, '');
}

/** Replace ThunderID-specific MDX components with readable text equivalents. */
function replaceThunderIDComponents(content) {
  let result = content;

  // Self-closing component tags
  result = result.replace(/<ProductName\s*\/>/g, PRODUCT_NAME);
  result = result.replace(/<ConsoleUrl\s*\/>/g, CONSOLE_URL);
  result = result.replace(/<WayFinderSampleUrl\s*\/>/g, WAYFINDER_SAMPLE_URL);
  result = result.replace(/<WayFinderMailUrl\s*\/>/g, WAYFINDER_MAIL_URL);

  // {{...}} frontmatter-style placeholders that may have leaked into body
  result = result.replace(/\{\{ProductName\}\}/g, PRODUCT_NAME);
  result = result.replace(/\{\{productSlug\}\}/g, PRODUCT_NAME.toLowerCase());
  result = result.replace(/\{\{ConsoleUrl\}\}/g, CONSOLE_URL);
  result = result.replace(/\{\{WayFinderSampleUrl\}\}/g, WAYFINDER_SAMPLE_URL);
  result = result.replace(/\{\{WayFinderMailUrl\}\}/g, WAYFINDER_MAIL_URL);

  // Stepper: keep inner content, remove wrapper
  result = result.replace(/<Stepper[^>]*>([\s\S]*?)<\/Stepper>/g, '$1');

  // Tabs / TabItem: strip wrappers, keep content
  result = result.replace(/<Tabs[^>]*>/g, '');
  result = result.replace(/<\/Tabs>/g, '');
  result = result.replace(/<TabItem[^>]*label="([^"]*)"[^>]*>/g, '\n**$1**\n\n');
  result = result.replace(/<\/TabItem>/g, '');

  // NextSteps / NextStepsCard: convert to bullet list.
  // Attributes may appear in any order and span multiple lines.
  result = result.replace(/<NextStepsCard([\s\S]*?)\/>/g, (_match, attrs) => {
    const get = (name) => { const m = attrs.match(new RegExp(`${name}="([^"]*)"`)); return m ? m[1] : ''; };
    const title = get('title');
    const href = get('href');
    const description = get('description');
    if (!title || !href) return '';
    return description ? `- [${title}](${href}) — ${description}` : `- [${title}](${href})`;
  });
  result = result.replace(/<NextSteps>/g, '');
  result = result.replace(/<\/NextSteps>/g, '');

  // SampleDownload: just remove (binary download, not useful in LLM context)
  result = result.replace(/<SampleDownload[^/]*\/>/g, '');

  // WayfinderCast / WayfinderArchitecture / WayfinderOrganization: descriptive fallback
  result = result.replace(/<WayfinderCast\s*\/>/g, '_[Cast diagram: Wayfinder sample users]_');
  result = result.replace(
    /<WayfinderArchitecture\s*\/>/g,
    '_[Architecture diagram: Wayfinder components]_',
  );
  result = result.replace(
    /<WayfinderOrganization\s*\/>/g,
    '_[Organization diagram: Wayfinder team structure]_',
  );

  // B2C / identity journey components: remove (interactive, not useful as text)
  result = result.replace(/<B2C\w+[^/]*\/>/g, '');
  result = result.replace(/<B2C\w+[^>]*>[\s\S]*?<\/B2C\w+>/g, '');
  result = result.replace(/<UseCaseBranchCards\s*\/>/g, '');

  return result;
}

/** Convert Docusaurus admonitions (:::note, :::tip, etc.) to blockquotes. */
function processAdmonitions(content) {
  return content.replace(
    /:::(\w+)(?:\[([^\]]*)\])?\n([\s\S]*?):::/g,
    (match, type, title, body) => {
      const label = title || type.charAt(0).toUpperCase() + type.slice(1);
      const lines = body
        .trim()
        .split('\n')
        .map((l) => '> ' + l);
      return `> **${label}**\n>\n${lines.join('\n')}\n`;
    },
  );
}

/** Strip any remaining JSX/HTML tags that weren't handled above. */
function stripRemainingJsx(content) {
  let result = content;

  // Self-closing tags
  result = result.replace(/<[A-Z][A-Za-z0-9]*[^>]*\/>/g, '');
  // Opening/closing pairs for known wrappers — keep inner text
  result = result.replace(/<[A-Z][A-Za-z0-9]*[^>]*>([\s\S]*?)<\/[A-Z][A-Za-z0-9]*>/g, '$1');
  // Remove leftover HTML div/span wrappers
  result = result.replace(/<div[^>]*>([\s\S]*?)<\/div>/g, '$1');
  result = result.replace(/<span[^>]*>([\s\S]*?)<\/span>/g, '$1');
  // Strip JSX expression blocks {/* comment */} and simple {expr}
  result = result.replace(/\{\/\*[\s\S]*?\*\/\}/g, '');
  result = result.replace(/\{[^}\n]{0,200}\}/g, '');

  return result;
}

/** Rewrite relative doc links to absolute .md URLs. */
function rewriteRelativeDocLinks(content, linkContext) {
  if (!linkContext?.docUrlPath || !linkContext?.siteUrl) return content;

  // Split on code blocks to avoid rewriting links inside them
  const parts = content.split(/(```[\s\S]*?```)/g);
  for (let i = 0; i < parts.length; i++) {
    if (i % 2 === 1) continue; // skip code blocks
    parts[i] = parts[i].replace(
      /(`[^`\n]*`)|(!?)\[([^\]]+)\]\(([^)\s]+)(\s+"[^"]*")?\)/g,
      (match, codeSpan, bang, text, url, title) => {
        if (codeSpan !== undefined) return codeSpan;
        if (bang === '!') return match; // keep images
        const newUrl = rewriteDocLinkUrl(url, linkContext);
        return `[${text}](${newUrl}${title || ''})`;
      },
    );
  }
  return parts.join('');
}

function rewriteDocLinkUrl(url, { docUrlPath, siteUrl }) {
  if (/^([a-z][a-z0-9+.-]*:|\/\/)/i.test(url)) return url; // absolute
  if (url.startsWith('#')) return url; // anchor

  const hashIdx = url.indexOf('#');
  const queryIdx = url.indexOf('?');
  let cutIdx = url.length;
  if (hashIdx !== -1) cutIdx = Math.min(cutIdx, hashIdx);
  if (queryIdx !== -1) cutIdx = Math.min(cutIdx, queryIdx);

  let pathPart = url.slice(0, cutIdx);
  const tail = url.slice(cutIdx);

  if (!pathPart) return url;

  if (/\.mdx$/i.test(pathPart)) {
    pathPart = pathPart.replace(/\.mdx$/i, '.md');
  } else if (!/\.md$/i.test(pathPart)) {
    if (/\.[a-z0-9]{1,5}$/i.test(pathPart)) return url; // non-doc file, leave alone
    pathPart = pathPart.replace(/\/$/, '') + '.md';
  }

  let absPath;
  if (pathPart.startsWith('/')) {
    absPath = pathPart;
  } else {
    const baseDir = docUrlPath.replace(/\/[^/]*$/, '');
    absPath = resolvePath(baseDir + '/' + pathPart);
  }

  return siteUrl.replace(/\/$/, '') + absPath + tail;
}

function resolvePath(p) {
  const segments = p.split('/');
  const out = [];
  for (const seg of segments) {
    if (seg === '' || seg === '.') {
      if (out.length === 0 && seg === '') out.push('');
      continue;
    }
    if (seg === '..') {
      if (out.length > 1) out.pop();
      continue;
    }
    out.push(seg);
  }
  return out.join('/') || '/';
}

function cleanupWhitespace(content) {
  let result = content;
  result = result.replace(/\n{4,}/g, '\n\n\n');
  result = result.replace(/[ \t]+$/gm, '');
  result = result.trim() + '\n';
  return result;
}

module.exports = { processMarkdownFile };
