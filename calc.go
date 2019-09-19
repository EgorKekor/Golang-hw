package main

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

var ErrNoData = errors.New("no data in stack"

type Stack interface {
	push(val int64)
	pop() (int64, error)
}


type OperandsStack struct {
	data	[]int64
}

func newOperandsStack(size int64) OperandsStack {
	buf := make([]int64, size)
	return OperandsStack{buf}
}

func (stack *OperandsStack) push(val int64)  {
	stack.data = append(stack.data, val)
}

func (stack *OperandsStack) pop() (int64, error)  {
	if len(stack.data) > 0 {
		retVal := stack.data[len(stack.data)-1]
		stack.data = stack.data[0:len(stack.data)-1]
		return retVal, nil
	}
	return 0, ErrNoData
}

// ============================================

type CommandsStack struct {
	data	[]int64
	openNum	int64
}

func newCommandsStack(size int64) CommandsStack {
	buf := make([]int64, size)
	return CommandsStack{buf, 0}
}

func (stack *CommandsStack) push(val int64)  {
	if val == 50 {
		stack.openNum++
	}
	stack.data = append(stack.data, val)
}

func (stack *CommandsStack) pop() (int64, error)  {
	if len(stack.data) > 0 {
		retVal := stack.data[len(stack.data)-1]
		if retVal == 50 {
			stack.openNum--
		}
		stack.data = stack.data[0:len(stack.data)-1]
		return retVal, nil
	}
	return 0, ErrNoData
}



func main()  {
	expression := strings.Join(os.Args[1:], "")
	exprReader := strings.NewReader(expression)

	operandsStack := newOperandsStack(2)
	commandsStack := newCommandsStack(int64 (len(expression)))


	for exprReader.Len() > 0 {
		if num, err := strconv.ParseInt(, 10, 64); err == nil {
			//sObj.sortSelect[i] = strconv.FormatInt(num, 10)
		} else {
			simb, err := exprReader.ReadByte()
			if err != io.EOF {
			}
		}
	}


	println(expression)
}
