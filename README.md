# download4

### Usage:

    $ ./four-download -u <url>

url: URL of the thread to download from.

*To download quicker:*

    $ ./four-download -u <url> -t <worker_count>

worker_count: Number of threads to use when concurrently downloading. Default is 1.

### Sample `config.json` file:

```
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