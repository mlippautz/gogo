// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// General code generation functions (registers, instructions, ...)
// Heavily depends on asm_out.go which represents the Plan9 assembly language
//

package main

import "./libgogo/_obj/libgogo"

// Currently register from R8-R15 are available for usage
var NumRegisters uint64 = 8;
var FreeRegisters [8]byte;

//
// Initialize the register to free-state.
//
func InitFreeRegisters() {
    var i uint64;
    for i = 0; i < NumRegisters; i = i + 1 {
        FreeRegisters[i] = 1;
    }
}

//
// Function returns a free register, BUT is not set to occupied.
//
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

//
// Occupy a given register.
//
func OccupyRegister(index uint64) {
    var realIndex uint64;
    realIndex = index-8;
    FreeRegisters[realIndex] = 0;
}

//
// Free a given register.
//
func FreeRegister(index uint64) {
    var realIndex uint64;
    realIndex = index-8;
    FreeRegisters[realIndex] = 1;
}

//
// Frees the register occupied by the given item if applicable.
// Freeing is only possible if the mode is registered.
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
// Derefers the item given if its type is a pointer
//
func DereferItemIfNecessary(item *libgogo.Item) {
    var oldA uint64;
    if item.PtrType == 1 {
        if item.Mode == libgogo.MODE_REG { //Item is already in a register => derefer register
            PrintInstruction_Reg_Reg("MOVQ", "R", item.R, 1, 0, 0, "R", item.R, 0, 0, 0); //MOVQ (item.R), item.R
        } else { //Item is not a register yet => make it a register by loading its value
            oldA = item.A; //Save value of A
            MakeRegistered(item, 0); //Don't load address as loading the value automatically derefers the item
            item.A = oldA; //Restore old value of A
        }
        item.PtrType = 0; //Item type is no longer a pointer
    }
}

//
// Simple wrapper to asm_out printing
//
func GenerateComment(msg string) {
    var str string = "  //--- ";
    var temp string;
    libgogo.StringAppend(&str, msg);
    libgogo.StringAppend(&str, " at ");
    temp = BuildHead();
    libgogo.StringAppend(&str, temp);
    libgogo.StringAppend(&str, "\n");
    PrintOutput(str);
}


func GenerateFieldAccess(item *libgogo.Item, offset uint64) {
    var offsetItem *libgogo.Item;
    if Compile != 0 {
        DereferItemIfNecessary(item); //Derefer address if item is a pointer
        if item.Mode == libgogo.MODE_VAR { //Variable
            item.A = item.A + offset; //Direct and indirect offset calculation
        } else { //Register
            offsetItem = libgogo.NewItem(); //For direct and indirect offset calculation
            libgogo.SetItem(offsetItem, libgogo.MODE_CONST, uint64_t, 0, offset, 0, 0); //Constant item for offset
            AddSubInstruction("ADDQ", item, offsetItem, 0, 1); //Add constant item (offset), calculating with addresses
        }
    }
}

func GenerateVariableFieldAccess(item *libgogo.Item, offsetItem *libgogo.Item, baseTypeSize uint64) {
    var sizeItem *libgogo.Item;
    if Compile != 0 {
        DereferItemIfNecessary(item); //Derefer address if item is a pointer
        DereferItemIfNecessary(offsetItem); //Derefer address if item is a pointer (should never be necessary here, but just to make sure it is present here)
        sizeItem = libgogo.NewItem();
        libgogo.SetItem(sizeItem, libgogo.MODE_CONST, uint64_t, 0, baseTypeSize, 0, 0); //Constant item
        DivMulInstruction("MULQ", offsetItem, sizeItem, 0, 1); //Multiply identifier value by array base type size => offsetItem now constains the field offset
        AddSubInstruction("ADDQ", item, offsetItem, 0, 1); //Add calculated offset to base address
    }
}

//
// Function converts a given item to registered mode if it is not already 
// a register. This function does not check whether the item is a pointer type!
//
func MakeRegistered(item *libgogo.Item, calculatewithaddresses uint64) {
    var reg uint64;
    if item.Mode != libgogo.MODE_REG {
        reg = GetFreeRegister();
        OccupyRegister(reg);

        if item.Mode == libgogo.MODE_CONST { // const item
            PrintInstruction_Imm_Reg("MOVQ", item.A, "R", reg, 0, 0, 0); // MOVQ $item.A, Rdone (soon to be item.R)
        } else { // var item
            if calculatewithaddresses == 0 {
                PrintInstruction_Var_Reg("MOVQ", item, "R", reg); // MOVQ item.A(SB), Rdone (soon to be item.R)
            } else {
                PrintInstruction_Var_Reg("LEAQ", item, "R", reg); // LEAQ item.A(SB), Rdone (soon to be item.R)
            }
        }

        item.Mode = libgogo.MODE_REG;
        item.R = reg; // item is now a register
        item.A = calculatewithaddresses; // item now contains a value if calculatewithaddresses is 0, or an address if calculatewithaddress is 1
    }
}

