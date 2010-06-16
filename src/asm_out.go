// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

//
// Variables holding the compiled code and data segment (as assembly code)
//
var Header string; //Header
var DataSegmentList libgogo.StringList;
var DataSegment string; //Data section for global variables
var InitCodeSegmentList libgogo.StringList;
var InitCodeSegment string; //Global variable initialization code
var CodeSegmentList libgogo.StringList;
var CodeSegment string; //Actual code

//
// String lists containing the symbol tables
//
var FunctionSymbolTable libgogo.StringList;
var TypeSymbolTable libgogo.StringList;

//
// Pointer used to decide to which string to append output
// There are functions like SwitchOutputToCodeSegment() to change this pointer in a clean way
// By default, it points to the code segment string
//
var OutputStringPtr *string = &CodeSegment;

//
// Size of the data segment
//
var DataSegmentSize uint64 = 0;

//
// Reseting the code before a new compile round starts
//
func ResetCode() {
    Header = "\
//\n\
// --------------------\n\
// GoGo compiler output\n\
// --------------------\n\
//\n\
";
    libgogo.StringAppend(&Header, fileInfo[curFileIndex].filename);
    libgogo.StringAppend(&Header,"\n\
// Syntax: Plan-9 assembler\n\
//\n\
// This code is automatically generated. DO NOT EDIT IT!\n\
//\n\
");
    //InspectorGadget();
    InitCodeSegment = "TEXT main·init(SB),0,$0-0\n";
    
    libgogo.InitializeStringList(&DataSegmentList);
    libgogo.InitializeStringList(&InitCodeSegmentList);
    libgogo.InitializeStringList(&CodeSegmentList);
    
    libgogo.InitializeStringList(&FunctionSymbolTable);
    libgogo.InitializeStringList(&TypeSymbolTable);
}

//
// Function printing the generated code (stored in Code) to a file called the
// same name as the input file + '.sog' extension
//
func PrintFile(Functions *libgogo.TypeDesc, Types *libgogo.TypeDesc) {
    var fd uint64;
    var outfile string = "_gogo_.sog";
    var DataSegmentSizeStr string;
    // The following line creates a new file for the assembler code
    // flags: O_WRONLY | O_CREAT | O_TRUNC => 577
    // mode: S_IWUSR | S_IRUSR | S_IRGRP => 416
    fd = libgogo.FileOpen2(outfile,577,416);

    libgogo.WriteString(fd, Header);
    libgogo.WriteString(fd, "\n"); //Separator
    if NeedsLink != 0 {
        libgogo.WriteString(fd, "__UNLINKED_CODE\n");
    }
    libgogo.WriteString(fd, "//Symbol table:\n"); //Separator
    libgogo.SymbolTableTypesToStringList(Types, &TypeSymbolTable);
    PrintStringList(fd, &TypeSymbolTable, 1); //Type symbol table
    libgogo.SymbolTableFunctionsToStringList(Functions, &FunctionSymbolTable);
    PrintStringList(fd, &FunctionSymbolTable, 1); //Function symbol table
    libgogo.WriteString(fd, "\n"); //Separator
    PrintStringList(fd, &DataSegmentList, 0);
    libgogo.WriteString(fd, DataSegment); //Data segment
    libgogo.WriteString(fd, "GLOBL data(SB),$"); //(Begin of) end of data segment
    DataSegmentSizeStr = libgogo.IntToString(DataSegmentSize);
    libgogo.WriteString(fd, DataSegmentSizeStr); //Size of data segment
    libgogo.WriteString(fd, "\n"); //End of data segment
    libgogo.WriteString(fd, "\n"); //Separator
    PrintStringList(fd, &InitCodeSegmentList, 0);
    libgogo.WriteString(fd, InitCodeSegment); //main.init
    libgogo.WriteString(fd, "  RET\n"); //End of function (main.init)
    libgogo.WriteString(fd, "\n"); //Separator
    PrintStringList(fd, &CodeSegmentList, 0);
    libgogo.WriteString(fd, CodeSegment); //Code segment
    libgogo.FileClose(fd);
}

func PrintStringList(fd uint64, list *libgogo.StringList, comment uint64) {
    var i uint64;
    var n uint64;
    var temp string;
    n = libgogo.GetStringListItemCount(list);
    for i = 0; i < n; i = i + 1 {
        temp = libgogo.GetStringItemAt(list, i);
        if comment != 0 {
            libgogo.WriteString(fd, "//");
        }
        libgogo.WriteString(fd, temp);
        if comment != 0 {
            libgogo.WriteString(fd, "\n");
        }
    }
}

