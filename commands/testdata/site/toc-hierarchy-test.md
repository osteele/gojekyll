---
title: TOC Heading Hierarchy Test
layout: default
---

# Heading Hierarchy Tests

## Test 1: Starting with H3

{:toc}

### First H3 Heading

This document starts with H3, skipping H1 and H2.

#### H4 Under First H3

#### Another H4

### Second H3 Heading

## Test 2: Multiple H1 Headings

# First Top Level

{::toc}

## Under First

### Subsection

# Second Top Level

## Under Second

### Another Subsection

## Test 3: Gaps in Heading Levels

# Level 1

{:toc}

### Level 3 (skipping 2)

###### Level 6 (skipping 4 and 5)

## Level 2 (after 6)

## Test 4: Deeply Nested Structure

### Level 3 Start

#### Level 4

##### Level 5

###### Level 6

##### Back to Level 5

#### Back to Level 4

### Another Level 3

## Test 5: Flat Structure (All Same Level)

### Item 1

### Item 2

### Item 3

### Item 4

## Summary

This page tests TOC generation with:
- Missing heading levels
- Multiple H1 headings
- Deeply nested hierarchies
- Flat (same-level) structures
- Heading levels appearing out of order
