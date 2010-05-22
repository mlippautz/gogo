// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo memory manager functions
//

package libgogo

//
// Default size increment in bytes when allocating more memory for the memory manager
//
var INC_SIZE uint64 = 1024; //Default 1 KB

//
// Pointers used for memory management
// Start pointer: Pointer to first address of memory manager in memory
// Bump pointer: Pointer to next free address in memory
// End pointer: Pointer to the end of the free memory available for use through the memory manager
//
var start_ptr uint64 = 0;
var bump_ptr uint64 = 0;
var end_ptr uint64 = 0;

//
// Sets the end address of the data segment to the one specified, thereby resizing it
// Brk returns 0 if the data segment could not be resized
// Implemented in assembler (see corresponding .s file)
//
func Brk(brk uint64) uint64;

//
// Gets the current end address of the data segment
// GetBrk returns 0 in case of an error
// Implemented in assembler (see corresponding .s file)
//
func GetBrk() uint64;

//
// Tests newly allocated memory by writing to the address given and verifying the result by re-reading it
// TestMem returns 0 in case of success; it is very improbable for the return value to be unequal 0 as an invalid address will most likely lead to a program crash before the function is able to return
// Implemented in assembler (see corresponding .s file)
//
func TestMem(address uint64) uint64;

//
// Initializes the memory manager by setting the end and bump pointer 
//
func InitMemoryManager() {
    start_ptr = GetBrk(); //Get current end of data segment
    if start_ptr == 0 { //Error check
        ExitError("GetBrk failed\n", 127);
    }
    start_ptr = start_ptr + 1; //First useable address in new memory (will be allocated below)
    end_ptr = start_ptr; //No free memory yet (will be allocated below)
    bump_ptr = end_ptr;
    MoreMemory(); //Allocate memory for the memory manager to work with
}

//
// Allocates more memory, more precisely the amount in bytes specified by INC_SIZE
//
func MoreMemory() {
    var errno uint64;

    end_ptr = end_ptr + INC_SIZE; //Increment end pointer
    errno = Brk(end_ptr); //Resize data segment
    if errno != 0 { //Error check
        PrintString("Brk failed while allocating ");
        PrintNumber(INC_SIZE);
        PrintString(" bytes. Errno: ");
        PrintNumber(errno);
        ExitError("\n", 127);
    }
    errno = TestMem(end_ptr - INC_SIZE + 1); //Sanity check for newly allocated memory
    if errno != 0 { //Error handling in case sanity check fails
        PrintString("Failed to write to newly allocated memory at ");
        PrintNumber(end_ptr - INC_SIZE + 1);
        ExitError("\n", 127);
    }
}

//
// Allocated size bytes of memory through the memory manager and returns the address of the newly allocated memory
//
func Alloc(size uint64) uint64 {
    var addr uint64;

    if (bump_ptr == 0) && (end_ptr == 0) { //First call => initialization required
        InitMemoryManager();
    }
    for ; bump_ptr + size >= end_ptr; { //Ensure there is enough memory; if there isn't, get more memory until there is enough
        MoreMemory();
    }
 
    addr = bump_ptr; //Return bump pointer (next free address)
    bump_ptr = bump_ptr + size; //Update bump pointer for next call
    return addr;
}

//
// Returns the number of bytes which are currently free and available through the memory manager without reallocation operations
//
func GetFreeMemory() uint64 {
    return end_ptr - bump_ptr; //Free memory are all addresses/bytes between the bump pointer and the end pointer
}
