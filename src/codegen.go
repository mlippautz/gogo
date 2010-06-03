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
// Saves the currently used registers to the stack
//
func SaveUsedRegisters() {
    var i uint64;
    for i = 0; i < NumRegisters; i = i + 1 {
        if FreeRegisters[i] == 0 {
            PrintInstruction_Reg("PUSH", 8, "R", i + 8, 0, 0, 0, "");
        }
    }
}

//
// Restores the currently used registers from the stack
//
func RestoreUsedRegisters() {
    var i uint64;
    for i = NumRegisters - 1; i >= 0; i = i - 1 { //Reverse order (stack!)
        if FreeRegisters[i] == 0 {
            PrintInstruction_Reg("POP", 8, "R", i + 8, 0, 0, 0, "");
        }
    }
}

//
// Frees the register occupied by the given item if applicable.
// Freeing is only possible if the mode is registered.
//
func FreeRegisterIfRequired(item *libgogo.Item) {
    if item.Mode == libgogo.MODE_REG {
        FreeRegister(item.R);
        if item.C != 0 {
            FreeRegister(item.C);
        }
    }
}

//
// Moves the value of the address a register is currently pointing to into the register itself
//
func DereferRegisterIfNecessary(item *libgogo.Item) {
    var opsize uint64;
    var retVal uint64;
    if (item.Mode == libgogo.MODE_REG) && (item.A != 0) { //Derefer register if it contains an address
        item.A = 0; //Register will soon contain a value; make sure op size calculation is based on the actual type size, not on the pointer size
        opsize = GetOpSize(item, "MOV");
        retVal = PrintInstruction_Reg_Reg("MOV", opsize, "R", item.R, 1, 0, 0, "", "R", item.R, 0, 0, 0, ""); //MOV (item.R), item.R
        if retVal != 0 {
            retVal = GetFreeRegister(); //Additional register required
            OccupyRegister(retVal);
            item.C = retVal;
            PrintInstruction_Reg_Reg("MOV", retVal, "R", item.R, 1, 8, 0, "", "R", item.C, 0, 0, 0, ""); //MOV 8(item.R), item.C
        }
    }
}

//
// Derefers the item given if its type is a pointer
//
func DereferItemIfNecessary(item *libgogo.Item) {
    var oldA uint64;
    var opsize uint64;
    var retVal uint64;
    if item.PtrType == 1 {
        if item.Mode == libgogo.MODE_REG { //Item is already in a register => derefer register
            opsize = GetOpSize(item, "MOV");
            retVal = PrintInstruction_Reg_Reg("MOV", opsize, "R", item.R, 1, 0, 0, "", "R", item.R, 0, 0, 0, ""); //MOV (item.R), item.R
            if retVal != 0 {
                retVal = GetFreeRegister(); //Additional register required
                OccupyRegister(retVal);
                item.C = retVal;
                PrintInstruction_Reg_Reg("MOV", retVal, "R", item.R, 1, 8, 0, "", "R", item.C, 0, 0, 0, ""); //MOV 8(item.R), item.C
            }
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
    var str string = "";
    var tmpPtr *string;
    var temp string;
    if DEBUG_LEVEL >= 10 {
        tmpPtr = &DataSegment;
        if (OutputStringPtr != tmpPtr) {
            str = "  //--- ";
        } else { //No indentation in data segment
            str = "//--- ";
        }
        libgogo.StringAppend(&str, msg);
        libgogo.StringAppend(&str, " at ");
        temp = BuildHead();
        libgogo.StringAppend(&str, temp);
        libgogo.StringAppend(&str, "\n");
        PrintCodeOutput(str);
    }
}


func GenerateFieldAccess(item *libgogo.Item, offset uint64) {
    var offsetItem *libgogo.Item;
    if Compile != 0 {
        DereferItemIfNecessary(item); //Derefer address if item is a pointer
        if item.Mode == libgogo.MODE_VAR { //Variable
            if item.Global == 0 { //Local variable offset calculation
                item.A = item.A - offset; //Reverse order due to sign (p.e. -24(SP) with offset 16 is to be -8(SP) and thus 24-16)
            } else { //Global variable offset calculation
                item.A = item.A + offset;
            }
        } else { //Register
            offsetItem = libgogo.NewItem(); //For direct and indirect offset calculation
            libgogo.SetItem(offsetItem, libgogo.MODE_CONST, uint64_t, 0, offset, 0, 0); //Constant item for offset
            AddSubInstruction("ADD", item, offsetItem, 0, 1); //Add constant item (offset), calculating with addresses
        }
    }
}

func GenerateVariableFieldAccess(item *libgogo.Item, offsetItem *libgogo.Item, baseTypeSize uint64) {
    var sizeItem *libgogo.Item;
    if Compile != 0 {
        DereferItemIfNecessary(item); //Derefer address if item is a pointer
        if (offsetItem.Itemtype != byte_t) && (offsetItem.Itemtype != uint64_t) {
            SymbolTableError("Invalid index type for", "", "array access:", offsetItem.Itemtype.Name);
        }
        sizeItem = libgogo.NewItem();
        libgogo.SetItem(sizeItem, libgogo.MODE_CONST, uint64_t, 0, baseTypeSize, 0, 0); //Constant item
        DivMulInstruction("MUL", offsetItem, sizeItem, 0, 1); //Multiply identifier value by array base type size => offsetItem now constains the field offset
        AddSubInstruction("ADD", item, offsetItem, 0, 1); //Add calculated offset to base address
    }
}

//
// Function converts a given item to registered mode if it is not already 
// a register. This function does not check whether the item is a pointer type!
//
func MakeRegistered(item *libgogo.Item, calculatewithaddresses uint64) {
    var reg uint64;
    var reg2 uint64;
    var opsize uint64;
    if item.Mode != libgogo.MODE_REG {
        reg = GetFreeRegister();
        OccupyRegister(reg);

        if item.Mode == libgogo.MODE_CONST { // const item
            opsize = GetOpSize(item, "MOV");
            PrintInstruction_Imm_Reg("MOV", opsize, item.A, "R", reg, 0, 0, 0, ""); // MOV $item.A, Rdone (soon to be item.R)
        } else { // var item
            if calculatewithaddresses == 0 {
                opsize = GetOpSize(item, "MOV");
                if opsize > 8 { //Occupy additional register if operand size is higher than 8 bytes
                    reg2 = GetFreeRegister();
                    OccupyRegister(reg2);
                    item.C = reg2;
                }
                PrintInstruction_Var_Reg("MOV", item, "R", reg, "R", reg2); // MOV item.A(SB), Rdone (soon to be item.R)
            } else {
                opsize = GetOpSize(item, "LEA");
                PrintInstruction_Var_Reg("LEA", item, "R", reg, "", 0); // LEA item.A(SB), Rdone (soon to be item.R)
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
