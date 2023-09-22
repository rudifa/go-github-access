# Access to github rest api

## Notes

Unauthenticated requests are limited to 60 per hour.
To get more requests (5000 per hour), you need to authenticate with a github token (classic).

This version requires a .env file with

```
GITHUB_TOKEN=ghp_...
```
