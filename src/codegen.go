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

/*func ItemToRegister(item *libgogo.Item) {
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
}*/

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

//-----------------------------------------------------------------------

//
// Frees the register occupied by the given item if applicable
//
func FreeRegisterIfRequired(item *libgogo.Item) {
    if item.Mode == libgogo.MODE_REG {
        FreeRegister(item.R);
    }
}

//
// Moves the value of the address a register is currently pointing to into the register itself
//
func DereferRegisterIfNecessary(item *libgogo.Item) {
    if (item.Mode == libgogo.MODE_REG) && (item.A != 0) { //Derefer register if it contains an address
        PrintInstruction_Reg_Reg("MOVQ", "R", item.R, 1, 0, 0, "R", item.R, 0, 0, 0); //MOVQ (item.R), item.R
        item.A = 0; //Register now contains a value
    }
}

func GenerateSimpleExpression(item1 *libgogo.Item, item2 *libgogo.Item, op uint64) {
    if Compile != 0 {
        if op == TOKEN_ARITH_PLUS { //Add
            TwoOperandInstruction("ADDQ", item1, item2, item1.A + item2.A, 0);
        } else { //Subtract
            TwoOperandInstruction("SUBQ", item1, item2, item1.A - item2.A, 0);
        }
    }
}

func GenerateAssignment(LHSItem *libgogo.Item, RHSItem *libgogo.Item) {
    var done uint64 = 0;
    if Compile != 0 {
        if LHSItem.Mode == libgogo.MODE_VAR { //Variable on LHS
            if (done == 0) && (RHSItem.Mode == libgogo.MODE_CONST) { //Const RHS
                PrintInstruction_Imm_Var("MOVQ", RHSItem.A, LHSItem); //MOVQ $RHSItem.A, LHSItem.A(SB)
                done = 1;
            }
            if (done == 0) && (RHSItem.Mode == libgogo.MODE_VAR) { //Var RHS
                done = GetFreeRegister();
                OccupyRegister(done);
                PrintInstruction_Var_Reg("MOVQ", RHSItem, "R", done); //MOVQ RHSItem.A(SB), Rdone (soon to be RHSItem.R)
                RHSItem.Mode = libgogo.MODE_REG;
                RHSItem.R = done; //RHS is now a register
                RHSItem.A = 0; //Register now contains RHS value
                PrintInstruction_Reg_Var("MOVQ", "R", RHSItem.R, LHSItem); //MOVQ RHSItem.R, LHSItem.A(SB)
                done = 1;
            }
            if (done == 0) && (RHSItem.Mode == libgogo.MODE_REG) { //Reg RHS
                DereferRegisterIfNecessary(RHSItem); //Make sure to work with the value, not the address
                PrintInstruction_Reg_Var("MOVQ", "R", RHSItem.R, LHSItem); //MOVQ RHSItem.R, LHSItem.A(SB)
                done = 1;
            }
        } else { //Register with address of variable on LGS; assertion: Register contains address and global/local flag is set correctly
            if (done == 0) && (RHSItem.Mode == libgogo.MODE_CONST) { //Const RHS
                PrintInstruction_Imm_Reg("MOVQ", RHSItem.A, "R", LHSItem.R, 1, 0, 0); //MOVQ $RHSItem.A, (LHSItem.R)
                done = 1;
            }
            if (done == 0) && (RHSItem.Mode == libgogo.MODE_VAR) { //Var RHS
                done = GetFreeRegister();
                OccupyRegister(done);
                PrintInstruction_Var_Reg("MOVQ", RHSItem, "R", done); //MOVQ RHSItem.A(SB), Rdone (soon to be RHSItem.R)
                RHSItem.Mode = libgogo.MODE_REG;
                RHSItem.R = done; //RHS is now a register
                RHSItem.A = 0; //Register now contains RHS value
                PrintInstruction_Reg_Reg("MOVQ", "R", RHSItem.R, 0, 0, 0, "R", LHSItem.R, 1, 0, 0); //MOVQ RHSItem.R, (LHSItem.R)
                done = 1;
            }
            if (done == 0) && (RHSItem.Mode == libgogo.MODE_REG) { //Reg RHS
                DereferRegisterIfNecessary(RHSItem); //Make sure to work with the value, not the address
                PrintInstruction_Reg_Reg("MOVQ", "R", RHSItem.R, 0, 0, 0, "R", LHSItem.R, 1, 0, 0); //MOVQ RHSItem.R, (LHSItem.R)
                done = 1;
            }
        }        
        FreeRegisterIfRequired(LHSItem);
        FreeRegisterIfRequired(RHSItem);
    }
}

