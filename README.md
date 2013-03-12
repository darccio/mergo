# Mergo

A helper to merge structs in Golang. Useful for configuration default values, avoiding messy if-statements.

Also a lovely [comune](http://en.wikipedia.org/wiki/Mergo) (municipality) in the Province of Ancona in the Italian region Marche.

![Mergo dall'alto](http://www.comune.mergo.an.it/Siti/Mergo/Immagini/Foto/mergo_dall_alto.jpg)

## Status

It is quick hack to scratch my own itch around how to handle configuration default values. It works fine but it needs a lot more of testing and real world usage.

[![Build Status][1]][2]

[1]: https://travis-ci.org/imdario/mergo.png
[2]: https://travis-ci.org/imdario/mergo

## Installation

    go get github.com/imdario/mergo

    // use in your .go code
    import (
        "github.com/imdario/mergo"
    )

## Usage

You can only merge same-type structs with exported fields initialized as zero value of their type. Mergo won't merge unexported (private) fields but will do recursively any exported one.

    if err := mergo.Merge(&dst, src); err != nil {
        // ...
    }

More information and examples in [godoc documentation](http://godoc.org/github.com/kless/modsql).

## Contact me

If I can help you, you have an idea or you are using Mergo in your projects, don't hesitate to drop me a line (or a pull request): [@im_dario](https://twitter.com/im_dario)

## About

Written by [Dario Castañé](http://dario.im).

## License

[BSD 3-Clause](http://opensource.org/licenses/BSD-3-Clause) license, as [Go language](http://golang.org/LICENSE).
