// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "fmt"
import "./libgogo/_obj/libgogo"


const TOKEN_IDENTIFIER uint64 = 1;
const TOKEN_STRING = 2;
const TOKEN_EOS = 3; // end of scan
const TOKEN_LBRAC = 4;
const TOKEN_RBRAC = 5;
const TOKEN_LSBRAC = 6;
const TOKEN_RSBRAC = 7;
const TOKEN_INTEGER = 8;
const TOKEN_LCBRAC = 9;
const TOKEN_RCBRAC = 10;
const TOKEN_PT = 11;
const TOKEN_NOT = 12;
const TOKEN_NOTEQUAL = 13;
const TOKEN_SEMICOLON = 14;
const TOKEN_COLON = 15;
const TOKEN_ASSIGN = 16;
const TOKEN_EQUALS = 17;


type Token struct {
    id uint64;
    value [255]byte;
    value_len uint64;

    // also need a second token, because we may scan it already when terminating the first one
    nextChar byte;    
};

func tmp_error ( s string) {
    fmt.Printf("%s\n",s);
    libgogo.Exit(1);
}

//
// Function getting the next token.
//
func GetNextToken(fd uint64, oldToken Token) Token {
    var singleChar byte;
    var inComment uint64;
    var done uint64;
    var space_done uint64;
    var newToken Token;

    done = 0;
    space_done =0; 
    newToken.id = 0;
    newToken.value_len = 0;
    newToken.nextChar = 0; 
    inComment = 0;   

    if oldToken.nextChar == 0 {       
        singleChar=libgogo.GetChar(fd)
    } else {
        singleChar = oldToken.nextChar;
    }

    // check if it is a valid read
    if singleChar == 0 {
        newToken.id = TOKEN_EOS;
        done = 1;
        space_done = 1;
    }

    // skip blank, newlines and comments
    // after function singleChar contains the latest valid char
    for ;space_done != 1; {
        if singleChar == '/' {
            if inComment == 0 {
                singleChar = libgogo.GetChar(fd); 
                if singleChar == '/' {
                    inComment = 1;
                } else {
                    tmp_error(">> Scanner: Unkown character combination for comments. Exiting.");
                }
            }
        } else {
            
            // newline
            if singleChar == 10 {
                if inComment == 1 {
                    inComment = 0;
                } 
            } else {
                // skip spaces and tabs
                if singleChar != ' ' && singleChar != 9 {
                    if inComment == 0 {
                        space_done = 1;
                    } 
                    if singleChar == 0 {
                        space_done = 1;
                        newToken.id = TOKEN_EOS;
                        done = 1;
                    }   
                }
            }
        }
        if space_done == 0 {        
            singleChar=libgogo.GetChar(fd);
        }
    }

    // get identifiers
    if (done != 1) && (singleChar > 64 && singleChar < 91) || (singleChar > 96 && singleChar < 122) || singleChar == '_' {     
        newToken.id = TOKEN_IDENTIFIER;
        for ; (singleChar > 64 && singleChar < 91) || (singleChar > 96 && singleChar < 122) || singleChar == '_'  || (singleChar > 47 && singleChar < 58); singleChar = libgogo.GetChar(fd) {
            newToken.value[newToken.value_len] = singleChar;
            newToken.value_len = newToken.value_len +1;
        }
        newToken.value[newToken.value_len] = 0;
        newToken.nextChar = singleChar;
        done = 1;
    }

    // string "..."
    if (done != 1) && singleChar == '"' {
        newToken.id = TOKEN_STRING;        
        for singleChar = libgogo.GetChar(fd); singleChar != '"' &&singleChar > 31 && singleChar < 127;singleChar = libgogo.GetChar(fd) {
            newToken.value[newToken.value_len] = singleChar;
            newToken.value_len = newToken.value_len +1;
        }
        newToken.value[newToken.value_len] = 0;
        if singleChar != '"' {
            tmp_error(">> Scanner: String not closing. Exiting.");
        }
        done = 1;
    }

    // left brace (
    if (done != 1) && singleChar == '(' {
        newToken.id = TOKEN_LBRAC;
        done = 1;
    }

    // right brace )
    if (done != 1) && singleChar == ')' {
        newToken.id = TOKEN_RBRAC;
        done = 1;
    }

    // left square bracket [
    if (done != 1) && singleChar == '[' {
        newToken.id = TOKEN_LSBRAC;
        done = 1;    
    }
    
    // right square bracket ]
    if (done != 1) && singleChar == ']' {
        newToken.id = TOKEN_RSBRAC;
        done = 1;
    }

    // integer
    if (done != 1) && singleChar > 47 && singleChar < 58 {
        newToken.id = TOKEN_INTEGER;
        for ; singleChar > 47 && singleChar < 58 ; singleChar = libgogo.GetChar(fd) {
            newToken.value[newToken.value_len] = singleChar;
            newToken.value_len = newToken.value_len +1;
        }
        newToken.value[newToken.value_len] = 0
        newToken.nextChar = singleChar;  
        done = 1;
    }

    if (done != 1) && singleChar == '{' {
        newToken.id = TOKEN_LCBRAC;
        done = 1;
    }
    
    if (done != 1) && singleChar == '}' {
        newToken.id = TOKEN_RCBRAC;
        done = 1;
    }

    if (done != 1) && singleChar == '.' {
        newToken.id = TOKEN_PT;
        done = 1;
    }

    // not or not equal
    if (done != 1) && singleChar == '!' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '=' {
            newToken.id = TOKEN_NOTEQUAL;
        } else {
            newToken.id = TOKEN_NOT;
            newToken.nextChar = singleChar;
        }
        done = 1;
    }

    if (done != 1) && singleChar == ';' {
        newToken.id = TOKEN_SEMICOLON;
        done = 1;
    }

    if (done != 1) && singleChar == ',' {
        newToken.id = TOKEN_COLON;
        done = 1;
    }

    if (done != 1) && singleChar == '=' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '=' {
            newToken.id = TOKEN_EQUALS;
        } else {
            newToken.id = TOKEN_ASSIGN;
            newToken.nextChar = singleChar;
        }
        done = 1;
    }

    if done != 1 {
        fmt.Printf("'%c'\n",singleChar);
        tmp_error(">> Scanner: Unkown char detected. Exiting");
    }

    return newToken;
}

func tmp_print(tok Token) {
    var i int;
    fmt.Printf("Token Id: %d\n",tok.id);
    if tok.id == TOKEN_IDENTIFIER || tok.id == TOKEN_STRING {
        fmt.Printf("Identifier/String value: ");
        for i=0;tok.value[i] != 0;i=i+1 {
            fmt.Printf("%c",tok.value[i]);
        }
        fmt.Printf("\n");
    }
}

func scanner_test(fd uint64) {  
    var tok Token;
    tok.id = 0;
    tok.nextChar = 0;

    for tok = GetNextToken(fd,tok); tok.id != TOKEN_EOS; tok = GetNextToken(fd,tok) {
        tmp_print(tok);
    }
}