func GenerateFieldAccess(item *libgogo.Item, offset uint64, indirect uint64) {
    var offsetItem *libgogo.Item;
    var temp uint64;
    if Compile != 0 {
        if item.Mode == libgogo.MODE_VAR { //Variable
            item.A = item.A + offset; //Direct and indirect offset calculation
            if indirect != 0 { //Indirect
                temp = GetFreeRegister();
                OccupyRegister(temp);
                PrintInstruction_Var_Reg("LEAQ", item, "R", temp); //LEAQ item.A(SB), Rtemp (soon to be item.R)
                item.Mode = libgogo.MODE_REG;
                item.R = temp;
                item.A = 1; //Register contains address
                DereferRegisterIfNecessary(item); //Indirection
                item.A = 1; //Register still contains address
            }
        } else { //Register
            offsetItem = libgogo.NewItem(); //For direct and indirect offset calculation
            libgogo.SetItem(offsetItem, libgogo.MODE_CONST, uint64_t, offset, 0, 0); //Constant item for offset
            TwoOperandInstruction("ADDQ", item, offsetItem, 0, 1); //Add constant item (offset), calculating with addresses
            if indirect != 0 { //Indirect
                DereferRegisterIfNecessary(item); //Indirection
                item.A = 1; //Register still contains address
            }
        }
    }
}

//
// item1 = item1 OP item2, or constvalue if both item1 and item2 are constants
// Side effect: The register item2 occupies is freed if applicable
// If calculatewithaddresses is 0, it is assumed that registers contain values, otherwise it is assumed that they contain addresses
//
func TwoOperandInstruction(op string, item1 *libgogo.Item, item2 *libgogo.Item, constvalue uint64, calculatewithaddresses uint64) {
    var done uint64 = 0;
    if (done == 0) && (item1.Mode == libgogo.MODE_CONST) && (item2.Mode == libgogo.MODE_CONST) { //Constant folding
        item1.A = constvalue; //item1 = item1 OP item2 (constvalue)
        done = 1;
    }
    if (done == 0) && (item1.Mode != libgogo.MODE_REG) { //item1 is not a register => make it a register
        done = GetFreeRegister();
        OccupyRegister(done);
        if item1.Mode == libgogo.MODE_CONST { //item1 is const
            PrintInstruction_Imm_Reg("MOVQ", item1.A, "R", done, 0, 0, 0); //MOVQ $item1.A, Rdone (soon to be item1.R)
        } else { //item1 is var
            if calculatewithaddresses == 0 {
                PrintInstruction_Var_Reg("MOVQ", item1, "R", done); //MOVQ item1.A(SB), Rdone (soon to be item1.R)
            } else {
                PrintInstruction_Var_Reg("LEAQ", item1, "R", done); //LEAQ item1.A(SB), Rdone (soon to be item1.R)
            }
        }
        item1.Mode = libgogo.MODE_REG;
        item1.R = done; //item1 is now a register; don't set done to 1 as the actual calculation has yet to be done
        item1.A = calculatewithaddresses; //item1 now contains a value if calculatewithaddresses is 0, or an address if calculatewithaddress is 1
    }
    if done == 0 { //item1 is now (or has even already been) a register => use it
		    if calculatewithaddresses == 0 { //Calculate with values
            DereferRegisterIfNecessary(item1); //Calculate with values
        }
				if (done == 0) && (item2.Mode == libgogo.MODE_CONST) {
				    PrintInstruction_Imm_Reg(op, item2.A, "R", item1.R, 0, 0, 0); //OP $item2.A, item1.R
				    done = 1;
				}
				if (done == 0) && (item2.Mode == libgogo.MODE_VAR) {
				    PrintInstruction_Var_Reg(op, item2, "R", item1.R); //OP item2.A(SB), item1.R
				    done = 1;
				}
				if (done == 0) && (item2.Mode == libgogo.MODE_REG) {
				    if calculatewithaddresses == 0 { //Calculate with values
				        DereferRegisterIfNecessary(item2);
				    }
				    PrintInstruction_Reg_Reg("ADDQ", "R", item2.R, 0, 0, 0, "R", item1.R, 0, 0, 0); //OP item2.R, item1.R
				    done = 1;
				}
    }
    FreeRegisterIfRequired(item2);
}
