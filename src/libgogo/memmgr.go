// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo memory manager functions
//

package libgogo

var INC_SIZE uint64 = 1024;

var bump_ptr uint64 = 0;
var end_ptr uint64 = 0;

func Brk(brk uint64) uint64;

func GetBrk() uint64;

func TestMem(address uint64) uint64;

func InitMemoryManager() {
    end_ptr = GetBrk();
    if end_ptr == 0 {
        ExitError("GetBrk failed\n", 127);
    }
    bump_ptr = end_ptr + 1; //First useable address in new memory (will be allocated in next line)
    MoreMemory();
}

func MoreMemory() {
    var errno uint64;

    end_ptr = end_ptr + INC_SIZE;
    errno = Brk(end_ptr);
    if errno != 0 {
        PrintString("Brk failed while allocating ");
        PrintNumber(INC_SIZE);
        PrintString(" bytes. Errno: ");
        PrintNumber(errno);
        ExitError("\n", 127);
    }
    errno = TestMem(end_ptr - INC_SIZE + 1); //Sanity check for newly allocated memory
    if errno != 0 {
        PrintString("Failed to write to newly allocated memory at ");
        PrintNumber(end_ptr - INC_SIZE + 1);
        ExitError("\n", 127);
    }
}

func Alloc(size uint64) uint64 {
    var addr uint64;

    if (bump_ptr == 0) && (end_ptr == 0) { //First call
        InitMemoryManager();
    }
    for ; bump_ptr + size >= end_ptr; { //Ensure there is enough memory
        MoreMemory();
    }
 
    addr = bump_ptr;
    bump_ptr = bump_ptr + size;
    return addr;
}

func GetFreeMemory() uint64 {
    return end_ptr - bump_ptr;
}
