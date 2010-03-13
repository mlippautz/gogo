// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "syscall"
import "unsafe"
import "fmt"

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
// Internal function handling the retrieval of the next valid Symbol. A symbol
// is basically a non-whitespace character that is not a comment.
// 
func getNextSymbol(fd int) byte {
    var singleChar byte;
    var inComment int;
    var done int;

    done = 0;
    inComment = 0;

    for singleChar=getNextChar(fd);done != 1; {
        if singleChar == '/' {
            singleChar = getNextChar(fd);
            if singleChar == '/' {
                inComment = 1;
            } else {
                // TODO: ERROR /<char> not in language
            }
        } else {
            if singleChar == '\n' {
                if inComment == 1 {
                    inComment = 0;
                }
            } else {
                if singleChar != ' ' {
                    if inComment == 0 {
                        done = 1;
                    }
                }
            }
        }

        if done == 0 {        
            singleChar=getNextChar(fd)
        }
    }

    return singleChar;
}

func GetNextToken(fd int) int {
    var sym byte;
    sym = getNextSymbol(fd);
    fmt.Printf("%c\n",sym);
    sym = getNextSymbol(fd);
    fmt.Printf("%c\n",sym);
    sym = getNextSymbol(fd);
    fmt.Printf("%c\n",sym);
    return 0;
}

