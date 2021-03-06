package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	//d "github.com/eugenefoxx/starLine/motivationUpdate/gsheets"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
)

// ErrorResponse exported
type ErrorResponse struct {
	Message string
	Err     error
}

// GeneralLogger exported
var GeneralLogger *log.Logger

// ErrorLogger exported
var ErrorLogger *log.Logger

const (
	// Составлена 1 УП для установщиков для изделий СЛ 1 раз в 3 месяца и чаще
	createNPM = "Создание программы для NPM"
	operator  = "Шергин Родион Олегович"
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
	SetupSelectivLineSEHOPRI      = "Настройка SEHO PRI"
	SetupSelectivLineSEHOSEC      = "Настройка SEHO SEC"
	SetupSelectivLineSolderingPRI = "Пайка компонентов PRI"
	SetupSelectivLineSolderingSEC = "Пайка компонентов SEC"

	// Выполнены работы по проверке и\или настройке на АОИ Modus 11 часов и более 50 заготовок
	VerifyAOIPRI = "Проверка на АОИ PRI"
	VerifyAOISEC = "Проверка на АОИ SEC"

	// Выполнены работы по загрузке и\или настройке трафаретного принтера 2 раза в месяц и более
	SetupTrafaretPrinterPRI  = "Настройка принтера Pri"
	SetupTrafaretPrinterSEC  = "Настройка принтера Sec"
	SetupTrafaretPrinterPRIM = "Настройка принтера Prim"

	// Выполнены работы по загрузке и\или настройке установщиков 22 часа в месяц и больше
	SetupNPM         = "Настройка установщиков"
	AssemblyLinePRIM = "Сборка на линии Prim"
	AssemblyLineSEC  = "Сборка на линии Sec"
	//	++++
	ChargingFeederPrim  = "Зарядка питателей Prim"
	ChargingFeederSec   = "Зарядка питателей Sec"
	BuildupOfEquipment  = "Наращивание комплектации"
	PreparingCompCharg  = "Подготовка компонентов к зарядке"
	DischargeFeederPrim = "Разрядка питателей Prim"
	DischargeFeederSec  = "Разрядка питателей Sec"
	VerifyCompToLine    = "Верификация компонентов на линию"

	// Выполнены работы по проверке и\или настройке на АОИ KY 11 часов в месяц и больше (более 50 заготовок)
	VerifyAOIKY    = "Проверка плат на АОИ"
	VerifyAOIKYPRI = "Проверка плат на АОИ Prim"
	VerifyAOIKYSEC = "Проверка плат на АОИ Sec"
	// добавил эти операции только по времени к функ checkVerifyAOIKYTime
	DebugAOI    = "Настойка первой платы на АОИ" // убрать - см выше
	DebugAOIPRI = "Настойка первой платы на АОИ PRI"
	DebugAOISEC = "Настойка первой платы на АОИ SEC" //	DebugAOI    = "Настойка первой платы на АОИ"
	//	DebugAOIPRI = "Настойка первой платы на АОИ PRI"
	//	DebugAOISEC = "Настойка первой платы на АОИ SEC"
	// Выполнены работы по разбраковке на ревью станции после АОИ KY 11 часов в месяц и больше (более 50 заготовок)
	ReviewStationPRI = "ReviewStation pri"
	ReviewStationSEC = "ReviewStation sec"
	ReviewStation    = "ReviewStation"

	// Проведение обучений для сотрудников 1 раз в месяц и чаще, бланк ознакомления сотрудников с подписями сдан администратору
	Training = "Проведение обучения"

	// Составление рабочих инструкций и документов 1 раз в 3 месяца и чаще
	WriteInstraction = "Написание инструкции"

	//Выполнена проверка программы установщиков 1 раз месяц и чаще
	VerifyProgrammInstaller = "Проверка программы установщиков"
	VerifyEquipment         = "Проверка комплектации" //"Проверка комплектации"

	// Выполнена проверка первой платы после сборки установщиками, оставлен комментарий в задаче 1 раз месяц и чаще
	VerifyPCBLine = "Проверка первой платы до оплавления"

	// Выполнена проверка первой спаянной платы после селективной пайки, оставлен комментарий в задаче 1 раз месяц и чаще
	VerifyPCBSolder = "Проверка первой платы после пайки"

	// Выполнена проверка первой платы после оплавления на ICT, сотавлен комментарий в задаче 1 раз месяц и чаще
	ICT = "Внутрисхемное тестирование ICT"

	// Выполнена отладка программы АОИ перед сборкой 1 раз в месяц и чаще.
	// Настойка первой платы на АОИ  Настойка первой платы на АОИ SEC Настойка первой платы на АОИ PRI
	//	DebugAOI    = "Настойка первой платы на АОИ" // убрать - см выше
	//	DebugAOIPRI = "Настойка первой платы на АОИ PRI"
	//	DebugAOISEC = "Настойка первой платы на АОИ SEC"

	// Выполнена отладка программы АОИ перед сборкой 1раз в месяц и чаще
	DebugProgrammAOIPRI = "Отладка программы на AOI PRI"
	DebugProgrammAOISEC = "Отладка программы на AOI SEC"

	// вычесление выполнения нормы
	/*
		DischargeFeederPrim = "Разрядка питателей Prim"
		DischargeFeederSec  = "Разрядка питателей Sec"
	*/
)

type Handler struct {
	*chi.Mux
}

