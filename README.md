dmv
===

Simple authentication for Martini. Does not handle state or make use of the sessions middleware. It only provides a means of initial authentication. External authentication mediums will provide profile information. For example, the OAuth 2.0 Facebook function provides information about the user including their name and email address.

[API Reference](http://godoc.org/github.com/tomsteele/dmv)


## Supported Mediums
- Local (Form)
- Github OAuth 2.0
- Facebook OAuth 2.0

## Usage
There is sample usage for each Auth* function in the docs. Also see [examples](https://github.com/tomsteele/dmv/tree/master/examples).
