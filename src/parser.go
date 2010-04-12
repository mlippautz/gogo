// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

//
// Main parsing function. Corresponds to the EBNF main structure called 
// go_program.
//
func Parse() {
    var tok Token;
    tok.id = 0;
    tok.nextChar = 0;
    tok.nextToken = 0;    

    ParsePackageStatement(&tok);    
    ParseImportStatementList(&tok);
    ParseStructDeclList(&tok);
    ParseVarDeclList(&tok);
    ParseFuncDeclList(&tok);

    AssertNextToken(&tok, TOKEN_EOS);
}

//
// Parses: package identifier
// This is enforced by the go language as first statement in a source file.
//
func ParsePackageStatement(tok *Token) {    
    PrintDebugString("Entering ParsePackageStatement()",1000);
    AssertNextToken(tok, TOKEN_PACKAGE);
    AssertNextToken(tok, TOKEN_IDENTIFIER);
    // package ok, value in tok.strValue
    PrintDebugString("Leaving ParsePackageStatement()",1000);
}

//
// Parses: { import_stmt }
// Function parsing the whole import block (which is optional) of a go program
//
func ParseImportStatementList(tok *Token) {
    var validImport uint64;
    PrintDebugString("Entering ParseImportStatementList()",1000);
    for validImport = ParseImportStatement(tok);
        validImport == 0;
        validImport = ParseImportStatement(tok) { }
    PrintDebugString("Leaving ParseImportStatementList()",1000);
}

//
// Parses: "import" string
// This function parses a single import line.
// Returning 0 if import statement is valid, 1 otherwise.
//
func ParseImportStatement(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseImportStatement()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_IMPORT {
        AssertNextToken(tok, TOKEN_STRING);
        // import ok, value in tok.strValue
        boolFlag = 0;
    } else {
        boolFlag = 1;
        SyncToken(tok);
    }    
    PrintDebugString("Leaving ParseImportStatement()",1000);
    return boolFlag;
}

//
// Parses: { struct_decl }
// A list of struct declarations.
//
func ParseStructDeclList(tok *Token) {
    var boolFlag uint64;
    PrintDebugString("Entering ParseStructDeclList()",1000);
    for boolFlag = ParseStructDecl(tok);
        boolFlag == 0;
        boolFlag = ParseStructDecl(tok) { }
    PrintDebugString("Leaving ParseStructDeclList()",1000);
}

