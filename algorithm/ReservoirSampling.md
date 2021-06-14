# Reservoir Sampling

- [wiki: 水塘抽样 Reservoir sampling](https://en.wikipedia.org/wiki/Reservoir_sampling)



## 已知大小

实现一个功能，从长度为n的数组中随机返回k个数（k < n），要求每个数被选中返回的概率一样。

```python
def sample(nums, k):
	res = []
	while k>0:
		n = len(nums)
		i = random(n)
		res.append(nums[i])
		swap(i, n-1)
		nums.pop()
```

## 数据流
当内存无法加载全部数据时，如何从包含未知大小的数据流中随机选取k个数据，并且要保证每个数据被抽取到的概率相等。

数据流中第i个数被保留的概率为 `1/i`。只要采取这种策略，只需要遍历一遍数据流就可以得到采样值，并且保证所有数被选取的概率均为 `1/N` 。

第`i`个数, 做 `j = random(i)`, if `j<k`, 那么把`i`放入`res[j]`中。
```python
def sample(stream, k):
	res = stream[0:k]
	count = k
	while stream[k+1:].hasnext:
		count+=1
		ran = random(count)
		if count<k:
			res[count] = stream[k+1:].next
```

- 推导:
![image](https://user-images.githubusercontent.com/11788053/121855230-f2338700-cca7-11eb-8c41-6c0e7a1460d7.png)
```
possibility = k/i * (1 - 1/(i+1)) * (1 - 1/(i+2)) * (1 - 1/(i+3)) * ... * (1 - 1/(n)) = k/n  

k/i 被放进的概率，(1-1/(i+1)) 不被选出去的概率
```

- [Leetcode tag: reservoir-sampling](https://leetcode.com/tag/reservoir-sampling/)
