// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

//
// Function printing a parse error using only libgogo.
// ue ... unexpected token
// e .... array of expected tokens
// eLen . actual length (items) of array
//
func parseError(ue uint64, e [255]uint64, eLen uint64) {
    var i uint64;
    var str string;

    libgogo.PrintString(filename);
    libgogo.PrintString(":");
    libgogo.PrintNumber(lineCounter);
    libgogo.PrintString(": syntax error: unexpected token '");
    str = TokenToString(ue);
    libgogo.PrintString(str);
    libgogo.PrintString("'");

    if eLen > 0 {
        libgogo.PrintString(", expecting one of: ");
        str = TokenToString(e[i]);
        libgogo.PrintString(str);
        for i = 1; i < eLen; i = i+1 {
            str = TokenToString(e[i]);
            libgogo.PrintString(str);
            libgogo.PrintString(", ");        
        }
        libgogo.PrintString("\n");
    }
    libgogo.Exit(2); 
}

//
//
//
func GetNextTokenSafe(fd uint64, tok *Token) {
    if tok.nextToken != 0 {
        tok.id = tok.nextToken;      
        tok.nextToken = 0;  
    } else {
        GetNextToken(fd, tok);
    }    
}


//
// Main parsing function
//
func Parse( fd uint64 ) {
    var tok Token;

    tok.id = 0;
    tok.nextChar = 0;
    tok.nextToken = 0;    

    ParsePackageStatement(fd, &tok);    
    ParseImportStatementList(fd, &tok);
    ParseStructDeclList(fd, &tok);
    ParseVarDeclList(fd, &tok);
    ParseFuncDeclList(fd, &tok);

    // Scan the rest for debugging purposes
    for GetNextToken(fd,&tok); tok.id != TOKEN_EOS; GetNextToken(fd,&tok) {
        debugToken(&tok);
    }
}

//
// Check for package "identifier" as the golang forces this.
//
func ParsePackageStatement(fd uint64, tok *Token) {    
    var es [255]uint64;

    GetNextTokenSafe(fd,tok);
    if tok.id != TOKEN_PACKAGE {
        es[0] = TOKEN_PACKAGE;
        parseError(tok.id,es,1);
    } else {
        GetNextTokenSafe(fd,tok);
        if tok.id != TOKEN_IDENTIFIER {
            es[0] = TOKEN_IDENTIFIER;
            parseError(tok.id,es,1);
        }
    }
    // package ok
}

//
// Parses: { import_stmt }
//
func ParseImportStatementList(fd uint64, tok *Token) {
    var boolFlag uint64;
    for boolFlag = ParseImportStatement(fd, tok);boolFlag == 0;boolFlag = ParseImportStatement(fd, tok) {
    }
}


//
// Parses: "import" string
//
func ParseImportStatement(fd uint64, tok *Token) uint64 {
    var es [255]uint64;
    var boolFlag uint64;

    GetNextTokenSafe(fd,tok);
    if tok.id == TOKEN_IMPORT {
        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_STRING  {
            es[0] = TOKEN_STRING;
            parseError(tok.id,es,1);
        } else {
            // import ok
        }   
        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }    
    return boolFlag;
}

//
// Parses: { struct_decl }
//
func ParseStructDeclList(fd uint64, tok *Token) {
    var boolFlag uint64;
    for boolFlag = ParseStructDecl(fd, tok);boolFlag == 0;boolFlag = ParseStructDecl(fd, tok) {
    }
}

//
// Parses: "type" identifier "struct" "{" struct_var_decl_list "}" ";"
//
func ParseStructDecl(fd uint64, tok *Token) uint64 {
    var es [255]uint64;
    var boolFlag uint64;

    GetNextTokenSafe(fd,tok);
    if tok.id == TOKEN_TYPE {

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_IDENTIFIER  {
            es[0] = TOKEN_IDENTIFIER;
            parseError(tok.id,es,1);
        }

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_STRUCT  {
            es[0] = TOKEN_STRUCT;
            parseError(tok.id,es,1);
        }    

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_LCBRAC  {
            es[0] = TOKEN_LCBRAC;
            parseError(tok.id,es,1);
        }

        ParseStructVarDeclList(fd, tok);

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_RCBRAC  {
            es[0] = TOKEN_RCBRAC;
            parseError(tok.id,es,1);
        }

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_SEMICOLON  {
            es[0] = TOKEN_SEMICOLON;
            parseError(tok.id,es,1);
        }
         
        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }    
    return boolFlag;
}

