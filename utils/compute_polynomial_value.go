package utils

import "math/big"

// ComputePolynomialValue 使用秦九韶算法计算多项式的值。
//
// coefficient: 多项式的系数，其中 coefficient[i] 是 x^i 的系数。
//
//	例如：P(x) = a_3*x^3 + a_2*x^2 + a_1*x + a_0，则 coefficient = {a_0, a_1, a_2, a_3}。
//
// x: 要求值的点。
// 返回值: P(x) 的计算结果。
func ComputePolynomialValue(coefficient []*big.Int, x *big.Int) big.Int {
	// 确保系数列表不为空
	if len(coefficient) == 0 {
		// 返回一个零
		return *big.NewInt(0)
	}

	// 从最高次系数开始，即秦九韶算法中的 a_n
	// 由于我们传入的系数数组是 {a_0, a_1, ..., a_n}，
	// 所以最高次系数是 coefficient[len(coefficient)-1]。

	// 使用一个新 big.Int 来存储最终结果或中间计算结果
	// 初始值设置为最高次系数 a_n
	result := new(big.Int).Set(coefficient[len(coefficient)-1])

	// 从倒数第二个系数开始 (a_{n-1}) 迭代到 a_0
	// i 的范围是 [len(coefficient)-2, 0]
	for i := len(coefficient) - 2; i >= 0; i-- {
		// 1. result = result * x
		// 使用 Multiply 方法将当前 result 乘以 x
		result.Mul(result, x)

		// 2. result = result + coefficient[i]
		// 使用 Add 方法将 coefficient[i] 加到 result 上
		result.Add(result, coefficient[i])
	}

	// 返回计算结果的副本
	return *result
}

/*
请注意：
1. 传入的 coefficient 数组应**从低次幂到高次幂**排列：{a_0, a_1, a_2, ..., a_n}。
2. 由于 big.Int 的方法通常是原地操作（修改接收者），因此需要使用 new(big.Int).Set()
   来初始化 result，以避免意外修改传入的 coefficient 数组中的值（尽管在这个实现中
   不会发生，但这是一个好的实践）。
*/
