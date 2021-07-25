# Elastic Search

## Elastic Search的使用

参见 [阮一峰：全文搜索引擎Elasticsearch入门教程](https://www.ruanyifeng.com/blog/2017/08/elasticsearch.html)



## Architecture



### Node

①主节点（Master Node）：也叫作主节点，主节点负责index update(create, delete)、sharding, admin job (monitor slave node status)等工作。Elasticsearch 中的主节点的工作量相对较轻。

用户的请求可以发往任何一个节点，并由该节点负责分发请求、收集结果等操作，而并不需要经过主节点转发。

通过在配置文件中设置 node.master=true 来设置该节点成为候选主节点（但该节点不一定是主节点，主节点是集群在候选节点中选举出来的），在 Elasticsearch 集群中只有候选节点能够用zen discovery协议(类Raft)来成为Master.

②数据节点（Data Node）：数据节点，负责数据的存储和相关具体操作，比如索引数据的创建、修改、删除、搜索、聚合。

所以，数据节点对机器配置要求比较高，首先需要有足够的磁盘空间来存储数据，其次数据操作对系统 CPU、Memory 和 I/O 的性能消耗都很大。

通常随着集群的扩大，需要增加更多的数据节点来提高可用性。通过在配置文件中设置 node.data=true 来设置该节点成为数据节点。

③客户端节点（Client Node）：就是既不做候选主节点也不做数据节点的节点，只负责请求的分发、汇总等，也就是下面要说到的协调节点的角色。

​	第一个serve user request的node称之为coordinate node, 它的主要任务是计算文件存在哪个位置，并转发请求，并汇总response, 发送给用户。



### ACID

- Isolation
  Elasticsearch 用乐观并发控制（Optimistic Concurrency Control）来保证新版本的数据不会被旧版本的数据覆盖。

- Consitency: **Quorum**

  写入前先检查有多少个replica可供写入，如果达到写入条件，则进行写操作，否则，Elasticsearch 会等待更多的replica出现，默认为一分钟。

  有如下三种设置来判断是否允许写操作：

  One：只要主分片可用，就可以进行写操作。

  All：只有当主分片和所有副本都可用时，才允许写操作。

  Quorum（k-wu-wo/reng，法定人数）：是 Elasticsearch 的默认选项。当有大部分的分片可用时才允许写操作。其中，对 “大部分” 的计算公式为 int ((primary+number_of_replicas)/2)+1。

- Durability

  首先，当有数据写入时，为了提升写入的速度，并没有数据直接写在磁盘上，而是先写入到内存中，但是为了防止数据的丢失，会追加一份数据到事务日志里(disk)。

  因为内存中的数据还会继续写入，所以内存中的数据并不是以段的形式存储的，是检索不到的。

  JVM->(refresh) OS buffer -> disk, 

  - JVM refresh to OS buffer一秒钟一次，JVM中的数据对user不可见。
  - WAL 写满时(清空log)再做一次referesh, 生成version id.





## Read Write



Write

![write](https://user-images.githubusercontent.com/11788053/126882726-c34f4bc2-2bac-4e13-b18a-bd6ef8da60d9.png)

- coordinator is the first one tha receive user's request, calculate which node contains the 

- 客户端向 Node 1（协调节点）发送写请求。

- Node 1 通过文档的 _id（默认是 _id，但不表示一定是 _id）确定文档属于哪个分片（在本例中是编号为 0 的分片）。请求会被转发到**master replica**所在的节点 Node 3 上。

- Node 3 在主分片上执行请求，如果成功，则将请求并行转发到 Node 1 和 Node 2 的副本分片上。

- 一旦Quorum的副本分片都报告成功（默认），则 Node 3 将向协调节点报告成功，协调节点向客户端报告成功。



Read

- Coordinate node收到请求后，随意从master和slave replica中选择一个，进行转发，收集数据，回复。



## Index in Elastic Search: Lucene

为什么要build index? 如果一个数据有很多tag, 可以很快根据tag检索到数据。并且tag之间还可以做And, Or, Nor等一系列操作。

1. Forward index

   等于创建表。如果是直接写入，可以变成以下的格式：

   ![](https://pic4.zhimg.com/80/v2-e6b81003803254b1d11b3384626c93ab_720w.jpg)

   如果是twitter post这种，可以通过文本操作，得到post id得到关键词。

2. Inverted Index

   通过field得到a list of docid that have this field

   ![](https://pic1.zhimg.com/80/v2-c1cf40e4c4218fd3e992258c08e4e334_720w.jpg)

   key为term, value为posting list.



3. Query

如何Build term index？Trie

假设我们有很多个term，比如：

Carla,Sara,Elin,Ada,Patty,Kate,Selena

如果按照这样的顺序排列，找出某个特定的term一定很慢，因为term没有排序，需要全部过滤一遍才能找出特定的term。排序之后就变成了：

Ada,Carla,Elin,Kate,Patty,Sara,Selena

这样我们可以用二分查找的方式，比全遍历更快地找出目标的term。这个就是 term dictionary。有了term dictionary之后，可以用 logN 次磁盘查找得到目标。但是磁盘的随机读操作仍然是非常昂贵的（一次random access大概需要10ms的时间）。所以尽量少的读磁盘，有必要把一些数据缓存到内存里。但是整个term dictionary本身又太大了，无法完整地放到内存里。于是就有了term index。term index有点像一本字典的大的章节表。比如：

A开头的term ……………. Xxx页

C开头的term ……………. Xxx页

E开头的term ……………. Xxx页

![](https://pic3.zhimg.com/80/v2-e4599b618e270df9b64a75eb77bfb326_720w.jpg)

Trie不会包含所有的term，它包含的是term的一些前缀。通过term index可以快速地定位到term dictionary的某个offset，然后从这个位置再往后顺序查找。比如，Allen在找到A，定位到Ada的offset之后，后面就是Allen的offset. 再加上一些压缩技术（搜索 Lucene Finite State Transducers） term index 的尺寸可以只有所有term的尺寸的几十分之一，使得用**内存缓存**整个term index变成可能。整体上来说就是这样的效果。

Term Dictionary可以用B+树继续优化。



4. **Union index**

   有时候有多个条件需要过滤。例如，age = 18, 女性。于是在得到两个postlist之后，要对他们做合并操作

   - **利用skip list合并**

     ![](https://static001.infoq.cn/resource/image/ea/9f/eafa46683272ff1b2081edbc8db5469f.jpg)

     AND合并这三条postlist. 从最短的list开始。然后从小到大遍历。遍历的过程可以跳过一些元素，比如我们遍历到绿色的 13 的时候，就可以跳过蓝色的 3 了，因为 3 比 13 要小。

     例如目前绿色遍历到13， 在红色和蓝色的skip list中寻找13. 可以把Skip List当成一个BST, 寻找元素的时间复杂度为O(lgn)

     ![](https://static001.infoq.cn/resource/image/a8/34/a8b78c8e861c34a1afd7891284852b34.png)

     利用Frame of reference编码对于postlist进行压缩

     ![](https://static001.infoq.cn/resource/image/9c/b7/9c03d3e449e3f8fb8182287048ad6db7.png)

     

   - **利用Bitset合并**

     如果postlist为`[1,3,4,7,10]`, 对应的 bitset 就是：`[1,0,1,1,0,0,1,0,0,1]`.

     其用一个 byte 就可以代表 8 个文档。所以 100 万个文档只需要 12.5 万个 byte。但是考虑到文档可能有数十亿之多，在内存里保存 bitset 仍然是很奢侈的事情。而且对于个每一个 filter 都要消耗一个 bitset，比如 age=18 缓存起来的话是一个 bitset，18<=age<25 是另外一个 filter 缓存起来也要一个 bitset。

     - 可以很压缩地保存上亿个 bit 代表对应的文档是否匹配 filter；
     - 这个压缩的 bitset 仍然可以很快地进行 AND 和 OR 的逻辑操作。

     Lucene 使用的这个数据结构叫做 Roaring Bitmap。

     ![](https://static001.infoq.cn/resource/image/94/7e/9482b84c4aa3fb77a959c1ead553037e.png)

     其压缩的思路其实很简单。与其保存 100 个 0，占用 100 个 bit。还不如保存 0 一次，然后声明这个 0 重复了 100 遍

   - **Summary**:  Elasticsearch 对其性能有详细的对比（[ https://www.elastic.co/blog/frame-of-reference-and-roaring-bitmaps ](https://www.elastic.co/blog/frame-of-reference-and-roaring-bitmaps)）。简单的结论是：因为 Frame of Reference 编码是如此 高效，对于简单的相等条件的过滤缓存成纯内存的 bitset 还不如需要访问磁盘的 skip list 的方式要快。



## Reference

1. [时间序列数据库的秘密 1,2,3](https://www.infoq.cn/article/database-timestamp-02/?utm_source=infoq&utm_medium=related_content_link&utm_campaign=relatedContent_articles_clk)

2. [掌握它才说明你真正懂 Elasticsearch - ES(三)](https://learnku.com/articles/40400)
3. [Elastic search的使用](https://www.ruanyifeng.com/blog/2017/08/elasticsearch.html)
