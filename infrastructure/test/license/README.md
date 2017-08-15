Verifying Licenses
==================

**tl;dr** Run `check_license`. Output indicates problems.

This is an automatic license checker aimed at ensuring compliance with
the Apache Software Foundation processes. It does not ensure compliance,
but it catches many common errors automatically.

`check_license`
---------------

This script will automatically build and run the license check on the
current directory. The exit code will be non-zero on error. There will
be no output on success.

`license`
---------

This is the binary with all the logic in it. `license` will list every
file in every subdirectory, starting from the current directory, and
determine the most likely license for that file. It errs on the side of
false positives, since the consequences of a false negative are
considerably more serious.

Run `license -q` to suppress the printing of non-problematic files.

`LICENSE`
---------

The `LICENSE` file is at the root of the project (from whence you ought
to run `license`). This must comply with the requirements of the Apache
Software Foundation and is intended for human consumption.

Nevertheless, with a bit of careful writing, it's possible to have
`license` help verify that everything gets covered.

Lines that begin with an `@`-symbol are interpreted as a path
specification that describes a set of files covered by the license.
`license` does not validate that `@` files are actually licensed
correctly, merely that they are mentioned. This covers the most common
case, which is adding a (even potentially correctly licensed!) file and
forgetting to mention it in the `LICENSE` file.

Likewise, it's impermissible to use an `@`-line that describes no files.
This usually happens when a dependency is removed and the `LICENSE` file
does not get updated properly.

`@`-lines are interpreted by
[path.Match](https://golang.org/pkg/path/#Match), the syntax for which
is:

    pattern:
        { term }
    term:
        '*'         matches any sequence of non-/ characters
        '?'         matches any single non-/ character
        '[' [ '^' ] { character-range } ']'
                    character class (must be non-empty)
        c           matches character c (c != '*', '?', '\\', '[')
        '\\' c      matches character c

    character-range:
        c           matches character c (c != '\\', '-', ']')
        '\\' c      matches character c
        lo '-' hi   matches character c for lo <= c <= hi

`.dependency_license`
---------------------

Sometimes, there's no reasonable way to automatically detect the
appropriate license for a file, especially files that don't support
comments or are binary. `.dependency_license` allows you to document the
exceptions so that new files show up clearly.

The `.dependency_license` must appear in the root of the project.

Each line should either be empty, a comment (prepended by an octothorp),
or a license exception line. A license exception line is a regular
expression, a comma, then the name of a license, then optionally an
octothorp followed by a comment (which may not contain a comma!).

    license-exception:
        regex ',' license-name [ '#' { commentable-char } ]       Associates the license with the file.        
        regex ',' '!' license-name [ '#' { commentable-char } ]   Disassociates the license from the file.

    regex: A regular expression accepted by golang regexps, described here: https://golang.org/s/re2syntax

    license-name:
        'Apache'    Apache License
        'BSD'       Berkeley Software Distribution License
        'MIT'       Massachusetts Institute of Technology License
        'GoBSD'     BSD-style license used by the GoLang team
        'ISC'       Internet Systems Consortium
        'X11'       MIT License, by an older name.
        'WTFPL'     Do What the Fuck You Want to Public License
        'GPL/LGPL'  Either the GNU General Public License or the GNU Lesser General Public License
        'Docs'      A documentation file
        'Empty'     An empty file
        'Ignored'   A file that ought not be analyzed for compliance

    commentable-char: Any character other than a ','

Best Practices
--------------

License management can be tricky at the best of times. The goal of this
tool is to automate as much of that as possible. Here are some best
practices:

-   **`check_license` before you commit.** If it prints anything, you
    probably need to add license headers.
-   **Do not `Ignore` files.** If it's reasonable to quiet `license`
    about a false positive or negative in another way, do that instead.
-   **If an unrecognized file has a header, update `license`, not
    `.dependency_license`.** It's relatively straightforward to add
    license recognition to `licenseList.go`. Doing it that way benefits
    future files as well.
-   **Run `check_license` as part of Continuous Integration.** Issues
    are not usually difficult to fix, but automatic running allows them
    to be fixed promptly.

