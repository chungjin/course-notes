# Step By Step Guide

[TOC]

<img src="https://user-images.githubusercontent.com/11788053/124994796-3b52ce80-dffb-11eb-994e-a14e7f1b917c.png" alt="Technology Stack" style="zoom:50%;" />





## Technology Stack

<img src="https://user-images.githubusercontent.com/11788053/124996683-3c392f80-dffe-11eb-98b3-c420d199547b.png" alt="Technology Stack" style="zoom:50%;" />



## Interview Process

### Functional Requirement - API (5 mins)

API

input: ProcessEvents(Id, eventtype, timerange)

- [extend] event: like, views, comments

- [extend] time range

How it will be used? real-time, or periodical system, eg: p99 latency, Write-to-read data delay

### Non-functional requirements - Qualities (5 mins)

consistency: eventually consistency

QPS

latency: real-time, or periodical system, eg: p99 latency, Write-to-read data delay

avaialbility: important,

read/write ratio: both heavily
Scalability

Traffic Spike

MISC - 这个就要看具体系统，比如金融交易这些，**idempotency**就很重要，要提到。比如用户facing的系统，latency就很重要



### Data Model

Design Princple:

- Scalability: scale write/read
- Delay: make both writes and read fast
- Consistency: what level of consistency
- Durability: not lose data in case of hardware faults and network partitions, and how to recover
- Extensible for data model changes in the future