func SwitchOutputToInitCodeSegment() {
    OutputStringPtr = &InitCodeSegment;
}

func SwitchOutputToCodeSegment() {
    OutputStringPtr = &CodeSegment;
}

func SwitchOutputToDataSegment() {
    OutputStringPtr = &DataSegment;
}

func PrintCodeOutput(output string) {
    var tempPtr *string;
    var oldLength uint64;
    var appendLength uint64;
    oldLength = libgogo.StringLength2(OutputStringPtr);
    appendLength = libgogo.StringLength(output);
    if oldLength + appendLength >= 65534 { //If code output exceeds max. string size allowed by Go runtime...
        tempPtr = &DataSegment;
        if OutputStringPtr == tempPtr {
            libgogo.AddStringItem(&DataSegmentList, DataSegment); //...save the current code to the code list...
            //DataSegment = "";
        }
        tempPtr = &InitCodeSegment;
        if OutputStringPtr == tempPtr {
            libgogo.AddStringItem(&InitCodeSegmentList, InitCodeSegment); //...save the current code to the code list...
            //InitCodeSegment = "";
        }
        tempPtr = &CodeSegment;
        if OutputStringPtr == tempPtr {
            libgogo.AddStringItem(&CodeSegmentList, CodeSegment); //...save the current code to the code list...
            //CodeSegment = ""; 
        }
        libgogo.ResetString(OutputStringPtr); //... and start over with a new output string
    }
    libgogo.StringAppend(OutputStringPtr, output);
}

func PrintCodeOutputChar(output byte) {
    libgogo.CharAppend(OutputStringPtr, output);
}

func PrintCodeOutputValue(value uint64) {
    var temp string;
    temp = libgogo.IntToString(value);
    PrintCodeOutput(temp);
}

func GenerateComment(msg string) {
    var str string = "";
    var tmpPtr *string;
    var temp string;
    var i uint64;
    var n uint64;
    if DEBUG_LEVEL >= 10 {
        tmpPtr = &DataSegment;
        if (OutputStringPtr != tmpPtr) {
            str = "  //--- ";
        } else { //No indentation in data segment
            str = "//--- ";
        }
        n = libgogo.StringLength(msg);
        for i = 0; i < n; i = i + 1 {
            if msg[i] == 10 { //Unescape line breaks in comments to avoid invalid assembly code
                libgogo.StringAppend(&str, "\\n"); //Literal \n, not actual \n
            } else {
                libgogo.CharAppend(&str, msg[i]);
            }
        }
        libgogo.StringAppend(&str, " at ");
        temp = BuildHead();
        libgogo.StringAppend(&str, temp);
        libgogo.StringAppend(&str, "\n");
        PrintCodeOutput(str);
    }
}

func PrintRegister(name string, number uint64, indirect uint64, offset uint64, negativeoffset uint64, optionaloffsetname string) {
    var R uint64;
    if indirect != 0 {
        PrintCodeOutput(optionaloffsetname);
        if offset != 0 {
            if negativeoffset != 0 {
                PrintCodeOutput("-");
            } else {
                R = libgogo.StringLength(optionaloffsetname);
                if R != 0 { //Print positive offset if there is an offset name
                    PrintCodeOutput("+");
                }
            }
            PrintCodeOutputValue(offset);
        }
        PrintCodeOutput("(");
    }
    PrintCodeOutput(name); //Print register name
    R = libgogo.StringCompare(name, "R");
    if R == 0 { //p.e. R8 => consider number; else: p.e. AX => ignore number
        PrintCodeOutputValue(number);
    }
    if indirect != 0 {
        PrintCodeOutput(")");
    }
}

func PrintImmediate(value uint64) {
    PrintCodeOutput("$");
    PrintCodeOutputValue(value);
}

func GetOpSize(item *libgogo.Item, op string) uint64 {
    var size uint64;
    size = libgogo.StringCompare(op, "LEA");
    if size == 0 { //Always use 64 bits when performing LEA (=> LEAQ) as reading something else than a 64 bit pointer on a 64 bit architecture does not make any sense
        size = 8;
    } else {
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
    }
    return size;
}

