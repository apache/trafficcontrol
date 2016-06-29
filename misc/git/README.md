Related files for use with this project and git.

## Hooks

A collection of useful pre-commit hooks can be found in the `pre-commit-hooks` directory.

#### Installing pre-commit hooks

In the `github.com/Comcast/traffic_control/.git/hooks` directory, create a symbolic link to the `pre-commit` executable contained in this directory. 

```bash
cd github.com/Comcast/traffic_control/.git/hooks
ln -s ../../misc/git/pre-commit .
```

Now, all executables in the `pre-commit-hooks` directory will be run on commit.

#### Adding pre-commit check

Once the pre-commit file is in place, all executables in the `pre-commit-hooks` directory will be run. Simply add an executable there. Exiting with non-zero status from this script causes the git commit to abort (the commit contents will be unaffected).

#### Skipping

To commit without running the hooks, use the `no-verify` flag.

```bash
git commit --no-verify
```
