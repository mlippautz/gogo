// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// This file holds the basic scanning routines that separate a source file
// into the various tokens

package main

import "./libgogo/_obj/libgogo"

//
// Function fetching a new character from the file that is compiled
// Adds counters for line and columns
//
func GetCharWrapped() byte {
    var singleChar byte;
    singleChar = libgogo.GetChar(fileInfo[curFileIndex].fd);
    if (singleChar == 10) {
        fileInfo[curFileIndex].charCounter = 1;
        fileInfo[curFileIndex].lineCounter = fileInfo[curFileIndex].lineCounter + 1;
    } else {
        fileInfo[curFileIndex].charCounter = fileInfo[curFileIndex].charCounter + 1;
    }
    return singleChar;
}

func GetNextTokenRaw() {
    var singleChar byte; // Byte holding the last read value
    // Flag indicating whether we are in a comment.
    // 0 for no comment
    // 1 for a single line comment 
    // 2 for a multi line comment
    var inComment uint64;
    var done uint64; // Flag indicating whether a cycle (Token) is finsihed 
    var spaceDone uint64; // Flag indicating whether an abolishment cycle is finished 
    var numBuf string;

    // Initialize variables
    done = 0;
    spaceDone = 0;
    inComment = 0;  
    tok.strValue = "";

    // If the previous cycle had to read the next char (and stored it), it is 
    // now used as first read
    if tok.nextChar == 0 {       
        singleChar = GetCharWrapped();
    } else {
        singleChar = tok.nextChar;
        tok.nextChar = 0;
    }

    // check if it is a valid read, or an EOF
    if singleChar == 0 {
        tok.id = TOKEN_EOS;
        done = 1;
        spaceDone = 1;
    }

    //
    // Cleaning Tasks
    // The next part strips out spaces, newlines, tabs, and comments
    // Comments can either be single line with double slashes (//) or multiline
    // using C++ syntax /* */ 
    //
    for ; spaceDone != 1; {

        // check whether a comment is starting
        if singleChar == '/' {
            // if we are in a comment skip the rest, get the next char otherwise
            if inComment == 0 {
                singleChar = GetCharWrapped(); 
                if singleChar == '/' {
                    // we are in a single line comment (until newline is found)
                    inComment = 1;
                } else {
                    if singleChar == '*' {
                        // we are in a multiline comment (until ending is found)
                        inComment = 2;
                    } else {
                        ScanErrorString("Unknown character combination for comments.");
                    }
                }
            }
        } 

        // check whether a multi-line comment is ending
        if singleChar == '*' {
            singleChar = GetCharWrapped();
            if singleChar == '/' {
                if inComment == 2 {
                    inComment = 0;
                    singleChar = GetCharWrapped();
                }
            }
        }

        // if character is a newline:
        //  *) if in a singleline comment, exit the comment
        //  *) skip otherwise
        if singleChar == 10 {
            if inComment == 1 {
                inComment = 0;
            } 
        } 

        // handle everything that is not a space,tab,newline
        if (singleChar != ' ') && (singleChar != 9) && (singleChar != 10) {
            // if not in a comment we have our current valid char
            if inComment == 0 {
                spaceDone = 1;
            } 

            // check if GetChar() returned EOF while skipping
            if singleChar == 0 {
                tok.id = TOKEN_EOS;
                spaceDone = 1;
                done = 1;
            }   
        }
    
        
        // if we are not done until now, get a new character and start another abolishing cycle        
        if spaceDone == 0 {        
            singleChar=GetCharWrapped();
        }
    }

    //
    // Actual scanning part starts here
    //

    // Catch identifiers
    // identifier = letter { letter | digit }.
    if (done != 1) && ((singleChar >= 'A') && (singleChar <= 'Z')) || ((singleChar >= 'a') && (singleChar <= 'z')) || (singleChar == '_') { // check for letter or _
        tok.id = TOKEN_IDENTIFIER;
        // preceding characters may be letter,_, or a number
        for ; ((singleChar >= 'A') && (singleChar <= 'Z')) || ((singleChar >= 'a') && (singleChar <= 'z')) || (singleChar == '_') || ((singleChar >= '0') && (singleChar <= '9')); singleChar = GetCharWrapped() {
            tmp_TokAppendStr(singleChar);
        }
        // save the last read character for the next GetNextToken() cycle
        tok.nextChar = singleChar;
        done = 1;
    }

    // string "..."
    if (done != 1) && (singleChar == '"') {
        tok.id = TOKEN_STRING;        
        for singleChar = GetCharWrapped(); (singleChar != '"') && (singleChar > 31) && (singleChar < 127);singleChar = GetCharWrapped() {
            tmp_TokAppendStr(singleChar);
        }
        if singleChar != '"' {
            ScanErrorString("String not closing.");
        }
        done = 1;
    }

    // Single Quoted Character
    if (done != 1) && singleChar == 39 {
        singleChar = GetCharWrapped();
        if (singleChar != 39) && (singleChar > 31) && (singleChar < 127) {
            tok.id = TOKEN_INTEGER;
            tok.intValue = libgogo.ToIntFromByte(singleChar);
        } else {
            ScanErrorString("Unknown character.");
        }
        singleChar = GetCharWrapped();
        if singleChar != 39 {
            ScanErrorString("Only single characters allowed. Use corresponding integer for special characters.");
        }
        done = 1;
    }

    // left brace (
    if (done != 1) && singleChar == '(' {
        tok.id = TOKEN_LBRAC;
        done = 1;
    }

    // right brace )
    if (done != 1) && singleChar == ')' {
        tok.id = TOKEN_RBRAC;
        done = 1;
    }

    // left square bracket [
    if (done != 1) && singleChar == '[' {
        tok.id = TOKEN_LSBRAC;
        done = 1;    
    }
    
    // right square bracket ]
    if (done != 1) && singleChar == ']' {
        tok.id = TOKEN_RSBRAC;
        done = 1;
    }

    // integer
    if (done != 1) && (singleChar > 47) && (singleChar < 58) {
        numBuf = "";
        
        for ; (singleChar > 47) && (singleChar < 58) ; singleChar = GetCharWrapped() {
            libgogo.CharAppend(&numBuf, singleChar);
        }

        tok.nextChar = singleChar;  
        tok.id = TOKEN_INTEGER;
        tok.intValue = libgogo.StringToInt(numBuf);

        done = 1;
    }

    // Left curly bracket '{'
    if (done != 1) && (singleChar == '{') {
        tok.id = TOKEN_LCBRAC;
        done = 1;
    }
    
    // Right curly bracket '}'
    if (done != 1) && (singleChar == '}') {
        tok.id = TOKEN_RCBRAC;
        done = 1;
    }

    // Point '.'
    if (done != 1) && (singleChar == '.') {
        tok.id = TOKEN_PT;
        done = 1;
    }

    // Not ('!') or Not Equal ('!=')
    if (done != 1) && (singleChar == '!') {
        singleChar = GetCharWrapped();
        if singleChar == '=' {
            tok.id = TOKEN_NOTEQUAL;
        } else {
            tok.id = TOKEN_NOT;
            tok.nextChar = singleChar;
        }
        done = 1;
    }

    // Semicolon ';'
    if (done != 1) && (singleChar == ';') {
        tok.id = TOKEN_SEMICOLON;
        done = 1;
    }

    // Colon ','
    if (done != 1) && (singleChar == ',') {
        tok.id = TOKEN_COLON;
        done = 1;
    }

    // Assignment '=' or Equals comparison '=='
    if (done != 1) && (singleChar == '=') {
        singleChar = GetCharWrapped();
        if singleChar == '=' {
            tok.id = TOKEN_EQUALS;
        } else {
            tok.id = TOKEN_ASSIGN;
            tok.nextChar = singleChar;
        }
        done = 1;
    }

    // AND Relation '&&'
    if (done != 1) && (singleChar == '&') {
        singleChar = GetCharWrapped();
        if singleChar == '&' {
            tok.id = TOKEN_REL_AND;
        } else {
            tok.id = TOKEN_OP_ADR;
            tok.nextChar = singleChar;
        }
        done = 1;
    }

    // OR Relation '||'
    if (done != 1) && (singleChar == '|') {
        singleChar = GetCharWrapped();
        if singleChar == '|' {
            tok.id = TOKEN_REL_OR;
        } else {    
            ScanErrorString("No binary OR (|) supported. Only ||.");
        }
        done = 1;
    } 

    // Greater and Greater-Than relation
    if (done != 1) && (singleChar == '>') {
        singleChar = GetCharWrapped();
        if singleChar == '=' {
            tok.id = TOKEN_REL_GTOE;
        } else {
            tok.id = TOKEN_REL_GT;
            tok.nextChar = singleChar;
        }            
        done = 1;
    }     

    // Less and Less-Than relation
    if (done != 1) && (singleChar == '<') {
        singleChar = GetCharWrapped();
        if singleChar == '=' {
            tok.id = TOKEN_REL_LTOE;
        } else {
            tok.id = TOKEN_REL_LT;
            tok.nextChar = singleChar;
        }            
        done = 1;
    }    

    if (done != 1) && (singleChar == '+') {
        tok.id = TOKEN_ARITH_PLUS;
        done = 1;
    }

    if (done != 1) && (singleChar == '-') {
        tok.id = TOKEN_ARITH_MINUS;
        done = 1;
    }

    if (done != 1) && (singleChar == '*') {
        tok.id = TOKEN_ARITH_MUL;
        done = 1;
    }

    if (done != 1) && (singleChar == '/') {
        tok.id = TOKEN_ARITH_DIV;
        done = 1;
    }

    if (done != 1) {
        ScanErrorChar(singleChar);
    }
}


