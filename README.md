# tun-container

## TUN

基于 wireguard 的网络库实现。

TODO:

- [] 转发 DNS 请求
- [] 解析 DNS UDP 包
- [] 内网流量转发
- [] 解析来源 IP

## Route

基于源策略的路由。

ref: https://superuser.com/questions/376667/how-to-route-only-specific-subnet-source-ip-to-a-particular-interface

```bash
ip rule add from <source>/<mask> table <name>
ip route add 1.2.3.4/24 via <router> dev eth4 table <name>
```

## Docker

创建 Docker 网络：

- **MacVlan**:

    创建 MacVlan 网络，二层设备流量转发，相比于 Bridge 模式性能更好。

    ```bash
    sudo docker network create -d network --gateway 192.168.10.1 --subnet 192.168.10.0/24 tunnet
    ```

- **Bridge**:

    Bridge 网络，Docker 的默认网络模式。

    ```bash
    sudo docker network create --gateway 192.168.20.1 --subnet 192.168.20.0/24 tunnet
    ```

删除 Docker network：

```bash
sudo docker network rm tunnet
```

## Debug

```bash
# 指定 Docker 网络
sudo docker run --rm -it --network tunnet golang:1.16
# 指定 IP 地址
sudo docker run --rm -it --network tunnet --ip 192.168.20.3 golang:1.16
```
