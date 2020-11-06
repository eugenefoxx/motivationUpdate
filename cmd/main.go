package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const (
	// Составлена 1 УП для установщиков для изделий СЛ 1 раз в 3 месяца и чаще
	createNPM = "Создание программы для NPM"
	operator  = "Ельцов Андрей Николаевич"
	// Александров Александр Викторович
	// Аникина Раиса Владимировна
	// Составлена 1 УП для установщиков для контрактных изделий 1 раз в месяц и чаще
	StarLine = "Starline"
	Contruct = "Контрактное пр-во"
	// Составлена 1 УП для машины селективной пайки 1 раз в месяц и чаще
	ProgrammCreateSEHOSEC = "Написание программы для SEHO SEC"
	ProgrammCreateSEHOPRI = "Написание программы для SEHO PRI"
	// Составлена 1 УП для AOI Modus 1 раз в месяц и чаще
	ProgrammCreateAOIModusPRI = "Написание программы для АОИ PRI"
	ProgrammCreateAOIModusSEC = "Написание программы для АОИ SEC"
	// Составлена 1 УП для AOI KohYoung 1 раз в месяц и чаще
	ProgrammCreateAOIKohYoungPRI = "Создание программы для AOI PRI"
	ProgrammCreateAOIKohYoungSEC = "Создание программы для AOI SEC"
	// Выполнены работы по загрузке и\или настройке машины селективной пайки 22 часа в месяц и больше
	SetupSelectivLineSEHOPRI   = "Настройка SEHO PRI"
	SetupSelectivLineSEHOSEC   = "Настройка SEHO SEC"
	SetupSelectivLineSoldering = "Пайка компонентов PRI"
	// Выполнены работы по проверке и\или настройке на АОИ Modus 11 часов и более 50 заготовок

	// Выполнены работы по загрузке и\или настройке трафаретного принтера 2 раза в месяц и более
	SetupTrafaretPrinterPRI = "Настройка принтера Pri"
	SetupTrafaretPrinterSEC = "Настройка принтера Sec"

	// Выполнены работы по загрузке и\или настройке установщиков 22 часа в месяц и больше

	// Выполнены работы по проверке и\или настройке на АОИ KY 11 часов в месяц и больше (более 50 заготовок)

	// Выполнены работы по разбраковке на ревью станции после АОИ KY 11 часов в месяц и больше (более 50 заготовок)

	// Проведение обучений для сотрудников 1 раз в месяц и чаще, бланк ознакомления сотрудников с подписями сдан администратору
	Training = "Проведение обучения"

	// Составление рабочих инструкций и документов 1 раз в 3 месяца и чаще
	WriteInstraction = "Написание инструкции"

	//Выполнена проверка программы установщиков 1 раз месяц и чаще
	VerifyProgrammInstaller = "Проверка программы установщиков"
	VerifyEquipment         = "Проверка комплектации"

	// Выполнена проверка первой платы после сборки установщиками, оставлен комментарий в задаче 1 раз месяц и чаще
	VerifyPCBLine = "Проверка первой платы до оплавления"

	// Выполнена проверка первой спаянной платы после селективной пайки, оставлен комментарий в задаче 1 раз месяц и чаще
	VerifyPCBSolder = "Проверка первой платы после пайки"

	// Выполнена проверка первой платы после оплавления на ICT, сотавлен комментарий в задаче 1 раз месяц и чаще
	ICT = "Внутрисхемное тестирование ICT"
)

