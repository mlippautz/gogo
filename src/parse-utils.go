// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// File containing utility functions that are used by the parser.
//

package main

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
    var index uint64;

    if tok.nextToken != 0 {
        tok.id = tok.nextToken;      
        tok.nextToken = 0;  
    } else {
        if tok.llCnt > 0 {
            index = tok.toRead;
            tok.id = tok.nextTokenId[index];
            tok.strValue = tok.nextTokenValStr[index];
            tok.llCnt = tok.llCnt-1;
            tok.toRead = index+1;
        } else {
            GetNextToken();
            tokenString = TokenToString(tok.id);
            PrintDebugString(tokenString, 1000);
        }
    }
}

func AssertNextTokenWeak(tokenNumber uint64) {
    GetNextTokenSafe();
    if tok.id != tokenNumber {
        ParseErrorWeak(tok.id, tokenNumber, 0, 1);
        tok.nextToken = tok.id;
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
    if tok.id != tokenNumber {
        ParseErrorWeak(tok.id, tokenNumber, 0, 1);
        //tok.nextToken = tok.id;
        ParserSync();
    }
}
