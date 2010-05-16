// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

var NumRegisters uint64 = 8;
var FreeRegisters [8]byte;

func InitFreeRegisters() {
    var i uint64;
    for i = 0; i < NumRegisters; i = i + 1 {
        FreeRegisters[i] = 1;
    }
}

func GetFreeRegister() uint64 {
    var i uint64;
    for i = 0; FreeRegisters[i] == 0; {
        i = i + 1;
        if i == NumRegisters {
            libgogo.ExitError("No more free registers available for code generation", 5);
        }
    }
    return i;
}

func OccupyRegister(index uint64) {
    FreeRegisters[index] = 0;
}

func FreeRegister(index uint64) {
    FreeRegisters[index] = 1;
}

func PrintFreeRegisters() {
    /*var i uint64;
    libgogo.PrintString("Free registers: ");
    for i = 0; i < NumRegisters; i++ {
        if FreeRegisters[i] == 1 {
            libgogo.PrintString("R");
            libgogo.PrintNumber(i);
            libgogo.PrintString(",");
        }
    }
    libgogo.PrintString("\b\n");*/
}

func ItemToRegister(item *libgogo.Item) {
    if item.Mode == libgogo.MODE_CONST {
        item.R = GetFreeRegister();
        OccupyRegister(item.R);
        item.Mode = libgogo.MODE_REG;
        libgogo.PrintString("MOVQ $");
        libgogo.PrintNumber(item.A);
        libgogo.PrintString(", R");
        libgogo.PrintNumber(item.R);
        libgogo.PrintString("\n");
    }
    if item.Mode == libgogo.MODE_VAR {
        item.R = GetFreeRegister();
        OccupyRegister(item.R);
        item.Mode = libgogo.MODE_REG;
        libgogo.PrintString("MOVQ ");
        if item.Global == 1 { //Global
            libgogo.PrintNumber(item.A);
            libgogo.PrintString("(SB)");
        } else { //Local
            libgogo.PrintString("-");
            libgogo.PrintNumber(item.A + 8); //SP = return address, start at address SP-8, decreasing
            libgogo.PrintString("(SP)");
        }
        libgogo.PrintString(", R");
        libgogo.PrintNumber(item.R);
        libgogo.PrintString("\n");
    }
    if item.Mode == libgogo.MODE_REG {
				; //Don't do anything - item is already a register
    }
}

func ItemAddressToRegister(item *libgogo.Item) {
    if item.Mode == libgogo.MODE_VAR {
        item.R = GetFreeRegister();
        OccupyRegister(item.R);
        item.Mode = libgogo.MODE_REG;
        libgogo.PrintString("LEAQ ");
        if item.Global == 1 { //Global
            libgogo.PrintNumber(item.A);
            libgogo.PrintString("(SB)");
        } else { //Local
            libgogo.PrintString("-");
            libgogo.PrintNumber(item.A + 8); //SP = return address, start at address SP-8, decreasing
            libgogo.PrintString("(SP)");
        }
        libgogo.PrintString(", R");
        libgogo.PrintNumber(item.R);
        libgogo.PrintString("\n");
    }
    if item.Mode == libgogo.MODE_REG {
				; //Don't do anything - item is already a register
    }
}

func GenerateTerm(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    var str string;
    if Compile != 0 {
				if (item1.Mode == libgogo.MODE_CONST) && (item2.Mode == libgogo.MODE_CONST) {
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
				}
    }
}

func GenerateSimpleExpression(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
		var str string;
    if Compile != 0 {
				if (item1.Mode == libgogo.MODE_CONST) && (item2.Mode == libgogo.MODE_CONST) {
						str = TokenToString(op);
				    libgogo.PrintString(";Constant folding: ");
				    libgogo.PrintString(str);
				    libgogo.PrintString("(");
				    libgogo.PrintNumber(item1.A);
				    libgogo.PrintString(",");
				    libgogo.PrintNumber(item2.A);
				    libgogo.PrintString(")=");
				    if op == TOKEN_ARITH_PLUS {
				        item1.A = item1.A + item2.A;
				    }
				    if op == TOKEN_ARITH_MINUS {
				        item1.A = item1.A - item2.A;
				    }
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
            if op == TOKEN_ARITH_PLUS {
                libgogo.PrintString("ADDQ AX, BX\n");
            }
            if op == TOKEN_ARITH_MINUS {
                libgogo.PrintString("SUBQ AX, BX\n");
            }
            libgogo.PrintString("MOVQ AX, R");
            libgogo.PrintNumber(item1.R);
            libgogo.PrintString("\n");
            FreeRegister(item2.R);
        }
    }
}

func GenerateFieldAccess(item *libgogo.Item, offset uint64, indirect uint64) {
    if Compile != 0 {
        if (indirect != 0) || (offset != 0) { //If offset 0 on direct access => no change
				    ItemAddressToRegister(item);
						libgogo.PrintString("MOVQ R");
						libgogo.PrintNumber(item.R);
						libgogo.PrintString(", AX\n");
				    if offset != 0 {
				        libgogo.PrintString("ADDQ AX, $");
						    libgogo.PrintNumber(offset); //Add offset
				        libgogo.PrintString("\n");
				    }
				    if indirect != 0 { //Indirect access
				        libgogo.PrintString("MOVQ (AX), AX\n");
				    }
				    libgogo.PrintString("MOVQ AX, R");
				    libgogo.PrintNumber(item.R);
				    libgogo.PrintString("\n");
        }
    }
}

func GenerateAssignment(LHSItem *libgogo.Item, RHSItem *libgogo.Item) {
    if Compile != 0 {
        ItemToRegister(RHSItem);
        ItemAddressToRegister(LHSItem);
				libgogo.PrintString("MOVQ R");
				libgogo.PrintNumber(RHSItem.R);
		    libgogo.PrintString(", (R");
				libgogo.PrintNumber(LHSItem.R);
		    libgogo.PrintString(")"); //Move content of RHS (in register RHS.R) to address of LHS (referred to by register LHS.R)
		    libgogo.PrintString("\n");
        FreeRegister(RHSItem.R);
        FreeRegister(LHSItem.R);
    }
}
