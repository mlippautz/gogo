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

    libgogo.PrintString(">> Parser: Unexpected token '");
    str = TokenToString(ue);
    libgogo.PrintString(str);
    libgogo.PrintString("' in line ");
    libgogo.PrintNumber(lineCounter);
    libgogo.PrintString(".\n");
    if eLen > 0 {
        libgogo.PrintString("           Expected one of: ");
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
