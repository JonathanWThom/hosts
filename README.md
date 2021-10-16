# Hosts

Simple API that gives a 👍 or 👎 to hostnames based on a stored obfuscated list.

### Usage

Hosts are looked up from a SQLite database at `hosts.db`. This database must be
populate and shipped with the deployment of this application. To populate it,
you need two things:

1) A URL contain a list of hostnames you want to block (such as [this](https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews/hosts)).

2) (Optionally) A key that will be used to hash the values as they are written
to the database. This key will need to be shared with any clients of the
service. `export HASH_KEY="your-key-value"` prior to running the populate task.

[Install Go](https://golang.org/doc/install) and then run `make popc
h=<your-hosts-url>`

On `make start`, a server will be available that responds to one API request:
`curl 'localhost:8080/allow?url=<your-hashed-and-encoded-hostname>\n' -i`

### Clients

Clients need to do the following:

1) Strip urls down to just their hosts, e.g. google.com, not
https://google.com?s=foo.

2) Hash the hostname using HMAC SHA256, passing the same key as used above.

3) Base64 encode that value, and include it as the `url` parameter in a `GET`
request to the `/allow` endpoint.

The response will be JSON of either `{"allow":true}` or `{"allow":false}.

### Hosting

`make build` and `make push` will build and push a container (including the
hosts.db file) to a repository of your choice.

### License

MIT
