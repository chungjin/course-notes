# Zookeeper

[TOC]

水平伸缩的方式来提升性能成为了主流。但在分布式架构下，当服务越来越多，规模越来越大时，对应的机器数量也越来越大，单靠人工来管理和维护服务及地址信息会越来越困难， 单点故障的问题也开始凸显出来，一旦服务路由或者负载均衡服务器宕机，依赖他的所有服务均将失效。

此时，需要一个能够动态注册和获取服务信息的地方。来统一管理服务名称和其对应的服务器列表信息，称之为服务配置中心，

ZooKeeper is a centralized service for maintaining configuration information, naming, providing distributed synchronization, and providing group services.

作为分布式数据库，要求Strong Consistency, 于是引入了2PC和Paxos



## 特性

- Tolerance of Single Point Failure: 

  如果要满足这样的一个高性能集群，我们最直观的想法应该是，每个节点都能接收到请求，并且每个节点的数据都必须要保持一致。要实现各个节点的数据一致性，就势必要一个 leader 节点负责协调和数据同步操作。这个我想大家都知道，如果在这样一个集群中没有 leader 节点，每个节点都可以接收所有请求，那么这个集群的数据同步的复杂度是非常大。
  结论：所以这个集群中涉及到数据同步以及会存在leader 节点，zookeeper 用了基于 paxos 理论所衍生出来的 ZAB 协议，来进行leader的选举和数据备份。

- Strong consistency: 2PC

  leader 节点如何和其他节点保证数据一致性，并且要求是强一致的。在分布式系统中，每一个机器节点虽然都能够明确知道自己进行的事务操作过程是成功和失败，但是却无法直接获取其他分布式节点的操作结果。所以当一个事务操作涉及到跨节点的时候，就需要用到分布式事务，分布式事务的数据一致性协议有 2PC 协议和3PC 协议。



## Read/Write

