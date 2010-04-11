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

func ParseExpression(fd uint64, tok *Token) {
    ParseSimpleExpression(fd, tok);    
    ParseCmpOp(fd, tok);
}

func ParseCmpOp(fd uint64, tok *Token) {
    GetNextTokenSafe(fd, tok);

    if tok.id == TOKEN_EQUALS || tok.id == TOKEN_NOTEQUAL || tok.id == TOKEN_REL_LT || tok.id == TOKEN_REL_LTOE || tok.id == TOKEN_REL_GT || tok.id == TOKEN_REL_GTOE {
        ParseSimpleExpression(fd, tok);
    } else {
        tok.nextToken = tok.id;
    }
}

func ParseSimpleExpression(fd uint64, tok *Token) {
    var boolFlag uint64;
   
    ParseUnaryArithOp(fd, tok);
    ParseTerm(fd, tok);

    for boolFlag = ParseSimpleExpressionOp(fd, tok);boolFlag == 0;boolFlag = ParseSimpleExpressionOp(fd, tok) {
    }
}

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

func ParseTerm(fd uint64, tok *Token) {
    var boolFlag uint64;
    ParseFactor(fd, tok);
    for boolFlag = ParseTermOp(fd, tok);boolFlag == 0;boolFlag = ParseTermOp(fd, tok) {
    }      
}

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

func ParseFactor(fd uint64, tok *Token) {
    var es [255]uint64;
    var doneFlag uint64 = 1;

    GetNextTokenSafe(fd, tok);

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

func ParseSelector(fd uint64, tok *Token) {
    var boolFlag uint64;
    for boolFlag = ParseSelectorSub(fd, tok);boolFlag == 0; boolFlag = ParseSelectorSub(fd, tok) {
    }
}

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
