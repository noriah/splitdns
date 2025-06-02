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

## configuration

example config.json

```json
{
  // listen on localhost at UDP port 8053
  "listen": "127.0.0.1:8053",
  "zones": [
    // root zone definition
    {
      "name": ".",
      // some sane default servers
      "servers": [ "8.8.8.8:53", "8.8.4.4:53" ]
    },
    {
      "name": "my.awesome.tld.",
      "servers": [
        "127.0.0.1:5353", // another dns server on localhost
        "10.12.34.56:53", // a dns server at 10.12.34.56
        "1.1.1.1:53" // cloudflare 1.1.1.1 dns
      ]
    },
    // more if desired
  ]
}
```

*remove comments. [golang json][go-json] follows [strict json spec][rfc7159]*

splitdns treats a `.` name as the most wild of wildcards. anything unmatched by
other zone definitions gets sent to the servers for this zone.
exclusion of this zone results in forwarding requests only for defined zones.

servers can be re-used in multiple zones.

zone/name order does not matter.

[go-json]: https://pkg.go.dev/encoding/json
[rfc7159]: https://www.ietf.org/rfc/rfc7159.txt
