---
title: AGENTS
description: AI agents should use this file when creating and reviewing documentation content for ThunderID. It contains the standards, guidelines, and requirements that must be followed to ensure high-quality documentation.
---

# ThunderID Documentation Creation Instructions

Follow these instructions when creating new documentation content for ThunderID. Adhere to all guidelines to ensure consistency, clarity, and quality.

## Scope and Boundaries


### Audience

- Your primary audience is ThunderID users, including system administrators, developers, and IT professionals.
- Assume the audience has a basic understanding of IT concepts but may be unfamiliar with ThunderID specifics.
- Avoid jargon and explain concepts clearly.

### What You Must Do

- Strictly adhere to the authoring standards outlined below.
- Choose the appropriate navigation location for the new content based on its topic and relevance.
- Create content that is clear, concise, and actionable for the intended audience.
- Ensure all technical details are accurate and up-to-date.
- Use the provided templates and formatting rules consistently.

## Authoring Standards

You must follow these standards when creating documentation content for ThunderID.

### Voice and Tone

- Use active voice and present tense. You can only use passive voice when the actor is unknown or unimportant.
- Use plain language and short sentences.
- Address the reader as “you.”
- Keep a professional, friendly, neutral tone.
- Avoid slang, jokes, sarcasm, and marketing language.

## Terminology and Consistency

### 1. Product and Feature Names

- Use official product and feature names exactly as defined.
- Do not invent shorthand names.
- Do not change capitalization.
- Do not alternate between long and short forms unless formally introduced.
- Never use the bare word `thunder`, `Thunder`, or `THUNDER` as a short form of the product name. The accepted forms are `ThunderID`, `thunderid`, and `THUNDERID` only.
- **PRs that introduce bare `thunder`/`Thunder`/`THUNDER` (where `thunder` is not immediately followed by `id`/`Id`/`ID`) must not be merged until corrected.**

**Correct:**
- ThunderID Console  
- Passkey Authentication  

**Incorrect:**
- ThunderID console  
- Console  
- Passkey auth
- Thunder (as a short form of ThunderID)

### 2. Acronyms and Abbreviations

- Define acronyms on first use unless universally known (API, URL, JSON, HTTP).
- After definition, use the acronym consistently.
- Do not redefine an acronym within the same document.
- Do not mix expanded and abbreviated forms randomly.

**Correct:**

> Multi-Factor Authentication (MFA)  
> Enable MFA for the application.

**Incorrect:**

> Multi-Factor Authentication (MFA)  
> Enable multi-factor authentication for the application.

### 3. Term Consistency

- Use one term per concept.
- Do not switch terminology mid-document.
- If two terms are synonymous, choose one and use it consistently.

**Incorrect examples:**
- application / app  
- organization / tenant  
- sign in / login (unless intentionally differentiated)

Consistency overrides preference.

### 4. Use Standard Technical Terminology

- Prefer established technical terms.
- Avoid inventing alternative phrases for common concepts.

**Prefer:**
- server  
- endpoint  
- token  
- request  
- response  
- database  
- session  

**Avoid:**
- backend machine  
- link point  
- data store system (unless specific)

### 5. Avoid Ambiguous Pronouns

- Avoid “it,” “this,” “that,” or “they” if the referent is unclear.
- Replace pronouns with explicit nouns when ambiguity exists.

**Ambiguous:**

> Configure the server and restart it.

**Clear:**

> Configure the server and restart the server.

### 6. Avoid Weak “Be” Verb Constructions

Reduce unnecessary use of:
- am  
- is  
- are  
- was  
- were  

Prefer direct verbs.

**Instead of:**

> The configuration is located in `deployment.toml`.

**Write:**

> The configuration file is `deployment.toml`.  
> Or:  
> Find the configuration in `deployment.toml`.

**Instead of:**

> The token is used to authenticate requests.

**Write:**

> The token authenticates requests.

Use “is” only when it improves clarity.

### 7. Prefer Concrete Language

- Use precise nouns and strong verbs.
- Avoid vague verbs such as:
  - handle  
  - manage  
  - deal with  
  - perform  
  - utilize  