func main() {
	reportCsv1 := readfileseeker("./docs/Ежедневный отчёт оператора (с 01.10.2020) (Ответы) (2).csv")

	//reportCsv2 := readseeker(reportCsv1)
	writeChange(reportCsv1)
	reportCsv2 := readfile("report.csv")
	//	for _, each := range reportCsv2 {
	//		fmt.Printf("%s\n", each[18])
	//		//	fmt.Println(each)
	//	}

	responseCheckNPMSL := checkCreatNPMStarLine(reportCsv2)
	fmt.Println("responseCheckNPM StarLine - ", responseCheckNPMSL)
	responseCheckNPMContruct := checkCreatNPMContruct(reportCsv2)
	fmt.Println("responseCheckNPM Contruct - ", responseCheckNPMContruct)
	responseCheckCreateSEHO := checkProgrammCreateSEHO(reportCsv2)
	fmt.Println("responseCheckCreateSEHO - ", responseCheckCreateSEHO)
	responseCheckCreateAOIModus := checkProgrammCreateAOIModus(reportCsv2)
	fmt.Println("responseCheckCreateAOIModus - ", responseCheckCreateAOIModus)
	responseCheckCreateAOIKohYoung := checkProgrammCreateAOIKohYoung(reportCsv2)
	fmt.Println("responseCheckCreateAOIKohYoung - ", responseCheckCreateAOIKohYoung)
	responseSetupSelectivLine := checkSetupSelectivLineSEHO(reportCsv2)
	fmt.Println("responseSetupSelectivLine - ", responseSetupSelectivLine)
	responseSetupTrafaretPrinter := checkSetupTrafaretPrinter(reportCsv2)
	fmt.Println("responseSetupTrafaretPrinter - ", responseSetupTrafaretPrinter)
	responseTraining := checkTraining(reportCsv2)
	fmt.Println("responseTraining - ", responseTraining)
	responseWriteInstraction := checkWriteInstraction(reportCsv2)
	fmt.Println("responseWriteInstraction - ", responseWriteInstraction)

	reponseVerifyProgrammInstaller := checkVerifyProgrammInstaller(reportCsv2)
	fmt.Println("reponseVerifyProgrammInstaller", reponseVerifyProgrammInstaller)
	reponseVerifyEquipment := checkVerifyEquipment(reportCsv2)
	fmt.Println("reponseVerifyEquipment", reponseVerifyEquipment)
	if reponseVerifyProgrammInstaller+reponseVerifyEquipment == 2 {
		fmt.Println("responseVerifyInstaller - OK")
	} else {
		fmt.Println("responseVerifyInstaller - NOK")
	}

	reponseVerifyPCBLine := checkVerifyPCBLine(reportCsv2)
	fmt.Println("reponseVerifyPCBLine - ", reponseVerifyPCBLine)

	responseVerifyPCBSolder := checkVerifyPCBSolder(reportCsv2)
	fmt.Println("responseVerifyPCBSolder - ", responseVerifyPCBSolder)

	responseICT := checkICT(reportCsv2)
	fmt.Println("responseICT - ", responseICT)

	d := time.Date(2020, 11, 3, 12, 30, 0, 0, time.UTC)
	//10.1.2020 19:48:34
	//d := time.Date("2020, 11, 3, 12, 30, 0, 0, time.UTC")
	year, month, day := d.Date()
	fmt.Printf("%v.%v.%v\n", year, month, day)

	now := time.Now()
	fmt.Println(now.Format("01.02.2006"))
	fmt.Println("Today:", now)
	after := now.AddDate(0, -3, 0)
	fmt.Println("Subtract 1 Month:", after)
	fmt.Println("Минус 3 месяца", after.Format("01.02.2006"))
	//	var f time.Time
	//f := "1:30:00 AM"
	//	fmt.Println("fff", f.Format("2:40:00"))
	g := "01.10.2020"
	layout := "01.02.2006"
	tt, _ := time.Parse(layout, g)
	//	fmt.Println("fff", tt.Format("01.02.2006"))
	fmt.Println("fff", tt.Format("01.02.2006"))
	after2 := tt.AddDate(0, 0, -3)
	fmt.Println("Минус 3 мес по отчету - ", after2.Format("01.02.2006"))

	inputTime := "1:30:00 AM"
	ww, _ := time.Parse("3:04:05 AM", inputTime)
	inputTime2 := "1:30:00 AM"
	ww2, _ := time.Parse("3:04:05 AM", inputTime2)
	fmt.Println("ww1 -", ww)
	fmt.Println("ww2 -", ww2)
	fmt.Println("ww -", ww.Format("2:04:00"))
	//	ww3 := ww.Format("2:04:00") + ww2.Format("2:04:00")
	//	start := time.Date(ww)
	//	afterTenSeconds := start.Add(time.Second * 10)
	//	fmt.Printf("start = %v\n", start)
	//	fmt.Printf("start.Add(time.Minute * 10) = %v\n", afterTenMinutes)
	//	newYY := yy.Add(ww * ww2)

	//	fmt.Println("newYY - ", newYY.Format("2:04:00"))
	//	ww4, _ := time.Parse("3:04:05 AM", ww3)
	//	fmt.Println("ww3 -", ww4.Format("2:04:00"))

	//responseVerifyInstaller := checkVerifyInstaller(reportCsv2)
	//fmt.Println("responseVerifyInstaller - ", responseVerifyInstaller)
	/*
		counterNPM := 0
		for _, each := range reportCsv2 {
			if each[18] == operator {
				if each[4] == createNPM {
					counterNPM++
					//	fmt.Println(counterNPM)
					if counterNPM > 7 {
						fmt.Println("OK")
						resalt := "OK"
						return resalt
					}
					if counterNPM < 7 {
						fmt.Println("NOK")
						return
					}
				}
			}
		}
	*/

}

