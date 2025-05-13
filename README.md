# CineChecker

**CineChecker** is a simple domain checker for content providers. It reads a list of provider URLs from a `providers.json` file, checks for redirects or updated URLs, and updates the JSON file if changes are detected.

## Features

- Reads provider URLs from `providers.json`.
- Checks each URL for valid response and redirections.
- Updates the provider URLs if redirected.
- Logs errors and results.
- Designed to run every 12 hours as a scheduled task.

## Usage

### Setup `providers.json`

The `providers.json` file should contain a JSON object where:

- Keys are **lowercase names** of the providers.
- URLs should **not have a trailing slash**.

Example:

```json
{
  "animepahe": { "url": "https://animepahe.ru" },
  "hdmovie2": { "url": "https://hdmovie2.mn" },
  "hianime": { "url": "https://hianime.sx" }
}
```

## How it Works

1. Opens and decodes `providers.json`.
2. For each provider:
   - Sends an HTTP GET request.
   - Follows redirects manually to detect changes.
   - Updates the URL if redirected.
3. Writes back the updated `providers.json` if changes are found.
4. Logs errors and final status.

## Notes

- Provider names **must be in lowercase**.
- Provider URLs **must not have a trailing slash**.
- Uses a 5-second timeout for each URL check.
- Only follows the first redirect, if any.

## License

MIT License
"""
