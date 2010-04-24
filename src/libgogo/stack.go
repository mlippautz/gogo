// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo stack functions
//

package libgogo

type Stack struct {
    baseAddress uint64; //Where the stack starts
    itemCount uint64; //How many items there currently are on the stack
    capacity uint64; //How many items there can be max. (not to be changed from outside the library)
};

func ToUint64FromUint64Ptr(value *uint64) uint64;

func InitializeStack(uninitializedStack *Stack) {
    uninitializedStack.baseAddress = Alloc(16 * 8); //Allocate 16 items by default
    uninitializedStack.capacity = 16;
    uninitializedStack.itemCount = 0; //Reset item count (to zero)
}

func Push(stack *Stack, value uint64) {
    var newAddress uint64;
    if (stack.capacity == stack.itemCount) { //Grow stack if its capacity doesn't suffice to push another item
        newAddress = Alloc(stack.capacity * 2 * 8); //Double the capacity
        CopyMem(stack.baseAddress, newAddress, stack.capacity * 8); //Copy old stack items
        stack.baseAddress = newAddress; //Set new address
        stack.capacity = stack.capacity * 2; //Update (increase) capacity
    }
    CopyMem(ToUint64FromUint64Ptr(&value), stack.baseAddress + 8 * stack.itemCount, 8); //Push the new value
    stack.itemCount = stack.itemCount + 1; //Update item count
}

func Pop(stack *Stack) uint64 {
    var returnValue uint64 = Peek(stack); //Peek in order to get the last value from the stack
    stack.itemCount = stack.itemCount - 1; //Update (decrease) item count
    return returnValue; //Return value taken from stack
}

func Peek(stack *Stack) uint64 {
    var returnValue uint64;
    if (stack.itemCount == 0) { //Check if there is an item on the stack
        ExitError("Tried to Peek() from an empty stack", 126);
    }
    CopyMem(stack.baseAddress + (stack.itemCount - 1) * 8, ToUint64FromUint64Ptr(&returnValue), 8); //Copy the last value from the stack into the return variable
    return returnValue; //Return peeked value
}

func GetStackItemCount(stack *Stack) uint64 {
    return stack.itemCount;
}

func GetStackCapacity(stack *Stack) uint64 {
    return stack.capacity;
}
