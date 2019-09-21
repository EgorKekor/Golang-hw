package main

import (
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var ErrNoData = errors.New("no data in stack")
var ErrNoOperator = errors.New("this is not operator")
var ErrInvalidExpression = errors.New("invalid expression")

const OPEN_BRAKE_PRIOR = 4


type expressionEntity struct {
	value 		int64
	isCommand	bool
}


func isDigitAscii(val byte) bool {
	return (val > 47 && val < 58)
}

func isOperation(val byte) bool {
	return (val > 39 && val < 44) || val == 47 || val == 45
}

func closeSyntax(ind int, str string) error {
	if ind == 0 {
		return ErrInvalidExpression
	}
	if ind > 0 && !isDigitAscii(str[ind - 1]) {
		return ErrInvalidExpression
	} else if ind < len(str) - 1 && !isOperation(str[ind + 1]) {
		return ErrInvalidExpression
	}
	return nil
}

func openSyntax(ind int, str string) error {
	if ind == len(str) - 1 {
		return ErrInvalidExpression
	}
	if ind > 0 && !isOperation(str[ind - 1]) {
		return ErrInvalidExpression
	} else if ind < len(str) - 1 && !isDigitAscii(str[ind + 1]) {
		return ErrInvalidExpression
	}
	return nil
}

func getPrioritet(simb byte) (byte, error) {
	switch simb {
	case '+': return 1, nil;
	case '-': return 1, nil;
	case '*': return 2, nil;
	case '/': return 2, nil;
	case ')': return 3, nil;
	case '(': return OPEN_BRAKE_PRIOR, nil;
	default: return 0, ErrNoOperator;
	}
}

func digitSyntax(ind int, str string) error {
	if ind == 0 || ind == len(str) - 1 {
		return nil
	} else if ind > 0 && !isOperation(str[ind - 1]) && str[ind - 1] != '(' {
		return ErrInvalidExpression
	} else if ind < len(str) - 1 && !isOperation(str[ind + 1]) && str[ind + 1] != ')' {
		return ErrInvalidExpression
	}
	return nil
}

func operationSyntax(ind int, str string) error {
	if ind == 0 || ind == len(str) - 1 {
		return ErrInvalidExpression
	} else if ind > 0 && !isDigitAscii(str[ind - 1]) && str[ind - 1] != ')' {
		return ErrInvalidExpression
	} else if ind < len(str) - 1 && !isDigitAscii(str[ind + 1]) && str[ind + 1] != '(' {
		return ErrInvalidExpression
	}
	return nil
}


//  ==================================================
type expressionDeck struct {
	data	[]expressionEntity
}

func newExpressionDeck(size int64) expressionDeck {
	buf := make([]expressionEntity, 0, size)
	return expressionDeck{buf}
}

func (stack *expressionDeck) push(ent expressionEntity)  {
	stack.data = append(stack.data, ent)
}

func (stack *expressionDeck) pop()(expressionEntity, error)  {
	if len(stack.data) > 0 {
		retVal := stack.data[len(stack.data) - 1]
		stack.data = stack.data[0 : len(stack.data) - 1]
		return retVal, nil
	}
	return expressionEntity{0, false}, ErrNoData
}

func (stack *expressionDeck) popFirst() (expressionEntity, error)  {
	if len(stack.data) > 0 {
		retVal := stack.data[0]
		for i, _ := range(stack.data) {
			if i + 1 < len(stack.data) {
				stack.data[i] = stack.data[i + 1]
			}
		}
		stack.data = stack.data[0 : len(stack.data) - 1]
		return retVal, nil
	}
	return expressionEntity{0, false}, ErrNoData
}


//  ==================================================
type commandsStack struct {
	data       []byte
}

func newCommandsStack(size int64) commandsStack {
	buf := make([]byte, 0, size)
	return commandsStack{buf}
}

func (stack *commandsStack) push(val byte)  {
	stack.data = append(stack.data, val)
}

func (stack *commandsStack) pop() (byte, error)  {
	if len(stack.data) > 0 {
		retVal := stack.data[len(stack.data) - 1]

		stack.data = stack.data[0 : len(stack.data) - 1]
		return retVal, nil
	}
	return 0, ErrNoData
}

func (stack *commandsStack) getTopPrioritet() (byte, error)  {
	if len(stack.data) > 0 {
		retVal := stack.data[len(stack.data) - 1]
		prior, _ := getPrioritet(retVal)
		return prior, nil
	}
	return 0, ErrNoData
}

func (stack *commandsStack) peek() (byte, error)  {
	if len(stack.data) > 0 {
		retVal := stack.data[len(stack.data) - 1]
		return retVal, nil
	}
	return 0, ErrNoData
}


//  ==================================================
type Calculator struct {
	expression	expressionDeck
	commands	commandsStack
	stringExpression	string
	err			error
}


func NewCalculator(stringExpr string) Calculator{
	return Calculator {
		expression: newExpressionDeck(int64(len(stringExpr))),
		commands:   newCommandsStack(int64(len(stringExpr)) / 2),
		stringExpression:	stringExpr,
		err:	nil,
	}
}


func (calc *Calculator) pushAll() {
	for val, err := calc.commands.pop(); err == nil; val, err = calc.commands.pop() {
		if val == '(' {
			break
		}
		if val == ')' {
			continue
		}
		calc.expression.push(expressionEntity{int64(val), true})
	}
}


func (calc *Calculator) addCommand (simb byte) {
	if len(calc.commands.data) == 0 {
		calc.commands.push(simb)
		return
	}

	if simb == ')' {
		calc.pushAll()
		return
	}

	var err error
	err = nil
	var topPrior byte
	var simbPrior byte
	for isTransaction := true; isTransaction && err == nil; {
		if topPrior, err = calc.commands.getTopPrioritet(); err == nil {
			if simbPrior, err = getPrioritet(simb); err == nil {
				if simbPrior <= topPrior && topPrior != OPEN_BRAKE_PRIOR {
					replace, _ := calc.commands.pop()
					calc.expression.push(expressionEntity{int64(replace), true})
				} else {
					calc.commands.push(simb)
					isTransaction = false
				}
			}
		} else {
			calc.commands.push(simb)
		}
	}
}


func (calc *Calculator) addOperand (val int64) {
	calc.expression.push(expressionEntity{int64(val), false})
}


func (calc *Calculator) Count () (float64, error) {
	if calc.err != nil {
		return 0, calc.err
	}

	result := make([]float64, 0, 10)
	for val, err := calc.expression.popFirst(); err == nil; val, err = calc.expression.popFirst() {
		if !val.isCommand {
			result = append(result, float64 (val.value))
		} else {
			switch val.value {
			case '+': result[len(result) - 2] += result[len(result) - 1]
			case '-': result[len(result) - 2] -= result[len(result) - 1]
			case '*': result[len(result) - 2] *= result[len(result) - 1]
			case '/':
				if result[len(result) - 1] == 0 {
					return 0, ErrInvalidExpression
				}
				result[len(result) - 2] /= result[len(result) - 1]
			}
			result = result[0:len(result) - 1]
		}
	}
	return result[0], nil
}


func (calc *Calculator) Parse() {
	openNum := 0
	closeNum := 0
	for i, simb := range(calc.stringExpression) {
		if (closeNum > openNum) {
			calc.err = ErrInvalidExpression
			return
		}
		if simb == '(' {
			if err := openSyntax(i, calc.stringExpression); err != nil {
				calc.err = err
				return
			} else {
				openNum++
			}
		} else if simb == ')' {
			if err := closeSyntax(i, calc.stringExpression); err != nil {
				calc.err = err
				return
			} else {
				closeNum++
			}
		} else if isOperation(byte(simb)) {
			if err := operationSyntax(i, calc.stringExpression); err != nil {
				calc.err = err
				return
			}
		} else if isDigitAscii(byte(simb)) {
			if err := digitSyntax(i, calc.stringExpression); err != nil {
				calc.err = err
				return
			}
		} else {
			calc.err = ErrInvalidExpression
			return
		}
	}
	if openNum != closeNum {
		calc.err = ErrInvalidExpression
		return
	}



	exprReader := strings.NewReader(calc.stringExpression)
	digitsReg := regexp.MustCompile("[0-9]+")
	digits := digitsReg.FindAllString(calc.stringExpression, -1)

	currentDigit := 0
	for exprReader.Len() > 0 {
		simb, err := exprReader.ReadByte()

		if err != io.EOF {
			if isDigitAscii(simb) {
				value, _ := strconv.ParseInt(digits[currentDigit], 10, 64)
				currentDigit++
				calc.addOperand(value)

				drop, _ := exprReader.ReadByte()
				for ; isDigitAscii(drop); drop, _ = exprReader.ReadByte() {}

				if exprReader.Len() > 0 || isOperation(drop) {
					exprReader.UnreadByte()
				}
			} else if isOperation(simb) {
				calc.addCommand(simb)
			}
		}
	}
	calc.pushAll()
}



//func main()  {
//	expression := strings.Join(os.Args[1:], "")
//	calculator := NewCalculator(expression)
//
//	calculator.Parse()
//	val, _ := calculator.Count()
//	fmt.Printf("%.5f", val)
//}