func main() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	//	r := mux.NewRouter()
	//	r.HandleFunc("/motivationCounter", pagemotivationRequest()).Methods("GET")
	//	r.HandleFunc("/motivationCounter", motivationRequest()).Methods("POST")
	//	http.Handle("/", r)
	initLog()
	SRCmain()
	time.Sleep(5 * time.Second)
	r := chi.NewRouter()
	h := &Handler{
		Mux: chi.NewMux(),
	}
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		//	r.HandleFunc("/motivationCounter", pagemotivationRequest)
		//	r.HandleFunc("/motivationCounter", motivationRequest)
		r.Get("/motivationCounter", h.PagemotivationRequest())
		r.Post("/motivationCounter", h.MotivationRequest())

	})

	// определяем среду запуска, если Linux google-chrome-stable, если нет -
	// то для Windows
	browser := isCommmandAvailable("google-chrome-stable")
	if browser == true {
		open.StartWith("http://localhost:3001/motivationCounter", "google-chrome-stable")
	}
	if browser == false {
		open.StartWith("http://localhost:3001/motivationCounter", "chrome.exe")
	}

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets"))
	FileServer(r, "/files", filesDir)

	http.ListenAndServe(":3001", r)
	/*
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

func (h *Handler) PagemotivationRequest() http.HandlerFunc {

	tpl, err := template.New("").ParseFiles(viper.GetString("web.html"))
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		err = tpl.ExecuteTemplate(w, "index.html", nil)

	}

}

func (h *Handler) MotivationRequest() http.HandlerFunc {
	type searchBy struct {
		Tabel string `json:"tabel"`
		Date1 string `json:"date1"`
		Date2 string `json:"date2"`
	}
	tpl, err := template.New("").ParseFiles(viper.GetString("web.html"))
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		search := &searchBy{}
		search.Tabel = r.FormValue("tabel")
		fmt.Println("tabel - ", search.Tabel)
		search.Date1 = r.FormValue("date1")
		fmt.Println("date1 - ", search.Date1)
		search.Date2 = r.FormValue("date2")
		fmt.Println("date2 - ", search.Date2)
		reportCsv1 := readfileseeker(viper.GetString("data.source"))

		//reportCsv2 := readseeker(reportCsv1)
		writeChange(reportCsv1)
		// reportCsv2 := readfile("report.csv")
		reportCsv2 := readfile(viper.GetString("update.updateReport"))

		// reportMotivation, err := os.Create("reportMotivation.csv")
		reportMotivation, err := os.Create(viper.GetString("update.updateReportMotivation"))
		if err != nil {
			log.Println(err)
		}
		defer reportMotivation.Close()

		writer := csv.NewWriter(reportMotivation)
		writer.Write([]string{"Критерии (сотрудник - " + search.Tabel + ")", "Балл", "Дополнительные сведения"})
		writer.Comma = ','
		writer.Flush()

		//	split, err := os.OpenFile("reportMotivation.csv", os.O_APPEND|os.O_WRONLY, 0644)
		//	viper.GetString("update.updateReportMotivation")
		split, err := os.OpenFile(viper.GetString("update.updateReportMotivation"), os.O_APPEND|os.O_WRONLY, 0644)

		if err != nil {
			log.Println(err)
			return
		}
		defer split.Close()
		//	for _, each := range reportCsv2 {
		//		fmt.Printf("%s\n", each[18])
		//		//	fmt.Println(each)
		//	}
		//	date1 := "01.10.2020"
		//	date2 := "31.10.2020"

		responseCheckNPMSL := checkCreatNPMStarLine(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responseCheckNPMSL == 1 {
			responseCheckNPMSL = 3
			result := []string{"Составлена 1 УП для установщиков для изделий СЛ 1 раз в 3 месяца и чаще" + "," + strconv.Itoa(responseCheckNPMSL)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}

			fmt.Println("Составлена 1 УП для установщиков для изделий СЛ 1 раз в 3 месяца и чаще, балл - ", responseCheckNPMSL)
		} else if responseCheckNPMSL != 1 {
			responseCheckNPMSL = 0
			result := []string{"Составлена 1 УП для установщиков для изделий СЛ 1 раз в 3 месяца и чаще" + "," + strconv.Itoa(responseCheckNPMSL)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Составлена 1 УП для установщиков для изделий СЛ 1 раз в 3 месяца и чаще, балл - ", responseCheckNPMSL)
		}
		//	fmt.Println("responseCheckNPM StarLine - ", responseCheckNPMSL)

		responseCheckNPMContruct := checkCreatNPMContruct(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responseCheckNPMContruct == 1 {
			responseCheckNPMContruct = 3
			result := []string{"Составлена 1 УП для установщиков для контрактных изделий 1 раз в месяц и чаще" + "," + strconv.Itoa(responseCheckNPMContruct)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Составлена 1 УП для установщиков для контрактных изделий 1 раз в месяц и чаще, балл - ", responseCheckNPMContruct)
		} else if responseCheckNPMContruct != 1 {
			responseCheckNPMContruct = 0
			result := []string{"Составлена 1 УП для установщиков для контрактных изделий 1 раз в месяц и чаще" + "," + strconv.Itoa(responseCheckNPMContruct)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Составлена 1 УП для установщиков для контрактных изделий 1 раз в месяц и чаще, балл - ", responseCheckNPMContruct)
		}
		//fmt.Println("responseCheckNPM Contruct - ", responseCheckNPMContruct)

		responseCheckCreateSEHO := checkProgrammCreateSEHO(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responseCheckCreateSEHO == 1 {
			responseCheckCreateSEHO = 3
			result := []string{"Cоставлена 1 УП для машины селективной пайки 1 раз в месяц и чаще" + "," + strconv.Itoa(responseCheckCreateSEHO)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Cоставлена 1 УП для машины селективной пайки 1 раз в месяц и чаще, балл -", responseCheckCreateSEHO)
		} else if responseCheckCreateSEHO != 1 {
			responseCheckCreateSEHO = 0
			result := []string{"Cоставлена 1 УП для машины селективной пайки 1 раз в месяц и чаще" + "," + strconv.Itoa(responseCheckCreateSEHO)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Cоставлена 1 УП для машины селективной пайки 1 раз в месяц и чаще, балл -", responseCheckCreateSEHO)
		}
		//fmt.Println("responseCheckCreateSEHO - ", responseCheckCreateSEHO)

		responseCheckCreateAOIModus := checkProgrammCreateAOIModus(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responseCheckCreateAOIModus == 1 {
			responseCheckCreateAOIModus = 3
			result := []string{"Составлена 1 УП для AOI Modus 1 раз в месяц и чаще" + "," + strconv.Itoa(responseCheckCreateAOIModus)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Составлена 1 УП для AOI Modus 1 раз в месяц и чаще, балл - ", responseCheckCreateAOIModus)
		} else if responseCheckCreateAOIModus != 1 {
			responseCheckCreateAOIModus = 0
			result := []string{"Составлена 1 УП для AOI Modus 1 раз в месяц и чаще" + "," + strconv.Itoa(responseCheckCreateAOIModus)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Составлена 1 УП для AOI Modus 1 раз в месяц и чаще, балл - ", responseCheckCreateAOIModus)
		}
		//fmt.Println("responseCheckCreateAOIModus - ", responseCheckCreateAOIModus)

		responseCheckCreateAOIKohYoung := checkProgrammCreateAOIKohYoung(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responseCheckCreateAOIKohYoung == 1 {
			responseCheckCreateAOIKohYoung = 3
			result := []string{"Составлена 1 УП для AOI KohYoung 1 раз в месяц и чаще" + "," + strconv.Itoa(responseCheckCreateAOIKohYoung)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Составлена 1 УП для AOI KohYoung 1 раз в месяц и чаще, балл - ", responseCheckCreateAOIKohYoung)
		} else if responseCheckCreateAOIKohYoung != 1 {
			responseCheckCreateAOIKohYoung = 0
			result := []string{"Составлена 1 УП для AOI KohYoung 1 раз в месяц и чаще" + "," + strconv.Itoa(responseCheckCreateAOIKohYoung)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Составлена 1 УП для AOI KohYoung 1 раз в месяц и чаще, балл - ", responseCheckCreateAOIKohYoung)
		}
		//fmt.Println("responseCheckCreateAOIKohYoung - ", responseCheckCreateAOIKohYoung)

		responseSetupSelectivLine, timeSetupSelectivLine := checkSetupSelectivLine(reportCsv2, search.Tabel, search.Date1, search.Date2)
		sumInDurationSelectivLine, _ := time.ParseDuration(fmt.Sprintf("%ds", timeSetupSelectivLine))
		if responseSetupSelectivLine == 1 {
			responseSetupSelectivLine = 3
			result := []string{"Выполнены работы по загрузке и или настройке машины селективной пайки 22 часа в месяц и больше" + "," + strconv.Itoa(responseSetupSelectivLine) + "," + sumInDurationSelectivLine.String()}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Printf("Выполнены работы по загрузке и или настройке машины селективной пайки 22 часа в месяц и больше, балл - %d, время - [%s]\n", responseSetupSelectivLine, sumInDurationSelectivLine.String())
		} else if responseSetupSelectivLine != 1 {
			responseSetupSelectivLine = 0
			result := []string{"Выполнены работы по загрузке и или настройке машины селективной пайки 22 часа в месяц и больше" + "," + strconv.Itoa(responseSetupSelectivLine) + "," + "время - " + sumInDurationSelectivLine.String()}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Printf("Выполнены работы по загрузке и или настройке машины селективной пайки 22 часа в месяц и больше, балл - %d, время - [%s]\n", responseSetupSelectivLine, sumInDurationSelectivLine.String())
		}

		//	responseSetupSelectivLineSoldering := checkSetupSelectivLineSoldering(reportCsv2, date1, date2)
		//	fmt.Println("responseSetupSelectivLineSoldering - ", responseSetupSelectivLineSoldering)
		//	levelSetupSelectivLine := responseSetupSelectivLineSEHO + responseSetupSelectivLineSoldering
		//	fmt.Println("Выполнены работы по загрузке и или настройке машины селективной пайки 22 часа в месяц и больше, баллов - ", levelSetupSelectivLine)

		//////////////////////////////////////////////////////////////////////////////////////////////////////
		reponseVerifyAOIModusTime, reponseVerifyAOIModusTimeHoure := checkVerifyAOIModusTime(reportCsv2, search.Tabel, search.Date1, search.Date2)
		sumInDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", reponseVerifyAOIModusTimeHoure))
		if reponseVerifyAOIModusTime == 1 {
			reponseVerifyAOIModusTime = 1
			fmt.Printf("Выполнены работы по проверке и или настройке на АОИ Modus 11 часов - OK %d, время - [%s]\n", reponseVerifyAOIModusTime, sumInDuration.String())
		} else if reponseVerifyAOIModusTime == 0 {
			reponseVerifyAOIModusTime = 0
			fmt.Printf("Выполнены работы по проверке и или настройке на АОИ Modus 11 часов - NOK %d, время - [%s]\n", reponseVerifyAOIModusTime, sumInDuration.String())
		}
		reponseVerifyAOIModusPCB, reponseVerifyAOIModusPCBQty := checkVerifyAOIModusPCB(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if reponseVerifyAOIModusPCB == 1 {
			reponseVerifyAOIModusPCB = 1
			fmt.Printf("Выполнены работы по проверке и или настройке более 50 заготовок - OK %d, количество -%d\n", reponseVerifyAOIModusPCB, reponseVerifyAOIModusPCBQty)
		} else if reponseVerifyAOIModusPCB == 0 {
			reponseVerifyAOIModusPCB = 0
			fmt.Println("Выполнены работы по проверке и или настройке более 50 заготовок - NOK %d, количество -%d\n", reponseVerifyAOIModusPCB, reponseVerifyAOIModusPCBQty)
		}
		reponseVerifyAOIModus := reponseVerifyAOIModusTime + reponseVerifyAOIModusPCB
		if reponseVerifyAOIModus == 2 {
			reponseVerifyAOIModus = 3
			result := []string{"Выполнены работы по проверке и или настройке на АОИ Modus 11 часов и более 50 заготовок" + "," + strconv.Itoa(reponseVerifyAOIModus) + "," + "время - " + sumInDuration.String() + " " + "сумма - " + strconv.Itoa(reponseVerifyAOIModusPCBQty)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		} else if reponseVerifyAOIModus != 2 {
			reponseVerifyAOIModus = 0
			result := []string{"Выполнены работы по проверке и или настройк на АОИ Modus 11 часов и более 50 заготовок" + "," + strconv.Itoa(reponseVerifyAOIModus) + "," + "время - " + sumInDuration.String() + " " + "количество - " + strconv.Itoa(reponseVerifyAOIModusPCBQty)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		}
		fmt.Println("Выполнены работы по проверке и или настройке на АОИ Modus 11 часов и более 50 заготовок, балл - ", reponseVerifyAOIModus)

		//////////////////////////////////////////////////////////////////////////////////////////////////////
		reponseReviewStationTime, reponseReviewStationHoure := checkReviewStationTime(reportCsv2, search.Tabel, search.Date1, search.Date2)
		sumInDurationReviewStation, _ := time.ParseDuration(fmt.Sprintf("%ds", reponseReviewStationHoure))
		if reponseReviewStationTime == 1 {
			reponseReviewStationTime = 1

			fmt.Printf("Выполнены работы по разбраковке на ревью станции после АОИ KY 11 часов в месяц и больше - OK %d, время - [%s]\n", reponseReviewStationTime, sumInDurationReviewStation.String())
		} else if reponseReviewStationTime == 0 {
			reponseReviewStationTime = 0

			fmt.Printf("Выполнены работы по разбраковке на ревью станции после АОИ KY 11 часов в месяц и больше - NOK %d, время - [%s]\n", reponseReviewStationTime, sumInDuration.String())
		}
		reponseReviewStationPCB, reponseReviewStationPCBQty := checkReviewStationPCB(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if reponseReviewStationPCB == 1 {
			reponseReviewStationPCB = 1

			fmt.Printf("Выполнены работы по проверке и или настройке более 50 заготовок - OK %d, количество -%d\n", reponseReviewStationPCB, reponseReviewStationPCBQty)
		} else if reponseReviewStationPCB == 0 {
			reponseReviewStationPCB = 0
			fmt.Println("Выполнены работы по проверке и или настройке более 50 заготовок - NOK %d, количество -%d\n", reponseReviewStationPCB, reponseReviewStationPCBQty)
		}
		reponseReviewStation := reponseReviewStationTime + reponseReviewStationPCB
		if reponseReviewStation == 2 {
			reponseReviewStation = 3
			result := []string{"Выполнены работы по разбраковке на ревью станции после АОИ KY 11 часов в месяц и больше и более 50 заготовок" + "," + strconv.Itoa(reponseReviewStation) + "," + "время - " + sumInDurationReviewStation.String() + " " + "количество - " + strconv.Itoa(reponseReviewStationPCBQty)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		} else if reponseReviewStation != 2 {
			reponseReviewStation = 0
			result := []string{"Выполнены работы по разбраковке на ревью станции после АОИ KY 11 часов в месяц и больше и более 50 заготовок" + "," + strconv.Itoa(reponseReviewStation) + "," + "время - " + sumInDurationReviewStation.String() + " " + "количество - " + strconv.Itoa(reponseReviewStationPCBQty)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}

		}
		fmt.Println("Выполнены работы по разбраковке на ревью станции после АОИ KY 11 часов в месяц и больше и более 50 заготовок, балл - ", reponseReviewStation)

		/////////////////////////////////////////////////////////////////////////////////
		reponseSetupNPM, timeSetupNPM := checkSetupNPM(reportCsv2, search.Tabel, search.Date1, search.Date2)
		sumInDurationSetupNPM, _ := time.ParseDuration(fmt.Sprintf("%ds", timeSetupNPM))
		if reponseSetupNPM == 1 {
			reponseSetupNPM = 3
			result := []string{"Выполнены работы по загрузке и или настройке установщиков 22 часа в месяц и больше" + "," + strconv.Itoa(reponseSetupNPM) + "," + "время - " + sumInDurationSetupNPM.String()}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		} else if reponseSetupNPM != 1 {
			reponseSetupNPM = 0
			result := []string{"Выполнены работы по загрузке и или настройке установщиков 22 часа в месяц и больше" + "," + strconv.Itoa(reponseSetupNPM) + "," + "время - " + sumInDurationSetupNPM.String()}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		}
		fmt.Printf("Выполнены работы по загрузке и или настройке установщиков 22 часа в месяц и больше, балл - %d, время - [%s]\n", reponseSetupNPM, sumInDurationSetupNPM.String())
		//////////////////////////////////////////////////////////////////////////////////////////////////////
		reponseVerifyAOIKYTime, reponseVerifyAOIKYTimeHoure := checkVerifyAOIKYTime(reportCsv2, search.Tabel, search.Date1, search.Date2)
		sumInDurationAOIKY, _ := time.ParseDuration(fmt.Sprintf("%ds", reponseVerifyAOIKYTimeHoure))
		if reponseVerifyAOIKYTime == 1 {
			reponseVerifyAOIKYTime = 1
			fmt.Printf("Выполнены работы по проверке и или настройке на АОИ KY 11 часов - OK %d, время - [%s]\n", reponseVerifyAOIKYTime, sumInDurationAOIKY.String())
		} else if reponseVerifyAOIModusTime != 1 {
			reponseVerifyAOIKYTime = 0
			fmt.Printf("Выполнены работы по проверке и или настройке на АОИ KY 11 часов - NOK %d, время - [%s]\n", reponseVerifyAOIKYTime, sumInDurationAOIKY.String())
		}
		reponseVerifyAOIKYPCB, reponseVerifyAOIKYPCBQty := checkVerifyAOIKYPCB(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if reponseVerifyAOIKYPCB == 1 {
			reponseVerifyAOIKYPCB = 1
			fmt.Printf("Выполнены работы по проверке и или настройке более 50 заготовок - OK %d, количество -%d\n", reponseVerifyAOIKYPCB, reponseVerifyAOIKYPCBQty)
		} else if reponseVerifyAOIKYPCB != 1 {
			reponseVerifyAOIKYPCB = 0
			fmt.Println("Выполнены работы по проверке и или настройке более 50 заготовок - NOK %d, количество -%d\n", reponseVerifyAOIKYPCB, reponseVerifyAOIKYPCBQty)
		}
		reponseVerifyAOIKY := reponseVerifyAOIKYTime + reponseVerifyAOIKYPCB
		if reponseVerifyAOIKY == 2 {
			reponseVerifyAOIKY = 3
			result := []string{"Выполнены работы по проверке и или настройке на АОИ KY 11 часов и более 50 заготовок" + "," + strconv.Itoa(reponseVerifyAOIKY) + "," + "время - " + sumInDurationAOIKY.String() + " " + "количество - " + strconv.Itoa(reponseVerifyAOIKYPCBQty)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		} else if reponseVerifyAOIKY != 2 {
			reponseVerifyAOIKY = 0
			result := []string{"Выполнены работы по проверке и или настройке на АОИ KY 11 часов и более 50 заготовок" + "," + strconv.Itoa(reponseVerifyAOIKY) + "," + "время - " + sumInDurationAOIKY.String() + " " + "количество - " + strconv.Itoa(reponseVerifyAOIKYPCBQty)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		}
		fmt.Println("Выполнены работы по проверке и или настройке на АОИ KY 11 часов и более 50 заготовок, балл - ", reponseVerifyAOIKY)

		responseSetupTrafaretPrinter := checkSetupTrafaretPrinter(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responseSetupTrafaretPrinter == 1 {
			responseSetupTrafaretPrinter = 3
			result := []string{"Выполнены работы по загрузке и или настройке трафаретного принтера 2 раза в месяц и более" + "," + strconv.Itoa(responseSetupTrafaretPrinter)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Выполнены работы по загрузке и или настройке трафаретного принтера 2 раза в месяц и более, балл - ", responseSetupTrafaretPrinter)
		} else if responseSetupTrafaretPrinter != 1 {
			responseSetupTrafaretPrinter = 0
			result := []string{"Выполнены работы по загрузке и или настройке трафаретного принтера 2 раза в месяц и более" + "," + strconv.Itoa(responseSetupTrafaretPrinter)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Выполнены работы по загрузке и или настройке трафаретного принтера 2 раза в месяц и более, балл - ", responseSetupTrafaretPrinter)
		}
		//	fmt.Println("responseSetupTrafaretPrinter - ", responseSetupTrafaretPrinter)

		responseTraining := checkTraining(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responseTraining == 1 {
			responseTraining = 3
			result := []string{"Проведение обучений для сотрудников 1 раз в месяц и чаще; бланк ознакомления сотрудников с подписями сдан администратору" + "," + strconv.Itoa(responseTraining)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Проведение обучений для сотрудников 1 раз в месяц и чаще, бланк ознакомления сотрудников с подписями сдан администратору, балл -", responseTraining)
		} else if responseTraining != 1 {
			responseTraining = 0
			result := []string{"Проведение обучений для сотрудников 1 раз в месяц и чаще; бланк ознакомления сотрудников с подписями сдан администратору" + "," + strconv.Itoa(responseTraining)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Проведение обучений для сотрудников 1 раз в месяц и чаще; бланк ознакомления сотрудников с подписями сдан администратору, балл -", responseTraining)
		}
		//	fmt.Println("responseTraining - ", responseTraining)

		// 3 мес инвервал
		responseWriteInstraction := checkWriteInstraction(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responseWriteInstraction == 1 {
			responseWriteInstraction = 3
			result := []string{"Составление рабочих инструкций и документов 1 раз в 3 месяца и чаще" + "," + strconv.Itoa(responseWriteInstraction)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Составление рабочих инструкций и документов 1 раз в 3 месяца и чаще, балл - ", responseWriteInstraction)
		} else if responseWriteInstraction != 1 {
			responseWriteInstraction = 0
			result := []string{"Составление рабочих инструкций и документов 1 раз в 3 месяца и чаще" + "," + strconv.Itoa(responseWriteInstraction)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Составление рабочих инструкций и документов 1 раз в 3 месяца и чаще, балл - ", responseWriteInstraction)
		}
		//	fmt.Println("responseWriteInstraction - ", responseWriteInstraction)

		reponseVerifyProgrammInstaller := checkVerifyProgrammInstaller(reportCsv2, search.Tabel, search.Date1, search.Date2)
		/*	if reponseVerifyProgrammInstaller == 1 {
				reponseVerifyProgrammInstaller = 3
				result := []string{"Выполнена проверка программы установщиков 1 раз месяц и чаще" + "," + strconv.Itoa(reponseVerifyProgrammInstaller)}
				for _, v := range result {
					_, err = fmt.Fprintln(split, v)
					if err != nil {
						split.Close()
						return
					}
				}
				fmt.Println("Выполнена проверка программы установщиков 1 раз месяц и чаще, балл - ", reponseVerifyProgrammInstaller)
			} else if reponseVerifyProgrammInstaller != 1 {
				reponseVerifyProgrammInstaller = 0
				result := []string{"Выполнена проверка программы установщиков 1 раз месяц и чаще" + "," + strconv.Itoa(reponseVerifyProgrammInstaller)}
				for _, v := range result {
					_, err = fmt.Fprintln(split, v)
					if err != nil {
						split.Close()
						return
					}
				}
				fmt.Println("Выполнена проверка программы установщиков 1 раз месяц и чаще, балл - ", reponseVerifyProgrammInstaller)
			}*/
		// fmt.Println("reponseVerifyProgrammInstaller", reponseVerifyProgrammInstaller)
		reponseVerifyEquipment := checkVerifyEquipment(reportCsv2, search.Tabel, search.Date1, search.Date2)
		//	fmt.Println("reponseVerifyEquipment", reponseVerifyEquipment)
		//	if reponseVerifyProgrammInstaller+reponseVerifyEquipment == 2 {
		//		fmt.Println("responseVerifyInstaller - OK")
		//	} else {
		//		fmt.Println("responseVerifyInstaller - NOK")
		//	}
		responseVerifyProgrammEquipment := reponseVerifyProgrammInstaller + reponseVerifyEquipment
		if responseVerifyProgrammEquipment >= 1 {
			responseVerifyProgrammEquipment = 3
			result := []string{"Выполнена проверка программы установщиков 1 раз месяц и чаще / Проверка комплектации 1 раз месяц и чаще" + "," + strconv.Itoa(responseVerifyProgrammEquipment)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		} else if responseVerifyProgrammEquipment == 0 {
			responseVerifyProgrammEquipment = 0
			result := []string{"Выполнена проверка программы установщиков 1 раз месяц и чаще / Проверка комплектации 1 раз месяц и чаще" + "," + strconv.Itoa(responseVerifyProgrammEquipment)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		}
		fmt.Println("Выполнена проверка программы установщиков 1 раз месяц и чаще / Проверка комплектации, балл - ", responseVerifyProgrammEquipment)

		reponseVerifyPCBLine := checkVerifyPCBLine(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if reponseVerifyPCBLine == 1 {
			reponseVerifyPCBLine = 3
			result := []string{"Выполнена проверка первой платы после сборки установщиками; оставлен комментарий в задаче 1 раз месяц и чаще" + "," + strconv.Itoa(reponseVerifyPCBLine)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Выполнена проверка первой платы после сборки установщиками, оставлен комментарий в задаче 1 раз месяц и чаще, балл - ", reponseVerifyPCBLine)
		} else if reponseVerifyPCBLine != 1 {
			reponseVerifyPCBLine = 0
			result := []string{"Выполнена проверка первой платы после сборки установщиками; оставлен комментарий в задаче 1 раз месяц и чаще" + "," + strconv.Itoa(reponseVerifyPCBLine)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Выполнена проверка первой платы после сборки установщиками; оставлен комментарий в задаче 1 раз месяц и чаще, балл - ", reponseVerifyPCBLine)
		}

		//fmt.Println("reponseVerifyPCBLine - ", reponseVerifyPCBLine)

		responseVerifyPCBSolder := checkVerifyPCBSolder(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responseVerifyPCBSolder == 1 {
			responseVerifyPCBSolder = 3
			result := []string{"Выполнена проверка первой спаянной платы после селективной пайки; оставлен комментарий в задаче 1 раз месяц и чаще" + "," + strconv.Itoa(responseVerifyPCBSolder)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Выполнена проверка первой спаянной платы после селективной пайки, оставлен комментарий в задаче 1 раз месяц и чаще, балл - ", responseVerifyPCBSolder)
		} else if responseVerifyPCBSolder != 1 {
			responseVerifyPCBSolder = 0
			result := []string{"Выполнена проверка первой спаянной платы после селективной пайки; оставлен комментарий в задаче 1 раз месяц и чаще" + "," + strconv.Itoa(responseVerifyPCBSolder)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Выполнена проверка первой спаянной платы после селективной пайки, оставлен комментарий в задаче 1 раз месяц и чаще, балл - ", responseVerifyPCBSolder)
		}

		//	fmt.Println("responseVerifyPCBSolder - ", responseVerifyPCBSolder)

		responseICT := checkICT(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responseICT == 1 {
			responseICT = 3
			result := []string{"Выполнена проверка первой платы после оплавления на ICT; сотавлен комментарий в задаче 1 раз месяц и чаще" + "," + strconv.Itoa(responseICT)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Выполнена проверка первой платы после оплавления на ICT, сотавлен комментарий в задаче 1 раз месяц и чаще, балл - ", responseICT)
		} else if responseICT != 1 {
			responseICT = 0
			result := []string{"Выполнена проверка первой платы после оплавления на ICT; сотавлен комментарий в задаче 1 раз месяц и чаще" + "," + strconv.Itoa(responseICT)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
			fmt.Println("Выполнена проверка первой платы после оплавления на ICT, сотавлен комментарий в задаче 1 раз месяц и чаще, балл - ", responseICT)
		}
		//
		// fmt.Println("responseICT - ", responseICT)
		/*
			reponseDebugAOI := checkDebugAOI(reportCsv2, search.Tabel, search.Date1, search.Date2)
			if reponseDebugAOI == 1 {
				reponseDebugAOI = 3
				result := []string{"Выполнена отладка программы АОИ перед сборкой 1 раз в месяц и чаще" + "," + strconv.Itoa(reponseDebugAOI)}
				for _, v := range result {
					_, err = fmt.Fprintln(split, v)
					if err != nil {
						split.Close()
						return
					}
				}
			} else if reponseDebugAOI != 1 {
				reponseDebugAOI = 0
				result := []string{"Выполнена отладка программы АОИ перед сборкой 1 раз в месяц и чаще" + "," + strconv.Itoa(reponseDebugAOI)}
				for _, v := range result {
					_, err = fmt.Fprintln(split, v)
					if err != nil {
						split.Close()
						return
					}
				}
			}
			fmt.Println("Выполнена отладка программы АОИ перед сборкой 1 раз в месяц и чаще, балл - ", reponseDebugAOI)
		*/
		responsedebugProgrammAOI := debugProgrammAOI(reportCsv2, search.Tabel, search.Date1, search.Date2)
		if responsedebugProgrammAOI == 1 {
			responsedebugProgrammAOI = 3
			result := []string{"Выполнена отладка программы АОИ перед сборкой 1раз в месяц и чаще" + "," + strconv.Itoa(responsedebugProgrammAOI)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		} else if responsedebugProgrammAOI != 1 {
			responsedebugProgrammAOI = 0
			result := []string{"Выполнена отладка программы АОИ перед сборкой 1раз в месяц и чаще" + "," + strconv.Itoa(responsedebugProgrammAOI)}
			for _, v := range result {
				_, err = fmt.Fprintln(split, v)
				if err != nil {
					split.Close()
					return
				}
			}
		}
		fmt.Println("Выполнена отладка программы АОИ перед сборкой 1раз в месяц и чаще, балл -", responsedebugProgrammAOI)
		// определяем среду запуска, если Linux - soffice, если нет -
		// то для Windows - scalc.exe
		calc := isCommmandAvailable("soffice")
		if calc == true {
			//	open.StartWith("reportMotivation.csv", "soffice")
			open.StartWith(viper.GetString("update.updateReportMotivation"), "soffice")
		}
		if calc == false {
			//	open.StartWith("reportMotivation.csv", "scalc.exe")
			open.StartWith(viper.GetString("update.updateReportMotivation"), "scalc.exe")
		}

		err = tpl.ExecuteTemplate(w, "index.html", nil)

	}
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
	//lines.Comma = '|'
	lines.Comma = ','
	lines.LazyQuotes = true
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
	//f, err := os.Create("report.csv")
	f, err := os.Create(viper.GetString("update.updateReport"))
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

func checkCreatNPMStarLine(rows [][]string, tabel, date1, date2 string) int {
	counterNPM := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"
	//	var dateFrom2 time.Time
	//	var dateTo2 time.Time
	fmt.Println("Tabel -", tabel)
	dateFrom, _ := time.Parse(layoutDate2, date1)
	fmt.Printf("Дата от -: %T\n", dateFrom.Format(layoutDate))
	fmt.Println("Дата от -: ", dateFrom.Format(layoutDate))
	dateFrom2 := dateFrom.Format(layoutDate)
	fmt.Printf("dateFrom - %T\n", dateFrom2)
	fmt.Println("dateFrom - ", dateFrom2)
	//	dateFrom2, _ := time.Parse(layoutDate, dateFromV)
	//	fmt.Println("dateFrom2", dateFrom2)
	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)
	fmt.Println("Дата до -:", dateTo.Format(layoutDate))
	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)
	//	dateTo3, _ := time.Parse(layoutDate, dateTo2)
	//	dateTo2 := dateTo.Format(layoutDate)
	dateCheckFrom := dateFrom3.AddDate(0, -2, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	fmt.Println("Дата по циклу -", dateFrom.AddDate(0, 0, -1))
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				//	fmt.Println("табель - ", each[2])
				if each[4] == createNPM && each[7] == StarLine {
					counterNPM++
					//	fmt.Println(counterNPM)
				}
			}
		}
	}
	fmt.Println("counterNPM - ", counterNPM)
	if counterNPM >= 1 {
		fmt.Println("OK")
		result := 1
		return result
	} else if counterNPM < 1 {
		fmt.Println("NOK")
		result := 0
		return result
	}
	return result
}

func checkCreatNPMContruct(rows [][]string, tabel, date1, date2 string) int {
	counterNPM := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[4] == createNPM && each[7] == Contruct {
					counterNPM++
					//	fmt.Println(counterNPM)
				}
			}
		}
	}
	fmt.Println("counterNPM - ", counterNPM)
	if counterNPM >= 1 {
		fmt.Println("OK")
		result := 1
		return result
		//	return resalt
	} else if counterNPM < 1 {
		fmt.Println("NOK")
		result := 0
		return result
	}
	return result
}

func checkProgrammCreateSEHO(rows [][]string, tabel, date1, date2 string) int {
	counterCreateSEHO := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {

				if each[5] == ProgrammCreateSEHOPRI || each[5] == ProgrammCreateSEHOSEC {
					counterCreateSEHO++
					//	fmt.Println(counterNPM)
				}
			}
		}
	}
	fmt.Println("counterCreateSEHO - ", counterCreateSEHO)
	if counterCreateSEHO >= 1 {
		fmt.Println("OK")
		result := 1
		return result
		//	return resalt
	} else if counterCreateSEHO < 1 {
		fmt.Println("NOK")
		result := 0
		return result
	}
	return result
}

//CreateAOIModus
func checkProgrammCreateAOIModus(rows [][]string, tabel, date1, date2 string) int {
	counterCreateAOIModus := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {

				if each[5] == ProgrammCreateAOIModusPRI || each[5] == ProgrammCreateAOIModusSEC && each[3] == "THT" {
					counterCreateAOIModus++
					//	fmt.Println(counterNPM)
				}
			}
		}
	}
	fmt.Println("counterCreateAOIModus - ", counterCreateAOIModus)
	if counterCreateAOIModus >= 1 {
		fmt.Println("OK")
		result := 1
		return result
		//	return resalt
	} else if counterCreateAOIModus < 1 {
		fmt.Println("NOK")
		result := 0
		return result
	}
	return result
}

