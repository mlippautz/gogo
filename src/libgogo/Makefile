# Copyright 2009 The GoGo Authors. All rights reserved.
# Use of this source code is governed by the MIT
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.$(GOARCH)

TARG=libgogo
GOFILES=\
        libgogo.go \
        convert.go \
        io.go \
        string.go \
        memmgr.go \
        stack.go \
        list.go \
        symbol.go
OFILES=\
        libgogo_$(GOARCH).$O \
        convert_$(GOARCH).$O \
        io_$(GOARCH).$O \
        string_$(GOARCH).$O \
        memmgr_$(GOARCH).$O \
        symbol_$(GOARCH).$O

include $(GOROOT)/src/Make.pkg
