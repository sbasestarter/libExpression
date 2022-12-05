package libexpression

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type opCls1 struct {
	S1, S2 string
}

type opCls2 struct {
	N1, N2 int
}

func TestBase1(t *testing.T) {
	UpdateOpValueTypeMap(reflect.TypeOf(opCls1{}).Name(), reflect.TypeOf(opCls1{}))
	UpdateOpValueTypeMap(reflect.TypeOf(opCls2{}).Name(), reflect.TypeOf(opCls2{}))

	v1 := Op{
		OpType:    OpTypeValue,
		ValueType: reflect.TypeOf(opCls1{}).Name(),
		Value:     &opCls1{"s1", "s2"},
	}

	v2 := Op{
		OpType:    OpTypeValue,
		ValueType: reflect.TypeOf(opCls2{}).Name(),
		Value:     &opCls2{1, 2},
	}

	v3 := Op{
		OpType:    OpTypeValue,
		ValueType: reflect.TypeOf(opCls1{}).Name(),
		Value:     &opCls1{"s31", "s32"},
	}

	or := Op{
		OpType: OpTypeOr,
		Values: []*Op{
			&v1, &v2,
		},
	}

	and := Op{
		OpType: OpTypeAnd,
		Values: []*Op{
			&or, &v3,
		},
	}

	d, err := json.Marshal(and)
	assert.Nil(t, err)

	t.Log(string(d))

	var and2 Op
	err = json.Unmarshal(d, &and2)
	assert.Nil(t, err)

	d, err = json.Marshal(and2)
	assert.Nil(t, err)

	t.Log(string(d))

	ops, result := Check(&and2)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 2, len(ops))
	ops[1].Flag = MarkFlagTrue

	ops, result = Check(&and2)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 1, len(ops))
	ops[0].Flag = MarkFlagTrue

	ops, result = Check(&and2)
	assert.EqualValues(t, CheckResultTrue, result)
	assert.EqualValues(t, 0, len(ops))
}

type utE int

const (
	utE1 utE = iota
	utE2
)

type utCls1 struct {
	S     string
	E     utE
	Inter interface{}
}

type utCls2 struct {
	I interface{}
}

func TestBase3(t *testing.T) {
	v := &utCls2{
		I: &utCls1{S: "s1", E: utE1, Inter: &utCls1{S: "s2"}},
	}

	d, err := json.Marshal(v)
	assert.Nil(t, err)

	t.Log(string(d))

	var v2 utCls2

	err = json.Unmarshal(d, &v2) // => map
	assert.Nil(t, err)

	d, err = json.Marshal(v)
	assert.Nil(t, err)

	t.Log(string(d))
}

func TestAnd1(t *testing.T) {
	and := &Op{
		OpType: OpTypeAnd,
		Values: []*Op{
			{
				OpType: OpTypeValue,
				Value:  "a1",
			},
			{
				OpType: OpTypeValue,
				Value:  "a2",
			},
		},
	}

	ops, result := Check(and)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 1, len(ops))
	assert.EqualValues(t, "a1", ops[0].Value)

	ops, result = Check(and)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 1, len(ops))
	assert.EqualValues(t, "a1", ops[0].Value)

	ops[0].Flag = MarkFlagTrue

	ops, result = Check(and)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 1, len(ops))
	assert.EqualValues(t, "a2", ops[0].Value)

	ops[0].Flag = MarkFlagTrue

	ops, result = Check(and)
	assert.EqualValues(t, CheckResultTrue, result)
	assert.EqualValues(t, 0, len(ops))
}

func TestAnd2(t *testing.T) {
	and := &Op{
		OpType: OpTypeAnd,
		Values: []*Op{
			{
				OpType: OpTypeValue,
				Value:  "a1",
			},
			{
				OpType: OpTypeValue,
				Value:  "a2",
			},
		},
	}

	ops, result := Check(and)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 1, len(ops))
	assert.EqualValues(t, "a1", ops[0].Value)

	ops[0].Flag = MarkFlagFalse

	ops, result = Check(and)
	assert.EqualValues(t, CheckResultFalse, result)
	assert.EqualValues(t, 0, len(ops))
}

func TestOr1(t *testing.T) {
	or := &Op{
		OpType: OpTypeOr,
		Values: []*Op{
			{
				OpType: OpTypeValue,
				Value:  "a1",
			},
			{
				OpType: OpTypeValue,
				Value:  "a2",
			},
		},
	}

	ops, result := Check(or)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 2, len(ops))
	assert.EqualValues(t, "a1", ops[0].Value)
	assert.EqualValues(t, "a2", ops[1].Value)

	ops, result = Check(or)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 2, len(ops))
	assert.EqualValues(t, "a1", ops[0].Value)
	assert.EqualValues(t, "a2", ops[1].Value)

	ops[0].Flag = MarkFlagTrue

	ops, result = Check(or)
	assert.EqualValues(t, CheckResultTrue, result)
	assert.EqualValues(t, 0, len(ops))
}

