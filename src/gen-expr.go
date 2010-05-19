// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// Expression related code generation
//

package main

import "./libgogo/_obj/libgogo"

//
// Called by parser (ParseSimpleExpression)
//
func GenerateSimpleExpression(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    if Compile != 0 {
        if op == TOKEN_ARITH_PLUS { //Add
            AddSubInstruction("ADDQ", item1, item2, item1.A + item2.A, 0);
        } else { //Subtract
            AddSubInstruction("SUBQ", item1, item2, item1.A - item2.A, 0);
        }
    }
}

//
// Called by parser (ParseTerm)
//
func GenerateTerm(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    if Compile != 0 {
        if op == TOKEN_ARITH_DIV { // Division
            if item2.Mode == libgogo.MODE_CONST {
                if item2.A == 0 {
                    GenErrorWeak("Division by zero.");
                }
            }
            if item2.A != 0 {
                DivMulInstruction("DIVQ", item1, item2, item1.A / item2.A, 0);
            } else {
                DivMulInstruction("DIVQ", item1, item2, 0, 0);
            }
        }
        if op == TOKEN_ARITH_MUL { // Multiplication
            DivMulInstruction("MULQ", item1, item2, item1.A * item2.A, 0);
        }
    }
}
