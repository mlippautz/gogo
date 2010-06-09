// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

var Operators libgogo.Stack;

var InsideFunction uint64 = 0;
var InsideStructDecl uint64 = 0;
var InsideFunctionVarDecl uint64 = 0;

//
// Package name of currently processed file
//
var CurrentPackage string = "<no package>";

//
// Pseudo object representing a function's return value
//
var ReturnValuePseudoObject *libgogo.ObjectDesc = nil;

//
// Temporary function object to reassign pointer in ParseSelectorSub_FunctionCall
//
var ReturnedFunction *libgogo.TypeDesc = nil;

//
// Main parsing function. Corresponds to the EBNF main structure called 
// go_program.
//
func Parse() {
    ResetToken();

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
    //TODO: Type check if function was forward decl.
    SetCurrentObjectType(typename, packagename, arraydim);
    PrintDebugString("Leaving ParseType()",1000);
}

//
// Parses: [ "[" integer "]" ] identifier
// Is completely optional. Only used to parse return value of a function
// declaration.
//
func ParseTypeOptional() {
    var arraydim uint64 = 0;
    var boolFlag uint64;
    var packagename string = CurrentPackage;
    var typename string;

    PrintDebugString("Entering ParseTypeOptional()",1000); 
    ReturnValuePseudoObject = nil;
    CurrentObject = libgogo.NewObject("return value", "", libgogo.CLASS_PARAMETER); //Return value object with a name which is impossible to declare (contains spaces) and therefore needs no additional checking
    GetNextTokenSafe();
    if tok.id == TOKEN_LSBRAC {
        AssertNextToken(TOKEN_INTEGER);        
        arraydim = tok.intValue;
        AssertNextToken(TOKEN_RSBRAC);
    } else {
        tok.nextToken = tok.id;
    }
    if tok.id == TOKEN_ARITH_MUL {
        // Return type is a pointer
        SetCurrentObjectTypeToPointer(); //Pointer type (indicated by *)
        GetNextTokenSafe();
    }

    GetNextTokenSafe();
    if tok.id != TOKEN_IDENTIFIER  {
        tok.nextToken = tok.id;
    } else {
        typename = tok.strValue;
        boolFlag = ParseSimpleSelector();
        if boolFlag == 0 { //Take selector into consideration if there is one
            packagename = typename; //Previously read namespace is actually the namespace
            typename = tok.strValue;
        }
        SetCurrentObjectType(typename, packagename, arraydim);
        ReturnValuePseudoObject = CurrentObject;
        if (Compile != 0) && (ReturnValuePseudoObject != nil) {
            //TODO: Type check if function was forward decl.
            libgogo.AddParameters(ReturnValuePseudoObject, CurrentFunction); //Treat return value like an additional parameter at the end of the parameter list
            CurrentFunction.Len = CurrentFunction.Len - 1; //Don't count parameter as input parameter
        }
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
    if ed != nil {
        ed.ExpressionDepth = ed.ExpressionDepth - 1;
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
    for boolFlag = ParseSimpleExpressionOp(item, tempItem2, ed);
        boolFlag == 0;
        boolFlag = ParseSimpleExpressionOp(item, tempItem2, ed) {
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
func ParseSimpleExpressionOp(item1 *libgogo.Item, item2 *libgogo.Item, ed *ExpressionDescriptor) uint64 {
    var boolFlag uint64 = 1;
    PrintDebugString("Entering ParseSimpleExpressionOp()",1000);
    boolFlag = ParseUnaryArithOp(); // +,-
    if boolFlag == 0 {
        ParseTerm(item2, ed);
    } else {
        GetNextTokenSafe();
        if tok.id == TOKEN_REL_OR {
            // ||
            GenerateRelative(item1, TOKEN_REL_OR, ed);
            libgogo.Push(&Operators, TOKEN_REL_OR);
            ParseTerm(item1, ed);
            boolFlag = 0;
        } else {
            tok.nextToken = tok.id;
        }
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
    for boolFlag = ParseTermOp(item, tempItem2, ed);
        boolFlag == 0;
        boolFlag = ParseTermOp(item, tempItem2, ed) {
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
func ParseTermOp(item1 *libgogo.Item, item2 *libgogo.Item, ed *ExpressionDescriptor) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseTermOp()",1000);
    boolFlag = ParseBinaryArithOp(); // *,/
    if boolFlag == 0 {
        ParseFactor(item2, ed); //Arith items
    } else {
        GetNextTokenSafe();
        if tok.id == TOKEN_REL_AND {
            // &&
            GenerateRelative(item1, TOKEN_REL_AND, ed);
            libgogo.Push(&Operators, TOKEN_REL_AND);
            ParseFactor(item1, ed); //Rel items
            boolFlag = 0;
        } else {
            tok.nextToken = tok.id;
        }
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
    var tempType *libgogo.TypeDesc;
    var tempByteArray *libgogo.ObjectDesc;
    var tempString *libgogo.ObjectDesc;
    var startAddress uint64;
    var LHSItem *libgogo.Item;
    var RHSItem *libgogo.Item;
    var s string;

    GetNextTokenSafe();
    if (doneFlag == 1) && (tok.id == TOKEN_IDENTIFIER) {
        FindIdentifierAndParseSelector(item);
        doneFlag = 0;
    } 
    if (doneFlag == 1) && (tok.id == TOKEN_INTEGER) {
        if Compile != 0 {
            if (tok.intValue <= 255) { //Value fits into byte_t
                libgogo.SetItem(item, libgogo.MODE_CONST, byte_t, 0, tok.intValue, 0, 0); //Constant item
            } else { //Value does not fit into byte_t
                libgogo.SetItem(item, libgogo.MODE_CONST, uint64_t, 0, tok.intValue, 0, 0); //Constant item
            }
        }
        doneFlag = 0;
    }
    if (doneFlag == 1) && (tok.id == TOKEN_STRING) {
        if Compile != 0 {
            boolFlag = libgogo.StringLength(tok.strValue); //Compute length
            tempType = libgogo.NewType("byteArray", ".internal", 0, boolFlag + 1, byte_t); //Create byte array type of according length (including trailing 0) to be able to address the characters
            tempByteArray = libgogo.NewObject("tempByteArray", ".internal", libgogo.CLASS_VAR); //Create object of previously declared byte array type
            tempByteArray.ObjType = tempType;
            libgogo.AppendObject(tempByteArray, GlobalObjects); //Add byte array to global objects
            startAddress = libgogo.GetObjectOffset(tempByteArray, GlobalObjects); //Calculate buffer start address
            
            SwitchOutputToDataSegment(); //Place content of byte array in data segment
            s = "String buffer start ('";
            libgogo.StringAppend(&s, tok.strValue);
            libgogo.StringAppend(&s, "')");
            GenerateComment(s); //Output string in comment
            for doneFlag = 0; doneFlag < boolFlag; doneFlag = doneFlag + 1 { //Set values in data segment accordingly
                PutDataByte(startAddress + doneFlag, tok.strValue[doneFlag]);
            }
            PutDataByte(startAddress + doneFlag, 0); //Add trailing 0
            GenerateComment("String buffer end");
            SwitchOutputToCodeSegment(); //Reset to default output
            
            tempString = libgogo.NewObject("tempString", ".internal", libgogo.CLASS_VAR); //Create object for actual string
            tempString.ObjType = string_t;
            libgogo.AppendObject(tempString, GlobalObjects); //Add string to global objects
            
            SwitchOutputToInitCodeSegment(); //Initialize strings globally
            GenerateComment("Assign byte buffer to new string constant start");
            LHSItem = libgogo.NewItem();
            VariableObjectDescToItem(tempString, LHSItem, 1); //Global variable
            LHSItem.Itemtype = byte_t; //Set appropriate type for byte (array) pointer
            LHSItem.PtrType = 1;
            RHSItem = libgogo.NewItem();
            VariableObjectDescToItem(tempByteArray, RHSItem, 1); //Global variable
            RHSItem.Itemtype = byte_t; //Set appropriate type for byte (array) to be pointed at
            GenerateAssignment(LHSItem, RHSItem, 1); //tempString{first qword} = &tempByteArray
            GenerateComment("Assign byte buffer to new string constant end");
            GenerateComment("Assign string length to new string constant start");
            LHSItem = libgogo.NewItem();
            VariableObjectDescToItem(tempString, LHSItem, 1); //Global variable
            LHSItem.A = LHSItem.A + 8; //Access second qword of string to place length in there
            LHSItem.Itemtype = uint64_t; //Set appropriate type for string length
            RHSItem = libgogo.NewItem();
            libgogo.SetItem(RHSItem, libgogo.MODE_CONST, uint64_t, 0, boolFlag, 0, 0); //Constant item containing the length of the string
            GenerateAssignment(LHSItem, RHSItem, 0); //tempString{second qword} = boolFlag (string length)
            GenerateComment("Assign string length to new string constant end");
            SwitchOutputToCodeSegment(); //Reset to default output
            
            VariableObjectDescToItem(tempString, item, 1); //Global variable to return (reference to properly initialized, global string)
        }
        doneFlag = 0;
    }
    if (doneFlag) == 1 && (tok.id == TOKEN_LBRAC) {
        ParseExpression(item, ed);
        AssertNextTokenWeak(TOKEN_RBRAC);
        doneFlag = 0;
    }
    if (doneFlag == 1) && (tok.id == TOKEN_NOT) {
        if ed.Not == 1 {
            ed.Not = 0;
        } else {
            ed.Not = 1;
        }
        SwapExpressionBranches(ed);
        ParseFactor(item, ed);
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

func ParseSelector_FunctionCall(FunctionCalled *libgogo.TypeDesc) *libgogo.TypeDesc {
   var boolFlag uint64;
   PrintDebugString("Entering ParseSelector_FunctionCall()",1000);
    for boolFlag = ParseSelectorSub_FunctionCall(FunctionCalled);
        boolFlag == 0; 
        boolFlag = ParseSelectorSub_FunctionCall(FunctionCalled) {
    }
    PrintDebugString("Leaving ParseSelector_FunctionCall()",1000);
    if ReturnedFunction != nil {
        FunctionCalled = ReturnedFunction;
    }
    return FunctionCalled;
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

func ParseSelectorSub_FunctionCall(FunctionCalled *libgogo.TypeDesc) uint64 {
    var boolFlag uint64;
    var tempFcn *libgogo.TypeDesc;
    PrintDebugString("Entering ParseSelectorSub_FunctionCall()",1000);
    GetNextTokenSafe();
    ReturnedFunction = nil; //No new FunctionCalled pointer by default
    if tok.id == TOKEN_PT {
        AssertNextToken(TOKEN_IDENTIFIER);
        if Compile != 0 {
            boolFlag = libgogo.StringLength(FunctionCalled.Name);
            if boolFlag != 0 {
                SymbolTableError("Cannot apply a selector to", "a", "function, function", FunctionCalled.Name);
            } else {
                tempFcn = libgogo.GetType(tok.strValue, FunctionCalled.PackageName, GlobalFunctions, 1); //Check global functions
   			    if tempFcn == nil { //New forward declaration
   			        tempFcn = NewFunction(tok.strValue, FunctionCalled.PackageName, 1);
   			        ReturnedFunction = tempFcn; //Assign new tempFcn pointer to FunctionCalled (outside this function)
			    }
		        FunctionCalled.Name = tempFcn.Name;
		        FunctionCalled.Len = tempFcn.Len;
		        FunctionCalled.Fields = tempFcn.Fields;
		        FunctionCalled.ForwardDecl = tempFcn.ForwardDecl;
		        FunctionCalled.Base = tempFcn.Base;
            }
        }
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        if Compile != 0 {
            tempFcn = libgogo.GetType(FunctionCalled.PackageName, CurrentPackage, GlobalFunctions, 1); //Check global functions
            if tempFcn == nil { //New forward declaration
                tempFcn = NewFunction(FunctionCalled.PackageName, CurrentPackage, 1);
    	        ReturnedFunction = tempFcn; //Assign new tempFcn pointer to FunctionCalled (outside this function)
            }
            FunctionCalled.Name = tempFcn.Name;
            FunctionCalled.PackageName = tempFcn.PackageName;
            FunctionCalled.Len = tempFcn.Len;
            FunctionCalled.Fields = tempFcn.Fields;
            FunctionCalled.ForwardDecl = tempFcn.ForwardDecl;
            FunctionCalled.Base = tempFcn.Base;
	    }
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
        NewFunction(tok.strValue, CurrentPackage, 0);
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
    var exprIndicator uint64;
    var ReturnValueItem *libgogo.Item = nil;
    var ReturnExpression *libgogo.Item;
    PrintDebugString("Entering ParseFuncDecl()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_LCBRAC {
        InsideFunction = 1;
        if Compile != 0 {
            if ReturnValuePseudoObject != nil {
                ReturnValueItem = libgogo.NewItem();
                VariableObjectDescToItem(ReturnValuePseudoObject, ReturnValueItem, 2); //Treat return value like an additional parameter at the end of the parameter list
            }
            PrintFunctionStart(CurrentPackage, CurrentFunction.Name);
        }
        ParseVarDeclList();
        ParseStatementSequence(nil);
        GetNextTokenSafe();
        if tok.id == TOKEN_RETURN {
            if ReturnValuePseudoObject == nil {
                SymbolTableError("Cannot return a value when there is no return type", "", "in function", CurrentFunction.Name);
            }
            if Compile != 0 {
                GenerateComment("Return value assignment start");
                ReturnExpression = libgogo.NewItem();
                GenerateComment("Return expression load start");
                exprIndicator = ParseExpression(ReturnExpression,nil); //Parse return expression
                GenerateComment("Return expression load end");
                GenerateAssignment(ReturnValueItem, ReturnExpression, exprIndicator); //Return value = Return value expression
                GenerateComment("Return value assignment end");
            } else {
                ParseExpression(nil, nil);
            }
            AssertNextTokenWeak(TOKEN_SEMICOLON);
        } else {
            tok.nextToken = tok.id;
        }
        AssertNextToken(TOKEN_RCBRAC);
        InsideFunction = 0;
        if Compile != 0 {
            PrintFunctionEnd();
        }
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

func ParseStatementSequence(ed *ExpressionDescriptor) {
    var boolFlag uint64; 
    PrintDebugString("Entering ParseStatementSequence()",1000);
    for boolFlag = ParseStatement(ed);
        boolFlag == 0;
        boolFlag = ParseStatement(ed) { }
    PrintDebugString("Leaving ParseStatementSequence()",1000);
}

func ParseStatement(ed *ExpressionDescriptor) uint64 {
    var boolFlag uint64;
    var doneFlag uint64;
    var funcIndicator uint64;
    PrintDebugString("Entering ParseStatement()",1000);
    doneFlag = 1;

    GetNextTokenSafe();
    if (doneFlag == 1) && (tok.id == TOKEN_IDENTIFIER) {
        funcIndicator = IsFunction();
        if funcIndicator == 1 {
            ParseFunctionCallStatement(nil, 0);
        } else {
            ParseAssignment(1); //Trailing semicolon
        }
        doneFlag = 0;
    }
    
    if (doneFlag == 1) && (tok.id == TOKEN_IF) {
        tok.nextToken = tok.id;
        ParseIfStatement(ed);
        doneFlag = 0;
    }

    if (doneFlag == 1) && (tok.id == TOKEN_FOR) {
        tok.nextToken = tok.id;
        ParseForStatement();
        doneFlag = 0;
    }

    if (doneFlag == 1) && (tok.id == TOKEN_BREAK) {
        // simple break statement
        if (ed == nil) || (ed.ForEd == nil) {
            GenErrorWeak("Can only generate code for 'break' in 'for' clause.");
        } else {
            GenerateBreak(ed);
        }
        AssertNextToken(TOKEN_SEMICOLON);
        doneFlag = 0;
    }

    if (doneFlag == 1) && (tok.id == TOKEN_CONTINUE) {
        if (ed == nil) || (ed.ForEd == nil) {
            GenErrorWeak("Can only generate code for 'break' in 'for' clause.");
        } else {
            GenerateContinue(ed);
        }
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
                RHSItem = ParseFunctionCallStatement(LHSItem.Itemtype, LHSItem.PtrType);
                if Compile != 0 {
                    if RHSItem == nil {
                        SymbolTableError("Function has", "no", "return value", "");
                    }
                    GenerateAssignment(LHSItem, RHSItem, 0); //LHS = RHS
                }
            } else { //Expression starting with an identifier
                if Compile != 0 {
                    RHSItem = libgogo.NewItem();
                    GenerateComment("Assignment RHS load start");
                    exprIndicator = ParseExpression(RHSItem, nil); //Parse RHS
                    GenerateComment("Assignment RHS load end");
                    GenerateAssignment(LHSItem, RHSItem, exprIndicator); //LHS = RHS
                } else {
                    ParseExpression(nil, nil);
                }
            }
        } else { //Expression
            tok.nextToken = tok.id;
            if Compile != 0 {
                RHSItem = libgogo.NewItem();
                GenerateComment("Assignment RHS load start");
                exprIndicator = ParseExpression(RHSItem, nil); //Parse RHS
                GenerateComment("Assignment RHS load end");
                GenerateAssignment(LHSItem, RHSItem, exprIndicator); //LHS = RHS
            } else {
                ParseExpression(nil, nil);
            }
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

func ParseFunctionCall(FunctionCalled *libgogo.TypeDesc) *libgogo.Item {
    var paramCount uint64 = 0;
    var FullFunctionName string;
    var tempString string;
    var TotalParameterSize uint64;
    var TotalLocalVariableSize uint64;
    var ReturnObject *libgogo.ObjectDesc;
    var ReturnItem *libgogo.Item;
    
    PrintDebugString("Entering ParseFunctionCall()",1000);
    AssertNextToken(TOKEN_LBRAC);
    GetNextTokenSafe();
    if Compile != 0 {
        SaveUsedRegisters();
        TotalParameterSize = libgogo.GetAlignedObjectListSize(FunctionCalled.Fields); //Get total size of parameters of function called
        TotalLocalVariableSize = libgogo.GetAlignedObjectListSize(LocalObjects); //Get total size of local variables of current function
    }
    if tok.id != TOKEN_RBRAC {
        tok.nextToken = tok.id;
        if FunctionCalled.ForwardDecl == 1 {
            paramCount = ParseExpressionList(FunctionCalled, TotalLocalVariableSize);
        } else {
            paramCount = ParseExpressionList(FunctionCalled, TotalParameterSize + TotalLocalVariableSize);
        }
        AssertNextTokenWeak(TOKEN_RBRAC);
    }
    if Compile != 0 {
        if FunctionCalled.Len > paramCount { //Compare number of actual parameters
            FullFunctionName = "";
            libgogo.StringAppend(&FullFunctionName, FunctionCalled.PackageName);
            libgogo.CharAppend(&FullFunctionName, '.');
            libgogo.StringAppend(&FullFunctionName, FunctionCalled.Name);
            tempString = libgogo.IntToString(FunctionCalled.Len);
            SymbolTableError("Expecting", tempString, "parameters (more than the actual ones) for function", FullFunctionName);
        }
        TotalParameterSize = libgogo.GetAlignedObjectListSize(FunctionCalled.Fields); //Recalculate sizes as function may have been forward declared
        if FunctionCalled.ForwardDecl == 1 {
            PrintFunctionCall(FunctionCalled.PackageName, FunctionCalled.Name, TotalLocalVariableSize, 1);
        } else {
            PrintFunctionCall(FunctionCalled.PackageName, FunctionCalled.Name, TotalParameterSize + TotalLocalVariableSize, 0);
        }
        RestoreUsedRegisters();
    }
    PrintDebugString("Leaving ParseFunctionCall()",1000);
    ReturnObject = libgogo.GetObject("return value", "", FunctionCalled.Fields); //Find return value
    if ReturnObject == nil {
        ReturnItem = nil;
    } else {
        if FunctionCalled.ForwardDecl == 1 {
            ReturnItem = ObjectToStackParameter(ReturnObject, FunctionCalled, TotalLocalVariableSize);
        } else {
            ReturnItem = ObjectToStackParameter(ReturnObject, FunctionCalled, TotalParameterSize + TotalLocalVariableSize);
        }
    }
    return ReturnItem;
}

func ParseExpressionList(FunctionCalled *libgogo.TypeDesc, TotalParameterSize uint64) uint64 {
    var ed ExpressionDescriptor;
    var boolFlag uint64;
    var paramCount uint64 = 1; //There has to be at least one expression
    var ExprItem *libgogo.Item;
    var ParameterLHSObject *libgogo.ObjectDesc;
    var Parameter *libgogo.Item;
    var TempObject *libgogo.ObjectDesc;
    var tempString string;
    var FullFunctionName string;
    
    PrintDebugString("Entering ParseExpressionList()",1000);
    if Compile != 0 {
        GenerateComment("First parameter expression start");
        ExprItem = libgogo.NewItem();
        GenerateComment("First parameter expression load start");
        boolFlag = ParseExpression(ExprItem, &ed);
        GenerateComment("First parameter expression load end");
        if (FunctionCalled.ForwardDecl == 1) && (FunctionCalled.Base == nil) { //Create artificial parameter from expression (based on the latter's type) if the function is called the first time without being declared
            TempObject = libgogo.NewObject("Artificial parameter", "", libgogo.CLASS_PARAMETER);
            TempObject.ObjType = ExprItem.Itemtype; //Derive type from expression
            TempObject.PtrType = ExprItem.PtrType; //Derive pointer type from expression
            if boolFlag != 0 { //& in expression forces pointer type
                if TempObject.PtrType == 0 {
                    TempObject.PtrType = 1;
                } else {
                    SymbolTableError("& operator on pointer type not allowed,", "", "type: pointer to", ExprItem.Itemtype.Name);
                }
            }
            libgogo.AddParameters(TempObject, FunctionCalled); //Add a new, artificial parameter
        }
        if FunctionCalled.Len == 0 { //Check if function expects parameters
            FullFunctionName = "";
            libgogo.StringAppend(&FullFunctionName, FunctionCalled.PackageName);
            libgogo.CharAppend(&FullFunctionName, '.');
            libgogo.StringAppend(&FullFunctionName, FunctionCalled.Name);
            SymbolTableError("Function expects", "no", "parameters:", FullFunctionName);
        }
        ParameterLHSObject = libgogo.GetParameterAt(paramCount, FunctionCalled);
        Parameter = ObjectToStackParameter(ParameterLHSObject, FunctionCalled, TotalParameterSize);
        GenerateAssignment(Parameter, ExprItem, boolFlag); //Assignment
        GenerateComment("First parameter expression end");
    } else  {
        ParseExpression(nil, &ed);
    }
    if Compile != 0 {
        for boolFlag = 0; boolFlag == 0; paramCount = paramCount + 1 {
            if FunctionCalled.Len < paramCount { //Compare number of actual parameters
                FullFunctionName = "";
                libgogo.StringAppend(&FullFunctionName, FunctionCalled.PackageName);
                libgogo.CharAppend(&FullFunctionName, '.');
                libgogo.StringAppend(&FullFunctionName, FunctionCalled.Name);
                tempString = libgogo.IntToString(FunctionCalled.Len);
                SymbolTableError("Expecting", tempString, "parameters (less than the actual ones) for function", FullFunctionName);
            }
            boolFlag = ParseExpressionListSub(FunctionCalled, TotalParameterSize, paramCount + 1);
        }
        if paramCount != 0 { //Correct param count if for loop has been entered
            paramCount = paramCount - 1;
        }
    } else {
        for boolFlag = ParseExpressionListSub(FunctionCalled, TotalParameterSize, paramCount + 1);
            boolFlag == 0;
            boolFlag = ParseExpressionListSub(FunctionCalled, TotalParameterSize, paramCount + 1) { }
    }
    PrintDebugString("Leaving ParseExpressionList()",1000);
    return paramCount;
}

func ParseExpressionListSub(FunctionCalled *libgogo.TypeDesc, TotalParameterSize uint64, ParameterIndex uint64) uint64 {
    var boolFlag uint64;
    var ed ExpressionDescriptor;
    var ExprItem *libgogo.Item;
    var ParameterLHSObject *libgogo.ObjectDesc;
    var Parameter *libgogo.Item;
    var TempObject *libgogo.ObjectDesc;
    
    PrintDebugString("Entering ParseExpressionListSub()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_COLON {
        if Compile != 0 {
            GenerateComment("Subsequent parameter expression start");
            ExprItem = libgogo.NewItem();
            GenerateComment("Subsequent parameter expression load start");
            boolFlag = ParseExpression(ExprItem, &ed);
            GenerateComment("Subsequent parameter expression load end");
            if (FunctionCalled.ForwardDecl == 1) && (FunctionCalled.Base == nil) { //Create artificial parameter from expression (based on the latter's type) if the function is called the first time without being declared
                TempObject = libgogo.NewObject("Artificial parameter", "", libgogo.CLASS_PARAMETER);
                TempObject.ObjType = ExprItem.Itemtype; //Derive type from expression
                TempObject.PtrType = ExprItem.PtrType; //Derive pointer type from expression
                if boolFlag != 0 { //& in expression forces pointer type
                    if TempObject.PtrType == 0 {
                        TempObject.PtrType = 1;
                    } else {
                        SymbolTableError("& operator on pointer type not allowed,", "", "type: pointer to", ExprItem.Itemtype.Name);
                    }
                }
                libgogo.AddParameters(TempObject, FunctionCalled); //Add a new, artificial parameter
            }
            ParameterLHSObject = libgogo.GetParameterAt(ParameterIndex, FunctionCalled);
            Parameter = ObjectToStackParameter(ParameterLHSObject, FunctionCalled, TotalParameterSize);
            GenerateAssignment(Parameter, ExprItem, boolFlag); //Assignment
            GenerateComment("Subsequent parameter expression end");
        } else {
            ParseExpression(nil, &ed);
        }
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseExpressionListSub()",1000);
    return boolFlag;   
}

func ParseFunctionCallStatement(ForwardDeclExpectedReturnType *libgogo.TypeDesc, ForwardDeclExpectedReturnPtrType uint64) *libgogo.Item {
    var FunctionCalled *libgogo.TypeDesc;
    var ReturnValue *libgogo.Item;
    var TotalLocalVariableSize uint64;
    PrintDebugString("Entering ParseFunctionCallStatement()",1000);
    AssertNextToken(TOKEN_IDENTIFIER);
    FunctionCalled = libgogo.NewType("", "", 0, 0, nil);
    FunctionCalled = FindIdentifierAndParseSelector_FunctionCall(FunctionCalled);
    ReturnValue = ParseFunctionCall(FunctionCalled);
    if Compile != 0 {
        if (FunctionCalled.ForwardDecl == 1) && (FunctionCalled.Base == nil) { //Create artifical return value if function is called the first time
            if ForwardDeclExpectedReturnType != nil { //Return type expected
                CurrentObject = libgogo.NewObject("return value", "", libgogo.CLASS_PARAMETER); //Create artificial return value
                CurrentObject.ObjType = ForwardDeclExpectedReturnType;
                CurrentObject.PtrType = ForwardDeclExpectedReturnPtrType;
                libgogo.AddParameters(CurrentObject, FunctionCalled); //Add a new, artificial return value
                FunctionCalled.Len = FunctionCalled.Len - 1; //Don't count parameter as input parameter
                TotalLocalVariableSize = libgogo.GetAlignedObjectListSize(LocalObjects); //Take local variable size into consideration for offset below
                ReturnValue = ObjectToStackParameter(CurrentObject, FunctionCalled, TotalLocalVariableSize);
            } else { //No return type expected
                ReturnValue = nil;
            }
        }
        FunctionCalled.Base = FunctionCalled; //Abuse Base field to indicate that the function has been called at least once
    }
    PrintDebugString("Leaving ParseFunctionCallStatement()",1000);
    return ReturnValue;
}

//
//
//
func ParseForStatement() {
    var ed ExpressionDescriptor;
    var item *libgogo.Item;
    var expr uint64 = 0;
    var postassign uint64 = 0;
    PrintDebugString("Entering ParseForStatement()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_FOR {
        GenerateComment("For start");
        SetExpressionDescriptor(&ed, "FOR_"); // Set the required descriptor parameters
        ed.ForEd = &ed;
        GetNextTokenSafe();

        if tok.id == TOKEN_SEMICOLON {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            GenerateComment("For [initial assignment] start");
            ParseAssignment(0); //No semicolon
            GenerateComment("For [initial assignment] end");
        }
        
        AssertNextToken(TOKEN_SEMICOLON);

        GenerateExpressionStart(&ed);

        GetNextTokenSafe();
        if tok.id == TOKEN_SEMICOLON {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            item = libgogo.NewItem();
            GenerateComment("For [expression] start");

            ParseExpression(item, &ed);
            GenerateForStart(item, &ed);
            GenerateComment("For [expression] end");
            expr = 1;
        }


        AssertNextToken(TOKEN_SEMICOLON);

        GetNextTokenSafe();
        if tok.id == TOKEN_LCBRAC {
            tok.nextToken = tok.id;
        } else {
            GenerateForBodyJump(&ed);
            GenerateForBodyExtended(&ed);
            tok.nextToken = tok.id;
            GenerateComment("For [post assignment] start");
            ParseAssignment(0); //No semicolon
            GenerateComment("For [post assignment] end");
            postassign = 1;
            ed.ForPost = 1;
        }

        GenerateForBody(&ed, postassign, expr);
        AssertNextToken(TOKEN_LCBRAC);        
        ParseStatementSequence(&ed);
        AssertNextToken(TOKEN_RCBRAC);
        GenerateForEnd(&ed, postassign);
    } else {
        tok.nextToken = tok.id;
    }   
    GenerateComment("For end");
    PrintDebugString("Leaving ParseForStatement()",1000);
}

//
// Parses: "if" expression "{" stmt_sequence [ "}" "else" else_stmt ].
// Represents an if statement. Else is parsed in a separate function.
//
func ParseIfStatement(oldEd *ExpressionDescriptor) {
    var item *libgogo.Item;
    var ed ExpressionDescriptor;
    PrintDebugString("Entering ParseIfStatement()",1000);
    GenerateComment("If start");
    SetExpressionDescriptor(&ed, "IF_"); // Set the required descriptor parameters
    if oldEd != nil {
        ed.ForEd = oldEd;
    }
    GetNextTokenSafe();
    if tok.id == TOKEN_IF {
        item = libgogo.NewItem();
        ParseExpression(item, &ed);
        GenerateIfStart(item, &ed);
        AssertNextToken(TOKEN_LCBRAC);
        ParseStatementSequence(&ed);
        AssertNextToken(TOKEN_RCBRAC);
        
        GetNextTokenSafe();
        if tok.id == TOKEN_ELSE {
            GenerateElseStart(&ed);
            ParseElseStatement(&ed);
            GenerateElseEnd(&ed);
        } else {
            GenerateIfEnd(&ed);
            GenerateComment("If end");
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
func ParseElseStatement(ed *ExpressionDescriptor) {
    PrintDebugString("Entering ParseElseStatement()",1000);
    GenerateComment("Else start");
    AssertNextTokenWeak(TOKEN_LCBRAC);
    ParseStatementSequence(ed);
    AssertNextToken(TOKEN_RCBRAC);
    GenerateComment("Else end");
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
		if (boolFlag == 0) && (CurrentFunction != nil) {
		    boolFlag = libgogo.FindPackageName(tok.strValue, CurrentFunction.Fields); //Check local parameters
		}
		if boolFlag == 0 { //Token is not package name, but identifier
			tempObject = libgogo.GetObject(tok.strValue, CurrentPackage, LocalObjects); //Check local objects
			tempList = LocalObjects;
			if tempObject == nil {
                if CurrentFunction != nil {
           			tempObject = libgogo.GetObject(tok.strValue, CurrentPackage, CurrentFunction.Fields); //Check local parameters
        			tempList = CurrentFunction.Fields;
    			}
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

func FindIdentifierAndParseSelector_FunctionCall(FunctionCalled *libgogo.TypeDesc) *libgogo.TypeDesc {
    var boolFlag uint64;
    var tempFcn *libgogo.TypeDesc;
    if Compile != 0 {
		//Token can be package name
		boolFlag = libgogo.FindTypePackageName(tok.strValue, GlobalFunctions); //Check global functions
		if boolFlag == 0 { //Token is not package name, but identifier
			tempFcn = libgogo.GetType(tok.strValue, CurrentPackage, GlobalFunctions, 1); //Check global functions
            if tempFcn == nil { //New forward declaration
                FunctionCalled.PackageName = tok.strValue; //Set package name
                tempFcn = ParseSelector_FunctionCall(FunctionCalled); //Parse selector for an undefined function in the given package
                boolFlag = libgogo.StringLength(tempFcn.Name);
                if boolFlag == 0 {
		            SymbolTableError("Cannot use package", "", "as a function:", tempFcn.PackageName);
		        }
            } //else: tempFcn is already the return value
        } else { //Token is package name
            FunctionCalled.PackageName = tok.strValue; //Set package name
            tempFcn = ParseSelector_FunctionCall(FunctionCalled); //Parse selector for an undefined function in the given package
            boolFlag = libgogo.StringLength(tempFcn.Name);
            if boolFlag == 0 {
		        SymbolTableError("Cannot use package", "", "as a function:", tempFcn.PackageName);
		    }
        }
    } else {
        tempFcn = ParseSelector_FunctionCall(FunctionCalled);
    }
    return tempFcn;
}
