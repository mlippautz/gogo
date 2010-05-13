// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

var tok Token;

var maxDepth uint64 = 10;
var curDepth uint64 = 1;
var compile uint64 = 1;

//
// List of global objects and declared types
//
var GlobalObjects *libgogo.ObjectDesc = nil;
var GlobalTypes *libgogo.TypeDesc = nil;

//
// List of function-local objects
//
var LocalObjects *libgogo.ObjectDesc = nil;

var CurrentType *libgogo.TypeDesc;
var CurrentObject *libgogo.ObjectDesc;

var InsideFunction uint64 = 0;
var InsideStructDecl uint64 = 0;
var InsideFunctionVarDecl uint64 = 0;

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
    var tempType *libgogo.TypeDesc;
    var dontAddType uint64 = 0;
    var temp byte;
    PrintDebugString("Entering ParseStructDecl()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_TYPE {
        AssertNextToken(TOKEN_IDENTIFIER);
        tempType = libgogo.GetType(tok.strValue, CurrentPackage, GlobalTypes, 1);
        if  tempType != nil { //Check for duplicates
            temp = libgogo.IsForwardDecl(tempType);
            if temp != 0 { //Separate handling of forward declarations => unset forward declaration flag
                libgogo.UnsetForwardDecl(tempType);
                CurrentType = tempType;
                dontAddType = 1;
            } else { //Real duplicate
                SymbolTableError("duplicate type", tok.strValue, "in package", CurrentPackage);
                CurrentType = tempType; //For stability/sanity of parser when continuing
                dontAddType = 1; //For stability/sanity of parser when continuing
            }
        } else {
            CurrentType = libgogo.NewType(tok.strValue, CurrentPackage, 0, 0, nil);
        }
        // identifier of struct in tok.strValue
        AssertNextToken(TOKEN_STRUCT);
        AssertNextToken(TOKEN_LCBRAC);
        InsideStructDecl = 1;
        ParseStructVarDeclList();
        InsideStructDecl = 0;
        AssertNextTokenWeak(TOKEN_RCBRAC);
        AssertNextTokenWeak(TOKEN_SEMICOLON);
        if dontAddType == 0 {
            GlobalTypes = libgogo.AppendType(CurrentType, GlobalTypes);
        }
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
    var temp uint64;
    PrintDebugString("Entering ParseStructVarDecl()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_IDENTIFIER {
        temp = libgogo.HasField(tok.strValue, CurrentType);
        if temp != 0 {
            SymbolTableError("duplicate", "", "field", tok.strValue);
        } else {
            CurrentObject = libgogo.NewObject(tok.strValue, "", libgogo.CLASS_FIELD); //A field has no package name
            libgogo.AddFields(CurrentObject, CurrentType);
        }
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
    var basetype *libgogo.TypeDesc;
    var temptype *libgogo.TypeDesc;
    var boolFlag uint64;
    var packagename string = CurrentPackage;
    var tempstr string = "";
    var typename string;
    var CurrentTypeName string;
    var temp byte;

    PrintDebugString("Entering ParseType()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_LSBRAC {   
        AssertNextToken(TOKEN_INTEGER);
        // value of integer in tok.intValue
        if InsideFunctionVarDecl == 0 {
            arraydim = tok.intValue;
        }
        AssertNextToken(TOKEN_RSBRAC);
    } else {
        tok.nextToken = tok.id;
    }
    GetNextTokenSafe();
    if tok.id != TOKEN_ARITH_MUL {
        tok.nextToken = tok.id;
    } else {
        if InsideFunctionVarDecl == 0 {
            libgogo.FlagObjectTypeAsPointer(CurrentObject); //Type is pointer
        }
    }
    AssertNextToken(TOKEN_IDENTIFIER);
    // typename in tok.strValue
    typename = tok.strValue;
    boolFlag = ParseSimpleSelector();
    if boolFlag == 0 { //Take selector into consideration if there is one
        packagename = typename; //Previously read namespace is actually the namespace
        typename = tok.strValue;
    }

    if InsideFunctionVarDecl == 0 {
        basetype = libgogo.GetType(typename, packagename, GlobalTypes, 0);
        if basetype == nil {
            if InsideStructDecl == 1 {
                CurrentTypeName = libgogo.GetTypeName(CurrentType);
                boolFlag = libgogo.StringCompare(typename, CurrentTypeName);
                if boolFlag == 0 {
                    temp = libgogo.IsPointerType(CurrentObject);
                    if temp == 1 { //Allow pointer to own type
                        basetype = CurrentType;
                    } else {
                        SymbolTableError("A type cannot contain itself,", "", "type", typename);
                    }
                } else { //Forward declaration
                    basetype = libgogo.NewType(typename, packagename, 1, 0, nil);
                    GlobalTypes = libgogo.AppendType(basetype, GlobalTypes); //Add forward declared type to global list
                }
            } else {
                libgogo.StringAppend(&tempstr, packagename);
                libgogo.CharAppend(&tempstr, '.');
                libgogo.StringAppend(&tempstr, typename);
                SymbolTableError("Unknown", "", "type", tempstr);
            }
        }
        if arraydim == 0 { //No array
            libgogo.SetObjType(CurrentObject, basetype);
        } else { //Array
            if basetype != nil {
                CurrentTypeName = libgogo.GetTypeName(basetype);
                libgogo.StringAppend(&tempstr, CurrentTypeName);
            }
            libgogo.StringAppend(&tempstr, "Array");
            CurrentTypeName = libgogo.IntToString(arraydim); //Reuse CurrentTypeName as string representation of arraydim
            libgogo.StringAppend(&tempstr, CurrentTypeName);
            temptype = libgogo.NewType(tempstr, packagename, 0, arraydim, basetype);
            libgogo.SetObjType(CurrentObject, temptype);
        }
    }
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
    var TempObject *libgogo.ObjectDesc;
    PrintDebugString("Entering ParseVarDecl()",1000);
    boolFlag = LookAheadAndCheck(TOKEN_VAR);
    if boolFlag == 0 {
        AssertNextToken(TOKEN_VAR);
        AssertNextToken(TOKEN_IDENTIFIER);
        CurrentObject = libgogo.NewObject(tok.strValue, CurrentPackage, libgogo.CLASS_VAR);
        if InsideFunction == 0 { //Global objects
            TempObject = libgogo.GetObject(tok.strValue, CurrentPackage, GlobalObjects);
            if TempObject != nil {
                SymbolTableError("duplicate", "global", "identifier", tok.strValue);
            }
            GlobalObjects = libgogo.AppendObject(CurrentObject, GlobalObjects);
        } else { //Function-local objects
            TempObject = libgogo.GetObject(tok.strValue, CurrentPackage, LocalObjects);
            if TempObject != nil {
                SymbolTableError("duplicate", "local", "identifier", tok.strValue);
            }
            LocalObjects = libgogo.AppendObject(CurrentObject, LocalObjects);
        }
        // variable name in tok.strValue
        ParseType();

        GetNextTokenSafe();
        if tok.id == TOKEN_ASSIGN {
            ParseExpression();        
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
func ParseExpression() {
    PrintDebugString("Entering ParseExpression()",1000);
    IncAndCheckDepth();
    ParseSimpleExpression();    
    ParseCmpOp();
    DecDepth();
    PrintDebugString("Leaving ParseExpression()",1000);
}

//
//
//
func ParseCmpOp() {
    PrintDebugString("Entering ParseCmpOp()",1000);
    GetNextTokenSafe();
    if (tok.id == TOKEN_EQUALS) || (tok.id == TOKEN_NOTEQUAL) || 
        (tok.id == TOKEN_REL_LT) || (tok.id == TOKEN_REL_LTOE) || 
        (tok.id == TOKEN_REL_GT) || (tok.id == TOKEN_REL_GTOE) {
        ParseSimpleExpression();
    } else {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseCmpOp()",1000);
}

//
//
//
func ParseSimpleExpression() {
    var boolFlag uint64;
    PrintDebugString("Entering ParseSimpleExpression()",1000);
    ParseUnaryArithOp();
    ParseTerm();
    for boolFlag = ParseSimpleExpressionOp();
        boolFlag == 0;
        boolFlag = ParseSimpleExpressionOp() { }
    PrintDebugString("Leaving ParseSimpleExpression()",1000);
}

//
//
//
func ParseSimpleExpressionOp() uint64 {
    var boolFlag uint64 = 1;
    PrintDebugString("Entering ParseSimpleExpressionOp()",1000);
    boolFlag = ParseUnaryArithOp(); // +,-
    if boolFlag != 0 {
        GetNextTokenSafe();
        if tok.id == TOKEN_REL_OR {
            // ||
            boolFlag = 0;
        } 
    }
    if boolFlag == 0 {
        ParseTerm();
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
    if tok.id == TOKEN_ARITH_PLUS {
        // +
        boolFlag = 0;
    }
    if tok.id == TOKEN_ARITH_MINUS {
        // -
        boolFlag = 0;
    }
    if boolFlag != 0 {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseUnaryArithOp()",1000);
    return boolFlag;   
}

//
//
//
func ParseTerm() {
    var boolFlag uint64;
    PrintDebugString("Entering ParseTerm()",1000);
    ParseFactor();
    for boolFlag = ParseTermOp();
        boolFlag == 0;
        boolFlag = ParseTermOp() { }      
    PrintDebugString("Leaving ParseTerm()",1000);
}

//
//
//
func ParseTermOp() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseTermOp()",1000);
    boolFlag = ParseBinaryArithOp(); // *,/
    if boolFlag != 0 {
        GetNextTokenSafe();
        if tok.id == TOKEN_REL_AND {
            // &&
            boolFlag = 0;
        }
    }
    if boolFlag == 0 {
        ParseFactor();
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
    if tok.id == TOKEN_ARITH_MUL {
        // *
        boolFlag = 0;
    }
    if tok.id == TOKEN_ARITH_DIV {
        // /
        boolFlag = 0;
    }
    if boolFlag != 0 {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseBinaryArithOp()",1000);
    return boolFlag; 
}

//
//
//
func ParseFactor() uint64 {
    var doneFlag uint64 = 1;
    var boolFlag uint64;
    var es [2]uint64;

    GetNextTokenSafe();
    if (doneFlag == 1) && (tok.id == TOKEN_OP_ADR) {
        AssertNextToken(TOKEN_IDENTIFIER);
        ParseSelector();
        doneFlag = 0;
    }
    if (doneFlag == 1) && (tok.id == TOKEN_IDENTIFIER) {
        ParseSelector();
        doneFlag = 0;
    } 
    if (doneFlag == 1) && (tok.id == TOKEN_INTEGER) {
        doneFlag = 0;
    }
    if (doneFlag) == 1 && (tok.id == TOKEN_STRING) {
        doneFlag = 0;
    }
    if (doneFlag) == 1 && (tok.id == TOKEN_LBRAC) {
        ParseExpression();
        AssertNextTokenWeak(TOKEN_RBRAC);
        doneFlag = 0;
    }
    if (doneFlag == 1) && (tok.id == TOKEN_NOT) {
        ParseFactor();
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
func ParseSelector() {
    var boolFlag uint64;
    PrintDebugString("Entering ParseSelector()",1000);
    for boolFlag = ParseSelectorSub();
        boolFlag == 0; 
        boolFlag = ParseSelectorSub() {
    }
    PrintDebugString("Leaving ParseSelector()",1000);
}

//
//
//
func ParseSelectorSub() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseSelectorSub()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_PT {
        AssertNextToken(TOKEN_IDENTIFIER);
        // value in tok.strValue
        boolFlag = 0;
    } else {
        if tok.id == TOKEN_LSBRAC {
            GetNextTokenSafe();
            if tok.id == TOKEN_INTEGER {
                
            } else {
                if tok.id == TOKEN_IDENTIFIER {  
                    ParseSelector();
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
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseFuncDeclRaw()",1000);
    return boolFlag;
}

func ParseFuncDecl() uint64 {
    var boolFlag uint64;
    var temp uint64;
    PrintDebugString("Entering ParseFuncDecl()",1000);
    InsideFunction = 1;
    GetNextTokenSafe();
    if tok.id == TOKEN_LCBRAC {
        ParseVarDeclList();
        ParseStatementSequence();
        GetNextTokenSafe();
        if tok.id == TOKEN_RETURN {
            ParseExpression();
            AssertNextTokenWeak(TOKEN_SEMICOLON);
        } else {
            tok.nextToken = tok.id;
        }
        AssertNextToken(TOKEN_RCBRAC);
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    InsideFunction = 0;
    temp = CheckDebugLevel(100);
    if temp == 1 { //Function-local symbol table
        libgogo.PrintString("\nFunction-local symbol table until line ");
        libgogo.PrintNumber(fileInfo[curFileIndex].lineCounter);
        libgogo.PrintString(" of ");
        libgogo.PrintString(fileInfo[curFileIndex].filename);
        libgogo.PrintString(":\n----------------------------------------------------------------------------\n");
        libgogo.PrintObjects(LocalObjects);
    }
    LocalObjects = nil; //Delete local objects
    PrintDebugString("Leaving ParseFuncDecl()",1000);
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
        GetNextTokenSafe();
        if tok.id == TOKEN_ARITH_MUL {

        } else {
            tok.nextToken = tok.id;
        }
        InsideFunctionVarDecl = 1;
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
            ParseAssignment();
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

func ParseAssignment() uint64 {
    var boolFlag uint64;
    var funcIndicator uint64;
    PrintDebugString("Entering ParseAssignment()",1000);
    AssertNextToken(TOKEN_IDENTIFIER);
    ParseSelector();
    GetNextTokenSafe();
    if tok.id == TOKEN_ASSIGN {
        GetNextTokenSafe();
        if tok.id == TOKEN_IDENTIFIER {
            funcIndicator = IsFunction();
            if funcIndicator == 1 {
                ParseFunctionCallStatement();
            } else {
                ParseExpression();
            }
        } else {
            tok.nextToken = tok.id;
            ParseExpression();
        }
        AssertNextTokenWeak(TOKEN_SEMICOLON);
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseAssignment()",1000);
    return boolFlag;
}

func ParseAssignmentWithoutSC() uint64 {
    var boolFlag uint64;
    var funcIndicator uint64;
    PrintDebugString("Entering ParseAssignmentWithoutSC()",1000);
    AssertNextToken(TOKEN_IDENTIFIER);
    ParseSelector();
    GetNextTokenSafe();
    if tok.id == TOKEN_ASSIGN {
        GetNextTokenSafe();
        if tok.id == TOKEN_IDENTIFIER {
            funcIndicator = IsFunction();
            if funcIndicator == 1 {
                ParseFunctionCallStatement();
            } else {
                ParseExpression();
            }
        } else {
            tok.nextToken = tok.id;
            ParseExpression();
        }
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseAssignmentWithoutSC()",1000);
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
    var boolFlag uint64;
    PrintDebugString("Entering ParseExpressionList()",1000);
    ParseExpression();
    for boolFlag = ParseExpressionListSub();
        boolFlag == 0;
        boolFlag = ParseExpressionListSub() { }   
    PrintDebugString("Leaving ParseExpressionList()",1000);
}

func ParseExpressionListSub() uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseExpressionListSub()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_COLON {
        ParseExpression();
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
    AssertNextToken(TOKEN_IDENTIFIER);
    ParseSelector();    
    ParseFunctionCall();                
    PrintDebugString("Leaving ParseFunctionCallStatement()",1000);
}

func ParseForStatement() {
    PrintDebugString("Entering ParseForStatement()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_FOR {
        GetNextTokenSafe();

        if tok.id == TOKEN_SEMICOLON {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            ParseAssignmentWithoutSC();
        }
        
        AssertNextToken(TOKEN_SEMICOLON);

        GetNextTokenSafe();
        if tok.id == TOKEN_SEMICOLON {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            ParseExpression();
        }

        AssertNextToken(TOKEN_SEMICOLON);

        GetNextTokenSafe();
        if tok.id == TOKEN_LCBRAC {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            ParseAssignmentWithoutSC();
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
    PrintDebugString("Entering ParseIfStatement()",1000);
    GetNextTokenSafe();
    if tok.id == TOKEN_IF {
        ParseExpression();
        AssertNextToken(TOKEN_LCBRAC);
        ParseStatementSequence();
        AssertNextToken(TOKEN_RCBRAC);

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
    compile = 0; // stop producing code
    for ;(tok.id != TOKEN_FUNC) && (tok.id != TOKEN_EOS); {
        GetNextTokenSafe();
    }
    if tok.id == TOKEN_FUNC {
        tok.nextToken = tok.id;
        ParseFuncDeclList();       
    }
    libgogo.Exit(3); // Exit with an error
}