//
// Parses: { struct_var_decl }
//
func ParseStructVarDeclList(fd uint64, tok *Token) {
    var boolFlag uint64;
    for boolFlag = ParseStructVarDecl(fd, tok);boolFlag == 0;boolFlag = ParseStructVarDecl(fd, tok) {
    }
}

//
// Parses: identifier type ";"
//
func ParseStructVarDecl(fd uint64, tok *Token) uint64 {
    var es [255]uint64;
    var boolFlag uint64;

    GetNextTokenSafe(fd,tok);
    if tok.id == TOKEN_IDENTIFIER {

        ParseType(fd, tok);

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_SEMICOLON  {
            es[0] = TOKEN_SEMICOLON;
            parseError(tok.id,es,1);
        }
        
        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }    
    return boolFlag;    
}

//
// Parses: [ "[" integer "]" ] identifier 
//
func ParseType(fd uint64, tok *Token) {
    var es [255]uint64;

    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_LSBRAC {        
        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_INTEGER  {
            es[0] = TOKEN_INTEGER;
            parseError(tok.id,es,1);
        }       

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_RSBRAC  {
            es[0] = TOKEN_RSBRAC;
            parseError(tok.id,es,1);
        }
    } else {
        tok.nextToken = tok.id;
    }

    GetNextTokenSafe(fd, tok);
    if tok.id != TOKEN_IDENTIFIER  {
        es[0] = TOKEN_IDENTIFIER;
        parseError(tok.id,es,0);
    } 
}

//
// Parses: [ "[" integer "]" ] identifier 
//
func ParseTypeOptional(fd uint64, tok *Token) {
    var es [255]uint64;

    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_LSBRAC {        
        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_INTEGER  {
            es[0] = TOKEN_INTEGER;
            parseError(tok.id,es,1);
        }       

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_RSBRAC  {
            es[0] = TOKEN_RSBRAC;
            parseError(tok.id,es,1);
        }
    } else {
        tok.nextToken = tok.id;
    }

    GetNextTokenSafe(fd, tok);
    if tok.id != TOKEN_IDENTIFIER  {
        tok.nextToken = tok.id;
    } 
}

//
//
//
func ParseVarDeclList(fd uint64, tok *Token) {
    var boolFlag uint64;
    for boolFlag = ParseVarDecl(fd, tok);boolFlag == 0;boolFlag = ParseVarDecl(fd, tok) {
    }
}

//
//
//
func ParseVarDecl(fd uint64, tok *Token) uint64 {
    var es [255]uint64;
    var boolFlag uint64;

    GetNextTokenSafe(fd,tok);
    if tok.id == TOKEN_VAR {
        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_IDENTIFIER {
            es[0] = TOKEN_IDENTIFIER;
            parseError(tok.id,es,1);
        } 
        
        ParseType(fd, tok);

        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_ASSIGN {
            ParseExpression(fd, tok);        
        } else {
            tok.nextToken = tok.id;
        } 

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_SEMICOLON {
            es[0] = TOKEN_SEMICOLON;
            parseError(tok.id,es,1);
        }

        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }    
    return boolFlag;
}

//
//
//
func ParseExpression(fd uint64, tok *Token) {
    ParseSimpleExpression(fd, tok);    
    ParseCmpOp(fd, tok);
}

//
//
//
func ParseCmpOp(fd uint64, tok *Token) {
    GetNextTokenSafe(fd, tok);

    if tok.id == TOKEN_EQUALS || tok.id == TOKEN_NOTEQUAL || tok.id == TOKEN_REL_LT || tok.id == TOKEN_REL_LTOE || tok.id == TOKEN_REL_GT || tok.id == TOKEN_REL_GTOE {
        ParseSimpleExpression(fd, tok);
    } else {
        tok.nextToken = tok.id;
    }
}

//
//
//
func ParseSimpleExpression(fd uint64, tok *Token) {
    var boolFlag uint64;
   
    ParseUnaryArithOp(fd, tok);
    ParseTerm(fd, tok);

    for boolFlag = ParseSimpleExpressionOp(fd, tok);boolFlag == 0;boolFlag = ParseSimpleExpressionOp(fd, tok) {
    }
}

//
//
//
func ParseSimpleExpressionOp(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    
    boolFlag = ParseUnaryArithOp(fd, tok);
    if boolFlag == 0 {
        // read +/-
    } else {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_REL_OR {
            // read ||
            boolFlag = 0;
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;   
        }
    }

    if boolFlag == 0 {
        ParseTerm(fd, tok);
    }

    return boolFlag;
}

