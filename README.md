dmv
===

[![](https://godoc.org/github.com/tomsteele/dmv?status.svg)](http://godoc.org/github.com/tomsteele/dmv)

Simple authentication for Martini. Does not handle state or make use of the sessions middleware. It only provides a means of initial authentication. Beacuse of this, it is up to the application to implement its own authorization. External authentication mediums will provide profile information. For example, the OAuth 2.0 Facebook function provides information about the user including their name and email address.

Authentication is handled on a per route basis, allowing applications to easily use multiple authentication mediums.

## Supported Mediums
- Local (Form)
- Local (Basic)
- Github OAuth 2.0
- Facebook OAuth 2.0
- Google OAuth 2.0

## Usage
There is sample usage for each Auth* function in the docs. Also see [examples](https://github.com/tomsteele/dmv/tree/master/examples).