//AOIKohYoung
func checkProgrammCreateAOIKohYoung(rows [][]string, tabel, date1, date2 string) int {
	counterCreateAOIKohYoung := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[4] == ProgrammCreateAOIKohYoungPRI || each[4] == ProgrammCreateAOIKohYoungSEC && each[3] == "SMT" {
					counterCreateAOIKohYoung++
					//	fmt.Println(counterNPM)
				}
			}
		}
	}
	fmt.Println("counterCreateAOIKohYoung - ", counterCreateAOIKohYoung)
	if counterCreateAOIKohYoung >= 1 {
		fmt.Println("OK")
		result := 1
		return result
		//	return resalt
	} else if counterCreateAOIKohYoung < 1 {
		fmt.Println("NOK")
		result := 0
		return result
	}
	return result
}

// SetupSelectivLine
// Выполнены работы по загрузке и\или настройке машины селективной пайки 22 часа в месяц и больше
// Настройка SEHO
func checkSetupSelectivLine(rows [][]string, tabel, date1, date2 string) (int, int) {
	//	counterCreateAOIKohYoung := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int
	var sumInSeconds int
	comparedDuration22, _ := time.ParseDuration("22h")
	comparedDurations22Sec := int(comparedDuration22.Seconds())
	/*	comparedDuration14, _ := time.ParseDuration("14h")
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
		comparedDurations112Sec := int(comparedDuration112.Seconds()) */
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	dateEachF := dateEach.Format(layoutDate)
		//if dateEachF >= dateFromV && dateEachF <= dateToV {
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == SetupSelectivLineSEHOPRI || each[17] == SetupSelectivLineSEHOSEC || each[17] == SetupSelectivLineSolderingPRI || each[17] == SetupSelectivLineSolderingSEC && each[3] == "THT" {
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
	}

	sumInDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", sumInSeconds))
	if sumInSeconds >= comparedDurations22Sec {
		fmt.Printf("Sum [%s] >= compareDuration22 [%s]\n", sumInDuration.String(), comparedDuration22.String())
		result := 1
		return result, sumInSeconds
	} else if sumInSeconds <= comparedDurations22Sec {
		fmt.Printf("Sum [%s] < compareDuration22 [%s]\n", sumInDuration.String(), comparedDuration22.String())
		result := 0
		return result, sumInSeconds
	} /*else if sumInSeconds <= comparedDurations21Sec {
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
	} */

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
	return result, sumInSeconds
}

