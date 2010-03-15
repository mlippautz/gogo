// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

//
// Set of recognized tokens 
//
const TOKEN_IDENTIFIER uint64 = 1;  // Identifier
const TOKEN_STRING = 2;             // String using "..."
const TOKEN_EOS = 3;                // End of Scan
const TOKEN_LBRAC = 4;              // Left bracket '('
const TOKEN_RBRAC = 5;              // Right bracket ')'
const TOKEN_LSBRAC = 6;             // Left square bracket '['
const TOKEN_RSBRAC = 7;             // Right square bracket ']'
const TOKEN_INTEGER = 8;            // Integer number
const TOKEN_LCBRAC = 9;             // Left curly bracket '{'
const TOKEN_RCBRAC = 10;            // Right curly bracket '}'
const TOKEN_PT = 11;                // Point '.'
const TOKEN_NOT = 12;               // Single not '!'
const TOKEN_NOTEQUAL = 13;          // Comparison, not equal '!='
const TOKEN_SEMICOLON = 14;         // Semi-colon ';'
const TOKEN_COLON = 15;             // Colon ','
const TOKEN_ASSIGN = 16;            // Assignment '='
const TOKEN_EQUALS = 17;            // Equal comparison '=='
const TOKEN_CHAR = 18;              // Single Quoted Character
const TOKEN_REL_AND = 19;           // AND Relation
const TOKEN_REL_OR = 20;            // OR Relation
const TOKEN_REL_GTOE = 21;          // Greather-Than or Equal
const TOKEN_REL_GT = 22;            // Greather-Than
const TOKEN_REL_LTOE = 23;          // Less-Than or Equal
const TOKEN_REL_LT = 24;            // Less-Than
const TOKEN_ARITH_PLUS = 25;        // Arith. Plus
const TOKEN_ARITH_MINUS = 26;       // Arith. Minus
const TOKEN_ARITH_MUL = 27;         // Arith. Multiplication
const TOKEN_ARITH_DIV = 28;         // Arith. Division

// Advanced tokens, that are generated in the 2nd step from identifiers
const TOKEN_FOR = 29;
const TOKEN_IF = 30;

//
// Token struct holding the relevant data of a parsed token.
//
type Token struct {
    id uint64; // The id. Is one of TOKEN_*
    value [255]byte; // If the id requires a value to be stored, it is found here 
    value_len uint64; // Length of the value stored in `value`
    intValue uint64;

    nextChar byte; // Sometime the next char is already read. It is stored here to be re-assigned in the next GetNextToken() round
};

/*
 * Function getting the next raw token. May contain token that must be converted
 * before the parser can work with them.
 */