func readfileseeker(name string) [][]string {
	f, err := os.Open(name)
	if err != nil {
		fmt.Println(err)

	}
	defer f.Close()

	//cr := csv.NewReader(f)
	cr, err := readseeker(f)
	if err != nil {
		log.Fatalf("error read", err)
	}
	//	cr.LazyQuotes = true
	//	cr.Comma = '|'
	/*
		CSVdata, err := cr.ReadAll()
		if err != nil {
			//	fmt.Println(err)
			//	os.Exit(1)
			log.Fatal(err)
		}

			for _, each := range CSVdata {
				//fmt.Printf("%s\n", each[0])
				fmt.Println(each)
			}
	*/
	return cr
}

func readseeker(rs io.ReadSeeker) ([][]string, error) {
	row1, err := bufio.NewReader(rs).ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	_, err = rs.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		return nil, err
	}

	lines := csv.NewReader(rs) //.ReadAll()
	lines.Comma = '|'
	//	if err != nil {
	//		return [][]string{}, err
	//	}
	CSVdata, err := lines.ReadAll()
	if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		log.Fatal(err)
	}

	return CSVdata, nil
}

func writeChange(rows [][]string) {
	f, err := os.Create("report.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	//	err = csv.NewWriter(f).WriteAll(rows)

	writer := csv.NewWriter(f)
	writer.Comma = '|'
	writer.WriteAll(rows)

	writer.Flush()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
}

func readfile(name string) [][]string {
	f, err := os.Open(name)
	if err != nil {
		fmt.Println(err)

	}
	defer f.Close()

	cr := csv.NewReader(f)
	//cr, err := readseeker(f)
	//if err != nil {
	//	log.Fatalf("error read", err)
	//}
	cr.LazyQuotes = true
	cr.Comma = '|'

	CSVdata, err := cr.ReadAll()
	if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		log.Fatal(err)
	}

	/*		for _, each := range CSVdata {
				//fmt.Printf("%s\n", each[0])
				fmt.Println(each)
			}
	*/
	return CSVdata
}

func checkCreatNPMStarLine(rows [][]string) string {
	counterNPM := 0
	var result string
	for _, each := range rows {
		if each[18] == operator {
			if each[4] == createNPM && each[7] == StarLine {
				counterNPM++
				//	fmt.Println(counterNPM)
				if counterNPM >= 1 {
					//	fmt.Println("OK")
					result := "OK"
					return result
					//	return resalt
				} else if counterNPM < 1 {
					//	fmt.Println("NOK")
					result := "NOK"
					return result
				}
			}
		}
	}
	return result
}

func checkCreatNPMContruct(rows [][]string) string {
	counterNPM := 0
	var result string
	for _, each := range rows {
		if each[18] == operator {
			if each[4] == createNPM && each[7] == Contruct {
				counterNPM++
				//	fmt.Println(counterNPM)
				if counterNPM >= 1 {
					//	fmt.Println("OK")
					result := "OK"
					return result
					//	return resalt
				} else if counterNPM < 1 {
					//	fmt.Println("NOK")
					result := "NOK"
					return result
				}
			}
		}
	}
	return result
}

