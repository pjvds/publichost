# publichost

Publichost is an easy way to make local server part of the world wide web.

    $ publichost -http localhost:3000
    localhost:3000 > afr.publichost.me

## Connection protocol

PROXY HTTP 127.0.0.1
ACK publicaddress || NACK reason

OPEN TCP 127.0.0.1
ACK streamid || NACK reason

DATA streamid 0x00 0x00 0x00 0x00
ACK || NACK reason

EOF streamid R || W
ACK || NACK reason

## Client

A client receives incomming requests to:

* Open new streams
* Proxy data to streams
* EOF streams
* Close streams

It actively monitors the open streams for incomming data and sends
this to the publichost front end server. Which will proxy it to
the origin connection.