/*
// Выполнены работы по загрузке и\или настройке машины селективной пайки 22 часа в месяц и больше
// Пайка компонентов
// SetupSelectivLineSoldering
func checkSetupSelectivLineSoldering(rows [][]string, date1, date2 string) int {
	//	counterCreateAOIKohYoung := 0
	//layout := "3:04:05"
	//	layoutPM := "3:04:05 PM"
	//	t0, _ := time.Parse(layout, "00:00:00")
	//	var sum time.Time
	layoutDate := "02.01.2006"
	//	layoutDate2 := "25.10.2006"
	//	date1 := "01.10.2020"
	dateFrom, _ := time.Parse(layoutDate, date1)
	fmt.Println("Дата от -:", dateFrom.Format(layoutDate))
	//	dateFromV := dateFrom.Format(layoutDate)

	dateTo, _ := time.Parse(layoutDate, date2)
	fmt.Println("Дата до -:", dateTo.Format(layoutDate))
	//	dateToV := dateTo.Format(layoutDate)

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
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	dateEachF := dateEach.Format(layoutDate)
		//if dateEachF >= dateFromV && dateEachF <= dateToV {
		if dateEach.After(dateFrom.AddDate(0, 0, -1)) && dateEach.Before(dateTo.AddDate(0, 0, +1)) {
			if each[18] == operator {
				if each[17] == SetupSelectivLineSolderingPRI || each[17] == SetupSelectivLineSolderingSEC && each[3] == "THT" {
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

	return result
}
*/

