package godis

import (
	"fmt"
	"github.com/hdt3213/godis/lib/utils"
	"github.com/hdt3213/godis/redis/reply"
	"github.com/hdt3213/godis/redis/reply/asserts"
	"strconv"
	"testing"
)

func TestGeoHash(t *testing.T) {
	FlushDB(testDB, utils.ToBytesList())
	key := utils.RandString(10)
	pos := utils.RandString(10)
	result := GeoAdd(testDB, utils.ToBytesList(key, "13.361389", "38.115556", pos))
	asserts.AssertIntReply(t, result, 1)
	result = GeoHash(testDB, utils.ToBytesList(key, pos))
	asserts.AssertMultiBulkReply(t, result, []string{"sqc8b49rnys00"})
}

func TestGeoRadius(t *testing.T) {
	FlushDB(testDB, utils.ToBytesList())
	key := utils.RandString(10)
	pos1 := utils.RandString(10)
	pos2 := utils.RandString(10)
	GeoAdd(testDB, utils.ToBytesList(key,
		"13.361389", "38.115556", pos1,
		"15.087269", "37.502669", pos2,
	))
	result := GeoRadius(testDB, utils.ToBytesList(key, "15", "37", "200", "km"))
	asserts.AssertMultiBulkReplySize(t, result, 2)
}

func TestGeoRadiusByMember(t *testing.T) {
	FlushDB(testDB, utils.ToBytesList())
	key := utils.RandString(10)
	pos1 := utils.RandString(10)
	pos2 := utils.RandString(10)
	pivot := utils.RandString(10)
	GeoAdd(testDB, utils.ToBytesList(key,
		"13.361389", "38.115556", pos1,
		"17.087269", "38.502669", pos2,
		"13.583333", "37.316667", pivot,
	))
	result := GeoRadiusByMember(testDB, utils.ToBytesList(key, pivot, "100", "km"))
	asserts.AssertMultiBulkReplySize(t, result, 2)
}

func TestGeoPos(t *testing.T) {
	FlushDB(testDB, utils.ToBytesList())
	key := utils.RandString(10)
	pos1 := utils.RandString(10)
	pos2 := utils.RandString(10)
	GeoAdd(testDB, utils.ToBytesList(key,
		"13.361389", "38.115556", pos1,
	))
	result := GeoPos(testDB, utils.ToBytesList(key, pos1, pos2))
	expected := "*2\r\n*2\r\n$18\r\n13.361386698670685\r\n$17\r\n38.11555536696687\r\n*0\r\n"
	if string(result.ToBytes()) != expected {
		t.Error("test failed")
	}
}

func TestGeoDist(t *testing.T) {
	FlushDB(testDB, utils.ToBytesList())
	key := utils.RandString(10)
	pos1 := utils.RandString(10)
	pos2 := utils.RandString(10)
	GeoAdd(testDB, utils.ToBytesList(key,
		"13.361389", "38.115556", pos1,
		"15.087269", "37.502669", pos2,
	))
	result := GeoDist(testDB, utils.ToBytesList(key, pos1, pos2, "km"))
	bulkReply, ok := result.(*reply.BulkReply)
	if !ok {
		t.Error(fmt.Sprintf("expected bulk reply, actually %s", result.ToBytes()))
		return
	}
	dist, err := strconv.ParseFloat(string(bulkReply.Arg), 10)
	if err != nil {
		t.Error(err)
		return
	}
	if dist < 166.274 || dist > 166.275 {
		t.Errorf("expected 166.274, actual: %f", dist)
	}

	result = GeoDist(testDB, utils.ToBytesList(key, pos1, pos2, "m"))
	bulkReply, ok = result.(*reply.BulkReply)
	if !ok {
		t.Error(fmt.Sprintf("expected bulk reply, actually %s", result.ToBytes()))
		return
	}
	dist, err = strconv.ParseFloat(string(bulkReply.Arg), 10)
	if err != nil {
		t.Error(err)
		return
	}
	if dist < 166274 || dist > 166275 {
		t.Errorf("expected 166274, actual: %f", dist)
	}
}