**Instead of:**

> The system handles authentication.

**Write:**

> The system validates credentials and issues tokens.

### 8. Formal Language Policy

Avoid informal shorthand in prose:

- config → configuration  
- dev → development  
- prod → production  
- env → environment  
- repo → repository  

These are allowed only inside code blocks, file paths, commands, or environment variable names.

### Heading Capitalization Rules

- Use **Title Case** for all headings (document titles).
  - Capitalize all major words (nouns, verbs, adjectives, adverbs, and pronouns).
  - Do not capitalize short articles, coordinating conjunctions, or prepositions unless they are the first or last word.

  Example:
  `# Configure Passwordless Authentication`

- Use consistent heading levels to reflect document structure.
- Make headings task-focused and descriptive. Do not use generic headings like “Introduction” or “Details.”

### Lists

- Use numbered lists for procedures and ordered steps.
- Use bulleted lists for non-sequential information.
- Keep list items parallel in grammar and structure.

### Formatting Rules

- UI labels, buttons, menu items: use **bold**.
  - Example: Select **Save**.

- Code elements, file names, paths, config keys, commands, URLs: use backticks.
  - Example: Update `deployment.toml`.

- Use descriptive link text. Do not paste raw URLs as link text.
  - Example: `[Microsoft Writing Style Guide](https://learn.microsoft.com/en-us/style-guide/welcome/)`

### Code Blocks

- Use fenced code blocks with a language tag when known.
- Keep code blocks focused.
- Do not include secrets, tokens, passwords, or realistic keys.

    Example:

    ```toml
    [server]
    hostname = "localhost"
    ```

### Configuration Guidance

When documenting configuration:

- Describe what the setting controls.
- State the default value.
- State constraints (type, valid range, allowed values).
- Provide a minimal example.
- Explain when the user should change it.

### Links and References

#### Internal Links

- Use descriptive link text.
- Prefer linking to canonical pages (overview or primary reference).
- Avoid linking to unstable or temporary resources.

#### External Links

- Use external links sparingly and only when they add clear value.
- Use descriptive link text.
- Prefer authoritative sources (official documentation or standards).

### Images and Screenshots

- Do not add, generate, or request new images or screenshots as part of documentation creation.
- Do not reference an image unless the user explicitly confirms it exists and is accessible.
- Do not make images required to complete a task. Provide text alternatives.


## Component Development Standards

Follow these standards when creating or modifying `.tsx` and `.ts` files under `docs/src/`.

### Product Name

Do not hardcode `ThunderID` or any product name string. Always derive it from Docusaurus site config:

```ts
const {siteConfig} = useDocusaurusContext();
const {project} = siteConfig.customFields?.product as DocusaurusProductConfig;
const productName = project.name;
```

In JSX output, prefer the `<ProductName />` component (import from `@site/src/components/ProductName`). It is also globally registered as an MDX component — no import needed in `.mdx` files.

### UI Components

Use `@wso2/oxygen-ui` components instead of raw HTML elements. The package is an MUI wrapper — prefer its components (`Box`, `Typography`, `Button`, `Chip`, `Card`, etc.) over native `<div>`, `<span>`, `<p>`, `<button>`, `<a>`, and similar elements.

Do not use native HTML elements for layout or structure when an equivalent Oxygen UI component exists.

### Icons

Use `@wso2/oxygen-ui-icons-react` (a Lucide React wrapper) for all icons. Do not import from other icon libraries (`lucide-react`, `react-icons`, `@heroicons/react`, etc.).

Exception: custom brand or technology logos (e.g., `AndroidLogo`, `ReactLogo`) that do not exist in the Oxygen UI icon set may remain as custom SVG components under `docs/src/components/icons/`.

### Styling

Use the `sx` prop for component styling. Use `styled()` from `@wso2/oxygen-ui` only when `sx` alone is insufficient (for example, complex descendant selectors or keyframe animations).

Do not use inline `style={{...}}` props. Do not add CSS class-based styles via `className` unless you are targeting Docusaurus-controlled elements where `sx` is not applicable.

