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
";
    libgogo.StringAppend(&Code, fileInfo[curFileIndex].filename);
    libgogo.StringAppend(&Code,"\n\
// Syntax: Plan-9 assembler\n\
//\n\
// This code is automatically generated. DO NOT EDIT IT!\n\
//\n\
\n\
");
    //InspectorGadget();

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

func GetOpSize(item *libgogo.Item) uint64 {
    var size uint64;
    if (item.PtrType == 1) || ((item.Mode == libgogo.MODE_REG) && (item.A != 0)) { //Pointer type or register with address
        size = 8; //Pointers always have a size of 64 bits
    } else { //Value type
        if item.Itemtype.Form == libgogo.FORM_SIMPLE { //Simple type
            size = item.Itemtype.Len; //Size of type
        } else {
            if item.Itemtype == string_t { //Special treatment for strings => 16 bytes
                size = 16;
            } else {
                size = 8; //Use 64 bits in all other cases (address calculations, records etc.)
            }
        }
    }
    return size;
}

func PrintInstructionStart(op string, opsize uint64) {
    PrintOutput("  ");
    PrintOutput(op);
    if opsize == 1 { //byte => B
        PrintOutput("B");
    } else {
        if opsize == 8 { //uint64 => Q
            PrintOutput("Q");
        } else {
            PrintOutput("?"); //TODO: Error: invalid opsize
        }
    }
    PrintOutput(" ");
}

func PrintInstructionOperandSeparator() {
    PrintOutput(", ");
}

func PrintInstructionEnd() {
    PrintOutput("\n");
}

func PrintInstruction_Reg(op string, opsize uint64, name string, number uint64, indirect uint64, offset uint64, offsetnegative uint64) {
    PrintInstructionStart(op, opsize);
    PrintRegister(name, number, indirect, offset, offsetnegative);
    PrintInstructionEnd();
}

func PrintInstruction_Reg_Reg(op string, opsize uint64, reg1name string, reg1number uint64, reg1indirect uint64, reg1offset uint64, reg1offsetnegative uint64, reg2name string, reg2number uint64, reg2indirect uint64, reg2offset uint64, reg2offsetnegative uint64) {
    if (opsize == 1) && (reg2indirect == 0) { //Clear upper bits using AND mask when operating on bytes; don't clear memory as op could be MOV and therefore set 7 bytes unrecoverably to zero
        PrintInstruction_Imm_Reg("AND", 8, 255, reg2name, reg2number, reg2indirect, reg2offset, reg2offsetnegative); //ANDQ $255, R
    }
    PrintInstructionStart(op, opsize);
    PrintRegister(reg1name, reg1number, reg1indirect, reg1offset, reg1offsetnegative);
    PrintInstructionOperandSeparator();
    PrintRegister(reg2name, reg2number, reg2indirect, reg2offset, reg2offsetnegative);
    PrintInstructionEnd();
}

func PrintInstruction_Reg_Imm(op string, opsize uint64, regname string, regnumber uint64, regindirect uint64, regoffset uint64, regoffsetnegative uint64, value uint64) {
    //No opcode check necessary here
    PrintInstructionStart(op, opsize);
    PrintRegister(regname, regnumber, regindirect, regoffset, regoffsetnegative);
    PrintInstructionOperandSeparator();
    PrintImmediate(value);
    PrintInstructionEnd();
}

func PrintInstruction_Imm_Reg(op string, opsize uint64, value uint64, regname string, regnumber uint64, regindirect uint64, regoffset uint64, regoffsetnegative uint64) {
    if (opsize == 1) && (regindirect == 0) { //Clear upper bits using AND mask when operating on bytes; don't clear memory as op could be MOV and therefore set 7 bytes unrecoverably to zero
        PrintInstruction_Imm_Reg("AND", 8, 255, regname, regnumber, regindirect, regoffset, regoffsetnegative); //ANDQ $255, R
    }
    PrintInstructionStart(op, opsize);
    PrintImmediate(value);
    PrintInstructionOperandSeparator();
    PrintRegister(regname, regnumber, regindirect, regoffset, regoffsetnegative);
    PrintInstructionEnd();
}

