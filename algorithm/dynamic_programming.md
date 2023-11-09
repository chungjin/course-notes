<!-- toc -->

# 单源,权重为正:Dijkstra算法 O(mlogn)
思路是，用priorityqueue存source到所有未探索到的点的距离。每次从中poll出最小的，标记为visited(已探索到的点)，然后update priorityqueue中的各个未探索到点的距离。
```java
class DijShortestPath{
    class Node{
        int index;
        int cost;
        Node(int index, int cost){
            this.index = index;
            this.cost = cost;
        }
    }

    void dijkstra(int[][] graph, int src){
        int n = graph.length;
        Map<Integer, Node> map = new HashMap<>();
        PriorityQueue<Node> pq = new PriorityQueue<Node>(10, (a, b)->a.cost-b.cost);
        boolean[] visited = new boolean[n];
        visited[src] = true;
        Node cur = null;
        for(int i=0; i<n; i++){
            if(i!=src){
                if(graph[src][i]!=0){
                    cur = new Node(i, graph[src][i]);
                }else{
                    cur = new Node(i, Integer.MAX_VALUE);
                }
                map.put(i, cur);
                pq.add(cur);
            }
        }

        while(!pq.isEmpty()){
            cur = pq.poll();
            visited[cur.index] = true;
            for(int i=0; i<n; i++){
                if(graph[cur.index][i]!=0 && !visited[i]){
                    Node node = map.get(i);
                    node.cost = Math.min(node.cost, cur.cost+graph[cur.index][i]);
                }
            }
        }
        
        System.out.println("here");
    }

    // Driver method
    public static void main (String[] args)
    {
        /* Let us create the example graph discussed above */
        int graph[][] = new int[][]{{0, 4, 0, 0, 0, 0, 0, 8, 0},
            {4, 0, 8, 0, 0, 0, 0, 11, 0},
            {0, 8, 0, 7, 0, 4, 0, 0, 2},
            {0, 0, 7, 0, 9, 14, 0, 0, 0},
            {0, 0, 0, 9, 0, 10, 0, 0, 0},
            {0, 0, 4, 14, 10, 0, 2, 0, 0},
            {0, 0, 0, 0, 0, 2, 0, 1, 6},
            {8, 11, 0, 0, 0, 0, 1, 0, 7},
            {0, 0, 2, 0, 0, 0, 6, 7, 0}
        };
        DijShortestPath t = new DijShortestPath();
        t.dijkstra(graph, 0);
        
        /** output
         * Vertex   Distance from Source
            0                0
            1                4
            2                12
            3                19
            4                21
            5                11
            6                9
            7                8
            8                14
         */
    }
}
```

