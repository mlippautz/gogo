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
    RestCounter uint64; // depr.
    IncCnt uint64;
    GlobalCounter uint64; // depr.
    Prefix string; // depr.
    CurFile string;
    CurLine uint64;
    T uint64;
    F uint64;
    TDepth uint64;
    FDepth uint64;
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
    var jmp string;

    if Compile != 0 {
        if item.Mode != libgogo.MODE_COND {
            GenErrorWeak("Can use relative operators only with conditionals.");
        }
        if op == TOKEN_REL_AND {
            labelString = GenerateSubLabel(ed,0,"END");
            jmp = GetJump(item.C, 1);
            PrintJump(jmp, labelString);
            if ed.T != 0 {
                if ed.TDepth >= ed.ExpressionDepth {
                    labelString = GetSubLabel(ed,1,"END");
                    PrintLabel(labelString);
                    ed.T = 0;
                }
            }
        } else {
            if op == TOKEN_REL_OR {
                labelString = GenerateSubLabel(ed,1,"END");
                jmp = GetJump(item.C, 0);
                PrintJump(jmp, labelString);
                if ed.F != 0 {
                    if ed.FDepth >= ed.ExpressionDepth {
                        labelString = GetSubLabel(ed,0,"END");
                        PrintLabel(labelString);
                        ed.F = 0;
                    }
                }
            } else {
                GenErrorWeak("Relative AND or OR expected.");
            }
        }
    }
}

//
// Called by parser (ParseExpression)
// Note: This function can only handle operands with a maximum of 8 bytes in size
//
func GenerateComparison(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    if Compile != 0 {   
        //if (item1.Itemtype != uint64_t) || (item2.Itemtype != uint64_t) {
        //    SymbolTableError("Cannot compare types", "", "other than", uint64_t.Name);
        //}
        if (item1.Itemtype != item2.Itemtype) && (item1.Itemtype != string_t) && (item2.Itemtype != string_t) {
            GenErrorWeak("Can only compare variables of same type.");
        }            
        if (item1.Itemtype == string_t) || (item2.Itemtype == string_t) {
            GenErrorWeak("Cannot compare string types.");
        }
        if item1.PtrType == 1 {
            if item2.PtrType == 1 {
                if (op != TOKEN_EQUALS) && (op != TOKEN_NOTEQUAL) {
                    GenErrorWeak("Can only compare '==' or '!=' on pointers'");
                }
            } else {
                GenErrorWeak("Pointer to non-pointer comparison.");
            }
        }
        if (item2.PtrType == 1) && (item1.PtrType != 1) {
            GenErrorWeak("Non-pointer to pointer comparison.");
        }

        // Generate CMP statements depending on items
        if item1.Mode == libgogo.MODE_CONST {
            if item2.Mode == libgogo.MODE_CONST {
                item1.Itemtype = bool_t;
                item1.A = GetConditionalBool(op, item1.A, item2.A);
            } else {
                if item1.PtrType == 0 {
                    MakeRegistered(item1, 0);
                } else {
                    MakeRegistered(item1, 1);
                }
                if item2.Mode == libgogo.MODE_REG {
                    PrintInstruction_Reg_Reg("CMP", 8, "R", item1.R, 0, 0, 0, "", "R", item2.R, 0, 0, 0, "");
                }
                if item2.Mode == libgogo.MODE_VAR {
                    if item2.PtrType == 1 {
                        MakeRegistered(item2, 1);
                        PrintInstruction_Reg_Reg("CMP", 8, "R", item1.R, 0, 0, 0, "", "R", item2.R, 0, 0, 0, "");
                    } else {
                        PrintInstruction_Reg_Var("CMP", "R", item1.R, "", 0, item2);
                    }
                }
            }
        }
        if item1.Mode == libgogo.MODE_REG {
            DereferRegisterIfNecessary(item1);
            if item2.Mode == libgogo.MODE_CONST {
                PrintInstruction_Reg_Imm("CMP", 8, "R", item1.R, 0, 0, 0, "", item2.A);
            }
            if item2.Mode == libgogo.MODE_REG {
                PrintInstruction_Reg_Reg("CMP", 8, "R", item1.R, 0, 0, 0, "", "R", item2.R, 0, 0, 0, "");
            }
            if item2.Mode == libgogo.MODE_VAR {
                PrintInstruction_Reg_Var("CMP", "R", item1.R, "", 0, item2);
            }
        }
        if item1.Mode == libgogo.MODE_VAR {
            if item2.Mode == libgogo.MODE_CONST {
                PrintInstruction_Var_Imm("CMP", item1, item2.A);
            }
            if item2.Mode == libgogo.MODE_REG {
                DereferRegisterIfNecessary(item1);
                if item1.PtrType == 1 {
                    MakeRegistered(item1, 1);
                    PrintInstruction_Reg_Reg("CMP", 8, "R", item1.R, 0, 0, 0, "", "R", item2.R, 0, 0, 0, "");
                } else {
                    PrintInstruction_Var_Reg("CMP", item1, "R", item2.R, "", 0);
                }
            }
            if item2.Mode == libgogo.MODE_VAR {
                if item2.PtrType == 0 {
                    MakeRegistered(item2, 0);
                } else {
                    MakeRegistered(item2, 1);
                }
                if item1.PtrType == 1 {
                    MakeRegistered(item1, 1);
                    PrintInstruction_Reg_Reg("CMP", 8, "R", item1.R, 0, 0, 0, "", "R", item2.R, 0, 0, 0, "");
                } else {
                    PrintInstruction_Var_Reg("CMP", item1, "R", item2.R, "", 0);
                }
            }
        }

        // Prepare item
        item1.Itemtype = bool_t;
        item1.Mode = libgogo.MODE_COND;
        if op == TOKEN_EQUALS {
            item1.C = libgogo.REL_EQ;
        }
        if op == TOKEN_NOTEQUAL {
            item1.C = libgogo.REL_NEQ;
        }
        if op == TOKEN_REL_LT {
            item1.C = libgogo.REL_LT;
        }
        if op == TOKEN_REL_LTOE {
            item1.C = libgogo.REL_LTE;
        }
        if op == TOKEN_REL_GT {
            item1.C = libgogo.REL_GT;
        }
        if op == TOKEN_REL_GTOE {
            item1.C = libgogo.REL_GTE;
        }

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
