# ibkr

Another unofficial Golang Interactive Brokers TWS/IB api client.

The official and unofficial SDKs uses a pattern of sending requests and raise events when a response is received.

The library acts differently. When you call a method to send a request, the function waits until the corresponding
response reaches.

For streaming methods, the response usually contains a channel to listen for incoming events.

## Acknowledgements

This library is strongly based on https://github.com/scmhub/ibapi and the official SDKs.

Thanks to the authors and collaborators for the excellent job.

## LICENSE

Copyright © 2025 SCM
Copyright © 2025 Mauro H. Leggieri

[MIT](/LICENSE)
