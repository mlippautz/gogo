// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// Expression related code generation
//

package main

import "./libgogo/_obj/libgogo"

func SwapExpressionBranches(ed *ExpressionDescriptor) {
    var stacksize uint64;
    var depth uint64;
    var tvalue uint64 = 0;
    var fvalue uint64 = 0;

    stacksize = libgogo.GetStackItemCount(&ed.FS);
    if stacksize > 0 {
        depth = libgogo.Peek(&ed.FDepthS);
        if depth >= ed.ExpressionDepth {
            fvalue = libgogo.Pop(&ed.FS);
        }
    }

    stacksize = libgogo.GetStackItemCount(&ed.TS);
    if stacksize > 0 {
        depth = libgogo.Peek(&ed.TDepthS);
        if depth >= ed.ExpressionDepth {
            tvalue = libgogo.Pop(&ed.TS);
        }
    }

    if tvalue != 0 {
        libgogo.Push(&ed.FS, tvalue);
    }
    if fvalue != 0 {
        libgogo.Push(&ed.TS, fvalue);
    }
}

//
// 
//
func SetExpressionDescriptor(ed *ExpressionDescriptor, labelPrefix string) {
    var strLen uint64;
    var singleChar byte;
    var i uint64;
    ed.CurFile = labelPrefix;
    strLen = libgogo.StringLength(fileInfo[curFileIndex].filename);
    for i=0;(i<strLen) && (fileInfo[curFileIndex].filename[i] != '.');i=i+1 {
        singleChar = fileInfo[curFileIndex].filename[i];
        if ((singleChar>=48) && (singleChar<=57)) || ((singleChar>=65) && (singleChar<=90)) || ((singleChar>=97) && (singleChar<=122)) {
            libgogo.CharAppend(&ed.CurFile, fileInfo[curFileIndex].filename[i]);
        } else {
            libgogo.CharAppend(&ed.CurFile, '_');
        }
    }
    ed.ExpressionDepth = 0;
    ed.CurLine = fileInfo[curFileIndex].lineCounter;
    ed.IncCnt = 1;
    ed.T = 0;
    ed.F = 0;
    ed.TDepth = 0;
    ed.FDepth = 0;
    ed.Not = 0;

    libgogo.InitializeStack(&ed.TS);
    libgogo.InitializeStack(&ed.FS);
    libgogo.InitializeStack(&ed.TDepthS);
    libgogo.InitializeStack(&ed.FDepthS);
}

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
    var depth uint64;
    var stacksize uint64;

    if Compile != 0 {
        if item.Mode != libgogo.MODE_COND {
            GenErrorWeak("Can use relative operators only with conditionals.");
        }
        if op == TOKEN_REL_AND {
            labelString = GenerateSubLabel(ed,0 /*negative*/,"END");
            if ed.Not == 0 {
                jmp = GetJump(item.C, 1);
            } else {
                ed.Not = 0;
                stacksize = libgogo.GetStackItemCount(&ed.TS);
                if stacksize > 0 {
                    depth = libgogo.Peek(&ed.TDepthS);
                    if depth >= ed.ExpressionDepth {
                        labelString = GenerateSubLabel(ed,1,"END");
                        jmp = GetJump(item.C, 0);
                        SwapExpressionBranches(ed);
                    }
                } else {
                    jmp = GetJump(item.C,0);
                }
            }
            PrintJump(jmp, labelString);

            stacksize = libgogo.GetStackItemCount(&ed.TS);
            if stacksize > 0 {
                depth = libgogo.Peek(&ed.TDepthS);
                if depth >= ed.ExpressionDepth {
                    PrintLabelWrapped(ed, 1 /*local*/, 1 /*positive*/, "END");
                    libgogo.Pop(&ed.TS);
                    libgogo.Pop(&ed.TDepthS);
                }
            }
        } else {
            if op == TOKEN_REL_OR {
                labelString = GenerateSubLabel(ed,1,"END");
                if ed.Not == 0 {
                    jmp = GetJump(item.C, 0);
                } else {
                    ed.Not = 0;
                    stacksize = libgogo.GetStackItemCount(&ed.FS);
                    if stacksize > 0 {
                        depth = libgogo.Peek(&ed.FDepthS);
                        if depth >= ed.ExpressionDepth {

                            labelString = GenerateSubLabel(ed,0,"END");
                            jmp = GetJump(item.C, 1);
                            SwapExpressionBranches(ed);
                        }
                    } else {
                        jmp = GetJump(item.C,1);
                    }   
                }
                PrintJump(jmp, labelString);
                stacksize = libgogo.GetStackItemCount(&ed.FS);
                if stacksize > 0 {
                    depth = libgogo.Peek(&ed.FDepthS);
                    if depth >= ed.ExpressionDepth {
                        PrintLabelWrapped(ed, 1 /*local*/, 0 /*negative*/, "END");
                        libgogo.Pop(&ed.FS);
                        libgogo.Pop(&ed.FDepthS);
                    }
                }
            } else {
                GenErrorWeak("Relative AND or OR expected.");
            }
        }
        item.C =0;
        FreeRegisterIfRequired(item);
    }
}

