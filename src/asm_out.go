// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

//
// Variable holding the compiled code
//
var Code string;

//
// Reseting the code before a new compile round starts
//
func ResetCode() {
    Code = "\
//\n\
// --------------------\n\
// GoGo compiler output\n\
// --------------------\n\
//\n\
// File: \
";
    libgogo.StringAppend(&Code, fileInfo[curFileIndex].filename);
    libgogo.StringAppend(&Code,"\n\
// Syntax: Plan-9 assembler\n\
//\n\
// This code is automatically generated. DO NOT EDIT IT!\n\
//\n\
\n\
");
    InspectorGadget();

    libgogo.StringAppend(&Code,"\n\
\n\
TEXT    main·init(SB),0,$0-0\n\
  // unused\n\
  RET\n\
");

    libgogo.StringAppend(&Code,"\n\
TEXT    main·main(SB),0,$0-24\n\
");
}

//
// Function printing the generated code (stored in Code) to a file called the
// same name as the input file + '.sog' extension
//
func PrintFile() {
    var fd uint64;    
    var outfile string = "_gogo_.sog";
    // The following line creates a new file for the assembler code
    // flags: O_WRONLY | O_CREAT | O_TRUNC => 577
    // mode: S_IWUSR | S_IRUSR | S_IRGRP => 416
    fd = libgogo.FileOpen2(outfile,577,416);

    libgogo.StringAppend(&Code,"\n\
  RET\n\
");

    libgogo.WriteString(fd,Code);
    libgogo.FileClose(fd);
}

func PrintOutput(output string) {
    libgogo.StringAppend(&Code,output);
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
    PrintOutput("  ");
    PrintOutput(op);
    PrintOutput(" ");
}

func PrintInstructionOperandSeparator() {
    PrintOutput(", ");
}

func PrintInstructionEnd() {
    PrintOutput("\n");
}

func PrintInstruction_Reg(op string, name string, number uint64, indirect uint64, offset uint64, offsetnegative uint64) {
    PrintInstructionStart(op);
    PrintRegister(name, number, indirect, offset, offsetnegative);
    PrintInstructionEnd();
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
        if variable.Global == 2 { //Parameter
            PrintInstruction_Imm_Reg(op, value, "SP", 0, 1, variable.A + 8, 0); //OP $value, [variable.A+8](SP)
        } else { //Local
            PrintInstruction_Imm_Reg(op, value, "SP", 0, 1, variable.A + 8, 1); //OP $value, -[variable.A+8](SP)
        }
    }
}

func PrintInstruction_Reg_Var(op string, regname string, regnumber uint64, variable *libgogo.Item) {
    if variable.Global == 1 { //Global
        PrintInstruction_Reg_Reg(op, regname, regnumber, 0, 0, 0, "SB", 0, 1, variable.A, 0); //OP regname_regnumber, variable.A(SB)
    } else { //Local
        if variable.Global == 2 { //Parameter
            PrintInstruction_Reg_Reg(op, regname, regnumber, 0, 0, 0, "SP", 0, 1, variable.A + 8, 0); //OP regname_regnumber, [variable.A+8](SP)
        } else { //Local
            PrintInstruction_Reg_Reg(op, regname, regnumber, 0, 0, 0, "SP", 0, 1, variable.A + 8, 1); //OP regname_regnumber, -[variable.A+8](SP)
        }
    }
}

func PrintInstruction_Var_Reg(op string, variable *libgogo.Item, regname string, regnumber uint64) {
    if variable.Global == 1 { //Global
        PrintInstruction_Reg_Reg(op, "SB", 0, 1, variable.A, 0, regname, regnumber, 0, 0, 0); //OP variable.A(SB), regname_regnumber
    } else { //Local
        if variable.Global == 2 { //Parameter
            PrintInstruction_Reg_Reg(op, "SP", 0, 1, variable.A + 8, 0, regname, regnumber, 0, 0, 0); //OP [variable.A+8](SP), regname_regnumber
        } else { //Local
            PrintInstruction_Reg_Reg(op, "SP", 0, 1, variable.A + 8, 1, regname, regnumber, 0, 0, 0); //OP -[variable.A+8](SP), regname_regnumber
        }
    }
}
