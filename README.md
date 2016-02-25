# gogo-garage-opener
Go implementation of a Raspberry Pi garage door opener

Also see [gogo-garage-opener-ui](https://github.com/benjefferies/gogo-garage-opener-ui) implemented using ionics framework

Also see [Compile for arm] (https://gist.github.com/steeve/6905542)

#### TODO
##### Testing
##### Reconsider database option
Currently using github.com/mattn/go-sqlite3 to store data. go-sqlite3
uses a c binding file which means you need to install gcc-arm when compiling making
the build more awkward. As what is being stored isn't much maybe it would be easier
to just manage a json file on disk or some simple alternative.