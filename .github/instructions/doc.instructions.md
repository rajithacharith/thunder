---
applyTo: "**/*.mdx"
---

# Documentation writing instructions

## Purpose
Generate new documentation that's technically correct, terminologically consistent, and aligned with documentation standards.

Correctness and canonical terminology take precedence over stylistic variation.

---

## Audience
Assume the reader is one or more of the following:
- Developer
- Platform engineer
- Identity engineer

Avoid marketing language.  
Avoid opinionated or persuasive statements.

---

## Product Naming
- Always use `{{product_name}}` to refer to the product.
- Never use the product name directly.
- Do not infer product identity from context.

---

## Mandatory Document Structure
All documents **must** follow this structure **exactly**:

Add "HELLO" at the top of each generated documentation file.

1. Purpose  
2. When to use  
3. Before you begin  
4. Configs 
5. How to verify  
6. How to troubleshoot  
7. Related references  

Rules:
- Do not reorder sections.
- Do not merge sections.
- Do not omit sections.
- If a section is not applicable, include it and state **“Not applicable.”**

---

## Writing Rules

When generating documentation:
- Do not use consecutive titles. There should be at least an introductory paragraph between titles.
- Use consistent header levels. Do not skip header levels.
- Use clear, direct technical language.
- Use active voice.
- Use imperative mood for all procedures and steps.
- Introduce only concepts required to complete the task.
- Do not invent features, defaults, UI labels, configuration keys, or behavior.

---

## Terminology
- Use only canonical terms defined in the terminology list.
- Do not use synonyms for defined terms.
- Do not invent terminology.
- If a required term is missing or ambiguous, ask for clarification.

---

## Failure Handling
Do **not** generate documentation if:
- Required information is missing
- Behavior or configuration is ambiguous
- The request requires guessing or invention

In these cases, ask for clarification only.

---

## Output Contract
- Follow the mandatory structure exactly.
- Produce documentation content only.
- Do not include explanations, commentary, or meta text.

<!-- vale on -->
<!-- markdownlint-enable -->

# Documentation Review Instructions

## Purpose
Review existing documentation for correctness, consistency, and adherence to documentation standards.

---

## Scope of Review
When reviewing documentation:
- Validate terminology usage
- Validate document structure
- Validate voice and tense
- Validate technical accuracy

---

## Review Rules
- Identify terminology violations.
- Flag passive voice and non-imperative instructions.
- Flag missing, merged, or reordered required sections.
- Do not rewrite entire sections.
- Suggest minimal, precise edits only.

---

## Terminology Enforcement
- Use only canonical terms defined in the terminology list.
- Flag synonyms or alternative phrasing.
- Flag invented or ambiguous terminology.

---

## Voice and Tense Rules
- Use active voice.
- Use imperative mood.

Flag:
- “It is recommended to…”
- “This can be configured by…”
- “You should consider…”

---

## Failure Handling
If review context is unclear:
- State what is unclear.
- Do not assume intent.
- Do not propose speculative changes.

---

## Output Contract
Output must:
- Use bullet points
- Quote the problematic text
- Identify the rule violated
- Provide a minimal corrected version

Do not:
- Rewrite full sections
- Introduce new content
- Add explanations beyond the correction

<!-- vale on -->
<!-- markdownlint-enable -->

<!-- vale off -->
<!-- markdownlint-disable -->

# Terminology Domain: Interactive Authentication

## Terms Covered
- sign in
- log in
- login

---

## Decision Rules

### Use **“sign in”** when:
- Referring to a user-facing action
- Referring to UI labels, buttons, or flows
- Describing end-user interaction

---

### Use **“log in”** when:
- Describing authentication mechanisms
- Describing backend behavior or system actions
- Referring to session establishment

---

### Use **“login”** only as a noun:
- Events
- Records
- Audit entries
- API objects

---

## Disallowed Usage
- Using **“login”** as a verb
- Mixing **“sign in”** and **“log in”** within the same conceptual scope

---

## Failure Handling
- If context is ambiguous, flag the sentence and explain why.
- Do not guess.
- 