// Выполнены работы по проверке и\или настройке на АОИ Modus 11 часов и более 50 заготовок
// VerifyAOI
// Выполнены работы по проверке и\или настройке на АОИ Modus 11 часов
func checkVerifyAOIModusTime(rows [][]string, tabel, date1, date2 string) (int, int) {
	//	counterCreateAOIKohYoung := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int

	var sumInSeconds int
	comparedDuration11, _ := time.ParseDuration("11h")
	comparedDurations11Sec := int(comparedDuration11.Seconds())
	/*	comparedDuration7, _ := time.ParseDuration("7h")
		comparedDurations7Sec := int(comparedDuration7.Seconds())
		comparedDuration10_5, _ := time.ParseDuration("10h:30m")
		comparedDurations10_5Sec := int(comparedDuration10_5.Seconds())
		comparedDuration14, _ := time.ParseDuration("14h")
		comparedDurations14Sec := int(comparedDuration14.Seconds())
		comparedDuration17_5, _ := time.ParseDuration("17h:30m")
		comparedDurations17_5Sec := int(comparedDuration17_5.Seconds())
		comparedDuration21, _ := time.ParseDuration("21h")
		comparedDurations21Sec := int(comparedDuration21.Seconds())
		comparedDuration24_5, _ := time.ParseDuration("24h:30m")
		comparedDurations24_5Sec := int(comparedDuration24_5.Seconds())
		comparedDuration28, _ := time.ParseDuration("28h")
		comparedDurations28Sec := int(comparedDuration28.Seconds())
		comparedDuration31_5, _ := time.ParseDuration("31h:30m")
		comparedDurations31_5Sec := int(comparedDuration31_5.Seconds())
		comparedDuration35, _ := time.ParseDuration("35h")
		comparedDurations35Sec := int(comparedDuration35.Seconds())
		comparedDuration38_5, _ := time.ParseDuration("38h:30m")
		comparedDurations38_5Sec := int(comparedDuration38_5.Seconds())
		comparedDuration42, _ := time.ParseDuration("42h")
		comparedDurations42Sec := int(comparedDuration42.Seconds())*/
	//	comparedDuration91, _ := time.ParseDuration("91h")
	//	comparedDurations91Sec := int(comparedDuration91.Seconds())
	//	comparedDuration98, _ := time.ParseDuration("98h")
	//	comparedDurations98Sec := int(comparedDuration98.Seconds())
	//	comparedDuration105, _ := time.ParseDuration("105h")
	//	comparedDurations105Sec := int(comparedDuration105.Seconds())
	//	comparedDuration112, _ := time.ParseDuration("112h")
	//	comparedDurations112Sec := int(comparedDuration112.Seconds())
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	dateEachF := dateEach.Format(layoutDate)
		//if dateEachF >= dateFromV && dateEachF <= dateToV {
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == VerifyAOIPRI || each[17] == VerifyAOISEC && each[3] == "THT" {
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
	}

	sumInDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", sumInSeconds))
	if sumInSeconds >= comparedDurations11Sec {
		fmt.Printf("Sum [%s] >= compareDuration11 [%s]\n", sumInDuration.String(), comparedDuration11.String())
		result := 1
		return result, sumInSeconds
	} else if sumInSeconds <= comparedDurations11Sec {
		fmt.Printf("Sum [%s] =< compareDuration11 [%s]\n", sumInDuration.String(), comparedDuration11.String())
		result := 0
		return result, sumInSeconds
	} /*else if sumInSeconds <= comparedDurations10_5Sec {
		fmt.Printf("Sum [%s] < compareDuration10,5 [%s]\n", sumInDuration.String(), comparedDuration10_5.String())
		result := 2
		return result
	} else if sumInSeconds <= comparedDurations14Sec {
		fmt.Printf("Sum [%s] < compareDuration14 [%s]\n", sumInDuration.String(), comparedDuration14.String())
		result := 3
		return result
	} else if sumInSeconds <= comparedDurations17_5Sec {
		fmt.Printf("Sum [%s] < compareDuration17,5 [%s]\n", sumInDuration.String(), comparedDuration17_5.String())
		result := 4
		return result
	} else if sumInSeconds <= comparedDurations21Sec {
		fmt.Printf("Sum [%s] < compareDuration21 [%s]\n", sumInDuration.String(), comparedDuration21.String())
		result := 5
		return result
	} else if sumInSeconds <= comparedDurations24_5Sec {
		fmt.Printf("Sum [%s] < compareDuration24,5 [%s]\n", sumInDuration.String(), comparedDuration24_5.String())
		result := 6
		return result
	} else if sumInSeconds <= comparedDurations28Sec {
		fmt.Printf("Sum [%s] < compareDuration28 [%s]\n", sumInDuration.String(), comparedDuration28.String())
		result := 7
		return result
	} else if sumInSeconds <= comparedDurations31_5Sec {
		fmt.Printf("Sum [%s] < compareDuration31,5 [%s]\n", sumInDuration.String(), comparedDuration31_5.String())
		result := 8
		return result
	} else if sumInSeconds <= comparedDurations35Sec {
		fmt.Printf("Sum [%s] < compareDuration35 [%s]\n", sumInDuration.String(), comparedDuration35.String())
		result := 9
		return result
	} else if sumInSeconds <= comparedDurations38_5Sec {
		fmt.Printf("Sum [%s] < compareDuration38,5 [%s]\n", sumInDuration.String(), comparedDuration38_5.String())
		result := 10
		return result
	} else if sumInSeconds <= comparedDurations42Sec {
		fmt.Printf("Sum [%s] < compareDuration42 [%s]\n", sumInDuration.String(), comparedDuration42.String())
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
	} */

	return result, sumInSeconds
}

