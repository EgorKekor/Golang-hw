package main

import (
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Stack interface {
	push(val int64)
	pop() (int64, error)
}


func isDigitAscii(val byte) bool{
	return (val > 47 && val < 58)
}

func isOperation(val byte) bool{
	return (val > 39 && val < 44) || val == 47 || val == 45
}


func processStacks(digits Stack, commands Stack) {

}




var ErrNoData = errors.New("no data in stack")



type OperandsStack struct {
	data	[]int64
}

func newOperandsStack(size int64) OperandsStack {
	buf := make([]int64, 0, size)
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
	buf := make([]int64, 0, size)
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


	digitsReg := regexp.MustCompile("[0-9]+")
	digits := digitsReg.FindAllString(expression, -1)

	currentDigit := 0
	for exprReader.Len() > 0 {
		simb, err := exprReader.ReadByte()

		if err != io.EOF {
			if isDigitAscii(simb) {
				value, _ := strconv.ParseInt(digits[currentDigit], 10, 64)
				currentDigit++
				operandsStack.push(value)
				for drop, _ := exprReader.ReadByte(); isDigitAscii(drop); drop, _ = exprReader.ReadByte() {
					println(drop)
				}
				exprReader.UnreadByte()
			} else if isOperation(simb) {
				commandsStack.push(int64 (simb))
			}

			processStacks(&operandsStack, &commandsStack)
		}


	}


	println(expression)
}
