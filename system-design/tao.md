

# Facebook Tao

[TOC]



最早期的版本和一般的web service一样，cache layer->DB(master-slave), 他的问题是：

- 一般的数据库最重要的性质是spatial locality, 但是这个对于FB的特性并不适用。因为FB最重要的特性是，最新的update更可能会被读到。
- FB还有另一个性质，最多的query是get_all_associate(ids, time_range), 但是kv cache并不是很适用list模型。并且对于其中的一个list update，时间复杂度是O(n). 并且concurrent update会make it more complicated.



## Data Model

Object:

- *object* are typed nodes

Assoc: 

- *associations* are typed directed edges between objects.
- inverted association: association id-> id2, then it will built an inverted association: id2->id1.



eg: 

![object and its association](https://user-images.githubusercontent.com/11788053/129439872-dc1e2537-3ef8-4d6d-b28f-9694223c8f0e.png)

## API

Object

- create/retrieve/update/delete object

Association(bidirectional):

- Write:
  - assoc_add(id1, atype, id2, time, (k→v)*) – Adds or overwrites the association (id1, atype,id2), and its inverse (id1, inv(atype), id2) if defined.
  - assoc_delete(id1, atype, id2) – Deletes the asso- ciation (id1, atype, id2) and the inverse if it exists.
  - assoc_change_type(id1, atype, id2, newtype) – Changes the association (id1, atype, id2) to (id1, newtype, id2), if (id1, atype, id2) exists.


- Read(time-ranged based list query):

  - assoc_get(id1, atype, id2set, high?, low?) – returns all of the associations (id1, atype, id2) and their time and data, where id2 ∈ id2set and high ≥ time ≥ low (if specified). The optional time bounds are to improve cacheability for large asso- ciation lists (see § 5).

  - assoc_count(id1, atype) – returns the size of the association list for (id1, atype), which is the num- ber of edges of type atype that originate at id1.

  - assoc_range(id1, atype, pos, limit) – returns el- ements of the (id1, atype) association list with in- dex *i* ∈ [pos, pos + limit).

  - assoc time range(id1, atype, high, low, limit)

## Architecture

### Storage Layer

DB shards: 

- We divide data into logical *shards*. Each shard is contained in a logical database. Database servers are responsible for one or more shards. 

- Object ids, and its association are stored in same shard



### Caching layer

Request look up the cache server by its ID.

Strategy: LRU

Read:

​	Read or cache missies.

Write:

- Write operations on an association with an inverse may involve two shards, No atomicity on write on the assoc and inverted assoc.
- If a failure occurs the forward may exist without an inverse; these *hanging* associations are scheduled for repair by an asynchronous job.

Sharding:

- Shards are mapped onto cache servers within a tier using consistent hashing
- how to resolve hot spot? *shard cloning*
  -  reads to a shard are served by multiple followers in a tier. 
  - also, good for read failover



#### Leader and follower

Leader: reading from and writing to the storage layer.

	- all write goes to leader: concistency
	- **read after write**: object update in the leader enqueues **invalidation messages** to each corresponding follower. The follower that issued the write is updated **synchronously** on reply from the leader; 
	- serialize concurrent writes that arrive from followers, so it is also ideally positioned to protect the database from **thundering herds**. 

Follower: serve read or forward read misses and writes to a leader.

Cache Maintenance(Eventual Consistency):  asynchronously sending cache maintenance messages rom the leader to the followers. 

- a version num- ber in the cache maintenance message allows it to be ig- nored when it arrives later than **invalidation messages**.



### Scaling Geographically

由于是read-heavy system，于是选择了master-slave的结构。

- write goes from slave leader to master leader
- TAO’s master/slave design ensures that all reads can be satisfied within a single region, at the expense of po- tentially returning stale data to clients. 
- Guarantee read after write: as long as a user consistently queries the same follower tier, the user will typically have a consistent view of TAO state. 





![](https://scontent-lax3-2.xx.fbcdn.net/v/t1.18169-9/1016494_10151647749827200_1788611608_n.png?_nc_cat=106&ccb=1-4&_nc_sid=abc084&_nc_ohc=GxXcMo0WZh4AX-PHsJn&_nc_ht=scontent-lax3-2.xx&oh=26993cd1eb1e2aa86bd97d3012e13be3&oe=613ADEEE)





## Failure tolerance

1. atomicity of association and inverted association

   TAO does not provide atomicity between the two updates. If a failure occurs the forward may exist without an inverse; these hanging associations are scheduled for repair by an asynchronous job.

2. read failure over

![read failover](https://user-images.githubusercontent.com/11788053/129435422-850791b4-552b-4b53-9486-87822395c821.png)





## Summary

- Cache: **write through**, it can support read after write.
- localiity: 
  - Objects are allocated to fixed “shards” via their object ID; these may move across databases, etc after creation 
  - Associations (id1, atype, id2) stored on same shard as **id1**
  - There is inverted association, eg: id1 comment on id2, there will be two association, `comment` and `comment by`
  - One leader cache server is responsible for each shard and its assocs
  - One leader is also ideal for thundering herd problem, it can mediates all of the requests for an id1.
- Two cache tiers
  - Quadratic growth in all-to-all connections for a single tier (because each server has to send writes to the server for their object ID) 
  - single tier is more prone to hot spots 
  - Suitable for read-dominated workload
- All write goes to leader region-> for consistency
  - Eventual consistency: 
    - Object update: follower receive, sync wait for leader response. leader async enqueues invalidation messages to each corresponding follower.
  - Read after write: 

## Reference

- [stanford TAO slides](https://cs.stanford.edu/~matei/courses/2015/6.S897/slides/tao.pdf)
- [TAO: The power of the graph](https://www.facebook.com/notes/10158791581867200/)
- [youtube: USENIX ATC '13 - TAO: Facebook’s Distributed Data Store for the Social Graph](https://www.youtube.com/watch?v=sNIvHttFjdI&t=20s)
- [paper](https://www.usenix.org/system/files/conference/atc13/atc13-bronson.pdf)
