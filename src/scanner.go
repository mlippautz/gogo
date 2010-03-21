// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// This file holds the basic scanning routines that separate a source file
// into the various tokens

package main

import "./libgogo/_obj/libgogo"
import "fmt"

// Token struct holding the relevant data of a parsed token.
type Token struct {
    id uint64; // The id. Is one of TOKEN_*
    intValue uint64; // value storing the integer value if the token is TOKEN_INTEGER
    strValue string; // Value storing the token string if the token is TOKEN_STRING or TOKEN_IDENTIFIER
    nextChar byte; // Sometime the next char is already read. It is stored here to be re-assigned in the next GetNextToken() round
};

func GetNextTokenRaw(fd uint64, tok *Token) {
    var singleChar byte; // Byte holding the last read value
    // Flag indicating whether we are in a comment.
    // 0 for no comment
    // 1 for a single line comment 
    // 2 for a multi line comment
    var inComment uint64;
    var done uint64; // Flag indicating whether a cycle (Token) is finsihed 
    var spaceDone uint64; // Flag indicating whether an abolishment cycle is finished 

    // Initialize variables
    done = 0;
    spaceDone = 0;
    inComment = 0;  

    tok.strValue = "";

    // If the previous cycle had to read the next char (and stored it), it is 
    // now used as first read
    if tok.nextChar == 0 {       
        singleChar = libgogo.GetChar(fd)
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
                singleChar = libgogo.GetChar(fd); 
                if singleChar == '/' {
                    // we are in a single line comment (until newline is found)
                    inComment = 1;
                } else {
                    if singleChar == '*' {
                        // we are in a multiline comment (until ending is found)
                        inComment = 2;
                    } else {
                        libgogo.ExitError(">> Scanner: Unkown character combination for comments. Exiting.",1);
                    }
                }
            }
        } 

        // check whether a multi-line comment is ending
        if singleChar == '*' {
            singleChar = libgogo.GetChar(fd);
            if singleChar == '/' {
                if inComment == 2 {
                    inComment = 0;
                    singleChar = libgogo.GetChar(fd);
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
        if singleChar != ' ' && singleChar != 9 && singleChar != 10 {
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
            singleChar=libgogo.GetChar(fd);
        }
    }

    //
    // Actual scanning part starts here
    //

    // Catch identifiers
    // identifier = letter { letter | digit }.
    if (done != 1) && (singleChar >= 'A' && singleChar <= 'Z') || (singleChar >= 'a' && singleChar <= 'z') || singleChar == '_' { // check for letter or _
        tok.id = TOKEN_IDENTIFIER;
        // preceding characters may be letter,_, or a number
        for ; (singleChar >= 'A' && singleChar <= 'Z') || (singleChar >= 'a' && singleChar <= 'z') || singleChar == '_' || (singleChar >= '0' && singleChar <= '9'); singleChar = libgogo.GetChar(fd) {
            tmp_TokAppendStr(tok,singleChar);
        }
        // save the last read character for the next GetNextToken() cycle
        tok.nextChar = singleChar;
        done = 1;
    }

    // string "..."
    if (done != 1) && singleChar == '"' {
        tok.id = TOKEN_STRING;        
        for singleChar = libgogo.GetChar(fd); singleChar != '"' &&singleChar > 31 && singleChar < 127;singleChar = libgogo.GetChar(fd) {
            tmp_TokAppendStr(tok,singleChar);
        }
        if singleChar != '"' {
            libgogo.ExitError(">> Scanner: String not closing. Exiting.",1);
        }
        done = 1;
    }

    // Single Quoted Character
    if (done != 1) && singleChar == 39 {
        singleChar = libgogo.GetChar(fd);
        if singleChar != 39 && singleChar > 31 && singleChar < 127 {
            tok.id = TOKEN_INTEGER;
            tok.intValue = libgogo.ToIntFromByte(singleChar);
        } else {
            libgogo.ExitError(">> Scanner: Unknown character. Exiting.",1);
        }
        singleChar = libgogo.GetChar(fd);
        if singleChar != 39 {
            libgogo.ExitError(">> Scanner: Only single characters allowed. Use corresponding integer for special characters. Exiting.",1);
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
    if (done != 1) && singleChar > 47 && singleChar < 58 {
        var byteBuf [255]byte;
        var i uint64;
        
        for i = 0; singleChar > 47 && singleChar < 58 ; singleChar = libgogo.GetChar(fd) {
            byteBuf[i] = singleChar;
            i = i +1;
        }

        tok.nextChar = singleChar;  
        tok.id = TOKEN_INTEGER;
        tok.intValue = libgogo.ByteBufToInt(byteBuf,i);

        done = 1;
    }

    // Left curly bracket '{'
    if (done != 1) && singleChar == '{' {
        tok.id = TOKEN_LCBRAC;
        done = 1;
    }
    
    // Right curly bracket '}'
    if (done != 1) && singleChar == '}' {
        tok.id = TOKEN_RCBRAC;
        done = 1;
    }

    // Point '.'
    if (done != 1) && singleChar == '.' {
        tok.id = TOKEN_PT;
        done = 1;
    }

    // Not ('!') or Not Equal ('!=')
    if (done != 1) && singleChar == '!' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '=' {
            tok.id = TOKEN_NOTEQUAL;
        } else {
            tok.id = TOKEN_NOT;
            tok.nextChar = singleChar;
        }
        done = 1;
    }

    // Semicolon ';'
    if (done != 1) && singleChar == ';' {
        tok.id = TOKEN_SEMICOLON;
        done = 1;
    }

    // Colon ','
    if (done != 1) && singleChar == ',' {
        tok.id = TOKEN_COLON;
        done = 1;
    }

    // Assignment '=' or Equals comparison '=='
    if (done != 1) && singleChar == '=' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '=' {
            tok.id = TOKEN_EQUALS;
        } else {
            tok.id = TOKEN_ASSIGN;
            tok.nextChar = singleChar;
        }
        done = 1;
    }

    // AND Relation '&&'
    if (done != 1) && singleChar == '&' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '&' {
            tok.id = TOKEN_REL_AND;
        } else {
            tok.id = TOKEN_OP_ADR;
            tok.nextChar = singleChar;
        }
        done = 1;
    }

    // OR Relation '||'
    if (done != 1) && singleChar == '|' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '|' {
            tok.id = TOKEN_REL_OR;
        } else {    
            libgogo.ExitError(">> Scanner: No binary OR (|) supported. Only ||.",1);
        }
        done = 1;
    } 

    // Greater and Greater-Than relation
    if (done != 1) && singleChar == '>' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '=' {
            tok.id = TOKEN_REL_GTOE;
        } else {
            tok.id = TOKEN_REL_GT;
            tok.nextChar = singleChar;
        }            
        done = 1;
    }     

    // Less and Less-Than relation
    if (done != 1) && singleChar == '<' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '=' {
            tok.id = TOKEN_REL_LTOE;
        } else {
            tok.id = TOKEN_REL_LT;
            tok.nextChar = singleChar;
        }            
        done = 1;
    }    

    if (done != 1) && singleChar == '+' {
        tok.id = TOKEN_ARITH_PLUS;
        done = 1;
    }

    if (done != 1) && singleChar == '-' {
        tok.id = TOKEN_ARITH_MINUS;
        done = 1;
    }

    if (done != 1) && singleChar == '*' {
        tok.id = TOKEN_ARITH_MUL;
        done = 1;
    }

    if (done != 1) && singleChar == '/' {
        tok.id = TOKEN_ARITH_DIV;
        done = 1;
    }

    if done != 1 {
        
        libgogo.PrintString(">> Scanner: Unkown char '");
        libgogo.PrintChar(singleChar);
        libgogo.PrintString("'. ");
        libgogo.ExitError("Exiting.",1);
    }
}


