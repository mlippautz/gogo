// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

func PrintOutput(output string) {
    libgogo.PrintString(output); //TODO: Output to file
}

func PrintOutputValue(value uint64) {
    var temp string;
    temp = libgogo.IntToString(value);
    PrintOutput(temp);
}

func PrintRegister(name string, number uint64, indirect uint64, offset uint64, negativeoffset uint64) {
    var R uint64;
    if indirect != 0 {
        if offset != 0 {
            if negativeoffset != 0 {
                PrintOutput("-");
            }
            PrintOutputValue(offset);
        }
        PrintOutput("(");
    }
    PrintOutput(name); //Print register name
    R = libgogo.StringCompare(name, "R");
    if R == 0 { //p.e. R8 => consider number; else: p.e. AX => ignore number
        PrintOutputValue(number);
    }
    if indirect != 0 {
        PrintOutput(")");
    }
}

func PrintImmediate(value uint64) {
    PrintOutput("$");
    PrintOutputValue(value);
}

func PrintInstructionStart(op string) {
    PrintOutput(op);
    PrintOutput(" ");
}

func PrintInstructionOperandSeparator() {
    PrintOutput(", ");
}

func PrintInstructionEnd() {
    PrintOutput("\n");
}

func PrintInstruction_Reg_Reg(op string, reg1name string, reg1number uint64, reg1indirect uint64, reg1offset uint64, reg1offsetnegative uint64, reg2name string, reg2number uint64, reg2indirect uint64, reg2offset uint64, reg2offsetnegative uint64) {
    PrintInstructionStart(op);
    PrintRegister(reg1name, reg1number, reg1indirect, reg1offset, reg1offsetnegative);
    PrintInstructionOperandSeparator();
    PrintRegister(reg2name, reg2number, reg2indirect, reg2offset, reg2offsetnegative);
    PrintInstructionEnd();
}

func PrintInstruction_Imm_Reg(op string, value uint64, regname string, regnumber uint64, regindirect uint64, regoffset uint64, regoffsetnegative uint64) {
    PrintInstructionStart(op);
    PrintImmediate(value);
    PrintInstructionOperandSeparator();
    PrintRegister(regname, regnumber, regindirect, regoffset, regoffsetnegative);
    PrintInstructionEnd();
}

func PrintInstruction_Imm_Var(op string, value uint64, variable *libgogo.Item) {
    if variable.Global == 1 { //Global
        PrintInstruction_Imm_Reg(op, value, "SB", 0, 1, variable.A, 0); //OP $value, variable.A(SB)
    } else { //Local
        PrintInstruction_Imm_Reg(op, value, "SP", 0, 1, variable.A + 8, 1); //OP $value, -[variable.A+8](SP)
    }
}

func PrintInstruction_Reg_Var(op string, regname string, regnumber uint64, variable *libgogo.Item) {
    if variable.Global == 1 { //Global
        PrintInstruction_Reg_Reg(op, regname, regnumber, 0, 0, 0, "SB", 0, 1, variable.A, 0); //OP regname_regnumber, variable.A(SB)
    } else { //Local
        PrintInstruction_Reg_Reg(op, regname, regnumber, 0, 0, 0, "SP", 0, 1, variable.A + 8, 1); //OP regname_regnumber, -[variable.A+8](SP)
    }
}

func PrintInstruction_Var_Reg(op string, variable *libgogo.Item, regname string, regnumber uint64) {
    if variable.Global == 1 { //Global
        PrintInstruction_Reg_Reg(op, "SB", 0, 1, variable.A, 0, regname, regnumber, 0, 0, 0); //OP variable.A(SB), regname_regnumber
    } else { //Local
        PrintInstruction_Reg_Reg(op, "SP", 0, 1, variable.A + 8, 1, regname, regnumber, 0, 0, 0); //OP -[variable.A+8](SP), regname_regnumber
    }
}
