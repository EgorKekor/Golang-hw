package main

type TestCalc struct {
	expression		string
	result 			float64
	err				error
	description 	string
}

//func TestCalculator(t *testing.T) {
//	tests := []TestCalc{
//		TestCalc{
//			expression: "1+(2+3)*4/5",
//			result:     5.0000,
//			description:	"Тест на приоритет",
//			err:		nil,
//		},
//		TestCalc{
//			expression: "(1+(2+3)*4/5)*(1+7)",
//			result:     40.0000,
//			description:	"Тест на приоритет",
//			err:		nil,
//		},
//		TestCalc{
//			expression: "1		+ ddd ajs (__dsp'2+3)*ss4/5\n",
//			result:     0,
//			description:	"Тест на ошибки в формуле",
//			err:		ErrInvalidExpression,
//		},
//		TestCalc{
//			expression: "4/0",
//			result:     0,
//			description:	"Тест на ошибки в формуле",
//			err:		ErrInvalidExpression,
//		},
//		TestCalc{
//			expression: "1+2(",
//			result:     0,
//			description:	"Тест на ошибки в формуле",
//			err:		ErrInvalidExpression,
//		},
//		TestCalc{
//			expression: "1)+2(",
//			result:     0,
//			description:	"Тест на ошибки в формуле",
//			err:		ErrInvalidExpression,
//		},
//		TestCalc{
//			expression: "1+2()",
//			result:     0,
//			description:	"Тест на ошибки в формуле",
//			err:		ErrInvalidExpression,
//		},
//		TestCalc{
//			expression: "1+(2)",
//			result:     3.0000,
//			description:	"Тест на ошибки в формуле",
//			err:		nil,
//		},
//	}
//
//	for testNum, item := range(tests) {
//		calculator := NewCalculator(item.expression)
//		calculator.Parse()
//		result, err := calculator.Count()
//
//		if result != item.result || err != item.err {
//			t.Errorf("[%d] wrong result: got %.5f, expected %.5f\n", testNum, result, item.result )
//		}
//	}
//
//}
