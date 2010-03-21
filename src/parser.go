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
func parseError ( ue uint64, e [255]uint64, eLen uint64) {
    var i uint64;

    libgogo.PrintString(">> Parser: Unexpected token '");
    libgogo.PrintString( TokenToString(ue) );   // maybe call in call will be supported
    libgogo.PrintString("'.\n");
    if eLen > 0 {
        libgogo.PrintString("           Expected one of: ");
        libgogo.PrintString( TokenToString(e[i]) );
        for i = 1; i < eLen; i = i+1 {
            libgogo.PrintString( TokenToString(e[i]) );
            libgogo.PrintString(", ");        
        }
        libgogo.PrintString("\n");
    }
    libgogo.Exit(2); 
}


//
// Main parsing function
//
func Parse( fd uint64 ) {
    var tok Token;
    var es [255]uint64;

    tok.id = 0;
    tok.nextChar = 0;
    

    //
    // Check for package "identifier" as the golang forces this.
    //
    GetNextToken(fd,&tok);
    if tok.id != TOKEN_PACKAGE {
        es[0] = TOKEN_PACKAGE;
        parseError(tok.id,es,1);
    } 

    // Parse the rest
    // TBD

    // Scan the rest for debugging purposes
    for GetNextToken(fd,&tok); tok.id != TOKEN_EOS; GetNextToken(fd,&tok) {
        debugToken(&tok);
    }
}
