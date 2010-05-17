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

func PrintRegister(name string, number uint64, indirect uint64, offset uint64) {
    var R uint64;
    if indirect != 0 {
        if offset != 0 {
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

func PrintAddress(addr uint64, indirect uint64, offset uint64) {
    if indirect != 0 {
        if offset != 0 {
            PrintOutputValue(offset);
        }
        PrintOutput("(");
    }
    PrintOutputValue(addr);
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

func PrintInstruction_Reg_Reg(op string, reg1name string, reg1number uint64, reg1indirect uint64, reg1offset uint64, reg2name string, reg2number uint64, reg2indirect uint64, reg2offset uint64) {
    PrintInstructionStart(op);
    PrintRegister(reg1name, reg1number, reg1indirect, reg1offset);
    PrintInstructionOperandSeparator();
    PrintRegister(reg2name, reg2number, reg2indirect, reg2offset);
    PrintInstructionEnd();
}

func PrintInstruction_Reg_Mem(op string, regname string, regnumber uint64, regindirect uint64, regoffset uint64, addr uint64, addrindirect uint64, addroffset uint64) {
    PrintInstructionStart(op);
    PrintRegister(regname, regnumber, regindirect, regoffset);
    PrintInstructionOperandSeparator();
    PrintAddress(addr, addrindirect, addroffset);
    PrintInstructionEnd();
}

func PrintInstruction_Mem_Reg(op string, addr uint64, addrindirect uint64, addroffset uint64, regname string, regnumber uint64, regindirect uint64, regoffset uint64) {
    PrintInstructionStart(op);
    PrintAddress(addr, addrindirect, addroffset);
    PrintInstructionOperandSeparator();
    PrintRegister(regname, regnumber, regindirect, regoffset);
    PrintInstructionEnd();
}

func PrintInstruction_Imm_Reg(op string, value uint64, regname string, regnumber uint64, regindirect uint64, regoffset uint64) {
    PrintInstructionStart(op);
    PrintImmediate(value);
    PrintInstructionOperandSeparator();
    PrintRegister(regname, regnumber, regindirect, regoffset);
    PrintInstructionEnd();
}

func PrintInstruction_Imm_Mem(op string, value uint64, addr uint64, addrindirect uint64, addroffset uint64) {
    PrintInstructionStart(op);
    PrintImmediate(value);
    PrintInstructionOperandSeparator();
    PrintAddress(addr, addrindirect, addroffset);
    PrintInstructionEnd();
}
