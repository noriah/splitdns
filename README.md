# splitdns

forwards queries to upstream dns servers based on domain

i need a way to query a specific dns server for one of my projects.
making the records public is unfavorable. i could use a local `bind` instance to
set authoritative NS for the domains i care about, but i rather explore this
method.

## usage

```shell
splitdns /path/to/config.json
```

see [config.example.json](/config.example.json)
