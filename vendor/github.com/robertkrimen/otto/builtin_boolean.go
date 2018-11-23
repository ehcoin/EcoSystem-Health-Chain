// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package otto

// Boolean

func builtinBoolean(call FunctionCall) Value {
	return toValue_bool(call.Argument(0).bool())
}

func builtinNewBoolean(self *_object, argumentList []Value) Value {
	return toValue_object(self.runtime.newBoolean(valueOfArrayIndex(argumentList, 0)))
}

func builtinBoolean_toString(call FunctionCall) Value {
	value := call.This
	if !value.IsBoolean() {
		// Will throw a TypeError if ThisObject is not a Boolean
		value = call.thisClassObject("Boolean").primitiveValue()
	}
	return toValue_string(value.string())
}

func builtinBoolean_valueOf(call FunctionCall) Value {
	value := call.This
	if !value.IsBoolean() {
		value = call.thisClassObject("Boolean").primitiveValue()
	}
	return value
}