//
// Called by parser (ParseExpression)
// Note: This function can only handle operands with a maximum of 8 bytes in size
//
func GenerateComparison(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    var opsize uint64;
    if Compile != 0 {   
        // Type/Pointer checking

        //byte op byte = byte, byte op uint64 = uint64, uint64 op byte = uint64, uint64 op uint64 = uint64
        if (item1.Itemtype == byte_t) && (item2.Itemtype == uint64_t) {
            if item1.Mode != libgogo.MODE_CONST { //No need to convert constants, as their upper bits are already implicitly zeroed
                MakeRegistered(item1, item1.PtrType); //Implicitly convert to uint64 by moving item1 to a register, thereby zeroing the upper bits if necessary
            }
            item1.Itemtype = uint64_t;
        }
        if (item2.Itemtype == byte_t) && (item1.Itemtype == uint64_t) {
            if item2.Mode != libgogo.MODE_CONST { //No need to convert constants, as their upper bits are already implicitly zeroed
                MakeRegistered(item2, item2.PtrType); //Implicitly convert to uint64 by moving item2 to a register, thereby zeroing the upper bits if necessary
            }
            item2.Itemtype = uint64_t;
        }

        if (item1.Itemtype != item2.Itemtype) && (item1.Itemtype != string_t) && (item2.Itemtype != string_t) {
            if (item1.PtrType == 1) && (item2.PtrType == 1) && ((item1.Itemtype != nil) || (item2.Itemtype != nil)) {
                ;
            } else {
                GenErrorWeak("Can only compare variables of same type.");
            }
        }            
        if (item1.Itemtype == string_t) || (item2.Itemtype == string_t) {
            if (item1.PtrType != 1) && (item2.PtrType != 1) {
                GenErrorWeak("Cannot compare string types.");
            }
        }
        if item1.PtrType == 1 {
            if item2.PtrType == 1 {
                if (op != TOKEN_EQUALS) && (op != TOKEN_NOTEQUAL) {
                    GenErrorWeak("Can only compare '==' or '!=' on pointers'");
                }
                if (item1.Mode == libgogo.MODE_CONST) || (item2.Mode == libgogo.MODE_CONST) {
                    GenErrorWeak("Const pointers not allowed. This should not happen.");
                }
            }
        }
        if (item2.PtrType == 1) && (item1.PtrType != 1) {
            GenErrorWeak("Non-pointer to pointer comparison.");
        }
        if (item1.PtrType ==1) && (item2.PtrType != 1) {
            GenErrorWeak("Pointer to non-pointer comparison.");
        }

        // Generate CMP statements depending on items
        if item1.Mode == libgogo.MODE_CONST {
            if item2.Mode == libgogo.MODE_CONST { // Values here, since Ptrs are not allowed
                // Move constvalue to register and compare it against 0                
                item1.Itemtype = bool_t;
                item1.A = GetConditionalBool(op, item1.A, item2.A);
                MakeRegistered(item1, 0);
                // CMP is handled by other if branch (free optimization, yey)
                item2.A = 0;
                op = TOKEN_NOTEQUAL; // Force != for comparison against 0
            } else {
                MakeRegistered(item1, 0);
                if item2.Mode == libgogo.MODE_REG {
                    opsize = GetOpSize(item2, "CMP");
                    PrintInstruction_Reg_Reg("CMP", opsize, "R", item1.R, 0, 0, 0, "", "R", item2.R, 0, 0, 0, "");
                }
                if item2.Mode == libgogo.MODE_VAR {
                    PrintInstruction_Reg_Var("CMP", "R", item1.R, "", 0, item2);
                }
            }
        }
        if item1.Mode == libgogo.MODE_REG {
            DereferRegisterIfNecessary(item1);
            if item2.Mode == libgogo.MODE_CONST {
                PrintInstruction_Reg_Imm("CMP", 8, "R", item1.R, 0, 0, 0, "", item2.A);
            }
            if item2.Mode == libgogo.MODE_REG {
                DereferRegisterIfNecessary(item2);
                opsize = GetOpSize(item2, "CMP");
                PrintInstruction_Reg_Reg("CMP", opsize, "R", item1.R, 0, 0, 0, "", "R", item2.R, 0, 0, 0, "");
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
                DereferRegisterIfNecessary(item2);
                PrintInstruction_Var_Reg("CMP", item1, "R", item2.R, "", 0);
            }
            if item2.Mode == libgogo.MODE_VAR {
                MakeRegistered(item2, 0);
                PrintInstruction_Var_Reg("CMP", item1, "R", item2.R, "", 0);
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

        FreeRegisterIfRequired(item1);
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