//
// GetNextToken should be called by the parser. It bascially fetches the next
// token by calling GetNextTokenRaw() and filters the identifiers for known
// keywords.
//
func GetNextToken() {
    GetNextTokenRaw();

    // Convert identifier to keyworded tokens
    if tok.id == TOKEN_IDENTIFIER {
        if libgogo.StringCompare("if",tok.strValue) == 0 {
            tok.id = TOKEN_IF;
        }
        if libgogo.StringCompare("else",tok.strValue) == 0 {
            tok.id = TOKEN_ELSE;
        }
        if libgogo.StringCompare("for",tok.strValue) == 0 {
            tok.id = TOKEN_FOR;
        }
        if libgogo.StringCompare("type",tok.strValue) == 0 {
            tok.id = TOKEN_TYPE;
        }
        if libgogo.StringCompare("const",tok.strValue) == 0 {
            tok.id = TOKEN_CONST;
        }
        if libgogo.StringCompare("var",tok.strValue) == 0 {
            tok.id = TOKEN_VAR;
        }
        if libgogo.StringCompare("struct", tok.strValue) == 0 {
            tok.id = TOKEN_STRUCT;
        }
        if libgogo.StringCompare("return", tok.strValue) == 0 {
            tok.id = TOKEN_RETURN;
        }
        if libgogo.StringCompare("func", tok.strValue) == 0 {
            tok.id = TOKEN_FUNC;
        }
        if libgogo.StringCompare("import", tok.strValue) == 0 {
            tok.id = TOKEN_IMPORT;
        }
        if libgogo.StringCompare("package", tok.strValue) == 0 {
            tok.id = TOKEN_PACKAGE;
        }
    }

    tok.nextToken = 0;
}

//
// Debugging and temporary functions
//

func debugToken(tok *Token) {
    libgogo.PrintString("---------------------\n");
    libgogo.PrintString("Token Id: ");
    libgogo.PrintNumber(tok.id);
    libgogo.PrintString("\n");
    if (tok.id == TOKEN_IDENTIFIER) || (tok.id == TOKEN_STRING) {
        libgogo.PrintString("Stored string: ");
        libgogo.PrintString(tok.strValue);
        libgogo.PrintString("\n");
    }
    if tok.id == TOKEN_INTEGER {
        libgogo.PrintString("Stored integer: ");
        libgogo.PrintNumber(tok.intValue);
        libgogo.PrintString("\n");
    }
}

func tmp_TokAppendStr(b byte) {
    libgogo.CharAppend(&tok.strValue, b);
}
