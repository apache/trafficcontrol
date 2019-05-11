# Building RPMs & Binaries
To build the RPMs and binaries you will need docker.
## Usage
```bash
docker-compose -f docker-compose.build_rpm.yml build --no-cache builder
docker-compose -f docker-compose.build_rpm.yml up --force-recreate --exit-code-from builder
```
This should dump all artifacts in the %repository root%/dist directory when it's done (like all other ATC build jobs).

## Version
If you need to manipulate the version information, you can export a few environment variables to supply your own overrides for the [VERSION](../version/VERSION) file.
* VER_MAJOR (integer)
* VER_MINOR (integer)
* VER_PATCH (integer)
* VER_DESC (short string)
* VER_COMMIT (git short hash string)
* BUILD_NUMBER (integer)
