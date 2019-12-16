# simple-blacklist
A simple tool to fetch and filter domain blacklists for use with tools like https://github.com/DNSCrypt/dnscrypt-proxy

## Usage
```shell script
➜ simple-blacklist --help
A simple tool to fetch, parse and merge domain blacklists

Usage:
  simple-blacklist [flags]

Flags:
  -h, --help            help for simple-blacklist
  -o, --output string   path to write final blacklist to
  -u, --url strings     comma-separated list of urls
```

## Examples
```shell script
# single list
➜ simple-blacklist -o /tmp/blacklist.txt -u 'https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts'
+------------------------------------------------------------------+---------+--------+
|                              SOURCE                              | DOMAINS | UNIQUE |
+------------------------------------------------------------------+---------+--------+
| https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts |   39711 |  39710 |
|                                                                  |
| /tmp/blacklist.txt                                               |   39711 |  39710 |
+------------------------------------------------------------------+---------+--------+

# multiple lists in single flags
➜ simple-blacklist -o /tmp/blacklist.txt -u 'https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts,https://hosts-file.net/ad_servers.txt'

+------------------------------------------------------------------+---------+--------+
|                              SOURCE                              | DOMAINS | UNIQUE |
+------------------------------------------------------------------+---------+--------+
| https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts |   39711 |  39710 |
| https://hosts-file.net/ad_servers.txt                            |   45736 |  45081 |
|                                                                  |
| /tmp/blacklist.txt                                               |   85447 |  78924 |
+------------------------------------------------------------------+---------+--------+

# multi list in single and repeated flag
➜ simple-blacklist -o /tmp/blacklist.txt \
  -u 'https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts' \
  -u 'https://mirror1.malwaredomains.com/files/justdomains' \
  -u 'http://sysctl.org/cameleon/hosts' \
  -u 'https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt,https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt' \
  -u 'https://hosts-file.net/ad_servers.txt'
+------------------------------------------------------------------+---------+--------+
|                              SOURCE                              | DOMAINS | UNIQUE |
+------------------------------------------------------------------+---------+--------+
| https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt |      34 |     34 |
| http://sysctl.org/cameleon/hosts                                 |   20567 |  20138 |
| https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts |   39711 |  39710 |
| https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt       |    2703 |   2703 |
| https://mirror1.malwaredomains.com/files/justdomains             |   26863 |  26863 |
| https://hosts-file.net/ad_servers.txt                            |   45736 |  45081 |
|                                                                  |
| /tmp/blacklist.txt                                               |  135614 | 113014 |
+------------------------------------------------------------------+---------+--------+

```