### Diagrams

Use Mermaid for architecture diagrams, flow diagrams, sequence diagrams, and similar visuals. Do not hand-build diagrams out of raw SVG elements (`<svg>`, `<rect>`, `<path>`, `<text>`, etc.) or ASCII art.

Only fall back to raw SVG when the diagram requires a layout Mermaid cannot express.

### custom.css

Do not add styles to `docs/src/css/custom.css`. That file is reserved for:

- Infima CSS variable overrides.
- Docusaurus structural adjustments with no Oxygen UI hook.
- Third-party overrides outside our control (Scalar API reference, etc.).

If a style can be expressed via `sx` or `styled()`, place it there instead.

## Documentation Structure Requirements

All task-based documentation must follow a logical, goal-oriented structure that guides the reader from start to finish. This should only apply to Guides and Tutorials. Community, Reference and API documentation may follow a different structure as appropriate.

The document must clearly communicate:

- What the reader will achieve.
- When the task is applicable.
- What prerequisites are required.
- How to complete the task (clear, sequential steps).
- How to confirm the outcome.
- How to troubleshoot common issues (if applicable).
- What to do next (related tasks or follow-up actions).

Each section must build on the previous one and move the reader toward successful task completion.

Avoid:
- Unnecessary background information.
- Repetition.
- Conceptual digressions unrelated to the task.
- Sections with no actionable value.

## Quality Checklist (Must Pass)

Before finalizing output, ensure:

- Headings are title case.
- Procedures use numbered lists.
- UI labels are **bold**.
- Code elements and paths are in backticks.
- Links use descriptive text.
- Content is concise, active voice, present tense.
- No unverified claims or placeholders remain.
- No secrets or sensitive data appear in examples.
- After creating content, run Vale locally and resolve all warnings.
- If Vale flags a word as a spelling error, check whether it is a legitimate product term, technical term, or widely accepted term. If yes, add it to `.vale/styles/config/vocabularies/vocab/accept.txt`. If not, fix the spelling instead.

## Output Requirements

- Output must be Markdown.
- Use a single top-level title (`# ...`).
- Use consistent section ordering and headings.
- If assumptions are made, include an **Assumptions** section near the top.
- End with a **Next steps** section when appropriate.

### Vale Verification Requirement

Before finalizing documentation output:

- If Vale output is provided, resolve all reported errors and warnings before finalizing.
- If Vale output is not available, remind the user to run Vale locally.

### CI Feedback Handling

When Vale feedback is provided through CI checks:

- Only respond to the **latest** Vale check run.
- Ignore resolved or outdated annotations from previous commits.
- Do NOT repeat or expand on previously addressed Vale findings.
- If the latest CI run is clean, do not comment on earlier issues.

## Vocabulary Guidelines

Strictly follow these vocabulary guidelines when writing ThunderID documentation.

### Use of "Multiple"

- Use multiple only when it adds clarity about behavior, constraints, or guarantees.
- Avoid multiple when the plural form already conveys the meaning.
- Use multiple when it expresses a real capability, constraint, or relationship.
  
  - Examples

    - A user can belong to multiple organizations.
    - A policy can include multiple conditions.
    - An application can have multiple identity providers.
    - A tenant may configure multiple authentication methods.

  In these cases, removing multiple would make the statement ambiguous or weaker.

### Use of 'Login' and 'Sign-In'

- Use login and sign-in consistently based on context.
- They are not interchangeable in documentation.

#### Login / Log In — System and Developer Perspective

Use login for system-level and developer-facing terminology, especially when the term is widely known, standardized, or protocol-defined.

Examples:
- social login
- login endpoint
- login_hint
- login URI
- last login time

Avoid using login to describe user-facing flows or actions.

#### Sign-in / Sign In — User-Facing Perspective

Use sign-in for end-user UI text, user actions, and user-facing flows or journeys.

Examples:
- Sign in with Google
- Sign in to the Console
- when the user signs in
- sign-in flow
- sign-in journey
