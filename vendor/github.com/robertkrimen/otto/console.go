// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package otto

import (
	"fmt"
	"os"
	"strings"
)

func formatForConsole(argumentList []Value) string {
	output := []string{}
	for _, argument := range argumentList {
		output = append(output, fmt.Sprintf("%v", argument))
	}
	return strings.Join(output, " ")
}

func builtinConsole_log(call FunctionCall) Value {
	fmt.Fprintln(os.Stdout, formatForConsole(call.ArgumentList))
	return Value{}
}

func builtinConsole_error(call FunctionCall) Value {
	fmt.Fprintln(os.Stdout, formatForConsole(call.ArgumentList))
	return Value{}
}

// Nothing happens.
func builtinConsole_dir(call FunctionCall) Value {
	return Value{}
}

func builtinConsole_time(call FunctionCall) Value {
	return Value{}
}

func builtinConsole_timeEnd(call FunctionCall) Value {
	return Value{}
}

func builtinConsole_trace(call FunctionCall) Value {
	return Value{}
}

func builtinConsole_assert(call FunctionCall) Value {
	return Value{}
}

func (runtime *_runtime) newConsole() *_object {

	return newConsoleObject(runtime)
}
