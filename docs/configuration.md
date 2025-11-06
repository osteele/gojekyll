---
layout: default
title: Configuration
---

# Configuration

This page details the configuration options available in Gojekyll's `_config.yml` file.

## Permalink Timezone

Gojekyll allows you to specify a timezone for generating dates in permalinks. This is useful for ensuring that your post URLs are consistent regardless of the server's local timezone.

**Option:** `permalink_timezone`

**Values:**
- A valid IANA Time Zone Database name (e.g., "UTC", "America/New_York", "Europe/Berlin").
- If left unset or empty, Gojekyll will use the server's local timezone. This matches the default behavior of Jekyll.

**Example `_config.yml`:**

```yaml
# _config.yml
title: My Awesome Blog
permalink_timezone: "UTC" # Generates permalink dates in UTC
```

**Behavior:**

- If `permalink_timezone` is set to a valid timezone, all date-based permalink variables (like `:year`, `:month`, `:day`) will be calculated based on that timezone.
- If `permalink_timezone` is invalid (e.g., "Invalid/Timezone"), Gojekyll will log a warning and fall back to using the server's local timezone.
- This setting affects how dates from your posts' front matter are interpreted for URL generation. For instance, a post dated `2023-01-01 02:00:00 +0500` (5 AM in a +0500 timezone) would use January 1st for permalink generation if `permalink_timezone` is "Asia/Kolkata" (UTC+5:30) or similar, but might use December 31st if `permalink_timezone` is "America/Los_Angeles" (UTC-8), depending on the exact date and time.

By default, if you have posts with explicit timezone offsets in their `date` front matter, Gojekyll (like Jekyll) first converts this date to the site's effective timezone (either local or the one specified by `permalink_timezone`) before extracting year, month, and day for the permalink.