//
//
//
func ParseUnaryArithOp(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_ARITH_PLUS {
        boolFlag = 0;
    } else {
        if tok.id == TOKEN_ARITH_MINUS {
            boolFlag = 0;            
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;
        }
    }

    return boolFlag;   
}

//
//
//
func ParseBinaryArithOp(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_ARITH_MUL {
        boolFlag = 0;
    } else {
        if tok.id == TOKEN_ARITH_DIV {
            boolFlag = 0;            
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;
        }
    }

    return boolFlag; 
}

//
//
//
func ParseTerm(fd uint64, tok *Token) {
    var boolFlag uint64;
    ParseFactor(fd, tok);
    for boolFlag = ParseTermOp(fd, tok);boolFlag == 0;boolFlag = ParseTermOp(fd, tok) {
    }      
}

//
//
//
func ParseTermOp(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    
    boolFlag = ParseBinaryArithOp(fd, tok);
    if boolFlag == 0 {
        // read *//
    } else {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_REL_AND {
            // read &&
            boolFlag = 0;
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;   
        }
    }

    if boolFlag == 0 {
        ParseFactor(fd, tok);
    }

    return boolFlag;
}

//
//
//
func ParseFactor(fd uint64, tok *Token) {
    var es [255]uint64;
    var doneFlag uint64 = 1;

    GetNextTokenSafe(fd, tok);

    if doneFlag == 1 && tok.id == TOKEN_OP_ADR {
        GetNextTokenSafe(fd, tok);        
        if tok.id == TOKEN_IDENTIFIER {
            ParseSelector(fd ,tok);
            ParseFunctionCallOptional(fd, tok);
            doneFlag = 0;
        } else {
            es[0] = TOKEN_IDENTIFIER;
            parseError(tok.id, es, 1);
        }
    }
    if doneFlag == 1 && tok.id == TOKEN_IDENTIFIER {
        ParseSelector(fd ,tok);
        doneFlag = 0;
    } 
    if doneFlag == 1 && tok.id == TOKEN_INTEGER {
        doneFlag = 0;
    }
    if doneFlag == 1 && tok.id == TOKEN_STRING {
        doneFlag = 0;
    }
    if doneFlag == 1 && tok.id == TOKEN_LBRAC {
        ParseExpression(fd, tok);
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_RBRAC {
            doneFlag = 0;
        } else {
            es[0] = TOKEN_RBRAC;
            parseError(tok.id,es,1);
        }
    }
    if doneFlag == 1 && tok.id == TOKEN_NOT {
        ParseFactor(fd, tok);
        doneFlag = 0;
    }

    if doneFlag != 0 {
        es[0] = TOKEN_IDENTIFIER;
        es[1] = TOKEN_INTEGER;
        es[2] = TOKEN_STRING;
        es[3] = TOKEN_LBRAC;
        es[4] = TOKEN_NOT;
        parseError(tok.id,es,5);
    }
}

//
//
//
func ParseSelector(fd uint64, tok *Token) {
    var boolFlag uint64;
    for boolFlag = ParseSelectorSub(fd, tok);boolFlag == 0; boolFlag = ParseSelectorSub(fd, tok) {
    }
}

//
//
//
func ParseSelectorSub(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    var es [255]uint64;

    GetNextTokenSafe(fd, tok);
    
    if tok.id == TOKEN_PT {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_IDENTIFIER  {
            boolFlag = 0;
        } else {
            es[0] = TOKEN_IDENTIFIER;
            parseError(tok.id,es,1);
        }
    } else {
        if tok.id == TOKEN_LSBRAC {
            GetNextTokenSafe(fd, tok);
            if tok.id == TOKEN_INTEGER {
                
            } else {
                if tok.id == TOKEN_IDENTIFIER {
                    ParseSelector(fd, tok);
                }
            } 

            GetNextTokenSafe(fd, tok);
            if tok.id == TOKEN_RSBRAC {
                boolFlag = 0;
            } else {
                es[0] = TOKEN_RSBRAC;
                parseError(tok.id,es,1);
            }
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;
        }
    }

    return boolFlag;
}

func ParseFuncDeclList(fd uint64, tok *Token) {
    var boolFlag uint64; 
    for boolFlag = ParseFuncDeclListSub(fd, tok);boolFlag == 0; boolFlag = ParseFuncDeclListSub(fd, tok) {
    }
}