# 单源,权重可以为负:Bellman-Ford O(mn)
![](./img/Bellman-Ford.png)
inner-loop要计算到点v的所有in-degree, 时间复杂度为O(m)，总的时间复杂度为O(mn)(m为边的数目，n为点的数目)
[来源:cnblogs](http://www.cnblogs.com/gaochundong/p/bellman_ford_algorithm.html)

```java
public class BellmanFord {
    
    public void bellmanFord(int[][] graph, int source) {
        //construct indegree graph, transpose the graph
        int n = graph.length;
        int[][] D = new int[n][n];
        for(int i=0; i<n; i++) {
            if(i != source) D[i][0] = Integer.MAX_VALUE;
        }
        
        for(int k=1; k<n; k++) {
            for(int v=0; v<n; v++) {
                int min = D[v][k-1];
                for(int indegree = 0; indegree <n; indegree++) {
                    if(graph[indegree][v]!=0 && D[indegree][k-1]!=Integer.MAX_VALUE) {
                        min = Math.min(min, D[indegree][k-1]+graph[indegree][v]);
                    }
                }
                D[v][k] = min;
            }
        }
        List<String> output = new ArrayList<>();
        for(int i=0; i<n; i++) {
            output.add(String.valueOf(D[i][n-1]));
        }
        String joined = String.join(" ", output);
        System.out.println(joined);
    }
    
    
    public static void main (String[] args)
    {
        /* Let us create the example graph discussed above */
        int graph[][] = new int[][]{
            {0, -1, 4, 0, 0},
            {0, 0, 3, 2, 2},
            {0, 0, 0, 0, 0},
            {0, 1, 5, 0, 0},
            {0, 0, 0, -3, 0},
        };
        BellmanFord t = new BellmanFord();
        t.bellmanFord(graph, 0);
        // output: 0 -1 2 -2 1
    }
}
```

# All Pairs: Floyd-Warshall O(n^3)
[](./img/Floyd-Warshall.png)
```java
public class FloydWarshall {
    final static int INF = 99999, V = 4;
    
    public void floydWarshall(int[][] graph) {
        //int[][] D = new int[V][V];
        
        for(int k=0; k<V; k++) {//add k to set
            for(int i=0; i<V; i++) {
                for(int j=0; j<V; j++) {
                    graph[i][j] = Math.min(graph[i][j], (graph[i][k] + graph[k][j]));
                }
            }
        }
        //print
        for(int i=0; i<V; i++) {
            System.out.println(Arrays.toString(graph[i]));
        }
    }
    
    
    public static void main (String[] args)
    {
        /* Let us create the example graph discussed above */
        int graph[][] = { {0,   5,  INF, 10},
                {INF, 0,   3, INF},
                {INF, INF, 0,   1},
                {INF, INF, INF, 0}
              };
        FloydWarshall t = new FloydWarshall();
        t.floydWarshall(graph);
        //answer
        /*
         * [0, 5, 8, 9]
         * [99999, 0, 3, 4]
         * [99999, 99999, 0, 1]
         * [99999, 99999, 99999, 0]
         */
    }
}
```

## traveling salesperson problem
[详细图解](https://jerkwin.github.io/2016/03/17/%E6%97%85%E8%A1%8C%E6%8E%A8%E9%94%80%E5%95%86%E9%97%AE%E9%A2%98TSP%E7%9A%84%E5%8A%A8%E6%80%81%E8%A7%84%E5%88%92%E8%A7%A3%E6%B3%95/)

一个人需要visit图上的所有点,并返回到原点, 并且保证路径sum最小。

### 有向图
假定我们从城市0出发, 城市1, 2, 3每个经过一次, 最后回到城市0, 那么求解的递归树可以表示如下:

这里有几个Trick:
1. 我们可以用bit来表示点的集合，如点0,1,2在集合中，可以表示为111,即为7
2. 空间复杂度O(n*2^n)
3. 对于空间每个位置，我们都进行了n次操作，时间复杂度O(n^2 * 2^n)
![](./img/TravelingSalesperson.png)

d[i, j]可以表示为从点i出发，经过集合j，回到0点的最短路径长度。

```java
public class TravelingSalesperson {
    int n;
    public void travelingSalesperson(int[][] graph) {
        n = graph.length;
        int[][] dp = new int[n][2<<(n-1)];
        for(int i=0; i<n; i++) {
            dp[i][0] = graph[i][0];//return back to 0
        }
        System.out.println(get(graph, dp, 0, (1<<(n-1))-1));
    }
    
    private int get(int[][] graph, int[][] dp, int source, int j) {
        if(dp[source][j]!=0) return dp[source][j];
        int min = Integer.MAX_VALUE;
        for(int k=0; k<n-1; k++) {
            if((j&(1<<k))!=0) {//set contains node k+1
                min = Math.min(min, graph[source][k+1]+get(graph, dp, k+1, j-(1<<k)));
            }
        }
        dp[source][j] = min;
        return dp[source][j];
    }
    
    public static void main (String[] args)
    {
        /* Let us create the example graph discussed above */
        int graph[][] = { {0, 3, 6, 7},
                {5, 0, 2, 3},
                {6, 4, 0, 2},
                {3, 7, 5, 0}
              };
        TravelingSalesperson t = new TravelingSalesperson();
        t.travelingSalesperson(graph);
        //最后的长度为10
    }
}
```