// Выполнены работы по проверке и\или настройке более 50 заготовок
func checkVerifyAOIModusPCB(rows [][]string, tabel, date1, date2 string) (int, int) {
	//	counterCreateAOIKohYoung := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int
	var sumpcbG int
	var sumpcbNG int

	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	dateEachF := dateEach.Format(layoutDate)
		//if dateEachF >= dateFromV && dateEachF <= dateToV {
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == VerifyAOIPRI || each[17] == VerifyAOISEC && each[3] == "THT" {
					//	t, _ := time.Parse(layout, each[16])
					//	sum = sum.Add(t.Sub(t0))
					//	fmt.Println("Time -", t)
					pcbG, _ := strconv.Atoi(each[20])
					sumpcbG += pcbG

					pcbNG, _ := strconv.Atoi(each[21])
					sumpcbNG += pcbNG

				}
			}
		}
	}
	sumpcbAll := sumpcbG + sumpcbNG
	fmt.Println("Сумма плат -", sumpcbAll)

	if sumpcbAll >= 50 {

		result := 1
		return result, sumpcbAll
	} else if sumpcbAll <= 50 {

		result := 0
		return result, sumpcbAll
	}
	return result, sumpcbAll
}

// Выполнены работы по проверке и\или настройке на АОИ KY 11 часов в месяц и больше (более 50 заготовок)
// VerifyAOIKY    = "Проверка плат на АОИ"
// VerifyAOIKYPRI = "Проверка плат на АОИ Prim"
// VerifyAOIKYSEC = "Проверка плат на АОИ Sec"
func checkVerifyAOIKYTime(rows [][]string, tabel, date1, date2 string) (int, int) {
	//	counterCreateAOIKohYoung := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int

	var sumInSeconds int
	comparedDuration11, _ := time.ParseDuration("11h")
	comparedDurations11Sec := int(comparedDuration11.Seconds())

	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	dateEachF := dateEach.Format(layoutDate)
		//if dateEachF >= dateFromV && dateEachF <= dateToV {
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == VerifyAOIKYPRI || each[17] == VerifyAOIKYSEC || each[17] == VerifyAOIKY || each[17] == DebugAOI || each[17] == DebugAOIPRI || each[17] == DebugAOISEC && each[3] == "SMT" {
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
	}

	sumInDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", sumInSeconds))
	if sumInSeconds >= comparedDurations11Sec {
		fmt.Printf("Sum [%s] >= compareDuration11 [%s]\n", sumInDuration.String(), comparedDuration11.String())
		result := 1
		return result, sumInSeconds
	} else if sumInSeconds <= comparedDurations11Sec {
		fmt.Printf("Sum [%s] =< compareDuration11 [%s]\n", sumInDuration.String(), comparedDuration11.String())
		result := 0
		return result, sumInSeconds
	}

	return result, sumInSeconds
}

// Выполнены работы по проверке и\или настройке более 50 заготовок
func checkVerifyAOIKYPCB(rows [][]string, tabel, date1, date2 string) (int, int) {
	//	counterCreateAOIKohYoung := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int
	var sumpcbG int
	var sumpcbNG int

	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	dateEachF := dateEach.Format(layoutDate)
		//if dateEachF >= dateFromV && dateEachF <= dateToV {
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == VerifyAOIKYPRI || each[17] == VerifyAOIKYSEC || each[17] == VerifyAOIKY && each[3] == "SMT" {
					//	t, _ := time.Parse(layout, each[16])
					//	sum = sum.Add(t.Sub(t0))
					//	fmt.Println("Time -", t)
					pcbG, _ := strconv.Atoi(each[20])
					sumpcbG += pcbG

					pcbNG, _ := strconv.Atoi(each[21])
					sumpcbNG += pcbNG

				}
			}
		}
	}
	sumpcbAll := sumpcbG + sumpcbNG
	fmt.Println("Сумма плат -", sumpcbAll)

	if sumpcbAll >= 50 {

		result := 1
		return result, sumpcbAll
	} else if sumpcbAll <= 50 {

		result := 0
		return result, sumpcbAll
	}
	return result, sumpcbAll
}

// Выполнены работы по разбраковке на ревью станции после АОИ KY 11 часов в месяц и больше (более 50 заготовок)
// ReviewStationPRI = "ReviewStation pri"
// ReviewStationSEC = "ReviewStation sec"
func checkReviewStationTime(rows [][]string, tabel, date1, date2 string) (int, int) {
	//	counterCreateAOIKohYoung := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int

	var sumInSeconds int
	comparedDuration11, _ := time.ParseDuration("11h")
	comparedDurations11Sec := int(comparedDuration11.Seconds())

	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	dateEachF := dateEach.Format(layoutDate)
		//if dateEachF >= dateFromV && dateEachF <= dateToV {
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == ReviewStationPRI || each[17] == ReviewStationSEC || each[17] == ReviewStation && each[3] == "SMT" {
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
	}

	sumInDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", sumInSeconds))
	if sumInSeconds >= comparedDurations11Sec {
		fmt.Printf("Sum [%s] >= compareDuration11 [%s]\n", sumInDuration.String(), comparedDuration11.String())
		result := 1
		return result, sumInSeconds
	} else if sumInSeconds <= comparedDurations11Sec {
		fmt.Printf("Sum [%s] =< compareDuration11 [%s]\n", sumInDuration.String(), comparedDuration11.String())
		result := 0
		return result, sumInSeconds
	}

	return result, sumInSeconds
}

