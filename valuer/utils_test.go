package easyconfmgrvaluer_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/derekparker/trie"
	"github.com/soyacen/goutils/sliceutils"

	easyconfmgrvaluer "github.com/soyacen/easyconfmgr/valuer"
)

func TestMapToTrie(t *testing.T) {
	key1 := "val1"
	key2 := 1
	key3 := true
	key4 := time.Now()
	key5 := time.Second
	key6 := 1.414
	subconf2key2 := []string{"a", "b", "c"}
	subconf2key3 := []int{1, 2, 3, 4, 5}
	subconf2subconf3key2 := "val2"

	subconfkey1 := "val1"
	subconf := map[string]interface{}{
		"key1": subconfkey1,
	}

	subconf2subconf3key1 := "val1"
	subconf2subconf3 := map[string]interface{}{
		"key1": subconf2subconf3key1,
		"key2": subconf2subconf3key2,
	}
	subconf2 := map[interface{}]interface{}{
		"key2":     subconf2key2,
		"key3":     subconf2key3,
		"subconf3": subconf2subconf3,
	}
	config := map[string]interface{}{
		"key1":     key1,
		"key2":     key2,
		"key3":     key3,
		"key4":     key4,
		"key5":     key5,
		"key6":     key6,
		"subconf":  subconf,
		"subconf2": subconf2,
	}

	exceptedKeys := []string{
		"",
		"key1",
		"key2",
		"key3",
		"key4",
		"key5",
		"key6",
		"subconf",
		"subconf.key1",
		"subconf2",
		"subconf2.key2",
		"subconf2.key3",
		"subconf2.subconf3",
		"subconf2.subconf3.key1",
		"subconf2.subconf3.key2",
	}

	exceptedVals := []interface{}{
		config,
		key1,
		key2,
		key3,
		key4,
		key5,
		key6,
		subconf,
		subconfkey1,
		subconf2,
		subconf2key2,
		subconf2key3,
		subconf2subconf3,
		subconf2subconf3key1,
		subconf2subconf3key2,
	}

	tree := trie.New()
	easyconfmgrvaluer.MapToTrie(config, tree)
	actualKeys := tree.Keys()
	for i := range exceptedKeys {
		key := exceptedKeys[i]
		if sliceutils.NotContainsString(actualKeys, key) {
			t.Fatal("not Contains", key)
		}
		val := exceptedVals[i]
		node, ok := tree.Find(key)
		if !ok {
			t.Fatal("not found", key)
		}
		if !reflect.DeepEqual(val, node.Meta()) {
			t.Fatalf("expected val is \n%v,\n but actual is \n%v", val, node.Meta())
		}
	}

}