func PrintInstructionStart(op string, opsize uint64) uint64 {
    var retVal uint64 = 0;
    PrintCodeOutput("  ");
    PrintCodeOutput(op);
    if opsize == 1 { //byte => B
        PrintCodeOutput("B");
    } else {
        if opsize == 8 { //uint64 => Q
            PrintCodeOutput("Q");
        } else {
            if opsize > 8 { //>uint64 => Q and return rest of size to be handled additionally
                PrintCodeOutput("Q");
                retVal = opsize - 8;
            } else { //Size between 1 and 8 which is neither byte nor uint64
                PrintCodeOutput("?"); //TODO: Error: invalid opsize
            }
        }
    }
    PrintCodeOutput(" ");
    return retVal;
}

func PrintInstructionOperandSeparator() {
    PrintCodeOutput(", ");
}

func PrintInstructionEnd() {
    PrintCodeOutput("\n");
}

func PrintInstruction_Reg(op string, opsize uint64, name string, number uint64, indirect uint64, offset uint64, offsetnegative uint64, optionaloffsetname string) uint64 {
    var retVal uint64;
    retVal = PrintInstructionStart(op, opsize);
    PrintRegister(name, number, indirect, offset, offsetnegative, optionaloffsetname);
    PrintInstructionEnd();
    return retVal;
}

func PrintInstruction_Reg_Reg(op string, opsize uint64, reg1name string, reg1number uint64, reg1indirect uint64, reg1offset uint64, reg1offsetnegative uint64, reg1optionaloffsetname string, reg2name string, reg2number uint64, reg2indirect uint64, reg2offset uint64, reg2offsetnegative uint64, reg2optionaloffsetname string) uint64 {
    var retVal uint64;
    if (opsize == 1) && (reg2indirect == 0) { //Clear upper bits using AND mask when operating on bytes; don't clear memory as op could be MOV and therefore set 7 bytes unrecoverably to zero
        PrintInstruction_Imm_Reg("AND", 8, 255, reg2name, reg2number, reg2indirect, reg2offset, reg2offsetnegative, reg2optionaloffsetname); //ANDQ $255, R
    }
    retVal = PrintInstructionStart(op, opsize);
    PrintRegister(reg1name, reg1number, reg1indirect, reg1offset, reg1offsetnegative, reg1optionaloffsetname);
    PrintInstructionOperandSeparator();
    PrintRegister(reg2name, reg2number, reg2indirect, reg2offset, reg2offsetnegative, reg2optionaloffsetname);
    PrintInstructionEnd();
    return retVal;
}

func PrintInstruction_Reg_Imm(op string, opsize uint64, regname string, regnumber uint64, regindirect uint64, regoffset uint64, regoffsetnegative uint64, regoptionaloffsetname string, value uint64) {
    //No opcode check necessary here
    PrintInstructionStart(op, opsize);
    PrintRegister(regname, regnumber, regindirect, regoffset, regoffsetnegative, regoptionaloffsetname);
    PrintInstructionOperandSeparator();
    PrintImmediate(value);
    PrintInstructionEnd();
}

func PrintInstruction_Imm_Reg(op string, opsize uint64, value uint64, regname string, regnumber uint64, regindirect uint64, regoffset uint64, regoffsetnegative uint64, regoptionaloffsetname string) {
    if (opsize == 1) && (regindirect == 0) { //Clear upper bits using AND mask when operating on bytes; don't clear memory as op could be MOV and therefore set 7 bytes unrecoverably to zero
        PrintInstruction_Imm_Reg("AND", 8, 255, regname, regnumber, regindirect, regoffset, regoffsetnegative, regoptionaloffsetname); //ANDQ $255, R
    }
    PrintInstructionStart(op, opsize);
    PrintImmediate(value);
    PrintInstructionOperandSeparator();
    PrintRegister(regname, regnumber, regindirect, regoffset, regoffsetnegative, regoptionaloffsetname);
    PrintInstructionEnd();
}

