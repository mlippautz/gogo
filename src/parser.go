// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

var maxDepth uint64 = 10;
var curDepth uint64 = 1;

var Compile uint64 = 0;

var Operators libgogo.Stack;

var InsideFunction uint64 = 0;
var InsideStructDecl uint64 = 0;
var InsideFunctionVarDecl uint64 = 0;

var Ifs uint64 = 0;

//
// Package name of currently processed file
//
var CurrentPackage string = "<no package>";

//
// Main parsing function. Corresponds to the EBNF main structure called 
// go_program.
//
func Parse() {
    tok.id = 0;
    tok.nextChar = 0;
    tok.nextToken = 0;   
    tok.llCnt = 0; 

    libgogo.InitializeStack(&Operators);

    ParsePackageStatement(); 
    ParseImportStatementList();
    ParseStructDeclList();
    ParseVarDeclList();
    ParseFuncDeclList();

    AssertNextToken(TOKEN_EOS);
}

//
// Parses: package identifier
// This is enforced by the go language as first statement in a source file.
//
func ParsePackageStatement() {    
    PrintDebugString("Entering ParsePackageStatement()",1000);
    AssertNextTokenWeak(TOKEN_PACKAGE);
    AssertNextTokenWeak(TOKEN_IDENTIFIER);
    // package ok, value in tok.strValue
    CurrentPackage = tok.strValue;
    PrintDebugString("Leaving ParsePackageStatement()",1000);
}

//
// Parses: { import_stmt }
// Function parsing the whole import block (which is optional) of a go program
//
func ParseImportStatementList() {
    var validImport uint64;
    PrintDebugString("Entering ParseImportStatementList()",1000);
    for validImport = ParseImportStatement();
        validImport == 0;
        validImport = ParseImportStatement() { }
    PrintDebugString("Leaving ParseImportStatementList()",1000);
}

//
// Parses: "import" string
// This function parses a single import line.
// Returning 0 if import statement is valid, 1 otherwise.
//
func ParseImportStatement() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseImportStatement()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_IMPORT {
        AssertNextToken(TOKEN_STRING);
        // import ok, value in tok.strValue
        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }    
    PrintDebugString("Leaving ParseImportStatement()",1000);
    return boolFlag;
}

//
// Parses: { struct_decl }
// A list of struct declarations.
//
func ParseStructDeclList() {
    var boolFlag uint64;
    PrintDebugString("Entering ParseStructDeclList()",1000);
    for boolFlag = ParseStructDecl();
        boolFlag == 0;
        boolFlag = ParseStructDecl() { }
    PrintDebugString("Leaving ParseStructDeclList()",1000);
}

//
// Parses: "type" identifier "struct" "{" struct_var_decl_list "}" ";"
// This is basically the skeleton of a struct.
//
func ParseStructDecl() uint64 {
    var boolFlag uint64;
    var dontAddType uint64;
    PrintDebugString("Entering ParseStructDecl()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_TYPE {
        AssertNextToken(TOKEN_IDENTIFIER);
        // identifier of struct in tok.strValue
        dontAddType = NewType(tok.strValue); //Create type
        AssertNextToken(TOKEN_STRUCT);
        AssertNextToken(TOKEN_LCBRAC);
        InsideStructDecl = 1;
        ParseStructVarDeclList();
        InsideStructDecl = 0;
        AssertNextTokenWeak(TOKEN_RCBRAC);
        AssertNextTokenWeak(TOKEN_SEMICOLON);
        AddType(dontAddType); //Add type
        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }    
    PrintDebugString("Leaving ParseStructDecl()",1000);
    return boolFlag;
}

//
// Parses: { struct_var_decl }
// The variable declaration list of a struct.
//
func ParseStructVarDeclList() {
    var boolFlag uint64;
    PrintDebugString("Entering ParseStructVarDeclList()",1000);
    for boolFlag = ParseStructVarDecl();
        boolFlag == 0;
        boolFlag = ParseStructVarDecl() { }
    PrintDebugString("Leaving ParseStructVarDeclList()",1000);
}

//
// Parses: identifier type ";"
// A single variable declaration in a struct.
//
func ParseStructVarDecl() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseStructVarDecl()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_IDENTIFIER {
        AddStructField(tok.strValue); //Add field to struct
        ParseType();
        AssertNextTokenWeak(TOKEN_SEMICOLON);
        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }    
    PrintDebugString("Leaving ParseStructVarDecl()",1000);
    return boolFlag;
}

