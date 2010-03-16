// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"
import "fmt"

//
// Token struct holding the relevant data of a parsed token.
//
type Token struct {
    id uint64; // The id. Is one of TOKEN_*
    /* value storing the integer value if the token is TOKEN_INTEGER */
    intValue uint64;
    /* Value that should be used instead of byte arrays */
    newValue string;

    nextChar byte; // Sometime the next char is already read. It is stored here to be re-assigned in the next GetNextToken() round
};

func GetNextTokenRaw(fd uint64, tok *Token) {
    var singleChar byte; // Byte holding the last read value
    /* 
     * Flag indicating whether we are in a comment.
     * 0 for no comment
     * 1 for a single line comment 
     * 2 for a multi line comment
     */
    var inComment uint64;
    var done uint64; // Flag indicating whether a cycle (Token) is finsihed 
    var spaceDone uint64; // Flag indicating whether an abolishment cycle is finished 

    // Initialize variables
    done = 0;
    spaceDone = 0;
    inComment = 0;  

    // If the old Token had to read the next char (and stored it), we can now
    // get it back
    if tok.nextChar == 0 {       
        singleChar=libgogo.GetChar(fd)
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
            tmp_StrAppend(tok.newValue,singleChar);
        }
        // save the last read character for the next GetNextToken() cycle
        tok.nextChar = singleChar;
        done = 1;
    }

    // string "..."
    if (done != 1) && singleChar == '"' {
        tok.id = TOKEN_STRING;        
        for singleChar = libgogo.GetChar(fd); singleChar != '"' &&singleChar > 31 && singleChar < 127;singleChar = libgogo.GetChar(fd) {
            tmp_StrAppend(tok.newValue,singleChar);
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
            tok.intValue = tmp_toInt(singleChar);
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

func GetNextToken(fd uint64, tok *Token) {
    GetNextTokenRaw(fd,tok)

    // Convert identifier to keyworded tokens
    if tok.id == TOKEN_IDENTIFIER {
        if tmp_StrCmp("if",tok.newValue) == 0 {
            tok.id = TOKEN_IF;
        }
        if tmp_StrCmp("for",tok.newValue) == 0 {
            tok.id = TOKEN_FOR;
        }
        if tmp_StrCmp("type",tok.newValue) == 0 {
            tok.id = TOKEN_TYPE;
        }
        if tmp_StrCmp("const",tok.newValue) == 0 {
            tok.id = TOKEN_CONST;
        }
        if tmp_StrCmp("var",tok.newValue) == 0 {
            tok.id = TOKEN_VAR;
        }
        if tmp_StrCmp("struct", tok.newValue) == 0 {
            tok.id = TOKEN_STRUCT;
        }
        if tmp_StrCmp("return", tok.newValue) == 0 {
            tok.id = TOKEN_RETURN;
        }
        if tmp_StrCmp("func", tok.newValue) == 0 {
            tok.id = TOKEN_FUNC;
        }
        if tmp_StrCmp("import", tok.newValue) == 0 {
            tok.id = TOKEN_IMPORT;
        }
        if tmp_StrCmp("package", tok.newValue) == 0 {
            tok.id = TOKEN_PACKAGE;
        }
    }
}

func debugToken (tok Token) {
    /*
    libgogo.PrintString("Token Id: ");
    libgogo.PrintNumber(tok.id);
    libgogo.PrintString("\n");

    if tok.id == TOKEN_IDENTIFIER || tok.id == TOKEN_STRING {
        libgogo.PrintString("Identifier/String value: ");
        libgogo.PrintByteBuf(tok.value);
        libgogo.PrintString("\n");
    }*/

}

// Temporary test function
func ScannerTest(fd uint64) {  
    var tok Token;

    tok.id = 0;
    tok.nextChar = 0;

    for GetNextToken(fd,&tok); tok.id != TOKEN_EOS; GetNextToken(fd,&tok) {
        fmt.Printf("%d\n",tok.id);
    }
}

// libgogo ...
func tmp_StrAppend(str string, b byte) {
    str += string(b);
}

// libgogo ...
func tmp_StrLen(str string) int {
    return len(str);
}

// libgogo ...
func tmp_StrCmp(str1 string, str2 string) uint64 {
    var ret uint64;    
    if str1 == str2 {
        ret = 0;
    } else {
        ret = 1;
    }
    return ret;
}

// libgogo ...
func tmp_toInt(b byte) uint64 {
    return uint64(b);
}
