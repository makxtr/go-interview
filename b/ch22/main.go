package main

import (
	"fmt"
	"sort"
)

func main() {
	// Standard U.S. coin denominations in cents

	denominations := []int{1, 3, 4}
	amount := 6

	coinCombo := CoinCombination(amount, denominations)
	fmt.Println(coinCombo)

	fmt.Println("DP MinCoins:", MinCoinsDP(amount, denominations))
	fmt.Println("DP Combination:", CoinCombinationDP(amount, denominations))

	//denominations := []int{1, 5, 10, 25, 50}
	//
	//// Test amounts
	//amounts := []int{87, 42, 99, 33, 7}
	//
	//for _, amount := range amounts {
	//	// Find minimum number of coins
	//	minCoins := MinCoins(amount, denominations)
	//
	//	// Find coin combination
	//	coinCombo := CoinCombination(amount, denominations)
	//
	//	// Print results
	//	fmt.Printf("Amount: %d cents\n", amount)
	//	fmt.Printf("Minimum coins needed: %d\n", minCoins)
	//	fmt.Printf("Coin combination: %v\n", coinCombo)
	//	fmt.Println("---------------------------")
	//}
}

// MinCoinsDP returns the minimum number of coins needed to make the given amount
// using bottom-up dynamic programming. Returns -1 if the amount cannot be made.
// Time: O(amount * len(denominations)). Space: O(amount).
func MinCoinsDP(amount int, denominations []int) int {
	if amount < 0 {
		return -1
	}
	const inf = int(^uint(0) >> 1)
	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		dp[i] = inf
	}
	for i := 1; i <= amount; i++ {
		for _, coin := range denominations {
			if coin <= i && dp[i-coin] != inf && dp[i-coin]+1 < dp[i] {
				dp[i] = dp[i-coin] + 1
			}
		}
	}
	if dp[amount] == inf {
		return -1
	}
	return dp[amount]
}

// CoinCombinationDP returns the optimal coin combination using DP. Keys are coin
// denominations, values are coin counts. Returns an empty map if the amount cannot
// be made. Tracks the chosen coin at each subproblem to reconstruct the solution.
func CoinCombinationDP(amount int, denominations []int) map[int]int {
	if amount < 0 {
		return map[int]int{}
	}
	const inf = int(^uint(0) >> 1)
	dp := make([]int, amount+1)
	pick := make([]int, amount+1) // coin chosen to reach amount i
	for i := 1; i <= amount; i++ {
		dp[i] = inf
	}
	for i := 1; i <= amount; i++ {
		for _, coin := range denominations {
			if coin <= i && dp[i-coin] != inf && dp[i-coin]+1 < dp[i] {
				dp[i] = dp[i-coin] + 1
				pick[i] = coin
			}
		}
	}
	if dp[amount] == inf {
		return map[int]int{}
	}
	combination := make(map[int]int)
	for i := amount; i > 0; {
		coin := pick[i]
		combination[coin]++
		i -= coin
	}
	return combination
}

// MinCoins returns the minimum number of coins needed to make the given amount.
// If the amount cannot be made with the given denominations, return -1.
func MinCoins(amount int, denominations []int) int {
	// Sort denominations in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(denominations)))

	coinCount := 0
	remainingAmount := amount

	for _, coin := range denominations {
		// Take as many coins of this denomination as possible
		count := remainingAmount / coin
		coinCount += count
		remainingAmount -= count * coin

		// If we've reached the target amount, we're done
		if remainingAmount == 0 {
			return coinCount
		}
	}

	// If we couldn't make the exact amount
	if remainingAmount > 0 {
		return -1
	}

	return coinCount
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	// Sort denominations in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(denominations)))

	combination := make(map[int]int)
	remainingAmount := amount

	for _, coin := range denominations {
		// Take as many coins of this denomination as possible
		count := remainingAmount / coin
		if count > 0 {
			combination[coin] = count
		}
		remainingAmount -= count * coin

		// If we've reached the target amount, we're done
		if remainingAmount == 0 {
			return combination
		}
	}

	// If we couldn't make the exact amount, return an empty map
	if remainingAmount > 0 {
		return map[int]int{}
	}

	return combination
}