func ParseFuncDeclListSub(fd uint64, tok *Token) uint64 {
    var es [255]uint64;
    var boolFlag uint64;    
    boolFlag = ParseFuncDeclHead(fd, tok);
    if boolFlag == 0 {
        boolFlag = ParseFuncDeclRaw(fd, tok);
        if boolFlag != 0 {
            boolFlag = ParseFuncDecl(fd, tok);
        }
        if boolFlag != 0 {
            es[0] = TOKEN_SEMICOLON;
            es[1] = TOKEN_LCBRAC;
            parseError(tok.id,es,2);
        }
    }
    return boolFlag;
}

func ParseFuncDeclHead(fd uint64, tok *Token) uint64 {
    var es [255]uint64; 
    var boolFlag uint64;

    GetNextTokenSafe(fd, tok);
    if tok.id != TOKEN_FUNC {
        boolFlag = 1;
        tok.nextToken = tok.id;
    } else  {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_IDENTIFIER  {

        } else {
            es[0] = TOKEN_IDENTIFIER;
            parseError(tok.id,es,1);        
        }

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_LBRAC  {
            es[0] = TOKEN_LBRAC;
            parseError(tok.id,es,1);        
        }

        ParseIdentifierTypeList(fd, tok);

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_RBRAC  {
            es[0] = TOKEN_RBRAC;
            parseError(tok.id,es,1);        
        }
        
        ParseTypeOptional(fd, tok);
        boolFlag = 0;
    }

    return boolFlag;
}

func ParseFuncDeclRaw(fd uint64, tok *Token) uint64 {
    var boolFlag uint64 = 1;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_SEMICOLON {
        boolFlag = 0;
    }
    return boolFlag;
}

func ParseFuncDecl(fd uint64, tok *Token) uint64 {
    var boolFlag uint64 = 1;
    var es [255]uint64; 
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_LCBRAC {
        boolFlag = 0;
    } else {
        ParseVarDeclList(fd, tok);
        ParseStatementSequence(fd, tok);
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_RETURN {
            ParseExpression(fd, tok);
            GetNextTokenSafe(fd, tok);
            if tok.id != TOKEN_SEMICOLON {
                es[0] = TOKEN_SEMICOLON;
                parseError(tok.id,es,1);
            }
        } else {
            tok.nextToken = tok.id;
        }
    }
    return boolFlag;
}

func ParseIdentifierTypeList(fd uint64, tok *Token) {
    var boolFlag uint64;

    boolFlag = ParseIdentifierType(fd, tok);
    if boolFlag == 0 {
        for boolFlag = ParseIdentifierTypeListSub(fd, tok);boolFlag == 0; boolFlag = ParseIdentifierTypeListSub(fd, tok) {
        }   
    }
}

func ParseIdentifierTypeListSub(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_COLON {
        boolFlag = ParseIdentifierType(fd, tok);
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }
    return boolFlag;
}

func ParseIdentifierType(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;

    GetNextTokenSafe(fd, tok);
    if tok.id != TOKEN_IDENTIFIER {
        tok.nextToken = tok.id;
        boolFlag = 1;        
    } else {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_ARITH_MUL {

        } else {
            tok.nextToken = tok.id;
        }
        ParseType(fd, tok);
        boolFlag = 0;
    }

    return boolFlag;
}

func ParseStatementSequence(fd uint64, tok *Token) {
    var boolFlag uint64; 
    for boolFlag = ParseStatement(fd, tok);boolFlag == 0; boolFlag = ParseStatement(fd, tok) {
    }
}

func ParseStatement(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    var doneFlag uint64;
    var es [255]uint64; 

    doneFlag = 1;

    GetNextTokenSafe(fd, tok);
    if doneFlag == 1 && tok.id == TOKEN_IDENTIFIER {
        tok.nextToken = tok.id;
        boolFlag = ParseAssignment(fd, tok);
        if boolFlag != 0 {
            boolFlag = ParseFunctionCallStatement(fd, tok);
        }        
        if boolFlag != 0 {
            parseError(tok.id,es,0);
        }
        doneFlag = 0;
    }
    
    if doneFlag == 1 && tok.id == TOKEN_IF {
        ParseIfStatement(fd, tok);
        doneFlag = 0;
    }

    if doneFlag == 1 && tok.id == TOKEN_FOR {
        ParseForStatement(fd, tok);
        doneFlag = 0;
    }

    if doneFlag == 1 && tok.id == TOKEN_SEMICOLON {
        // NOP    
        doneFlag = 0;
    }

    if doneFlag != 0 {
        boolFlag = 1;
    } else {
        boolFlag = 0;
    }

    return boolFlag;
}