//
// Constant folding function. If both items are constants the operation can be
// done in the compiler.
//
func ConstFolding(item1 *libgogo.Item, item2 *libgogo.Item, constvalue uint64) uint64 {
    var boolFlag uint64 = 0;
    if (item1.Mode == libgogo.MODE_CONST) && (item2.Mode == libgogo.MODE_CONST) {
        item1.A = constvalue;
        boolFlag = 1;
    }
    return boolFlag;
}

//
// item1 = item1 OP item2, or constvalue if both item1 and item2 are constants
// Side effect: The register item2 occupies is freed if applicable
// If calculatewithaddresses is 0, it is assumed that registers contain values, 
// otherwise it is assumed that they contain addresses
//
func AddSubInstruction(op string, item1 *libgogo.Item, item2 *libgogo.Item, constvalue uint64, calculatewithaddresses uint64) {
    var done uint64 = 0;

    done = ConstFolding(item1, item2, constvalue);
    
    if done == 0 {
        DereferItemIfNecessary(item1); //Derefer address if item is a pointer
        DereferItemIfNecessary(item2); //Derefer address if item is a pointer
    }

    if (done == 0) && (item1.Mode != libgogo.MODE_REG) { //item1 is not a register => make it a register
        MakeRegistered(item1, calculatewithaddresses);
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
            if calculatewithaddresses == 0 { // Calculate with values
                DereferRegisterIfNecessary(item2);
            }
            PrintInstruction_Reg_Reg(op, "R", item2.R, 0, 0, 0, "R", item1.R, 0, 0, 0); //OP item2.R, item1.R
            done = 1;
        }
    }

    FreeRegisterIfRequired(item2); // item2 should be useless by now
}

//
// item1 = item1 OP item2, or constvalue if both item1 and item2 are constants
// Difference here is that it uses a one operand assembly instruction which 
// operates on AX as first operand
//
func DivMulInstruction(op string, item1 *libgogo.Item, item2 *libgogo.Item, constvalue uint64, calculatewithaddresses uint64) {
    var done uint64 = 0;

    done = ConstFolding(item1, item2, constvalue);
    
    if done == 0 {
        DereferItemIfNecessary(item1); //Derefer address if item is a pointer
        DereferItemIfNecessary(item2); //Derefer address if item is a pointer
    }

    if done == 0 { // item1 is now (or has even already been) a register => use it
        if calculatewithaddresses == 0 { // Calculate with values
            DereferRegisterIfNecessary(item1); // Calculate with values
        }

        if item1.Mode == libgogo.MODE_CONST {
            PrintInstruction_Imm_Reg("MOVQ", item1.A, "AX", 0, 0, 0, 0) // move $item1.A into AX
        }
        if item1.Mode == libgogo.MODE_VAR {
            PrintInstruction_Var_Reg("MOVQ", item1, "AX", 0); // move item2.A(SB), AX
        }
        if item1.Mode == libgogo.MODE_REG {
            PrintInstruction_Reg_Reg("MOVQ", "R", item1.R, 0, 0, 0, "AX", 0, 0, 0, 0) // move item1.R into AX
        }

        if item2.Mode != libgogo.MODE_REG {
            // item2 needs to be registered as the second operand of a DIV/MUL
            // instruction always needs to be a register
            MakeRegistered(item2, calculatewithaddresses);
        }

        // OP item2.R
        if calculatewithaddresses == 0 { // Calculate with values
            DereferRegisterIfNecessary(item2);
        }
        done = libgogo.StringCompare(op,"DIVQ");
        if done == 0 { //Set DX to zero to avoid 128 bit division as DX is "high" part of DX:AX 128 bit register
            PrintInstruction_Reg_Reg("XORQ", "DX", 0, 0, 0, 0, "DX", 0, 0, 0, 0); //XORQ DX, DX is equal to MOVQ $0, DX
        }
        PrintInstruction_Reg(op, "R", item2.R, 0, 0, 0); //op item2.R
        PrintInstruction_Reg_Reg("MOVQ", "AX", 0, 0, 0, 0, "R", item2.R, 0, 0, 0) // move AX into item2.R
    }

    // Since item2 already had to be converted to a register, we now assign 
    // item2 to item1 after freeing item1 first (if necessary)
    FreeRegisterIfRequired(item1);
    item1.Mode = item2.Mode;
    item1.R = item2.R;
    item1.A = item2.A;
    item1.Itemtype = item2.Itemtype;
    item1.PtrType = item2.PtrType; //Should always be 0
    item1.Global = item2.Global;
}
