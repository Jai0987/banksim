## CNSE-IDP

A simple oAuth 2 server that we can use for various testing and research purposes to change user roles/permissions/etc as we test misconfiguration scenarios.  This package is written in typescript becuase it relies on and extend the `oauth2-mock-server` project which is here: https://github.com/axa-group/oauth2-mock-server

#### Requirements

1. You must have a newer version of nodejs installed.  I am using Node 20
2. You must have yarn installed as a package manager.  You can get it here: https://yarnpkg.com/
3. The GitHub repo will not have the dependencies included.  You run `yarn` the first time to down load them.
4. You can then run `yarn run exec` to start the server