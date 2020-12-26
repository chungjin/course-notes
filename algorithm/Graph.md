# Graph

[最短路径—Dijkstra算法和Floyd算法](https://www.cnblogs.com/biyeymyhjob/archive/2012/07/31/2615833.html)


## Dijkstra
单源最短路径算法，用于计算一个节点到其他所有节点的最短路径.主要特点是以起始点为中心向外层层扩展，直到扩展到终点为止。限制条件是，路径的cost**权值为正**，所以才可以用**greedy**算法来做。

设G=(V,E)是一个带权有向图，把图中顶点集合V分成两组，第一组为已求出最短路径的顶点集合（用S表示，初始时S中只有一个源点，以后每求得一条最短路径 , 就将加入到集合S中，直到全部顶点都加入到S中，算法就结束了），第二组为其余未确定最短路径的顶点集合（用U表示），按最短路径长度的递增次序依次把第二组的顶点加入S中。在加入的过程中，总保持从源点v到S中各顶点的最短路径长度不大于从源点v到U中任何顶点的最短路径长度。

1. 初始时，S只包含源点，即`S＝{v}`，v的距离为0。U包含除v外的其他顶点，即:`U={其余顶点}`，若v与U中顶点u有边，则`<u,v>`正常有权值，若u不是v的出边邻接点，则`<u,v>`权值为`∞`。

2. 从U中选取一个距离v最小的顶点k，把k，加入S中（该选定的距离就是v到k的最短路径长度）。

3. 以k为新考虑的中间点，修改U中各顶点的距离；若从源点v到顶点u的距离（经过顶点k）比原来距离（不经过顶点k）短，则修改顶点u的距离值，修改后的距离值的顶点k的距离加上边上的权。

4. 重复步骤2和3直到所有顶点都包含在S中。


![](https://pic002.cnblogs.com/images/2012/426620/2012073019540660.gif)


### Leetcode
- [LC787. Cheapest Flights Within K Stops](https://leetcode.com/problems/cheapest-flights-within-k-stops/)
- [LC743. Network Delay Time](https://leetcode.com/problems/network-delay-time/)
- [LC568. Maximum Vacation Days](https://leetcode.com/problems/maximum-vacation-days/)

## bellman-ford算法,允许负权边的单源最短路径算法

## Floyd
[【最短路径Floyd算法详解推导过程】](https://juejin.im/post/5cc79c93f265da035b61a42e)

1. 算法思想
是解决**任意两点间**的最短路径的一种算法，可以正确处理**有向图或负权**的最短路径问题，同时也被用于计算有向图的传递闭包。Floyd-Warshall算法的时间复杂度为O(N3)，空间复杂度为O(N2)。

是一种动态规划算法。

从任意节点i到任意节点j的最短路径不外乎2种可能，1是直接从i到j，2是从i经过若干个节点k到j。所以，我们假设Dis(i,j)为节点u到节点v的最短路径的距离，对于每一个节点k，我们检查Dis(i,k) + Dis(k,j) < Dis(i,j)是否成立，如果成立，证明从i到k再到j的路径比i直接到j的路径短，我们便设置Dis(i,j) = Dis(i,k) + Dis(k,j)，这样一来，当我们遍历完所有节点k，Dis(i,j)中记录的便是i到j的最短路径的距离。

2. 算法描述
- 从任意一条单边路径开始。所有两点之间的距离是边的权，如果两点之间没有边相连，则权为无穷大。


### Leetcode
- [399. Evaluate Division](https://leetcode.com/problems/evaluate-division/)
