// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// File containing utility functions that are used by the parser.
//

package main

import "./libgogo/_obj/libgogo"

//
// Function used to check whether a certain stackdepth has been reached.
// Throws a fatal error!
//
func IncAndCheckDepth() {
    curDepth = curDepth +1;
    if curDepth > maxDepth {
        libgogo.ExitError("Max code depth reached. Please modify your code.",2);
    }
}

//
// Function decrementing the current depth
// To be used after some codeblock that has been "protected" by IncAndCheckDepth
//
func DecDepth() {
    curDepth = curDepth -1;
}

//
// Function looks 1 token ahead and checks if it equals the provided one
// Returns: 0 if tokens match, 1 otherwise
//
func LookAheadAndCheck(tokenNumber uint64) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe();
    if (tok.id == tokenNumber) {
        boolFlag = 0;
    } else {
        boolFlag = 1;
    }
    tok.nextToken = tok.id;
    return boolFlag;
}

//
// Safely gets the next token and stores it in the supplied token.
//
func GetNextTokenSafe() {
    var tokenString string;
    if tok.nextToken != 0 {
        tok.id = tok.nextToken;      
        tok.nextToken = 0;  
    } else {
        GetNextToken();
        tokenString = TokenToString(tok.id);
        PrintDebugString(tokenString, 1000);
    }
}

//
// Checks if the next token matches a tokenNumber, or produces an error.
//
func AssertNextToken(tokenNumber uint64) {
    GetNextTokenSafe();
    AssertToken(tokenNumber);
}

//
// Asserts an already read token. (no GetNextToken() before the check)
//
func AssertToken(tokenNumber uint64) {
    var expectedTokens [2]uint64;
    if tok.id != tokenNumber {
        expectedTokens[0] = tokenNumber;
        ParseError(tok.id, expectedTokens, 1);        
    }
}
