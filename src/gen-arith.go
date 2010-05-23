// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// General code generation functions (registers, instructions, ...)
// Heavily depends on asm_out.go which represents the Plan9 assembly language
//

package main

import "./libgogo/_obj/libgogo"

//
// item1 = item1 OP item2, or constvalue if both item1 and item2 are constants
// Side effect: The register item2 occupies is freed if applicable
// If calculatewithaddresses is 0, it is assumed that registers contain values, 
// otherwise it is assumed that they contain addresses
// Note: This function does not perform type checking, but converts one item to
// type uint64 if necessary
//
func AddSubInstruction(op string, item1 *libgogo.Item, item2 *libgogo.Item, constvalue uint64, calculatewithaddresses uint64) {
    var done uint64 = 0;
    var opsize1 uint64;
    var opsize2 uint64;

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

        //byte + byte = byte, byte + uint64 = uint64, uint64 + byte = uint64, uint64 + uint64 = uint64
        if (item1.Itemtype == byte_t) && (item2.Itemtype == uint64_t) {
            item1.Itemtype = uint64_t;
        }
        if (item2.Itemtype == byte_t) && (item1.Itemtype == uint64_t) {
            item2.Itemtype = uint64_t;
        } //TODO: Perform real conversion (setting the unused bits to zero etc.)
        opsize1 = GetOpSize(item1);
        opsize2 = GetOpSize(item2);
        if opsize1 > opsize2 {
            opsize2 = opsize1;
        } else {
            opsize1 = opsize2;
        }
        
        if (done == 0) && (item2.Mode == libgogo.MODE_CONST) {
            PrintInstruction_Imm_Reg(op, opsize2, item2.A, "R", item1.R, 0, 0, 0); //OP $item2.A, item1.R
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
            opsize2 = GetOpSize(item2);
            if opsize1 > opsize2 {
                opsize2 = opsize1;
            } else {
                opsize1 = opsize2;
            }
            PrintInstruction_Reg_Reg(op, opsize2, "R", item2.R, 0, 0, 0, "R", item1.R, 0, 0, 0); //OP item2.R, item1.R
            done = 1;
        }
    }

    FreeRegisterIfRequired(item2); // item2 should be useless by now
}

//
// item1 = item1 OP item2, or constvalue if both item1 and item2 are constants
// Difference here is that it uses a one operand assembly instruction which 
// operates on AX as first operand
// Note: This function does not perform type checking
//
func DivMulInstruction(op string, item1 *libgogo.Item, item2 *libgogo.Item, constvalue uint64, calculatewithaddresses uint64) {
    var done uint64 = 0;
    var opsize1 uint64;
    var opsize2 uint64;

    done = ConstFolding(item1, item2, constvalue);
    
    if done == 0 {
        DereferItemIfNecessary(item1); //Derefer address if item is a pointer
        DereferItemIfNecessary(item2); //Derefer address if item is a pointer
    }

    if done == 0 { // item1 is now (or has even already been) a register => use it
        if calculatewithaddresses == 0 { // Calculate with values
            DereferRegisterIfNecessary(item1); // Calculate with values
        }

        opsize1 = GetOpSize(item1);
        if item1.Mode == libgogo.MODE_CONST {
            PrintInstruction_Imm_Reg("MOV", opsize1, item1.A, "AX", 0, 0, 0, 0) // move $item1.A into AX
        }
        if item1.Mode == libgogo.MODE_VAR {
            PrintInstruction_Var_Reg("MOV", item1, "AX", 0); // move item2.A(SB), AX
        }
        if item1.Mode == libgogo.MODE_REG {
            PrintInstruction_Reg_Reg("MOV", opsize1, "R", item1.R, 0, 0, 0, "AX", 0, 0, 0, 0) // move item1.R into AX
        }
        
        //byte * byte = byte, byte * uint64 = uint64, uint64 * byte = uint64, uint64 * uint64 = uint64
        if (item1.Itemtype == byte_t) && (item2.Itemtype == uint64_t) {
            item1.Itemtype = uint64_t;
        }
        if (item2.Itemtype == byte_t) && (item1.Itemtype == uint64_t) {
            item2.Itemtype = uint64_t;
        } //TODO: Perform real conversion (setting the unused bits to zero etc.)

        if item2.Mode != libgogo.MODE_REG {
            // item2 needs to be registered as the second operand of a DIV/MUL
            // instruction always needs to be a register
            MakeRegistered(item2, calculatewithaddresses);
        }

        // OP item2.R
        if calculatewithaddresses == 0 { // Calculate with values
            DereferRegisterIfNecessary(item2);
        }
        done = libgogo.StringCompare(op, "DIV");
        if done == 0 { //Set DX to zero to avoid 128 bit division as DX is "high" part of DX:AX 128 bit register
            PrintInstruction_Reg_Reg("XOR", 8, "DX", 0, 0, 0, 0, "DX", 0, 0, 0, 0); //XORQ DX, DX is equal to MOVQ $0, DX
        }
        
        opsize1 = GetOpSize(item1);
        opsize2 = GetOpSize(item2);
        if opsize1 > opsize2 {
            opsize2 = opsize1;
        } else {
            opsize1 = opsize2;
        }
        
        PrintInstruction_Reg(op, opsize2, "R", item2.R, 0, 0, 0); //op item2.R
        PrintInstruction_Reg_Reg("MOV", opsize2, "AX", 0, 0, 0, 0, "R", item2.R, 0, 0, 0) // move AX into item2.R
        
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
}
