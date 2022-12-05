package libexpression

type CheckResult int8

const (
	CheckResultUnspecified CheckResult = iota
	CheckResultNeedOp
	CheckResultTrue
	CheckResultFalse
	CheckResultInvalidExpression
	CheckResultInvalidLogic
)

type ExResultCalc func(valueType string, value interface{}) CheckResult

func Check(op *Op) ([]*Op, CheckResult) {
	return CheckEx(op, nil)
}

func CheckEx(op *Op, calc ExResultCalc) ([]*Op, CheckResult) {
	return check(op, calc)
}

func check(op *Op, calc ExResultCalc) (ops []*Op, r CheckResult) {
	if op == nil {
		r = CheckResultInvalidExpression

		return
	}

	switch op.OpType {
	case OpTypeAnd:
		return checkAnd(op, calc)
	case OpTypeOr:
		return checkOr(op, calc)
	case OpTypeValue:
		r = markFlag2CheckResult(op.Flag)
		if r == CheckResultNeedOp {
			if calc != nil {
				r = calc(op.ValueType, op.Value)
				if r != CheckResultNeedOp {
					return
				}
			}

			ops = append(ops, op)
		}

		return
	}

	r = CheckResultInvalidExpression

	return
}

func markFlag2CheckResult(f MarkFlag) CheckResult {
	switch f {
	case MarkFlagTrue:
		return CheckResultTrue
	case MarkFlagFalse:
		return CheckResultFalse
	case MarkFlagUnspecified:
		return CheckResultNeedOp
	}

	return CheckResultUnspecified
}

func checkAnd(op *Op, calc ExResultCalc) (ops []*Op, r CheckResult) {
	if len(op.Values) == 0 {
		r = CheckResultInvalidExpression

		return
	}

	for _, value := range op.Values {
		ops, r = check(value, calc)
		switch r {
		case CheckResultTrue:
			continue
		case CheckResultNeedOp:
			if len(ops) == 0 {
				r = CheckResultInvalidExpression
			}

			return
		default:
			ops = nil

			return
		}
	}

	if r != CheckResultTrue {
		r = CheckResultInvalidLogic
	}

	ops = nil

	return
}

func checkOr(op *Op, calc ExResultCalc) (ops []*Op, r CheckResult) {
	if len(op.Values) == 0 {
		r = CheckResultInvalidExpression

		return
	}

	for _, value := range op.Values {
		var tmpOps []*Op

		tmpOps, r = check(value, calc)
		switch r {
		case CheckResultTrue:
			ops = nil

			return
		case CheckResultNeedOp:
			ops = append(ops, tmpOps...)
		case CheckResultFalse:
			r = CheckResultNeedOp
		default:
			ops = nil

			return
		}
	}

	if r != CheckResultNeedOp {
		r = CheckResultInvalidLogic
	} else {
		if len(ops) == 0 {
			r = CheckResultFalse
		}
	}

	return
}
