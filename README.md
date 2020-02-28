# helm3-upload

用法：
```shell script

helm3-upload  -l https://172.27.139.109:30003//chartrepo/kubernetes-sigs -u "" -p "" -c ~/Downloads/charts/grafana-0.0.1.tgz

```

参数

|参数 | 简写 | 默认值 | 含义  |
|:-: | :-: | :-: | :-:  |
|user | u | admin | harbor账号 |
|password | p| Harbor12345 | harbor密码 |
|url | l|  | harbor的chart的地址 |
|chart | c|  | chart的本地路径 |
