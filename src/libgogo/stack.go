// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo stack functions
//

package libgogo

//
// Stack data structure
//
type Stack struct {
    baseAddress uint64; //Where the stack starts
    itemCount uint64; //How many items there currently are on the stack
    capacity uint64; //How many items there can be max. (not to be changed from outside the library)
};

//
// Initializes the given stack
//
func InitializeStack(uninitializedStack *Stack) {
    uninitializedStack.baseAddress = Alloc(16 * 8); //Allocate 16 items by default
    uninitializedStack.capacity = 16;
    uninitializedStack.itemCount = 0; //Reset item count (to zero)
}

//
// Pushes an item onto the given stack, increasing its capacity if required
//
func Push(stack *Stack, value uint64) {
    var newAddress uint64;
    if (stack.capacity == stack.itemCount) { //Grow stack if its capacity doesn't suffice to push another item
        newAddress = Alloc(stack.capacity * 2 * 8); //Double the capacity
        CopyMem(stack.baseAddress, newAddress, stack.capacity * 8); //Copy old stack items
        stack.baseAddress = newAddress; //Set new address as base address
        stack.capacity = stack.capacity * 2; //Update (increase) capacity
    }
    CopyMem(ToUint64FromUint64Ptr(&value), stack.baseAddress + 8 * stack.itemCount, 8); //Push the new value by copying its value into the memory of the corresponding stack item
    stack.itemCount = stack.itemCount + 1; //Update item count
}

//
// Takes and returns an item from the given stack
// Note that this function call fails if there are no items on the stack, e.g. the stack is empty
//
func Pop(stack *Stack) uint64 {
    var returnValue uint64 = Peek(stack); //Peek in order to get the last value from the stack
    stack.itemCount = stack.itemCount - 1; //Update (decrease) item count
    return returnValue; //Return value taken from stack
}

//
// Returns the topmost item of a given stack without removing/taking it
// Note that this function call fails if there are no items on the stack, e.g. the stack is empty
//
func Peek(stack *Stack) uint64 {
    var returnValue uint64;
    if (stack.itemCount == 0) { //Check if there is an item on the stack
        ExitError("Tried to Peek() from an empty stack", 126);
    }
    CopyMem(stack.baseAddress + (stack.itemCount - 1) * 8, ToUint64FromUint64Ptr(&returnValue), 8); //Copy the last value from the stack (e.g. its corresponding memory) into the return variable
    return returnValue; //Return peeked value
}

//
// Returns the number of items on the given stack
//
func GetStackItemCount(stack *Stack) uint64 {
    return stack.itemCount;
}

//
// Returns the current capacity of the given stack
//
func GetStackCapacity(stack *Stack) uint64 {
    return stack.capacity;
}
