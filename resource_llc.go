package main

import (
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
    // "github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"

    "github.com/desertbit/fillpdf"
    "net/url"
    "strings"
    // "errors"
    "github.com/PuerkitoBio/goquery"
    "os"
    "io"
    "encoding/json"
    "net/http"
    // "net/url"
    "io/ioutil"
    "log"
    "bytes"
    "mime/multipart"
    "path/filepath"
)


type Address struct {
    AddressCity    string           `json:"address_city,omitempty"`
    AddressCountry string           `json:"address_country,omitempty"`
    AddressLine1   string            `json:"address_line1,omitempty"`
    AddressLine2   string           `json:"address_line2,omitempty"`
    AddressState   string           `json:"address_state,omitempty"`
    AddressZip     string           `json:"address_zip,omitempty"`
    Company        string           `json:"company,omitempty"`
    DateCreated    string            `json:"date_created,omitempty"`
    DateModified   string            `json:"date_modified,omitempty"`
    Deleted        bool             `json:"deleted,omitempty"`
    Description    string           `json:"description,omitempty"`
    Email          string           `json:"email,omitempty"`
    ID             string            `json:"id,omitempty"`
    Metadata       map[string]string `json:"metadata,omitempty"`
    Name           string           `json:"name,omitempty"`
    Object         string            `json:"object,omitempty"`
    Phone          string           `json:"phone,omitempty"`
}


func resourceLLC() *schema.Resource {
    return &schema.Resource{
        Create: resourceLLCCreate,
        Read:   resourceLLCRead,
        Update: resourceLLCUpdate,
        Delete: resourceLLCDelete,

        Schema: map[string]*schema.Schema{
            "owner_name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
            },
            "name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
            },
        },

        Importer: &schema.ResourceImporter{
      State: resourceLLCImport,
    },

       //  CustomizeDiff: customdiff.All(
       //      customdiff.ValidateChange("name", func (old, new, meta interface{}) error {
       //          // If we are increasing "size" then the new value must be
       //          // a multiple of the old value.
       //          // if new.(int) <= old.(int) {
       //          //     return nil
       //          // }
       //          // if (new.(int) % old.(int)) != 0 {
       //          //     return fmt.Errorf("new size value must be an integer multiple of old value %d", old.(int))
       //          // }
       //          // return nil
       //      }),
       //      // customdiff.ForceNewIfChange("size", func (old, new, meta interface{}) bool {
       //      //     // "size" can only increase in-place, so we must create a new resource
       //      //     // if it is decreased.
       //      //     return new.(int) < old.(int)
       //      // }),
       // ),
    }
}

func resourceLLCImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
    entity_name := d.Id()

    var customURL url.URL
    customURL.Scheme = "https"
    customURL.Host = "businesssearch.sos.ca.gov"
    customURL.Path = "/CBS/SearchResults"
    newQueryValues := customURL.Query()
    newQueryValues.Set("filing", "")
    newQueryValues.Set("SearchType", "LPLLC")
    newQueryValues.Set("SearchCriteria", entity_name)
    newQueryValues.Set("SearchSubType", "Keyword")
    customURL.RawQuery = newQueryValues.Encode()


    response, err := http.Get(customURL.String())
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    // Create a goquery document from the HTTP response
    document, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal("Error loading HTTP response body. ", err)
    }

    // Find all links and process them with the function
    // defined earlier
    entityID := strings.TrimSpace(document.Find(".EntityTable tbody tr td:first-child").Text())
    
    // log.Printf("[WARN] Creating membership: %s", address)
    d.SetId(entityID)

    d.Set("name", strings.TrimSpace(document.Find(".EntityTable tbody tr td:nth-child(4) button").Text()))

    return []*schema.ResourceData{d}, nil
}


