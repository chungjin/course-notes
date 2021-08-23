# SQL and Comparison

[TOC]



## Main Concept



### Storage Engine

#### B+ tree

- Feature: Random I/O, range query.

![B+ tree](https://user-images.githubusercontent.com/11788053/126915521-9466e544-bdca-4131-984b-a0446a85442e.png)

以InnoDB为例，B+ Tree每个节点对应一个page.



Write:

1. append to WAL(disk)
2. update B+ tree
3. 如果mem中被修改的页面超过了一定比例，flush to disk



Cache:

1. LRU算法，第一级Cache.
2. Cache分层，因为有时候全表扫描会产生大量脏数据，这个时候只是放在第二级Cache中。除非在短时间内被visit多次，再放入第一级。



#### LSM Tree(Log Structured Merge Tree)

是Cassandra和levelDB的storage engine.

Featured：

- 可以有效避免对磁盘的Random Access.

- 适合write heavy system

- 但是在读取时可能要访问较多的磁盘文件。miss的时候需要访问磁盘。



![LSM Tree](https://user-images.githubusercontent.com/11788053/126915800-0bd79223-bda2-4e28-b5e9-bbc59d4659f7.png)

- Write/Read在MemTable中进行，满后再flush to disk.
- MemTable满后，转为immutable Memtable. 后台程序把immutable memtable转为SST数据。
- SSTable中的文件按照pk排序。每个文件由自己的pk范围，例如`a-c`
- Manifest file(清单文件): key所在的SST位置。



Write Access:

​	- Update Memtable

Read Access:

​	- check Memtable, check immutable Memtable, check SST.



### Transaction and Concurrency Control(ACID)

Deadlock:

	- 引入lock超时
	- Deadlock检测



Copy on write:

​	这是个为了提高读效率的方案，因为大多数数据库, read/write > 6. 

	- No lock on read(configurable)
	- Lock on write
	- 在对节点修改的时候，整条branch复制，修改完毕之前，read读取的是stale data.



MVCC(multi-version concurrency control):

- No lock on read, 用版本号来实现类lock的optimistic lock



### Failure Tolerance

- WAL + Checkpoint





## Data Proxy



What if the dataset is too large for one single machine to hold? For MySQL, the answer is to use a DB proxy to distribute data, [either by clustering or by sharding](http://dba.stackexchange.com/questions/8889/mysql-sharding-vs-mysql-cluster)



- Clustering is a decentralized solution. Everything is automatic. Data is distributed, moved, rebalanced automatically. Nodes gossip with each other, (though it may cause group isolation).
- Sharding is a centralized solution. If we get rid of properties of clustering that we don’t like, sharding is what we get. Data is distributed manually and does not move. Nodes are not aware of each other.

## Comparision

|                 | mysql                              | Cassandra<br />Dynamo                                        | Redis                                      | Kafka                   | ElasticSearch                   |
| --------------- | ---------------------------------- | ------------------------------------------------------------ | ------------------------------------------ | ----------------------- | ------------------------------- |
| Architecture    | M-S                                | serverless                                                   | M-S                                        | M-S(多个master-slave组) | M-S(多个master-slave组)         |
| Shard           |                                    | Consistent Hashing+virtual node                              | x                                          | Topic                   | Consistent Hashing+virtual node |
| Durability      | WAL                                | ...                                                          | ...                                        | ...                     | ...                             |
| Consistency     | lock on write or read(configrable) | W+R>N(保证读到)<br />多master的结构，引入时钟向量避免冲突(不完全准确) | Eventual                                   | Eventual                | Eventual                        |
| Xsaction        | Full and join support              | 支持batch, LWT                                               | 支持batch(script or watch)                 | x                       | x                               |
| Master failover | manual swich to Slave              | Not applicable<br />短时unreachable, 等待，恢复后用Markle树来更新stale的部分<br /> | Sentinel(类Raft, 监督master状态和彼此状态) | ZK                      | 类Raft                          |



