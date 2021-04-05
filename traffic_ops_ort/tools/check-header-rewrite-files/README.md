# check-header-rewrite-files

check-header-rewrite-files is a utility to verify all header rewrite directives in remap.config point to files that exist.

The config generation around creating and injecting header rewrites is complex, and it's difficult to verify that all scenarios which result in a header rewrite directive on a remap line also result in a header rewrite file being generated.

It's also possible there may be valid scenarios where it's possible to create data resulting in a header rewrite remap line without a corresponding file (such as with Profiles and Parameters), but which isn't a code bug. So it may not be possible to detect in the generation itself. Moreover, any detection would have to be after creating remap.config (or else we would have prevented it there), and thus it isn't obvious what ORT should do if it did detect it.

It's also possible a user may have a header rewrite file not managed by Traffic Control, and thus ORT can't safely simply create an empty file if one doesn't exist.

Therefore, this tool gives users and developers a way to detect code and data bugs resulting in malformed config, so the bugs can then be fixed.

# Usage

check-header-rewrite-files takes the output of atstccfg as input. If no missing files were found, it outputs nothing and returns exit code 0. If any files are missing, their names are printed and a nonzero exit code is returned.

# Examples

```sh
./atstccfg --traffic-ops-url="trafficops.example.net" --cache-host-name="my-cache" > allfiles.txt
go run check-header-rewrite-files < allfiles.txt
```

```sh
atstccfg -u "trafficops.example.net" -n "my-cache" > check-header-rewrite-files
```

```sh
go run check-header-rewrite-files < allfiles.txt
```

```sh
go build;
./check-header-rewrite-files allfiles.txt
```