func resourceLLCCreate(d *schema.ResourceData, m interface{}) error {
    file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    log.SetOutput(file)


    // entity_name := d.Get("name").(string)

    // var customURL url.URL
    // customURL.Scheme = "https"
    // customURL.Host = "businesssearch.sos.ca.gov"
    // customURL.Path = "/CBS/SearchResults"
    // newQueryValues := customURL.Query()
    // newQueryValues.Set("filing", "")
    // newQueryValues.Set("SearchType", "LPLLC")
    // newQueryValues.Set("SearchCriteria", entity_name)
    // newQueryValues.Set("SearchSubType", "Keyword")
    // customURL.RawQuery = newQueryValues.Encode()


    // response, err := http.Get(customURL.String())
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer response.Body.Close()

    // // Create a goquery document from the HTTP response
    // document, err := goquery.NewDocumentFromReader(response.Body)
    // if err != nil {
    //     log.Fatal("Error loading HTTP response body. ", err)
    // }

    // // Find all links and process them with the function
    // // defined earlier
    // entityID := strings.TrimSpace(document.Find(".EntityTable tbody tr td:first-child").Text())
    
    // // log.Printf("[WARN] Creating membership: %s", address)
    // d.SetId(entityID)
    return resourceLLCRead(d, m)
}

func resourceLLCRead(d *schema.ResourceData, m interface{}) error {
    return nil
}


func resourceLLCUpdate(d *schema.ResourceData, m interface{}) error {
    file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    log.SetOutput(file)


    if d.HasChange("name") {
        d.Partial(true)

        name_parts := strings.Split(d.Get("owner_name").(string), " ")

        old_name, new_name := d.GetChange("name")

        form := fillpdf.Form{
            "1FirstName": name_parts[0],
            "1LastName": name_parts[1],
            "1PhoneNumber": "",
            "2EntityName": old_name.(string),
            "2EntityNumber": d.Id(),
            "1CorporationName": old_name.(string),
            "3bNewCorporationName": new_name.(string),
            "2CommentsLine1": "$ terraform -v",
            "2CommentsLine2": "Terraform v0.12.19",
            "2CommentsLine3": "+ provider.terracorp (unversioned)",

            "3ReturneeName": "Kevin Kwok",
            "3ReturneeCompanyName": new_name.(string),
            "3ReturneeAddress": "683 Brannan St, Unit 204",
            "3ReturneeCityStateZip": "San Francisco, CA 94107",
        }

        // Fill the form PDF with our values.
        err = fillpdf.Fill(form, "name-change.pdf", "name-change-filled.pdf", true)

        if err != nil {
            return err
        }

        fromAddy := Address {
            Name: "Kevin Kwok",
            AddressLine1: "683 Brannan St",
            AddressLine2: "Unit 204",
            AddressState: "CA",
            AddressCity: "San Francisco",
            AddressZip: "94107",
        }

        toAddy := Address {
            Name: "Kevin Kwok",
            AddressLine1: "228 CLIPPER STREET",
            // AddressLine2: "Unit 204",
            AddressState: "CA",
            AddressCity: "San Francisco",
            AddressZip: "94114",
        }

        toAddyJSON, err := json.Marshal(toAddy)
        if err != nil {
            log.Print(err)
        }
        log.Print(string(toAddyJSON))

        fromAddyJSON, err := json.Marshal(fromAddy)
        if err != nil {
            log.Print(err)
        }
        log.Print(string(fromAddyJSON))


        filename := "name-change-filled.pdf"
        url := "https://api.lob.com/v1/letters"

        file, err = os.Open(filename)

        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()


        body := &bytes.Buffer{}
        writer := multipart.NewWriter(body)

        writer.WriteField("description", "Demo Letter")
        writer.WriteField("color", "false")
        writer.WriteField("double_sided", "false")
        writer.WriteField("address_placement", "insert_blank_page")
        writer.WriteField("to", string(toAddyJSON))
        writer.WriteField("from", string(fromAddyJSON))

        part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))

        if err != nil {
            log.Fatal(err)
        }

        io.Copy(part, file)
        writer.Close()
        request, err := http.NewRequest("POST", url, body)

        if err != nil {
            log.Fatal(err)
        }

        request.SetBasicAuth("live_71b0e2f396ec63f8cc2fc927f7de6d15a2e", "")

        request.Header.Add("Content-Type", writer.FormDataContentType())
        client := &http.Client{}

        response, err := client.Do(request)

        if err != nil {
            log.Fatal(err)
        }
        defer response.Body.Close()

        content, err := ioutil.ReadAll(response.Body)

        log.Printf("%s", content)

        // return errors.New("what")
        d.Partial(false)

    }
    return resourceLLCRead(d, m)
}

func resourceLLCDelete(d *schema.ResourceData, m interface{}) error {
    log.Printf("Creating membership")
    return nil
}