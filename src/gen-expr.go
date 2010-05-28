// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// Expression related code generation
//

package main

import "./libgogo/_obj/libgogo"

type ExpressionDescriptor struct {
    ExpressionDepth uint64;
    RestCounter uint64;
    GlobalCounter uint64;
    Prefix string;
};

//
// Called by parser (ParseSimpleExpression)
//
func GenerateSimpleExpressionArith(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    if Compile != 0 {
        if (item1.Itemtype != byte_t) && (item1.Itemtype != uint64_t) {
            SymbolTableError("Invalid left operand type for", "", "addition/subtraction:", item1.Itemtype.Name);
        }
        if (item2.Itemtype != byte_t) && (item2.Itemtype != uint64_t) {
            SymbolTableError("Invalid right operand type for", "", "addition/subtraction:", item2.Itemtype.Name);
        }        
        if op == TOKEN_ARITH_PLUS { //Add
            AddSubInstruction("ADD", item1, item2, item1.A + item2.A, 0);
        } else { //Subtract
            AddSubInstruction("SUB", item1, item2, item1.A - item2.A, 0);
        }
    }
}

//
// Called by parser (ParseTerm)
//
func GenerateTermArith(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    if Compile != 0 {
        if (op == TOKEN_ARITH_DIV) || op == (TOKEN_ARITH_MUL) {
            if (item1.Itemtype != byte_t) && (item1.Itemtype != uint64_t) {
                SymbolTableError("Invalid left operand type for", "", "multiplication/division:", item1.Itemtype.Name);
            }
            if (item2.Itemtype != byte_t) && (item2.Itemtype != uint64_t) {
                SymbolTableError("Invalid right operand type for", "", "multiplication/division:", item2.Itemtype.Name);
            }
        }
        if op == TOKEN_ARITH_DIV { // Division
            if item2.Mode == libgogo.MODE_CONST {
                if item2.A == 0 {
                    GenErrorWeak("Division by zero.");
                }
            }
            if item2.A != 0 { //Avoid division by zero for constvalue parameter
                DivMulInstruction("DIV", item1, item2, item1.A / item2.A, 0);
            } else {
                DivMulInstruction("DIV", item1, item2, 0, 0);
            }
        }
        if op == TOKEN_ARITH_MUL { // Multiplication
            DivMulInstruction("MUL", item1, item2, item1.A * item2.A, 0);
        }
    }
}

func GenerateRelative(item *libgogo.Item, op uint64, ed *ExpressionDescriptor) {
    var labelString string;
    if op == TOKEN_REL_AND {
        labelString = GenerateIfLabel(ed.Prefix,ed.GlobalCounter,ed.ExpressionDepth,"END");
        PrintJump("JNZ", labelString);
    } else {
        if op == TOKEN_REL_OR {
            labelString = GenerateIfLabel(ed.Prefix,ed.GlobalCounter,0,"OK");
            PrintJump("JZ", labelString);
            labelString = GenerateIfLabel(ed.Prefix,ed.GlobalCounter,ed.ExpressionDepth-1,"END");
            PrintLabel(labelString);
        } else {
            GenErrorWeak("Relative AND expected.");
        }
    }
}

//
// Called by parser (ParseExpression)
//
func GenerateComparison(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    if Compile != 0 {   
        if (item1.Itemtype != uint64_t) || (item2.Itemtype != uint64_t) {
            SymbolTableError("Cannot compare types", "", "other than", uint64_t.Name);
        }    
        DereferItemIfNecessary(item1); //Derefer address if item is a pointer (should not be necessary here)
        DereferItemIfNecessary(item2); //Derefer address if item is a pointer (should not be necessary here)

        // Generate CMP statements depending on items
        if item1.Mode == libgogo.MODE_CONST {
            if item2.Mode == libgogo.MODE_CONST {
                item1.Itemtype = bool_t;
                item1.A = GetConditionalBool(op, item1.A, item2.A);
            } else {
                MakeRegistered(item1, 0);
                if item2.Mode == libgogo.MODE_REG {
                    PrintInstruction_Reg_Reg("CMP", 8, "R", item1.R, 0, 0, 0, "", "R", item2.R, 0, 0, 0, "");
                }
                if item2.Mode == libgogo.MODE_VAR {
                    PrintInstruction_Reg_Var("CMP", "R", item1.R, item2);
                }
            }
        }
        if item1.Mode == libgogo.MODE_REG {
            if item2.Mode == libgogo.MODE_CONST {
                PrintInstruction_Reg_Imm("CMP", 8, "R", item1.R, 0, 0, 0, "", item2.A);
            }
            if item2.Mode == libgogo.MODE_REG {
                PrintInstruction_Reg_Reg("CMP", 8, "R", item1.R, 0, 0, 0, "", "R", item2.R, 0, 0, 0, "");
            }
            if item2.Mode == libgogo.MODE_VAR {
                PrintInstruction_Reg_Var("CMP", "R", item1.R, item2);
            }
        }
        if item1.Mode == libgogo.MODE_VAR {
            if item2.Mode == libgogo.MODE_CONST {
                PrintInstruction_Var_Imm("CMP", item1, item2.A);
            }
            if item2.Mode == libgogo.MODE_REG {
                PrintInstruction_Var_Reg("CMP", item1, "R", item2.R);
            }
            if item2.Mode == libgogo.MODE_VAR {
                MakeRegistered(item2, 0);
                PrintInstruction_Var_Reg("CMP", item1, "R", item2.R);
            }
        }

        item1.Itemtype = bool_t;
        FreeRegisterIfRequired(item2);    
    }
}

func GetConditionalBool(op uint64, val1 uint64, val2 uint64) uint64 {
    var ret uint64;    
    if op == TOKEN_EQUALS {
        if val1 == val2 {
            ret = 1;
        } else {
            ret = 0;
        }
    }
    if op == TOKEN_NOTEQUAL {
        if val1 == val2 {
            ret = 0;
        } else {
            ret = 1;
        }
    }
    if op == TOKEN_REL_GTOE {
        if val1 >= val2 {
            ret = 1;
        } else {
            ret = 0;
        }
    }
    if op == TOKEN_REL_LTOE {
        if val1 <= val2 {
            ret = 1;
        } else {
            ret = 0;
        }
    }
    if op == TOKEN_REL_GT {
        if val1 > val2 {
            ret = 1;
        } else {
            ret = 0;
        }
    }
    if op == TOKEN_REL_LT {
        if val1 < val2 {
            ret = 1;
        } else {
            ret = 0;
        }
    }
    return ret;
}