// Выполнены работы по проверке и\или настройке более 50 заготовок
func checkReviewStationPCB(rows [][]string, tabel, date1, date2 string) (int, int) {
	//	counterCreateAOIKohYoung := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int
	var sumpcbG int
	var sumpcbNG int

	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	dateEachF := dateEach.Format(layoutDate)
		//if dateEachF >= dateFromV && dateEachF <= dateToV {
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == ReviewStationPRI || each[17] == ReviewStationSEC || each[17] == ReviewStation && each[3] == "SMT" {
					//	t, _ := time.Parse(layout, each[16])
					//	sum = sum.Add(t.Sub(t0))
					//	fmt.Println("Time -", t)
					pcbG, _ := strconv.Atoi(each[20])
					sumpcbG += pcbG

					pcbNG, _ := strconv.Atoi(each[21])
					sumpcbNG += pcbNG

				}
			}
		}
	}
	sumpcbAll := sumpcbG + sumpcbNG
	fmt.Println("Сумма плат -", sumpcbAll)

	if sumpcbAll >= 50 {

		result := 1
		return result, sumpcbAll
	} else if sumpcbAll <= 50 {

		result := 0
		return result, sumpcbAll
	}
	return result, sumpcbAll
}

// Выполнены работы по загрузке и\или настройке установщиков 22 часа в месяц и больше
// SetupNPM = "Настройка установщиков"
func checkSetupNPM(rows [][]string, tabel, date1, date2 string) (int, int) {
	//	counterCreateAOIKohYoung := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int

	var sumInSeconds int
	comparedDuration22, _ := time.ParseDuration("22h")
	comparedDurations22Sec := int(comparedDuration22.Seconds())

	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	dateEachF := dateEach.Format(layoutDate)
		//if dateEachF >= dateFromV && dateEachF <= dateToV {
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel && each[18] != "SMT" && each[18] != "THT" {
				if each[17] == SetupNPM || each[17] == AssemblyLinePRIM || each[17] == AssemblyLineSEC || each[17] == ChargingFeederPrim ||
					each[17] == ChargingFeederSec || each[17] == BuildupOfEquipment || each[17] == PreparingCompCharg ||
					each[17] == DischargeFeederPrim || each[17] == DischargeFeederSec || each[17] == VerifyCompToLine && each[3] == "SMT" {
					//if each[17] == SetupNPM || each[17] == AssemblyLinePRIM || each[17] == AssemblyLineSEC && each[3] == "SMT" {
					//t, _ := time.Parse("15:04:05", each[16])
					//	sum = sum.Add(t.Sub(t0))
					//fmt.Println("Time Check NPM - ; SetupNPM - ", t, each[17], each[19])

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
	}

	sumInDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", sumInSeconds))
	if sumInSeconds >= comparedDurations22Sec {
		fmt.Printf("Sum [%s] >= compareDuration22 [%s]\n", sumInDuration.String(), comparedDuration22.String())
		result := 1
		return result, sumInSeconds
	} else if sumInSeconds <= comparedDurations22Sec {
		fmt.Printf("Sum [%s] =< compareDuration22 [%s]\n", sumInDuration.String(), comparedDuration22.String())
		result := 0
		return result, sumInSeconds
	}

	return result, sumInSeconds
}

// SetupTrafaretPrinterPRI = "Настройка принтера Prim / Sec"
func checkSetupTrafaretPrinter(rows [][]string, tabel, date1, date2 string) int {
	counterSetupTrafaretPrinter := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == SetupTrafaretPrinterPRI || each[17] == SetupTrafaretPrinterSEC || each[17] == SetupTrafaretPrinterPRIM && each[3] == "SMT" {
					counterSetupTrafaretPrinter++
				}
			}
		}
	}
	fmt.Println("counterSetupTrafaretPrinter -", counterSetupTrafaretPrinter)
	if counterSetupTrafaretPrinter >= 2 {
		fmt.Println("OK")
		result := 1
		return result
		//	return resalt
	} else if counterSetupTrafaretPrinter < 2 {
		fmt.Println("NOK")
		result := 0
		return result
	}
	return result
}

// Training
func checkTraining(rows [][]string, tabel, date1, date2 string) int {
	counterTraining := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == Training {
					counterTraining++
				}
			}
		}
	}
	if counterTraining >= 1 {
		fmt.Println("OK")
		result := 1
		return result
		//	return resalt
	} else if counterTraining < 1 {
		fmt.Println("NOK")
		result := 0
		return result
	}
	return result
}

// WriteInstraction
func checkWriteInstraction(rows [][]string, tabel, date1, date2 string) int {
	counterWriteInstraction := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, -2, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		//	fmt.Println("dateFrom", dateEach.After(dateFrom.AddDate(0, -1, -1)))
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {

			if each[2] == tabel {
				if each[17] == WriteInstraction {
					//	intr := each[17]
					//	fmt.Println("Инструкция - ", intr)
					counterWriteInstraction++
					//	fmt.Println(counterNPM)

				}
			}

		}
	}
	fmt.Println("Инструкция - ")
	if counterWriteInstraction >= 1 {
		fmt.Println("OK")
		result := 1
		return result
		//	return resalt
	} else if counterWriteInstraction < 1 {
		fmt.Println("NOK")
		result := 0
		return result
	}
	return result
}

// Выполнена проверка программы установщиков 1 раз месяц и чаще
// VerifyProgrammInstaller = "Проверка программы установщиков"
// VerifyEquipment         = "Проверка комплектации"
func checkVerifyProgrammInstaller(rows [][]string, tabel, date1, date2 string) int {

	counterVerifyProgrammInstaller := 0

	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"
	dateFrom, _ := time.Parse(layoutDate2, date1)
	dateFrom2 := dateFrom.Format(layoutDate)
	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)
	dateTo, _ := time.Parse(layoutDate2, date2)
	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)
	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == VerifyProgrammInstaller {
					//	intr := each[17]
					//	fmt.Println("Инструкция - ", intr)
					fmt.Println("VerifyProgrammInstaller -", each[17])
					counterVerifyProgrammInstaller++
					//	fmt.Println(counterNPM)

				}
			}

		}

	}
	fmt.Println("Summa VerifyProgrammInstaller: ", counterVerifyProgrammInstaller)
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

	return result
}

func checkVerifyEquipment(rows [][]string, tabel, date1, date2 string) int {
	counterVerifyInstaller := 0

	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"
	dateFrom, _ := time.Parse(layoutDate2, date1)
	dateFrom2 := dateFrom.Format(layoutDate)
	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)
	dateTo, _ := time.Parse(layoutDate2, date2)
	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)
	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {

				if each[17] == VerifyEquipment { // && each[3] == "SMT"
					fmt.Println("VerifyEquipment: ", each[17])
					counterVerifyInstaller++
				}

			}

		}
	}

	fmt.Println("Summa resultVerifyEquipment: ", counterVerifyInstaller)
	if counterVerifyInstaller >= 1 {
		result := 1
		fmt.Println("resultVerifyEquipment OK -", result)
		return result
	} else if counterVerifyInstaller == 0 {
		result := 0
		fmt.Println("resultVerifyEquipment NOK -", result)
		return result
	}

	return result
}

// VerifyPCBLine
func checkVerifyPCBLine(rows [][]string, tabel, date1, date2 string) int {
	counterVerifyPCBLine := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {

				if each[17] == VerifyPCBLine && each[3] == "SMT" {

					counterVerifyPCBLine++
				}
			}
		}

	}
	if counterVerifyPCBLine >= 1 {
		result := 1
		fmt.Println("resultVerifyPCBLine OK -", result)
		return result
	} else if counterVerifyPCBLine < 1 {
		result := 0
		fmt.Println("resultVerifyPCBLine NOK -", result)
		return result
	}

	return result
}

// VerifyPCBSolder = "Проверка первой платы после пайки"
func checkVerifyPCBSolder(rows [][]string, tabel, date1, date2 string) int {
	counterVerifyPCBSolder := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {

				if each[17] == VerifyPCBSolder && each[3] == "THT" {

					counterVerifyPCBSolder++
				}
			}
		}
	}
	if counterVerifyPCBSolder >= 1 {
		result := 1
		fmt.Println("resultVerifyPCBSolder OK -", result)
		return result
	} else if counterVerifyPCBSolder < 1 {
		result := 0
		fmt.Println("resultVerifyPCBSolder NOK -", result)
		return result
	}

	return result
}