func checkProgrammCreateSEHO(rows [][]string) string {
	counterCreateSEHO := 0
	var result string
	for _, each := range rows {
		if each[18] == operator {
			if each[5] == ProgrammCreateSEHOPRI || each[5] == ProgrammCreateSEHOSEC {
				counterCreateSEHO++
				//	fmt.Println(counterNPM)
				if counterCreateSEHO >= 1 {
					//	fmt.Println("OK")
					result := "OK"
					return result
					//	return resalt
				} else if counterCreateSEHO < 1 {
					//	fmt.Println("NOK")
					result := "NOK"
					return result
				}
			}
		}
	}
	return result
}

//CreateAOIModus
func checkProgrammCreateAOIModus(rows [][]string) string {
	counterCreateAOIModus := 0
	var result string
	for _, each := range rows {
		if each[18] == operator {
			if each[5] == ProgrammCreateAOIModusPRI || each[5] == ProgrammCreateAOIModusSEC && each[3] == "THT" {
				counterCreateAOIModus++
				//	fmt.Println(counterNPM)
				if counterCreateAOIModus >= 1 {
					//	fmt.Println("OK")
					result := "OK"
					return result
					//	return resalt
				} else if counterCreateAOIModus < 1 {
					//	fmt.Println("NOK")
					result := "NOK"
					return result
				}
			}
		}
	}
	return result
}

//AOIKohYoung
func checkProgrammCreateAOIKohYoung(rows [][]string) string {
	counterCreateAOIKohYoung := 0
	var result string
	for _, each := range rows {
		if each[18] == operator {
			if each[4] == ProgrammCreateAOIKohYoungPRI || each[4] == ProgrammCreateAOIKohYoungSEC && each[3] == "SMT" {
				counterCreateAOIKohYoung++
				//	fmt.Println(counterNPM)
				if counterCreateAOIKohYoung >= 1 {
					//	fmt.Println("OK")
					result := "OK"
					return result
					//	return resalt
				} else if counterCreateAOIKohYoung < 1 {
					//	fmt.Println("NOK")
					result := "NOK"
					return result
				}
			}
		}
	}
	return result
}

