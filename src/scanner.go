// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "syscall"
import "unsafe"
import "fmt"

type Token struct {
    id int;
    value [255]byte;
    value_len int;
};

//
// Function retrieving the next byte of a file pointer.
// Should/Must be replaced by our own library functions.
//
func getNextChar(fd int) byte {
	var b [1]byte;
	_, _, _ = syscall.Syscall(syscall.SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(&b)), 1);
	return b[0];
}

//
// Function getting the next token.
//
func GetNextToken(fd int) Token {
    var singleChar byte;
    var inComment int;
    var done int;
    var tok Token;

    tok.id = 0;
    tok.value_len = 0;
    done = 0;

    for singleChar=getNextChar(fd);done != 1; {
        if singleChar == '/' { 
            // Handling comments

            singleChar = getNextChar(fd);
            if singleChar == '/' {
                inComment = 1;
            } else {
                // TODO: ERROR /<char> not in language
            }
        } else {
            // Handling rest

            if singleChar == '\n' {
                if inComment == 1 {
                    inComment = 0;
                } else {
                    if tok.id != 0 {
                        done = 1;
                    }
                }
            } else {
                if inComment == 0 {
                    if singleChar == ' ' {
                        done = 1;
                    } else {

                        // TODO: Seperate different tokens.
                        tok.id = 1;
                        tok.value[tok.value_len] = singleChar;
                        tok.value_len = tok.value_len + 1;

                    }
                }                                
            }
        }

        if done == 0 {        
            singleChar=getNextChar(fd)
        }
    }

    tok.value[tok.value_len] = 0;

    return tok;
}

func tmp_print(b [255]byte) {
    var i int;
    for i=0;b[i] != 0;i=i+1 {
        fmt.Printf("%c",b[i]);
    }
    fmt.Printf("\n");
}

func scanner_test(fd int) {
    var tok Token;
    tok = GetNextToken(fd);
    if tok.id == 1 {
        fmt.Printf("Found 'identifier'\n");
        tmp_print(tok.value); 
    }
    tok = GetNextToken(fd);
    if tok.id == 1 {
        fmt.Printf("Found 'identifier'\n");
        tmp_print(tok.value); 
    }
}