func PrintInstruction_Imm_Var(op string, value uint64, variable *libgogo.Item) {
    var opsize uint64;
    var temp uint64;
    opsize = GetOpSize(variable, op);
    if variable.Global == 1 { //Global
        PrintInstruction_Imm_Reg(op, opsize, value, "SB", 0, 1, variable.A, 0, "data"); //OP $value, data+variable.A(SB)
    } else { //Local
        if variable.Global == 2 { //Parameter
            PrintInstruction_Imm_Reg(op, opsize, value, "SP", 0, 1, variable.A + 8, 0, ""); //OP $value, [variable.A+8](SP)
        } else { //Local
            temp = libgogo.StringLength(variable.LinkerInformation);
            if temp != 0 {
                GenerateComment(variable.LinkerInformation);
            }
            PrintInstruction_Imm_Reg(op, opsize, value, "SP", 0, 1, variable.A + 8, 1, ""); //OP $value, -[variable.A+8](SP)
        }
    }
}

func PrintInstruction_Var_Imm(op string, variable *libgogo.Item, value uint64) {
    var opsize uint64;
    var temp uint64;
    opsize = GetOpSize(variable, op);
    if variable.Global == 1 { //Global
        PrintInstruction_Reg_Imm(op, opsize, "SB", 0, 1, variable.A, 0, "data", value); // OP data+variable.A(SB), $value
    } else { //Local
        if variable.Global == 2 { //Parameter
            PrintInstruction_Reg_Imm(op, opsize, "SP", 0, 1, variable.A + 8, 0, "", value); // OP [variable.A+8](SP), $value
        } else { //Local
            temp = libgogo.StringLength(variable.LinkerInformation);
            if temp != 0 {
                GenerateComment(variable.LinkerInformation);
            }
            PrintInstruction_Reg_Imm(op, opsize, "SP", 0, 1, variable.A + 8, 1, "", value); // OP -[variable.A+8](SP), $value
        }
    }
}

func PrintInstruction_Reg_Var(op string, regname string, regnumber uint64, optregname string, optregnumber uint64, variable *libgogo.Item) {
    var opsize uint64;
    var retVal uint64;
    var temp uint64;
    opsize = GetOpSize(variable, op);
    if variable.Global == 1 { //Global
        retVal = PrintInstruction_Reg_Reg(op, opsize, regname, regnumber, 0, 0, 0, "", "SB", 0, 1, variable.A, 0, "data"); //OP regname_regnumber, data+variable.A(SB)
        if retVal != 0 { //Handle operands > 8 bytes
            PrintInstruction_Reg_Reg(op, retVal, optregname, optregnumber, 0, 0, 0, "", "SB", 0, 1, variable.A + 8, 0, "data"); //OP optregname_optregnumber, data+variable.A+8(SB)
        }
    } else { //Local
        if variable.Global == 2 { //Parameter
            retVal = PrintInstruction_Reg_Reg(op, opsize, regname, regnumber, 0, 0, 0, "", "SP", 0, 1, variable.A + 8, 0, ""); //OP regname_regnumber, [variable.A+8](SP)
            if retVal != 0 { //Handle operands > 8 bytes
                PrintInstruction_Reg_Reg(op, retVal, optregname, optregnumber, 0, 0, 0, "", "SP", 0, 1, variable.A + 16, 0, ""); //OP optregname_optregnumber, [variable.A+8+8](SP)
            }
        } else { //Local
            temp = libgogo.StringLength(variable.LinkerInformation);
            if temp != 0 {
                GenerateComment(variable.LinkerInformation);
            }
            retVal = PrintInstruction_Reg_Reg(op, opsize, regname, regnumber, 0, 0, 0, "", "SP", 0, 1, variable.A + 8, 1, ""); //OP regname_regnumber, -[variable.A+8](SP)
            if retVal != 0 { //Handle operands > 8 bytes
                if temp != 0 {
                    GenerateComment(variable.LinkerInformation);
                }
                PrintInstruction_Reg_Reg(op, retVal, optregname, optregnumber, 0, 0, 0, "", "SP", 0, 1, variable.A, 1, ""); //OP optregname_optregnumber, -[variable.A+8-8](SP)
            }
        }
    }
}

