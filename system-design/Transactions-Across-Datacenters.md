# Transactions Across Datacenters

[TOC]

## Consistency?

Weak Consistency

- After a write, reads may or may not see it
- Best effort only
- Eg: memcache, video, realtime multiplayer games

Eventual Consistency

- After a write, reads *will eventually* see it
- Eg: Mail, Search Engine Indexing, DNS, Amazon S3

Strong Consistency

- After a write, read will see it
- Eg: Datastore, file systems, RDBMSes, Azure tables.



## 常见机制

![](https://lotabout.me/2019/Data-Synchronization-in-Distributed-System/sync-trade-off.svg)

1. Backup，即定期备份，对现有的系统的性能基本没有影响，但节点宕机时只能勉强恢复(如果数据写入master, 但还没来得及备份，会有损失数据的风险。)

2. Master-Slave，主从复制，异步复制每个指令，可以看作是粒度更细的定期备份。Write through Master, read through slave(Commonly used in DB, read heavy system).

3. Multi-Muster，多主，也称“主主”，MS 的加强版，可以在多个节点上写，事后再想办法同步。但是依赖全局clock来保证order, 但很难做到，一般不支持transaction的操作。

4. 2 Phase-Commit，二阶段提交，同步先确保**通知到所有节点**并得到他们同意后再写入，性能容易卡在“主”节点上。每个transaction都需要在所有server中选取一个作为Master来执行这个xsaction操作。类似于支持T ransaction的M-M,2PC协议保证多台服务器上的操作要么全部成功，要么全部失败。

5. Paxos，类似 2PC，同一时刻有多个节点可以写入，也只需要通知到**大多数**节点，有更高的吞吐。Paxos协议用于保证**同一个数据分片**的**多个副本**之间的数据一致性。通常这些副本存在不同的数据中心。

   2PC协议最大的缺陷在于无法处理协调者宕机问题。如果协调者宕机，那么2PC协议中的每个参与者可能都不知道事务应该提交还是回滚，整个协议被阻塞，执行过程中申请的资源都无法释放。 **因此常见做法是将2PC和Paxos协议结合起来** 通过2PC保证多个数据分片上操作的原子性，通过Paxos协议实现同一个数据分片的多个副本之间的一致性，另外通过Paxos协议解决2PC协议中协调者宕机问题。当2PC协议中的协调者出现故障时，通过Paxos协议选举出新的协调者继续提供服务。eg: zookeeper



## 常见系统的Trade off

### Redis

edis 3.0 开始引入 Redis Cluster 支持集群模式，个人认为它的设计很漂亮，大家可以看看[官方文档](https://redis.io/topics/cluster-spec)。

- 采用的是主从复制，**异步**同步消息，极端情况会丢数据
- 只能从主节点读写数据，从节点只会拒绝并让客户端重定向，不会转发请求
- 如果主节点宕机一段时间，从节点中会自动选主
- 如果期间有数据不一致，以最新选出的主节点的数据为准。

一些设计细节：

- Redis 的 Key 会被分配(分片/分桶？)到 16384 个 slot 中，每个节点提供部分 slot 的数据
- 分配的算法为 `HASH_SLOT = CRC16(Key) mod 16384`
- 集群的节有一个随机生成的唯一 ID，节点的 IP 可以变，但 ID 不会变
- 新节点加入时先执行 `MEET` 来认识集群中的某个节点，集群节点间相互“八卦(gossip)”，最终相互认识
- 主从的粒度是节点，不是 slot。
- 自动选主，使用类似 Raft 的选主机制。
- 也提供了 `WAIT` 指令来来保证写入时同步复制到从节点。

### Kafka

Kafka 的分片粒度是 Partition，每个 Partition 可以有多个副本。副本同步设计参考[ 官方文档](https://cwiki.apache.org/confluence/display/KAFKA/Kafka+Replication)

- 类似于 2PC[[1\]](https://lotabout.me/2019/Data-Synchronization-in-Distributed-System/#fn1)，节点分主从，**同步**更新消息，除非节点全挂，否则不会丢消息
- 消息发到主节点，主节点写入后等待“所有”从节点拉取该消息，之后通知客户端写入完成
- “所有”节点指的是 In-Sync Replica(ISR)，响应太慢或宕机的从节点会被踢除
- 主节点宕机后，从节点选举成为新的主节点，继续提供服务
- 主节点宕机时正在提交的修改没有做保证（消息可能没有 ACK 却提交了[[2\]](https://lotabout.me/2019/Data-Synchronization-in-Distributed-System/#fn2)）

一些设计细节：

- 当前消费者只能从主节点读取数据，未来可能会改变[[3\]](https://lotabout.me/2019/Data-Synchronization-in-Distributed-System/#fn3)
- 主从的粒度是 partition，每个 broker 对于某些 Partition 而言是主节点，对于另一些而言是从节点
- Partition 创建时，Kafka 会尽量让 preferred replica 均匀分布在各个 broker
- 选主由一个 controller 跟 zookeeper 交互后“内定”，再通过 RPC 通知具体的主节点，此举能防止 partition 过多，同时选主导致 zk 过载。

### Elastic Search

ElasticSearch 对数据的存储需求和 Kafka 很类似，设计也很类似，详细可见[官方文档](https://www.elastic.co/guide/en/elasticsearch/guide/current/distributed-docs.html)。

Master node 的概念，它实际的作用是对**集群状态进行管理**，跟数据的请求无关。

Coordinating node（协调）节点，负责接收客户端请求，分发请求到其他节点，最后再将数据汇集响应给客户端。集群中得任何节点都可以作为协调节点，包括Master Node节点，每个节点都知道任意文档所处的位置。

- 类似于 2PC[[4\]](https://lotabout.me/2019/Data-Synchronization-in-Distributed-System/#fn4)，节点分主从，同步更新消息，除非节点全挂，否则不会丢消息
- 消息发到主节点，主节点写入成功后并行发给从节点，等到从节点全部写入成功，通知客户端写入完成
- 管理节点会维护每个分片需要写入的从节点列表，称为 in-sync copies
- 主节点宕机后，从节点选举成为新的主节点，继续提供服务
- 提交阶段从节点不可用的话，主节点会要求管理节点将从节点从 in-sync copies 中移除

一些设计细节：

- 写入只能通过只主节点进行，读取可以从任意从节点进行[[5\]](https://lotabout.me/2019/Data-Synchronization-in-Distributed-System/#fn5)
- 每个节点均可提供服务，它们会转发请求到数据分片所在的节点，但建议循环访问各个节点以平衡负载
- 数据做分片：`shard = hash(routing) % number_of_primary_shards`
- primary shard 的数量是需要在创建 index 的时候就确定好的
- 主从的粒度是 shard，每个节点对于某些 shard 而言是主节点，对于另一些而言是从节点
- 选主算法使用了 ES 自己的 Zen Discovery[[6\]](https://lotabout.me/2019/Data-Synchronization-in-Distributed-System/#fn6)

### Hadoop

Hadoop 使用的是链式复制，参考 [Replication Pipelining](http://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-hdfs/HdfsDesign.html#Replication_Pipelining)

- 数据的多个复本写入多个 datanode，只要有一个存活数据就不会丢失
- 数据拆分成多个 block，每个 block 由 namenode 决定数据写入哪几个 datanode
- 链式复制要求数据发往一个节点，该节点发往下一节点，待下个节点返回及本地写入成功后返回，以此类推形成一条写入链。
- 写入过程中的宕机节点会被移除 pineline，不一致的数据之后由 namenode 处理。

实现细节：

- 实现中优化了链式复制：block 拆分成多个 packet，节点 1 收到 packet, 写入本地的同时发往节点 2，等待节点 2 完成及本地完成后返回 ACK。节点 2 以此类推将 packet 写入本地及发往节点 3……

### TiKV

TiKV 使用的是 Raft 协议来实现写入数据时的一致性。参考 [三篇文章了解 TiDB 技术内幕——说存储](https://zhuanlan.zhihu.com/p/26967545)

- 使用 Raft，写入时需要半数以上的节点写入成功才返回，宕机节点不超过半数则数据不丢失。
- TiKV 将数据的 key 按 range 分成 region，写入时以 region 为粒度进行同步。
- 写入和读取都通过 leader 进行。每个 region 形成自己的 raft group，有自己的 leader。

### Zookeeper

Zookeeper 使用的是 Zookeeper 自己的 Zab 算法(Paxos 的变种？)，参考 [Zookeeper Internals](https://zookeeper.apache.org/doc/r3.5.5/zookeeperInternals.html)

- 数据只可以通过主节点写入（请求会被转发到主节点进行），可以通过任意节点读取
- 主节点写入数据后会广播给所有节点，超过半数节点写入后返回客户端
- Zookeeper 不保证数据读取为最新，但通过“单一视图”保证读取的数据版本不“回退”



### Cassandra

类似于Multi-master, 可以选择任何一个节点Write/Read, 但是由于很难保证时钟同步，于是不支持transaction.

利于Quorum机制来保证Strong Consistency.

### Summary

如果系统对性能要求高以至于能容忍数据的丢失(Redis)，则显然异步的同步方式是一种好的选择。

而当系统要保证不丢数据，则几乎只能使用同步复制的机制，看到 Kafka 和 Elasticsearch 不约而同地使用了 PacificA 算法（个人认为可以看成是 2PC 的变种），当然这种方法的响应制约于最慢的副本，因此 Kafka 和 Elasticsearch 都有相关的机制将慢的副本移除。

当然看起来 Paxos, Raft, Zab 等新的算法比起 2PC 还是要好的：一致性保证更强，只要半数节点写入成功就可以返回，Paxos 还支持多点写入。只不过这些算法也很难正确实现和优化。

## Reference

- [Transactions Across Datacenters](https://www.youtube.com/watch?v=srOgpXECblk), [slides](https://snarfed.org/transactions_across_datacenters_io.html)
- [分布式系统常见同步机制](https://lotabout.me/2019/Data-Synchronization-in-Distributed-System/)