// ICT = "Внутрисхемное тестирование ICT"
func checkICT(rows [][]string, tabel, date1, date2 string) int {
	counterICT := 0
	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {

				if each[17] == ICT {

					counterICT++
				}
			}
		}
	}
	if counterICT >= 1 {
		result := 1
		fmt.Println("resultICT OK -", result)
		return result
	} else if counterICT < 1 {
		result := 0
		fmt.Println("resultICT NOK -", result)
		return result
	}

	return result
}

// Выполнена отладка программы АОИ перед сборкой 1 раз в месяц и чаще.
// Настойка первой платы на АОИ  Настойка первой платы на АОИ SEC Настойка первой платы на АОИ PRI
func checkDebugAOI(rows [][]string, tabel, date1, date2 string) int {
	counterDebugAOI := 0

	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"

	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {

				if each[17] == DebugAOI || each[17] == DebugAOIPRI || each[17] == DebugAOISEC && each[3] == "SMT" {
					intr := each[19]
					fmt.Println("Изделие debugAOI - ", intr)
					counterDebugAOI++
				}
			}
		}
	}
	if counterDebugAOI >= 1 {
		result := 1
		fmt.Println("resultDebugAOI OK -", result)
		return result
	} else if counterDebugAOI < 1 {
		result := 0
		fmt.Println("resultDebugAOI NOK -", result)
		return result
	}

	return result
}

// Выполнена отладка программы АОИ перед сборкой 1раз в месяц и чаще
func debugProgrammAOI(rows [][]string, tabel, date1, date2 string) int {
	counterDebugProgrammAOI := 0

	layoutDate := "02.01.2006" //"02.01.2006"
	layoutDate2 := "2006-01-02"
	dateFrom, _ := time.Parse(layoutDate2, date1)
	dateFrom2 := dateFrom.Format(layoutDate)
	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)
	dateTo, _ := time.Parse(layoutDate2, date2)
	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)
	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)
	var result int

	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])
		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {
				if each[17] == DebugProgrammAOIPRI || each[17] == DebugProgrammAOISEC && each[3] == "SMT" {
					counterDebugProgrammAOI++
				}
			}
		}
	}

	fmt.Println("Summa counterDebugProgrammAOI: ", counterDebugProgrammAOI)
	if counterDebugProgrammAOI >= 1 {
		result := 1
		fmt.Println("resultDebugProgrammAOI OK -", result)
		return result
	} else if counterDebugProgrammAOI < 1 {
		result := 0
		fmt.Println("resultDebugProgrammAOI NOK -", result)
		return result
	}

	return result

}

/*
func rateDischargeFeeder(rows [][]string, tabel, date1, date2 string) int {
	counterDischargeFeeder := 0

	layoutDate := "02.01.2006"
	layoutDate2 := "2006-01-02"
	dateFrom, _ := time.Parse(layoutDate2, date1)

	dateFrom2 := dateFrom.Format(layoutDate)

	dateFrom3, _ := time.Parse(layoutDate, dateFrom2)

	dateTo, _ := time.Parse(layoutDate2, date2)

	dateTo2 := dateTo.Format(layoutDate)
	dateTo3, _ := time.Parse(layoutDate, dateTo2)

	dateCheckFrom := dateFrom3.AddDate(0, 0, -1).Format(layoutDate)
	dateCheckFrom2, _ := time.Parse(layoutDate, dateCheckFrom)
	dateCheckTo := dateTo3.AddDate(0, 0, +1).Format(layoutDate)
	dateCheckTo2, _ := time.Parse(layoutDate, dateCheckTo)

	var result int
	for _, each := range rows {
		dateEach, _ := time.Parse(layoutDate, each[15])

		if dateEach.After(dateCheckFrom2) && dateEach.Before(dateCheckTo2) {
			if each[2] == tabel {

			}
		}
	}
	return result
}
*/

func isCommmandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command", "-v", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func initConfig() error {
	MyDir := DirExists("/home/eugenearch/Code/github.com/eugenefoxx/starLine/motivationUpdate/configs")
	MyDirWin := DirExists("C:/Users/Евгений/Code/github.com/eugenefoxx/starline/motivationUpdate/configs")
	WinSLDir := DirExists("Z:/1_Планирование производства Победит1/motivationUpdate/configs")
	if MyDir == true {
		viper.AddConfigPath("/home/eugenearch/Code/github.com/eugenefoxx/starLine/motivationUpdate/configs")
		viper.SetConfigName("config")
		return viper.ReadInConfig()
	}
	if MyDirWin == true {
		viper.AddConfigPath("C:/Users/Евгений/Code/github.com/eugenefoxx/starline/motivationUpdate/configs")
		viper.SetConfigName("config")
		return viper.ReadInConfig()
	}
	if WinSLDir == true {
		viper.AddConfigPath("Z:/1_Планирование производства Победит1/motivationUpdate/configs")
		viper.SetConfigName("config")
		return viper.ReadInConfig()
	}

	return viper.ReadInConfig()

}

func DirExists(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}
	return false
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func SRCmain() {
	// Strting the Application
	GeneralLogger.Println("Starting Extracting Language Files from GoogleSheet - downloading csv approach..")

	//content, err := ioutil.ReadFile("/home/eugenearch/Code/github.com/eugenefoxx/starLine/motivationUpdate/docs/url.txt")
	content, err := ioutil.ReadFile(viper.GetString("data.url"))
	if err != nil {
		log.Fatal(err)
	}
	// Convert []byte to string and print to screen
	text := string(content)
	// Delete line break from text
	// text = strings.TrimSuffix(text, "\r\n")
	text = strings.TrimSuffix(text, "\n")
	fmt.Println(text)
	//csvFilePath := "./outputs/gsheet.csv"
	//csvFilePath := "./outputs/source.csv"
	csvFilePath := viper.GetString("data.source")
	errorResponse := Download(
		//	"https://docs.google.com/spreadsheets/d/e/2PACX-1vS3Im-UdiFgf2aqvrTm14rXlExQOHYPKq9TTDXp0qjBU4Tnb4e2C5_kFQ8u1U0gfWuV70Ybe1gUO8fe/pub?gid=558395991&sin
		text,
		csvFilePath,
		5000,
	)
	if errorResponse.Err != nil {
		ErrorLogger.Println(errorResponse.Message, errorResponse.Err)
		os.Exit(1)
	}
	WriteLanguageFiles(csvFilePath)
	GeneralLogger.Println("Completed Execution..")
}

func initLog() {
	//	absPath, err := filepath.Abs("./outputs/log")
	//	if err != nil {
	//		fmt.Println("Error reading given path:", err)
	//	}

	//	generalLog, err := os.OpenFile(absPath+"/general-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//generalLog, err := os.OpenFile(viper.GetString("log.logfile"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	generalLog, err := os.OpenFile(viper.GetString("log.logfile"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//generalLog, err := os.OpenFile("/home/eugenearch/Code/github.com/eugenefoxx/starLine/motivationUpdate/log/general-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	GeneralLogger = log.New(generalLog, "General Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(generalLog, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
}

// ReturnErrorResponse exported
func ReturnErrorResponse(err error, message string) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Err:     err,
	}
}

// Should point to the folder where language json files are kept
// const outputPath = "./outputs/"

// Download exported
func Download(url string, filename string, timeout int64) *ErrorResponse {
	GeneralLogger.Println("Downloading", url, "...")
	client := http.Client{
		Timeout: time.Duration(timeout * int64(time.Second)),
	}
	resp, err := client.Get(url)
	if err != nil {
		ErrorLogger.Println("Cannot download file from the given url", err)
		return ReturnErrorResponse(err, "Cannot download file from the given url")
	}

	if resp.StatusCode != 200 {
		ErrorLogger.Printf("Response from the URL was %d, but expecting 200", resp.StatusCode)
		return ReturnErrorResponse(
			errors.New("Response returned with a status different from 200"),
			"Response returned with a status different from 200",
		)
	}
	if resp.Header["Content-Type"][0] != "text/csv" {
		ErrorLogger.Printf("The file downloaded has content type '%s', expected 'text/csv'.", resp.Header["Content-Type"])
		return ReturnErrorResponse(
			errors.New("Downloaded file didn't contain the expected content-type: 'text/csv'"),
			"Downloaded file didn't contain the expected content-type: 'text/csv'",
		)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorLogger.Println("Cannot read Body of Response", err)
		return ReturnErrorResponse(err, "Cannot read Body of Response")
	}

	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		ErrorLogger.Println("Cannot write to file", err)
		return ReturnErrorResponse(err, "Cannot write to file")
	}

	GeneralLogger.Println("Doc downloaded in ", filename)

	return ReturnErrorResponse(nil, "")
}

// WriteLanguageFiles exported
func WriteLanguageFiles(csvFilePath string) *ErrorResponse {
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		ErrorLogger.Println("Cannot open file:"+csvFilePath, err)
		return ReturnErrorResponse(err, "Cannot open file:"+csvFilePath)

	}

	//csvFileContent, err := csv.NewReader(csvFile).ReadAll()
	_, errCSV := csv.NewReader(csvFile).ReadAll()
	if errCSV != nil {
		return ReturnErrorResponse(err, "Cannot read file:"+csvFilePath)
	}

	return ReturnErrorResponse(nil, "")
}
