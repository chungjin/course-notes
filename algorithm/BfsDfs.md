# DfsBfs

- [DFS BFS](#dfsbfs)
  * [DFS 和 BFS 的比较](#dfs-和-bfs-的比较)
  * [Leetcode](#leetcode)
    + [遍历搜索类](#遍历搜索类)
    + [有向图cycle detect](#有向图cycle-detect)
    + [无向图Detect Cycle](#无向图detect-cycle)
    + [有向图Topological Sort](#有向图topological-sort)
  * [Topological Sorting](#topological-sorting)
  * [Tarjan's strongly connected components algorithm - DFS](#tarjan-s-strongly-connected-components-algorithm---dfs)

LeetCode 上很多问题都可以抽象成 “图” ，比如搜索类问题，树类问题，迷宫问题，矩阵路径问题，等等。

## DFS 和 BFS 的比较
- BFS 的时间空间占用以 branching factor 为底， 到解的距离 d 为指数增长；空间占用上 Queue 是不会像 DFS 一样只存一条路径的，而是从起点出发越扩越大，因此会有空间不够的风险，空间占用为 O(b^d)。(其中d为到解的深度，b为每个节点的节点的数目).
- DFS 的时间占用以 branching factor 为底，树的深度 m 为指数增长；而空间占用上，却只是 O(bm)，可视化探索过程中只把每个 Node 的所有子节点存在 Stack 上， 探索完了再 pop 出来接着探，因此储存的节点数为 O(bm)。
- 时间复杂度均为O(n), n为所有节点的数目


## Leetcode

### 遍历搜索类
DFS, BFS 均可，一般采用DFS，在空间复杂度上较低，并且写起来比较简单。  
[79. Word Search](https://leetcode.com/problems/word-search/)  
[200. Number of Islands](https://leetcode.com/problems/number-of-islands/)  
[130. Surrounded Regions](https://leetcode.com/problems/surrounded-regions/)
### 有向图cycle detect
  - DFS
    - 暴力解法：DFS + Backtracking，寻找“所有从当前节点的” path，如果试图访问 visited 则有环；缺点是，同一个点会被探索多次，而且要从所有点作为起点保证算法正确性，时间复杂度非常高
    - 最优解法是 CLRS 上用三种状态表示每个节点：
      - "0" 还未访问过;
      - "1" 代表正在访问；
      - "2" 代表已经访问过；
    - DFS 开始时把当前节点设为 "1";
    - 在从任意节点出发的时候，如果我们试图访问一个状态为 "1" 的节点，都说明图上有环。
    - 当前节点的 DFS 结束时，设为 "2";
    - 在找环上，DFS 比起 BFS 最大的优点是"对路径有记忆"，DFS 记得来时的路径和导致环的位置，BFS 却不行。
  - BFS
    - 扫一遍所有 edge，记录每个节点的 indegree.
    - 在有向无环图中，一定会存在至少一个 indegree 为 0 的起点，将所有这样的点加入queue。
    - 依次处理queue里的节点，把每次poll出来的节点的 children indegree -1. 减完之后如果 child 的 indegree = 0 了，就也放入队列。
    - 如果图真的没有环，可以顺利访问完所有节点，如果还有剩的，说明图中有环，因为环上节点的 indegree 没法归 0.

  [207. Course Schedule](https://leetcode.com/problems/course-schedule/description/)

### 无向图Detect Cycle
  - DFS
    - 依然记录每个点的状态，0 代表“未访问”；1 代表“访问中”；2 代表“已访问”；
    - DFS call里面要传入prev节点这个参数，避免出现原路返回，或者回到前一个节点误判为有环。(和directed graph DFS唯一的不同之处)。
    - 其他情况下，如果我们试图访问一个状态为 “1” 的节点，都可以说明图中有环。
  - BFS
     + 方法1: 一层一层的扫，并且当前层结束时，把当前层所有的点再iterate一遍，全部标记为已访问。避免扫到下一层的时候，寻找相邻点，会误判有环。
       - 初始化标记所有点的状态为0.
       - 随便扔一个点进 queue，标记 "1"，然后 BFS，所有 child = "0" 的都加入队列，队列中的点都标记为1.
       - 当 node 的所有 child 点都检查完并加入queue后，立刻把当前 node = 2，不然下一层 BFS 会回头去看自己然后误报。
       - 如果遇到 child = "1" 的说明有环
    + 方法2: 在常规BFS基础上，记录访问次序，比如`a->b, parent[b] = a`. 下一次从`c->b`,如果`c!=a`, 则说明有环 
       - 初始化标记所有点的状态为0.
       - 随便扔一个点 a 进 queue, 把它所有child都加入队列。如果child c被visit过，并且不是a->c, 那么证明环
       - 然后扫描下一层
    ```python
    # Python3 program to detect cycle in  
    # an undirected graph using BFS.
    from collections import deque

    def addEdge(adj: list, u, v):
        adj[u].append(v)
        adj[v].append(u)

    def isCyclicConnected(adj: list, s, V,  
                          visited: list):

        # Set parent vertex for every vertex as -1.
        parent = [-1] * V

        # Create a queue for BFS
        q = deque()

        # Mark the current node as  
        # visited and enqueue it
        visited[s] = True
        q.append(s)

        while q != []:

            # Dequeue a vertex from queue and print it
            u = q.pop()

            # Get all adjacent vertices of the dequeued
            # vertex u. If a adjacent has not been visited,
            # then mark it visited and enqueue it. We also
            # mark parent so that parent is not considered
            # for cycle.
            for v in adj[u]:
                if not visited[v]:
                    visited[v] = True
                    q.append(v)
                    parent[v] = u
                # 如果访问到一个已经visit过的点确不是它的parent，说明有环
                elif parent[u] != v:
                    return True

        return False

    def isCyclicDisconnected(adj: list, V):

        # Mark all the vertices as not visited
        visited = [False] * V

        for i in range(V):
            if not visited[i] and \
                   isCyclicConnected(adj, i, V, visited):
                return True
        return False
    ```

  [261. Graph Valid Tree](https://leetcode.com/problems/graph-valid-tree/description/)

### 有向图Topological Sort
  - BFS
    - 假设L是存放结果的列表，先找到那些入度为零的节点，把这些节点放到L中，因为这些节点没有任何的父节点。然后把与这些节点相连的边从图中去掉，再寻找图中的入度为零的节点。对于新找到的这些入度为零的节点来说，他们的父节点已经都在L中了，所以也可以放入L。重复上述操作，直到找不到入度为零的节点。如果此时L中的元素个数和节点总数相同，说明排序完成；如果L中的元素个数和节点总数不同，说明原图中存在环，无法进行拓扑排序。
- Count # of connected components
  同时记录下到底做了几次 BFS/DFS 才扫遍全图，图上就有几个 connected components  
  这时直接用DFS就可以，会比较好写  
  [323. Number of Connected Components in an Undirected Graph](https://leetcode.com/problems/number-of-connected-components-in-an-undirected-graph/description/)
- BFS来做树的层次遍历  
  应用：来判断树是否为完全二叉树[958. Check Completeness of a Binary Tree](https://leetcode.com/contest/weekly-contest-115/problems/check-completeness-of-a-binary-tree/)


## Topological Sorting

[LC 207. Course Schedule](https://leetcode.com/problems/course-schedule/discuss/)

求Course Schedule，等同问题是**有向图检测环**，vertex是course， edge是prerequisite。我觉得一般会使用Topological Sorting拓扑排序来检测。一个有向图假如有环则不存在Topological Order。一个DAG的Topological Order可以有大于1种。 常用的Topological Sorting算法有两种

1. Kahn's Algorithms (wiki)： BFS based， **start from with vertices with 0 incoming edge**，insert them into list S，at the same time we remove all their outgoing edges，after that find new vertices with 0 incoming edges and go on. 详细过程见Reference里Brown大学的课件。

其实就是不断的寻找有向图中没有前驱(入度为0)的顶点，将之输出。然后从有向图中删除所有以此顶点为尾的弧。重复操作，直至图空，或者找不到没有前驱的顶点为止。

该算法还可以判断有向图是否存在环(存在环的有向图肯定没有拓扑序列)，通过一个**count**记录visit过的顶点个数，如果少于N则说明存在环使剩余的顶点的入度不为0。因为环内的点永远无法满足`indegree = 0`（degree数组记录每个点的入度数）


对于BFS, 注意用把输入图List of Edges的表达方式转变为Adjacency Lists的表达方式。

```java
// BFS detect loop
public class Solution {
    public boolean canFinish(int numCourses, int[][] prerequisites) {
        if (numCourses < 0 || prerequisites == null) return false;
        if (prerequisites.length == 0) return true;
        List<List<Integer>> adjacencyListsGraph = new ArrayList<>();
        for (int i = 0; i < numCourses; i++) adjacencyListsGraph.add(new ArrayList<>());

        int[] inDegrees = new int[numCourses];
        for (int[] prerequisite : prerequisites) {
            adjacencyListsGraph.get(prerequisite[1]).add(prerequisite[0]);
            inDegrees[prerequisite[0]]++;
        }

        Queue<Integer> q = new LinkedList<>();
        for (int i = 0; i < numCourses; i++) {
            if (inDegrees[i] == 0) q.offer(i);
        }

        List<Integer> res = new ArrayList<>();
        while (!q.isEmpty()) {
            int src = q.poll();
            res.add(src);
            for (int dest : adjacencyListsGraph.get(src)) {
                inDegrees[dest]--;
                if (inDegrees[dest] == 0) q.offer(dest);
            }
        }
        return res.size() == numCourses;
    }
}
```

2. Tarjan's Algorithms (wiki)： DFS based， loop through each node of the graph in an arbitrary order，initiating a depth-first search that terminates when it hits any node that has already been visited since the beginning of the topological sort or the node has no outgoing edges (i.e. a leaf node). 详细过程见Reference里 NYU的课件。

```java
// DFS detect loop
public class Solution {
    public boolean canFinish(int numCourses, int[][] prerequisites) {
        if (numCourses < 0 || prerequisites == null) return false;
        if (prerequisites.length == 0) return true;
        List<List<Integer>> adjListsGraph = new ArrayList<>();
        for (int i = 0; i < numCourses; i++) adjListsGraph.add(new ArrayList<>());
        for (int[] prerequisite : prerequisites) adjListsGraph.get(prerequisite[1]).add(prerequisite[0]);

        List<Integer> res = new ArrayList<>();
        int[] visited = new int[numCourses];

        for (int i = 0; i < numCourses; i++) {
            // if 没有visited过i,
            if (visited[i]==0 && !canFinish(i, adjListsGraph, visited)) return false;
        }
        return true;
    }

    private boolean canFinish(int courseNum, List<List<Integer>> adjListsGraph, int[] visited) {
        if (visited[courseNum]==2) return true;
        visited[courseNum] = 1;
        for (int dependent : adjListsGraph.get(courseNum)) {
            if (visited[dependent]==1 || (visited[dependent] == 0 && !canFinish(dependent, adjListsGraph, visited))) {
                return false;
            }
        }

        //如果一条路径不满足（有loop）,则可以结束了，所以不需要把状态再重置。
        visited[courseNum] = 2;
        return true;
    }
}
```


## Tarjan's strongly connected components algorithm - DFS

[链接](https://www.byvoid.com/zhs/blog/scc-tarjan)有向图强连通分量的Tarjan算法在图论中，一个有向图被成为是强连通的（strongly connected）当且仅当每一对不相同结点 u 和 v 间既存在从 u 到 v 的路径也存在从 v 到 u 的路径。即，从图内任意一点出发都可以到达其他所有点。有向图的极大强连通子图（这里指点数极大）被称为强连通分量（strongly connected component）。

篇幅太大，放在另外一篇
- [Tarjan 有向图强连通分量算法, 无向图割边和割点](https://github.com/chungjin/course-notes/blob/master/algorithm/TSCC.md)
