# ip2region
go版本的ip2region https://github.com/lionsoul2014/ip2region

## 库特点
1. 舍弃原始库的文件搜索模式，使用golang的embed特性将数据文件直接嵌入到go程序中，只提供内存搜索
2. 支持多线程同时查询
3. 使用string形式ip查询地区信息，查询速度约为 400纳秒一次(0.0000004秒)，单线程每秒支持百万qps级别，如果使用uint32查地区信息，性能还可提升一倍
