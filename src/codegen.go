// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// General code generation functions (registers, instructions, ...)
//

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
    return i+8;
}

func OccupyRegister(index uint64) {
    var realIndex uint64;
    realIndex = index-8;
    FreeRegisters[realIndex] = 0;
}

func FreeRegister(index uint64) {
    var realIndex uint64;
    realIndex = index-8;
    FreeRegisters[realIndex] = 1;
}

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

//
// Simple wrapper to asm_out printing
//
func GenerateComment(msg string) {
    var str string = "  // >>> ";
    libgogo.StringAppend(&str, msg);
    libgogo.StringAppend(&str,"\n");
    PrintOutput(str);
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
        done = 0;
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
            PrintInstruction_Reg_Reg(op, "R", item2.R, 0, 0, 0, "R", item1.R, 0, 0, 0); //OP item2.R, item1.R
            done = 1;
        }
    }
    FreeRegisterIfRequired(item2);
}

