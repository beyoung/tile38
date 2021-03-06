package tests

import (
	"testing"
	"time"
)

func subTestKeys(t *testing.T, mc *mockServer) {
	runStep(t, mc, "BOUNDS", keys_BOUNDS_test)
	runStep(t, mc, "DEL", keys_DEL_test)
	runStep(t, mc, "DROP", keys_DROP_test)
	runStep(t, mc, "EXPIRE", keys_EXPIRE_test)
	runStep(t, mc, "FSET", keys_FSET_test)
	runStep(t, mc, "GET", keys_GET_test)
	runStep(t, mc, "KEYS", keys_KEYS_test)
	runStep(t, mc, "PERSIST", keys_PERSIST_test)
	runStep(t, mc, "SET", keys_SET_test)
	runStep(t, mc, "STATS", keys_STATS_test)
	runStep(t, mc, "TTL", keys_TTL_test)
}

func keys_BOUNDS_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"SET", "mykey", "myid1", "POINT", 33, -115}, {"OK"},
		{"BOUNDS", "mykey"}, {"[[-115 33 0] [-115 33 0]]"},
		{"SET", "mykey", "myid2", "POINT", 34, -112}, {"OK"},
		{"BOUNDS", "mykey"}, {"[[-115 33 0] [-112 34 0]]"},
		{"DEL", "mykey", "myid2"}, {1},
		{"BOUNDS", "mykey"}, {"[[-115 33 0] [-115 33 0]]"},
		{"SET", "mykey", "myid3", "OBJECT", `{"type":"Point","coordinates":[-130,38,10]}`}, {"OK"},
		{"SET", "mykey", "myid4", "OBJECT", `{"type":"Point","coordinates":[-110,25,-8]}`}, {"OK"},
		{"BOUNDS", "mykey"}, {"[[-130 25 -8] [-110 38 10]]"},
	})
}
func keys_DEL_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"SET", "mykey", "myid", "POINT", 33, -115}, {"OK"},
		{"GET", "mykey", "myid", "POINT"}, {"[33 -115]"},
		{"DEL", "mykey", "myid"}, {"1"},
		{"GET", "mykey", "myid"}, {nil},
	})
}
func keys_DROP_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"SET", "mykey", "myid1", "HASH", "9my5xp7"}, {"OK"},
		{"SET", "mykey", "myid2", "HASH", "9my5xp8"}, {"OK"},
		{"SCAN", "mykey", "COUNT"}, {2},
		{"DROP", "mykey"}, {1},
		{"SCAN", "mykey", "COUNT"}, {0},
		{"DROP", "mykey"}, {0},
		{"SCAN", "mykey", "COUNT"}, {0},
	})
}
func keys_EXPIRE_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"SET", "mykey", "myid", "STRING", "value"}, {"OK"},
		{"EXPIRE", "mykey", "myid", 1}, {1},
		{time.Second / 4}, {}, // sleep
		{"GET", "mykey", "myid"}, {"value"},
		{time.Second}, {}, // sleep
		{"GET", "mykey", "myid"}, {nil},
	})
}
func keys_FSET_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"SET", "mykey", "myid", "HASH", "9my5xp7"}, {"OK"},
		{"GET", "mykey", "myid", "WITHFIELDS", "HASH", 7}, {"[9my5xp7]"},
		{"FSET", "mykey", "myid", "f1", 105.6}, {1},
		{"GET", "mykey", "myid", "WITHFIELDS", "HASH", 7}, {"[9my5xp7 [f1 105.6]]"},
		{"FSET", "mykey", "myid", "f1", 0}, {1},
		{"GET", "mykey", "myid", "WITHFIELDS", "HASH", 7}, {"[9my5xp7]"},
		{"FSET", "mykey", "myid", "f1", 0}, {0},
		{"DEL", "mykey", "myid"}, {"1"},
		{"GET", "mykey", "myid"}, {nil},
	})
}
func keys_GET_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"SET", "mykey", "myid", "STRING", "value"}, {"OK"},
		{"GET", "mykey", "myid"}, {"value"},
		{"SET", "mykey", "myid", "STRING", "value2"}, {"OK"},
		{"GET", "mykey", "myid"}, {"value2"},
		{"DEL", "mykey", "myid"}, {"1"},
		{"GET", "mykey", "myid"}, {nil},
	})
}
func keys_KEYS_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"SET", "mykey11", "myid4", "STRING", "value"}, {"OK"},
		{"SET", "mykey22", "myid2", "HASH", "9my5xp7"}, {"OK"},
		{"SET", "mykey22", "myid1", "OBJECT", `{"type":"Point","coordinates":[-130,38,10]}`}, {"OK"},
		{"SET", "mykey11", "myid3", "OBJECT", `{"type":"Point","coordinates":[-110,25,-8]}`}, {"OK"},
		{"SET", "mykey42", "myid2", "HASH", "9my5xp7"}, {"OK"},
		{"SET", "mykey31", "myid4", "STRING", "value"}, {"OK"},
		{"KEYS", "*"}, {"[mykey11 mykey22 mykey31 mykey42]"},
		{"KEYS", "*key*"}, {"[mykey11 mykey22 mykey31 mykey42]"},
		{"KEYS", "mykey*"}, {"[mykey11 mykey22 mykey31 mykey42]"},
		{"KEYS", "mykey4*"}, {"[mykey42]"},
		{"KEYS", "mykey*1"}, {"[mykey11 mykey31]"},
		{"KEYS", "mykey*2"}, {"[mykey22 mykey42]"},
		{"KEYS", "*2"}, {"[mykey22 mykey42]"},
		{"KEYS", "*1*"}, {"[mykey11 mykey31]"},
	})
}
func keys_PERSIST_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"SET", "mykey", "myid", "STRING", "value"}, {"OK"},
		{"EXPIRE", "mykey", "myid", 2}, {1},
		{"PERSIST", "mykey", "myid"}, {1},
		{"PERSIST", "mykey", "myid"}, {0},
	})
}
func keys_SET_test(mc *mockServer) error {
	return mc.DoBatch(
		"point", [][]interface{}{
			{"SET", "mykey", "myid", "POINT", 33, -115}, {"OK"},
			{"GET", "mykey", "myid", "POINT"}, {"[33 -115]"},
			{"GET", "mykey", "myid", "BOUNDS"}, {"[[33 -115] [33 -115]]"},
			{"GET", "mykey", "myid", "OBJECT"}, {`{"type":"Point","coordinates":[-115,33]}`},
			{"GET", "mykey", "myid", "HASH", 7}, {"9my5xp7"},
			{"DEL", "mykey", "myid"}, {"1"},
			{"GET", "mykey", "myid"}, {nil},
		},
		"object", [][]interface{}{
			{"SET", "mykey", "myid", "OBJECT", `{"type":"Point","coordinates":[-115,33]}`}, {"OK"},
			{"GET", "mykey", "myid", "POINT"}, {"[33 -115]"},
			{"GET", "mykey", "myid", "BOUNDS"}, {"[[33 -115] [33 -115]]"},
			{"GET", "mykey", "myid", "OBJECT"}, {`{"type":"Point","coordinates":[-115,33]}`},
			{"GET", "mykey", "myid", "HASH", 7}, {"9my5xp7"},
			{"DEL", "mykey", "myid"}, {"1"},
			{"GET", "mykey", "myid"}, {nil},
		},
		"bounds", [][]interface{}{
			{"SET", "mykey", "myid", "BOUNDS", 33, -115, 33, -115}, {"OK"},
			{"GET", "mykey", "myid", "POINT"}, {"[33 -115]"},
			{"GET", "mykey", "myid", "BOUNDS"}, {"[[33 -115] [33 -115]]"},
			{"GET", "mykey", "myid", "OBJECT"}, {`{"type":"Polygon","coordinates":[[[-115,33],[-115,33],[-115,33],[-115,33],[-115,33]]]}`},
			{"GET", "mykey", "myid", "HASH", 7}, {"9my5xp7"},
			{"DEL", "mykey", "myid"}, {"1"},
			{"GET", "mykey", "myid"}, {nil},
		},
		"hash", [][]interface{}{
			{"SET", "mykey", "myid", "HASH", "9my5xp7"}, {"OK"},
			{"GET", "mykey", "myid", "HASH", 7}, {"9my5xp7"},
			{"DEL", "mykey", "myid"}, {"1"},
			{"GET", "mykey", "myid"}, {nil},
		},
		"field", [][]interface{}{
			{"SET", "mykey", "myid", "FIELD", "f1", 33, "FIELD", "a2", 44.5, "HASH", "9my5xp7"}, {"OK"},
			{"GET", "mykey", "myid", "WITHFIELDS", "HASH", 7}, {"[9my5xp7 [a2 44.5 f1 33]]"},
			{"FSET", "mykey", "myid", "f1", 0}, {1},
			{"FSET", "mykey", "myid", "f1", 0}, {0},
			{"GET", "mykey", "myid", "WITHFIELDS", "HASH", 7}, {"[9my5xp7 [a2 44.5]]"},
			{"DEL", "mykey", "myid"}, {"1"},
			{"GET", "mykey", "myid"}, {nil},
		},
		"string", [][]interface{}{
			{"SET", "mykey", "myid", "STRING", "value"}, {"OK"},
			{"GET", "mykey", "myid"}, {"value"},
			{"SET", "mykey", "myid", "STRING", "value2"}, {"OK"},
			{"GET", "mykey", "myid"}, {"value2"},
			{"DEL", "mykey", "myid"}, {"1"},
			{"GET", "mykey", "myid"}, {nil},
		},
	)
}