#### What we store
![data model](https://user-images.githubusercontent.com/11788053/125008859-3ea88300-e018-11eb-93a6-1932f14eca4a.png)





#### Prerequisite: SQL vs Nosql

##### SQL


###### Performance finetune

Focus on scalability, availability(replica)

<img src="https://user-images.githubusercontent.com/11788053/125004643-d1dcbb00-e00e-11eb-81ba-82cb2b725f86.png" alt="SQL" style="zoom:67%;" />
**Load balancer**: Zookeeper provide the registry service

- monitor health of master, replica, handle master failover

**Replica** provide data replication and increase read throuput

**Shard proxy** provide cache of the query, monitor instance health, publish metrics, terminate queries that take too long to return

###### Relational

![relational](https://user-images.githubusercontent.com/11788053/125006515-300b9d00-e013-11eb-95f6-e27db0437b6e.png)

SQL can store the relation between the table, and decrease data duplication, avoid inconsistency data.



##### NoSQL: Cassandra

**Feature**: fault-tolerant, scalable(read/write throuput increase linearly as new machines are added), multi datacenter replication, works well with time-series data.



**Gossip protocol**: every second, the node share the info with 3 nodes, so **each node about the info of whole cluster**. -> don't need proxy

**Quorum read/write**, quorum level is configurable.

<img src="https://user-images.githubusercontent.com/11788053/125006119-3d745780-e012-11eb-899b-6506efaacefa.png" alt="NoSQL" style="zoom:67%;" />



Once each server query to the node(maybe the nearest, the fastest), the node will return the value that can achieve the quorum.

###### Wide Column

The field in cassandra support thread safe list, map, continuing append new data into the list.

![cassandra](https://user-images.githubusercontent.com/11788053/125006651-7cef7380-e013-11eb-9aad-bc1279b9a671.png)



#### Data Aggregation

##### Pre-processing?

![preprocess](https://user-images.githubusercontent.com/11788053/125007399-f6d42c80-e014-11eb-9e8f-11b460a12f78.png)

First one will contineouing update the DB once it receive the query.

Second one, store the data in memory, but update DB later(single point failure)

##### Push or Pull

<img src="https://user-images.githubusercontent.com/11788053/125007513-33078d00-e015-11eb-878b-fd309d0838a0.png" alt="push or pull" style="zoom:70%;" />

1. Push data into server, but if server failed, the data may lose
2. server pull data from MQ(message queue), so can recalculate it(with offset)

##### Partitioning(Scalability):

<img src="https://user-images.githubusercontent.com/11788053/125008174-d0af8c00-e016-11eb-9d0c-ca8c27822940.png" alt="MQ partition" style="zoom:80%;" />

### High Level Design - Diagram (重点)



#### Processing service

<img src="https://user-images.githubusercontent.com/11788053/125009530-8d0a5180-e019-11eb-9ebb-89b21db4388d.png" alt="processing service" style="zoom:80%;" />

Consumer will have a **cache**(with lease) to store last 10 mins events, so as to deduplicate the event.

**Aggregator** helps to calculate the group the events, and calculate the sum.

**Internal Queue** is to decouple the aggregator to the db, async way to write to the database, and improve the **throughput**(mutliple thread in Database writer).

**Dead-letter Queue** can help when database is unreachable or slow, or **database writer can store locally(in disk)**

**State Store**: save the state, as a checkpoint, to avoid reproduce a lot data if there is failure.



#### Data Ingestion Path

<img src="https://user-images.githubusercontent.com/11788053/125013726-33a62080-e021-11eb-95de-aa98b7b9b7f5.png" alt="Data Ingestion Path" style="zoom:67%;" />



Strategy to choose from: pros and cons

![Ingestion](https://user-images.githubusercontent.com/11788053/125014008-b5964980-e021-11eb-8053-986f3f408b4b.png)



##### Partitioner Service Client

Blocking vs non-blocking I/O:

- Blocking: one thread per request, synced, easy to debug/track, lower throuput.
- Unblocking, multiplex: one thread pile up the request, increase the complexity of the system

Buffering and batching: gateway, reduce the traffic, improve the efficiency. But it increase complexity, and data may be unrecoverable if batch processing failed in the middle.

Timeout: Connection timeout and request timeout:

 - Retry after timeout
   	- Change to another available machine
   	- Retry storm event when a lot client retry at the same time, it will overload server
    - Improvement: 
      	- **Exponential backoff and Jitter(random increase interval)**.
       - Circuit breaker: stop a client from repeatedly trying to execute a operation that is likely to fail. Theory is calculate how many request is failed **recently(a timer)**, if exceed the threshold, stop service.
         	- drawbacks: harder to test, hard to set properly error threshold and timers.



##### Load Balancer

- Hardware
  - pros: powerful, can handle big throuput
  - Cons: Expensive

- Software
  - pros: open-source, free, a lot type/algorithm to choose from
  - TCP: cannot inspect the pkg,
  - HTTP: can check the content, make decision upon that.
  - Algorithm: round robin, lowest live connection, server with the fastest response time, hash-based.

- DNS
  - paritioner service is registerd to DNS
  - Monitor the health of paritioner
  - Primary-Secondary node, reside in different data-center, to improva the availability. Primary take request, secondary work as backup.



##### Partitioner service and Partitions

Hot partition, with unbalance load 

	- hash by video id, event time is not a good choice.
	- split hot partition into multiple new partition.
	- allocated dedicated paritition just for some hot video.

Service Discovery -> zookeeper

	- maintain the info of replication
	- monitor the health, replace by backup.

Replication: leader and follower.

Message format

	- textual formats, xml, csv, json, simple and readable, but it waste some space to encode the field name
	- Binary format, save the space to encode the filed name, but server, client side should be aware of the data strucure in advance.



#### Data retrieval path: design for extension

<img src="https://user-images.githubusercontent.com/11788053/125020969-a073e780-e02e-11eb-9ed0-441ff6a83bdf.png" alt="data retrieve path" style="zoom:67%;" />

Extension:

	- Save the hisotry data for analysis use
	- CDN speed up
	- Pre-cache some results with lease.



Save the hisotry data for analysis use

- Improve the granduality of the table, for time serires data.

  hour-> day-> month-> year

- store the **cold data** in object storage, like AWS S3. Cold data store the archived and infrequently accessed data.



#### Data Flow simulation

![data flow simulation](https://user-images.githubusercontent.com/11788053/125021130-e630b000-e02e-11eb-81ac-4455c9e45e97.png)



## Reference

- [System Design Interview – Step By Step Guide](https://www.youtube.com/watch?v=bUHFg8CZFws&t=4235s)



