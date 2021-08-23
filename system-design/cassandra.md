# Cassandra



## 单机

Storage Engine与LevelDB相同，可以非常快速的进行写操作，因为不用触发磁盘I/O. 

参考这篇：[SQL and Comparison](storage-enigine.md#LSMTreeLogStructuredMergeTree))



## 多机

- Sharding: Consistent hashing来确定节点所在位置。并通过virtual node来解决hot spot问题。
- High availability -> replica
- Decentralized, 所有的节点都可以serve query
- Configurable Consistency-> Quorum
  - W+R>Number of replica

## Reference

- [读过本文才算真正了解 Cassandra 数据库](https://www.infoq.cn/news/underlying-storage-of-uber-change-from-mysql-to-postgres?utm_source=related_read_bottom&utm_medium=article)
- 大规模分布式存储系统：原理解析与架构实战, Chapter5.1: Amazon Dynamo