func ParseAssignment(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    var es [255]uint64;

    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_IDENTIFIER {
        ParseSelector(fd, tok);
        
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_ASSIGN {
            ParseExpression(fd, tok);
            GetNextTokenSafe(fd, tok);
            if tok.id != TOKEN_SEMICOLON {
                es[0] = TOKEN_SEMICOLON;
                parseError(tok.id, es, 1);
            }   
        } else {
            tok.nextToken = tok.id;
            boolFlag = 1;
        }
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }

    return boolFlag;
}

func ParseFunctionCallOptional(fd uint64, tok *Token) {
    var es [255]uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_LBRAC {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_RBRAC {

        } else {
            ParseExpressionList(fd, tok);
            GetNextTokenSafe(fd, tok);
            if tok.id != TOKEN_RBRAC {
                es[0] = TOKEN_RBRAC;
                parseError(tok.id, es, 1);
            }        
        }
    }
}

func ParseFunctionCall(fd uint64, tok *Token) {
    var es [255]uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_LBRAC {
        if tok.id == TOKEN_RBRAC {

        } else {
            ParseExpressionList(fd, tok);
            GetNextTokenSafe(fd, tok);
            if tok.id != TOKEN_RBRAC {
                es[0] = TOKEN_RBRAC;
                parseError(tok.id, es, 1);
            }        
        }     
    } else {
        es[0] = TOKEN_LBRAC;
        parseError(tok.id,es,1);
    }
}

func ParseExpressionList(fd uint64, tok *Token) {
    var boolFlag uint64;

    ParseExpression(fd, tok);
    for boolFlag = ParseExpressionListSub(fd, tok);boolFlag == 0; boolFlag = ParseExpressionListSub(fd, tok) {
    }   
}

func ParseExpressionListSub(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_COLON {
        ParseExpression(fd, tok);
        boolFlag = 0;
    } else {
        boolFlag = 1;
        tok.nextToken = tok.id;
    }
    return boolFlag;   
}

func ParseFunctionCallStatement(fd uint64, tok *Token) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_IDENTIFIER {
        ParseSelector(fd, tok);
        ParseFunctionCall(fd, tok);                
    } else {
        tok.nextToken = tok.id;
        boolFlag = 1;
    }
    return boolFlag;
}

func ParseForStatement(fd uint64, tok *Token) {
    var es [255]uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_FOR {
        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_SEMICOLON {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            ParseAssignment(fd, tok);
        }
        
        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_SEMICOLON {
            es[0] = TOKEN_SEMICOLON;
            parseError(tok.id,es,1);
        }

        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_SEMICOLON {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            ParseExpression(fd, tok);
        }

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_SEMICOLON {
            es[0] = TOKEN_SEMICOLON;
            parseError(tok.id,es,1);
        }

        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_LCBRAC {
            tok.nextToken = tok.id;
        } else {
            tok.nextToken = tok.id;
            ParseAssignment(fd, tok);
        }

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_LCBRAC {
            es[0] = TOKEN_LCBRAC;
            parseError(tok.id,es,1);
        }
        
        ParseStatementSequence(fd, tok);

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_RCBRAC {
            es[0] = TOKEN_RCBRAC;
            parseError(tok.id,es,1);
        }        

    } else {
        tok.nextToken = tok.id;
    }   
}

func ParseIfStatement(fd uint64, tok *Token) {
    var es [255]uint64;
    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_IF {
        ParseExpression(fd, tok);

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_LCBRAC {
            es[0] = TOKEN_LCBRAC;
            parseError(tok.id,es,1);
        } 

        ParseStatementSequence(fd, tok);

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_RCBRAC {
            es[0] = TOKEN_RCBRAC;
            parseError(tok.id,es,1);
        } 

        GetNextTokenSafe(fd, tok);
        if tok.id == TOKEN_ELSE {
            tok.nextToken = tok.id;
            ParseElseStatement(fd, tok);
        } else {
            tok.nextToken = tok.id;
        }

    } else {
        tok.nextToken = tok.id;
    }
}

func ParseElseStatement(fd uint64, tok *Token) {
    var es [255]uint64;

    GetNextTokenSafe(fd, tok);
    if tok.id == TOKEN_ELSE {
        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_LCBRAC {
            es[0] = TOKEN_LCBRAC;
            parseError(tok.id,es,1);
        } 

        ParseStatementSequence(fd, tok);

        GetNextTokenSafe(fd, tok);
        if tok.id != TOKEN_RCBRAC {
            es[0] = TOKEN_RCBRAC;
            parseError(tok.id,es,1);
        } 
    }
}
