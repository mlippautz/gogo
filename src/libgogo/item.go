// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package libgogo

type Item struct {
    mode uint64;
    itemtype *TypeDesc;
    a uint64;
    r uint64;
};

//
// Pseudo constants that specify the descriptor sizes 
//
var ITEM_SIZE uint64 = 32; //4*8 bytes space for an object

//
// Modes for items
//
var MODE_VAR uint64 = 1;
var MODE_CONST uint64 = 2;
var MODE_REG uint64 = 3;

//
// Convert the uint64 value (returned from malloc) to a real item address
//
func Uint64ToItemPtr(adr uint64) *Item;

//
// Creates a new item
//
func NewItem(mode uint64, itemtype *TypeDesc, a uint64, r uint64) *Item {
    var adr uint64;
    var item *Item;
    adr = Alloc(ITEM_SIZE);
    item = Uint64ToItemPtr(adr);
    item.mode = mode;
    item.itemtype = itemtype;
    item.a = a;
    item.r = r;
    return item;
}
