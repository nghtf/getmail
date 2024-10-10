# getmail

Expects getmail.rc like this

```
[retriever]
type = SimplePOP3SSLRetriever
server = pop.gmail.com
username = user@gmail.com
password = P@$$W0rd

[destination]
type = Maildir
path = ./maildir/

[options]
delete = false
```