package main

import (
	"log"
	"os"
	"strconv"
	"testing"
)

type TestSort struct {
		flagF bool
		flagU bool
		flagR bool
		flagN bool
		flagK int
}


func TestFileIO(t *testing.T) {
	keys := []TestSort {
		{flagF:false, flagU:false, flagR:false, flagN:false, flagK:0},
		{flagF:false, flagU:false, flagR:false, flagN:false, flagK:0},
		{flagF:true, flagU:false, flagR:false, flagN:false, flagK:0},		// Ignore register
		{flagF:false, flagU:true, flagR:false, flagN:false, flagK:0},		// Unique
		{flagF:true, flagU:true, flagR:false, flagN:false, flagK:0},		// Ignore + Unique
		{flagF:true, flagU:true, flagR:true, flagN:false, flagK:0},	// Ignore + Unique + Reverse
		{flagF:true, flagU:true, flagR:true, flagN:false, flagK:2},		// Ignore + Unique + Reverse + Column
		{flagF:true, flagU:true, flagR:true, flagN:true, flagK:2},			// Ignore + Unique + Reverse + Column + Numeric
	}



	for test, _ := range(keys) {
		sourceFile, errTest := os.OpenFile("sort_cases/test" + strconv.Itoa(test + 1) + ".dat", os.O_RDONLY, 0755)
		if errTest != nil {
			log.Fatal(errTest)
			return
		}


		etalonFile, err := os.OpenFile("sort_cases/test" + strconv.Itoa(test + 1) + "_etalon.dat", os.O_RDONLY, 0755)
		if err != nil {
			log.Fatal(err)
			return
		}


		sortingObj := createSortObject(sourceFile, keys[test].flagK);
		etalonObj := createSortObject(etalonFile, keys[test].flagK);

		simple := func(str1, str2 *string) bool {
			return *str1 < *str2
		}

		reverse := func(str1, str2 *string) bool {
			return *str1 > *str2
		}

		numeric_revers := func(str1, str2 *string) bool {
			left, _ := strconv.Atoi(*str1)
			right, _ := strconv.Atoi(*str2)
			return left < right
		}

		numeric := func(str1, str2 *string) bool {
			left, _ := strconv.Atoi(*str1)
			right, _ := strconv.Atoi(*str2)
			return left < right
		}



		if keys[test].flagF {
			sortingObj.setLowerCaseMode()
		}

		if keys[test].flagU {
			sortingObj.setUniqueMode()
		}

		if keys[test].flagN {
			sortingObj.setNumericMode()
		}

		if keys[test].flagR && keys[test].flagN {
			By(numeric_revers).Sort(sortingObj)
		} else if keys[test].flagN{
			By(numeric).Sort(sortingObj)
		} else if keys[test].flagR {
			By(reverse).Sort(sortingObj)
		} else {
			By(simple).Sort(sortingObj)
		}


		if len(etalonObj.lines) != len(sortingObj.lines) {
			t.Errorf("%d  %d\n", len(etalonObj.lines), len(sortingObj.lines))
			t.Errorf("[" + strconv.Itoa(test + 1) + "]" + "wrong result: got:")
			for _, str := range (sortingObj.lines) {
				t.Errorf("%s\n", str)
			}

			t.Errorf("\nexpected:")
			for _, str := range (etalonObj.lines) {
				t.Errorf("%s\n", str)
			}
			return
		}

		for i, _ := range (etalonObj.lines) {
			if sortingObj.lines[i] != etalonObj.lines[i] {
				t.Errorf("[" + strconv.Itoa(test + 1) + "]" + "wrong result: got:\n%s\n", sortingObj.lines[i])
				t.Errorf("expected:\n%s\n", etalonObj.lines[i])

				t.Errorf("\nRESULT")
				for _, str := range(sortingObj.lines) {
					t.Errorf("%s\n", str)
				}

				t.Errorf("\nRESULT")
				for _, str := range(sortingObj.sortSelect) {
					t.Errorf("%s\n", str)
				}


				t.Errorf("\nETALON")
				for _, str := range(etalonObj.lines) {
					t.Errorf("%s\n", str)
				}
				return
			}
		}

		sourceFile.Close();
		etalonFile.Close();
	}

}
