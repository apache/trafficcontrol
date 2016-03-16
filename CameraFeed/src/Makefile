
SRCS := CommandLine.cpp \
        CurlInit.cpp \
        GetVideo.cpp \
        JPEGFileWriter.cpp \
        main.cpp

EXE := AmcrestIPM-721S_StreamReader

MAKEFLAGS := --no-print-directory
DEPEND_FILE := .dependlist

# TODO: need to install libcurl4-gnutls-dev package
CURL_C_FLAGS := $(shell curl-config --cflags)
CURL_LD_FLAGS := $(shell curl-config --libs)

CC := g++
CFLAGS := --std=c++11 -g -Wall $(CURL_C_FLAGS) $(INC_DIRS)

LD := g++
LDFLAGS :=
LIBS := $(CURL_LD_FLAGS)

OBJS := $(SRCS:%.cpp=%.o)

all: $(EXE)

$(EXE): $(OBJS)
	@echo "Linking $(EXE)"
	@$(LD) $(LDFLAGS) -o $(EXE) $(OBJS) $(LIBS)

%.o:%.cpp
	@echo "Compiling $<"
	@$(CC) $(CFLAGS) -o $@ -c $<

.PHONY: clean
clean:
	@echo "Cleaning $(EXE)"
	@$(RM) $(OBJS) $(EXE) $(DEPEND_FILE) *~

.PHONY: depend
depend:
	@echo "Building dependencies for: $(SRCS)"
	@/bin/cat < /dev/null > $(DEPEND_FILE); \
	 for file in $(SRCS) ; do \
	  srcDepend=$${file}.d; \
	  objFile=`echo $$file | /bin/sed -e 's~\.cpp~.o~'`; \
	  $(CC) $(CFLAGS) -MM -MT $$objFile \
             -MF $$srcDepend $$file; \
	  /bin/cat < $$srcDepend >> $(DEPEND_FILE); \
	  $(RM) $$srcDepend; \
	done

-include $(DEPEND_FILE)
