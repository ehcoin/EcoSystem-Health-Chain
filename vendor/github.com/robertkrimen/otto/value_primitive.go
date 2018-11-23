// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package otto

func toStringPrimitive(value Value) Value {
	return _toPrimitive(value, defaultValueHintString)
}

func toNumberPrimitive(value Value) Value {
	return _toPrimitive(value, defaultValueHintNumber)
}

func toPrimitive(value Value) Value {
	return _toPrimitive(value, defaultValueNoHint)
}

func _toPrimitive(value Value, hint _defaultValueHint) Value {
	switch value.kind {
	case valueNull, valueUndefined, valueNumber, valueString, valueBoolean:
		return value
	case valueObject:
		return value._object().DefaultValue(hint)
	}
	panic(hereBeDragons(value.kind, value))
}
