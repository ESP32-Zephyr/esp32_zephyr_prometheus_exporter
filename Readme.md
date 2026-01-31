
## Build & Install

```console
mkdir build && cd build && cmake ..
cmake --build . --target build-go-binary
sudo make install
```

## Start service

```console
sudo systemctl daemon-reload
sudo systemctl start esp32_prometheus_exporter.service
```
