# Tarjan 有向图强连通分量算法, 无向图割边和割点

## Terminology
- [reference: Tarjan算法：求解图的割点与桥（割边）](https://www.cnblogs.com/nullzx/p/7968110.html)
- connected component:
	a subgraph in which any two vertices are connected to each other by paths
- [strong connected component](https://www.byvoid.com/zhs/blog/scc-tarjan)
	在有向图G中，如果两个顶点vi,vj间（vi>vj）有一条从vi到vj的有向路径，同时还有一条从vj到vi的有向路径，则称两个**顶点强连通**(strongly connected)。如果有向图G的每两个顶点都强连通，称G是一个**强连通图**。
- [图的割点、桥与双连通分支](https://www.byvoid.com/zhs/blog/biconnect)
	对于无向图，删除顶点 和其相连的边后所包含的连通分量增多，则称为割点 (cut point)。同理，删除边 和其相连的顶点后图包含的连通分量增多，则 是割边 (cut edge) 或桥 (bridge)。
- 桥两边的一定是割点
- 除“桥两边的一定是割点”这个条件以外，还需要遍历所有的点，用U顶点的dnf值和它的所有的孩子顶点的low值进行比较，如果存在至少一个孩子顶点V满足`low[v] >= dnf[u]`, 那么u为割点。
    eg: 3为割点，没有桥
    ```
    1---2
     \ /
      3
     / \
    4---5
    ```

## Tarjan算法 O(N+M)

### 应用：求有向图的strong connected component
- [有向图强连通分量的Tarjan算法](https://www.byvoid.com/zhs/blog/scc-tarjan)

对图做`DFS`，定义`DFN(u)`为节点`u`搜索的次序编号(时间戳)，`Low(u)`为u或u的子树能够追溯到的最早的栈中节点的次序号。由定义可以得出，
```
Low(u)=Min
{
    DFN(u),
    Low(v),(u,v)为树枝边，u为v的父节点
    DFN(v),(u,v)为指向栈中节点的后向边(非横叉边)
}
```
当`DFN(u)=Low(u)`时，以u为根的搜索子树上所有节点是一个强连通分量。


伪码如下:
```
tarjan(u)
{
    DFN[u]=Low[u]=++Index                      // 为节点u设定次序编号和Low初值
    Stack.push(u)                              // 将节点u压入栈中
    for each (u, v) in E                       // 枚举每一条边
        if (v is not visted)               // 如果节点v未被访问过
            tarjan(v)                  // 继续向下找
            Low[u] = min(Low[u], Low[v])
        else if (v in S)                   // 如果节点v还在栈内
            Low[u] = min(Low[u], DFN[v])
    if (DFN[u] == Low[u])                      // 如果节点u是强连通分量的根
        repeat
            v = S.pop                  // 将v退栈，为该强连通分量中一个顶点
            print v
        until (u== v)
}
```

每个顶点都被访问了一次，且只进出了一次堆栈，每条边也只被访问了一次，所以该算法的时间复杂度为O(N+M)。


- test case
![](https://www.byvoid.com/upload/wp/2009/04/image5.png)
expected: {6: {6}, 5: {5}, 1: {1, 2, 3, 4}}), 有3个强连通分量
核心在于, only update low[u] when v are not visited or in stack(current subgraph). u和v在同一搜索树上，且方向是从u到v.

```python
'''
Tarjan algorithm to find the strong connected component

result format: a dict with node as key and set as value.
set is the strong connected component with the node as the root of the subgraph
'''
class Solution:
    def criticalConnections(self, edges):


        def dfs(n):
            nonlocal timestamp
            visited.add(n)
            dfn[n] = timestamp
            low[n] = timestamp
            timestamp += 1
            stack.append(n)
            for v in graph[n]:
                if v not in visited:
                    dfs(v)
                    # after do tarjan starts with v, uptate with low[v]
                    low[n] = min(low[n], low[v])
                elif v in stack:
                    low[n] = min(low[n], dfn[v])
                # notice: cannot update with the node visited but not in stack,
                # it may already removed from the graph, like 5->6, 6 already removed
                # abandon this edge.
            # n is the cutting vetex
            if low[n] == dfn[n]:
                print(stack)
                # pop the stack until n, which forms a scc, with root n
                while len(stack)>0:
                    tmp = stack.pop()
                    res[n].add(tmp)
                    if tmp == n:
                        break
        res = collections.defaultdict(set)
        dfn = collections.defaultdict(int)
        low = collections.defaultdict(int)
        timestamp = 1
        visited = set()
        stack = []
        graph = collections.defaultdict(set)
        for x, y in edges:
            graph[x].add(y)
        dfs(1)
        print(low)
        print(res)

def main():
    a = Solution()
    graph = [(1,3), (1,2), (2,4), (3,4), (3,5), (4,1), (4,6), (5,6)]
    a.criticalConnections(graph)

main()
```


### 应用：undirected 图的割点、桥与双连通分支
- [图的割点、桥通过找到strong connected component](https://www.byvoid.com/zhs/blog/biconnect)
- connected graph: the subgraph where any two vertices are connected.
- cut vertex
	A cut, vertex cut, or separating set of a connected graph G is a set of vertices whose **removal renders G disconnected**. The connectivity or vertex connectivity κ(G) (where G is not a complete graph) is the size of a minimal vertex cut. A graph is called k-connected or k-vertex-connected if its vertex connectivity is k or greater.
- cut edge
	cutting an edge, will make the **G disconnected**, The edge-connectivity λ(G) is the size of a smallest edge cut.

1. 求cutting vertex and bridge
用求有向图的tarjan算法。对图深度优先搜索，定义DFS(u)为u在搜索树（以下简称为树）中被遍历到的次序号。定义Low(u)为u或u的子树中能通过非父子边追溯到的最早的节点，即DFS序号最小的节点。根据定义，则有：

Low(u)=Min { DFS(u) DFS(v) (u,v)为后向边(返祖边) 等价于 DFS(v)<DFS(u)且v不为u的父亲节点 Low(v) (u,v)为树枝边(父子边) }

一个顶点u是**割点**，当且仅当满足(1)或(2) (1) u为树根，且u有多于一个子树(相邻点数目大于1)。 (2) u不为树根，且满足存在(u,v)为树枝边(或称父子边，即u为v在搜索树中的父亲)，使得DFS(u)<=Low(v)。

一条无向边(u,v)是**桥**，当且仅当(u,v)为edge，且满足DFS(u)<Low(v)。


- Summary: 即应用tarjan算法，得到每个点的dfn和low了之后，分析每一条边，如果有dfn[parent]<=low[child], 那么parent就是cutting vertex.

2. Biconnected component
在图G的所有子图G'中，如果G'是双连通的，则称G'为双连通子图。如果一个双连通子图G'它不是任何一个双连通子图的真子集，则G'为极大双连通子图。双连通分支(biconnected component)，或重连通分支，就是图的极大双连通子图。特殊的，点双连通分支又叫做块。
做Tarjan算法时，遇到`u->...->v,` 满足DFS(u)<=Low(v),说明u是个cutting vetex, 同时把边从栈顶一个个取出，直到遇到了边(u,v)，取出的这些边与其关联的点，组成一个点双连通分支。cutting vetex用来连接两个subcomponent.

## Leetcode
- 1192, follow up: 割点, 割边, Biconnected component

### 无向图求brige
Tarjan在运用时，注意dfs的时候传入parent节点，这样避免原路返回，造成error
- test case:
```
4
[[0,1],[1,2],[2,0],[1,3]]

6
[[0,1],[1,2],[2,0],[1,3],[3,4],[4,5],[5,3]]
```

- Solution:
```python
class Solution:
    def criticalConnections(self, n: int, connections: List[List[int]]) -> List[List[int]]:
        # tarjan
        def dfs(n, parent):
            nonlocal timestamp
            visited.add(n)
            dfn[n] = timestamp
            low[n] = timestamp
            timestamp += 1
            stack.append(n)
            stack_set.add(n)
            for v in graph[n]:
                if v == parent:
                    continue
                if v not in visited:
                    dfs(v, n)
                    # after do tarjan starts with v, uptate with low[v]
                    low[n] = min(low[n], low[v])
                elif v in stack_set:
                    low[n] = min(low[n], dfn[v])
                # notice: cannot update with the node visited but not in stack,
                # it may already removed from the graph, like 5->6, 6 already removed
                # abandon this edge.
            # n is the cutting vetex
            if low[n] == dfn[n]:
                #print(stack)
                # pop the stack until n, which forms a scc, with root n
                while len(stack)>0:
                    tmp = stack.pop()
                    stack_set.remove(tmp)
                    if tmp == n:
                        break
        res = []
        dfn = collections.defaultdict(int)
        low = collections.defaultdict(int)
        timestamp = 1
        visited = set()
        stack = []
        stack_set = set()
        graph = collections.defaultdict(set)
        for x, y in connections:
            graph[x].add(y)
            graph[y].add(x)
        dfs(0, -1)
        #print(low)
        for v in connections:
            if low[v[0]]>dfn[v[1]] or low[v[1]]>dfn[v[0]]:
                res.append(v)
        return res

```

### 无向图求Articulation Point
即割边上的点，有大于或等于两个边相连。