func PrintInstruction_Var_Reg(op string, variable *libgogo.Item, regname string, regnumber uint64, optregname string, optregnumber uint64) {
    var opsize uint64;
    var retVal uint64;
    var temp uint64;
    opsize = GetOpSize(variable, op);
    if variable.Global == 1 { //Global
        retVal = PrintInstruction_Reg_Reg(op, opsize, "SB", 0, 1, variable.A, 0, "data", regname, regnumber, 0, 0, 0, ""); //OP data+variable.A(SB), regname_regnumber
        if retVal != 0 { //Handle operands > 8 byte
            PrintInstruction_Reg_Reg(op, retVal, "SB", 0, 1, variable.A + 8, 0, "data", optregname, optregnumber, 0, 0, 0, ""); //OP data+variable.A+8(SB), optregname_optregnumber
        }
    } else { //Local
        if variable.Global == 2 { //Parameter
            retVal = PrintInstruction_Reg_Reg(op, opsize, "SP", 0, 1, variable.A + 8, 0, "", regname, regnumber, 0, 0, 0, ""); //OP [variable.A+8](SP), regname_regnumber
            if retVal != 0 { //Handle operands > 8 byte
                PrintInstruction_Reg_Reg(op, retVal, "SP", 0, 1, variable.A + 16, 0, "", optregname, optregnumber, 0, 0, 0, ""); //OP [variable.A+8+8](SP), optregname_optregnumber
            }
        } else { //Local
            temp = libgogo.StringLength(variable.LinkerInformation);
            if temp != 0 {
                GenerateComment(variable.LinkerInformation);
            }
            retVal = PrintInstruction_Reg_Reg(op, opsize, "SP", 0, 1, variable.A + 8, 1, "", regname, regnumber, 0, 0, 0, ""); //OP -[variable.A+8](SP), regname_regnumber
            if retVal != 0 { //Handle operands > 8 byte
                if temp != 0 {
                    GenerateComment(variable.LinkerInformation);
                }
                PrintInstruction_Reg_Reg(op, retVal, "SP", 0, 1, variable.A, 1, "", optregname, optregnumber, 0, 0, 0, ""); //OP -[variable.A+8-8](SP), optregname_optregnumber
            }
        }
    }
}

func PrintJump(jump string, label string) {
    PrintCodeOutput("  ");
    PrintCodeOutput(jump);
    PrintCodeOutput(" ");
    PrintCodeOutput(label);
    PrintInstructionEnd();
}

func PrintLabel(label string) {
    PrintCodeOutput(label);
    PrintCodeOutput(":");
    PrintInstructionEnd();
}

func PrintFunctionCall(packagename string, label string, stackoffset uint64, unknownoffset uint64) {
    var comment string;
    GenerateComment("Stack pointer offset before function call for local variables start");
    if unknownoffset != 0 { //Output linker information
        comment = "##2##";
        libgogo.StringAppend(&comment, packagename);
        libgogo.StringAppend(&comment, "·");
        libgogo.StringAppend(&comment, label);
        libgogo.StringAppend(&comment, "##");
        GenerateComment(comment);
    }
    PrintInstruction_Imm_Reg("SUB", 8, stackoffset, "SP", 0, 0, 0, 0, ""); //SUBQ $stackoffset, SP
    GenerateComment("Stack pointer offset before function call for local variables end");
    PrintCodeOutput("  CALL ");
    PrintCodeOutput(packagename);
    PrintCodeOutput("·");
    PrintCodeOutput(label);
    PrintCodeOutput("(SB)");
    PrintInstructionEnd();
    GenerateComment("Stack pointer offset after function call for local variables start");
    if unknownoffset != 0 { //Output linker information
        GenerateComment(comment);
    }
    PrintInstruction_Imm_Reg("ADD", 8, stackoffset, "SP", 0, 0, 0, 0, ""); //ADDQ $stackoffset, SP
    GenerateComment("Stack pointer offset after function call for local variables end");
}

func PrintFunctionStart(packagename string, label string) {
    PrintCodeOutput("TEXT ");
    PrintCodeOutput(packagename);
    PrintCodeOutput("·");
    PrintCodeOutput(label);
    PrintCodeOutput("(SB),0,$0-0"); //Stack is managed manually
    PrintInstructionEnd();
}

func PrintFunctionEnd() {
    PrintCodeOutput("  RET");
    PrintInstructionEnd();
    PrintInstructionEnd(); //Additional new line separation
}

func SetDataSegmentSize(size uint64) {
    DataSegmentSize = size;
}

func PutDataByte(offset uint64, value byte) {
    var temp uint64;
    PrintCodeOutput("DATA ");
    PrintRegister("SB", 0, 1, offset, 0, "data");
    PrintCodeOutput("/1");
    PrintInstructionOperandSeparator();
    temp = libgogo.ToIntFromByte(value);
    PrintImmediate(temp);
    PrintInstructionEnd();
}
