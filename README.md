# download4

A downloader for a website which has neglected to upgrade its PHP since 2016.

> [!WARNING]
> Besides the small update in 2024, this software has not been maintained since 2023.
>
> Due to external circumstances, this software will not be maintained at all anymore.

### Usage:

```console
$ four-download -u <url>
```

url: URL of the thread to download from.

*To download slower:*

```console
$ four-download -u <url> -t <wait_time>
```

wait_time: Number of seconds to wait in-between requests. Default is 1.

### Sample `config.json` file:

```json
{
    "log_path": "./path/to/log/dir/",
    "download_path": "./path/to/downloads/dir/"
}
```

*Having a config file is optional. In a situation where the config file is not being used or cannot be found, the default paths will be used.*

### Default paths:

*Linux / MacOS:*

- Logs: `$HOME/.cache/download4/logs/`
- Downloads: `$HOME/Documents/download4/downloads/`
- Config: `$HOME/.config/download4/config/`

*Windows:*

- Logs: `%userprofile%/AppData/Roaming/download4/logs/`
- Downloads: `%userprofile%/Documents/download4/downloads/`
- Config: `%userprofile%/AppData/Roaming/download4/config/`
