# gomd
Cross-platform Markdown editor written in Go.

- Edit files in your browser: where Markdown usually ends up anyway.
- No internet connection needed though. It stays all on your computer.

## Installation
    $ go get -u github.com/nochso/gomd

## Usage
Open an existing file and edit it:

    $ gomd todo.md

See the command line help for more:
```
usage: gomd [<flags>] <file>

Flags:
      --help            Show context-sensitive help (also try --help-long and
                        --help-man).
  -a, --protocol=https  Application protocol (http or https) used by webserver
  -p, --port=10101      Listening port used by webserver
  -d, --daemon          Run in daemon mode (don't open browser)
      --version         Show application version.

Args:
  <file>  Markdown file
```

## Changes
See the [CHANGELOG](CHANGELOG.md) for the full history of changes between releases.

## License
This project is licensed under the MIT license. See the [LICENSE](LICENSE.md) file for the full license text.

## Credits
* [SimpleMDE](https://github.com/sparksuite/simplemde-markdown-editor) - WYSIWYG*ish* MD editor written in JS
* [gomd](https:github.com/nochso/gomd) - This is a fork of the original
* [minica](https://github.com/jsha/minica) - A simple Certificate Authority (CA) used to generate the certificates used for TLS in https protocol.