# System Design

## Content

- Overiew: [System Design Interview - Step By Step Guide](https://www.youtube.com/watch?v=bUHFg8CZFws&t=3399s) \[1\]
  - Notes: [step-by-step-guide.md](step-by-step-guide.md)
- Concept
  - [http request/response protocol](basic/http_request_response.md)
  - http://www.ruanyifeng.com/blog/2017/05/websocket.html)
  - [Transaction across datacenters](basic/Transactions-Across-Datacenters.md)
  - [How nodes communicate](node-communicate.md)
- Database, Schema Design
  - Basic
    - [Cache Strategy](cache.md)
    - [Storage Engine](storage-engine.md)
  - NoSQL
    - [Cassandra](cassandra.md)
    - [HBase store msg](hbase-msg.md)
  - Distributed File System
    - (todo)GFSM
    - [Facebook Haystack](haystack.md)
  - [Facebook TAO](tao.md)
  - [Time-series datagase TSDB](ts-db.md)
  - Summary: 几乎所有类型的数据库都用了consistent hashing来做sharding.
- Feature Design
  + Global ranking and user ranking, [Ranking System](ranking.md)
  + Redis, zookeeper实现全局锁
- Large Scale System
  + Netflix long-live push system, [Zuul2: The Netflix Journey to Asynchronous, Non-Blocking Systems](async-non-blocking-system.md)
  + [Slack](slack.md)
  + [Kafka](kafka.md)
  + [Redis: eventual consistency](redis.md)
  + [Zookeeper: strong consistency](zookeeper.md)
  + [Elasticsearch + logstash + kibana ELK](elk.md)
  + [Webhook](webhook.md)

## Foundation
第一步首先要知道基本概念, 比如从Designing data intensive application就很好。但书籍太厚，可以从读别人笔记入手.
Educative最大的问题就是离实践很远。如果出现了不在题库内的问题，就会不知所措。

- Designing data intensive application
  + https://github.com/keyvanakbary/learning-notes/blob/master/books/designing-data-intensive-applications.md
  + https://timilearning.com/posts/ddia/notes/
  + https://www.jianshu.com/nb/28197258
  + [DDIA Seminar](https://www.youtube.com/c/ScottShiCS/playlists)


## Real Industry Solution

How to search for resources

对每家onsite的公司 -> 在地里翻完最近两年所有onsite面筋aggregate所有design题目 -> 对每个题目找工业界实现的blog -> 阅读每个blog，选中最好的一到两个 -> 读到烂熟，整理出我当面试官的话会问的所有问题不停考自己 -> 白板英文自行mock 3遍 -> over

看这些真实案例的技术分享，一方面是学习他们怎么基于当时的情况逐步演化。另一方面，学习他们怎么把这个问题提出，探索，解决，反思的。里面的tradeoff，bottleneck，scalability等等探讨与问题，他们如何present的。我相信，如果我们面试中的presentation能做到跟这些视频的主讲者一样，条理清晰，得分不会低。



[Facebook TAO, Amazon Dynamo, Kafka, RabbitMQ, Zookeeper]这几篇paper建议读一下，可以学到很多实际分布式系统设计的模式，以及如何容错。read repair？consisten hashing? how to choose partition key?master election?这些问题读完这些paper以后就迎刃而解了。对于面试本身也有直接的帮助。比如high scale系统就一定得上NoSQL么？SQL database更改schema很麻烦么？Facebook TAO就是搭载MySQL上的，value就是一个json，schema都是application定的。我就在某一轮面试时说用SQL，面试官追问那schema不好更改怎么办，我就直接说如果expect schema会经常变，用json就行，TAO就是这么干的。面试官表示OK。另外有些公司的面试题就有design key value store，那基本上就按照dynamo来答就好了。

- real industry case sharing
  + [Scaling Facebook Live Videos to a Billion Users](https://youtu.be/IO4teCbHvZw)
  + [Scaling Instagram Infrastructure](https://youtu.be/hnpzNAPiC0E)
  + [Scaling Push Messaging for Millions of Devices @Netflix](https://youtu.be/6w6E_B55p0E)
  + [Zuul](https://github.com/netflix/zuul)
  + [Building a reliable and scalable metrics aggregation and monitoring system](https://youtu.be/UEJ6xq4frEw)
  + [Scaling Slack - The Good, the Unexpected, and the Road Ahead](Scaling Slack - The Good, the Unexpected, and the Road Ahead)
  + [Mastering Chaos - A Netflix Guide to Microservices](https://www.youtube.com/watch?v=CZ3wIuvmHeM)


## Interview process
Tl; dr

时间管理非常重要

Data Mode(schema, relationship, sql/nosql, sharding, replication),  API design, capacity estimation,  high-level architecture design和"other topics"

这些就是系统设计可能会讨论的大方面，提前写好可以保证你在主导对话的同时记得涵盖所有的点，不然自己一直说，可能一兴奋忘了说重点。列出这些点之后我会说"now let's start with API design, what do you think?"然后面试官就会回答好或者不好，你就可以继续说下去了。这里要注意的是，我们其实只是在假装"drive the conversation"，最后先说哪个再说哪个其实还是面试官决定的，除非在少数情况下面试官没有偏好，告诉你先说哪个都行，那你就可以自己随便挑。注意，如果面试官让你自己挑顺序，你也要按照合理的顺序，或者让自己最舒服的顺序，来最大化自己的收益。下面具体说说买个方面应该怎么应对。

### Functional Requirement

搞清楚核心需求，主要的提问时间。这部分争取在5分钟以内解决

- 系统的使用者是谁，是用户，还是其他服务，还是其他系统。
- 如何被使用，是有UI界面，还是定期从哪里拿请求，是提供一个接口。
- REST API: GET/CREATE/DELETE/UPDATE

### Non-functional requirement
次要的提问时间。系统需求之外，有哪些额外的需求需要满足？按照你对系统的理解去offer自己心里觉得最关键的几个点，这部分也是5分钟以内，别展开讨论解决方案，就只是收集信息。

- Consistency - 是否会有数据不一致导致的问题？是否能接受部分数据copy不一致？电商类肯定就要strong consistency，like button能做到read your own write就可以了，后台型的服务eventual足矣。
- Availability - 这个系统要有多可用，1年能承受多少down time？3个9是1年8.8小时(1天88秒），4个9是1年53分钟（1天8秒），5个9是1年5.3分钟（1天1秒不到）。
- Performance/Throughput - 整体有多少用户量，大概率多少活跃用户，此处的用户可以是其他系统。QPS大概有多少，peak factor是多少。
- Read/Write ratio - 读为主还是写为主，用户使用pattern是否是非对称的（比如微博），读写是否有seasonal/recency pattern（比如读的都是最近写的），这个会让后面数据库选型讨论很方便。
- Scalability - 这个和上面的吞吐的区别是，可以展开问面试官希望这个系统有多scalable，因为这个概念本身其实是有一定模糊性的。单机百万websocket和集群数万分布式事务都可以说scalable。
- Misc - 这个就要看具体系统，比如金融交易这些，**idempotency**就很重要，要提到。比如用户facing的系统，latency就最好控制在200ms里（先问latency，别自己主动offer，如果给了一个很大的值，就说明有很大的空间）。
  等等 - 这些就需要自己通过上面的学习来总结，我能想到的上面这些肯定还只是其中一部分。



### Data Model

上面两个提问之后，就是你开始主导了。此处建议先问面试官偏好，说清楚自己接下来要聊data和api然后画图，data可以和API的顺序进行调转。这部分可能要10分钟左右，但是也不要花太多时间。

基于上面的信息，你心里应该大概知道这个系统里核心的几个data model （面向对象编程101）是哪些方面。然后它们之间的互动，简单粗暴可以用一个junction table去存互动的关系（就像一个edge，比如a like b），复杂点的就是transaction（a pay b）。把这些data model对应的“表”写出来，此处你**不用强调它是个sql表还是nosql的结构**，就**只说我们这个data model里有哪些信息。根据你的理解去填基本的那些信息**，比如某一种对象的id，对应的基本信息，根据上面的需求产生的额外信息（比如餐馆系统，用户表里可以放会员状态，打个比方）。

基本data model完了，再聊下他们的**关系表**，此处依然不要说关系表是sql还是nosql，就说某个地方我们放着这个关系。但是如果是比较经典的事务型关系，比如交易，比如订餐，比如发货，并且你了解为什么很多时候用NoSql类系统存这个信息（写为主，读会带着primary key，immutable records，量大，等等），那么就可以主动提一下这一类数据适合放在比如说cassandra里（因为我自己工作里会用到cassandra）。

数据分类的过程中，不用追求完美，哪怕你提前看到过类似系统的设计方案，也不要上来就把n个field都给定义了，或者说你定义某个field要能对应到之前的需求。比如餐馆系统，每2个小时的一张桌子可以算一条记录，但是桌子本身应该有人数或者类型的信息，这个会在谈需求的时候谈到。**表与表之间产生联系（比如foreign key那种）一定要说，而且说到这个联系就可以自然地说indexing和sharding，或者说data locality你要怎么做，然后延展到已经定义的field里有哪些可以用来做composite key，secondary index，sorting key还有sharding key等等**。同样的，别过于深入，这一步重点不是系统的scalability，而是把data model定义清楚，为了后面服务。

### API

这一部分根据面试官偏好，可能非常快，2分钟就过，也可能有些人喜欢和你深入探讨一个好的API设计。但是如果数据部分已经理得很清晰，我是觉得这部分一般最多也就5分钟就能讨论完。

首先都是自己主导，按照需求写出基本的几个API设计，**输入是什么，输出是什么，他们是怎么被call到的（同步，异步，批量）**，以及会返回**那几类的结果**（成功/失败/有去重效果的成功或者失败, retryable error等等）。然后同个API可以根据请求格式的不同有不同的行为。这方面一般来说不会谈太复杂，因为时间关系不可能都涉及到，但是面试官会希望能看到你把已经谈清楚的需求，清晰地反映到API层面，而不是漏掉了一些已经明确的需求（比如job scheduler里能改变优先级这种需求）。

此处一个潜在考点就是你对于API设计的经验积累，能够让一套API围绕一个资源或者目的来服务，而不是说你想到一个API就列举一个API。Google的API设计指南可以帮助你理解一些里面的关键要素。REST API和RPC API的基础知识可能会被考到，所以也最好了解清楚。



### Diagram(重点)

好了，时间过半了，我们终于来到了画图部分。其实在这之前的很多讨论，基本也都已经写在了白板上。此处要重提一下就是，如果1小时面试，上面的4部分消耗20分钟是没问题的，如果是45分钟，那么就必须要压缩，比如functional和nonfunctional放在一起快速过，API部分只提关键核心的部分等等，争取在15分钟这个时间点进入画图部分。

和其他帖子里说的略微不同，我个人觉得，hotspot也好，是否要sharding也好，系统瓶颈分析也好，这些都是只有把图画出来了，明白请求都是怎么hit到哪里的前提下，才好开始分析的。毕竟有时候，你套个cache，或者挂一层LB，实战里就已经很足够，甚至都不用sharding那些。

然后很多同学在这一步可能就会遇到一个很常见的矛盾点 - 假设我知道某个系统大概是这么做的，是否要直接画成那样。我的观点是根据你的理解深度来判断：你是否能够defend自己这个架构？

比如说大家都知道推特类的系统要考虑pull和push，但是其实如果你深入看，类似的系统有做纯pull的，有做纯push的，也有做hybrid的，并没有存在正确答案。push虽然有write amplifcation但是他的实现其实要简单很多，反过来说pull也有一样的看起来美好但是很麻烦的点（不只是单纯数据库压力）。如果你碰巧只看过其中一类的架构图，就原样招呼上去，那么面试官问起来你为什么不用其他的方案，就会非常尴尬。但是如果你的理解广度+深度能到比如这篇分析推特架构的知乎专栏这样，那么你就可以随意把自己知道的正确答案往上招乎了。

### Failure

Tl;dr: Other topics（cache/how to scale/push vs. pull/monitoring/rate limiting/failure handling/logging）

大部分情况下，面试应该就会在上面的部分消耗完时间，但是如果你特别娴熟，飞快地把面试官的各种问题都解决了的话，那么可能他就会问你，是否还有什么别的想谈的。此处其实是一个非常好的信号，说明你已经cover到了他所有的考点，接下来就是纯加分题了。我一般在这个地方，就会开始谈整个系统可能的问题点，然后在一些极端情况下，会产生什么样的问题。也有其他同学觉得可以聊聊logging，monitoring，我本身也可以从系统故障的角度求聊这些。但是你要是有不同的偏好，我觉得是完全okay的，这个也可以算作是Other。

此时就是考验你之前知识积累的广度的时候了（抄答案的时候到了），你如果看过类似的系统的分析和博客，就能想到一些很不错而且不容易考虑到的点（比如对数据做了冷存储的情况下，突然大量访问老数据，又比如某个api突然被ddos导致系统故障）。此处就不太需要担心被人考太深入，不像上面的pull push那样如果准备不足贸然说一个方案反而会被问出问题来。打地鼠一样，多打一个是一个。但是如果你到此刻才发现有些自己应该cover的关键点之前居然没谈到，那就走正式流程把这个点在最后几分钟里好好说清楚，不用慌张，这个时间本来就是作为redundancy来给你的。

## Template

```
Functional Requirement
  Target users (phone, web etc)
  Agree on p0 (mvp) requirements

non-functional requirements
  Consistency or Availability?
  Real-time or batch?
  Read heavy or write heavy?
  QPS
	how much data is queried per request
	can there be spikes in traffic
	performance: expected write-to-read data delay
	             p00 latency for read queries

Database schema
	Data need to store and basic relationship
  How to partition and sharding
  How to build index
  Avoid hotspot
  Avoid consistently updating one record
  Use SQL unless a clear advantage for NoSQL, write heavy, total distributed
  Main/Follower
  Async        vs sync replication
  Write once, read everywhere; write everywhere, read once

High level diagram
API
  Rest vs GraphQL or SOAP
  Talk about throttling, auth


Deep dive
  Scale each component
  Consistency vs Availability
  Components to be considered; how they could help
  DNS
  Load Balancer
  Cache(cdn, data buffer)
    Cacheability. I.e. data repetition.
    Memory requirement
    Staleness effect
```



## Reference
- https://www.1point3acres.com/bbs/thread-776466-1-1.html
- https://www.1point3acres.com/bbs/thread-771667-1-1.html
- [mock interview](https://youtu.be/TUhbXHRLf0o)
- [System Design Interview: Step By Step](https://youtu.be/bUHFg8CZFws)
- [tianpan-system-design](https://tianpan.co/hacking-the-software-engineer-interview/?isAdmin=true)
