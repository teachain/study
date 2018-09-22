##算法（golang）##

####冒泡排序####

```
package algorithm

/*
 * 冒泡排序
 * 待排序的整型数组
 * 是否递增排序
 *
 */
func BubbleSort(arr []int, increase bool) {
	len := len(arr) - 1 
	for i := 0; i < len; i++ {
		for j := 0; j < len-i; j++ {
			if increase {
				if arr[j] > arr[j+1] {
					tmp := arr[j]
					arr[j] = arr[j+1]
					arr[j+1] = tmp
				}
			} else {
				if arr[j] < arr[j+1] {
					tmp := arr[j]
					arr[j] = arr[j+1]
					arr[j+1] = tmp
				}
			}
		}
	}
}

```
冒泡排序的基本思想是：每次比较两个相邻的元素，如果他们的顺序错误就把他们交换过来。


####快速排序####

```
/*
 * 快速排序
 * 待排序的整型数组
 * 是否递增排序
 *
 */
func QuickSort(arr []int, left int, right int, increase bool) {

	if left > right {
		return
	}
	var i, j, temp int
   //我们始终以左边第一个数为基准数
	temp = arr[left] //temp中存的就是基准数

	i = left

	j = right

	for i != j {

		if increase {
			//顺序很重要，要先从右边往左找
			for arr[j] >= temp && i < j {
				j--
			}

			for arr[i] <= temp && i < j {
				i++
			}
		} else {
			//顺序很重要，要先从右边往左找
			for arr[j] <= temp && i < j {
				j--
			}

			for arr[i] >= temp && i < j {
				i++
			}
		}

		//当哨兵i和哨兵j没有相遇时，交换两个数在数组中的位置
		if i < j {

			change := arr[i]

			arr[i] = arr[j]

			arr[j] = change
		}

	}
	//最终将基准数归位（也就是基准数最终找到他的最终位置）
	arr[left] = arr[i]

	arr[i] = temp
   
	QuickSort(arr, left, i-1, increase)

	QuickSort(arr, i+1, right, increase)
}


```


##直接插入排序##

####（类似玩牌时整理手中纸牌的过程）####

方法：每次从无序表中取出第一个元素，把它插入到有序表的合适位置，使有序表仍然有序。

第一趟比较前两个数，然后把第二个数按大小插入到有序表中，第二趟把第三个数据与前两个数从后向前扫描，把第三个数按大小插入到有序表中，依次进行下去，进行了(n-1)趟扫描以后就完成了整个排序过程。

直接插入排序属于稳定的排序，时间复杂度为o(n^2),空间复杂度为o(1)。
