func TestOr2(t *testing.T) {
	or := &Op{
		OpType: OpTypeOr,
		Values: []*Op{
			{
				OpType: OpTypeValue,
				Value:  "a1",
			},
			{
				OpType: OpTypeValue,
				Value:  "a2",
			},
		},
	}

	ops, result := Check(or)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 2, len(ops))
	assert.EqualValues(t, "a1", ops[0].Value)
	assert.EqualValues(t, "a2", ops[1].Value)

	ops, result = Check(or)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 2, len(ops))
	assert.EqualValues(t, "a1", ops[0].Value)
	assert.EqualValues(t, "a2", ops[1].Value)

	ops[0].Flag = MarkFlagFalse

	ops, result = Check(or)
	assert.EqualValues(t, CheckResultNeedOp, result)
	assert.EqualValues(t, 1, len(ops))
	assert.EqualValues(t, "a2", ops[0].Value)

	ops[0].Flag = MarkFlagTrue

	ops, result = Check(or)
	assert.EqualValues(t, CheckResultTrue, result)
	assert.EqualValues(t, 0, len(ops))
}

// nolint
func TestX(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	fnRandBool := func() bool {
		return rand.Int63()%2 == 0 // nolint: gosec
	}

	fnRandN := func(n int) (t bool, op *Op) {
		andOp := fnRandBool()
		if andOp {
			op = &Op{
				OpType: OpTypeAnd,
			}
			t = true
		} else {
			op = &Op{
				OpType: OpTypeOr,
			}
		}

		for idx := 0; idx < n; idx++ {
			curOp := &Op{
				OpType: OpTypeValue,
			}

			if fnRandBool() {
				curOp.Flag = MarkFlagTrue

				if !andOp {
					t = true
				}
			} else {
				curOp.Flag = MarkFlagFalse

				if andOp {
					t = false
				}
			}

			op.Values = append(op.Values, curOp)
		}

		return
	}

	for idx := 0; idx < 1000; idx++ {
		r, op := fnRandN(4)
		_, cr := Check(op)
		assert.Equal(t, r, cr == CheckResultTrue)
	}

	ds := make([]struct {
		t  bool
		op *Op
	}, 0, 100)

	fnUpdate := func(t bool, op *Op) {
		if len(ds) >= 100 {
			ds = ds[1:]
		}
		ds = append(ds, struct {
			t  bool
			op *Op
		}{t: t, op: op})
	}

	fnSelect := func(n int) (tds []struct {
		t  bool
		op *Op
	}) {
		if n > len(ds) {
			n = len(ds)
		}

		for idx := 0; idx < n; idx++ {
			tds = append(tds, ds[rand.Intn(len(ds))])
		}

		return
	}

	for idx := 0; idx < 1000000; idx++ {
		n := rand.Intn(4)
		tsd := fnSelect(n)

		r, op := fnRandN(4 - n)

		for _, ts := range tsd {
			op.Values = append(op.Values, ts.op)

			switch op.OpType {
			case OpTypeAnd:
				r = r && ts.t
			case OpTypeOr:
				r = r || ts.t
			default:
				t.Fatal(op.OpType)
			}
		}

		_, cr := Check(op)
		assert.Equal(t, r, cr == CheckResultTrue)
		fnUpdate(r, op)
	}
}

// nolint
func TestX2(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	fnCheck := func(valueType string, value interface{}) CheckResult {
		f, ok := value.(bool)
		if !ok {
			t.Fatal("invalidBool")
		}
		if f {
			return CheckResultTrue
		}
		return CheckResultFalse
	}

	fnRandBool := func() bool {
		return rand.Int63()%2 == 0 // nolint: gosec
	}

	fnRandN := func(n int) (t bool, op *Op) {
		andOp := fnRandBool()
		if andOp {
			op = &Op{
				OpType: OpTypeAnd,
			}
			t = true
		} else {
			op = &Op{
				OpType: OpTypeOr,
			}
		}

		for idx := 0; idx < n; idx++ {
			curOp := &Op{
				OpType: OpTypeValue,
			}

			if fnRandBool() {
				curOp.Value = true

				if !andOp {
					t = true
				}
			} else {
				curOp.Value = false

				if andOp {
					t = false
				}
			}

			op.Values = append(op.Values, curOp)
		}

		return
	}

	for idx := 0; idx < 1000; idx++ {
		r, op := fnRandN(4)
		_, cr := CheckEx(op, fnCheck)
		assert.Equal(t, r, cr == CheckResultTrue)
	}

	ds := make([]struct {
		t  bool
		op *Op
	}, 0, 100)

	fnUpdate := func(t bool, op *Op) {
		if len(ds) >= 100 {
			ds = ds[1:]
		}
		ds = append(ds, struct {
			t  bool
			op *Op
		}{t: t, op: op})
	}

	fnSelect := func(n int) (tds []struct {
		t  bool
		op *Op
	}) {
		if n > len(ds) {
			n = len(ds)
		}

		for idx := 0; idx < n; idx++ {
			tds = append(tds, ds[rand.Intn(len(ds))])
		}

		return
	}

	for idx := 0; idx < 1000000; idx++ {
		n := rand.Intn(4)
		tsd := fnSelect(n)

		r, op := fnRandN(4 - n)

		for _, ts := range tsd {
			op.Values = append(op.Values, ts.op)

			switch op.OpType {
			case OpTypeAnd:
				r = r && ts.t
			case OpTypeOr:
				r = r || ts.t
			default:
				t.Fatal(op.OpType)
			}
		}

		_, cr := CheckEx(op, fnCheck)
		assert.Equal(t, r, cr == CheckResultTrue)
		fnUpdate(r, op)
	}
}
