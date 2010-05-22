// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// Assignment generation code
// 

package main

import "./libgogo/_obj/libgogo"

func GenerateAssignment(LHSItem *libgogo.Item, RHSItem *libgogo.Item) {
    var done uint64 = 0;
    if Compile != 0 {
        if LHSItem.PtrType == 1 { //Pointer assignment
            if RHSItem.PtrType != 1 {
                SymbolTableError("Cannot assign a value type to a pointer", "", "type:", RHSItem.Itemtype.Name);
            }
            if LHSItem.Itemtype != RHSItem.Itemtype {
                SymbolTableError("Incompatible pointer types:", LHSItem.Itemtype.Name, "and", RHSItem.Itemtype.Name);
            }
            DereferItemIfNecessary(LHSItem); //Derefer address if item is a pointer
            DereferItemIfNecessary(RHSItem); //Derefer address if item is a pointer
        } else { //Value assignment
            if RHSItem.PtrType == 1 {
                SymbolTableError("Cannot assign a pointer type to a value", "", "type:", RHSItem.Itemtype.Name);
            }
            if LHSItem.Itemtype != RHSItem.Itemtype {
                SymbolTableError("Incompatible types:", LHSItem.Itemtype.Name, "and", RHSItem.Itemtype.Name);
            }
            if (LHSItem.Itemtype != byte_t) && (LHSItem.Itemtype != uint64_t) && (LHSItem.Itemtype != string_t) {
                SymbolTableError("Cannot assign to", "", "type", LHSItem.Itemtype.Name);
            }
            //TODO: Type-dependent op code sizes (MOVB, MOVQ, 2xMOVQ etc.)
        }
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
