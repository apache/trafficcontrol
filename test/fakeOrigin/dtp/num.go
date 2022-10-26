package dtp

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"fmt"
	"math/rand"
	"strconv"
	"unicode"
)

/* Operators:
 *   a,+: add
 *   s,-: subtract
 *   m,*: multiply
 *   d: divide
 *   u: modulo
 *   ,: separate
 *   k: multiply by 1024
 *   M: multiply by 1024^2
 *   G: multiply by 1024^3
 *   T: multiply by 1024^4
 *   P: multiply by 1024^5
 *   r: Linear Distribution
 *   (: Start Selection
 *   ): Perform Selection
 *   w: Add weight
 */

func EvalNumber(numStr string, randSeed int64) int64 {
	rnd := rand.New(rand.NewSource(randSeed))
	ops := LexNumberStr(numStr)
	DebugLogf("Evaluating %s: %v\n", numStr, ops)
	var stack []Number
	for _, op := range ops {
		DebugLogf("Evaluating %v on %v ...", op, stack)
		pop, push := op.Evaluate(stack, rnd)
		DebugLogf("Popping %d, Pushing %v\n", pop, push)

		posend := len(stack) - pop
		if posend < 0 {
			return 0
		}

		stack = stack[:posend]
		stack = append(stack, push...)
	}
	if len(stack) == 0 {
		return 0
	}
	return stack[len(stack)-1].Value
}

type Evaluable interface {
	Evaluate(stack []Number, rnd *rand.Rand) (pop int, push []Number)
}

type Operator byte
type Literal int64

type Number struct {
	Value    int64
	Weight   int64
	Sentinel bool
}

func (o Operator) String() string {
	return string(byte(o))
}

func (n Number) String() string {
	if n.Sentinel {
		return "!"
	}
	if n.Weight != 1 {
		return fmt.Sprintf("%d,%dw", n.Value, n.Weight)
	}
	return fmt.Sprint(n.Value)
}

func (o Operator) Evaluate(stack []Number, rnd *rand.Rand) (pop int, push []Number) {
	rd := func(i int) Number {
		if i > len(stack) || i < 1 {
			return Number{Weight: 1}
		}
		return stack[len(stack)-i]
	}
	var c int64 = 1
	switch o {
	case 'a', '+':
		a, b := rd(1), rd(2)
		if a.Sentinel || b.Sentinel {
			return 0, nil
		}
		return 2, []Number{
			{
				Value:  b.Value + a.Value,
				Weight: a.Weight,
			},
		}
	case 's', '-':
		a, b := rd(1), rd(2)
		if a.Sentinel || b.Sentinel {
			return 0, nil
		}
		return 2, []Number{
			{
				Value:  b.Value - a.Value,
				Weight: a.Weight,
			},
		}
	case 'm', '*':
		a, b := rd(1), rd(2)
		if a.Sentinel || b.Sentinel {
			return 0, nil
		}
		return 2, []Number{
			{
				Value:  b.Value * a.Value,
				Weight: a.Weight,
			},
		}
	case 'd':
		a, b := rd(1), rd(2)
		if a.Sentinel || b.Sentinel {
			return 0, nil
		}
		return 2, []Number{
			{
				Value:  b.Value / a.Value,
				Weight: a.Weight,
			},
		}
	case 'u':
		a, b := rd(1), rd(2)
		if a.Sentinel || b.Sentinel {
			return 0, nil
		}
		return 2, []Number{
			{
				Value:  b.Value % a.Value,
				Weight: a.Weight,
			},
		}
	case ',':
		return 0, nil
	case 'P':
		c *= 1024
		fallthrough
	case 'T':
		c *= 1024
		fallthrough
	case 'G':
		c *= 1024
		fallthrough
	case 'M':
		c *= 1024
		fallthrough
	case 'k':
		c *= 1024
		a := rd(1)
		if a.Sentinel {
			return 0, nil
		}
		return 1, []Number{
			{
				Value:  a.Value * c,
				Weight: a.Weight,
			},
		}
	case 'r':
		a, b := rd(1), rd(2)
		if a.Sentinel || b.Sentinel {
			return 0, nil
		}
		var min, max int64
		if a.Value < b.Value {
			min = a.Value
			max = b.Value
		} else {
			min = b.Value
			max = a.Value
		}
		return 2, []Number{
			{
				Value:  rnd.Int63n(max-min) + min,
				Weight: a.Weight,
			},
		}
	case '(':
		return 0, []Number{{Sentinel: true}}
	case ')':
		var wt int64
		var ct int
		var sentinelFound int
		for i := 1; i <= len(stack); i++ {
			x := rd(i)
			if x.Sentinel {
				sentinelFound++
				break
			}

			wt += x.Weight
			ct++
		}
		wt_r := rnd.Int63n(wt)
		DebugLogf("(w%d/%d*%d) ", wt_r, wt, ct)
		var n int64
		for i := int(1); i <= ct; i++ {
			x := rd(i)
			wt_r -= x.Weight
			DebugLogf("(w%d/%d*%d %d) ", wt_r, wt, i, x.Value)
			if wt_r < 0 {
				n = x.Value
				break
			}
		}
		return ct + sentinelFound, []Number{{Value: n, Weight: 1}}
	case 'w':
		return 2, []Number{{Value: rd(2).Value, Weight: rd(1).Value}}
	}
	return 0, nil
}

func (n Literal) Evaluate(stack []Number, rnd *rand.Rand) (pop int, push []Number) {
	return 0, []Number{
		{
			Value:  int64(n),
			Weight: 1,
		},
	}
}

func LexNumberStr(s string) (ops []Evaluable) {
	n := ``
	addNum := func() {
		if len(n) > 0 {
			num, _ := strconv.ParseInt(n, 10, 64)
			ops = append(ops, Literal(num))
			n = ``
		}
	}
	addOp := func(c byte) {
		ops = append(ops, Operator(c))
	}
	for _, c := range s {
		if c&^rune(0xff) != 0 {
			continue
		}
		if !unicode.IsDigit(c) {
			addNum()
			addOp(byte(c))
		} else {
			n += string(c)
		}
	}
	addNum()
	return
}
