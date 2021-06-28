# Zuul 2 : The Netflix Journey to Asynchronous, Non-Blocking Systems

![Netflix Architecture](img/acyn/Netflix-Architecture.png)

## What is Zuul
![Zuul](img/acyn/zuul.png)

Zuul act as gateway, put after the AWS elastic load balancer. Need to route the traffic, filter request, authentication, encrypt, etc.

Zuul core is composed by a bunch of filters, each filter is independent, communicate through the requestContext.

![Request lifecycle](img/acyn/request-lifecycle.png)

## Differences Between Blocking vs. Non-Blocking Systems, From Zuul 1 to Zuul 2

![Zuul1: Multithreaded System Architecture](https://miro.medium.com/max/1000/0*kPzgZrACokyPJJfy.png)
Zuul 1 was built on the Servlet framework. Such systems are blocking and multithreaded, which means they process requests by using one thread per connection. I/O operations are done by choosing a worker thread from a thread pool to execute the I/O, and the request thread is **blocked** until the worker thread completes.

Problem:
- Estimate Capacity: 100 concurrent connections each instance.
  + cost of each connection: a thread (with heavy memory and system overhead, like ctx switch)
- high latency, especially when error happens and retry, eg, backend latency increases or device retries due to errors, the count of active connections and threads increases.
  + Zuul 1 solution: bandwidth throttling.

![Zuul2: Asynchronous and Non-blocking System Architecture](https://miro.medium.com/max/1000/0*jrG2ldEVRRJcgpkj.png)

One thread per CPU core handling all requests and responses.
- Estimate Capacity: 10 thousand concurrent connections each instance.
  + cost of each connection: a file descriptor, and the addition of a listener.

## Which one to choose
CPU bound task: zuul1
- Highly CPU-bound work loads
- Desire operational simplicity
- Desire development simplicity
- Run legacy systems that are blocking

I/O bound task: zuul2
- Highly I/O bound workloads, most time is waiting for response
- Long requests and large files
- Streaming data from queues
- Massive amounts of connections


## Reference
- [Announcing Zuul: Edge Service in the Cloud](https://netflixtechblog.com/announcing-zuul-edge-service-in-the-cloud-ab3af5be08ee)
- [Zuul 2 : The Netflix Journey to Asynchronous, Non-Blocking Systems](https://netflixtechblog.com/zuul-2-the-netflix-journey-to-asynchronous-non-blocking-systems-45947377fb5c)
- [Zuul's Journey to Non-Blocking" by Arthur Gonigberg](https://www.youtube.com/watch?v=2oXqbLhMS_A)