func keys_STATS_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"STATS", "mykey"}, {"[nil]"},
		{"SET", "mykey", "myid", "STRING", "value"}, {"OK"},
		{"STATS", "mykey"}, {"[[in_memory_size 9 num_objects 1 num_points 0 num_strings 1]]"},
		{"SET", "mykey", "myid2", "STRING", "value"}, {"OK"},
		{"STATS", "mykey"}, {"[[in_memory_size 19 num_objects 2 num_points 0 num_strings 2]]"},
		{"SET", "mykey", "myid3", "OBJECT", `{"type":"Point","coordinates":[-115,33]}`}, {"OK"},
		{"STATS", "mykey"}, {"[[in_memory_size 40 num_objects 3 num_points 1 num_strings 2]]"},
		{"DEL", "mykey", "myid"}, {1},
		{"STATS", "mykey"}, {"[[in_memory_size 31 num_objects 2 num_points 1 num_strings 1]]"},
		{"DEL", "mykey", "myid3"}, {1},
		{"STATS", "mykey"}, {"[[in_memory_size 10 num_objects 1 num_points 0 num_strings 1]]"},
		{"STATS", "mykey", "mykey2"}, {"[[in_memory_size 10 num_objects 1 num_points 0 num_strings 1] nil]"},
		{"DEL", "mykey", "myid2"}, {1},
		{"STATS", "mykey"}, {"[nil]"},
		{"STATS", "mykey", "mykey2"}, {"[nil nil]"},
	})
}
func keys_TTL_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"SET", "mykey", "myid", "STRING", "value"}, {"OK"},
		{"EXPIRE", "mykey", "myid", 2}, {1},
		{time.Second / 4}, {}, // sleep
		{"TTL", "mykey", "myid"}, {1},
	})
}