func GetNextTokenRaw(fd uint64, oldToken Token) Token {
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
    var newToken Token; // The new token that is returned

    // Initialize variables
    done = 0;
    spaceDone = 0;
    newToken.id = 0;
    newToken.value_len = 0;
    newToken.nextChar = 0; 
    inComment = 0;   

    // If the old Token had to read the next char (and stored it), we can now
    // get it back
    if oldToken.nextChar == 0 {       
        singleChar=libgogo.GetChar(fd)
    } else {
        singleChar = oldToken.nextChar;
    }

    // check if it is a valid read, or an EOF
    if singleChar == 0 {
        newToken.id = TOKEN_EOS;
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
                newToken.id = TOKEN_EOS;
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
        newToken.id = TOKEN_IDENTIFIER;
        // preceding characters may be letter,_, or a number
        for ; (singleChar >= 'A' && singleChar <= 'Z') || (singleChar >= 'a' && singleChar <= 'z') || singleChar == '_' || (singleChar >= '0' && singleChar <= '9'); singleChar = libgogo.GetChar(fd) {
            newToken.value[newToken.value_len] = singleChar;
            newToken.value_len = newToken.value_len +1;
        }
        newToken.value[newToken.value_len] = 0;
        // save the last read character for the next GetNextToken() cycle
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
            libgogo.ExitError(">> Scanner: String not closing. Exiting.",1);
        }
        done = 1;
    }

    // Single Quoted Character
    // Needs to be converted later on
    if (done != 1) && singleChar == 39 {
        newToken.id = TOKEN_CHAR;
        for singleChar = libgogo.GetChar(fd); singleChar != 39 && singleChar > 31 && singleChar < 127;singleChar = libgogo.GetChar(fd) {
            newToken.value[newToken.value_len] = singleChar;
            newToken.value_len = newToken.value_len +1;
        }        
        newToken.value[newToken.value_len] = 0;
        if singleChar != 39 {
            libgogo.ExitError(">> Scanner: Character not closing. Exiting.",1);
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

    // Left curly bracket '{'
    if (done != 1) && singleChar == '{' {
        newToken.id = TOKEN_LCBRAC;
        done = 1;
    }
    
    // Right curly bracket '}'
    if (done != 1) && singleChar == '}' {
        newToken.id = TOKEN_RCBRAC;
        done = 1;
    }

    // Point '.'
    if (done != 1) && singleChar == '.' {
        newToken.id = TOKEN_PT;
        done = 1;
    }

    // Not ('!') or Not Equal ('!=')
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

    // Semicolon ';'
    if (done != 1) && singleChar == ';' {
        newToken.id = TOKEN_SEMICOLON;
        done = 1;
    }

    // Colon ','
    if (done != 1) && singleChar == ',' {
        newToken.id = TOKEN_COLON;
        done = 1;
    }

    // Assignment '=' or Equals comparison '=='
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

    // AND Relation '&&'
    if (done != 1) && singleChar == '&' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '&' {
            newToken.id = TOKEN_REL_AND;
        } else {
            libgogo.ExitError(">> Scanner: No address operator supported.",1);
        }
        done = 1;
    }

    // OR Relation '||'
    if (done != 1) && singleChar == '|' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '|' {
            newToken.id = TOKEN_REL_OR;
        } else {    
            libgogo.ExitError(">> Scanner: No binary OR (|) supported. Only ||.",1);
        }
        done = 1;
    } 

    // Greater and Greater-Than relation
    if (done != 1) && singleChar == '>' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '=' {
            newToken.id = TOKEN_REL_GTOE;
        } else {
            newToken.id = TOKEN_REL_GT;
            newToken.nextChar = singleChar;
        }            
        done = 1;
    }     

    // Less and Less-Than relation
    if (done != 1) && singleChar == '<' {
        singleChar = libgogo.GetChar(fd);
        if singleChar == '=' {
            newToken.id = TOKEN_REL_LTOE;
        } else {
            newToken.id = TOKEN_REL_LT;
            newToken.nextChar = singleChar;
        }            
        done = 1;
    }    

    if (done != 1) && singleChar == '+' {
        newToken.id = TOKEN_ARITH_PLUS;
        done = 1;
    }

    if (done != 1) && singleChar == '-' {
        newToken.id = TOKEN_ARITH_MINUS;
        done = 1;
    }

    if (done != 1) && singleChar == '*' {
        newToken.id = TOKEN_ARITH_MUL;
        done = 1;
    }

    if (done != 1) && singleChar == '/' {
        newToken.id = TOKEN_ARITH_DIV;
        done = 1;
    }

    if done != 1 {
        
        libgogo.PrintString(">> Scanner: Unkown char '");
        libgogo.PrintChar(singleChar);
        libgogo.PrintString("'. ");
        libgogo.ExitError("Exiting.",1);
    }

    return newToken;
}

func GetNextToken(fd uint64, oldToken Token) Token {
    var newToken Token;
    newToken = GetNextTokenRaw(fd,oldToken)
    
    // Convert integer from byte array to integer value
    if newToken.id == TOKEN_INTEGER {
        newToken.intValue = libgogo.ByteBufToInt(newToken.value,newToken.value_len);
    }

    // Convert single quoted characters to integer
    if newToken.id == TOKEN_CHAR {
        if newToken.value_len != 1 {
            libgogo.ExitError (">> Scanner: Only single characters are supported!",1);
        } else {
            newToken.id = TOKEN_INTEGER;
            newToken.intValue = libgogo.ToIntFromByte(newToken.value[0]);
        }
    }

    // Convert identifier to keyworded tokens (if possible)
    // <TODO>

    return newToken;
}

func debugToken (tok Token) {
    libgogo.PrintString("Token Id: ");
    libgogo.PrintNumber(tok.id);
    libgogo.PrintString("\n");

    if tok.id == TOKEN_IDENTIFIER || tok.id == TOKEN_STRING {
        libgogo.PrintString("Identifier/String value: ");
        libgogo.PrintByteBuf(tok.value);
        libgogo.PrintString("\n");
    }
}

// Temporary test function
func ScannerTest(fd uint64) {  
    var tok Token;
    tok.id = 0;
    tok.nextChar = 0;

    for tok = GetNextToken(fd,tok); tok.id != TOKEN_EOS; tok = GetNextToken(fd,tok) {
        debugToken(tok);
    }
}
