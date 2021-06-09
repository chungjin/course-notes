# 高效算法设计

<!--ts-->
* [高效算法设计](#高效算法设计)
   * [sort and search](#sort-and-search)
      * [merge sort](#merge-sort)
      * [Quick sort](#quick-sort)
      * [二分查找](#二分查找)
   * [External Sort](#external-sort)
      * [External Merge Sort](#external-merge-sort)
   * [External Selection](#external-selection)
   * [[todo]divide and conquer more](#tododivide-and-conquer-more)
   * [[todo]Leetcode相关](#todoleetcode相关)

<!-- Added by: Jin Zhang, at: Tue Jun  8 23:21:34 PDT 2021 -->

<!--te-->


## sort and search

### merge sort
- merge sort  
  第一种Divide and Conquer是merge sort
  - 划分问题：把序列分成元素个数尽量相等的两半
  - 递归求解：把两半元素分别排序
  - 合并问题：把两个有序表合并成一个

  由于每次需要一个新表来存放结果，所以附加空间为n. 时间复杂度为O(nlgn)

  ```cpp
  //merge sort in A, range [x, y), extra space is T
  void merge_sort(int* A, int x, int y, int* T){
    if(y-x>1){
      int m = x + (y-x)/2;//向0递减，比如0, 1->0. -1, 0->0.
      int p = x, q = m, i = x;
      merge_sort(A, x, m , T);//[x, m)
      merge_sort(A, m, y, T);//[m, y)
      while(p<m || q<y){
        if(q>=y || (p<m && A[p]<=A[q])) T[i++] = A[p++];
        else T[i++] = A[q++];
      }
      for(i=x; i<y; i++) A[i] = T[i];
    }
  }
  ```

- 逆序对  
  然后用它可以解决另一个divide and conquer的问题: **逆序对问题**  
  给一列数，a1, a2, ..., an, 求它的逆序对，即有多少个有序对(i, j),使得i<j但ai>aj, 暴力求解时间复杂度为O(n^2)
  - 划分问题：把序列分成元素个数尽量相等的两半
  - 递归求解：把两半元素分别求出逆序对
  - 合并问题：然后求i, j分别在左半边和右半边的逆序对
    对于右边的每个j, 统计左边比它大的元素个数f(j), 则所有f(j)之和则是答案。刚好，merge-sort在合并问题阶段就有这样的特点

  ```cpp
  //merge sort in A, range [x, y), extra space is T
  void merge_sort(int* A, int x, int y, int* T){
    if(y-x>1){
      int m = x + (y-x)/2;//向0递减，比如0, 1->0. -1, 0->0.
      int p = x, q = m, i = x;
      int a = merge_sort(A, x, m , T);//[x, m)
      int b = merge_sort(A, m, y, T);//[m, y)
      int cnt = 0;
      while(p<m || q<y){
        if(q>=y || (p<m && A[p]<=A[q])) T[i++] = A[p++];
        else {
          T[i++] = A[q++];
          cnt += m-p;
        }
      }
      for(i=x; i<y; i++) A[i] = T[i];
      return a+b+cnt;
    }
    return 0;
  }
  ```

### Quick sort
- Quick sort  
  - 划分问题：把数组的各个元素重排后分成左右两个部分，使得左边的任意元素都小于或等于右边的任意元素
  - 递归求解：对左右两边分别排序
  - 合并问题：不用合并。此时已经完全有序

  ```java
    public class QuickSort {
        public static void main(String[] args) {
            int[] x = { 9, 2, 4, 7, 3, 7, 10 };
            System.out.println(Arrays.toString(x));

            int low = 0;
            int high = x.length - 1;

            quickSort(x, low, high);
            System.out.println(Arrays.toString(x));
        }

        public static void quickSort(int[] arr, int low, int high) {
            if (arr == null || arr.length == 0)
                return;

            if (low >= high)
                return;

            // pick the pivot
            int middle = low + (high - low) / 2;
            int pivot = arr[middle];

            // make left < pivot and right > pivot
            int i = low, j = high;
            while (i <= j) {
                while (arr[i] < pivot) {
                    i++;
                }

                while (arr[j] > pivot) {
                    j--;
                }

                if (i <= j) {
                    int temp = arr[i];
                    arr[i] = arr[j];
                    arr[j] = temp;
                    i++;
                    j--;
                }
            }

            // recursively sort two sub parts
            if (low < j)
                quickSort(arr, low, j);

            if (high > i)
                quickSort(arr, i, high);
        }
    }
  ```

- Quick select
  输入n个整数个一个正整数k, 输出这些整数从小到大排序后的第k个.
  - 划分问题：把数组的各个元素重排后分成左右两个部分，使得左边的任意元素都小于或等于右边的任意元素
  - 左半部分size为a, 右半部分为n-a. 如果a<k, 那么在右半部分继续寻找，否则在左半部分继续寻找
  - 时间复杂度分析:
      - T(N) = T(a1N) + T(a2N) + ... + T(anN) + O(n), if a1+a2+...+an==1, 时间复杂度为O(NlgN)
      - if a1+a2+...+an<1, 时间复杂度为O(N)

### 二分查找
- 二分查找
  ```cpp
  //search v in range [x, y]
  int bsearch(int* A, int x, int y, int v){
      int m;
      while(x<y){
          m = x + (y-x)/2;
          if(A[m]==v) return m;
          else if(A[m]>v) y = m;
          else x = m+1;
      }
      return -1;
  }
  ```

- 二分查找求下界  
  当v存在时返回它出现的**第一个位置**，如果不存在，返回这样一个下标i: 在此处插入v(原来的元素A[i], A[i+1],..., 全部往后移动一个位置)后序列仍然有序。
  ```cpp
    //求第一个出现的位置
    int lower_bound(int* A, int x, int y, int v){
        while(x<y){
          int m = x + (y-x)/2;
          if(A[m]>=v){
            y=m;
          }else{
            x = m+1;
          }
        }
        return x
    }
  ```
  分析:
  - A[m] = v, 至少已经找到一个，而左边可能还有，因此区间变为[x, m]
  - A[m] > v, 所求位置不可能在右边，但是有可能是m, 因此区间变为[x, m]
  - A[m] < v, 则m和m左边的位置都不可以，区间变为[m+1, y]
  **注意**潜在的危险：最后返回的区间如果和[x,y]相同，会陷入死循环。这种写法不会有。  
  **不会有这种情况的要诀是，永远不要设置x=m**

- 二分查找求上界  
  当v存在时返回它出现的**最后一个位置的后面一个位置**，如果不存在，返回这样一个下标i, 在此插入v, 数组仍然有序
  ```cpp
    //求最后一个出现的位置
    int upper_bound(int* A, int x, int y, int v){
        while(x<y){
          int m = x + (y-x)/2;
          if(A[m] > v){
            y = m;
          }else(A[m]<=v){
            x = m+1;
          }
        }
        return x
    }
  ```

## External Sort

### External Merge Sort
- [reference](https://zh.wikipedia.org/wiki/%E5%A4%96%E6%8E%92%E5%BA%8F)
做法是，map/reduce. 把数据分成N份，每分进行sort，存入临时文件中。
然后对文件头部创建一个指针，然后做类似于K位merge sort

## External Selection
类似于External Merge Sort, 不同的是，对每份sorted文件，保留一个头部指针。同时有个counter, 如果ith文件头部元素最小，counter++, file[i]++, 直到`counter == median`

## [todo]divide and conquer more

## [todo]Leetcode相关
