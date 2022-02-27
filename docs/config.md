# 配置方式

- 基础配置说明

```yaml
service:
  name: kratos-layout
  version: 1.0.0
  grpc:
    addr: 0.0.0.0:9000
    timeout: 3s
  http:
    addr: 0.0.0.0:8000
    timeout: 3s
data:
  casdoor:
    endpoint: https://eof.com:8080
    client_id: test
    client_secret: test
    jwt_secret: test
    organization_name: test
    application_name: test
  mysql:
    source: root:123456@tcp(eof.com)/kratos-layout?charset=utf8mb4&parseTime=True&loc=Local
    max_idl: 10
    max_open: 100
    conn_max_lift: 3600s
  redis:
    addr: eof.com:6379
    read_timeout: 3s
    write_timeout: 3s
    password: 123456
    pool_size: 10
    min_idle_conn: 3
    max_conn_lifetime: 1h
    db: 1
  grpc:
    timeout: 5s
  nacos:
    addr: eof.com
    port: 8848
    namespace_id: namepace_id
    log_rotate_time: 1h
    log_max_age: 7
    log_level: debug
    cluster_name: DEFAULT
    group_name: DEFAULT_GROUP
    weight: 100
  sentinel:
    enabled: false
    group_name: sentinel-go
    data_id_flow: example-flow.yaml
    data_id_cb: example-cb.yaml
  skyWalking:
    addr: eof.com:11800
    enabled: false

```