// SetupSelectivLine
// Выполнены работы по загрузке и\или настройке машины селективной пайки 22 часа в месяц и больше
func checkSetupSelectivLineSEHO(rows [][]string) int {
	//	counterCreateAOIKohYoung := 0
	//layout := "3:04:05"
	//	layoutPM := "3:04:05 PM"
	//	t0, _ := time.Parse(layout, "00:00:00")
	//	var sum time.Time
	var result int
	var sumInSeconds int
	comparedDuration7, _ := time.ParseDuration("7h")
	comparedDurations7Sec := int(comparedDuration7.Seconds())
	comparedDuration14, _ := time.ParseDuration("14h")
	comparedDurations14Sec := int(comparedDuration14.Seconds())
	comparedDuration21, _ := time.ParseDuration("21h")
	comparedDurations21Sec := int(comparedDuration21.Seconds())
	comparedDuration28, _ := time.ParseDuration("28h")
	comparedDurations28Sec := int(comparedDuration28.Seconds())
	comparedDuration35, _ := time.ParseDuration("35h")
	comparedDurations35Sec := int(comparedDuration35.Seconds())
	comparedDuration42, _ := time.ParseDuration("42h")
	comparedDurations42Sec := int(comparedDuration42.Seconds())
	comparedDuration49, _ := time.ParseDuration("49h")
	comparedDurations49Sec := int(comparedDuration49.Seconds())
	comparedDuration56, _ := time.ParseDuration("56h")
	comparedDurations56Sec := int(comparedDuration56.Seconds())
	comparedDuration63, _ := time.ParseDuration("63h")
	comparedDurations63Sec := int(comparedDuration63.Seconds())
	comparedDuration70, _ := time.ParseDuration("70h")
	comparedDurations70Sec := int(comparedDuration70.Seconds())
	comparedDuration77, _ := time.ParseDuration("77h")
	comparedDurations77Sec := int(comparedDuration77.Seconds())
	comparedDuration84, _ := time.ParseDuration("84h")
	comparedDurations84Sec := int(comparedDuration84.Seconds())
	comparedDuration91, _ := time.ParseDuration("91h")
	comparedDurations91Sec := int(comparedDuration91.Seconds())
	comparedDuration98, _ := time.ParseDuration("98h")
	comparedDurations98Sec := int(comparedDuration98.Seconds())
	comparedDuration105, _ := time.ParseDuration("105h")
	comparedDurations105Sec := int(comparedDuration105.Seconds())
	comparedDuration112, _ := time.ParseDuration("112h")
	comparedDurations112Sec := int(comparedDuration112.Seconds())
	for _, each := range rows {
		if each[18] == operator {
			if each[17] == SetupSelectivLineSEHOPRI || each[17] == SetupSelectivLineSEHOSEC && each[3] == "THT" {
				//	t, _ := time.Parse(layout, each[16])
				//	sum = sum.Add(t.Sub(t0))
				//	fmt.Println("Time -", t)

				pt := strings.Split(each[16], ":") // parsed time by ":"
				if len(pt) != 3 {
					log.Fatalf("input format mismatch.\nExpecting H:M:S\nHave: %v", pt)
				}

				h, m, s := pt[0], pt[1], pt[2] // hours, minutes, seconds
				formattedDuration := fmt.Sprintf("%sh%sm%ss", h, m, s)

				duration, err := time.ParseDuration(formattedDuration)
				if err != nil {
					log.Fatalf("Failed to parse duration: %v", formattedDuration)
				}
				sumInSeconds += int(duration.Seconds())

			}
		}
	}

	sumInDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", sumInSeconds))
	if sumInSeconds <= comparedDurations7Sec {
		fmt.Printf("Sum [%s] < compareDuration7 [%s]\n", sumInDuration.String(), comparedDuration7.String())
		result := 0
		return result
	} else if sumInSeconds <= comparedDurations14Sec {
		fmt.Printf("Sum [%s] < compareDuration14 [%s]\n", sumInDuration.String(), comparedDuration14.String())
		result := 1
		return result
	} else if sumInSeconds <= comparedDurations21Sec {
		fmt.Printf("Sum [%s] < compareDuration21 [%s]\n", sumInDuration.String(), comparedDuration21.String())
		result := 2
		return result
	} else if sumInSeconds <= comparedDurations28Sec {
		fmt.Printf("Sum [%s] < compareDuration28 [%s]\n", sumInDuration.String(), comparedDuration28.String())
		result := 3
		return result
	} else if sumInSeconds <= comparedDurations35Sec {
		fmt.Printf("Sum [%s] < compareDuration35 [%s]\n", sumInDuration.String(), comparedDuration35.String())
		result := 4
		return result
	} else if sumInSeconds <= comparedDurations42Sec {
		fmt.Printf("Sum [%s] < compareDuration42 [%s]\n", sumInDuration.String(), comparedDuration42.String())
		result := 5
		return result
	} else if sumInSeconds <= comparedDurations49Sec {
		fmt.Printf("Sum [%s] < compareDuration49 [%s]\n", sumInDuration.String(), comparedDuration49.String())
		result := 6
		return result
	} else if sumInSeconds <= comparedDurations56Sec {
		fmt.Printf("Sum [%s] < compareDuration56 [%s]\n", sumInDuration.String(), comparedDuration56.String())
		result := 7
		return result
	} else if sumInSeconds <= comparedDurations63Sec {
		fmt.Printf("Sum [%s] < compareDuration63 [%s]\n", sumInDuration.String(), comparedDuration63.String())
		result := 8
		return result
	} else if sumInSeconds <= comparedDurations70Sec {
		fmt.Printf("Sum [%s] < compareDuration70 [%s]\n", sumInDuration.String(), comparedDuration70.String())
		result := 9
		return result
	} else if sumInSeconds <= comparedDurations77Sec {
		fmt.Printf("Sum [%s] < compareDuration77 [%s]\n", sumInDuration.String(), comparedDuration77.String())
		result := 10
		return result
	} else if sumInSeconds <= comparedDurations84Sec {
		fmt.Printf("Sum [%s] < compareDuration84 [%s]\n", sumInDuration.String(), comparedDuration84.String())
		result := 11
		return result
	} else if sumInSeconds <= comparedDurations91Sec {
		fmt.Printf("Sum [%s] < compareDuration91 [%s]\n", sumInDuration.String(), comparedDuration91.String())
		result := 12
		return result
	} else if sumInSeconds <= comparedDurations98Sec {
		fmt.Printf("Sum [%s] < compareDuration98 [%s]\n", sumInDuration.String(), comparedDuration98.String())
		result := 13
		return result
	} else if sumInSeconds <= comparedDurations105Sec {
		fmt.Printf("Sum [%s] < compareDuration105 [%s]\n", sumInDuration.String(), comparedDuration105.String())
		result := 14
		return result
	} else if sumInSeconds <= comparedDurations112Sec {
		fmt.Printf("Sum [%s] < compareDuration112 [%s]\n", sumInDuration.String(), comparedDuration112.String())
		result := 15
		return result
	}

	/*
		fmt.Println("Sum time - ", sum)

		checkTimeL0 := "07:00:00 AM"
		tcheckTimeL0, _ := time.Parse(layoutPM, checkTimeL0)
		fmt.Printf("tcheckTimeL1: %v\n", tcheckTimeL0)
		checkTimeL1 := "02:00:00 PM"
		tcheckTimeL1, _ := time.Parse(layoutPM, checkTimeL1)
		fmt.Printf("tcheckTimeL1: %v\n", tcheckTimeL1)
		checkTimeL2 := "09:00:00 PM"
		tcheckTimeL2, _ := time.Parse(layoutPM, checkTimeL2)
		fmt.Printf("tcheckTimeL2: %v\n", tcheckTimeL2)
		checkTimeL3 := "02:00:00 PM"
		tcheckTimeL3, _ := time.Parse(layoutPM, checkTimeL3)
		fmt.Printf("tcheckTimeL3: %v\n", tcheckTimeL3)

		level0 := tcheckTimeL0.Before(sum)
		fmt.Println("level0 is:", level0)
		level1 := tcheckTimeL0.After(sum) && tcheckTimeL1.Before(sum)
		fmt.Println("level1 is:", level1)
		level2 := tcheckTimeL1.After(sum) && tcheckTimeL2.Before(sum)
		fmt.Println("level2 is:", level2)
	*/

	/*
		if level == true {
			fmt.Println("OK")
			result := "OK"
			return result
		} else {
			fmt.Println("NOK")
			result := "NOK"
			return result
		}
	*/
	return result
}

