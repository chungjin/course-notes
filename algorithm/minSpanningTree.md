# minimum spanning tree

在一给定的无向图 G = (V, E) 中，(u, v) 代表连接顶点 u 与顶点 v 的边（即 (u, v) 属于 E），而 w(u, v) 代表此边的权重，若存在 T 为 E 的子集， 且 (V, T) 为树，使得 sum W(T) 的 w(T) 最小，则此 T 为 G 的最小生成树。

![](https://zh.wikipedia.org/wiki/File:Minimum_spanning_tree.svg)

- leetcode: [lc 1135 Connecting Cities With Minimum Cost](https://leetcode.com/problems/connecting-cities-with-minimum-cost/)

## prime 时间复杂度: O(ElogV)
从单一顶点开始，普里姆算法按照以下步骤逐步扩大树中所含顶点的数目，直到遍及连通图的所有顶点。

1. 输入：一个加权连通图，其中顶点集合为V，边集合为E；
2. 初始化：Vnew = {x}，其中x为集合V中的任一节点（起始点），Enew = {}；
3. 重复下列操作，直到Vnew = V：
    1. 在集合E中选取**权值最小**的边（u, v），其中u为集合Vnew中的元素，而v则是V中没有加入Vnew的顶点（如果存在有多条满足前述条件即具有相同权值的边，则可任意选取其中之一）；
    2. 将v加入集合Vnew中，将（u, v）加入集合Enew中；
4. 输出：使用集合Vnew和Enew来描述所得到的最小生成树。

代码实现：Priority Queue
```python
class Solution:
    '''Connecting Cities with Minimum Cost == Find Minimum Spanning Tree'''
    def minimumCost(self, N: int, connections: List[List[int]]) -> int:
        '''
        Prim's Algorithm:
        1) Initialize a tree with a single vertex, chosen
        arbitrarily from the graph.
        2) Grow the tree by one edge: of the edges that
        connect the tree to vertices not yet in the tree,
        find the minimum-weight edge, and transfer it to the tree.
        3) Repeat step 2 (until all vertices are in the tree).
        '''
        # city1 <-> city2 may have multiple different cost connections,
        # so use a list of tuples. Nested dict will break algorithm.
        G = collections.defaultdict(list)
        for city1, city2, cost in connections:
            G[city1].append((cost, city2))
            G[city2].append((cost, city1))

        queue = [(0, N)]  # [1] Arbitrary starting point N costs 0.
        visited = set()
        total = 0
        while queue and len(visited) < N: # [3] Exit if all cities are visited.
            # cost is always least cost connection in queue.
            cost, city = heapq.heappop(queue)
            if city not in visited:
                visited.add(city)
                total += cost # [2] Grow tree by one edge.
                for edge_cost, next_city in G[city]:
                    heapq.heappush(queue, (edge_cost, next_city))
        return total if len(visited) == N else -1
```



## Kruskal 时间复杂度: O(ElogV)
1. 新建图G，G中拥有原图中相同的节点，但没有边
2. 将原图中所有的边按权值从小到大排序
3. 从权值最小的边开始，如果这条边连接的两个节点于图G中不在同一个连通分量中，则添加这条边到图G中
4. 重复3，直至图G中所有的节点都在同一个连通分量中

代码实现：union and find

```python
class Solution:
    def minimumCost(self, N: int, connections: List[List[int]]) -> int:

        def find(city):
            while parent[city]!=city:
                parent[city] = parent[parent[city]]
                city = parent[city]
            return city

        def union(c1, c2):
            root1, root2 = find(c1), find(c2)
            if root1 == root2:
                return False
            parent[root2] = root1
            return True



        # keeps the root for each node, initially is itself
        parent = {city: city for city in range(1, N+1)}

        connections.sort(key=lambda x:x[2])
        total = 0
        for city1, city2, cost in connections:
            if union(city1, city2): #belongs to different component
                total += cost

        # check that all cities are connected
        root = find(N)
        return total if all(root == find(city) for city in range(1, N+1)) else -1
```

### Related question
- amazon OA: [Min Cost to Connect All Nodes](https://leetcode.com/discuss/interview-question/356981)  
    题解： 这题只能用Kruskal而不能用Prim, 因为prime的原理是，不停地把新的节点加入唯一的一个component中。而Kruskal则是，我们可能会在中途形成很多的component, 每次选取的边一定是：可以连接两个component并且权值最小。
    ```python
    import collections

    '''
    Given an undirected graph with n nodes labeled 1..n. Some of the nodes are already connected. The i-th edge connects nodes edges[i][0] and edges[i][1] together. Your task is to augment this set of edges with additional edges to connect all the nodes. Find the minimum cost to add new edges between the nodes such that all the nodes are accessible from each other.

    Input:

    n, an int representing the total number of nodes.
    edges, a list of integer pair representing the nodes already connected by an edge.
    newEdges, a list where each element is a triplet representing the pair of nodes between which an edge can be added and the cost of addition, respectively (e.g. [1, 2, 5] means to add an edge between node 1 and 2, the cost would be 5).

    Example 1:
    Input: n = 6, edges = [[1, 4], [4, 5], [2, 3]], newEdges = [[1, 2, 5], [1, 3, 10], [1, 6, 2], [5, 6, 5]]
    Output: 7
    Explanation:
    There are 3 connected components [1, 4, 5], [2, 3] and [6].
    We can connect these components into a single component by connecting node 1 to node 2 and node 1 to node 6 at a minimum cost of 5 + 2 = 7.
    '''

    '''
    union and find
    '''
    def minCost(n, edges, newEdges):
        def find(city):
            while parent[city]!=city:
                parent[city] = parent[parent[city]]
                city = parent[city]
            return city

        def union(c1, c2):
            # return true if it is belongs to different component
            r1, r2 = find(c1), find(c2)
            if r1!=r2:
                parent[r1] = r2
                return True



        # sort by its cost
        newEdges.sort(key=lambda x:x[2])
        parent = {city:city for city in range(1, n+1)}
        res = 0
        for i, j in edges:
            union(i, j)
        for c1, c2, cost in newEdges:
            if union(c1, c2):
                res += cost
        r = find(1)
        return res if all(find(city) == r for city in range(2, n+1)) else -1




    def main():
        n = 6
        edges = [[1, 4], [4, 5], [2, 3]]
        newEdges = [[1, 2, 5], [1, 3, 10], [1, 6, 2], [5, 6, 5]]
        res = minCost(n, edges, newEdges)
        print(res)

    main()
    ```



## Reference
- [wiki: 最小生成树](https://zh.wikipedia.org/wiki/%E6%9C%80%E5%B0%8F%E7%94%9F%E6%88%90%E6%A0%91)
- [wiki: prime](https://zh.wikipedia.org/wiki/%E6%99%AE%E9%87%8C%E5%A7%86%E7%AE%97%E6%B3%95)
- [wiki: Kruskal](https://zh.wikipedia.org/wiki/%E5%85%8B%E9%B2%81%E6%96%AF%E5%85%8B%E5%B0%94%E6%BC%94%E7%AE%97%E6%B3%95)
***
