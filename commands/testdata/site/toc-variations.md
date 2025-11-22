---
title: TOC Syntax Variations Test
layout: default
---

# TOC Syntax Variations

This page tests different ways to specify the TOC marker.

## Test 1: Unordered list with {:toc}

This should work (Jekyll supports this):

* This will be replaced by TOC
{:toc}

## Test 2: Ordered list with {:toc}

This should NOT work (Jekyll doesn't support this):

1. This should remain
{:toc}

## Test 3: Block syntax {::toc}

This should NOT work ({::toc} is invalid kramdown syntax):

{::toc}

## Test 4: Block syntax in a list

This should NOT replace the list item (Jekyll only supports {:toc} for list replacement):

* This should remain
{::toc}

## Test 5: Whitespace variations

Testing with extra whitespace:

* Whitespace test
{: toc }

## Test 6: Heading with .no_toc

### This heading should be excluded
{:.no_toc}

The above heading has `{:.no_toc}` and should not appear in any TOC on this page.

### This heading should appear

This is a normal heading.

## Summary of Jekyll TOC Behavior

Based on Jekyll/Kramdown documentation and testing with Jekyll 4.4.1:

1. **{:toc} in unordered list**: Replaces the `<ul>` containing the marker with the TOC ✓
2. **{::toc} anywhere**: NOT valid kramdown syntax, remains as literal text ✗
3. **{:toc} in ordered list**: NOT supported, marker remains as literal text ✗
4. **{:.no_toc} after heading**: Excludes that heading from all TOCs ✓

Note: `{::toc}` is NOT part of the kramdown specification and is not processed by Jekyll.
