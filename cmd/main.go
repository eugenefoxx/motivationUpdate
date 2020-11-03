package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	// Составлена 1 УП для установщиков для изделий СЛ 1 раз в 3 месяца и чаще
	createNPM = "Создание программы для NPM"
	operator  = "Александров Александр Викторович"
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
	SetupSelectivLineSEHO      = "Настройка SEHO PRI"
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
)

func main() {
	reportCsv1 := readfileseeker("./docs/Ежедневный отчёт оператора (с 01.10.2020) (Ответы) (1).csv")

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
	responseSetupSelectivLine := checkSetupSelectivLine(reportCsv2)
	fmt.Println("responseSetupSelectivLine - ", responseSetupSelectivLine)
	responseSetupTrafaretPrinter := checkSetupTrafaretPrinter(reportCsv2)
	fmt.Println("responseSetupTrafaretPrinter - ", responseSetupTrafaretPrinter)
	responseTraining := checkTraining(reportCsv2)
	fmt.Println("responseTraining - ", responseTraining)
	responseWriteInstraction := checkWriteInstraction(reportCsv2)
	fmt.Println("responseWriteInstraction - ", responseWriteInstraction)
	responseVerifyInstaller := checkVerifyInstaller(reportCsv2)
	fmt.Println("responseVerifyInstaller - ", responseVerifyInstaller)
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
func checkSetupSelectivLine(rows [][]string) string {
	//	counterCreateAOIKohYoung := 0
	layout := "10:30:40 AM"
	//	t, _ := time.ParseDuration("00:00:00 AM")
	//time, _ := time.ParseDuration("22:00:00")
	var result string
	for _, each := range rows {
		if each[18] == operator {
			if each[17] == SetupSelectivLineSEHO && each[3] == "THT" {
				t, _ := time.Parse(layout, each[16])
				fmt.Println("Time -", t)
				//	counterCreateAOIKohYoung++
				//	t++
				//	fmt.Println(counterNPM)
				/*
					if t >= time {
						//	fmt.Println("OK")
						fmt.Println("Время - ", t)
						result := "OK"
						return result
						//	return resalt
					} else if t < time {
						//	fmt.Println("NOK")
						result := "NOK"
						return result
					} */
			}
		}
	}
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
func checkVerifyInstaller(rows [][]string) string {
	counterVerifyInstaller := 0
	counterVerifyProgrammInstaller := 0
	var resultVerifyInstaller string
	var resultVerifyEquipment string
	var result string
	for _, each := range rows {
		if each[18] == operator {
			if each[17] == VerifyProgrammInstaller {
				//	intr := each[17]
				//	fmt.Println("Инструкция - ", intr)
				counterVerifyInstaller++
				//	fmt.Println(counterNPM)
				if counterVerifyInstaller >= 1 {
					//	fmt.Println("OK")
					resultVerifyInstaller := "OK"
					fmt.Println("resultVerifyInstaller OK -", resultVerifyInstaller)
					return resultVerifyInstaller

					//	return resalt
				} else if counterVerifyInstaller < 1 {
					//	fmt.Println("NOK")
					resultVerifyInstaller := "NOK"
					fmt.Println("resultVerifyInstaller NOK -", resultVerifyInstaller)
					return resultVerifyInstaller
				}
			}

		}

	}

	for _, each := range rows {
		if each[18] == operator {
			if each[17] == VerifyEquipment {
				counterVerifyProgrammInstaller++
				if counterVerifyProgrammInstaller >= 1 {
					resultVerifyEquipment := "OK"
					fmt.Println("resultVerifyEquipment OK -", resultVerifyEquipment)
					return resultVerifyEquipment
				} else if counterVerifyProgrammInstaller < 1 {
					resultVerifyEquipment := "NOK"
					fmt.Println("resultVerifyEquipment NOK -", resultVerifyEquipment)
					return resultVerifyEquipment
				} else if counterVerifyProgrammInstaller == 0 {
					resultVerifyEquipment := "NOK"
					fmt.Println("resultVerifyEquipment NOK -", resultVerifyEquipment)
					return resultVerifyEquipment
				}
			}
		}

	}
	if resultVerifyInstaller == "OK" && resultVerifyEquipment == "OK" {
		return result
	} else if resultVerifyInstaller == "OK" && resultVerifyEquipment == "NOK" {
		result := "NOK"
		return result
	}
	return result
}
