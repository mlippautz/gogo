# Copyright 2009 The GoGo Authors. All rights reserved.
# Use of this source code is governed by the MIT
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.$(GOARCH)

TARG=libgogo
GOFILES=\
        libgogo.go \
        convert.go \
        string.go \
        memmgr.go \
        stack.go
OFILES=\
        libgogo_$(GOARCH).$O \
        convert_$(GOARCH).$O \
        string_$(GOARCH).$O \
        memmgr_$(GOARCH).$O \

include $(GOROOT)/src/Make.pkg
