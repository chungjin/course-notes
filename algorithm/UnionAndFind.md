# Union and Find

## Data structure
- 假设有N个点，M个Union操作
- Integer array id[] of size N.
- Interpretation: p and q are connected if they have the same id.
```
    i 0 1 2 3 4 5 6 7 8 9
id[i] 0 1 9 9 9 6 6 7 8 9
```
5 and 6 are connected. 2, 3, 4, and 9 are connected

## Basic Union and Find method
- Find. Check if p and q have the same id
- Union. To merge components containing p and q, change all entries with root id[p] to id[q].

## Improvement1: weighting
尽量让所形成的树扁平，于是Balance by linking small tree below large one.
- Implement:
    Maintain extra array sz[] to count number of elements in the tree rooted at i.
- Union:
    - merge smaller tree into larger tree
    - updat4e the sz[] array
    ```cpp
    if(sz[i]<sz[j]{id[i] = j, sz[j] += sz[i];}
    else {id[j] = i, sz[i] += sz[j];}
    ```
- time complexity  
    Union O(lgN), Find O(lgN)

## Improvement 2: Path compression
Just after computing the root of i, set the id of **each examined node** to root(i).
- Implement(in join method):
    - 标准做法：得到root后，重新加一个循环，将所有**examined node**都设置为root.
    - 简单one-pass做法：将所有**examined node**的root设置为它的grandparent.
    ```cpp
    public int root(int i)
    {
        while (i != id[i])
        {
            id[i] = id[id[i]];
            i = id[i];
        }
        return i;
    }
    ```
- Weighted + path worst-case time complexity
    (M+N)lgN

## Leetcode tag
- [leetcode union and find tag](https://leetcode.com/tag/union-find/)

## Reference
-  [Princeton AlgsDS07](https://www.cs.princeton.edu/~rs/AlgsDS07/01UnionFind.pdf)
