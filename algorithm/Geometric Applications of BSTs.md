# Geometric Applications of BSTs

- [slides](https://www.cs.princeton.edu/courses/archive/spr15/cos226/lectures/99GeometricSearch.pdf)

## 1d range search
- [youtube video](https://www.youtube.com/watch?v=aRDnG4pMhrU)

- Application: Database queries
- append一个ordered的list
- 在任意时刻都可能做range query
	+ 找到all keys between k1 and k2
	+ count all keys between k1 and k2
- eg: append `[A, B, D, F, H, I, P]`
- search G to K -> `H, I`
- count from G to K -> 2
```
               S(4)
          /        \
        E(2)        X(5)
     /     \
    A(0)   R(3)
     \     
      C(1) 

//时间复杂度O(lgn)
public int size(Key lo, Key hi)
{
// rank is to find the largest element smaller than key	
 if (contains(hi)) return rank(hi) - rank(lo) + 1;
 else return rank(hi) - rank(lo);
} 



// how to find the largest element smaller than key
if root.left == None and root.right == None:
	return root
elif root.val < val:
	root = root.right
elif root.val>root:
	root = root.left	


// 输出所有[k1, k2] 之间的key, 例如[F, T]
// 类似于inorder
- Recursively find all keys in left subtree (if any could fall in range).
- Check key in current node.
- Recursively find all keys in right subtree (if any could fall in range).
- 时间复杂度 R + lgN, R为k1到k2的点的数目
```

### Exercise 1: Prefix Max Query
- [reference](https://www.cs.cmu.edu/afs/cs/academic/class/15210-f11/www/resources/recis/rec12.pdf)
- 输入是一串 unordered list(x, y), 要求输入是一个query(x), 输出是所有 pair(x', y'), 其中x'<\x,, 的所有y'中的最大值。
- 构建一个二叉树，node的value是x。
- 二叉树节点node中同时存有一个max, 保存自己左子树的最大y.

### Exercise 2: Sweep-line algorithm
- [reference](https://cse.taylor.edu/~jdenning/classes/cos265/slides/10_BSTGeometry.html)
![](https://cse.taylor.edu/~jdenning/classes/cos265/slides/10_BSTGeometry/images/lineintersection1.png)
- 将interval按照x从小到大排列
- sweep vertical line from left to right
- h-segment (left endpoint): insert y-coordinate into BST
- h-segment (right endpoint): remove y-coordinate from BST
- v-segment: 碰到竖线的时候range search for interval of y-endpoints



## kd trees
- [youtube video](https://www.youtube.com/watch?v=1OoM0phlO_U)

### 2d orthogonal range search
- Insert a 2D key
- Search for a 2D key
- Delete a 2D key
- Range search: find all keys that lie in a 2D range
- Range count: number of keys that lie in a 2D range


#### Problem 1: check the points in the rectangle
![](img/Geometric_0.png)

- Check if point in node lies in given rectangle.
- 如果线穿过rectangle, 则需要check both size
- 如果线不穿过rectangle, 只需要check一边。


#### Problem 2: nearest neighnor search in a 2d tree
![](img/Geometric_1.png)
inorder traverse, 主要是剪枝操作，例如做完root的左子树的traverse之后，发现1到query point的vertical距离已经超过了目前的最短距离，则不用再traverse右子树。


## interval search trees
- [youtube video](https://www.youtube.com/watch?v=q0QOYtSsTg4)


### Interval search API
```java
public class IntervalST<Key extends Comparable<Key>, Value>
IntervalST() //create interval search tree
void put(Key lo, Key hi, Value val) //put interval-value pair into ST
Value get(Key lo, Key hi) //value paired with given interval
void delete(Key lo, Key hi) //delete the given interval
Iterable<Value> intersects(Key lo, Key hi) //all intervals that intersect (lo, hi)
```
- Nondegeneracy assumption. No two intervals have the same left endpoint.

### Interval search trees
![](img/Geometric_2.png)
- Build Tree or Insert an interval
	Create BST, where each node stores an interval `( lo, hi )`.
	- Use left endpoint as BST key.
	- Store/Update max endpoint in subtree rooted at node.
- To search for any interval that intersects query interval `(lo, hi)`
	- If interval in node intersects query interval, return it.
	- Else if left subtree is null, go right.
	- Else if max endpoint in left subtree is less than `lo`, go right.
	- Else go left.


## rectangle intersection
- [youtube video](https://www.youtube.com/watch?v=p9cChQlgx08)

### Orthogonal rectangle intersection: sweep-line algorithm

Sweep vertical line from left to right.
- x-coordinates of left and right endpoints define events.
- Maintain set of rectangles that intersect the sweep line in an interval
search tree (using y-intervals of rectangle).
- Left endpoint: **interval search** for y-interval of rectangle; insert y-interval.
- Right endpoint: remove y-interval.
