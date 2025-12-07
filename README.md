# FancySpaces

A platform to distribute artifacts like Minecraft plugins or other files.
This platform is only available to FancyInnovations products at the moment, but might be opened to everyone in the future.

### API usage

**Fetch latest version info:**

```http
GET https://fancyspaces.net/api/v1/spaces/{space_id}/versions/latest
```

Append the `?channel={channel}` query parameter to receive the latest version from a specific channel.
Append the `?platform={platform}` query parameter to receive the latest version for a specific platform.

**Download latest version:**

```http
GET https://fancyspaces.net/api/v1/spaces/{space_id}/versions/latest/files/{file_name}
```

Append the `?channel={channel}` query parameter to receive the latest version from a specific channel.
Append the `?platform={platform}` query parameter to receive the latest version for a specific platform.