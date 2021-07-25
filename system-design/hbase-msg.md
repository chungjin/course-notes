# HBase store msg

[TOC]



## HBase data model

![](https://img-blog.csdnimg.cn/20210528231817288.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQxNzA5NTc3,size_16,color_FFFFFF,t_70)

Table is sorted in HBase, sorted by the row_key.[1]

## Terminology

1. Name Space（相当于mysql的库

 命名空间，类似于关系型数据库的DataBase概念，每个命名空间下有多个表。HBase有两个自带的命名空间，分别是“hbase”和“default”，“hbase”中存放的是HBase内置的表，“default”表是用户默认使用的命名空间。

2. Region（为了实现数据分布式存储

 类似于关系型数据库的表概念。不同的是，HBase定义表时只需要声明Column Family即可，不需要声明具体的列。这意味着，往HBase写入数据时，字段可以动态、按需指定。因此，和关系型数据库相比，HBase能够轻松应对字段变更的场景。

3. Row（指逻辑结构的一行 eg：rowkey 1001有多行，但是逻辑结构就1行

 HBase表中的每行数据都由一个RowKey和多个Column（列）组成，数据是按照RowKey的字典顺序存储的，并且查询数据时只能根据RowKey进行检索，所以RowKey的设计十分重要。

4. Column

 HBase中的每个列都由Column Family(列族)和Column Qualifier（列限定符）进行限定，例如info：name，info：age。建表时，只需指明列族，而列限定符无需预先定义。

5. Time Stamp

 用于标识数据的不同版本（version），每条数据写入时，如果不指定时间戳，系统会自动为其加上该字段，其值为写入HBase的时间。

6. Cell

 由{rowkey, column Family：column Qualifier, time Stamp} 唯一确定的单元。

 cell中的数据是没有类型的，全部是字节码形式存储。



## Architecture

![](https://img-blog.csdnimg.cn/20210528231906782.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQxNzA5NTc3,size_16,color_FFFFFF,t_70)

store的个数 = column family的个数

## HBase存储Message信息

Facebook Messenger采用HBase来存储信息。



### Rowkey Format

会话ID和消息ID采用snowflake算法生成，RowKey包括了三部分内容。

*会话hash值 | 会话id | 逆序消息id*  

- 会话hash值的目的为数据分区（region）存储，预分区能够分摊数据读写压力；

- 会话id确定唯一会话，一个**群**里的所有消息拥有相同的会话id；

- 逆序消息id可以优先显示最新消息。当往上面拉的时候，才显示更多的消息。



### Region Partition Design

同一会话的消息，一般会集中读取（用户查看某个聊天的消息就是这种场景）。因此需要把同一会话的消息存储在一个分区。我们采用会话id的hash值来做分区字段, 能够确保同一会话的消息一定在同一分区。



### Read & Write

1、依照上述设计格式，我们用传参后的会话Id，002|***|***取模128—以此分散到不同的region；

2、确定具体region后依照rowkey的后续***|312312312312312312312|***的会话Id确定唯一的会话；

3、确定唯一会话后依照rowkey的后续消息Id确定某一个具体消息***|***|8896232141957373907，注意这个消息Id已经被逆序处理（Long.MAX_VALUE-消息Id），用来做拉取最邻近的消息。 

## Reference

1. [聊聊HBase-1](https://blog.csdn.net/qq_41709577/article/details/117375943?utm_medium=distribute.pc_relevant.none-task-blog-baidujs_title-0&spm=1001.2101.3001.4242)
2. [HBase存储IM消息，RowKey该怎么设计？](https://cloud.tencent.com/developer/article/1525569)