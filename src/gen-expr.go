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
        if (item1.Itemtype != byte_t) && (item1.Itemtype != uint64_t) {
            SymbolTableError("Invalid left operand type for", "", "addition/subtraction:", item1.Itemtype.Name);
        }
        if (item2.Itemtype != byte_t) && (item2.Itemtype != uint64_t) {
            SymbolTableError("Invalid right operand type for", "", "addition/subtraction:", item2.Itemtype.Name);
        }
        //TODO: Consider special cases with byte_t requring ADDB and MOVB op codes
        
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
        if (item1.Itemtype != byte_t) && (item1.Itemtype != uint64_t) {
            SymbolTableError("Invalid left operand type for", "", "multiplication/division:", item1.Itemtype.Name);
        }
        if (item2.Itemtype != byte_t) && (item2.Itemtype != uint64_t) {
            SymbolTableError("Invalid right operand type for", "", "multiplication/division:", item2.Itemtype.Name);
        }
        //TODO: Consider special cases with byte_t requring ADDB and MOVB op codes
        
        if op == TOKEN_ARITH_DIV { // Division
            if item2.Mode == libgogo.MODE_CONST {
                if item2.A == 0 {
                    GenErrorWeak("Division by zero.");
                }
            }
            if item2.A != 0 { //Avoid division by zero for constvalue parameter
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

//
// Called by parser (ParseExpression)
//
func GenerateRelation(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    if Compile != 0 {   
        if (item1.Itemtype != uint64_t) || (item2.Itemtype != uint64_t) {
            SymbolTableError("Cannot compare types", "", "other than", uint64_t.Name);
        }
        //TODO: Consider byte_t
    
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
                    PrintInstruction_Reg_Reg("CMPQ", "R", item1.R, 0, 0, 0, "R", item2.R, 0, 0, 0);
                }
                if item2.Mode == libgogo.MODE_VAR {
                    PrintInstruction_Reg_Var("CMPQ", "R", item1.R, item2);
                }
            }
        }
        if item1.Mode == libgogo.MODE_REG {
            if item2.Mode == libgogo.MODE_CONST {
                PrintInstruction_Reg_Imm("CMPQ", "R", item1.R, 0, 0, 0, item2.A);
            }
            if item2.Mode == libgogo.MODE_REG {
                PrintInstruction_Reg_Reg("CMPQ", "R", item1.R, 0, 0, 0, "R", item2.R, 0, 0, 0);
            }
            if item2.Mode == libgogo.MODE_VAR {
                PrintInstruction_Reg_Var("CMPQ", "R", item1.R, item2);
            }
        }
        if item1.Mode == libgogo.MODE_VAR {
            if item2.Mode == libgogo.MODE_CONST {
                PrintInstruction_Var_Imm("CMPQ", item1, item2.A);
            }
            if item2.Mode == libgogo.MODE_REG {
                PrintInstruction_Var_Reg("CMPQ", item1, "R", item2.R);
            }
            if item2.Mode == libgogo.MODE_VAR {
                MakeRegistered(item2, 0);
                PrintInstruction_Var_Reg("CMPQ", item1, "R", item2.R);
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
