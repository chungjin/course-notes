# Elasticsearch + logstash + kibana ELK



### 为什么使用 ELK ？

ELK是一个应用套件，由Elasticsearch、Logstash、Kibana三部分组成，简称ELK。它是一套开源免费、功能强大的日志分析管理系统。ELK可以将我们的系统日志、网站日志、应用系统日志等各种日志进行收集、过滤、清洗、然后进行集中存放并可用于实时检索、分析。这三款软件都是开源软件，通常配合使用，而且又先后归于Elastic.co公司名下，故又称为ELK Stack。那么接下来，我们将从此开源的架构ELK说起，然后一步步推出我们自己的产品Tencent ES系列。下图为ELK Stack的基本组成。



# Prerequisite 



### Logstash

集中、转换和存储数据

Logstash 是开源的服务器端数据处理管道，能够同时从多个来源采集数据，转换数据，然后将数据发送到所选择的DB中，eg: Elasticsearch.

![logstash](https://user-images.githubusercontent.com/11788053/126909959-178a51f1-85a8-4ffe-9ba8-de3b952a69ce.png)

其中，每个部分含义如下：

-  **Shipper**：主要用来收集日志数据，负责监控本地日志文件的变化，及时把日志文件的最新内容收集起来，然后经过加工、过滤，输出到Broker

- **Broker**：相当于日志Hub，用来连接多个Shipper和多个Indexer。这个broker起数据缓存的作用，通过这个缓存器可以提高Logstash shipper发送日志到Logstash indexer的速度，同时避免由于突然断电等导致的数据丢失。[Redis](https://cloud.tencent.com/product/crs?from=10680)服务器是logstash官方推荐的broker，如果数据量大还可以用Kafka, Redis较为昂贵。

- **Indexer**：从Broker读取文本(pull)，经过加工、过滤，输出到指定的介质（可以是文件、网络、elasticsearch等）中。

这里需要说明的是，在实际应用中，LogStash自身并没有什么角色，只是根据不同的功能、不同的配置给出不同的称呼而已，无论是Shipper还是Indexer，始终只做前面提到的三件事。这里需要重点掌握的是logstash中Shipper和Indexer的作用，因为这两个部分是logstash功能的核心，后面会陆续介绍到这两个部分实现的功能细节。

### Elastic Search

[Elastic Search](elasticsearch.md)

### Dashboard: kibana





## Reference

- [南非骆驼说大数据: ELK Stack系列1-7](https://cloud.tencent.com/developer/article/1584012)