func PrintInstruction_Imm_Var(op string, value uint64, variable *libgogo.Item) {
    var opsize uint64;
    opsize = GetOpSize(variable);
    if variable.Global == 1 { //Global
        PrintInstruction_Imm_Reg(op, opsize, value, "SB", 0, 1, variable.A, 0); //OP $value, variable.A(SB)
    } else { //Local
        if variable.Global == 2 { //Parameter
            PrintInstruction_Imm_Reg(op, opsize, value, "SP", 0, 1, variable.A + 8, 0); //OP $value, [variable.A+8](SP)
        } else { //Local
            PrintInstruction_Imm_Reg(op, opsize, value, "SP", 0, 1, variable.A + 8, 1); //OP $value, -[variable.A+8](SP)
        }
    }
}

func PrintInstruction_Var_Imm(op string, variable *libgogo.Item, value uint64) {
    var opsize uint64;
    opsize = GetOpSize(variable);
    if variable.Global == 1 { //Global
        PrintInstruction_Reg_Imm(op, opsize, "SB", 0, 1, variable.A, 0, value); // OP variable.A(SB), $value
    } else { //Local
        if variable.Global == 2 { //Parameter
            PrintInstruction_Reg_Imm(op, opsize, "SP", 0, 1, variable.A + 8, 0, value); // OP [variable.A+8](SP), $value
        } else { //Local
            PrintInstruction_Reg_Imm(op, opsize, "SP", 0, 1, variable.A + 8, 1, value); // OP -[variable.A+8](SP), $value
        }
    }
}

func PrintInstruction_Reg_Var(op string, regname string, regnumber uint64, variable *libgogo.Item) {
    var opsize uint64;
    opsize = GetOpSize(variable);
    if variable.Global == 1 { //Global
        PrintInstruction_Reg_Reg(op, opsize, regname, regnumber, 0, 0, 0, "SB", 0, 1, variable.A, 0); //OP regname_regnumber, variable.A(SB)
    } else { //Local
        if variable.Global == 2 { //Parameter
            PrintInstruction_Reg_Reg(op, opsize, regname, regnumber, 0, 0, 0, "SP", 0, 1, variable.A + 8, 0); //OP regname_regnumber, [variable.A+8](SP)
        } else { //Local
            PrintInstruction_Reg_Reg(op, opsize, regname, regnumber, 0, 0, 0, "SP", 0, 1, variable.A + 8, 1); //OP regname_regnumber, -[variable.A+8](SP)
        }
    }
}

func PrintInstruction_Var_Reg(op string, variable *libgogo.Item, regname string, regnumber uint64) {
    var opsize uint64;
    opsize = GetOpSize(variable);
    if variable.Global == 1 { //Global
        PrintInstruction_Reg_Reg(op, opsize, "SB", 0, 1, variable.A, 0, regname, regnumber, 0, 0, 0); //OP variable.A(SB), regname_regnumber
    } else { //Local
        if variable.Global == 2 { //Parameter
            PrintInstruction_Reg_Reg(op, opsize, "SP", 0, 1, variable.A + 8, 0, regname, regnumber, 0, 0, 0); //OP [variable.A+8](SP), regname_regnumber
        } else { //Local
            PrintInstruction_Reg_Reg(op, opsize, "SP", 0, 1, variable.A + 8, 1, regname, regnumber, 0, 0, 0); //OP -[variable.A+8](SP), regname_regnumber
        }
    }
}

func PrintJump(jump string, label string) {
    PrintOutput("  ");
    PrintOutput(jump);
    PrintOutput(" ");
    PrintOutput(label);
    PrintInstructionEnd();
}

func PrintLabel(label string) {
    PrintOutput(label);
    PrintOutput(":");
    PrintInstructionEnd();
}
