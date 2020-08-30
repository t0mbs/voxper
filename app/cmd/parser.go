package main

import (
  "fmt"
  "strings"
  "encoding/json"
)

type Snyk struct {
  Findings []struct {
    Score string `json:"cvssScore"`
    Severity string `json:"severity"`
    CreationTime string `json:"creationTime"`
    Descr string `json:"description"`
    FixedIn []string `json:"fixedIn"`
    SnykId string `json:"id"`
    Identifiers struct {
      Cve []string `json"CVE"`
      Cwe []string `json"CWE"`
    } `json:"identifiers"`
    PackageName string `json:"packageName"`
    Version string `json:"version"`
    References []string `json:"references"`
    Upgradable bool `json:"isUpgradable"`
    UpgradePath []string `json:"upgradePath"`
    Name string `json:"title"`
  } `json:"vulnerabilities"`
}

type Brakeman struct {
  Findings[] struct {
    Name string `json:"warning_type"`
    Message string `json:"message"`
    File string `json:"file"`
    Line int `json:"line"`
    Ref string `json:"link"`
    Code string `json:"code"`
  } `json:"warnings"`
}

type Parser struct {}

func (p *Parser) Snyk(scanId uint64, js string) []Finding {
  var findings []Finding
  var snyk Snyk
  json.Unmarshal([]byte(js), &snyk)

  for _, f := range snyk.Findings {
    var remediation string
    if f.Upgradable {
      // Strip out empty strings (caused by false boolean in Snyk upgradepath)
      var upgrades []string
      for _, u := range f.UpgradePath {
        if len(u) != 0 {
          upgrades = append(upgrades, u)
        }
      }

      remediation = fmt.Sprintf("Currently running %s %s, this package should be upgraded to one of the followiong versions or higher: %s", f.PackageName, f.Version, strings.Join(upgrades, ","))
    } else {
      remediation = fmt.Sprintf("This package is not upgradeable.")
    }
    ref := "https://snyk.io/vuln/" + f.SnykId
    findings = append(findings, Finding{ScanId: scanId, Name: f.Name, Severity: f.Severity, Descr: f.Descr, Mitigation: remediation, Cve: f.Identifiers.Cve[0], Cwe: f.Identifiers.Cwe[0], Ref: ref})
  }
  return findings
}

func (p *Parser) Brakeman(scanId uint64, js string) []Finding {
  var findings []Finding
  var brakeman Brakeman
  json.Unmarshal([]byte(js), &brakeman)

  for _, f := range brakeman.Findings {
    descr := fmt.Sprintf("%s in file %s l.%d\nRelevant code\n`%s`", f.Message, f.File, f.Line, f.Code)
    mitigation := fmt.Sprintf("Mitigation steps are detailed in the relevant brakeman (documentation page|%s).", f.Ref)
    findings = append(findings, Finding{ScanId: scanId, Name: f.Name, Severity: "medium", Mitigation: mitigation, Descr: descr, Ref: f.Ref})
  }
  return findings
}
