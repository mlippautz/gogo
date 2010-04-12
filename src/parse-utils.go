// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

//
// Function looks 1 token ahead and checks if it equals the provided one
// Returns: 0 if tokens match, 1 otherwise
//
func LookAheadAndCheck(tok *Token, tokenNumber uint64) uint64 {
    var boolFlag uint64;
    GetNextTokenSafe(tok);
    if (tok.id == tokenNumber) {
        boolFlag = 0;
    } else {
        boolFlag = 1;
    }
    tok.id = tok.nextToken;
    return boolFlag;
}

//
// Safely gets the next token and stores it in the supplied token.
//
func GetNextTokenSafe(tok *Token) {
    var tokenString string;
    if tok.nextToken != 0 {
        tok.id = tok.nextToken;      
        tok.nextToken = 0;  
    } else {
        GetNextToken(tok);
        tokenString = TokenToString(tok.id);
        PrintDebugString(tokenString, 1000);
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
func AssertNextToken(tok *Token, tokenNumber uint64) {
    GetNextTokenSafe(tok);
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
