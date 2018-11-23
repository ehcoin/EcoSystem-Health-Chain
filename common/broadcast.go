// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package common

var (
	broadcastInterval  = uint64(20)
	reelectionInterval = uint64(60)
)

func IsBroadcastNumber(number uint64) bool {
	if number%broadcastInterval == 0 {
		return true
	}
	return false
}

func IsReElectionNumber(number uint64) bool {
	if number%reelectionInterval == 0 {
		return true
	}
	return false
}

func GetLastBroadcastNumber(number uint64) uint64 {
	if IsBroadcastNumber(number) {
		return number
	}
	ans := (number / broadcastInterval) * broadcastInterval
	return ans
}

func GetLastReElectionNumber(number uint64) uint64 {
	if IsReElectionNumber(number) {
		return number
	}
	ans := (number / reelectionInterval) * reelectionInterval
	return ans
}

func GetNextBroadcastNumber(number uint64) uint64 {
	if IsBroadcastNumber(number) {
		return number
	}
	ans := (number/broadcastInterval + 1) * broadcastInterval
	return ans
}

func GetNextReElectionNumber(number uint64) uint64 {
	if IsReElectionNumber(number) {
		return number
	}
	ans := (number/reelectionInterval + 1) * reelectionInterval
	return ans
}

func GetBroadcastInterval() uint64 {
	return broadcastInterval
}
func GetReElectionInterval() uint64 {
	return reelectionInterval
}