//
// Parses: [ "[" integer "]" ] identifier 
// Use for variable declarations in a struct and for declarations in function
// heads and functions
//
func ParseType() {
    var arraydim uint64 = 0;
    var boolFlag uint64;
    var packagename string = CurrentPackage;
    var typename string;

    PrintDebugString("Entering ParseType()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_LSBRAC {   
        AssertNextToken(TOKEN_INTEGER);
        // value of integer in tok.intValue
        arraydim = tok.intValue;
        AssertNextToken(TOKEN_RSBRAC);
    } else {
        tok.nextToken = tok.id;
    }
    GetNextTokenSafe();
    if tok.id != TOKEN_ARITH_MUL {
        tok.nextToken = tok.id;
    } else {
        SetCurrentObjectTypeToPointer(); //Pointer type (indicated by *)
    }
    AssertNextToken(TOKEN_IDENTIFIER);
    // typename in tok.strValue
    typename = tok.strValue;
    boolFlag = ParseSimpleSelector();
    if boolFlag == 0 { //Take selector into consideration if there is one
        packagename = typename; //Previously read namespace is actually the namespace
        typename = tok.strValue;
    }
    SetCurrentObjectType(typename, packagename, arraydim);
    PrintDebugString("Leaving ParseType()",1000);
}

//
// Parses: [ "[" integer "]" ] identifier
// Is completelly optional. Only used to parse return value of a function
// declaration.
//
func ParseTypeOptional() {
    PrintDebugString("Entering ParseTypeOptional()",1000); 
    GetNextTokenSafe();
    if tok.id == TOKEN_LSBRAC {
        AssertNextToken(TOKEN_INTEGER);        
        AssertNextToken(TOKEN_RSBRAC);
    } else {
        tok.nextToken = tok.id;
    }
    if tok.id == TOKEN_ARITH_MUL {
        // Return type is a pointer
        GetNextTokenSafe();
    }

    GetNextTokenSafe();
    if tok.id != TOKEN_IDENTIFIER  {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseTypeOptional()",1000);
}

//
// Parses: { var_decl }
// Is used for a list of variable declarations. Can either be global or in 
// functions.
//
func ParseVarDeclList() {
    var boolFlag uint64;
    PrintDebugString("Entering ParseVarDeclList()",1000);
    for boolFlag = ParseVarDecl();
        boolFlag == 0;
        boolFlag = ParseVarDecl() { }
    PrintDebugString("Leaving ParseVarDeclList()",1000);
}

//
// Parses: "var" identifier type [ "=" expression ];
// Is used to parse a single variable declaration with optional initializer.
//
func ParseVarDecl() uint64 {
    var boolFlag uint64;
    var ed ExpressionDescriptor;
    var exprIndicator uint64;
    var LHSItem *libgogo.Item;
    var RHSItem *libgogo.Item;
    PrintDebugString("Entering ParseVarDecl()",1000);
    boolFlag = LookAheadAndCheck(TOKEN_VAR);
    if boolFlag == 0 {
        AssertNextToken(TOKEN_VAR);
        AssertNextToken(TOKEN_IDENTIFIER);
        NewVariable(tok.strValue); //New object
        // variable name in tok.strValue
        ParseType();

        GetNextTokenSafe();
        if tok.id == TOKEN_ASSIGN {
            if Compile != 0 {
                LHSItem = libgogo.NewItem();
                RHSItem = libgogo.NewItem();
                if InsideFunction != 0 { //Local variable
                   GenerateComment("Local variable assignment start");
                   VariableObjectDescToItem(CurrentObject, LHSItem, 0); //Local variable
                   GenerateComment("Local variable assignment RHS load start");
                   exprIndicator = ParseExpression(RHSItem, &ed); //Parse RHS
                   GenerateComment("Local variable assignment RHS load end");
                   GenerateAssignment(LHSItem, RHSItem, exprIndicator); //LHS = RHS
                   GenerateComment("Local variable assignment end");
                } else { //Global variable
                   SwitchOutputToInitCodeSegment(); //Write code to main.init in order to make sure that the variables are initialized globaly
                   GenerateComment("Global variable assignment start");
                   VariableObjectDescToItem(CurrentObject, LHSItem, 1); //Global variable
                   GenerateComment("Global variable assignment RHS load start");
                   exprIndicator = ParseExpression(RHSItem, &ed); //Parse RHS
                   GenerateComment("Global variable assignment RHS load end");
                   GenerateAssignment(LHSItem, RHSItem, exprIndicator); //LHS = RHS
                   GenerateComment("Global variable assignment end");
                   SwitchOutputToCodeSegment(); //Write the rest of the code to the code segment
                }
            } else {
                ParseExpression(nil, &ed);
            }
        } else {
            tok.nextToken = tok.id;
        } 

        AssertNextTokenWeak(TOKEN_SEMICOLON);
        boolFlag = 0;
    }
    PrintDebugString("Leaving ParseVarDecl()",1000);
    return boolFlag;
}

//
//
//
func ParseExpression(item *libgogo.Item, ed *ExpressionDescriptor) uint64 {
    var boolFlag uint64;
    var op uint64;
    var tempItem2 *libgogo.Item;
    var retValue uint64 = 0;
    if ed != nil {
        ed.RestCounter = 0;
        ed.ExpressionDepth = ed.ExpressionDepth + 1;
    }
    PrintDebugString("Entering ParseExpression()",1000);
    if item == nil {
        item = libgogo.NewItem();
    }
    GetNextTokenSafe();
    if tok.id == TOKEN_OP_ADR {
        AssertNextToken(TOKEN_IDENTIFIER);
        FindIdentifierAndParseSelector(item);
        retValue = 1;
    } else {
        tok.nextToken = tok.id;   
        ParseSimpleExpression(item, ed);
        boolFlag = ParseCmpOp();
        if boolFlag == 0 {
            tempItem2 = libgogo.NewItem();
            ParseSimpleExpression(tempItem2, ed);
            op = libgogo.Pop(&Operators);
            GenerateComparison(item, tempItem2, op);
	    }
    }
    PrintDebugString("Leaving ParseExpression()",1000);
    return retValue;
}

//
//
//
func ParseCmpOp() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseCmpOp()",1000);
    GetNextTokenSafe();
    if (tok.id == TOKEN_EQUALS) || (tok.id == TOKEN_NOTEQUAL) || 
        (tok.id == TOKEN_REL_LT) || (tok.id == TOKEN_REL_LTOE) || 
        (tok.id == TOKEN_REL_GT) || (tok.id == TOKEN_REL_GTOE) {
        libgogo.Push(&Operators, tok.id);
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseCmpOp()",1000);
    return boolFlag;
}

//
//
//
func ParseSimpleExpression(item *libgogo.Item, ed *ExpressionDescriptor) {
    var boolFlag uint64;
    var tempItem2 *libgogo.Item;
    var op uint64;
    PrintDebugString("Entering ParseSimpleExpression()",1000);
    ParseUnaryArithOp();
    tempItem2 = libgogo.NewItem();
    ParseTerm(item, ed);
    for boolFlag = ParseSimpleExpressionOp(tempItem2, ed);
        boolFlag == 0;
        boolFlag = ParseSimpleExpressionOp(tempItem2, ed) {
        op = libgogo.Pop(&Operators);
        if op != TOKEN_REL_OR {
            GenerateSimpleExpressionArith(item, tempItem2, op);
        }
    }
    PrintDebugString("Leaving ParseSimpleExpression()",1000);
}

//
//
//
func ParseSimpleExpressionOp(item *libgogo.Item, ed *ExpressionDescriptor) uint64 {
    var boolFlag uint64 = 1;
    PrintDebugString("Entering ParseSimpleExpressionOp()",1000);
    boolFlag = ParseUnaryArithOp(); // +,-
    if boolFlag != 0 {
        GetNextTokenSafe();
        if tok.id == TOKEN_REL_OR {
            // ||
            GenerateRelative(item, TOKEN_REL_OR, ed);
            libgogo.Push(&Operators, TOKEN_REL_OR);
            boolFlag = 0;
        } 
    }
    if boolFlag == 0 {
        ParseTerm(item, ed);
    } else {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseSimpleExpressionOp()",1000);
    return boolFlag;
}

//
// Function parsing the unary arithmetic ops PLUS (+) and MINUS (-)
// Returns: 0 if matched, 1 otherwise.
//
func ParseUnaryArithOp() uint64 {
    var boolFlag uint64 = 1;
    PrintDebugString("Entering ParseUnaryArithOp()",1000);
    GetNextTokenSafe();
    if (tok.id == TOKEN_ARITH_PLUS) || (tok.id == TOKEN_ARITH_MINUS) {
        libgogo.Push(&Operators, tok.id);
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseUnaryArithOp()",1000);
    return boolFlag;   
}

//
//
//
func ParseTerm(item *libgogo.Item, ed *ExpressionDescriptor) {
    var boolFlag uint64;
    var tempItem2 *libgogo.Item;
    var op uint64;
    PrintDebugString("Entering ParseTerm()",1000);
    ParseFactor(item, ed);
    tempItem2 = libgogo.NewItem();
    for boolFlag = ParseTermOp(tempItem2, ed);
        boolFlag == 0;
        boolFlag = ParseTermOp(tempItem2, ed) {
        op = libgogo.Pop(&Operators);
        if op != TOKEN_REL_AND {
            GenerateTermArith(item, tempItem2, op);
        }
    }      
    PrintDebugString("Leaving ParseTerm()",1000);
}

//
//
//
func ParseTermOp(item *libgogo.Item, ed *ExpressionDescriptor) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseTermOp()",1000);
    boolFlag = ParseBinaryArithOp(); // *,/
    if boolFlag != 0 {
        GetNextTokenSafe();
        if tok.id == TOKEN_REL_AND {
            // &&
            GenerateRelative(item, TOKEN_REL_AND, ed);
            libgogo.Push(&Operators, TOKEN_REL_AND);
            boolFlag = 0;
        }
    }
    if boolFlag == 0 {
        ParseFactor(item, ed);
    } else {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseTermOp()",1000);
    return boolFlag;
}

//
// Function parsing the binary airthmetic ops MUL (*) and DIV (/)
// Returns: 0 if matched, 1 otherwise.
//
func ParseBinaryArithOp() uint64 {
    var boolFlag uint64 = 1;
    PrintDebugString("Entering ParseBinaryArithOp()",1000);
    GetNextTokenSafe();
    if (tok.id == TOKEN_ARITH_MUL) || (tok.id == TOKEN_ARITH_DIV) {
        libgogo.Push(&Operators, tok.id);
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseBinaryArithOp()",1000);
    return boolFlag; 
}

//
//
//
func ParseFactor(item *libgogo.Item, ed *ExpressionDescriptor) uint64 {
    var doneFlag uint64 = 1;
    var boolFlag uint64;
    var es [2]uint64;
    var tempString *libgogo.ObjectDesc;

    GetNextTokenSafe();
    if (doneFlag == 1) && (tok.id == TOKEN_IDENTIFIER) {
        FindIdentifierAndParseSelector(item);
        doneFlag = 0;
    } 
    if (doneFlag == 1) && (tok.id == TOKEN_INTEGER) {
        if (tok.intValue <= 255) { //Value fits into byte_t
            libgogo.SetItem(item, libgogo.MODE_CONST, byte_t, 0, tok.intValue, 0, 0); //Constant item
        } else { //Value does not fit into byte_t
            libgogo.SetItem(item, libgogo.MODE_CONST, uint64_t, 0, tok.intValue, 0, 0); //Constant item
        }
        doneFlag = 0;
    }
    if (doneFlag) == 1 && (tok.id == TOKEN_STRING) {
        //TODO: Actually declare string in data segment and refer to it in form of a variable (see below)
        tempString = libgogo.NewObject("unique tempstring (TODO)", ".internal", libgogo.CLASS_VAR);
        tempString.ObjType = string_t;
        libgogo.AppendObject(tempString, GlobalObjects);
        VariableObjectDescToItem(tempString, item, 1); //Global variable
        doneFlag = 0;
    }
    if (doneFlag) == 1 && (tok.id == TOKEN_LBRAC) {
        ParseExpression(item, ed);
        AssertNextTokenWeak(TOKEN_RBRAC);
        doneFlag = 0;
    }
    if (doneFlag == 1) && (tok.id == TOKEN_NOT) {
        //libgogo.Push(&Operators, TOKEN_NOT);
        ParseFactor(item, ed);
        //TODO: Generate code
        doneFlag = 0;
    }

    if doneFlag != 0 {
        boolFlag = 1;
        tok.nextToken = tok.id;
        // Fix (?) empty factor, which should not be possible.
        ParseErrorWeak(tok.id,es,0);
        ParserSync();
    } else {
        boolFlag = 0;
    }
    PrintDebugString("Leaving ParseFactor()",1000);
    return boolFlag;
}

func ParseSimpleSelector() uint64 {
    var boolFlag uint64;
    GetNextTokenSafe();
    if tok.id == TOKEN_PT {
        AssertNextToken(TOKEN_IDENTIFIER);
        // value in tok.strValue
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    return boolFlag;
}

//
//
//
func ParseSelector(item *libgogo.Item, packagename string) {
    var boolFlag uint64;
    if item == nil { //TODO: Remove. This is only temporary until all calls to ParseSelector are fully implemented and correct
        item = libgogo.NewItem();
        libgogo.SetItem(item, 0, nil, 0, 0, 0, 0); //Mark item as not being set
    }
    PrintDebugString("Entering ParseSelector()",1000);
    for boolFlag = ParseSelectorSub(item, packagename);
        boolFlag == 0; 
        boolFlag = ParseSelectorSub(item, packagename) {
    }
    PrintDebugString("Leaving ParseSelector()",1000);
}

func ParseSelector_FunctionCall() {
   var boolFlag uint64;
   PrintDebugString("Entering ParseSelector_FunctionCall()",1000);
    for boolFlag = ParseSelectorSub_FunctionCall();
        boolFlag == 0; 
        boolFlag = ParseSelectorSub_FunctionCall() {
    }
    PrintDebugString("Leaving ParseSelector_FunctionCall()",1000);
}

//
//
//
func ParseSelectorSub(item *libgogo.Item, packagename string) uint64 {
    var boolFlag uint64;
    var tempObject *libgogo.ObjectDesc;
    var tempList *libgogo.ObjectDesc;
    var tempItem *libgogo.Item;
    PrintDebugString("Entering ParseSelectorSub()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_PT {
        AssertNextToken(TOKEN_IDENTIFIER);
        // value in tok.strValue
        if Compile != 0 {
            if (item.Itemtype == nil) && (item.A == 0) && (item.R == 0) { //Item undefined, only package known => Find object
                tempObject = libgogo.GetObject(tok.strValue, packagename, LocalObjects); //Check local objects
			    tempList = LocalObjects;
			    if tempObject == nil {
			        tempObject = libgogo.GetObject(tok.strValue, packagename, GlobalObjects); //Check global objects
			        tempList = GlobalObjects;
			    }
			    if tempObject == nil {
			        SymbolTableError("Package", packagename, "has no variable named", tok.strValue);
			    }
			    if tempList == LocalObjects { //Local
			        VariableObjectDescToItem(tempObject, item, 0); //Local variable
			    } else { //Global
			        VariableObjectDescToItem(tempObject, item, 1); //Global variable
		  	    }
            } else { //Field access
                if Compile != 0 {
                    if item.Itemtype == nil {
                        SymbolTableError("Type has", "no", "fields:", "?");
                    }
                    if item.Itemtype.Form != libgogo.FORM_STRUCT { //Struct check
                        SymbolTableError("Type is", "not a", "struct:", item.Itemtype.Name);
                    } else {
                        boolFlag = libgogo.HasField(tok.strValue, item.Itemtype);
                        if boolFlag == 0 { //Field check
                            SymbolTableError("Type", item.Itemtype.Name, "has no field named", tok.strValue);
                        } else {
                            tempObject = libgogo.GetField(tok.strValue, item.Itemtype);
                            boolFlag = libgogo.GetFieldOffset(tempObject, item.Itemtype); //Calculate offset
                            GenerateFieldAccess(item, boolFlag);
                            item.Itemtype = tempObject.ObjType; //Set item type to field type
                            item.PtrType = tempObject.PtrType;
                        }
                    }
                }
            }
        }
        boolFlag = 0;
    } else {
        if tok.id == TOKEN_LSBRAC {
            if Compile != 0 {
                if (item.Itemtype == nil) && (item.A == 0) && (item.R == 0) {
                    SymbolTableError("No array access possible to", "", "package", packagename);
                }
                if item.Itemtype.Form != libgogo.FORM_ARRAY { //Array check
                    SymbolTableError("Type is", "not an", "array:", item.Itemtype.Name);
                }
            }
            GetNextTokenSafe();
            if (Compile != 0) && (item.Itemtype == string_t) { //Derefer string address at offset 0 to access actual byte array of characters
                boolFlag = item.PtrType; //Save old value of PtrType
                item.PtrType = 1; //Force deref.
                DereferItemIfNecessary(item); //Actual deref.
                item.PtrType = boolFlag; //Restore old value of PtrType
            }
            if tok.id == TOKEN_INTEGER {
                if Compile != 0 {
                    boolFlag = libgogo.GetTypeSize(item.Itemtype.Base); //Get unaligned array base type size
                    GenerateFieldAccess(item, boolFlag * tok.intValue); //Direct field access to field offset tok.intValue times base size tok.intValue
                    item.Itemtype = item.Itemtype.Base; //Set item type to array base type
                    item.PtrType = 0; //No pointers to array allowed by EBNF
                }
            } else {
                if tok.id == TOKEN_IDENTIFIER {
                    if Compile != 0 {
                        tempItem = libgogo.NewItem();
                        FindIdentifierAndParseSelector(tempItem); //Load identifier's value
                        boolFlag = libgogo.GetTypeSize(item.Itemtype.Base); //Get unaligned array base type size
                        GenerateVariableFieldAccess(item, tempItem, boolFlag); //Direct field access to field offset identifier value * boolflag (base type size)
                        item.Itemtype = item.Itemtype.Base; //Set item type to array base type
                        item.PtrType = 0; //No pointers to array allowed by EBNF
                    } else {
                        ParseSelector(item, CurrentPackage);
                    }
                }
            } 

            AssertNextToken(TOKEN_RSBRAC);
            boolFlag = 0;
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;
        }
    }
    PrintDebugString("Leaving ParseSelectorSub()",1000);
    return boolFlag;
}

func ParseSelectorSub_FunctionCall() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseSelectorSub_FunctionCall()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_PT {
        AssertNextToken(TOKEN_IDENTIFIER);
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseSelectorSub_FunctionCall()",1000);
    return boolFlag;
}

//
//
//
func ParseFuncDeclList() {
    var boolFlag uint64; 
    PrintDebugString("Entering ParseFuncDeclList()",1000);
    for boolFlag = ParseFuncDeclListSub();
        boolFlag == 0; 
        boolFlag = ParseFuncDeclListSub() { }
    PrintDebugString("Leaving ParseFuncDeclList()",1000);
}

//
//
//
func ParseFuncDeclListSub() uint64 {
    var es [2]uint64;
    var boolFlag uint64;    
    PrintDebugString("Entering ParseFuncDeclListSub()",1000);
    boolFlag = ParseFuncDeclHead();
    if boolFlag == 0 {
        boolFlag = ParseFuncDeclRaw();
        if boolFlag != 0 {
            boolFlag = ParseFuncDecl();
        }
        if boolFlag != 0 {
            es[0] = TOKEN_SEMICOLON;
            es[1] = TOKEN_LCBRAC;
            ParseErrorFatal(tok.id,es,2);
        }
    }
    PrintDebugString("Leaving ParseFuncDeclListSub()",1000);
    return boolFlag;
}

func ParseFuncDeclHead() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseFuncDeclHead()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_FUNC {
        AssertNextToken(TOKEN_IDENTIFIER);
        // function name in tok.strValue
        AssertNextToken(TOKEN_LBRAC);
        ParseIdentifierTypeList();
        AssertNextTokenWeak(TOKEN_RBRAC);
        ParseTypeOptional();
        boolFlag = 0;
    } else {    
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseFuncDeclHead()",1000);
    return boolFlag;
}

func ParseFuncDeclRaw() uint64 {
    var boolFlag uint64 = 1;
    PrintDebugString("Entering ParseFuncDeclRaw()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_SEMICOLON {
        EndOfFunction(); //Delete local variables etc.
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseFuncDeclRaw()",1000);
    return boolFlag;
}

func ParseFuncDecl() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseFuncDecl()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_LCBRAC {
        InsideFunction = 1;
        ParseVarDeclList();
        ParseStatementSequence();
        GetNextTokenSafe();
        if tok.id == TOKEN_RETURN {
            ParseExpression(nil,nil); //TODO
            AssertNextTokenWeak(TOKEN_SEMICOLON);
        } else {
            tok.nextToken = tok.id;
        }
        AssertNextToken(TOKEN_RCBRAC);
        InsideFunction = 0;
        EndOfFunction(); //Delete local variables etc.
        PrintDebugString("Leaving ParseFuncDecl()",1000);
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    return boolFlag;
}

func ParseIdentifierTypeList() {
    var boolFlag uint64;
    PrintDebugString("Entering ParseIdentifierTypeList()",1000);
    boolFlag = ParseIdentifierType();
    if boolFlag == 0 {
        for boolFlag = ParseIdentifierTypeListSub();
            boolFlag == 0; 
            boolFlag = ParseIdentifierTypeListSub() { }   
    }
    PrintDebugString("Leaving ParseIdentifierTypeList()",1000);
}

func ParseIdentifierTypeListSub() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseIdentifierTypeListSub()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_COLON {
        boolFlag = ParseIdentifierType();
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseIdentifierTypeListSub()",1000);
    return boolFlag;
}

func ParseIdentifierType() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseIdentifierType()",1000);
    GetNextTokenSafe();
    if tok.id != TOKEN_IDENTIFIER {
        tok.nextToken = tok.id;
        boolFlag = 1;        
    } else {
        InsideFunctionVarDecl = 1;
        NewVariable(tok.strValue);
        GetNextTokenSafe();
        if tok.id == TOKEN_ARITH_MUL {
            SetCurrentObjectTypeToPointer();
        } else {
            tok.nextToken = tok.id;
        }
        ParseType();
        InsideFunctionVarDecl = 0;
        boolFlag = 0;
    }
    PrintDebugString("Leaving ParseIdentifierType()",1000);
    return boolFlag;
}

func ParseStatementSequence() {
    var boolFlag uint64; 
    PrintDebugString("Entering ParseStatementSequence()",1000);
    for boolFlag = ParseStatement();
        boolFlag == 0;
        boolFlag = ParseStatement() { }
    PrintDebugString("Leaving ParseStatementSequence()",1000);
}

func ParseStatement() uint64 {
    var boolFlag uint64;
    var doneFlag uint64;
    var funcIndicator uint64;
    PrintDebugString("Entering ParseStatement()",1000);
    doneFlag = 1;

    GetNextTokenSafe();
    if (doneFlag == 1) && (tok.id == TOKEN_IDENTIFIER) {
        funcIndicator = IsFunction();
        if funcIndicator == 1 {
            ParseFunctionCallStatement();
        } else {
            ParseAssignment(1); //Trailing semicolon
        }
        doneFlag = 0;
    }
    
    if (doneFlag == 1) && (tok.id == TOKEN_IF) {
        tok.nextToken = tok.id;
        ParseIfStatement();
        doneFlag = 0;
    }

    if (doneFlag == 1) && (tok.id == TOKEN_FOR) {
        tok.nextToken = tok.id;
        ParseForStatement();
        doneFlag = 0;
    }

    if (doneFlag == 1) && (tok.id == TOKEN_BREAK) {
        // simple break statement
        AssertNextToken(TOKEN_SEMICOLON);
        doneFlag = 0;
    }

    if (doneFlag == 1) && (tok.id == TOKEN_CONTINUE) {
        AssertNextToken(TOKEN_SEMICOLON);
        doneFlag = 0;
    }

    if (doneFlag == 1) && (tok.id == TOKEN_SEMICOLON) {
        // NOP    
        doneFlag = 0;
    }

    if doneFlag != 0 {
        tok.nextToken = tok.id;
        boolFlag = 1;
    } else {
        boolFlag = 0;
    }
    PrintDebugString("Leaving ParseStatement()",1000);
    return boolFlag;
}

//
// This function checks whether an identifier (and its selectors)
// is part of a function call or an assignment.
// Due to Go's library hierachy and syntax this has to be done using LL3
// A '.' + <identifier> + <deciding token>
//
func IsFunction() uint64 {
    var returnValue uint64 = 0;
    var cnt uint64 = 1;
    PrintDebugString("Entering IsFunction()",1000);
    // Current token HAS to be an identifier
    tok.nextTokenId[0] = tok.id;
    tok.nextTokenValStr[0] = tok.strValue;
    tok.toRead = 0;

    GetNextTokenSafe();
    if tok.id == TOKEN_PT { // LL1, but we need another look-ahead
        tok.nextTokenId[1] = tok.id;
        cnt = 2;

        GetNextTokenSafe();
        if tok.id == TOKEN_IDENTIFIER { // LL2, need more
            tok.nextTokenId[2] = tok.id;
            tok.nextTokenValStr[2] = tok.strValue;
            cnt = 3;

            GetNextTokenSafe();
            if tok.id == TOKEN_LBRAC { // LL3, finally can decide it
                returnValue = 1;
            }
        } 
    } else  { // Still LL1 here, so no additional handling needed
        if tok.id == TOKEN_LBRAC { // LL1
            returnValue = 1;
        }
    }
    tok.nextTokenId[cnt] = tok.id;
    tok.nextTokenValStr[cnt]  = tok.strValue;
    tok.llCnt = cnt + 1;
    PrintDebugString("Leaving IsFunction()",1000);
    return returnValue;
}

func ParseAssignment(semicolon uint64) uint64 {
    var boolFlag uint64;
    var exprIndicator uint64;
    var funcIndicator uint64;
    var LHSItem *libgogo.Item;
    var RHSItem *libgogo.Item;
    PrintDebugString("Entering ParseAssignment()",1000);
    GenerateComment("Assignment start");
    AssertNextToken(TOKEN_IDENTIFIER);
    LHSItem = libgogo.NewItem();
    GenerateComment("Assignment LHS load start");
    FindIdentifierAndParseSelector(LHSItem); //Parse LHS
    GenerateComment("Assignment LHS load end");
    GetNextTokenSafe();
    if tok.id == TOKEN_ASSIGN {
        GetNextTokenSafe();
        if tok.id == TOKEN_IDENTIFIER {
            funcIndicator = IsFunction();
            if funcIndicator == 1 { //Function call
                ParseFunctionCallStatement();
            } else { //Expression starting with an identifier
                RHSItem = libgogo.NewItem();
                GenerateComment("Assignment RHS load start");
                exprIndicator = ParseExpression(RHSItem,nil); //Parse RHS
                GenerateComment("Assignment RHS load end");
                GenerateAssignment(LHSItem, RHSItem, exprIndicator); //LHS = RHS
            }
        } else { //Expression
            tok.nextToken = tok.id;
            RHSItem = libgogo.NewItem();
            GenerateComment("Assignment RHS load start");
            exprIndicator = ParseExpression(RHSItem,nil); //Parse RHS
            GenerateComment("Assignment RHS load end");
            GenerateAssignment(LHSItem, RHSItem, exprIndicator); //LHS = RHS
        }
        if semicolon != 0 {
            AssertNextTokenWeak(TOKEN_SEMICOLON);
        }
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    GenerateComment("Assignment end");
    PrintDebugString("Leaving ParseAssignment()",1000);
    return boolFlag;
}

func ParseFunctionCallOptional() {
    PrintDebugString("Entering ParseFunctionCallOptional()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_LBRAC {
        GetNextTokenSafe();
        if tok.id == TOKEN_RBRAC {

        } else {
            tok.nextToken = tok.id;
            ParseExpressionList();
            AssertNextToken(TOKEN_RBRAC);
        }
    } else {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseFunctionCallOptional()",1000);
}

func ParseFunctionCall() {
    PrintDebugString("Entering ParseFunctionCall()",1000);
    AssertNextToken(TOKEN_LBRAC);
    GetNextTokenSafe();
    if tok.id != TOKEN_RBRAC {
        tok.nextToken = tok.id;
        ParseExpressionList();
        AssertNextTokenWeak(TOKEN_RBRAC);     
    }     
    PrintDebugString("Leaving ParseFunctionCall()",1000);
}

func ParseExpressionList() {
    var ed ExpressionDescriptor;
    var boolFlag uint64;
    PrintDebugString("Entering ParseExpressionList()",1000);
    ParseExpression(nil, &ed); //TODO
    for boolFlag = ParseExpressionListSub();
        boolFlag == 0;
        boolFlag = ParseExpressionListSub() { }   
    PrintDebugString("Leaving ParseExpressionList()",1000);
}

func ParseExpressionListSub() uint64 {
    var boolFlag uint64;
    var ed ExpressionDescriptor;
    PrintDebugString("Entering ParseExpressionListSub()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_COLON {
        ParseExpression(nil, &ed); //TODO
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseExpressionListSub()",1000);
    return boolFlag;   
}

func ParseFunctionCallStatement() {
    PrintDebugString("Entering ParseFunctionCallStatement()",1000);
    AssertNextToken(TOKEN_IDENTIFIER); //TODO
    ParseSelector_FunctionCall(); //TODO
    ParseFunctionCall();
    PrintDebugString("Leaving ParseFunctionCallStatement()",1000);
}

func ParseForStatement() {
    var ed ExpressionDescriptor;
    PrintDebugString("Entering ParseForStatement()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_FOR {
        GetNextTokenSafe();

        if tok.id == TOKEN_SEMICOLON {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            ParseAssignment(0); //No semicolon
        }
        
        AssertNextToken(TOKEN_SEMICOLON);

        GetNextTokenSafe();
        if tok.id == TOKEN_SEMICOLON {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            ParseExpression(nil, &ed); //TODO
        }

        AssertNextToken(TOKEN_SEMICOLON);

        GetNextTokenSafe();
        if tok.id == TOKEN_LCBRAC {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            ParseAssignment(0); //No semicolon
        }

        AssertNextToken(TOKEN_LCBRAC);        
        ParseStatementSequence();
        AssertNextToken(TOKEN_RCBRAC);

    } else {
        tok.nextToken = tok.id;
    }   
    PrintDebugString("Leaving ParseForStatement()",1000);
}

//
// Parses: "if" expression "{" stmt_sequence [ "}" "else" else_stmt ].
// Represents an if statement. Else is parsed in a separate function.
//
func ParseIfStatement() {
    var item *libgogo.Item;
    var ed ExpressionDescriptor;
    Ifs = Ifs +1;
    ed.GlobalCounter = Ifs;
    ed.Prefix = "";
    libgogo.StringAppend(&ed.Prefix,fileInfo[curFileIndex].filename);
    ed.ExpressionDepth = 0;
    PrintDebugString("Entering ParseIfStatement()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_IF {
        item = libgogo.NewItem();
        ParseExpression(item, &ed);
        GenerateIfStart(item, ed);
        AssertNextToken(TOKEN_LCBRAC);
        ParseStatementSequence();
        AssertNextToken(TOKEN_RCBRAC);
        GenerateIfEnd(item, ed);
        GetNextTokenSafe();
        if tok.id == TOKEN_ELSE {
            ParseElseStatement();
        } else {
            tok.nextToken = tok.id;
        }

    } else {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseIfStatement()",1000);
}

//
// Parses: "{" stmt_sequence "}"
// The optional else branch of an if.
//
func ParseElseStatement() {
    PrintDebugString("Entering ParseElseStatement()",1000);
    AssertNextTokenWeak(TOKEN_LCBRAC);
    ParseStatementSequence();
    AssertNextToken(TOKEN_RCBRAC);
    PrintDebugString("Leaving ParseElseStatement()",1000);
}

func ParserSync() {
    Compile = 0; // stop producing code
    for ;(tok.id != TOKEN_FUNC) && (tok.id != TOKEN_EOS); {
        GetNextTokenSafe();
    }
    if tok.id == TOKEN_FUNC {
        tok.nextToken = tok.id;
        ParseFuncDeclList();       
    }
    libgogo.Exit(3); // Exit with an error
}

func FindIdentifierAndParseSelector(item *libgogo.Item) {
    var boolFlag uint64;
    var tempObject *libgogo.ObjectDesc;
    var tempList *libgogo.ObjectDesc;
    var packagename string;
    if Compile != 0 {
		//Token can be package name
		boolFlag = libgogo.FindPackageName(tok.strValue, GlobalObjects); //Check global objects
		if boolFlag == 0 {
		    boolFlag = libgogo.FindPackageName(tok.strValue, LocalObjects); //Check local objects
		}
		if boolFlag == 0 {
		    boolFlag = libgogo.FindPackageName(tok.strValue, LocalParameters); //Check local parameters
		}
		if boolFlag == 0 { //Token is not package name, but identifier
			tempObject = libgogo.GetObject(tok.strValue, CurrentPackage, LocalObjects); //Check local objects
			tempList = LocalObjects;
			if tempObject == nil {
       			tempObject = libgogo.GetObject(tok.strValue, CurrentPackage, LocalParameters); //Check local parameters
    			tempList = LocalParameters;
    			if tempObject == nil {
    				tempObject = libgogo.GetObject(tok.strValue, CurrentPackage, GlobalObjects); //Check global objects
	    			tempList = GlobalObjects;
	    		}
			}
			if tempObject == nil {
				SymbolTableError("Undefined", "", "variable", tok.strValue);
			}
			if tempList == LocalObjects { //Local
    			VariableObjectDescToItem(tempObject, item, 0); //Local variable
			} else { //Global or parameter
			    if tempList == GlobalObjects { //Global
    				VariableObjectDescToItem(tempObject, item, 1); //Global variable
    			} else { //Parameter
    				VariableObjectDescToItem(tempObject, item, 2); //Local parameter
    	        }
			}
		    ParseSelector(item, CurrentPackage); //Parse selectors for an object in the current package
		} else { //Token is package name
		    libgogo.SetItem(item, 0, nil, 0, 0, 0, 0); //Mark item as not being set
		    packagename = tok.strValue; //Save package name
		    ParseSelector(item, tok.strValue); //Parse selectors for an undefined object in the given package
		    if (item.Itemtype == nil) && (item.A == 0) && (item.R == 0) {
		        SymbolTableError("Cannot use package", "", "as a variable:", packagename);
		    }
		}
    } else {
        ParseSelector(item, CurrentPackage);
    }
}