//
// GetNextToken should be called by the parser. It bascially fetches the next
// token by calling GetNextTokenRaw() and filters the identifiers for known
// keywords.
//
func GetNextToken(fd uint64, tok *Token) {
    GetNextTokenRaw(fd,tok)

    // Convert identifier to keyworded tokens
    if tok.id == TOKEN_IDENTIFIER {
        if libgogo.StringCompare("if",tok.strValue) != 0 {
            tok.id = TOKEN_IF;
       }
        if libgogo.StringCompare("for",tok.strValue) != 0 {
            tok.id = TOKEN_FOR;
        }
        if libgogo.StringCompare("type",tok.strValue) != 0 {
            tok.id = TOKEN_TYPE;
        }
        if libgogo.StringCompare("const",tok.strValue) != 0 {
            tok.id = TOKEN_CONST;
        }
        if libgogo.StringCompare("var",tok.strValue) != 0 {
            tok.id = TOKEN_VAR;
        }
        if libgogo.StringCompare("struct", tok.strValue) != 0 {
            tok.id = TOKEN_STRUCT;
        }
        if libgogo.StringCompare("return", tok.strValue) != 0 {
            tok.id = TOKEN_RETURN;
        }
        if libgogo.StringCompare("func", tok.strValue) != 0 {
            tok.id = TOKEN_FUNC;
        }
        if libgogo.StringCompare("import", tok.strValue) != 0 {
            tok.id = TOKEN_IMPORT;
        }
        if libgogo.StringCompare("package", tok.strValue) != 0 {
            tok.id = TOKEN_PACKAGE;
        }
    }
}

//
// Debugging and temporary functions
//

func debugToken(tok *Token) {
    libgogo.PrintString("---------------------\n");
    libgogo.PrintString("Token Id: ");
    libgogo.PrintNumber(tok.id);
    libgogo.PrintString("\n");
    if tok.id == TOKEN_IDENTIFIER || tok.id == TOKEN_STRING {
        libgogo.PrintString("Stored string: ");
        fmt.Printf(tok.strValue);
        libgogo.PrintString("\n");
    }
    if tok.id == TOKEN_INTEGER {
        libgogo.PrintString("Stored integer: ");
        libgogo.PrintNumber(tok.intValue);
        libgogo.PrintString("\n");
    }
}

// Temporary test function
func ScannerTest(fd uint64) {  
    var tok Token;

    tok.id = 0;
    tok.nextChar = 0;

    for GetNextToken(fd,&tok); tok.id != TOKEN_EOS; GetNextToken(fd,&tok) {
        debugToken(&tok);
    }
}

func tmp_TokAppendStr(tok *Token, b byte) {
    libgogo.StringAppend(&tok.strValue, b);
}
