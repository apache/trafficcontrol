# cache-config

Traffic Control cache configuration is done via the `t3c` app and its ecosystem of sub-apps.

These are provided in the RPM `trafficcontrol-cache-config`.

To apply Traffic Control configuration and changes to caches, users will typically run `t3c` periodically via `cron` or some other system automation mechanism. See [t3c](./t3c/README.md).

The `t3c` app is an ecosystem of apps that work together, similar to `git` and other Linux tools. The `t3c` app itself has commands to proxy the other apps, as well as a mode to generate and apply the entire configuration.