//
// Parses: "type" identifier "struct" "{" struct_var_decl_list "}" ";"
// This is basically the skeleton of a struct.
//
func ParseStructDecl(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseStructDecl()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_TYPE {
        AssertNextToken(tok, TOKEN_IDENTIFIER);
        // identifier of struct in tok.strValue
        AssertNextToken(tok, TOKEN_STRUCT);
        AssertNextToken(tok, TOKEN_LCBRAC);
        ParseStructVarDeclList(tok);
        AssertNextToken(tok, TOKEN_RCBRAC);
        AssertNextToken(tok, TOKEN_SEMICOLON);
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
func ParseStructVarDeclList(tok *Token) {
    var boolFlag uint64;
    PrintDebugString("Entering ParseStructVarDeclList()",1000);
    for boolFlag = ParseStructVarDecl(tok);
        boolFlag == 0;
        boolFlag = ParseStructVarDecl(tok) { }
    PrintDebugString("Leaving ParseStructVarDeclList()",1000);
}

//
// Parses: identifier type ";"
// A single variable declaration in a struct.
//
func ParseStructVarDecl(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseStructVarDecl()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_IDENTIFIER {
        ParseType(tok);
        AssertNextToken(tok, TOKEN_SEMICOLON);
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
func ParseType(tok *Token) {
    PrintDebugString("Entering ParseType()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_LSBRAC {   
        AssertNextToken(tok, TOKEN_INTEGER);
        // value of integer in tok.intValue
        AssertNextToken(tok, TOKEN_RSBRAC);        
    } else {
        SyncToken(tok);
    }

    AssertNextToken(tok, TOKEN_IDENTIFIER);
    // typename in tok.strValue
    PrintDebugString("Leaving ParseType()",1000);
}

//
// Parses: [ "[" integer "]" ] identifier
// Is completelly optional. Only used to parse return value of a function
// declaration.
//
func ParseTypeOptional(tok *Token) {
    PrintDebugString("Entering ParseTypeOptional()",1000); 
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_LSBRAC {
        AssertNextToken(tok, TOKEN_INTEGER);        
        AssertNextToken(tok, TOKEN_RSBRAC);
    } else {
        tok.nextToken = tok.id;
    }
    GetNextTokenSafe(tok);
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
func ParseVarDeclList(tok *Token) {
    var boolFlag uint64;
    PrintDebugString("Entering ParseVarDeclList()",1000);
    for boolFlag = ParseVarDecl(tok);
        boolFlag == 0;
        boolFlag = ParseVarDecl(tok) { }
    PrintDebugString("Leaving ParseVarDeclList()",1000);
}

//
// Parses: "var" identifier type [ "=" expression ];
// Is used to parse a single variable declaration with optional initializer.
//
func ParseVarDecl(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseVarDecl()",1000);
    boolFlag = LookAheadAndCheck(tok, TOKEN_VAR);
    if boolFlag == 0 {
        AssertNextToken(tok, TOKEN_VAR);
        AssertNextToken(tok, TOKEN_IDENTIFIER);
        // variable name in tok.strValue
        ParseType(tok);

        GetNextTokenSafe(tok);
        if tok.id == TOKEN_ASSIGN {
            ParseExpression(tok);        
        } else {
            SyncToken(tok);
        } 

        AssertNextToken(tok, TOKEN_SEMICOLON);
        boolFlag = 0;
    }
    PrintDebugString("Leaving ParseVarDecl()",1000);
    return boolFlag;
}

//
//
//
func ParseExpression(tok *Token) {
    PrintDebugString("Entering ParseExpression()",1000);
    ParseSimpleExpression(tok);    
    ParseCmpOp(tok);
    PrintDebugString("Leaving ParseExpression()",1000);
}

//
//
//
func ParseCmpOp(tok *Token) {
    PrintDebugString("Entering ParseCmpOp()",1000);
    GetNextTokenSafe(tok);
    if (tok.id == TOKEN_EQUALS) || (tok.id == TOKEN_NOTEQUAL) || 
        (tok.id == TOKEN_REL_LT) || (tok.id == TOKEN_REL_LTOE) || 
        (tok.id == TOKEN_REL_GT) || (tok.id == TOKEN_REL_GTOE) {
        ParseSimpleExpression(tok);
    } else {
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseCmpOp()",1000);
}

//
//
//
func ParseSimpleExpression(tok *Token) {
    var boolFlag uint64;
    PrintDebugString("Entering ParseSimpleExpression()",1000);
    ParseUnaryArithOp(tok);
    ParseTerm(tok);
    for boolFlag = ParseSimpleExpressionOp(tok);
        boolFlag == 0;
        boolFlag = ParseSimpleExpressionOp(tok) { }
    PrintDebugString("Leaving ParseSimpleExpression()",1000);
}

//
//
//
func ParseSimpleExpressionOp(tok *Token) uint64 {
    var boolFlag uint64 = 1;
    PrintDebugString("Entering ParseSimpleExpressionOp()",1000);
    boolFlag = ParseUnaryArithOp(tok); // +,-
    if boolFlag != 0 {
        GetNextTokenSafe(tok);
        if tok.id == TOKEN_REL_OR {
            // ||
            boolFlag = 0;
        } 
    }
    if boolFlag == 0 {
        ParseTerm(tok);
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
func ParseUnaryArithOp(tok *Token) uint64 {
    var boolFlag uint64 = 1;
    PrintDebugString("Entering ParseUnaryArithOp()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_ARITH_PLUS {
        // *
        boolFlag = 0;
    }
    if tok.id == TOKEN_ARITH_MINUS {
        // /
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
func ParseTerm(tok *Token) {
    var boolFlag uint64;
    PrintDebugString("Entering ParseTerm()",1000);
    ParseFactor(tok);
    for boolFlag = ParseTermOp(tok);
        boolFlag == 0;
        boolFlag = ParseTermOp(tok) { }      
    PrintDebugString("Leaving ParseTerm()",1000);
}

//
//
//
func ParseTermOp(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseTermOp()",1000);
    boolFlag = ParseBinaryArithOp(tok); // *,/
    if boolFlag != 0 {
        GetNextTokenSafe(tok);
        if tok.id == TOKEN_REL_AND {
            // &&
            boolFlag = 0;
        }
    }
    if boolFlag == 0 {
        ParseFactor(tok);
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
func ParseBinaryArithOp(tok *Token) uint64 {
    var boolFlag uint64 = 1;
    PrintDebugString("Entering ParseBinaryArithOp()",1000);
    GetNextTokenSafe(tok);
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
func ParseFactor(tok *Token) uint64 {
    var doneFlag uint64 = 1;
    var boolFlag uint64;
    PrintDebugString("Entering ParseFactor()",1000);
    GetNextTokenSafe(tok);
    if (doneFlag == 1) && (tok.id == TOKEN_OP_ADR) {
        AssertNextToken(tok, TOKEN_IDENTIFIER);
        ParseSelector(tok);
        ParseFunctionCallOptional(tok);
        doneFlag = 0;
    }
    if (doneFlag == 1) && (tok.id == TOKEN_IDENTIFIER) {
        ParseSelector(tok);
        ParseFunctionCallOptional(tok);
        doneFlag = 0;
    } 
    if (doneFlag == 1) && (tok.id == TOKEN_INTEGER) {
        doneFlag = 0;
    }
    if (doneFlag) == 1 && (tok.id == TOKEN_STRING) {
        doneFlag = 0;
    }
    if (doneFlag) == 1 && (tok.id == TOKEN_LBRAC) {
        ParseExpression(tok);
        AssertNextToken(tok, TOKEN_RBRAC);
        doneFlag = 0;
    }
    if (doneFlag == 1) && (tok.id == TOKEN_NOT) {
        ParseFactor(tok);
        doneFlag = 0;
    }

    if doneFlag != 0 {
        boolFlag = 1;
        tok.nextToken = tok.id;
    } else {
        boolFlag = 0;
    }
    PrintDebugString("Leaving ParseFactor()",1000);
    return boolFlag;
}

//
//
//
func ParseSelector(tok *Token) {
    var boolFlag uint64;
    PrintDebugString("Entering ParseSelector()",1000);
    for boolFlag = ParseSelectorSub(tok);
        boolFlag == 0; 
        boolFlag = ParseSelectorSub(tok) {
    }
    PrintDebugString("Leaving ParseSelector()",1000);
}

//
//
//
func ParseSelectorSub(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseSelectorSub()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_PT {
        AssertNextToken(tok, TOKEN_IDENTIFIER);
        // value in tok.strValue
        boolFlag = 0;
    } else {
        if tok.id == TOKEN_LSBRAC {
            GetNextTokenSafe(tok);
            if tok.id == TOKEN_INTEGER {
                
            } else {
                if tok.id == TOKEN_IDENTIFIER {  
                    ParseSelector(tok);
                }
            } 

            AssertNextToken(tok, TOKEN_RSBRAC);
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
func ParseFuncDeclList(tok *Token) {
    var boolFlag uint64; 
    PrintDebugString("Entering ParseFuncDeclList()",1000);
    for boolFlag = ParseFuncDeclListSub(tok);
        boolFlag == 0; 
        boolFlag = ParseFuncDeclListSub(tok) { }
    PrintDebugString("Leaving ParseFuncDeclList()",1000);
}

//
//
//
func ParseFuncDeclListSub(tok *Token) uint64 {
    var es [255]uint64;
    var boolFlag uint64;    
    PrintDebugString("Entering ParseFuncDeclListSub()",1000);
    boolFlag = ParseFuncDeclHead(tok);
    if boolFlag == 0 {

        boolFlag = ParseFuncDeclRaw(tok);
        if boolFlag != 0 {
            boolFlag = ParseFuncDecl(tok);
        }
        if boolFlag != 0 {
            es[0] = TOKEN_SEMICOLON;
            es[1] = TOKEN_LCBRAC;
            ParseError(tok.id,es,2);
        }
    }
    PrintDebugString("Leaving ParseFuncDeclListSub()",1000);
    return boolFlag;
}

func ParseFuncDeclHead(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseFuncDeclHead()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_FUNC {
        AssertNextToken(tok, TOKEN_IDENTIFIER);
        // function name in tok.strValue
        AssertNextToken(tok, TOKEN_LBRAC);
        ParseIdentifierTypeList(tok);
        AssertNextToken(tok, TOKEN_RBRAC);
        ParseTypeOptional(tok);
        boolFlag = 0;
    } else {    
        SyncToken(tok);
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseFuncDeclHead()",1000);
    return boolFlag;
}

func ParseFuncDeclRaw(tok *Token) uint64 {
    var boolFlag uint64 = 1;
    PrintDebugString("Entering ParseFuncDeclRaw()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_SEMICOLON {
        boolFlag = 0;
    } else {
        SyncToken(tok);
    }
    PrintDebugString("Leaving ParseFuncDeclRaw()",1000);
    return boolFlag;
}

func ParseFuncDecl(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseFuncDecl()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_LCBRAC {
        ParseVarDeclList(tok);
        ParseStatementSequence(tok);
        GetNextTokenSafe(tok);
        if tok.id == TOKEN_RETURN {
            ParseExpression(tok);
            AssertNextToken(tok, TOKEN_SEMICOLON);
        } else {
            SyncToken(tok);
        }
        AssertNextToken(tok, TOKEN_RCBRAC);
        boolFlag = 0;
    } else {
        SyncToken(tok);
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseFuncDecl()",1000);
    return boolFlag;
}

func ParseIdentifierTypeList(tok *Token) {
    var boolFlag uint64;
    PrintDebugString("Entering ParseIdentifierTypeList()",1000);
    boolFlag = ParseIdentifierType(tok);
    if boolFlag == 0 {
        for boolFlag = ParseIdentifierTypeListSub(tok);
            boolFlag == 0; 
            boolFlag = ParseIdentifierTypeListSub(tok) { }   
    }
    PrintDebugString("Leaving ParseIdentifierTypeList()",1000);
}

func ParseIdentifierTypeListSub(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseIdentifierTypeListSub()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_COLON {
        boolFlag = ParseIdentifierType(tok);
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }
    PrintDebugString("Leaving ParseIdentifierTypeListSub()",1000);
    return boolFlag;
}

func ParseIdentifierType(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseIdentifierType()",1000);
    GetNextTokenSafe(tok);
    if tok.id != TOKEN_IDENTIFIER {
        tok.nextToken = tok.id;
        boolFlag = 1;        
    } else {
        GetNextTokenSafe(tok);
        if tok.id == TOKEN_ARITH_MUL {

        } else {
            tok.nextToken = tok.id;
        }
        ParseType(tok);
        boolFlag = 0;
    }
    PrintDebugString("Leaving ParseIdentifierType()",1000);
    return boolFlag;
}

func ParseStatementSequence(tok *Token) {
    var boolFlag uint64; 
    PrintDebugString("Entering ParseStatementSequence()",1000);
    for boolFlag = ParseStatement(tok);
        boolFlag == 0;
        boolFlag = ParseStatement(tok) { }
    PrintDebugString("Leaving ParseStatementSequence()",1000);
}

func ParseStatement(tok *Token) uint64 {
    var boolFlag uint64;
    var doneFlag uint64;
    var es [255]uint64; 
    PrintDebugString("Entering ParseStatement()",1000);
    doneFlag = 1;

    GetNextTokenSafe(tok);
    if (doneFlag == 1) && (tok.id == TOKEN_IDENTIFIER) {
        // Could be assignment or a function call.
        // Cannnot be resolved until selectors are all parsed
        // To be improved!
        ParseSelector(tok);
        tok.nextToken = tok.id;
        boolFlag = ParseAssignment(tok);
        if boolFlag != 0 {
            tok.nextToken = tok.id;
            boolFlag = ParseFunctionCallStatement(tok);
        }        
        if boolFlag != 0 {
            ParseError(tok.id,es,0);
        }
        doneFlag = 0;
    }
    
    if (doneFlag == 1) && (tok.id == TOKEN_IF) {
        tok.nextToken = tok.id;
        ParseIfStatement(tok);
        doneFlag = 0;
    }

    if (doneFlag == 1) && (tok.id == TOKEN_FOR) {
        tok.nextToken = tok.id;
        ParseForStatement(tok);
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

func ParseAssignment(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseAssignment()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_ASSIGN {
        ParseExpression(tok);
        AssertNextToken(tok, TOKEN_SEMICOLON);
        boolFlag = 0;
    } else {
        SyncToken(tok);
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseAssignment()",1000);
    return boolFlag;
}

func ParseAssignmentWithoutSC(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseAssignmentWithoutSC()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_ASSIGN {      
        ParseExpression(tok);
        boolFlag = 0;
    } else {
        SyncToken(tok);
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseAssignmentWithoutSC()",1000);
    return boolFlag;
}

func ParseFunctionCallOptional(tok *Token) {
    PrintDebugString("Entering ParseFunctionCallOptional()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_LBRAC {
        GetNextTokenSafe(tok);
        if tok.id == TOKEN_RBRAC {

        } else {
            tok.nextToken = tok.id;
            ParseExpressionList(tok);
            AssertNextToken(tok, TOKEN_RBRAC);    
        }
    } else {
        SyncToken(tok);
    }
    PrintDebugString("Leaving ParseFunctionCallOptional()",1000);
}

func ParseFunctionCall(tok *Token) {
    PrintDebugString("Entering ParseFunctionCall()",1000);
    AssertNextToken(tok, TOKEN_LBRAC);
    GetNextTokenSafe(tok);
    if tok.id != TOKEN_RBRAC {
        SyncToken(tok);
        ParseExpressionList(tok);
        AssertNextToken(tok, TOKEN_RBRAC);     
    }     
    PrintDebugString("Leaving ParseFunctionCall()",1000);
}

func ParseExpressionList(tok *Token) {
    var boolFlag uint64;
    PrintDebugString("Entering ParseExpressionList()",1000);
    ParseExpression(tok);
    for boolFlag = ParseExpressionListSub(tok);
        boolFlag == 0;
        boolFlag = ParseExpressionListSub(tok) { }   
    PrintDebugString("Leaving ParseExpressionList()",1000);
}

func ParseExpressionListSub(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseExpressionListSub()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_COLON {
        ParseExpression(tok);
        boolFlag = 0;
    } else {
        SyncToken(tok);
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseExpressionListSub()",1000);
    return boolFlag;   
}

func ParseFunctionCallStatement(tok *Token) uint64 {
    var boolFlag uint64;
    PrintDebugString("Entering ParseFunctionCallStatement()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_LBRAC {
        tok.nextToken = tok.id;
        ParseFunctionCall(tok);                
        boolFlag = 0;
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    PrintDebugString("Leaving ParseFunctionCallStatement()",1000);
    return boolFlag;
}

func ParseForStatement(tok *Token) {
    PrintDebugString("Entering ParseForStatement()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_FOR {
        GetNextTokenSafe(tok);
        if tok.id == TOKEN_SEMICOLON {
            SyncToken(tok);
        } else {
            SyncToken(tok);
            AssertNextToken(tok, TOKEN_IDENTIFIER);
            // tok.strValue
            ParseSelector(tok);
            ParseAssignmentWithoutSC(tok);
        }
        
        AssertNextToken(tok, TOKEN_SEMICOLON);

        GetNextTokenSafe(tok);
        if tok.id == TOKEN_SEMICOLON {
            SyncToken(tok);
        } else {
            SyncToken(tok);
            ParseExpression(tok);
        }

        AssertNextToken(tok, TOKEN_SEMICOLON);

        GetNextTokenSafe(tok);
        if tok.id == TOKEN_LCBRAC {
            SyncToken(tok);
        } else {
            SyncToken(tok);
            AssertNextToken(tok, TOKEN_IDENTIFIER);
            // tok.strValue
            ParseSelector(tok);
            ParseAssignmentWithoutSC(tok);
        }

        AssertNextToken(tok, TOKEN_LCBRAC);        
        ParseStatementSequence(tok);
        AssertNextToken(tok, TOKEN_RCBRAC);

    } else {
        SyncToken(tok);
    }   
    PrintDebugString("Leaving ParseForStatement()",1000);
}

func ParseIfStatement(tok *Token) {
    PrintDebugString("Entering ParseIfStatement()",1000);
    GetNextTokenSafe(tok);
    if tok.id == TOKEN_IF {
        ParseExpression(tok);
        AssertNextToken(tok, TOKEN_LCBRAC);
        ParseStatementSequence(tok);
        AssertNextToken(tok, TOKEN_RCBRAC);

        GetNextTokenSafe(tok);
        if tok.id == TOKEN_ELSE {
            ParseElseStatement(tok);
        } else {
            SyncToken(tok);
        }

    } else {
        SyncToken(tok);
    }
    PrintDebugString("Leaving ParseIfStatement()",1000);
}

func ParseElseStatement(tok *Token) {
    PrintDebugString("Entering ParseElseStatement()",1000);
    AssertNextToken(tok, TOKEN_LCBRAC);
    ParseStatementSequence(tok);
    AssertNextToken(tok, TOKEN_RCBRAC);
    PrintDebugString("Leaving ParseElseStatement()",1000);
}
