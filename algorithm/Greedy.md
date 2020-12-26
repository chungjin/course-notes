# Greedy

## 区间相关问题
1. 尽量选择更多的不相交区间。  
数轴上有n个开区间(ai, bi). 选择**尽量多个区间**，使得这些区间两两没有公共点。**贪心策略是，一定要选第一个区间。**  
假设有两个区间x,y. 如果x完全包含y, 那么则选择y, 放弃x.   
接下来按照bi从小到大的顺序给区间排序。现在区间已经排序成b1<=b2<=b3...了。考虑a1, a2的大小关系。
- 情况1: a1>a2, 已经完全被剔除
- 情况2: 一定有a1<a2<a3... 如果区间2和1完全不相交，需要选第一个。否则，因为a1<a2, 所以也要选**第一个**。不能选第二个。
  选择了区间1以后，需要把所有和区间1相交的区间都排除在外。

[LC435. Non-overlapping Intervals](https://leetcode.com/problems/non-overlapping-intervals/)

```python
# Definition for an interval.
# class Interval:
#     def __init__(self, s=0, e=0):
#         self.start = s
#         self.end = e

class Solution:
    def eraseOverlapIntervals(self, intervals):
        """
        :type intervals: List[Interval]
        :rtype: int
        """
        res = 0
        if len(intervals)<=1:
            return res
        intervals.sort(key=lambda x:x.end)
        start =  intervals[0].start   
        end = intervals[0].end
        for i in range(1, len(intervals)):
            if start>=intervals[i].start:
                res+=1
            elif end>intervals[i].start:
                res +=1
            else:
                end = intervals[i].end
        return res
```

2. 区间选点问题
数轴上有n个闭区间[ai, bi]. 取尽量少的点，使得每个区间内都至少有一个点(不同区间内含的点可以是同一个)。  
受上一题启发，由于小区间被满足时，大区间一定满足。于是可以剔除掉小区间。
- 把所有区间按b从小到大排序(b相同时a从大到小排序). 如果出现区间包含的情况，小区间一定排在前面。**第一个区间应该，取最后一个点**

同理，如果拓展成两个点，每次尽量取最后一个点/两个点。
- [LC757. Set Intersection Size At Least Two](https://leetcode.com/problems/set-intersection-size-at-least-two/)

3. 区间覆盖问题
数轴上有n个闭区间[ai,bi], 选择尽量少的区间覆盖一条指定线段[s, t].  
- 首先是预处理，每个区间在[s,t]外的部分都应该被预先被切掉，因为它们的存在毫无意义。
- 预处理后，在互相包含的情况下，小区间不被考虑
- 把各区间按照a从小到大排序。如果区间1的起点不是s, 无解(因为其他区间的起点更大，更不会覆盖到s).
- 选择此区间[ai, bi]后，新的起点设置为bi, 并且忽略所有区间在bi之前的部分。
