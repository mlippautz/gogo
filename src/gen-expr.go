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
            TwoOperandInstruction("ADDQ", item1, item2, item1.A + item2.A, 0);
        } else { //Subtract
            TwoOperandInstruction("SUBQ", item1, item2, item1.A - item2.A, 0);
        }
    }
}

//
// Called by parser (ParseTerm)
//
func GenerateTerm(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    /*var str string;*/
    if Compile != 0 {
				/*if (item1.Mode == libgogo.MODE_CONST) && (item2.Mode == libgogo.MODE_CONST) {
				    libgogo.PrintString(";Constant folding: ");
						str = TokenToString(op);
				    libgogo.PrintString(str);
				    libgogo.PrintString("(");
				    libgogo.PrintNumber(item1.A);
				    libgogo.PrintString(",");
				    libgogo.PrintNumber(item2.A);
				    if op == TOKEN_ARITH_MUL {
				        item1.A = item1.A * item2.A;
				    }
				    if op == TOKEN_ARITH_DIV {
				        item1.A = item1.A / item2.A;
				    }
				    libgogo.PrintString(")=");
				    libgogo.PrintNumber(item1.A);
				    libgogo.PrintString("\n");
				} else {
				    ItemToRegister(item1);
				    libgogo.PrintString("MOVQ R");
				    libgogo.PrintNumber(item1.R);
				    libgogo.PrintString(", AX\n");
				    ItemToRegister(item2);
				    libgogo.PrintString("MOVQ R");
				    libgogo.PrintNumber(item2.R);
				    libgogo.PrintString(", BX\n");
				    if op == TOKEN_ARITH_MUL {
				        libgogo.PrintString("MULQ BX\n");
				    }
				    if op == TOKEN_ARITH_DIV {
				        libgogo.PrintString("DIVQ BX\n");
				    }
				    libgogo.PrintString("MOVQ AX, R");
				    libgogo.PrintNumber(item1.R);
				    libgogo.PrintString("\n");
				    FreeRegister(item2.R);
				}*/
    }
}


