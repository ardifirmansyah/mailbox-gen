package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var emailInput string
var licenseInput string

var validLicenseType = map[string]string{
	"E1":    "astrainternational:STANDARDPACK",
	"E3":    "astrainternational:ENTERPRISEPACK",
	"F1":    "astrainternational:DESKLESSPACK",
	"KIOSK": "astrainternational:EXCHANGEDESKLESS",
}

const reMail = `^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`

const targetDomain = "astrainternational.onmicrosoft.com"
const setRegionQuery = "Set-MsolUser -UserPrincipalName %s -UsageLocation ID"
const setLicenseQuery = "Set-MsolUserLicense -UserPrincipalName %s -addLicenses %s"
const setMailboxQuery = "Enable-RemoteMailbox %s -RemoteRoutingAddress %s -PrimarySMTPAddress %s"
const setIcalQuery = "Get-Mailbox %s | Set-CASMailbox -PopUseProtocolDefaults $false"
const setIcalQuery2 = "Get-Mailbox %s | Set-CASMailbox -PopForceICalForCalendarRetrievalOption $true"

// 1. Set Region
// Set-MsolUser -UserPrincipalName ronny.septo@dso.astra.co.id -UsageLocation ID
// 2. Set License
// Set-MsolUserLicense -UserPrincipalName ronny.septo@dso.astra.co.id -addLicenses astrainternational:EXCHANGEDESKLESS
// 3. Remote Mailbox
// Enable-RemoteMailbox ronny.septo@dso.astra.co.id -RemoteRoutingAddress ronny.septo@astrainternational.onmicrosoft.com -PrimarySMTPAddress ronny.septo@dso.astra.co.id
// 4.  Set iCal
// Get-Mailbox ronny.septo@dso.astra.co.id | Set-CASMailbox -PopUseProtocolDefaults $false	Get-Mailbox ronny.septo@astrainternational.onmicrosoft.com | Set-CASMailbox -PopForceICalForCalendarRetrievalOption $true

// E3 = astrainternational:ENTERPRISEPACK
// E1 = astrainternational:STANDARDPACK
// F1 = astrainternational:DESKLESSPACK
// KIOSK = astrainternational:EXCHANGEDESKLESS

func main() {
	r, _ := regexp.Compile(reMail)
	flag.StringVar(&emailInput, "e", "", "email input for generator")
	flag.StringVar(&licenseInput, "l", "", "license type")

	// read input flag
	flag.Parse()

	if emailInput == "" {
		log.Fatal("email not defined")
	}
	if licenseInput == "" {
		log.Fatal("license type not defined")
	}

	if valid := r.MatchString(emailInput); !valid {
		log.Fatal("email format is not correct")
	}

	licenseType, ok := validLicenseType[strings.ToUpper(licenseInput)]
	if !ok {
		log.Fatal("license type not valid")
	}

	vars := strings.Split(emailInput, "@")

	if len(vars) != 2 {
		log.Fatal("email format is not correct")
	}

	nickName := vars[0]
	targetMail := fmt.Sprintf("%s@%s", nickName, targetDomain)

	resultRegion := fmt.Sprintf(setRegionQuery, emailInput)
	resultLicense := fmt.Sprintf(setLicenseQuery, emailInput, licenseType)
	resultMailbox := fmt.Sprintf(setMailboxQuery, emailInput, targetMail, emailInput)
	resultIcal := fmt.Sprintf(setIcalQuery, emailInput)
	resultIcal2 := fmt.Sprintf(setIcalQuery2, targetMail)

	result := fmt.Sprintf("%s\n============\n%s\n============\n%s\n============\n",
		resultRegion,
		resultLicense,
		resultMailbox,
	)

	if strings.ToUpper(licenseInput) == "KIOSK" || strings.ToUpper(licenseInput) == "F1" {
		result = fmt.Sprintf("%s\n============\n%s\n============\n%s\n============\n%s\n============\n",
			resultRegion,
			resultLicense,
			resultMailbox,
			fmt.Sprintf("%s\n%s", resultIcal, resultIcal2),
		)
	}

	fmt.Printf(result)

	outputFile, _ := os.Create("output.txt")
	defer outputFile.Close()

	outputFile.WriteString(result)
}
