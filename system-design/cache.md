# Cache



The main reason we want to use a key-value cache is to reduce latency for accessing active data. Achieve an O(1) read/write performance on a fast and expensive media (like memory or SSD), instead of a traditional O(logn) read/write on a slow and cheap media (typically hard drive).

There are three major factors to consider when we design the cache.

1. Pattern: How to cache? is it read-through/write-through/write-around/write-back/cache-aside?

   - Read
     - Read-aside(最常用)
       - if cache miss, read from db, return to user and update DB simultaneously.
       - 好处: 使用Cache-aside的系统对缓存失效具有一定的resilience。如果缓存集群宕机，系统仍然可以通过直接访问数据库进行操作。
       - 一般client负责处理
     - Read-through
       - 缓存与数据库保持一致。当缓存丢失时，它从数据库加载相应的数据，填充缓存并将其返回给应用程序。
       - 通常由DB或Cache来处理。
     - 如何解决thundring herd问题：当key有cache miss，从数据库取数据的时候，following request o n `GetValue(key)`会等在这里，而不会重复去取。
   - Write-through
     - Only return when the update store in both Cache and DB
     - 与read-through结合使用时，可以保证**read after write**
     - eg: Dynamo DB
   - Write-behind
     - for insert and update, return once add to update queue, async to update cache and DB
     - for delete, sync update to cache and DB
     - pros: 
       - high performance,
       - reduce DB loads, for read is to read from cache instead; for write, in the queue, multiple write are "coalesced" within the write-behind interval. 

2. Placement: Where to place the cache? client-side/distinct layer/server side?

3. Replacement: When to expire/replace the data? LRU/LFU/ARC?

   - LRU(Least Recently Used): check time, and evict the most recently used entries and keep the most recently used ones.

   - LFU(Least Frequently Used): check frequency, and evict the most frequently used entries and keep the most frequently used ones.

   - ARC(Adaptive replacement cache): it has a better performance than LRU. It is achieved by keeping both the most frequently and frequently used entries, as well as a history for eviction. (Keeping MRU+MFU+eviction history.)

Out-of-box choices: Redis/Memcache? Redis supports data persistence while Memcache does not. Riak, Berkeley DB, HamsterDB, Amazon Dynamo, Project Voldemort, etc.



## Reference

- [常用缓存策略的优劣对比](https://www.boxuegu.com/news/2860.html)
- [Read-Through, Write-Through, Write-Behind, and Refresh-Ahead Caching](https://docs.oracle.com/cd/E15357_01/coh.360/e15723/cache_rtwtwbra.htm#COHDG197)

