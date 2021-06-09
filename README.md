# 💓 mebeats ![Go](https://github.com/wuhan005/mebeats/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/mebeats)](https://goreportcard.com/report/github.com/wuhan005/mebeats)

小米手环实时心率数据采集 - Your Soul, Your Beats!

* `cmd/mebeats-client`: the mebeats client. It collects the heart rate data from Mi Band and reports to server.
* `cmd/mebeats-server`: the mebeats sever. It receives the heart rate data and generate the badge.

## Requirement

* MiBand (2, 3, 4, 5, 6)
* macOS 11.3.1 or higher

## Run server

```bash
git clone git@github.com:wuhan005/mebeats.git

cd mebeats/cmd/mebeats-server

go build . && ./mebeats-server --key=<your_secret_key>
```

Server runs on `0.0.0.0:2830`.

### Run client

```bash
git clone git@github.com:wuhan005/mebeats.git

cd mebeats/cmd/mebeats-client

go build . && ./mebeats-client --addr=<mi_band_addr> --auth-key=<mi_band_auth_key> --server-addr=<mebeats_server_addr> --server-key=<your_secret_key>
```

Server runs on `0.0.0.0:2830`.

## License

MIT