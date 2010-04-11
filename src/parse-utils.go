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
func ParseError(ue uint64, e [255]uint64, eLen uint64) {
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
    } 
    libgogo.PrintString("\n");
    libgogo.Exit(2); 
}

//
// Safely gets the next token and stores it in the supplied token.
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
// Syncing a token after a look-ahead (LL1) has taken place.
//
func SyncToken(tok *Token) {
    if tok.id != 0 {
        tok.nextToken = tok.id;
    }
}

//
//
//
func AssertNextToken(fd uint64,tok *Token, tokenNumber uint64) {
    GetNextTokenSafe(fd, tok);
    AssertToken(tok, tokenNumber);
}

//
//
//
func AssertToken(tok *Token, tokenNumber uint64) {
    var expectedTokens [255]uint64;
    if tok.id != tokenNumber {
        expectedTokens[0] = tokenNumber;
        ParseError(tok.id, expectedTokens, 1);        
    }
}
