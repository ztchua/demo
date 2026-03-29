---
name: github-pr
description: This skill is used to concisely summarize changes from the existing git branch against another branch (typically `main`), and to create a PR on GitHub.
disable-model-invocation: true
agent: Explore
allowed-tools: Read, Grep, Glob, Bash(git *)
---

# GitHub - Pull Request (PR) workflow

This skill concisely summarize differences into a GitHub pull request. Follow the steps below to effectively review code changes in 4 steps: context gathering, code analysis, summary creation and PR creation via GitHub MCP.

## When to offer this workflow

This skill should be used when the user requests to write a summary for merge request, and the underlying repository resides on any variant of GitHub's server.  

## Workflow steps

1. Context Gathering: Ask the user for details about the merge request, such as the purpose of the changes, key features, and any specific areas they want feedback on if necessary.

2. Code Analysis: Review the code changes in the merge request, understanding the modifications made, the rationale behind them, and how they fit into the overall codebase.

3. Summary Creation: Summarize the key points of the merge request, including the main changes, benefits, and any remaining concerns or questions that need to be addressed before merging.

4. PR Creation: Using existing GitHub MCP, authenticate and create the pull request on the existing branch against the target (by default: `main`) branch with the details prepared in point 3.

## Tips for creating effective merge request summaries

- Be concise and clear, focusing on the most important aspects of the changes.
- Be direct and procedural, outlining the changes in a logical order.
- Use bullet points or numbered lists to organize information when appropriate.
- Only gather context that are committed, not uncommitted changes, to ensure the summary reflects the actual changes being proposed.
- Use backticks to format code snippets, resources or bespoke terms for clarity.
- Goal is to create a easy to comprehend summary that can be quickly understood by reviewers and stakeholders.
- Output must always be in markdown compliant format so that it is also easy to copy and paste directly into a merge request description if needed. Ensure the summary is well-structured and visually clear when rendered.
- Output into a merge-request.md file in the root of the repository, which can be referenced against when creating the merge request via MCP.

Always extend on the following template for merge request summaries:
# Context
<Briefly describe the context and purpose of the merge request.>

## Major changes
<List and describe any changes that may introduce a breaking change, or is a new feature.>

## Minor changes
<List and describe any changes that are already present, and the change does not affect major functionality>

## Fixes
<List and describe any fixes in this section>