// SetupTrafaretPrinterPRI = "Настройка принтера Prim / Sec"
func checkSetupTrafaretPrinter(rows [][]string) string {
	counterSetupTrafaretPrinter := 0
	var result string
	for _, each := range rows {
		if each[18] == operator {
			if each[17] == SetupTrafaretPrinterPRI || each[17] == SetupTrafaretPrinterSEC && each[3] == "SMT" {
				product := each[19]
				fmt.Println("настройка трафаретного принтера, изделия - ", product)
				counterSetupTrafaretPrinter++
				//	fmt.Println(counterNPM)
				if counterSetupTrafaretPrinter >= 2 {
					//	fmt.Println("OK")
					result := "OK"
					return result
					//	return resalt
				} else if counterSetupTrafaretPrinter < 2 {
					//	fmt.Println("NOK")
					result := "NOK"
					return result
				}
			}
		}
	}
	return result
}

// Training
func checkTraining(rows [][]string) string {
	counterTraining := 0
	var result string
	for _, each := range rows {
		if each[18] == operator {
			if each[17] == Training {

				counterTraining++
				//	fmt.Println(counterNPM)
				if counterTraining >= 1 {
					//	fmt.Println("OK")
					result := "OK"
					return result
					//	return resalt
				} else if counterTraining < 1 {
					//	fmt.Println("NOK")
					result := "NOK"
					return result
				}
			}
		}
	}
	return result
}

