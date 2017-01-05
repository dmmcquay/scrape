# Contributing

We love pull requests from everyone. By participating in this project, you
agree to abide by the golang [code of conduct].

[code of conduct]: https://golang.org/conduct

The basic steps begin with a clone this repo:

    $ git clone https://github.com/dmmcquay/scrape $GOPATH/src/github.com/dmmcquay

add a feature and some tests then run the tests, check your formatting:

    $ go test github.com/dmmcquay/scrape
    $ go vet github.com/dmmcquay/scrape
    $ golint github.com/dmmcquay/scrape

If things look good and tests pass commit and push to your remote:

    $ git add (files you changed)
    $ git commit -m "Job's done"
    $ git push mine feature

Push to your fork and [submit a pull request][pr].

[pr]: https://github.com/thoughtbot/factory_girl_rails/compare/

At this point you're waiting on us. We will comment on the pull request request
within three business days (and, typically, one business day). We may suggest
some changes or improvements or alternatives.

Some things that will increase the chance that your pull request is accepted:

* Write tests.
* follow good Go style, including [effective go], running [go vet] and [golint].
* Write a [good commit message][commit].

[effective go]: https://golang.org/doc/effective_go.html
[go vet]: https://golang.org/cmd/vet/
[golint]: https://github.com/golang/lint
[commit]: http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html

contribution guidelines borrowed from [factory girl rails].

[factory girl rails]: https://github.com/thoughtbot/factory_girl_rails