![](https://img-blog.csdnimg.cn/20190116112438533.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3p5MzQ1MjkzNzIx,size_16,color_FFFFFF,t_70)

在 zookeeper 中，客户端会随机连接到 zookeeper 集群中的一个节点，如果是读请求，就直接从当前节点中读取数据，如果是**写**请求，那么请求会被转发给**leader** 提交事务，然后 leader 会广播事务，只要有**超过**半数节点写入成功，那么写请求就会被提交（类 2PC 事务）。

## Leader/Follower

集群角色
Leader 角色
       Leader 服务器是整个 zookeeper 集群的核心，主要的工作任务有两项

- 事物请求的唯一调度和处理者，保证集群事物处理的顺序性
- 集群内部各服务器的调度者

Follower 角色

- 处理客户端非事物请求、转发事物请求给 leader 服务器

- 参与事物请求 Proposal 的投票（需要半数以上服务器通过才能通知 leader commit 数据; Leader 发起的提案，要求 Follower 投票）

- 参与 Leader 选举的投票 

- 备份Leader数据, 做leader的failover, 只有保存有up-to-date的信息的follower能被选为leader



## Zookeeper Data Model

<img src="https://images0.cnblogs.com/blog/671563/201411/301534562152768.png" style="zoom:67%;" />

ZooKeeper拥有一个层次的naming space，这个和标准的文件系统非常相似.

**(1) 引用方式**

Zonde通过**路径引用**，如同Unix中的文件路径。路径必须是**absolute path**，因此他们必须由斜杠字符来**开头**。除此以外，他们必须是唯一的，也就是说每一个路径只有一个表示，因此这些路径不能改变。在ZooKeeper中，路径由Unicode字符串组成，并且有一些限制。字符串"/zookeeper"用以保存管理信息，比如关键配额信息。

**(2)** **Znode结构**

ZooKeeper命名空间中的Znode，兼具文件和目录两种特点。既像文件一样维护着数据、元信息、ACL、时间戳等数据结构，又像目录一样可以作为路径标识的一部分。图中的每个节点称为一个Znode。 每个Znode由3部分组成:

**①** stat：此为状态信息, 描述该Znode的版本, 权限等信息

**②** data：与该Znode关联的数据

**③** children：该Znode下的子节点

ZooKeeper虽然可以关联一些数据，但并没有被设计为常规的数据库或者大数据存储，相反的是，它用来**管理调度数据**，比如分布式应用中的配置文件信息、状态信息、汇集位置等等。这些数据的共同特性就是它们都是很小的数据，通常以KB为大小单位。ZooKeeper的服务器和客户端都被设计为严格检查并限制每个Znode的数据大小至多1M，但常规使用中应该远小于此值。

**(3) 数据访问**

ZooKeeper中的每个节点存储的数据要被**atomicity的操作**。也就是说读操作将获取与节点相关的所有数据，写操作也将替换掉节点的所有数据。另外，每一个节点都拥有自己的ACL(访问控制列表)，这个列表规定了用户的权限，即限定了特定用户对目标节点可以执行的操作。

**(4) 节点类型**

ZooKeeper中的节点有两种，分别为**临时节点**和**永久节点**。节点的类型在创建时即被确定，并且不能改变。

**① 临时节点：**该节点的生命周期依赖于创建它们的会话。一旦会话(Session)结束，临时节点将被自动删除，当然可以也可以手动删除。虽然每个临时的Znode都会绑定到一个客户端会话，但他们对所有的客户端还是可见的。另外，ZooKeeper的临时节点不允许拥有子节点。

**② 永久节点：**该节点的生命周期不依赖于会话，并且只有在客户端显示执行删除操作的时候，他们才能被删除。

**(5)** **顺序节点**

当创建Znode的时候，用户可以请求在ZooKeeper的路径结尾添加一个**递增的计数**。这个计数**对于此节点的父节点来说**是唯一的，它的格式为"%10d"(10位数字，没有数值的数位用0补充，例如"0000000001")。当计数值大于232-1时，计数器将溢出。

**(6) 观察**

客户端可以在节点上设置watch，我们称之为**监视器**。当节点状态发生改变时(Znode的增、删、改)将会触发watch所对应的操作。当watch被触发时，ZooKeeper将会向客户端发送且仅发送一条通知，因为watch只能被触发一次，这样可以减少网络流量。

**(7) 操作: Optimistic Lock**

新ZooKeeper操作是有限制的。delete或setData必须明确要更新的Znode的版本号，我们可以调用exists找到。如果版本号不匹配，更新将会失败。

更新ZooKeeper操作是非阻塞式的。因此客户端如果失去了一个更新(由于另一个进程在同时更新这个Znode)，他可以在不阻塞其他进程执行的情况下，选择重新尝试或进行其他操作。



## 应用：分布式锁解决集群中的单点故障

1. **Master启动**

   在引入了Zookeeper以后我们启动了两个主节点，"主节点-A"和"主节点-B"他们启动以后，都向ZooKeeper去注册一个节点。我们假设"主节点-A"锁注册地节点是"master-00001"，"主节点-B"注册的节点是"master-00002"，注册完以后进行选举，编号最小的节点将在选举中获胜获得锁成为主节点，也就是我们的"主节点-A"将会获得锁成为主节点，然后"主节点-B"将被阻塞成为一个备用节点。那么，通过这种方式就完成了对两个Master进程的调度。

   ![](https://images0.cnblogs.com/blog/671563/201411/301535008567950.png)

2. **Master故障**
   如果"主节点-A"挂了，这时候他所注册的节点将被自动删除，ZooKeeper会自动感知节点的变化，然后再次发出选举，这时候"主节点-B"将在选举中获胜，替代"主节点-A"成为主节点。
   ![](https://images0.cnblogs.com/blog/671563/201411/301535008567950.png)

3. **Master 恢复**
   如果主节点恢复了，他会再次向ZooKeeper注册一个节点，这时候他注册的节点将会是"master-00003"，ZooKeeper会感知节点的变化再次发动选举，这时候"主节点-B"在选举中会再次获胜继续担任"主节点"，"主节点-A"会担任备用节点。
   ![](https://images0.cnblogs.com/blog/671563/201411/301535016997293.png)

## Reference

-  [分布式系统基础--注册中心Zookeeper](https://blog.csdn.net/zy345293721/category_9284922.html)
- [ZooKeeper学习第一期---Zookeeper简单介绍](https://www.cnblogs.com/sunddenly/p/4033574.html)