// WriteInstraction
func checkWriteInstraction(rows [][]string) string {
	counterWriteInstraction := 0
	var result string
	for _, each := range rows {
		if each[18] == operator {
			if each[17] == WriteInstraction {
				//	intr := each[17]
				//	fmt.Println("Инструкция - ", intr)
				counterWriteInstraction++
				//	fmt.Println(counterNPM)
				if counterWriteInstraction >= 1 {
					//	fmt.Println("OK")
					result := "OK"
					return result
					//	return resalt
				} else if counterWriteInstraction < 1 {
					//	fmt.Println("NOK")
					result := "NOK"
					return result
				}
			}
		}
	}
	return result
}

// Выполнена проверка программы установщиков 1 раз месяц и чаще
// VerifyProgrammInstaller = "Проверка программы установщиков"
// VerifyEquipment         = "Проверка комплектации"
func checkVerifyProgrammInstaller(rows [][]string) int {

	counterVerifyProgrammInstaller := 0

	var result int
	for _, each := range rows {
		if each[18] == operator {
			if each[17] == VerifyProgrammInstaller {
				//	intr := each[17]
				//	fmt.Println("Инструкция - ", intr)
				counterVerifyProgrammInstaller++
				//	fmt.Println(counterNPM)
				if counterVerifyProgrammInstaller >= 1 {
					//	fmt.Println("OK")
					result := 1
					fmt.Println("resultVerifyInstaller OK -", result)
					return result

					//	return resalt
				} else if counterVerifyProgrammInstaller < 1 {
					//	fmt.Println("NOK")
					result := 0
					fmt.Println("resultVerifyInstaller NOK -", result)
					return result
				}

			}

		}

	}

	return result
}

func checkVerifyEquipment(rows [][]string) int {
	counterVerifyInstaller := 0

	var result int
	for _, each := range rows {
		if each[18] == operator {

			if each[17] == VerifyEquipment {

				counterVerifyInstaller++
				if counterVerifyInstaller >= 1 {
					result := 1
					fmt.Println("resultVerifyEquipment OK -", result)
					return result
				} else if counterVerifyInstaller < 1 {
					result := 0
					fmt.Println("resultVerifyEquipment NOK -", result)
					return result
				}
			}

		}

	}

	return result
}

// VerifyPCBLine
func checkVerifyPCBLine(rows [][]string) string {
	counterVerifyPCBLine := 0

	var result string
	for _, each := range rows {
		if each[18] == operator {

			if each[17] == VerifyPCBLine && each[3] == "SMT" {

				counterVerifyPCBLine++
				if counterVerifyPCBLine >= 1 {
					result := "OK"
					//	fmt.Println("resultVerifyEquipment OK -", result)
					return result
				} else if counterVerifyPCBLine < 1 {
					result := "NOK"
					//	fmt.Println("resultVerifyEquipment NOK -", result)
					return result
				}
			}

		}

	}

	return result
}

// VerifyPCBSolder = "Проверка первой платы после пайки"
func checkVerifyPCBSolder(rows [][]string) string {
	counterVerifyPCBSolder := 0

	var result string
	for _, each := range rows {
		if each[18] == operator {

			if each[17] == VerifyPCBSolder && each[3] == "THT" {

				counterVerifyPCBSolder++
				if counterVerifyPCBSolder >= 1 {
					result := "OK"
					//	fmt.Println("resultVerifyEquipment OK -", result)
					return result
				} else if counterVerifyPCBSolder < 1 {
					result := "NOK"
					//	fmt.Println("resultVerifyEquipment NOK -", result)
					return result
				}
			}

		}

	}

	return result
}

// ICT = "Внутрисхемное тестирование ICT"
func checkICT(rows [][]string) string {
	counterICT := 0

	var result string
	for _, each := range rows {
		if each[18] == operator {

			if each[17] == ICT {

				counterICT++
				if counterICT >= 1 {
					result := "OK"
					//	fmt.Println("resultVerifyEquipment OK -", result)
					return result
				} else if counterICT < 1 {
					result := "NOK"
					//	fmt.Println("resultVerifyEquipment NOK -", result)
					return result
				}
			}

		}

	}

	return result
}
