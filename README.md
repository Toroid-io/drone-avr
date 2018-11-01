[![Build Status](https://bianca.toroid.io/api/badges/Toroid-io/drone-avr/status.svg?branch=master)](https://bianca.toroid.io/Toroid-io/drone-avr)

## drone-avr

`drone-avr` is a [drone](https://github.com/drone/drone) plugin for
building AVR C code. By default, it runs `make` in the project
directory. This may be overridden passing a `command` option.

## Example configuration

```yml
pipeline:
  kicad:
    image: toroid/drone-avr
    dependencies:
      - https://github.com/toroid-io/toroid-c-library   # Clone in current directory
    projects:
      - source: Project1                                # Makefile folder
        dependencies:
          - "https://github.com/myuser/awesome-lib"     # Clone this repo in the source directory
        arguments: "arg1 arg2"                          # These are passed to make as commmand line arguments
      - source: Project1                                # Makefile folder
        command: "make my_variant"                      # Custom command
        store: "cp build/output FW/"                    # Do something with your build artifacts
```

## Contributing

Don't hesitate to submit issues or pull requests.

## License

This project is made available under the GNU General Public License(GPL) version 3 or grater.
